package server

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"

	netHttp "net/http"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/google/uuid"
	"github.com/harrison-minibucks/github-api-demo/internal/biz"
	"github.com/harrison-minibucks/github-api-demo/internal/conf"
	"github.com/harrison-minibucks/github-api-demo/internal/model"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/github"
	"gorm.io/gorm"
)

func initGoth(conf *conf.GitHubApp) {
	if conf == nil || conf.ClientId == "" || conf.ClientSecret == "" || conf.CallbackUrl == "" {
		fmt.Println("GitHubApp will not be initialized")
		return
	}
	goth.UseProviders(
		github.New(conf.ClientId, conf.ClientSecret, conf.CallbackUrl, "user:email"),
	)
	// Convert sessions into database sessions
	// gothic.Store = sessions.NewCookieStore([]byte("SESSION_SECRET"))
	gothic.GetProviderName = func(*netHttp.Request) (string, error) {
		return "github", nil
	}
	gothic.SetState = func(req *http.Request) string {
		return uuid.New().String()
	}
}

// Validates session
func validateSession(ghRepo biz.GitHubRepo, session string) (bool, error) {
	if _, err := ghRepo.FindByID(context.Background(), session); err != nil {
		if err == gorm.ErrRecordNotFound {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// Filter specific requests to make use of gothic auth handler
func GitHubApp(ghRepo biz.GitHubRepo) http.FilterFunc {
	return func(next netHttp.Handler) netHttp.Handler {
		return netHttp.HandlerFunc(func(w netHttp.ResponseWriter, r *netHttp.Request) {
			if r.URL.Path == "/github/login" {
				gothic.BeginAuthHandler(w, r)
				return
			}
			if strings.HasPrefix(r.URL.Path, "/github/callback") {
				user, err := gothic.CompleteUserAuth(w, r)
				if err != nil {
					fmt.Fprintln(w, err)
					return
				}
				ghUser := &model.GitHubUser{}
				userJson, _ := json.Marshal(user.RawData)
				json.Unmarshal(userJson, ghUser)
				ghUser.Email = user.Email
				dataUser, err := ghRepo.FindUserByID(context.Background(), ghUser.Id)
				if err != nil && err != gorm.ErrRecordNotFound {
					fmt.Fprintf(w, "Failed: %s\n", err.Error())
					return
				} else {
					if dataUser != nil {
						if _, err := ghRepo.DeleteByGhId(context.Background(), dataUser.Id); err != nil && err != gorm.ErrRecordNotFound {
							fmt.Fprintf(w, "Failed to delete session: %s\n", err.Error())
							return
						}
					} else {
						if _, err := ghRepo.SaveUser(context.Background(), ghUser); err != nil {
							fmt.Fprintf(w, "Failed to save user: %s\n", err.Error())
							return
						}
					}
					if res, err := ghRepo.Save(context.Background(), &model.Session{
						Id:   uuid.NewString(),
						GhId: ghUser.Id,
					}); err != nil {
						fmt.Fprintf(w, "Failed to save session: %s\n", err.Error())
					} else {
						fmt.Fprintf(w, "Logged in, session id: %s\n", res.Id)
					}
				}
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

// Authenticate user with either GitHub Token / Session
func GitHubAuthenticator(config *conf.Config, ghRepo biz.GitHubRepo) middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (reply interface{}, err error) {
			if config.Env == "prod" {
				if tr, ok := http.RequestFromServerContext(ctx); ok {
					sessionToken := tr.Header.Get("Session")
					if !strings.HasPrefix(tr.URL.Path, "/github") {
						authToken := tr.Header.Get("Authorization")
						valid := false
						if sessionToken != "" {
							valid, err = validateSession(ghRepo, sessionToken)
							if err != nil {
								return nil, err
							}
						}
						// Checks GitHub Token if session is not valid
						if !valid {
							ghUser := &model.GitHubUser{}
							if err := callGitHubAPI(authToken, ghUser); err != nil {
								return nil, errors.New("please include a valid Authorization/Session header")
							}
						}
					}
					if sessionToken != "" {
						ctx = context.WithValue(ctx, model.SessionKey("session"), sessionToken)
					}
				} else {
					return nil, errors.New("failed context retrieval")
				}
			}
			return handler(ctx, req)
		}
	}
}

// Helps to log errors
func LoggingMiddleware(logger log.Logger) middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (reply interface{}, err error) {
			reply, err = handler(ctx, req)
			if err != nil {
				log.WithContext(ctx, logger).Log(log.LevelError, "error", err)
			}
			return reply, err
		}
	}
}

// Perform a call to GitHub API and retrieve user info
func callGitHubAPI(token string, ghUser *model.GitHubUser) error {
	client := &netHttp.Client{}
	req, err := netHttp.NewRequest("GET", "https://api.github.com/user", nil)
	if err != nil {
		return err
	}
	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Add("Accept", "application/vnd.github+json")
	req.Header.Add("X-GitHub-Api-Version", "2022-11-28")

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != netHttp.StatusOK {
		return fmt.Errorf("failed status code: %d", resp.StatusCode)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	err = json.Unmarshal(body, ghUser)
	if err != nil {
		return err
	}
	return nil
}
