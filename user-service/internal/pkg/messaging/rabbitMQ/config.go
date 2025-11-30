package rabbitMQ

import "time"

type RBMQConfig struct {
	URL          string
	ExchangeName string
	RetryCount   int
	RetryDelay   time.Duration
	HeartBeat    time.Duration
}
