package services

import (
	"github.com/WalterPaes/go-rest-api-crud/internal/domain"
	"github.com/WalterPaes/go-rest-api-crud/pkg/logger"
	"go.uber.org/zap"
)

type UserService interface{}

type userSvc struct{}

func NewUserService() *userSvc {
	return &userSvc{}
}

func (us *userSvc) CreateUser(user domain.User) (*domain.User, error) {
	logger.Info("Starting Create User Service", zap.String("stacktrace", "create-user"))

	return nil, nil
}
