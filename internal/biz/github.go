package biz

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/harrison-minibucks/github-api-demo/internal/model"
	"gorm.io/gorm"
)

type GitHubRepo interface {
	Save(context.Context, *model.Session) (*model.Session, error)
	SaveUser(context.Context, *model.GitHubUser) (*model.GitHubUser, error)
	Update(context.Context, *model.Session) (*model.Session, error)
	DeleteByGhId(context.Context, uint32) (*model.Session, error)
	FindByID(context.Context, string) (*model.Session, error)
	FindUserByID(context.Context, uint32) (*model.GitHubUser, error)
	ListAll(context.Context) ([]*model.Session, error)
	ListAllUsers(context.Context) ([]*model.GitHubUser, error)
}

type GitHubItem struct {
	Id string
}

type GitHubUsecase struct {
	repo GitHubRepo
	log  *log.Helper
}

func NewGitHubUsecase(repo GitHubRepo, logger log.Logger) *GitHubUsecase {
	return &GitHubUsecase{repo: repo, log: log.NewHelper(logger)}
}

func (uc *GitHubUsecase) ListSessions(ctx context.Context) ([]*model.Session, error) {
	return uc.repo.ListAll(ctx)
}

func (uc *GitHubUsecase) ListUsers(ctx context.Context) ([]*model.GitHubUser, error) {
	return uc.repo.ListAllUsers(ctx)
}

func (uc *GitHubUsecase) GetAvatar(ctx context.Context) (*model.GitHubUser, error) {
	session := ctx.Value(model.SessionKey("session")).(string)
	res, err := uc.repo.FindByID(ctx, session)
	if err != nil {
		return nil, err
	}
	ghUser, err := uc.repo.FindUserByID(ctx, res.GhId)
	if err != nil {
		return nil, err
	}
	return &model.GitHubUser{
		AvatarURL: ghUser.AvatarURL,
	}, nil
}

func (uc *GitHubUsecase) Logout(ctx context.Context) (bool, error) {
	session := ctx.Value(model.SessionKey("session")).(string)
	res, err := uc.repo.FindByID(ctx, session)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return true, nil
		}
		return false, err
	}
	_, err = uc.repo.DeleteByGhId(ctx, res.GhId)
	if err != nil {
		return false, err
	}
	return true, nil
}
