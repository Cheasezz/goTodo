package v1

import (
	_ "github.com/Cheasezz/goTodo/docs"
	"github.com/Cheasezz/goTodo/internal/service"
	"github.com/Cheasezz/goTodo/pkg/auth"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	*AuthHandler
	*TodoListHandler
	*TodoItemHandler
	// TokenManager auth.TokenManager
}

func NewHandlers(s *service.Services, tm auth.TokenManager) *Handler {
	return &Handler{
		AuthHandler:     NewAuthHandler(s.Auth, tm),
		TodoListHandler: NewTodoListHandler(s.TodoList),
		TodoItemHandler: NewTodoItemHandler(s.TodoItem),
		// TokenManager: tm,
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
