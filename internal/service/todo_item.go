package service

import (
	"context"

	"github.com/Cheasezz/goTodo/internal/core"
)

type TodoItemRepo interface {
	Create(ctx context.Context, listId int, item core.TodoItem) (int, error)
	GetAll(ctx context.Context, userId, listId int) ([]core.TodoItem, error)
	GetById(ctx context.Context, userId, itemId int) (core.TodoItem, error)
	Delete(ctx context.Context, userId, itemId int) error
	Update(ctx context.Context, userId, itemId int, input core.UpdateItemInput) error
}

type TodoItem struct {
	repo     TodoItemRepo
	listRepo TodoListRepo
}

func NewTodoItemService(repo TodoItemRepo, listRepo TodoListRepo) *TodoItem {
	return &TodoItem{repo: repo, listRepo: listRepo}
}

func (s *TodoItem) Create(ctx context.Context, userId, listId int, item core.TodoItem) (int, error) {
	_, err := s.listRepo.GetById(ctx, userId, listId)
	if err != nil {
		return 0, err
	}

	return s.repo.Create(ctx, listId, item)
}

func (s *TodoItem) GetAll(ctx context.Context, userId, listId int) ([]core.TodoItem, error) {
	return s.repo.GetAll(ctx, userId, listId)
}

func (s *TodoItem) GetById(ctx context.Context, userId, itemId int) (core.TodoItem, error) {
	return s.repo.GetById(ctx, userId, itemId)
}

func (s *TodoItem) Delete(ctx context.Context, userId, itemId int) error {
	return s.repo.Delete(ctx, userId, itemId)
}

func (s *TodoItem) Update(ctx context.Context, userId, itemId int, input core.UpdateItemInput) error {
	return s.repo.Update(ctx, userId, itemId, input)
}
