package services

import (
	"context"
	"errors"

	"github.com/WalterPaes/go-rest-api-crud/internal/domain"
	"github.com/WalterPaes/go-rest-api-crud/internal/repositories"
	"github.com/WalterPaes/go-rest-api-crud/pkg/logger"
	resterrors "github.com/WalterPaes/go-rest-api-crud/pkg/rest_errors"
	"go.uber.org/zap"
)

var (
	stacktraceCreateUserService   = zap.String("stacktrace", "create-user-service")
	stacktraceFindUserByIdService = zap.String("stacktrace", "find-user-by-id-service")
	stacktraceUpdateUserService   = zap.String("stacktrace", "update-user-service")
	stacktraceDeleteUserService   = zap.String("stacktrace", "delete-user-service")
)

type UserService interface {
	CreateUser(context.Context, *domain.User) (*domain.User, *resterrors.RestErr)
	FindUserById(ctx context.Context, userID string) (*domain.User, *resterrors.RestErr)
	UpdateUser(ctx context.Context, userID string, user *domain.User) (*domain.User, *resterrors.RestErr)
	DeleteUser(ctx context.Context, userID string) *resterrors.RestErr
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

	if err := s.checkIfEmailIsAlreadyRegistered(ctx, user.Email, user.ID); err != nil {
		return nil, err
	}

	createdUser, err := s.userRepository.CreateUser(ctx, user)
	if err != nil {
		logger.Error("Error when trying call repository", err, stacktraceCreateUserService)
		return nil, err
	}

	logger.Info("CreateUser service executed successfully", zap.String("user_id", createdUser.ID), stacktraceCreateUserService)
	return createdUser, nil
}

func (s *userSvc) FindUserById(ctx context.Context, userID string) (*domain.User, *resterrors.RestErr) {
	logger.Info("Starting Find User By Id", stacktraceFindUserByIdService)

	user, err := s.userRepository.FindUserById(ctx, userID)
	if err != nil {
		logger.Error("Error when trying call repository", err, stacktraceFindUserByIdService)
		return nil, err
	}

	logger.Info("FindById service executed successfully", zap.String("user_id", user.ID), stacktraceFindUserByIdService)
	return user, nil
}

func (s *userSvc) UpdateUser(ctx context.Context, userID string, user *domain.User) (*domain.User, *resterrors.RestErr) {
	logger.Info("Starting Update User", stacktraceUpdateUserService)

	if err := s.checkIfEmailIsAlreadyRegistered(ctx, user.Email, userID); err != nil {
		return nil, err
	}

	updatedUser, err := s.userRepository.UpdateUser(ctx, userID, user)
	if err != nil {
		logger.Error("Error when trying call repository", err, stacktraceUpdateUserService)
		return nil, err
	}

	logger.Info("UpdateUser service executed successfully", zap.String("user_id", updatedUser.ID), stacktraceUpdateUserService)
	return updatedUser, nil
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

func (s *userSvc) checkIfEmailIsAlreadyRegistered(ctx context.Context, email, userID string) *resterrors.RestErr {
	resultUser, _ := s.userRepository.FindUserByEmail(ctx, email)

	if resultUser != nil && resultUser.ID != userID {
		errMsg := errors.New("email is already registered")
		logger.Error(errMsg.Error(), errMsg)
		return resterrors.NewBadRequestError(errMsg.Error())
	}

	return nil
}
