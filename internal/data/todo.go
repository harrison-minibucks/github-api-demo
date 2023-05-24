package data

import (
	"context"

	"github.com/google/uuid"
	"github.com/harrison-minibucks/github-api-demo/internal/biz"
	"github.com/harrison-minibucks/github-api-demo/internal/model"
	"gorm.io/gorm"

	"github.com/go-kratos/kratos/v2/log"
)

type todoRepo struct {
	data *Data
	log  *log.Helper
}

func NewTodoRepo(data *Data, logger log.Logger) biz.TodoRepo {
	return &todoRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (r *todoRepo) Save(ctx context.Context, item *biz.TodoItem) (*biz.TodoItem, error) {
	todoItem := &model.Item{
		Id:          uuid.NewString(),
		Title:       item.Title,
		Description: item.Description,
	}
	if err := r.data.db.Transaction(func(tx *gorm.DB) error {
		return tx.Create(todoItem).Error
	}); err != nil {
		return nil, err
	}
	return item, nil
}

func (r *todoRepo) Update(ctx context.Context, todoitem *biz.TodoItem) (*biz.TodoItem, error) {
	item := &model.Item{
		Id: todoitem.Id,
	}
	if err := r.data.db.First(&item).Error; err != nil {
		return nil, err
	}
	item.Marked = todoitem.Marked
	if err := r.data.db.Transaction(func(tx *gorm.DB) error {
		return tx.Save(item).Error
	}); err != nil {
		return nil, err
	}
	return mapTodoItem(item), nil
}

func (r *todoRepo) FindByID(ctx context.Context, id string) (*biz.TodoItem, error) {
	item := &model.Item{
		Id: id,
	}
	if err := r.data.db.First(&item).Error; err != nil {
		return nil, err
	}
	return mapTodoItem(item), nil
}

func (r *todoRepo) ListByHello(context.Context, string) ([]*biz.TodoItem, error) {
	return nil, nil
}

func (r *todoRepo) ListAll(context.Context) ([]*biz.TodoItem, error) {
	items := []*model.Item{}
	todoList := []*biz.TodoItem{}
	if err := r.data.db.Find(&items).Error; err != nil {
		return nil, err
	}
	for i := 0; i < len(items); i++ {
		todoList = append(todoList, mapTodoItem(items[i]))
	}
	return todoList, nil
}

func mapTodoItem(item *model.Item) *biz.TodoItem {
	return &biz.TodoItem{
		Id:          item.Id,
		Title:       item.Title,
		Description: item.Description,
		Marked:      item.Marked,
		CreatedAt:   item.CreatedAt,
	}
}
