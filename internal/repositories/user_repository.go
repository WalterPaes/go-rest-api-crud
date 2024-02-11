package repositories

import (
	"context"

	"github.com/WalterPaes/go-rest-api-crud/internal/domain"
	"github.com/WalterPaes/go-rest-api-crud/internal/repositories/entities/converter"
	"github.com/WalterPaes/go-rest-api-crud/pkg/logger"
	resterrors "github.com/WalterPaes/go-rest-api-crud/pkg/rest_errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

const (
	errInsertUser = "Error When Trying Insert User"
	errUpdateUser = "Error When Trying Update User"
	errDeleteUser = "Error When Trying Delete User"
)

var (
	stacktraceCreateUserRepository = zap.String("stacktrace", "create-user-repository")
	stacktraceUpdateUserRepository = zap.String("stacktrace", "update-user-repository")
	stacktraceDeleteUserRepository = zap.String("stacktrace", "delete-user-repository")
)

type UserRepository interface {
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

func (us *userRepo) UpdateUser(parentCtx context.Context, userID string, userDomain *domain.User) (*domain.User, *resterrors.RestErr) {
	logger.Info("Starting Update User Repository", stacktraceUpdateUserRepository)

	ctx, cancel := context.WithCancel(parentCtx)
	defer cancel()

	userObjectId, _ := primitive.ObjectIDFromHex(userID)
	userEntity := converter.UserDomainToUserEntity(userDomain)

	filter := bson.D{{Key: "_id", Value: userObjectId}}
	updateData := bson.D{{Key: "$set", Value: userEntity}}

	_, err := us.collection.UpdateOne(ctx, filter, updateData)
	if err != nil {
		logger.Error(errUpdateUser, err, stacktraceUpdateUserRepository)
		return nil, resterrors.NewInternalServerError(errUpdateUser)
	}
	userEntity.ID = userObjectId

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
		logger.Error(errDeleteUser, err, stacktraceCreateUserRepository)
		return resterrors.NewInternalServerError(errDeleteUser)
	}

	logger.Info("User was Delete Successfully", zap.String("userId", userID), stacktraceDeleteUserRepository)
	return nil
}
