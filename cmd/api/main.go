package main

import (
	"context"
	"log"
	"net/http"

	"github.com/WalterPaes/go-rest-api-crud/configs"
	"github.com/WalterPaes/go-rest-api-crud/internal/handlers"
	"github.com/WalterPaes/go-rest-api-crud/internal/repositories"
	"github.com/WalterPaes/go-rest-api-crud/internal/services"
	"github.com/WalterPaes/go-rest-api-crud/pkg/logger"
	"github.com/WalterPaes/go-rest-api-crud/pkg/mongodb"
	"github.com/gin-gonic/gin"
)

func main() {
	cfg, err := configs.Load()
	if err != nil {
		log.Fatal(err)
	}

	logger.Init(cfg.LogLevel, cfg.LogOutput)
	logger.Info("Start Application")

	dbClient := mongodb.NewMongoDBClient(context.Background(), cfg.MongoDBTimeout, cfg.MongoDBUri)

	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	userRepository := repositories.NewUserRepository(dbClient, cfg.MongoDBDatabase, cfg.MongoDBCollection)
	userService := services.NewUserService(userRepository)
	userHandler := handlers.NewUserHandler(userService)

	r.POST("/users", userHandler.CreateUser)
	r.PUT("/users/:id", userHandler.UpdateUser)
	r.DELETE("/users/:id", userHandler.DeleteUser)

	r.Run(cfg.ApiPort)
}
