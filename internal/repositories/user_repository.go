package repositories

import (
	"context"

	"github.com/WalterPaes/go-rest-api-crud/internal/domain"
	"github.com/WalterPaes/go-rest-api-crud/internal/repositories/entities/converter"
	"github.com/WalterPaes/go-rest-api-crud/pkg/logger"
	resterrors "github.com/WalterPaes/go-rest-api-crud/pkg/rest_errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

const (
	ErrInsertData = "Error When Trying Insert Data"
)

var (
	stacktraceCreateUserRepository = zap.String("stacktrace", "create-user-repository")
)

type UserRepository interface {
	CreateUser(context.Context, *domain.User) (*domain.User, *resterrors.RestErr)
}

type userRepo struct {
	collection *mongo.Collection
}

func NewUserRepository(collection *mongo.Collection) *userRepo {
	return &userRepo{
		collection: collection,
	}
}

func (us *userRepo) CreateUser(parentCtx context.Context, userDomain *domain.User) (*domain.User, *resterrors.RestErr) {
	logger.Info("Starting Create User Repository", stacktraceCreateUserRepository)

	ctx, cancel := context.WithCancel(parentCtx)
	defer cancel()

	userEntity := converter.UserDomainToUserEntity(userDomain)

	res, err := us.collection.InsertOne(ctx, userEntity)
	if err != nil {
		logger.Error(ErrInsertData, err, zap.String("stacktrace", "create-user"))
		return nil, resterrors.NewInternalServerError(ErrInsertData)
	}

	logger.Info("User was Created Successfully", stacktraceCreateUserRepository)

	userEntity.ID = res.InsertedID.(primitive.ObjectID)
	return converter.UserEntityToUserDomain(*userEntity), nil
}
