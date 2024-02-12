package repositories

import (
	"context"
	"fmt"
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

	currentPage  = 1
	itemsPerPage = 10
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

func Test_userRepo_FindUserById(t *testing.T) {
	mtestDB := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mtestDB.Run("Should Find an User By Id Successfully", func(mtestDB *mtest.T) {
		mtestDB.AddMockResponses(
			mtest.CreateCursorResponse(
				1,
				fmt.Sprintf("%s.%s", dbName, collectionName),
				mtest.FirstBatch,
				bson.D{
					{Key: "_id", Value: userEntity.ID},
					{Key: "email", Value: userEntity.Email},
					{Key: "password", Value: userEntity.Password},
					{Key: "name", Value: userEntity.Name},
				},
			),
		)

		userRepository := NewUserRepository(mtestDB.Client, dbName, collectionName)

		result, err := userRepository.FindUserById(ctx, userEntity.ID.Hex())

		assert.Nil(t, err)
		assert.Equal(t, result.Name, user.Name)
		assert.Equal(t, result.Email, user.Email)
	})

	mtestDB.Run("Should return an error when try Find an user by id", func(mtestDB *mtest.T) {
		mtestDB.AddMockResponses(bson.D{
			{Key: "ok", Value: 0},
		})

		userRepository := NewUserRepository(mtestDB.Client, dbName, collectionName)

		result, err := userRepository.FindUserById(ctx, userEntity.ID.Hex())

		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.Equal(t, err.HttpStatusCode, http.StatusInternalServerError)
		assert.Equal(t, err.Message, errFindByIdUser)
	})

	mtestDB.Run("Should return an error when user not found", func(mtestDB *mtest.T) {
		mtestDB.AddMockResponses(
			mtest.CreateCursorResponse(
				0,
				fmt.Sprintf("%s.%s", dbName, collectionName),
				mtest.FirstBatch,
			),
		)

		userRepository := NewUserRepository(mtestDB.Client, dbName, collectionName)

		result, err := userRepository.FindUserById(ctx, userEntity.ID.Hex())

		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.Equal(t, err.HttpStatusCode, http.StatusNotFound)
	})
}

func Test_userRepo_FindUserByEmail(t *testing.T) {
	mtestDB := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mtestDB.Run("Should Find an User By Email Successfully", func(mtestDB *mtest.T) {
		mtestDB.AddMockResponses(
			mtest.CreateCursorResponse(
				1,
				fmt.Sprintf("%s.%s", dbName, collectionName),
				mtest.FirstBatch,
				bson.D{
					{Key: "_id", Value: userEntity.ID},
					{Key: "email", Value: userEntity.Email},
					{Key: "password", Value: userEntity.Password},
					{Key: "name", Value: userEntity.Name},
				},
			),
		)

		userRepository := NewUserRepository(mtestDB.Client, dbName, collectionName)

		result, err := userRepository.FindUserByEmail(ctx, user.Email)

		assert.Nil(t, err)
		assert.Equal(t, result.Name, user.Name)
		assert.Equal(t, result.Email, user.Email)
	})

	mtestDB.Run("Should return an error when try Find an user by email", func(mtestDB *mtest.T) {
		mtestDB.AddMockResponses(bson.D{
			{Key: "ok", Value: 0},
		})

		userRepository := NewUserRepository(mtestDB.Client, dbName, collectionName)

		result, err := userRepository.FindUserByEmail(ctx, user.Email)

		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.Equal(t, err.HttpStatusCode, http.StatusInternalServerError)
		assert.Equal(t, err.Message, errFindByEmailUser)
	})

	mtestDB.Run("Should return an error when user not found", func(mtestDB *mtest.T) {
		mtestDB.AddMockResponses(
			mtest.CreateCursorResponse(
				0,
				fmt.Sprintf("%s.%s", dbName, collectionName),
				mtest.FirstBatch,
			),
		)

		userRepository := NewUserRepository(mtestDB.Client, dbName, collectionName)

		result, err := userRepository.FindUserByEmail(ctx, user.Email)

		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.Equal(t, err.HttpStatusCode, http.StatusNotFound)
	})
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

func Test_userRepo_FindAll(t *testing.T) {
	mtestDB := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mtestDB.Run("Should Find all Users Successfully", func(mtestDB *mtest.T) {
		first := mtest.CreateCursorResponse(
			1,
			fmt.Sprintf("%s.%s", dbName, collectionName),
			mtest.FirstBatch,
			bson.D{
				{Key: "_id", Value: userEntity.ID},
				{Key: "email", Value: userEntity.Email},
				{Key: "password", Value: userEntity.Password},
				{Key: "name", Value: userEntity.Name},
			},
		)

		getMore := mtest.CreateCursorResponse(
			1,
			fmt.Sprintf("%s.%s", dbName, collectionName),
			mtest.NextBatch,
			bson.D{
				{Key: "_id", Value: userEntity.ID},
				{Key: "email", Value: userEntity.Email},
				{Key: "password", Value: userEntity.Password},
				{Key: "name", Value: userEntity.Name},
			},
		)

		killCursors := mtest.CreateCursorResponse(
			0,
			fmt.Sprintf("%s.%s", dbName, collectionName),
			mtest.NextBatch,
		)

		mtestDB.AddMockResponses(first, getMore, killCursors)

		userRepository := NewUserRepository(mtestDB.Client, dbName, collectionName)

		result, err := userRepository.FindAll(ctx, itemsPerPage, currentPage)

		assert.Nil(t, err)
		assert.LessOrEqual(t, len(result), itemsPerPage)
	})

	mtestDB.Run("Should return an error when try find users", func(mtestDB *mtest.T) {
		mtestDB.AddMockResponses(bson.D{
			{Key: "ok", Value: 0},
		})

		userRepository := NewUserRepository(mtestDB.Client, dbName, collectionName)

		result, err := userRepository.FindAll(ctx, itemsPerPage, currentPage)

		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.Equal(t, err.HttpStatusCode, http.StatusInternalServerError)
		assert.Equal(t, err.Message, errFindAllUsers)
	})
}
