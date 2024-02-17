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

// List Users godoc
// @Summary list all users
// @Description list all users
// @Tags users
// @Produce json
// @Param page query string false "page number"
// @Param per_page query string false "items per page number"
// @Success 200 {object} dtos.UsersListResponse
// @Failure 400 {object} resterrors.RestErr
// @Failure 500 {object} resterrors.RestErr
// @Router /users [get]
// @Security ApiKeyAuth
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

		if value > 0 {
			currentPage = value
		}
	}

	perPage, exists := c.GetQuery("per_page")
	if exists {
		value, err := strconv.Atoi(perPage)
		if err != nil {
			restErr := resterrors.NewBadRequestError(`Param "per_page" must be a int value`)
			c.JSON(restErr.HttpStatusCode, restErr)
			return
		}

		if value > 0 {
			itemsPerPage = value
		}
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

// Create User godoc
// @Summary create an user
// @Description create an user
// @Tags users
// @Accept json
// @Produce json
// @Param request body dtos.UserRequest true "user request"
// @Success 201 {object} dtos.UsersListResponse
// @Failure 400 {object} resterrors.RestErr
// @Failure 500 {object} resterrors.RestErr
// @Router /users [post]
// @Security ApiKeyAuth
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

// Get User godoc
// @Summary get an user
// @Description get an user
// @Tags users
// @Produce json
// @Param id path string true "user id"
// @Success 200 {object} dtos.UsersListResponse
// @Failure 400 {object} resterrors.RestErr
// @Failure 404 {object} resterrors.RestErr
// @Failure 500 {object} resterrors.RestErr
// @Router /users/{id} [get]
// @Security ApiKeyAuth
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

// Update User godoc
// @Summary update an user
// @Description update an user
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "user id"
// @Param request body dtos.UserRequest true "user request"
// @Success 200 {object} dtos.UsersListResponse
// @Failure 400 {object} resterrors.RestErr
// @Router /users/{id} [put]
// @Security ApiKeyAuth
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

// Delete User godoc
// @Summary delete an user
// @Description delete an user
// @Tags users
// @Produce json
// @Param id path string true "user id"
// @Success 204
// @Failure 400 {object} resterrors.RestErr
// @Failure 500 {object} resterrors.RestErr
// @Router /users/{id} [delete]
// @Security ApiKeyAuth
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
