package service

import (
	"github.com/Cheasezz/goTodo/internal/repository"
)

type Service struct {
	*Auth
	*TodoList
	*TodoItem
}

func NewServices(r *repository.Repository) *Service {
	return &Service{
		Auth:     newAuthService(r.Auth),
		TodoList: NewTodoListService(r.TodoList),
		TodoItem: NewTodoItemService(r.TodoItem, r.TodoList),
	}
}
