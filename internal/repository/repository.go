package repository

import (
	"github.com/jmoiron/sqlx"
)

type Repository struct {
	*Auth
	*TodoList
	*TodoItem
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Auth: NewAuthPostgres(db),
		TodoList:      NewTodoListPostgres(db),
		TodoItem:      NewTodoItemPostgres(db),
	}
}
