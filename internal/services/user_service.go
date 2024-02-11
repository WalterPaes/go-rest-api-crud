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
	stacktraceDeleteUserService = zap.String("stacktrace", "delete-user-service")
)

type UserService interface {
	CreateUser(context.Context, *domain.User) (*domain.User, *resterrors.RestErr)
	DeleteUser(context.Context, string) *resterrors.RestErr
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

	if err := user.EncryptPassword(); err != nil {
		logger.Error("Error when trying Encrypt Password", err, stacktraceCreateUserService)
		return nil, resterrors.NewInternalServerError("Error when try create user")
	}

	createdUser, err := s.userRepository.CreateUser(ctx, user)
	if err != nil {
		logger.Error("Error when trying call repository", err, stacktraceCreateUserService)
		return nil, err
	}

	logger.Info("CreateUser service executed successfully", zap.String("user_id", createdUser.ID), stacktraceCreateUserService)
	return createdUser, nil
}

func (s *userSvc) DeleteUser(ctx context.Context, userID string) *resterrors.RestErr {
	logger.Info("Starting Delete User", stacktraceDeleteUserService)

	err := s.userRepository.DeleteUser(ctx, userID)
	if err != nil {
		logger.Error("Error when trying call repository", err, stacktraceDeleteUserService)
		return err
	}

	logger.Info("DeleteUser service executed successfully", zap.String("user_id", userID), stacktraceDeleteUserService)
	return nil
}
