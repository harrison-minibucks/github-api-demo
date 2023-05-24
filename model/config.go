package model

type Config struct {
	GithubApp struct {
		ClientId     string `yaml:"client_id"`
		ClientSecret string `yaml:"client_secret"`
	} `yaml:"github_app"`
}
