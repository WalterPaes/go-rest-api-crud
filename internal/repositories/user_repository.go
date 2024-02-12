package repositories

import (
	"context"
	"fmt"

	"github.com/WalterPaes/go-rest-api-crud/internal/domain"
	"github.com/WalterPaes/go-rest-api-crud/internal/repositories/entities"
	"github.com/WalterPaes/go-rest-api-crud/internal/repositories/entities/converter"
	"github.com/WalterPaes/go-rest-api-crud/pkg/logger"
	resterrors "github.com/WalterPaes/go-rest-api-crud/pkg/rest_errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

const (
	errFindAllUsers    = "Error When Try Find All Users"
	errFindByIdUser    = "Error When Try Find User By ID"
	errFindByEmailUser = "Error When Try Find User By Email"
	errInsertUser      = "Error When Try Insert User"
	errUpdateUser      = "Error When Try Update User"
	errDeleteUser      = "Error When Try Delete User"
)

var (
	stacktraceFindAllUserRepository     = zap.String("stacktrace", "find-all-user-repository")
	stacktraceFindUserByIdRepository    = zap.String("stacktrace", "find-user-by-id-repository")
	stacktraceFindUserByEmailRepository = zap.String("stacktrace", "find-user-by-email-repository")
	stacktraceCreateUserRepository      = zap.String("stacktrace", "create-user-repository")
	stacktraceUpdateUserRepository      = zap.String("stacktrace", "update-user-repository")
	stacktraceDeleteUserRepository      = zap.String("stacktrace", "delete-user-repository")
)

type UserRepository interface {
	FindUserById(parentCtx context.Context, userID string) (*domain.User, *resterrors.RestErr)
	FindUserByEmail(parentCtx context.Context, email string) (*domain.User, *resterrors.RestErr)
	CreateUser(context.Context, *domain.User) (*domain.User, *resterrors.RestErr)
	UpdateUser(ctx context.Context, userID string, user *domain.User) (*domain.User, *resterrors.RestErr)
	DeleteUser(ctx context.Context, userID string) *resterrors.RestErr
	FindAll(parentCtx context.Context, itemsPerPage, currentPage int) ([]*domain.User, *resterrors.RestErr)
}

type userRepo struct {
	collection *mongo.Collection
}

func NewUserRepository(client *mongo.Client, databaseName, collectionName string) *userRepo {
	return &userRepo{
		collection: client.Database(databaseName).Collection(collectionName),
	}
}

func (us *userRepo) FindAll(parentCtx context.Context, itemsPerPage, currentPage int) ([]*domain.User, *resterrors.RestErr) {
	logger.Info("Starting Find All Users", stacktraceFindAllUserRepository)

	ctx, cancel := context.WithCancel(parentCtx)
	defer cancel()

	limit := int64(itemsPerPage)
	skip := int64(currentPage*itemsPerPage - itemsPerPage)

	curr, err := us.collection.Find(ctx, bson.D{{}}, &options.FindOptions{Limit: &limit, Skip: &skip})
	if err != nil {
		logger.Error(errFindAllUsers, err, stacktraceFindAllUserRepository)
		return nil, resterrors.NewInternalServerError(errFindAllUsers)
	}

	var usersList []*domain.User

	for curr.Next(ctx) {
		var user *entities.UserEntity
		if err := curr.Decode(&user); err != nil {
			logger.Error("Error when try decode user", err, stacktraceFindAllUserRepository)
		}
		usersList = append(usersList, converter.UserEntityToUserDomain(*user))
	}

	logger.Info(
		"List Users Successfully",
		zap.Int("items_per_page", itemsPerPage),
		zap.Int("current_page", currentPage),
		stacktraceFindUserByIdRepository,
	)
	return usersList, nil
}

func (us *userRepo) CreateUser(parentCtx context.Context, userDomain *domain.User) (*domain.User, *resterrors.RestErr) {
	logger.Info("Starting Create User", stacktraceCreateUserRepository)

	ctx, cancel := context.WithCancel(parentCtx)
	defer cancel()

	userEntity := converter.UserDomainToUserEntity(userDomain)

	res, err := us.collection.InsertOne(ctx, userEntity)
	if err != nil {
		logger.Error(errInsertUser, err, stacktraceCreateUserRepository)
		return nil, resterrors.NewInternalServerError(errInsertUser)
	}
	userEntity.ID = res.InsertedID.(primitive.ObjectID)

	logger.Info("User Created Successfully", zap.String("user_id", userEntity.ID.Hex()), stacktraceCreateUserRepository)

	return converter.UserEntityToUserDomain(*userEntity), nil
}

func (us *userRepo) FindUserById(parentCtx context.Context, userID string) (*domain.User, *resterrors.RestErr) {
	logger.Info("Starting Find User by Id", stacktraceFindUserByIdRepository)

	ctx, cancel := context.WithCancel(parentCtx)
	defer cancel()

	userObjectId, _ := primitive.ObjectIDFromHex(userID)

	userEntity := &entities.UserEntity{}

	err := us.collection.FindOne(ctx, bson.D{{Key: "_id", Value: userObjectId}}).Decode(userEntity)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			errorMsg := fmt.Sprintf("No users found with this id: %s", userID)
			logger.Error(errorMsg, err, stacktraceFindUserByIdRepository)
			return nil, resterrors.NewNotFoundError(errorMsg)
		}

		logger.Error(errFindByIdUser, err, stacktraceFindUserByIdRepository)
		return nil, resterrors.NewInternalServerError(errFindByIdUser)
	}

	logger.Info("User Found Successfully", zap.String("user_id", userID), stacktraceFindUserByIdRepository)
	return converter.UserEntityToUserDomain(*userEntity), nil
}

