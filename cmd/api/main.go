package main

import (
	"context"
	"log"
	"net/http"

	"github.com/WalterPaes/go-rest-api-crud/configs"
	"github.com/WalterPaes/go-rest-api-crud/internal/handlers"
	"github.com/WalterPaes/go-rest-api-crud/internal/repositories"
	"github.com/WalterPaes/go-rest-api-crud/internal/services"
	"github.com/WalterPaes/go-rest-api-crud/pkg/jwt"
	"github.com/WalterPaes/go-rest-api-crud/pkg/logger"
	"github.com/WalterPaes/go-rest-api-crud/pkg/mongodb"
	"github.com/gin-gonic/gin"

	docs "github.com/WalterPaes/go-rest-api-crud/docs"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Go User's API
// @version 1.0
// @description User API with authentication
// @host localhost:8000
// @BasePath /
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	cfg, err := configs.Load()
	if err != nil {
		log.Fatal(err)
	}

	logger.Init(cfg.LogLevel, cfg.LogOutput)
	logger.Info("Start Application")

	jwtAuth := jwt.NewJwtAuth(cfg.JwtSecret, cfg.JwtExpTime)

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

	loginService := services.NewLoginService(userRepository, jwtAuth)
	loginHandler := handlers.NewLoginHandler(loginService)

	r.POST("/login", loginHandler.Login)

	r.GET("/users", jwtAuth.VerifyTokenMiddleware, userHandler.ListAll)
	r.POST("/users", jwtAuth.VerifyTokenMiddleware, userHandler.CreateUser)
	r.GET("/users/:id", jwtAuth.VerifyTokenMiddleware, userHandler.GetUserById)
	r.PUT("/users/:id", jwtAuth.VerifyTokenMiddleware, userHandler.UpdateUser)
	r.DELETE("/users/:id", jwtAuth.VerifyTokenMiddleware, userHandler.DeleteUser)

	docs.SwaggerInfo.BasePath = "/"
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	r.Run(cfg.ApiPort)
}
