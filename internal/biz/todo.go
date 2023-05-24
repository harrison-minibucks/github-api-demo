package biz

import (
	"context"
	"time"

	"github.com/go-kratos/kratos/v2/log"
)

type TodoItem struct {
	Id          string
	Title       string
	Description string
	Marked      bool
	CreatedAt   time.Time
}

type TodoRepo interface {
	Save(context.Context, *TodoItem) (*TodoItem, error)
	Update(context.Context, *TodoItem) (*TodoItem, error)
	FindByID(context.Context, string) (*TodoItem, error)
	ListByHello(context.Context, string) ([]*TodoItem, error)
	ListAll(context.Context) ([]*TodoItem, error)
}

type TodoUsecase struct {
	repo TodoRepo
	log  *log.Helper
}

func NewTodoUsecase(repo TodoRepo, logger log.Logger) *TodoUsecase {
	return &TodoUsecase{repo: repo, log: log.NewHelper(logger)}
}

func (uc *TodoUsecase) AddItem(ctx context.Context, item *TodoItem) (*TodoItem, error) {
	uc.log.WithContext(ctx).Infof("Create TODO Item: %v", item.Title)
	return uc.repo.Save(ctx, item)
}

func (uc *TodoUsecase) List(ctx context.Context) ([]*TodoItem, error) {
	return uc.repo.ListAll(ctx)
}
