package repositories

import (
	"github.com/Cheasezz/goTodo/internal/repository/psql"
	"github.com/Cheasezz/goTodo/pkg/postgres"
)

type Repositories struct {
	Psql *psql.Repository
}

func NewRepositories(postgres *postgres.Postgres) *Repositories {
	return &Repositories{
		Psql: psql.NewPsqlRepository(postgres),
	}
}
