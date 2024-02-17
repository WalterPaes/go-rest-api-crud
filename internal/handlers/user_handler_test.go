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

func getUserServiceFindAllSuccess(t *testing.T) services.UserService {
	t.Helper()
	m := mocks.NewUserService(t)
	m.On("FindAll", ctx, itemsPerPage, currentPage).
		Return([]*domain.User{}, nil)
	return m
}

func Test_userHandler_ListAll(t *testing.T) {
	t.Run("Should return a list of users when not send pagination", func(t *testing.T) {
		userService := getUserServiceFindAllSuccess(t)
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
		userService := getUserServiceFindAllSuccess(t)
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
}
