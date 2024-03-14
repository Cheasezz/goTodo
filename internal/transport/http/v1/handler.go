package v1

import (
	_ "github.com/Cheasezz/goTodo/docs"
	"github.com/Cheasezz/goTodo/internal/service"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	*AuthHandler
	*TodoListHandler
	*TodoItemHandler
}

func NewHandlers(s *service.Services) *Handler {
	return &Handler{
		AuthHandler:     NewAuthHandler(s.Auth),
		TodoListHandler: NewTodoListHandler(s.TodoList),
		TodoItemHandler: NewTodoItemHandler(s.TodoItem),
	}
}

func (h *Handler) InitRoutes(router *gin.RouterGroup) {

	v1 := router.Group("/v1")
	{
		h.initAuthRoutes(v1)

		todo := router.Group("/v1/todo", h.userIdentity)
		{
			h.initListRoutes(todo)
			h.initItemRoutes(todo)
		}
	}
}
