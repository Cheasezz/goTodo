package service

import repositories "github.com/Cheasezz/goTodo/internal/repository"

type Services struct {
	*Auth
	*TodoList
	*TodoItem
}

func NewServices(r *repositories.Repositories) *Services {
	return &Services{
		Auth:     newAuthService(r.Psql.Auth),
		TodoList: NewTodoListService(r.Psql.TodoList),
		TodoItem: NewTodoItemService(r.Psql.TodoItem, r.Psql.TodoList),
	}
}
