package service

import (
	"github.com/Cheasezz/goTodo/internal/repository"
)

type Service struct {
	*AuthService
	*TodoListService
	*TodoItemService
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		AuthService:     newAuthService(repos.Authorization),
		TodoListService: NewTodoListService(repos.TodoList),
		TodoItemService: NewTodoItemService(repos.TodoItem, repos.TodoList),
	}
}
