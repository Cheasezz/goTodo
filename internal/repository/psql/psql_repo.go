package psql

import (
	"github.com/Cheasezz/goTodo/pkg/postgres"
)

const (
	userTable       = "users"
	todoListsTable  = "todo_lists"
	usersListsTable = "users_lists"
	todoItemsTable  = "todo_items"
	listsItemsTable = "lists_items"
)

type Repository struct {
	*Auth
	*TodoList
	*TodoItem
}

func NewPsqlRepository(db *postgres.Postgres) *Repository {
	return &Repository{
		Auth:     NewAuthPostgres(db),
		TodoList: NewTodoListPostgres(db),
		TodoItem: NewTodoItemPostgres(db),
	}
}
