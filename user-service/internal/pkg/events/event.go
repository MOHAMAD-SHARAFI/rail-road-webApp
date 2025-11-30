package events

import (
	"context"
	"time"
)

type Event interface {
	GetName() string
	GetPayload() interface{}
	GetTimestamp() time.Time
}

type EventHandler interface {
	Handle(ctx context.Context, event Event) error
}
