package main

import (
	"log"
	"user-service/internal/bootstrap"
	"user-service/internal/config"
	"user-service/internal/logger"
	"user-service/internal/models"
	"user-service/internal/repositories"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func main() {
	config.LoadConfig()
	addr := viper.GetString("addr")
	//port := viper.GetInt("port")
	db, err := bootstrap.InitDB()
	if err != nil {
		log.Fatalln(models.ErrFailedConnectDB)
	}
	err = bootstrap.MigrateDB(db)
	if err != nil {
		log.Fatalln(models.ErrFailedMigrateDB)
	}
	userRepo := repositories.NewUSerRepository(db)

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
