package services

import (
	"context"

	"github.com/WalterPaes/go-rest-api-crud/internal/domain"
	"github.com/WalterPaes/go-rest-api-crud/internal/repositories"
	"github.com/WalterPaes/go-rest-api-crud/pkg/logger"
	resterrors "github.com/WalterPaes/go-rest-api-crud/pkg/rest_errors"
	"go.uber.org/zap"
)

var (
	stacktraceCreateUserService = zap.String("stacktrace", "create-user-service")
)

type UserService interface {
	CreateUser(context.Context, *domain.User) (*domain.User, *resterrors.RestErr)
}

type userSvc struct {
	userRepository repositories.UserRepository
}

func NewUserService(userRepository repositories.UserRepository) *userSvc {
	return &userSvc{
		userRepository: userRepository,
	}
}

func (s *userSvc) CreateUser(ctx context.Context, user *domain.User) (*domain.User, *resterrors.RestErr) {
	logger.Info("Starting Create User", stacktraceCreateUserService)

	// Encrypt Password
	createdUser, err := s.userRepository.CreateUser(ctx, user)
	if err != nil {
		logger.Error("Error when trying call repository", err, stacktraceCreateUserService)
		return nil, err
	}

	logger.Info("CreateUser service executed successfully", zap.String("user_id", createdUser.ID), stacktraceCreateUserService)
	return createdUser, nil
}
