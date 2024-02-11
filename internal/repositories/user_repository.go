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
	"go.uber.org/zap"
)

const (
	errFindByIdUser = "Error When Trying Find User By ID"
	errInsertUser   = "Error When Trying Insert User"
	errUpdateUser   = "Error When Trying Update User"
	errDeleteUser   = "Error When Trying Delete User"
)

var (
	stacktraceFindUserByIdRepository = zap.String("stacktrace", "find-user-by-id-repository")
	stacktraceCreateUserRepository   = zap.String("stacktrace", "create-user-repository")
	stacktraceUpdateUserRepository   = zap.String("stacktrace", "update-user-repository")
	stacktraceDeleteUserRepository   = zap.String("stacktrace", "delete-user-repository")
)

type UserRepository interface {
	FindUserById(parentCtx context.Context, userID string) (*domain.User, *resterrors.RestErr)
	CreateUser(context.Context, *domain.User) (*domain.User, *resterrors.RestErr)
	UpdateUser(ctx context.Context, userID string, user *domain.User) (*domain.User, *resterrors.RestErr)
	DeleteUser(ctx context.Context, userID string) *resterrors.RestErr
}

type userRepo struct {
	collection *mongo.Collection
}

func NewUserRepository(client *mongo.Client, databaseName, collectionName string) *userRepo {
	return &userRepo{
		collection: client.Database(databaseName).Collection(collectionName),
	}
}

func (us *userRepo) CreateUser(parentCtx context.Context, userDomain *domain.User) (*domain.User, *resterrors.RestErr) {
	logger.Info("Starting Create User Repository", stacktraceCreateUserRepository)

	ctx, cancel := context.WithCancel(parentCtx)
	defer cancel()

	userEntity := converter.UserDomainToUserEntity(userDomain)

	res, err := us.collection.InsertOne(ctx, userEntity)
	if err != nil {
		logger.Error(errInsertUser, err, stacktraceCreateUserRepository)
		return nil, resterrors.NewInternalServerError(errInsertUser)
	}
	userEntity.ID = res.InsertedID.(primitive.ObjectID)

	logger.Info("User was Created Successfully", zap.String("userId", userEntity.ID.Hex()), stacktraceCreateUserRepository)

	return converter.UserEntityToUserDomain(*userEntity), nil
}

func (us *userRepo) FindUserById(parentCtx context.Context, userID string) (*domain.User, *resterrors.RestErr) {
	logger.Info("Starting Find User by Id Repository", stacktraceFindUserByIdRepository)

	ctx, cancel := context.WithCancel(parentCtx)
	defer cancel()

	userObjectId, _ := primitive.ObjectIDFromHex(userID)

	userEntity := &entities.UserEntity{}

	err := us.collection.FindOne(ctx, bson.D{{Key: "_id", Value: userObjectId}}).Decode(userEntity)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			errorMsg := fmt.Sprintf("User not found with this id: %s", userID)
			logger.Error(errorMsg, err, stacktraceFindUserByIdRepository)
			return nil, resterrors.NewNotFoundError(errorMsg)
		}

		logger.Error(errFindByIdUser, err, stacktraceFindUserByIdRepository)
		return nil, resterrors.NewInternalServerError(errFindByIdUser)
	}

	logger.Info("User was Found Successfully", zap.String("userId", userID), stacktraceFindUserByIdRepository)
	return converter.UserEntityToUserDomain(*userEntity), nil
}

func (us *userRepo) UpdateUser(parentCtx context.Context, userID string, userDomain *domain.User) (*domain.User, *resterrors.RestErr) {
	logger.Info("Starting Update User Repository", stacktraceUpdateUserRepository)

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

	logger.Info("User was Update Successfully", zap.String("userId", userEntity.ID.Hex()), stacktraceUpdateUserRepository)

	return converter.UserEntityToUserDomain(*userEntity), nil
}

func (us *userRepo) DeleteUser(parentCtx context.Context, userID string) *resterrors.RestErr {
	logger.Info("Starting Delete User Repository", stacktraceDeleteUserRepository)

	ctx, cancel := context.WithCancel(parentCtx)
	defer cancel()

	userObjectId, _ := primitive.ObjectIDFromHex(userID)

	_, err := us.collection.DeleteOne(ctx, bson.D{{Key: "_id", Value: userObjectId}})
	if err != nil {
		logger.Error(errDeleteUser, err, stacktraceDeleteUserRepository)
		return resterrors.NewInternalServerError(errDeleteUser)
	}

	logger.Info("User was Delete Successfully", zap.String("userId", userID), stacktraceDeleteUserRepository)
	return nil
}
