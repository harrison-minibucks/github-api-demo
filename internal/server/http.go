package server

import (
	ghV1 "github.com/harrison-minibucks/github-api-demo/api/github/v1"
	v1 "github.com/harrison-minibucks/github-api-demo/api/todo/v1"
	"github.com/harrison-minibucks/github-api-demo/internal/biz"
	"github.com/harrison-minibucks/github-api-demo/internal/conf"
	"github.com/harrison-minibucks/github-api-demo/internal/service"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport/http"
)

// NewHTTPServer new an HTTP server.
func NewHTTPServer(c *conf.Server, config *conf.Config, conf *conf.GitHubApp, ghRepo biz.GitHubRepo, todo *service.TodoService, gh *service.GitHubService, logger log.Logger) *http.Server {
	initGoth(conf)
	var opts = []http.ServerOption{
		http.Middleware(
			recovery.Recovery(),
			LoggingMiddleware(logger),
			GitHubAuthenticator(config, ghRepo),
		),
		http.Filter(
			GitHubApp(ghRepo),
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
	ghV1.RegisterGitHubHTTPServer(srv, gh)
	return srv
}
