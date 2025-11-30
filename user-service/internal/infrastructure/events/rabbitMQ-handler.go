package events

import (
	"context"
	"encoding/json"
	"user-service/internal/pkg/events"

	"user-service/internal/pkg/messaging"
	"user-service/pkg/logger"

	"github.com/sirupsen/logrus"
)

type RabbitMQEventHandler struct {
	producer messaging.MessageProducer
}

func NewRabbitMQEventHandler(producer messaging.MessageProducer) *RabbitMQEventHandler {
	return &RabbitMQEventHandler{
		producer: producer,
	}
}

func (h *RabbitMQEventHandler) Handle(ctx context.Context, event events.Event) error {
	messageBytes, err := json.Marshal(event.GetPayload())
	if err != nil {
		logger.Log.WithFields(logrus.Fields{
			"event-name":  event.GetName(),
			"Operation :": "HandleEvent",
			"Error":       err.Error(),
		}).Error("failed to marshal event payload")

		return err
	}

	err = h.producer.Publish(ctx, "user_events", event.GetName(), messageBytes)
	if err != nil {
		logger.Log.WithFields(logrus.Fields{
			"event-name":  event.GetName(),
			"Operation :": "PublishEvent",
			"Error":       err.Error(),
		}).Error("failed to publish event to rabbitMQ")

		return err
	}

	logger.Log.WithFields(logrus.Fields{
		"event-name":  event.GetName(),
		"Operation :": "HandleEvent",
	}).Info("event published to rabbitMQ successfully")

	return nil
}