func (us *userRepo) FindUserByEmail(parentCtx context.Context, email string) (*domain.User, *resterrors.RestErr) {
	logger.Info("Starting Find User by Email", stacktraceFindUserByEmailRepository)

	ctx, cancel := context.WithCancel(parentCtx)
	defer cancel()

	userEntity := &entities.UserEntity{}

	err := us.collection.FindOne(ctx, bson.D{{Key: "email", Value: email}}).Decode(userEntity)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			errorMsg := fmt.Sprintf("No users found with this email: %s", email)
			logger.Error(errorMsg, err, stacktraceFindUserByEmailRepository)
			return nil, resterrors.NewNotFoundError(errorMsg)
		}

		logger.Error(errFindByEmailUser, err, stacktraceFindUserByEmailRepository)
		return nil, resterrors.NewInternalServerError(errFindByEmailUser)
	}

	logger.Info("User Found Successfully", zap.String("user_email", userEntity.Email), stacktraceFindUserByEmailRepository)
	return converter.UserEntityToUserDomain(*userEntity), nil
}

func (us *userRepo) UpdateUser(parentCtx context.Context, userID string, userDomain *domain.User) (*domain.User, *resterrors.RestErr) {
	logger.Info("Starting Update User", stacktraceUpdateUserRepository)

	ctx, cancel := context.WithCancel(parentCtx)
	defer cancel()

	userObjectId, _ := primitive.ObjectIDFromHex(userID)
	userEntity := &entities.UserEntity{}

	filter := bson.D{{Key: "_id", Value: userObjectId}}
	updateData := bson.D{{Key: "$set", Value: converter.UserDomainToUserEntity(userDomain)}}

	err := us.collection.FindOneAndUpdate(ctx, filter, updateData).Decode(userEntity)
	if err != nil {
		logger.Error(errUpdateUser, err, stacktraceUpdateUserRepository)
		return nil, resterrors.NewInternalServerError(errUpdateUser)
	}

	logger.Info("User Updated Successfully", zap.String("user_id", userEntity.ID.Hex()), stacktraceUpdateUserRepository)

	return converter.UserEntityToUserDomain(*userEntity), nil
}

func (us *userRepo) DeleteUser(parentCtx context.Context, userID string) *resterrors.RestErr {
	logger.Info("Starting Delete User", stacktraceDeleteUserRepository)

	ctx, cancel := context.WithCancel(parentCtx)
	defer cancel()

	userObjectId, _ := primitive.ObjectIDFromHex(userID)

	_, err := us.collection.DeleteOne(ctx, bson.D{{Key: "_id", Value: userObjectId}})
	if err != nil {
		logger.Error(errDeleteUser, err, stacktraceDeleteUserRepository)
		return resterrors.NewInternalServerError(errDeleteUser)
	}

	logger.Info("User Delete Successfully", zap.String("user_id", userID), stacktraceDeleteUserRepository)
	return nil
}
