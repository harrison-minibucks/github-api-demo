package data

import (
	"fmt"

	"github.com/harrison-minibucks/github-api-demo/internal/conf"
	"github.com/harrison-minibucks/github-api-demo/internal/model"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(NewData, NewTodoRepo, NewGitHubRepo, NewDB)

// Data .
type Data struct {
	db  *gorm.DB
	log *log.Helper
}

// NewData .
func NewData(c *conf.Data, db *gorm.DB, logger log.Logger) (*Data, func(), error) {
	l := log.NewHelper(log.With(logger, "module", "transaction/data"))
	d := &Data{
		db:  db,
		log: l,
	}
	return d, func() {
	}, nil
}

// NewDB gorm Connecting to a Database
func NewDB(conf *conf.Data, logger log.Logger) *gorm.DB {
	log := log.NewHelper(log.With(logger, "module", "order-service/data/gorm"))
	dsn := fmt.Sprintf("postgresql://%s:%s@%s:%d/%s",
		conf.Database.Username,
		conf.Database.Password,
		conf.Database.Source,
		conf.Database.Port,
		conf.Database.DbName,
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed opening connection to postgres: %v", err)
	}
	if err := db.AutoMigrate(&model.Item{}, &model.Session{}, &model.GitHubUser{}); err != nil {
		log.Fatal(err)
	}
	return db
}
