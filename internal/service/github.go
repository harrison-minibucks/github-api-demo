package service

import (
	"context"

	v1 "github.com/harrison-minibucks/github-api-demo/api/github/v1"
	"github.com/harrison-minibucks/github-api-demo/internal/biz"
)

type GitHubService struct {
	v1.UnimplementedGitHubServer

	uc *biz.GitHubUsecase
}

func NewGitHubService(uc *biz.GitHubUsecase) *GitHubService {
	return &GitHubService{uc: uc}
}

func (s *GitHubService) ListUsers(ctx context.Context, in *v1.ListRequest) (*v1.ListUsersReply, error) {
	res, err := s.uc.ListUsers(ctx)
	if err != nil {
		return nil, err
	}
	reply := &v1.ListUsersReply{}
	for i := 0; i < len(res); i++ {
		reply.User = append(reply.User, &v1.User{
			Id:    res[i].Id,
			Email: res[i].Email,
			Login: res[i].Login,
		})
	}
	return reply, nil
}

func (s *GitHubService) ListSessions(ctx context.Context, in *v1.ListRequest) (*v1.ListSessionsReply, error) {
	res, err := s.uc.ListSessions(ctx)
	if err != nil {
		return nil, err
	}
	reply := &v1.ListSessionsReply{}
	for i := 0; i < len(res); i++ {
		reply.Session = append(reply.Session, &v1.Session{
			Id:   res[i].Id,
			GhId: res[i].GhId,
		})
	}
	return reply, nil
}
