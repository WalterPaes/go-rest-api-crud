package repositories

import (
	"context"
	"net/http"
	"testing"

	"github.com/WalterPaes/go-rest-api-crud/internal/domain"
	"github.com/WalterPaes/go-rest-api-crud/internal/repositories/entities"
	"github.com/WalterPaes/go-rest-api-crud/pkg/logger"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
)

const (
	dbName         = "users"
	collectionName = "users"
)

var (
	ctx = context.Background()

	user = &domain.User{
		Name:     "Teste",
		Email:    "teste@email.com",
		Password: "abcdef",
	}

	userEntity = &entities.UserEntity{
		ID:       primitive.NewObjectID(),
		Name:     user.Name,
		Email:    user.Email,
		Password: user.Password,
	}
)

func init() {
	logger.Init("debug", "stdout")
}

func Test_userRepo_CreateUser(t *testing.T) {
	mtestDB := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mtestDB.Run("Should Create an User Successfully", func(mtestDB *mtest.T) {
		mtestDB.AddMockResponses(bson.D{
			{Key: "ok", Value: 1},
			{Key: "value", Value: userEntity},
		})

		userRepository := NewUserRepository(mtestDB.Client, dbName, collectionName)

		result, err := userRepository.CreateUser(ctx, user)

		assert.Nil(t, err)
		assert.Equal(t, result.Name, user.Name)
		assert.Equal(t, result.Email, user.Email)
	})

	mtestDB.Run("Should return an error when try create an user", func(mtestDB *mtest.T) {
		mtestDB.AddMockResponses(bson.D{
			{Key: "ok", Value: 0},
		})

		userRepository := NewUserRepository(mtestDB.Client, dbName, collectionName)

		result, err := userRepository.CreateUser(ctx, user)

		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.Equal(t, err.HttpStatusCode, http.StatusInternalServerError)
		assert.Equal(t, err.Message, errInsertUser)
	})
}

func Test_userRepo_UpdateUser(t *testing.T) {
	mtestDB := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mtestDB.Run("Should Update an User Successfully", func(mtestDB *mtest.T) {
		mtestDB.AddMockResponses(bson.D{
			{Key: "ok", Value: 1},
			{Key: "value", Value: userEntity},
		})

		userRepository := NewUserRepository(mtestDB.Client, dbName, collectionName)

		result, err := userRepository.UpdateUser(ctx, userEntity.ID.Hex(), user)

		assert.Nil(t, err)
		assert.Equal(t, result.Name, user.Name)
		assert.Equal(t, result.Email, user.Email)
	})

	mtestDB.Run("Should return an error when try update an user", func(mtestDB *mtest.T) {
		mtestDB.AddMockResponses(bson.D{
			{Key: "ok", Value: 0},
		})

		userRepository := NewUserRepository(mtestDB.Client, dbName, collectionName)

		result, err := userRepository.UpdateUser(ctx, userEntity.ID.Hex(), user)

		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.Equal(t, err.HttpStatusCode, http.StatusInternalServerError)
		assert.Equal(t, err.Message, errUpdateUser)
	})
}

func Test_userRepo_DeleteUser(t *testing.T) {
	mtestDB := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mtestDB.Run("Should Delete an User Successfully", func(mtestDB *mtest.T) {
		mtestDB.AddMockResponses(bson.D{
			{Key: "ok", Value: 1},
			{Key: "n", Value: 1},
			{Key: "acknowledged", Value: true},
		})

		userRepository := NewUserRepository(mtestDB.Client, dbName, collectionName)

		err := userRepository.DeleteUser(ctx, userEntity.ID.Hex())

		assert.Nil(t, err)
	})

	mtestDB.Run("Should return an error when try delete an user", func(mtestDB *mtest.T) {
		mtestDB.AddMockResponses(bson.D{
			{Key: "ok", Value: 0},
		})

		userRepository := NewUserRepository(mtestDB.Client, dbName, collectionName)

		err := userRepository.DeleteUser(ctx, userEntity.ID.Hex())

		assert.NotNil(t, err)
		assert.Equal(t, err.HttpStatusCode, http.StatusInternalServerError)
		assert.Equal(t, err.Message, errDeleteUser)
	})
}
