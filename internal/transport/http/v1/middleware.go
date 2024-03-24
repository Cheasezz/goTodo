package v1

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

const (
	authorizationHeader = "Authorization"
	userCtx             = "userId"
)

func (h *AuthHandler) userIdentity(c *gin.Context) {
	logrus.Print("From start userIdentity")
	header := c.GetHeader(authorizationHeader)
	if header == "" {
		newErrorResponse(c, http.StatusUnauthorized, "empty auth header")
		return
	}
	headerParts := strings.Split(header, " ")
	if len(headerParts) != 2 {
		newErrorResponse(c, http.StatusUnauthorized, "invalid auth header")
		return
	}

	// userId, err := h.service.ParseToken(headerParts[1])
	userId, err := h.TokenManager.Parse(headerParts[1])
	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}
	logrus.Printf("From start userIdentity and before c.Set, userId: %s", userId)
	c.Set(userCtx, userId)
	logrus.Print("From start userIdentity and after c.Set")

}

func getUserId(c *gin.Context) (int, error) {
	id, ok := c.Get(userCtx)
	logrus.Printf("From start getUserId, userId from ctx: %s", id)
	if !ok {
		return 0, errors.New("user id not found")
	}

	idInt, err := strconv.Atoi(id.(string))
	logrus.Printf("From start getUserId, idInt: %d", idInt)
	if err !=nil {
		return 0, errors.New("user id not found")
	}

	return idInt, nil
}
