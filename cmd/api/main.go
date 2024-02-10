package main

import (
	"log"
	"net/http"

	"github.com/WalterPaes/go-rest-api-crud/configs"
	"github.com/WalterPaes/go-rest-api-crud/pkg/logger"
	"github.com/gin-gonic/gin"
)

func main() {
	cfg, err := configs.Load()
	if err != nil {
		log.Fatal(err)
	}

	logger := logger.NewLogger(cfg.LogLevel, cfg.LogOutput)
	logger.Info("Start Application")

	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	r.Run(cfg.ApiPort)
}
