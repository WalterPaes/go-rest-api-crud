package repositories

import (
	"github.com/WalterPaes/go-rest-api-crud/internal/repositories/entities"
	"github.com/WalterPaes/go-rest-api-crud/pkg/logger"
	"go.uber.org/zap"
)

type UserRepository interface{}

type userRepo struct{}

func NewUserRepository() *userRepo {
	return &userRepo{}
}

func (us *userRepo) CreateUser(user entities.UserEntity) (*entities.UserEntity, error) {
	logger.Info("Starting Create User Repository", zap.String("stacktrace", "create-user"))

	return nil, nil
}
