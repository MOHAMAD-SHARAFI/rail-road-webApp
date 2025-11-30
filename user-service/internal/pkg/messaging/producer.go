package messaging

import "context"

type MessageProducer interface {
	Connect() error
	Publish(ctx context.Context, exchange string, routingKey string, message []byte) error
	Close() error
}
