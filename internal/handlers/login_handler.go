package handlers

import (
	"net/http"

	"github.com/WalterPaes/go-rest-api-crud/internal/handlers/dtos"
	"github.com/WalterPaes/go-rest-api-crud/internal/handlers/dtos/converter"
	"github.com/WalterPaes/go-rest-api-crud/internal/services"
	"github.com/WalterPaes/go-rest-api-crud/pkg/logger"
	"github.com/WalterPaes/go-rest-api-crud/pkg/validation"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

var (
	stacktraceLoginUserHandler = zap.String("stacktrace", "login-user-handler")
)

type loginHandler struct {
	loginService services.LoginService
}

func NewLoginHandler(loginService services.LoginService) *loginHandler {
	return &loginHandler{
		loginService: loginService,
	}
}

// User Login godoc
// @Summary Login an user
// @Description Login an user
// @Tags login
// @Accept json
// @Produce json
// @Param request body dtos.LoginRequest true "Login Request"
// @Success 200 {object} dtos.LoginResponse
// @Failure 400 {object} resterrors.RestErr
// @Failure 401
// @Failure 404 {object} resterrors.RestErr
// @Router /login [post]
func (h *loginHandler) Login(c *gin.Context) {
	logger.Info("Starting Login User Handler", stacktraceLoginUserHandler)

	var loginRequest dtos.LoginRequest

	if err := c.ShouldBindJSON(&loginRequest); err != nil {
		logger.Error("Login Request Validation Error", err, stacktraceLoginUserHandler)

		restErr := validation.ValidationUserError(err)
		c.JSON(restErr.HttpStatusCode, restErr)
		return
	}

	jwtToken, err := h.loginService.LoginUser(c.Request.Context(), converter.LoginRequestToUserDomain(loginRequest))
	if err != nil {
		logger.Error("Error when trying call service", err, stacktraceLoginUserHandler)

		c.JSON(err.HttpStatusCode, err)
		return
	}

	logger.Info("User was logged Successfully", zap.String("user_email", loginRequest.Email), stacktraceLoginUserHandler)
	c.JSON(http.StatusOK, dtos.LoginResponse{Token: jwtToken})
}
