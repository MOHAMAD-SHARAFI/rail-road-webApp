package events

import "time"

type BaseEvent struct {
	name      string
	payload   interface{}
	timestamp time.Time
}

func NewBaseEvent(name string, payload any) BaseEvent {
	return BaseEvent{
		name:      name,
		payload:   payload,
		timestamp: time.Now(),
	}
}

func (e BaseEvent) GetName() string {
	return e.name
}

func (e BaseEvent) GetPayload() interface{} {
	return e.payload
}

func (e BaseEvent) GetTimestamp() time.Time {
	return e.timestamp
}
