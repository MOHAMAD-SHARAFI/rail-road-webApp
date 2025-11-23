package main

import (
	"log"
	"payment-service/internal/bootstrap"
	"payment-service/internal/config"
	"payment-service/internal/models"

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
	err = bootstrap.Migration(db)
	if err != nil {
		log.Fatalln(models.ErrFailedMigrateDB)
	}

	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	r.Run(addr)
}
