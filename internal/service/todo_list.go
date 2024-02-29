package service

import (
	"github.com/Cheasezz/goTodo/internal/core"
)

type TodoListRepo interface {
	Create(userId int, list core.TodoList) (int, error)
	GetAll(userId int) ([]core.TodoList, error)
	GetById(userId, listId int) (core.TodoList, error)
	Delete(userId, listId int) error
	Update(userId, listId int, input core.UpdateListInput) error
}

type TodoList struct {
	repo TodoListRepo
}

func NewTodoListService(repo TodoListRepo) *TodoList {
	return &TodoList{repo: repo}
}

func (s *TodoList) Create(userId int, list core.TodoList) (int, error) {
	return s.repo.Create(userId, list)
}

func (s *TodoList) GetAll(userId int) ([]core.TodoList, error) {
	return s.repo.GetAll(userId)
}

func (s *TodoList) GetById(userId, listId int) (core.TodoList, error) {
	return s.repo.GetById(userId, listId)
}

func (s *TodoList) Delete(userId, listId int) error {
	return s.repo.Delete(userId, listId)
}

func (s *TodoList) Update(userId, listId int, input core.UpdateListInput) error {
	if err := input.Validate(); err != nil {
		return err
	}
	return s.repo.Update(userId, listId, input)
}
