package handlers

import (
	"net/http"

	"github.com/WalterPaes/go-rest-api-crud/internal/handlers/dtos"
	"github.com/WalterPaes/go-rest-api-crud/internal/handlers/dtos/converter"
	"github.com/WalterPaes/go-rest-api-crud/internal/services"
	"github.com/WalterPaes/go-rest-api-crud/pkg/logger"
	resterrors "github.com/WalterPaes/go-rest-api-crud/pkg/rest_errors"
	"github.com/WalterPaes/go-rest-api-crud/pkg/validation"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"
)

var (
	stacktraceCreateUserHandler = zap.String("stacktrace", "create-user-handler")
	stacktraceDeleteUserHandler = zap.String("stacktrace", "delete-user-handler")
)

type userHandler struct {
	userService services.UserService
}

func NewUserHandler(userService services.UserService) *userHandler {
	return &userHandler{
		userService: userService,
	}
}

func (h *userHandler) CreateUser(c *gin.Context) {
	logger.Info("Starting Create User", stacktraceCreateUserHandler)

	var userRequest dtos.UserRequest

	if err := c.ShouldBindJSON(&userRequest); err != nil {
		logger.Error("User Request Validation Error", err, stacktraceCreateUserHandler)

		restErr := validation.ValidationUserError(err)
		c.JSON(restErr.HttpStatusCode, restErr)
		return
	}

	userResult, err := h.userService.CreateUser(c.Request.Context(), converter.UserRequestToUserDomain(userRequest))
	if err != nil {
		logger.Error("Error when trying call service", err, stacktraceCreateUserHandler)

		c.JSON(err.HttpStatusCode, err)
		return
	}

	logger.Info("User Created Successfully", zap.String("user_id", userResult.ID), stacktraceCreateUserHandler)
	c.JSON(http.StatusCreated, converter.UserDomainToUserResponse(userResult))
}

func (h *userHandler) DeleteUser(c *gin.Context) {
	logger.Info("Starting Delete User", stacktraceDeleteUserHandler)

	userID := c.Param("id")
	if _, err := primitive.ObjectIDFromHex(userID); err != nil {
		restErr := resterrors.NewBadRequestError("Invalid userID, must be a hex value")
		c.JSON(restErr.HttpStatusCode, restErr)
		return
	}

	err := h.userService.DeleteUser(c.Request.Context(), userID)
	if err != nil {
		logger.Error("Error when trying call service", err, stacktraceDeleteUserHandler)

		c.JSON(err.HttpStatusCode, err)
		return
	}

	logger.Info("User Deleted Successfully", zap.String("user_id", userID), stacktraceDeleteUserHandler)
	c.Status(http.StatusNoContent)
}
