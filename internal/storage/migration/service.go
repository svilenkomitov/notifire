package migration

import (
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/svilenkomitov/notifire/internal/storage"
)

const (
	dbURLTemplate  = "postgres://%s:%s@%s:%d/%s?sslmode=disable"
	migrationsPath = "file://db/migrations"
)

type Service interface {
	Up() error
}

type service struct {
	client *migrate.Migrate
}

func New(config *storage.Config) Service {
	client, err := migrate.New(migrationsPath, buildDBURL(config))
	if err != nil {
		panic(err)
	}
	return service{
		client: client,
	}
}

func (s service) Up() error {
	return s.client.Up()
}

func buildDBURL(config *storage.Config) string {
	return fmt.Sprintf(dbURLTemplate, config.User, config.Password,
		config.Host, config.Port, config.DbName)
}
