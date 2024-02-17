package handlers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/WalterPaes/go-rest-api-crud/internal/domain"
	"github.com/WalterPaes/go-rest-api-crud/internal/services"
	"github.com/WalterPaes/go-rest-api-crud/mocks"
	"github.com/WalterPaes/go-rest-api-crud/pkg/logger"
	resterrors "github.com/WalterPaes/go-rest-api-crud/pkg/rest_errors"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func init() {
	logger.Init("debug", "stdout")
}

var (
	ctx = context.Background()
)

func getContext(recorder *httptest.ResponseRecorder) *gin.Context {
	gin.SetMode(gin.TestMode)

	c, _ := gin.CreateTestContext(recorder)
	c.Request = &http.Request{
		Header: make(http.Header),
		URL:    &url.URL{},
	}

	return c
}

func getUserHandler(userService services.UserService) *userHandler {
	return NewUserHandler(userService)
}

func getUserServiceFindAll(t *testing.T, usersList []*domain.User, err error) services.UserService {
	t.Helper()
	m := mocks.NewUserService(t)
	m.On("FindAll", ctx, itemsPerPage, currentPage).
		Return(usersList, err)
	return m
}

func Test_userHandler_ListAll(t *testing.T) {
	t.Run("Should return a list of users when not send pagination", func(t *testing.T) {
		userService := getUserServiceFindAll(t, []*domain.User{}, nil)
		userHandler := getUserHandler(userService)

		recorder := httptest.NewRecorder()
		ctx := getContext(recorder)

		ctx.Request.Method = http.MethodGet
		ctx.Request.Header.Set("Content-Type", "application/json")

		url := &url.Values{}

		ctx.Request.URL.RawQuery = url.Encode()

		userHandler.ListAll(ctx)

		assert.EqualValues(t, http.StatusOK, recorder.Code)
	})

	t.Run("Should return an error when try call user service", func(t *testing.T) {
		userService := getUserServiceFindAll(t, nil, resterrors.NewInternalServerError("error"))
		userHandler := getUserHandler(userService)

		recorder := httptest.NewRecorder()
		ctx := getContext(recorder)

		ctx.Request.Method = http.MethodGet
		ctx.Request.Header.Set("Content-Type", "application/json")

		url := &url.Values{}

		ctx.Request.URL.RawQuery = url.Encode()

		userHandler.ListAll(ctx)

		assert.EqualValues(t, http.StatusInternalServerError, recorder.Code)
	})

	t.Run("Should return a list of users when not send pagination", func(t *testing.T) {
		userService := getUserServiceFindAll(t, []*domain.User{}, nil)
		userHandler := getUserHandler(userService)

		recorder := httptest.NewRecorder()
		ctx := getContext(recorder)

		ctx.Request.Method = http.MethodGet
		ctx.Request.Header.Set("Content-Type", "application/json")

		url := &url.Values{}

		ctx.Request.URL.RawQuery = url.Encode()

		userHandler.ListAll(ctx)

		assert.EqualValues(t, http.StatusOK, recorder.Code)
	})

	t.Run("Should return a list of users when send pagination", func(t *testing.T) {
		userService := getUserServiceFindAll(t, []*domain.User{}, nil)
		userHandler := getUserHandler(userService)

		recorder := httptest.NewRecorder()
		ctx := getContext(recorder)

		ctx.Request.Method = http.MethodGet
		ctx.Request.Header.Set("Content-Type", "application/json")

		url := &url.Values{}
		url.Add("page", "1")
		url.Add("per_page", "10")

		ctx.Request.URL.RawQuery = url.Encode()

		userHandler.ListAll(ctx)

		assert.EqualValues(t, http.StatusOK, recorder.Code)
	})

	t.Run("Should return an error when send not int values to 'per_page' query param", func(t *testing.T) {
		userService := mocks.NewUserService(t)
		userHandler := getUserHandler(userService)

		recorder := httptest.NewRecorder()
		ctx := getContext(recorder)

		ctx.Request.Method = http.MethodGet
		ctx.Request.Header.Set("Content-Type", "application/json")

		url := &url.Values{}
		url.Add("page", "1")
		url.Add("per_page", "invalid")

		ctx.Request.URL.RawQuery = url.Encode()

		userHandler.ListAll(ctx)

		assert.EqualValues(t, http.StatusBadRequest, recorder.Code)
	})

	t.Run("Should return an error when send not int values to 'page' query param", func(t *testing.T) {
		userService := mocks.NewUserService(t)
		userHandler := getUserHandler(userService)

		recorder := httptest.NewRecorder()
		ctx := getContext(recorder)

		ctx.Request.Method = http.MethodGet
		ctx.Request.Header.Set("Content-Type", "application/json")

		url := &url.Values{}
		url.Add("page", "invalid")
		url.Add("per_page", "10")

		ctx.Request.URL.RawQuery = url.Encode()

		userHandler.ListAll(ctx)

		assert.EqualValues(t, http.StatusBadRequest, recorder.Code)
	})
}
