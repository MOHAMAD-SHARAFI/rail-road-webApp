package main

import (
	"payment-service/internal/config"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	config.LoadConfig()
	addr := viper.GetString("addr")
	dsn := viper.GetString("dsn")
	//port := viper.GetInt("port")

	_, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	r.Run(addr)
}
