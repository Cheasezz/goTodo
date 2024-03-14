package http

import (
	"github.com/Cheasezz/goTodo/internal/service"
	v1 "github.com/Cheasezz/goTodo/internal/transport/http/v1"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Handler struct {
	services *service.Services
}

func NewHandlers(services *service.Services) *Handler {
	return &Handler{
		services: services,
	}
}

func (h *Handler) Init() *gin.Engine {
	router := gin.New()

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	h.InitRoutes(router)

	return router
}

func (h *Handler) InitRoutes(router *gin.Engine) {
	handlerV1 := v1.NewHandlers(h.services)
	api := router.Group("/api")
	{
		handlerV1.InitRoutes(api)
	}
}
