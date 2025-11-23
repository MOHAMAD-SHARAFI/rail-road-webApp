package main

import (
	"log"
	"user-service/internal/bootstrap"
	"user-service/internal/config"

	"user-service/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func main() {
	config.LoadConfig()
	addr := viper.GetString("addr")
	//port := viper.GetInt("port")
	db, err := bootstrap.InitDB()
	if err != nil {
		log.Fatal("Failed to initialize database", err)
	}

	logger.InitLogger()

	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	err = r.Run(addr)
	if err != nil {
		log.Fatalf("cannot start server: %v", err)
	}
}
