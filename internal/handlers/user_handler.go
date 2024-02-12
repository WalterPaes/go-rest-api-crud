package handlers

import (
	"net/http"
	"strconv"

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

const (
	errUserRequestValidation = "User Request Validation Error"
	errTryCallService        = "Error when try call service"
)

var (
	stacktraceFindAllUsersHandler = zap.String("stacktrace", "find-all-users-handler")
	stacktraceCreateUserHandler   = zap.String("stacktrace", "create-user-handler")
	stacktraceFindUserByIdHandler = zap.String("stacktrace", "find-user-by-id-handler")
	stacktraceUpdateUserHandler   = zap.String("stacktrace", "update-user-handler")
	stacktraceDeleteUserHandler   = zap.String("stacktrace", "delete-user-handler")
)

type userHandler struct {
	userService services.UserService
}

func NewUserHandler(userService services.UserService) *userHandler {
	return &userHandler{
		userService: userService,
	}
}

func (h *userHandler) ListAll(c *gin.Context) {
	logger.Info("Starting Find User By Id", stacktraceFindAllUsersHandler)

	var (
		currentPage  int = 1
		itemsPerPage int = 10
	)

	page, exists := c.GetQuery("page")
	if exists {
		value, err := strconv.Atoi(page)
		if err != nil {
			restErr := resterrors.NewBadRequestError(`Param "page" must be a int value`)
			c.JSON(restErr.HttpStatusCode, restErr)
			return
		}
		currentPage = value
	}

	perPage, exists := c.GetQuery("per_page")
	if exists {
		value, err := strconv.Atoi(perPage)
		if err != nil {
			restErr := resterrors.NewBadRequestError(`Param "per_page" must be a int value`)
			c.JSON(restErr.HttpStatusCode, restErr)
			return
		}
		itemsPerPage = value
	}

	userResult, err := h.userService.FindAll(c.Request.Context(), itemsPerPage, currentPage)
	if err != nil {
		logger.Error(errTryCallService, err, stacktraceFindUserByIdHandler)

		c.JSON(err.HttpStatusCode, err)
		return
	}

	logger.Info(
		"User Found Successfully",
		zap.Int("items_per_page", itemsPerPage),
		zap.Int("current_page", currentPage),
		stacktraceFindAllUsersHandler,
	)
	c.JSON(http.StatusOK, converter.UsersDomainListToUserListResponse(userResult, currentPage, itemsPerPage))
}

func (h *userHandler) CreateUser(c *gin.Context) {
	logger.Info("Starting Create User", stacktraceCreateUserHandler)

	var userRequest dtos.UserRequest

	if err := c.ShouldBindJSON(&userRequest); err != nil {
		logger.Error(errUserRequestValidation, err, stacktraceCreateUserHandler)

		restErr := validation.ValidationUserError(err)
		c.JSON(restErr.HttpStatusCode, restErr)
		return
	}

	user, convertionErr := converter.UserRequestToUserDomain(userRequest)
	if convertionErr != nil {
		restErr := resterrors.NewInternalServerError("Error when try create user")
		c.JSON(restErr.HttpStatusCode, restErr)
		return
	}

	userResult, err := h.userService.CreateUser(c.Request.Context(), user)
	if err != nil {
		logger.Error(errTryCallService, err, stacktraceCreateUserHandler)

		c.JSON(err.HttpStatusCode, err)
		return
	}

	logger.Info("User Created Successfully", zap.String("user_id", userResult.ID), stacktraceCreateUserHandler)
	c.JSON(http.StatusCreated, converter.UserDomainToUserResponse(userResult))
}

func (h *userHandler) GetUserById(c *gin.Context) {
	logger.Info("Starting Find User By Id", stacktraceFindUserByIdHandler)

	userID, err := h.getIdFromParam(c)
	if err != nil {
		c.JSON(err.HttpStatusCode, err)
		return
	}

	userResult, err := h.userService.FindUserById(c.Request.Context(), userID)
	if err != nil {
		logger.Error(errTryCallService, err, stacktraceFindUserByIdHandler)

		c.JSON(err.HttpStatusCode, err)
		return
	}

	logger.Info("User Found Successfully", zap.String("user_id", userResult.ID), stacktraceFindUserByIdHandler)
	c.JSON(http.StatusOK, converter.UserDomainToUserResponse(userResult))
}

func (h *userHandler) UpdateUser(c *gin.Context) {
	logger.Info("Starting Update User", stacktraceUpdateUserHandler)

	userID, err := h.getIdFromParam(c)
	if err != nil {
		c.JSON(err.HttpStatusCode, err)
		return
	}

	var userRequest dtos.UserRequest
	if err := c.ShouldBindJSON(&userRequest); err != nil {
		logger.Error(errUserRequestValidation, err, stacktraceCreateUserHandler)

		restErr := validation.ValidationUserError(err)
		c.JSON(restErr.HttpStatusCode, restErr)
		return
	}

	user, convertionErr := converter.UserRequestToUserDomain(userRequest)
	if convertionErr != nil {
		restErr := resterrors.NewInternalServerError("Error when try update user")
		c.JSON(restErr.HttpStatusCode, restErr)
		return
	}

	userResult, err := h.userService.UpdateUser(c.Request.Context(), userID, user)
	if err != nil {
		logger.Error(errTryCallService, err, stacktraceUpdateUserHandler)

		c.JSON(err.HttpStatusCode, err)
		return
	}

	logger.Info("User Updated Successfully", zap.String("user_id", userResult.ID), stacktraceUpdateUserHandler)
	c.JSON(http.StatusCreated, converter.UserDomainToUserResponse(userResult))
}

func (h *userHandler) DeleteUser(c *gin.Context) {
	logger.Info("Starting Delete User", stacktraceDeleteUserHandler)

	userID, err := h.getIdFromParam(c)
	if err != nil {
		c.JSON(err.HttpStatusCode, err)
		return
	}

	err = h.userService.DeleteUser(c.Request.Context(), userID)
	if err != nil {
		logger.Error(errTryCallService, err, stacktraceDeleteUserHandler)

		c.JSON(err.HttpStatusCode, err)
		return
	}

	logger.Info("User Deleted Successfully", zap.String("user_id", userID), stacktraceDeleteUserHandler)
	c.Status(http.StatusNoContent)
}

func (*userHandler) getIdFromParam(c *gin.Context) (string, *resterrors.RestErr) {
	userID := c.Param("id")
	if _, err := primitive.ObjectIDFromHex(userID); err != nil {
		restErr := resterrors.NewBadRequestError("Invalid userID, must be a hex value")
		logger.Error(restErr.Message, restErr, stacktraceDeleteUserHandler)
		return "", restErr
	}
	return userID, nil
}
