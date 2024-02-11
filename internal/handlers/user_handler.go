package handlers

import (
	"github.com/WalterPaes/go-rest-api-crud/pkg/logger"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type UserHandler interface {
	ListUsers(c *gin.Context)
	FindUserByID(c *gin.Context)
	CreateUser(c *gin.Context)
	UpdateUser(c *gin.Context)
	DeleteUser(c *gin.Context)
}

type userHandler struct {
	logger *logger.Logger
}

func NewUserHandler(logger *logger.Logger) *userHandler {
	return &userHandler{
		logger: logger,
	}
}

func (h *userHandler) CreateUser(c *gin.Context) {
	h.logger.Info("Starting Create User", zap.String("stacktrace", "create-user"))
}
