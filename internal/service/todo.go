package service

import (
	"context"
	"errors"

	v1 "github.com/harrison-minibucks/github-api-demo/api/todo/v1"
	"github.com/harrison-minibucks/github-api-demo/internal/biz"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type TodoService struct {
	v1.UnimplementedTodoServer

	uc *biz.TodoUsecase
}

func NewTodoService(uc *biz.TodoUsecase) *TodoService {
	return &TodoService{uc: uc}
}

func (s *TodoService) Add(ctx context.Context, in *v1.AddRequest) (*v1.AddReply, error) {
	if in == nil || in.Item == nil {
		return nil, errors.New("please provide todo item")
	}
	res, err := s.uc.AddItem(ctx, &biz.TodoItem{
		Title:       in.Item.Title,
		Description: in.Item.Description,
	})
	if err != nil {
		return nil, err
	}
	return &v1.AddReply{
		Message: "Added " + res.Title,
		Item:    mapItem(res),
	}, nil
}

func (s *TodoService) List(ctx context.Context, in *v1.ListRequest) (*v1.ListReply, error) {
	items, err := s.uc.List(ctx)
	if err != nil {
		return nil, err
	}
	replyList := []*v1.Item{}
	for i := 0; i < len(items); i++ {
		replyList = append(replyList, mapItem(items[i]))
	}
	return &v1.ListReply{Items: replyList}, nil
}

func (s *TodoService) Delete(ctx context.Context, in *v1.DeleteRequest) (*v1.DeleteReply, error) {
	items, err := s.uc.Delete(ctx, &biz.TodoItem{
		Id:    in.Id,
		Title: in.Title,
	})
	if err != nil {
		return nil, err
	}
	replyList := []*v1.Item{}
	for i := 0; i < len(items); i++ {
		replyList = append(replyList, mapItem(items[i]))
	}
	return &v1.DeleteReply{
		Message: "Deleted the following items",
		Items:   replyList,
	}, nil
}

func (s *TodoService) Mark(ctx context.Context, in *v1.MarkRequest) (*v1.MarkReply, error) {
	res, err := s.uc.Mark(ctx, &biz.TodoItem{
		Id: in.Id,
	})
	if err != nil {
		return nil, err
	}
	return &v1.MarkReply{
		Message: "Marked Item as complete",
		Item:    mapItem(res),
	}, nil
}

func mapItem(item *biz.TodoItem) *v1.Item {
	return &v1.Item{
		Id:          item.Id,
		Title:       item.Title,
		Description: item.Description,
		Marked:      item.Marked,
		CreatedAt:   timestamppb.New(item.CreatedAt),
	}
}
