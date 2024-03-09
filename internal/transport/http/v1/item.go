package v1

import (
	"net/http"
	"strconv"

	"github.com/Cheasezz/goTodo/internal/core"
	"github.com/gin-gonic/gin"
)

type TodoItemService interface {
	Create(userId, listId int, item core.TodoItem) (int, error)
	GetAll(userId, listId int) ([]core.TodoItem, error)
	GetById(userId, itemId int) (core.TodoItem, error)
	Delete(userId, itemId int) error
	Update(userId, itemId int, input core.UpdateItemInput) error
}

type TodoItemHandler struct {
	service TodoItemService
}

func (h *Handler) initItemRoutes(router *gin.RouterGroup) {
	items := router.Group("items")
	{
		items.GET("/:id", h.getItemById)
		items.PUT("/:id", h.updateItem)
		items.DELETE("/:id", h.deleteItem)
	}
}

func NewTodoItemHandler(s TodoItemService) *TodoItemHandler {
	return &TodoItemHandler{service: s}
}

func (h *TodoItemHandler) createItem(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		return
	}

	listId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid list id param")
		return
	}

	var input core.TodoItem
	if err := c.ShouldBindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	id, err := h.service.Create(userId, listId, input)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id": id,
	})
}

func (h *TodoItemHandler) getAllItems(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		return
	}

	listId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid list id param")
		return
	}

	items, err := h.service.GetAll(userId, listId)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, items)
}

func (h *TodoItemHandler) getItemById(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		return
	}

	itemId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid list id param")
		return
	}

	item, err := h.service.GetById(userId, itemId)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, item)
}

func (h *TodoItemHandler) updateItem(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid id param")
		return
	}

	var input core.UpdateItemInput
	if err := c.ShouldBindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.service.Update(userId, id, input); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, statusResponse{
		Status: "ok",
	})
}

func (h *TodoItemHandler) deleteItem(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		return
	}

	itemId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid list id param")
		return
	}

	err = h.service.Delete(userId, itemId)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, statusResponse{"ok"})
}
