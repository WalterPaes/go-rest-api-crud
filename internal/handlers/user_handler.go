package handlers

import (
	"net/http"

	"github.com/WalterPaes/go-rest-api-crud/internal/domain"
	"github.com/WalterPaes/go-rest-api-crud/internal/handlers/dtos"
	"github.com/WalterPaes/go-rest-api-crud/pkg/logger"
	"github.com/WalterPaes/go-rest-api-crud/pkg/validation"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type userHandler struct{}

func NewUserHandler() *userHandler {
	return &userHandler{}
}

func (h *userHandler) CreateUser(c *gin.Context) {
	logger.Info("Starting Create User", zap.String("stacktrace", "create-user"))

	var userRequest dtos.UserRequest

	if err := c.ShouldBindJSON(&userRequest); err != nil {
		logger.Error("User Request Validation Error", err, zap.String("stacktrace", "create-user"))

		restErr := validation.ValidationUserError(err)
		c.JSON(restErr.HttpStatusCode, restErr)
		return
	}

	user := domain.User{
		Name:     userRequest.Name,
		Email:    userRequest.Email,
		Password: userRequest.Password,
	}

	logger.Info("User Created Successfully", zap.String("stacktrace", "create-user"))
	c.JSON(http.StatusCreated, dtos.UserResponse{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	})
}
