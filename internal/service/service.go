package service

import (
	repositories "github.com/Cheasezz/goTodo/internal/repository"
	"github.com/Cheasezz/goTodo/pkg/auth"
	"github.com/Cheasezz/goTodo/pkg/hash"
)

type Services struct {
	*Auth
	*TodoList
	*TodoItem
}

type Deps struct {
	Repos        *repositories.Repositories
	Hasher       hash.PasswordHasher
	TokenManager auth.TokenManager
}

func NewServices(d Deps) *Services {
	return &Services{
		Auth:     newAuthService(d.Repos.Psql.Auth, d.Hasher, d.TokenManager),
		TodoList: NewTodoListService(d.Repos.Psql.TodoList),
		TodoItem: NewTodoItemService(d.Repos.Psql.TodoItem, d.Repos.Psql.TodoList),
	}
}
