package services

import (
	"context"
	"net/http"

	"github.com/WalterPaes/go-rest-api-crud/internal/domain"
	"github.com/WalterPaes/go-rest-api-crud/internal/repositories"
	"github.com/WalterPaes/go-rest-api-crud/pkg/jwt"
	"github.com/WalterPaes/go-rest-api-crud/pkg/logger"
	resterrors "github.com/WalterPaes/go-rest-api-crud/pkg/rest_errors"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

const errInvalidCredentials = "Credentials are Invalid"

var (
	stacktraceLoginService = zap.String("stacktrace", "login-service")
)

type LoginService interface {
	LoginUser(ctx context.Context, user *domain.User) (string, *resterrors.RestErr)
}

type loginSvc struct {
	userRepository repositories.UserRepository
	jwtAuth        jwt.JwtAuth
}

func NewLoginService(userRepository repositories.UserRepository, jwtAuth jwt.JwtAuth) *loginSvc {
	return &loginSvc{
		userRepository: userRepository,
		jwtAuth:        jwtAuth,
	}
}

func (s *loginSvc) LoginUser(ctx context.Context, user *domain.User) (string, *resterrors.RestErr) {
	logger.Info("Starting Login User", stacktraceLoginService)

	resultUser, err := s.userRepository.FindUserByEmail(ctx, user.Email)
	if err != nil {
		if err.HttpStatusCode == http.StatusNotFound {
			logger.Error(errInvalidCredentials, err, stacktraceLoginService)
			return "", resterrors.NewUnauthorizedError(errInvalidCredentials)
		}

		logger.Error("Error when trying call repository", err, stacktraceLoginService)
		return "", err
	}

	if !s.validatePassword(user.Password, resultUser.Password) {
		logger.Error(errInvalidCredentials, err, stacktraceLoginService)
		return "", resterrors.NewUnauthorizedError(errInvalidCredentials)
	}

	token, err := s.jwtAuth.GenerateToken(map[string]any{
		"id":    resultUser.ID,
		"email": resultUser.Email,
		"name":  resultUser.Name,
	})
	if err != nil {
		logger.Error(err.Error(), err, stacktraceLoginService)
		return "", err
	}

	logger.Info("User was logged successfully", zap.String("user_id", resultUser.ID), stacktraceLoginService)
	return token, nil
}

func (s *loginSvc) validatePassword(password, userPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(userPassword), []byte(password))
	return err == nil
}
