package repositories

import (
	"context"
	"net/http"
	"testing"

	"github.com/WalterPaes/go-rest-api-crud/internal/domain"
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
)

func init() {
	logger.Init("debug", "stdout")
}

func Test_userRepo_CreateUser(t *testing.T) {
	mtestDB := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mtestDB.Run("Should Create an User Successfully", func(mtestDB *mtest.T) {
		mtestDB.AddMockResponses(bson.D{
			{Key: "ok", Value: 1},
			{Key: "value", Value: bson.D{
				{Key: "_id", Value: primitive.NewObjectID()},
				{Key: "name", Value: user.Name},
				{Key: "email", Value: user.Email},
				{Key: "password", Value: user.Password},
			}},
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
		assert.Equal(t, err.Message, ErrInsertData)
	})
}
