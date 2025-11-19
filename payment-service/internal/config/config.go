package config

import (
	"log"

	"github.com/spf13/viper"
)

func LoadConfig() {
	viper.AddConfigPath("~/rail-road-webApp/payment-service/internal/config")
	viper.SetConfigName("app")
	viper.SetConfigType("env")
	viper.AutomaticEnv()
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}
}
