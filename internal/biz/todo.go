package biz

import (
	"context"
	"errors"
	"time"

	"github.com/go-kratos/kratos/v2/log"
)

type TodoItem struct {
	Id          string
	Title       string
	Description string
	Marked      bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type TodoRepo interface {
	Save(context.Context, *TodoItem) (*TodoItem, error)
	Update(context.Context, *TodoItem) (*TodoItem, error)
	DeleteByID(context.Context, string) (*TodoItem, error)
	DeleteByTitle(context.Context, string) ([]*TodoItem, error)
	FindByID(context.Context, string) (*TodoItem, error)
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
	uc.log.WithContext(ctx).Infof("Add TODO Item: %v", item.Title)
	if item.Title == "" || item.Description == "" {
		return nil, errors.New("please provide a title and description")
	}
	return uc.repo.Save(ctx, item)
}

func (uc *TodoUsecase) List(ctx context.Context) ([]*TodoItem, error) {
	uc.log.WithContext(ctx).Infof("List TODO Item")
	return uc.repo.ListAll(ctx)
}

func (uc *TodoUsecase) Delete(ctx context.Context, item *TodoItem) ([]*TodoItem, error) {
	// Check if Item ID / title is present
	if item.Id != "" {
		uc.log.WithContext(ctx).Infof("Delete TODO Item: %v", item.Id)
		res, err := uc.repo.DeleteByID(ctx, item.Id)
		if err != nil {
			return nil, err
		}
		return []*TodoItem{res}, nil
	}
	if item.Title == "" {
		return nil, errors.New("please provide item ID or Title")
	}
	uc.log.WithContext(ctx).Infof("Delete TODO Item: %v", item.Title)
	return uc.repo.DeleteByTitle(ctx, item.Title)
}

func (uc *TodoUsecase) Mark(ctx context.Context, item *TodoItem) (*TodoItem, error) {
	uc.log.WithContext(ctx).Infof("Mark TODO Item: %v", item.Id)
	if item.Id == "" {
		return nil, errors.New("please provide an ID to mark an item complete")
	}
	item.Marked = true
	return uc.repo.Update(ctx, item)
}
