package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/google/uuid"
	"gopkg.in/yaml.v3"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/harrison-minibucks/github-api-demo/model"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/github"
)

const userSession = "user"

func init() {
	data, err := os.ReadFile("secret.yml")
	if err != nil {
		panic(err)
	}
	config := &model.Config{}
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		panic(err)
	}
	goth.UseProviders(
		github.New(config.GithubApp.ClientId, config.GithubApp.ClientSecret, "http://localhost:3000/auth/github/callback"),
	)

	// TODO: Convert sessions into database sessions
	gothic.Store = sessions.NewCookieStore([]byte("SESSION_SECRET"))

	gothic.SetState = func(req *http.Request) string {
		return uuid.New().String()
	}
}

func callbackHandler(res http.ResponseWriter, req *http.Request) {
	user, err := gothic.CompleteUserAuth(res, req)
	if err != nil {
		fmt.Fprintln(res, err)
		return
	}
	session, _ := gothic.Store.Get(req, "github_login_session")
	session.Values[userSession] = user
	session.Save(req, res)

	// user is now authenticated with github and available for use
	userJson, _ := json.Marshal(user)
	fmt.Fprintf(res, "User: %s\n", string(userJson))
}

func manualCallbackHandler(res http.ResponseWriter, req *http.Request) {
	var user goth.User
	err := json.NewDecoder(req.Body).Decode(&user)
	if err != nil {
		fmt.Fprintln(res, err)
		return
	}
	session, _ := gothic.Store.Get(req, "github_login_session")
	session.Values[userSession] = user
	session.Save(req, res)

	// user is now authenticated with github and available for use
	fmt.Fprintf(res, "User: %s\n", user)
}

func tokenHandler(res http.ResponseWriter, req *http.Request) {
	token := req.Header.Get("Authorization")
	ghUser := &model.GitHubUser{}
	err := callGitHubAPI(token, ghUser)
	if err != nil {
		fmt.Fprintln(res, err)
		return
	}
	userJson, _ := json.Marshal(ghUser)
	fmt.Fprintf(res, "User: %s\n", string(userJson))
}

func callGitHubAPI(token string, ghUser *model.GitHubUser) error {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://api.github.com/user", nil)
	if err != nil {
		fmt.Println("Error creating HTTP request:", err)
		return err
	}
	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Add("Accept", "application/vnd.github+json")
	req.Header.Add("X-GitHub-Api-Version", "2022-11-28")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error making HTTP request:", err)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed status code: %d", resp.StatusCode)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading HTTP response:", err)
		return err
	}
	err = json.Unmarshal(body, ghUser)
	if err != nil {
		fmt.Println("Error decoding user:", err)
		return err
	}
	fmt.Println(ghUser)

	return nil
}

func logoutHandler(res http.ResponseWriter, req *http.Request) {
	session, err := gothic.Store.Get(req, "github_login_session")
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
	// Check if the user is in the session
	if _, ok := session.Values[userSession]; !ok {
		http.Error(res, "User not logged in", http.StatusBadRequest)
		return
	}
	session.Values = nil
	session.Save(req, res)
	err = gothic.Logout(res, req)
	if err != nil {
		fmt.Fprintln(res, err)
		return
	}
	fmt.Fprintln(res, "<p>Logout Successfully</p><a href=\"/auth/github\">Login to GitHub</a>")
}

func avatarHandler(res http.ResponseWriter, req *http.Request) {
	session, err := gothic.Store.Get(req, "github_login_session")
	if err != nil {
		fmt.Fprintln(res, err)
		return
	}
	if session.Values[userSession] != nil {
		user := session.Values[userSession].(goth.User)
		fmt.Fprintf(res, "<p>Avatar: </p><img src=\"%s\"/>", user.AvatarURL)
	} else {
		fmt.Fprintln(res, "<a href=\"/auth/github\">Login to GitHub First</a>")
	}
}

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/auth/github/callback", callbackHandler)
	r.HandleFunc("/auth/github/set", manualCallbackHandler)
	r.HandleFunc("/auth/github/logout", logoutHandler)
	r.HandleFunc("/auth/github/token", tokenHandler)
	r.HandleFunc("/get/avatar", avatarHandler)
	r.HandleFunc("/auth/{provider}", func(res http.ResponseWriter, req *http.Request) {
		gothic.BeginAuthHandler(res, req)
	})

	server := &http.Server{
		Addr:    ":3000",
		Handler: r,
	}

	fmt.Println("listening on localhost:3000")
	server.ListenAndServe()
}
