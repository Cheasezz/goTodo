package service

import (
	"github.com/Cheasezz/goTodo/internal/core"
)

type TodoItemRepo interface {
	Create(listId int, item core.TodoItem) (int, error)
	GetAll(userId, listId int) ([]core.TodoItem, error)
	GetById(userId, itemId int) (core.TodoItem, error)
	Delete(userId, itemId int) error
	Update(userId, itemId int, input core.UpdateItemInput) error
}

type TodoItem struct {
	repo     TodoItemRepo
	listRepo TodoListRepo
}

func NewTodoItemService(repo TodoItemRepo, listRepo TodoListRepo) *TodoItem {
	return &TodoItem{repo: repo, listRepo: listRepo}
}

func (s *TodoItem) Create(userId, listId int, item core.TodoItem) (int, error) {
	_, err := s.listRepo.GetById(userId, listId)
	if err != nil {
		return 0, err
	}

	return s.repo.Create(listId, item)
}

func (s *TodoItem) GetAll(userId, listId int) ([]core.TodoItem, error) {
	return s.repo.GetAll(userId, listId)
}

func (s *TodoItem) GetById(userId, itemId int) (core.TodoItem, error) {
	return s.repo.GetById(userId, itemId)
}

func (s *TodoItem) Delete(userId, itemId int) error {
	return s.repo.Delete(userId, itemId)
}

func (s *TodoItem) Update(userId, itemId int, input core.UpdateItemInput) error {
	return s.repo.Update(userId, itemId, input)
}
