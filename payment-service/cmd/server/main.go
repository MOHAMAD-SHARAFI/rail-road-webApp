package main

import (
	"log"
	"payment-service/internal/bootstrap"
	"payment-service/internal/config"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.LoadConfig()

	db, err := bootstrap.InitDB(cfg.DatabaseURL)
	if err != nil {
		log.Fatalln("Could not connect to database:", err)
	}

	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	err = r.Run(cfg.Port)
	if err != nil {
		return
	}
}
