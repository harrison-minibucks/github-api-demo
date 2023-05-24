package biz

import (
	"context"
	"os"
	"reflect"
	"testing"

	"github.com/go-kratos/kratos/v2/log"
)

var logger = log.With(log.NewStdLogger(os.Stdout),
	"ts", log.DefaultTimestamp,
	"caller", log.DefaultCaller,
)

type MockTodoRepo struct {
	MockSave          func(context.Context, *TodoItem) (*TodoItem, error)
	MockUpdate        func(context.Context, *TodoItem) (*TodoItem, error)
	MockDeleteByID    func(context.Context, string) (*TodoItem, error)
	MockDeleteByTitle func(context.Context, string) ([]*TodoItem, error)
	MockFindByID      func(context.Context, string) (*TodoItem, error)
	MockListAll       func(context.Context) ([]*TodoItem, error)
}

func (m *MockTodoRepo) Save(c context.Context, i *TodoItem) (*TodoItem, error) {
	return m.MockSave(c, i)
}
func (m *MockTodoRepo) Update(c context.Context, i *TodoItem) (*TodoItem, error) {
	return m.MockUpdate(c, i)
}
func (m *MockTodoRepo) DeleteByID(c context.Context, i string) (*TodoItem, error) {
	return m.MockDeleteByID(c, i)
}
func (m *MockTodoRepo) DeleteByTitle(c context.Context, i string) ([]*TodoItem, error) {
	return m.MockDeleteByTitle(c, i)
}
func (m *MockTodoRepo) FindByID(c context.Context, i string) (*TodoItem, error) {
	return m.MockFindByID(c, i)
}
func (m *MockTodoRepo) ListAll(c context.Context) ([]*TodoItem, error) {
	return m.MockListAll(c)
}

type MockLogger struct{}

func TestTodoUsecase_AddItem(t *testing.T) {
	type fields struct {
		repo TodoRepo
		log  *log.Helper
	}
	type args struct {
		ctx  context.Context
		item *TodoItem
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *TodoItem
		wantErr bool
	}{
		{
			name: "Test AddItem",
			fields: fields{
				repo: &MockTodoRepo{
					MockSave: func(ctx context.Context, ti *TodoItem) (*TodoItem, error) {
						return ti, nil
					},
				},
				log: log.NewHelper(logger),
			},
			args: args{
				ctx: context.Background(),
				item: &TodoItem{
					Title:       "item 1",
					Description: "desc 1",
				},
			},
			want: &TodoItem{
				Title:       "item 1",
				Description: "desc 1",
			},
			wantErr: false,
		},
		{
			name: "Test AddItem - Missing title",
			fields: fields{
				repo: &MockTodoRepo{
					MockSave: func(ctx context.Context, ti *TodoItem) (*TodoItem, error) {
						return ti, nil
					},
				},
				log: log.NewHelper(logger),
			},
			args: args{
				ctx: context.Background(),
				item: &TodoItem{
					Description: "desc 1",
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Test AddItem - Missing title",
			fields: fields{
				repo: &MockTodoRepo{
					MockSave: func(ctx context.Context, ti *TodoItem) (*TodoItem, error) {
						return ti, nil
					},
				},
				log: log.NewHelper(logger),
			},
			args: args{
				ctx: context.Background(),
				item: &TodoItem{
					Title: "item 1",
				},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uc := &TodoUsecase{
				repo: tt.fields.repo,
				log:  tt.fields.log,
			}
			got, err := uc.AddItem(tt.args.ctx, tt.args.item)
			if (err != nil) != tt.wantErr {
				t.Errorf("TodoUsecase.AddItem() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TodoUsecase.AddItem() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTodoUsecase_List(t *testing.T) {
	type fields struct {
		repo TodoRepo
		log  *log.Helper
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []*TodoItem
		wantErr bool
	}{
		{
			name: "Test ListItem",
			fields: fields{
				repo: &MockTodoRepo{
					MockListAll: func(ctx context.Context) ([]*TodoItem, error) {
						return []*TodoItem{}, nil
					},
				},
				log: log.NewHelper(logger),
			},
			args: args{
				ctx: context.Background(),
			},
			want:    []*TodoItem{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uc := &TodoUsecase{
				repo: tt.fields.repo,
				log:  tt.fields.log,
			}
			got, err := uc.List(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("TodoUsecase.List() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TodoUsecase.List() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTodoUsecase_Delete(t *testing.T) {
	type fields struct {
		repo TodoRepo
		log  *log.Helper
	}
	type args struct {
		ctx  context.Context
		item *TodoItem
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []*TodoItem
		wantErr bool
	}{
		{
			name: "Test DeleteItem",
			fields: fields{
				repo: &MockTodoRepo{
					MockDeleteByID: func(ctx context.Context, s string) (*TodoItem, error) {
						return &TodoItem{
							Id: "id1",
						}, nil
					},
				},
				log: log.NewHelper(logger),
			},
			args: args{
				ctx: context.Background(),
				item: &TodoItem{
					Id: "id1",
				},
			},
			want: []*TodoItem{
				{
					Id: "id1",
				},
			},
			wantErr: false,
		},
		{
			name: "Test DeleteItem - Title",
			fields: fields{
				repo: &MockTodoRepo{
					MockDeleteByTitle: func(ctx context.Context, s string) ([]*TodoItem, error) {
						return []*TodoItem{
							{
								Id:    "id1",
								Title: "title 1",
							},
						}, nil
					},
				},
				log: log.NewHelper(logger),
			},
			args: args{
				ctx: context.Background(),
				item: &TodoItem{
					Title: "title 1",
				},
			},
			want: []*TodoItem{
				{
					Id:    "id1",
					Title: "title 1",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uc := &TodoUsecase{
				repo: tt.fields.repo,
				log:  tt.fields.log,
			}
			got, err := uc.Delete(tt.args.ctx, tt.args.item)
			if (err != nil) != tt.wantErr {
				t.Errorf("TodoUsecase.Delete() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TodoUsecase.Delete() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTodoUsecase_Mark(t *testing.T) {
	type fields struct {
		repo TodoRepo
		log  *log.Helper
	}
	type args struct {
		ctx  context.Context
		item *TodoItem
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *TodoItem
		wantErr bool
	}{
		{
			name: "Test MarkItem",
			fields: fields{
				repo: &MockTodoRepo{
					MockUpdate: func(ctx context.Context, ti *TodoItem) (*TodoItem, error) {
						return &TodoItem{
							Id: "id1",
						}, nil
					},
				},
				log: log.NewHelper(logger),
			},
			args: args{
				ctx: context.Background(),
				item: &TodoItem{
					Id: "id1",
				},
			},
			want: &TodoItem{
				Id: "id1",
			},
			wantErr: false,
		},
		{
			name: "Test MarkItem - Missing ID",
			fields: fields{
				repo: &MockTodoRepo{},
				log:  log.NewHelper(logger),
			},
			args: args{
				ctx:  context.Background(),
				item: &TodoItem{},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uc := &TodoUsecase{
				repo: tt.fields.repo,
				log:  tt.fields.log,
			}
			got, err := uc.Mark(tt.args.ctx, tt.args.item)
			if (err != nil) != tt.wantErr {
				t.Errorf("TodoUsecase.Mark() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TodoUsecase.Mark() = %v, want %v", got, tt.want)
			}
		})
	}
}
