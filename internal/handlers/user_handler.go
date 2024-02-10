package handlers

import "github.com/gin-gonic/gin"

type UserHandler interface {
	ListUsers(c *gin.Context)
	FindUserByID(c *gin.Context)
	CreateUser(c *gin.Context)
	UpdateUser(c *gin.Context)
	DeleteUser(c *gin.Context)
}

type userHandler struct{}

func NewUserHandler() *userHandler {
	return &userHandler{}
}
