package server

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"

	v1 "github.com/harrison-minibucks/github-api-demo/api/todo/v1"
	"github.com/harrison-minibucks/github-api-demo/internal/conf"
	"github.com/harrison-minibucks/github-api-demo/internal/model"
	"github.com/harrison-minibucks/github-api-demo/internal/service"

	netHttp "net/http"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport/http"
)

// NewHTTPServer new an HTTP server.
func NewHTTPServer(c *conf.Server, config *conf.Config, todo *service.TodoService, logger log.Logger) *http.Server {
	var opts = []http.ServerOption{
		http.Middleware(
			recovery.Recovery(),
			LoggingMiddleware(logger),
			GitHubAuthenticator(config),
		),
	}
	if c.Http.Network != "" {
		opts = append(opts, http.Network(c.Http.Network))
	}
	if c.Http.Addr != "" {
		opts = append(opts, http.Address(c.Http.Addr))
	}
	if c.Http.Timeout != nil {
		opts = append(opts, http.Timeout(c.Http.Timeout.AsDuration()))
	}
	srv := http.NewServer(opts...)
	v1.RegisterTodoHTTPServer(srv, todo)
	return srv
}

// TODO: Validate github token
func GitHubAuthenticator(config *conf.Config) middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (reply interface{}, err error) {
			if config.Env == "prod" {
				if tr, ok := http.RequestFromServerContext(ctx); ok {
					authToken := tr.Header.Get("Authorization")
					ghUser := &model.GitHubUser{}
					if err := callGitHubAPI(authToken, ghUser); err != nil {
						return nil, errors.New("please include a valid Authorization header")
					}
				} else {
					return nil, errors.New("failed context retrieval")
				}
			}
			return handler(ctx, req)
		}
	}
}

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
	// fmt.Println(ghUser)
	return nil
}
