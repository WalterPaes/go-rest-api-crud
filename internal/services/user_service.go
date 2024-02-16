package services

import (
	"context"
	"errors"
	"net/http"

	"github.com/WalterPaes/go-rest-api-crud/internal/domain"
	"github.com/WalterPaes/go-rest-api-crud/internal/repositories"
	"github.com/WalterPaes/go-rest-api-crud/pkg/logger"
	resterrors "github.com/WalterPaes/go-rest-api-crud/pkg/rest_errors"
	"go.uber.org/zap"
)

const (
	errCallRepositoy          = "Error when try call repository"
	errEmailAlreadyRegistered = "Email is already registered"
)

var (
	stacktraceFindAllUsersService = zap.String("stacktrace", "find-all-users-service")
	stacktraceCreateUserService   = zap.String("stacktrace", "create-user-service")
	stacktraceFindUserByIdService = zap.String("stacktrace", "find-user-by-id-service")
	stacktraceUpdateUserService   = zap.String("stacktrace", "update-user-service")
	stacktraceDeleteUserService   = zap.String("stacktrace", "delete-user-service")
)

type UserService interface {
	FindAll(ctx context.Context, itemsPerPage, currentPage int) ([]*domain.User, *resterrors.RestErr)
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

func (s *userSvc) FindAll(ctx context.Context, itemsPerPage, currentPage int) ([]*domain.User, *resterrors.RestErr) {
	logger.Info("Starting FindAll", stacktraceFindAllUsersService)

	users, err := s.userRepository.FindAll(ctx, itemsPerPage, currentPage)
	if err != nil {
		logger.Error(errCallRepositoy, err, stacktraceFindAllUsersService)
		return nil, err
	}

	logger.Info(
		"FindAll executed successfully",
		zap.Int("items_per_page", itemsPerPage),
		zap.Int("current_page", currentPage),
		stacktraceFindAllUsersService,
	)
	return users, nil
}

func (s *userSvc) CreateUser(ctx context.Context, user *domain.User) (*domain.User, *resterrors.RestErr) {
	logger.Info("Starting CreateUser", stacktraceCreateUserService)

	if err := s.checkIfEmailIsAlreadyRegistered(ctx, user.Email, user.ID); err != nil {
		return nil, err
	}

	createdUser, err := s.userRepository.CreateUser(ctx, user)
	if err != nil {
		logger.Error(errCallRepositoy, err, stacktraceCreateUserService)
		return nil, err
	}

	logger.Info("CreateUser executed successfully", zap.String("user_id", createdUser.ID), stacktraceCreateUserService)
	return createdUser, nil
}

func (s *userSvc) FindUserById(ctx context.Context, userID string) (*domain.User, *resterrors.RestErr) {
	logger.Info("Starting FindUserById", stacktraceFindUserByIdService)

	user, err := s.userRepository.FindUserById(ctx, userID)
	if err != nil {
		logger.Error(errCallRepositoy, err, stacktraceFindUserByIdService)
		return nil, err
	}

	logger.Info("FindUserById executed successfully", zap.String("user_id", user.ID), stacktraceFindUserByIdService)
	return user, nil
}

func (s *userSvc) UpdateUser(ctx context.Context, userID string, user *domain.User) (*domain.User, *resterrors.RestErr) {
	logger.Info("Starting UpdateUser", stacktraceUpdateUserService)

	if err := s.checkIfEmailIsAlreadyRegistered(ctx, user.Email, userID); err != nil {
		return nil, err
	}

	updatedUser, err := s.userRepository.UpdateUser(ctx, userID, user)
	if err != nil {
		logger.Error(errCallRepositoy, err, stacktraceUpdateUserService)
		return nil, err
	}

	logger.Info("UpdateUser executed successfully", zap.String("user_id", updatedUser.ID), stacktraceUpdateUserService)
	return updatedUser, nil
}

func (s *userSvc) DeleteUser(ctx context.Context, userID string) *resterrors.RestErr {
	logger.Info("Starting DeleteUser", stacktraceDeleteUserService)

	err := s.userRepository.DeleteUser(ctx, userID)
	if err != nil {
		logger.Error(errCallRepositoy, err, stacktraceDeleteUserService)
		return err
	}

	logger.Info("DeleteUser executed successfully", zap.String("user_id", userID), stacktraceDeleteUserService)
	return nil
}

func (s *userSvc) checkIfEmailIsAlreadyRegistered(ctx context.Context, email, userID string) *resterrors.RestErr {
	resultUser, err := s.userRepository.FindUserByEmail(ctx, email)
	if err != nil && err.HttpStatusCode != http.StatusNotFound {
		return err
	}

	if resultUser != nil && resultUser.ID != userID {
		errMsg := errors.New(errEmailAlreadyRegistered)
		logger.Error(errMsg.Error(), errMsg)
		return resterrors.NewBadRequestError(errMsg.Error())
	}

	return nil
}
