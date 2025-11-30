package messaging

import "time"

type Message struct {
	ID         string
	Exchanges  string
	RoutingKey string
	Body       []byte
	Timestamp  time.Time
}

type ProducerConfig struct {
	URL          string
	ExchangeName string
	RetryCount   int
	RetryDelay   time.Duration
}
