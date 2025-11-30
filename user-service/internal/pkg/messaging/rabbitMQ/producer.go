package rabbitMQ

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/sirupsen/logrus"

	"user-service/pkg/logger"
)

type RMQProducer struct {
	conn      *amqp.Connection
	channel   *amqp.Channel
	config    RBMQConfig
	mu        *sync.Mutex
	connected bool
}

func NewRabbitMQProducer(config RBMQConfig) (*RMQProducer, error) {
	return &RMQProducer{
		config:    config,
		mu:        &sync.Mutex{},
		connected: false,
	}, nil
}

func (r *RMQProducer) Connect() error {
	logger.Log.WithFields(logrus.Fields{
		"Operation :": "RabbitMQ_Connect",
		"URL :":       r.config.URL,
		"Exchange :":  r.config.ExchangeName,
	}).Info("Connecting to RabbitMQ")
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.connected {
		logger.Log.WithFields(logrus.Fields{
			"Operation :": "RabbitMQ_Connect",
			"Status :":    "Already connected",
		}).Debug("Already connected to RabbitMQ")
		return nil
	}

	conn, err := amqp.Dial(r.config.URL)
	if err != nil {
		logger.Log.WithFields(logrus.Fields{
			"Operation :": "RabbitMQ_Connect",
			"Error :":     err.Error(),
			"URL :":       r.config.URL,
		}).Error("Failed to connect to RabbitMQ")

		return fmt.Errorf("failed to connect to RabbitMQ : %w", err)
	}

	channel, err := conn.Channel()
	if err != nil {
		logger.Log.WithFields(logrus.Fields{
			"Operation :": "RabbitMQ_Connect",
			"Error :":     err.Error(),
		}).Error("Failed to create channel")

		err := conn.Close()
		if err != nil {
			return err
		}
		return fmt.Errorf("failed to OPEN channel : %w", err)
	}

	err = channel.ExchangeDeclare(
		r.config.ExchangeName,
		"topic",
		true,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		logger.Log.WithFields(logrus.Fields{
			"Operation :": "RabbitMQ_Connect",
			"Error :":     err.Error(),
			"Exchange :":  r.config.ExchangeName,
		}).Error("Failed to declare exchange")

		err = channel.Close()
		if err != nil {
			return err
		}
		err = conn.Close()
		if err != nil {
			return err
		}
		return fmt.Errorf("failed to declare exchange : %w", err)
	}

	r.conn = conn
	r.channel = channel
	r.connected = true
	logger.Log.WithFields(logrus.Fields{
		"Operation :": "RabbitMQ_Connect",
		"Status :":    "Connected",
		"Exchange :":  r.config.ExchangeName,
	}).Info("Successfully Connected to RabbitMQ")

	return nil
}

func (r *RMQProducer) Publish(ctx context.Context, exchange string, routingKey string, message []byte) error {
	logger.Log.WithFields(logrus.Fields{
		"Operation :":  "RabbitMQ_Publish",
		"Exchange :":   r.config.ExchangeName,
		"RoutingKey :": routingKey,
		"MessageSize":  len(message),
	}).Debug("Publishing message to RabbitMQ")

	// if producer not connected try to connect it
	if !r.connected {
		logger.Log.WithFields(logrus.Fields{
			"Operation :": "RabbitMQ_Publish",
		}).Warn("RabbitMQ is not connected , attempting to reconnect")

		if err := r.Connect(); err != nil {
			logger.Log.WithFields(logrus.Fields{
				"Operation :": "RabbitMQ_Publish",
				"Error :":     err.Error(),
			}).Error("Failed to reconnect to RabbitMQ")
			return err
		}
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	// if exchange is empty use default exchange
	if exchange == "" {
		exchange = r.config.ExchangeName
	}

	var lastErr error

	// publishing setting
	publishing := amqp.Publishing{
		ContentType:  "application/json",
		Body:         message,
		DeliveryMode: amqp.Persistent,
		Timestamp:    time.Now(),
	}

	//  retry logic for sending message
	for i := 0; i < r.config.RetryCount; i++ {
		//	if this is not first try to get delay
		if i > 0 {
			logger.Log.WithFields(logrus.Fields{
				"Operation :":     "RabbitMQ_Publish",
				"Retry_attempt :": i,
				"Retry_Delay :":   r.config.RetryDelay.String(),
			}).Warn("Retrying message publish")

			time.Sleep(r.config.RetryDelay)
		}

		err := r.channel.Publish(
			exchange,
			routingKey,
			false,
			false,
			publishing,
		)

		if err == nil {
			logger.Log.WithFields(logrus.Fields{
				"Operation :":   "RabbitMQ_Publish",
				"Exchange :":    r.config.ExchangeName,
				"RoutingKey :":  routingKey,
				"Retry_Count :": i,
			}).Info("Successfully published message")
		}

		lastErr = err

		logger.Log.WithFields(logrus.Fields{
			"Operation :":   "RabbitMQ_Publish",
			"Error :":       err.Error(),
			"Retry_Count :": i + 1,
		}).Warn("failed to publish message, will retry")

		if errors.Is(err, amqp.ErrClosed) {
			logger.Log.WithFields(logrus.Fields{
				"Operation :": "RabbitMQ_Publish",
			}).Warn("RabbitMQ connection is closed, reconnecting")

			r.connected = false
			if connectErr := r.Connect(); connectErr != nil {
				return connectErr
			}
		}
	}

	logger.Log.WithFields(logrus.Fields{
		"Operation :":   "RabbitMQ_Publish",
		"Exchange :":    r.config.ExchangeName,
		"RoutingKey :":  routingKey,
		"Error :":       lastErr.Error(),
		"Retry_Count :": r.config.RetryCount,
	}).Error("Failed to publish message after all retry attempts")

	return fmt.Errorf("failed to publish message after %d attempts : %w", r.config.RetryCount, lastErr)
}

func (r *RMQProducer) Close() error {
	logger.Log.WithFields(logrus.Fields{
		"Operation :": "RabbitMQ_Close",
	}).Info("Closing RabbitMQ connection")

	r.mu.Lock()
	defer r.mu.Unlock()

	var errs []string

	//first close the channel
	if r.channel != nil {
		if err := r.channel.Close(); err != nil {
			logger.Log.WithFields(logrus.Fields{
				"Operation :": "RabbitMQ_Close",
				"Error :":     err.Error(),
			}).Warn("Failed to close RabbitMQ channel")
			errs = append(errs, fmt.Sprintf("Channel: %s", err.Error()))
		} else {
			logger.Log.WithFields(logrus.Fields{
				"Operation :": "RabbitMQ_Close",
			}).Debug("Successfully closed RabbitMQ channel")
		}
		r.channel = nil
	}

	//	now closing the connection
	if r.conn != nil {
		if err := r.conn.Close(); err != nil {
			logger.Log.WithFields(logrus.Fields{
				"Operation :": "RabbitMQ_Close",
				"Error :":     err.Error(),
			}).Warn("Failed to close RabbitMQ connection")
			errs = append(errs, fmt.Sprintf("Connection: %s", err.Error()))
		} else {
			logger.Log.WithFields(logrus.Fields{
				"Operation :": "RabbitMQ_Close",
			}).Debug("Successfully closed RabbitMQ connection")
		}
		r.conn = nil
	}

	r.connected = false

	if len(errs) > 0 {
		return fmt.Errorf("errors closing RsbbitMQ: %s", strings.Join(errs, "; "))
	}

	logger.Log.WithFields(logrus.Fields{
		"Operation :": "RabbitMQ_Close",
	}).Info("RabbitMQ connection closed successfully")
	return nil
}
