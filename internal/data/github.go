package data

import (
	"context"

	"github.com/harrison-minibucks/github-api-demo/internal/biz"
	"github.com/harrison-minibucks/github-api-demo/internal/model"
	"gorm.io/gorm"

	"github.com/go-kratos/kratos/v2/log"
)

type githubRepo struct {
	data *Data
	log  *log.Helper
}

func NewGitHubRepo(data *Data, logger log.Logger) biz.GitHubRepo {
	return &githubRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (r *githubRepo) Save(ctx context.Context, item *model.Session) (*model.Session, error) {
	if err := r.data.db.Transaction(func(tx *gorm.DB) error {
		return tx.Create(item).Error
	}); err != nil {
		return nil, err
	}
	return item, nil
}

func (r *githubRepo) SaveUser(ctx context.Context, item *model.GitHubUser) (*model.GitHubUser, error) {
	if err := r.data.db.Transaction(func(tx *gorm.DB) error {
		return tx.Create(item).Error
	}); err != nil {
		return nil, err
	}
	return item, nil
}

func (r *githubRepo) Update(ctx context.Context, session *model.Session) (*model.Session, error) {
	item := &model.Session{
		Id: session.Id,
	}
	if err := r.data.db.First(&item).Error; err != nil {
		return nil, err
	}
	if err := r.data.db.Transaction(func(tx *gorm.DB) error {
		return tx.Save(item).Error
	}); err != nil {
		return nil, err
	}
	return item, nil
}

func (r *githubRepo) DeleteByGhId(ctx context.Context, ghId uint32) (*model.Session, error) {
	item := &model.Session{}
	if err := r.data.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("gh_id = ?", ghId).First(&item).Error; err != nil {
			return err
		}
		return tx.Delete(&item).Error
	}); err != nil {
		return nil, err
	}
	return item, nil
}

func (r *githubRepo) FindByID(ctx context.Context, id string) (*model.Session, error) {
	item := &model.Session{
		Id: id,
	}
	if err := r.data.db.First(&item).Error; err != nil {
		return nil, err
	}
	return item, nil
}

func (r *githubRepo) FindUserByID(ctx context.Context, id uint32) (*model.GitHubUser, error) {
	item := &model.GitHubUser{
		Id: id,
	}
	if err := r.data.db.First(&item).Error; err != nil {
		return nil, err
	}
	return item, nil
}

func (r *githubRepo) ListAll(context.Context) ([]*model.Session, error) {
	items := []*model.Session{}
	if err := r.data.db.Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

func (r *githubRepo) ListAllUsers(context.Context) ([]*model.GitHubUser, error) {
	items := []*model.GitHubUser{}
	if err := r.data.db.Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}
