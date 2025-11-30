package config

import (
	"os"
	"time"
)

type Config struct {
	DatabaseURL           string
	Port                  string
	UserServiceAddr       string
	RabbitMQURL           string
	OpenTelemetryEndpoint string
	WebSocketPort         string
	FeePercentage         float64
	MinFee                float64
}

func LoadConfig() *Config {
	return &Config{
		DatabaseURL: getenv("database_URL", "host=localhost user=postgres dbname=payment_service port=5432 sslmode=disable"),
		Port:        getenv("PORT", "8081"),
	}
}

func getenv(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}

func getdurationfromenv(key string, defaultval time.Duration) time.Duration {
	if val := os.Getenv(key); val != "" {
		if duration, err := time.ParseDuration(val); err == nil {
			return duration
		}
	}
	return defaultval
}
