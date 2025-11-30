package events

import (
	"context"
	"sync"

	"user-service/pkg/logger"

	"github.com/sirupsen/logrus"
)

type EventDispatcher interface {
	Register(eventName string, handler EventHandler)
	Dispatch(ctx context.Context, event Event) error
}
type eventDispatcher struct {
	handlers map[string][]EventHandler
	mutex    sync.RWMutex
}

func NewEventDispatcher() EventDispatcher {
	return &eventDispatcher{
		handlers: make(map[string][]EventHandler),
	}
}

func (d *eventDispatcher) Register(eventName string, handler EventHandler) {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	d.handlers[eventName] = append(d.handlers[eventName], handler)

	logger.Log.WithFields(logrus.Fields{
		"eventName":   eventName,
		"Operation :": "RegistereventHandler",
	}).Debug("Registered event handler")
}

func (d *eventDispatcher) Dispatch(ctx context.Context, event Event) error {
	d.mutex.RLock()
	handlers, exists := d.handlers[event.GetName()]
	d.mutex.RUnlock()
	if !exists {
		return nil
	}

	var errors []error

	for _, handler := range handlers {
		err := handler.Handle(ctx, event)
		if err != nil {
			errors = append(errors, err)

			logger.Log.WithFields(logrus.Fields{
				"eventName":   event.GetName(),
				"Operation :": "DispatchEvent",
				"Error":       err.Error(),
			}).Error("failed handling event")
		}
	}

	if len(errors) > 0 {
		return errors[0]
	}

	return nil
}
