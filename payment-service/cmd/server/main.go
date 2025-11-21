package main

import (
	"log"
	"payment-service/internal/bootstrap"
	"payment-service/internal/config"
	"payment-service/internal/logger"
	"payment-service/internal/models"
	"payment-service/internal/repositories"

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
	bootstrap.Migration(db)
	paymentRepo := repositories.NewPaymentRepository(db)

	logger.InitLogger()

	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	r.Run(addr)
}
