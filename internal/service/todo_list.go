package service

import (
	"context"

	"github.com/Cheasezz/goTodo/internal/core"
)

type TodoListRepo interface {
	Create(ctx context.Context, userId int, list core.TodoList) (int, error)
	GetAll(ctx context.Context, userId int) ([]core.TodoList, error)
	GetById(ctx context.Context, userId, listId int) (core.TodoList, error)
	Delete(ctx context.Context, userId, listId int) error
	Update(ctx context.Context, userId, listId int, input core.UpdateListInput) error
}

type TodoList struct {
	repo TodoListRepo
}

func NewTodoListService(repo TodoListRepo) *TodoList {
	return &TodoList{repo: repo}
}

func (s *TodoList) Create(ctx context.Context, userId int, list core.TodoList) (int, error) {
	return s.repo.Create(ctx, userId, list)
}

func (s *TodoList) GetAll(ctx context.Context, userId int) ([]core.TodoList, error) {
	return s.repo.GetAll(ctx, userId)
}

func (s *TodoList) GetById(ctx context.Context, userId, listId int) (core.TodoList, error) {
	return s.repo.GetById(ctx, userId, listId)
}

func (s *TodoList) Delete(ctx context.Context, userId, listId int) error {
	return s.repo.Delete(ctx, userId, listId)
}

func (s *TodoList) Update(ctx context.Context, userId, listId int, input core.UpdateListInput) error {
	if err := input.Validate(); err != nil {
		return err
	}
	return s.repo.Update(ctx, userId, listId, input)
}
