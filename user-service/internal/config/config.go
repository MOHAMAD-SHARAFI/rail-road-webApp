package config

import (
	"os"
	"time"
)

type Config struct {
	DatabaseURL           string
	Port                  string
	JWTSecret             string
	RefreshTokenSecret    string
	AccessTokenExpiry     time.Duration
	RefreshTokenExpiry    time.Duration
	RabbitMQURL           string
	OpenTelemetryEndPoint string
	WebSocketPort         string
	EventDispatcherType   string
}

func LoadConfig() *Config {
	return &Config{
		Port:                  getenv("PORT", "8080"),
		DatabaseURL:           getenv("DATABASE_URL", "host=localhost user=postgres dbname=user_service port=5432 sslmode=disable"),
		JWTSecret:             getenv("JWT_SECRET", "holy jesus christ*$%^"),
		RefreshTokenSecret:    getenv("REFRESH_TOKEN_SECRET", "WHat da fuck @#$"),
		AccessTokenExpiry:     getdurationfromenv("ACCESS_TOKEN_EXPIRY", time.Hour*24),
		RefreshTokenExpiry:    getdurationfromenv("REFRESH_TOKEN_EXPIRY", time.Hour*24*2),
		RabbitMQURL:           getenv("RABBITMQ_URL", "ampq://localhost:5672"),
		OpenTelemetryEndPoint: getenv("OPEN_TELEMETRY_ENDPOINT", "localhost:4317"),
		WebSocketPort:         getenv("WEBSOCKET_PORT", ":8082"),
		EventDispatcherType:   getenv("EVENT_DISPATCHER_TYPE", "sync"),
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
