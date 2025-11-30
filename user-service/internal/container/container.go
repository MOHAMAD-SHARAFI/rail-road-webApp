package container

import (
	"time"

	"user-service/internal/bootstrap"
	"user-service/internal/config"
	"user-service/internal/handler/grpc"
	"user-service/internal/handler/http"
	evs "user-service/internal/infrastructure/events"
	ev "user-service/internal/pkg/events"
	"user-service/internal/pkg/messaging"
	rbq "user-service/internal/pkg/messaging/rabbitMQ"
	"user-service/internal/repositories"
	"user-service/internal/services"
	"user-service/pkg/logger"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type Container struct {
	DB               *gorm.DB
	Config           *config.Config
	UserRepo         repositories.UserRepository
	TokenRepo        repositories.TokenRepository
	RefreshTokenRepo repositories.RefreshTokenRepository
	AuthService      *services.AuthService
	PasswordService  services.PasswordService
	HttpAuthHandler  *http.AuthHandler
	PasswordHandler  *http.PasswordHandler
	GRPCAuthHandler  *grpc.AuthHandler
	EventDispatcher  ev.EventDispatcher
	MessageProducer  messaging.MessageProducer
	EventHandler     evs.RabbitMQEventHandler
	Logger           *logrus.Logger
}

func NewContainer(cfg *config.Config) (*Container, error) {
	container := &Container{Config: cfg}
	logger.InitLogger()
	container.Logger = logger.Log
	db, err := bootstrap.InitDB(cfg.DatabaseURL)
	if err != nil {
		container.Logger.WithFields(logrus.Fields{
			"Operation :": "InitDB",
			"Error":       err.Error(),
		}).Error("Failed to init DB")
		return nil, err
	}

	container.DB = db
	container.UserRepo = repositories.NewUSerRepository(db)
	container.TokenRepo = repositories.NewTokenRepository(db)
	container.RefreshTokenRepo = repositories.NewRefreshTokenRepository(db)
	dispatcher := ev.NewEventDispatcher()
	container.EventDispatcher = dispatcher
	rabbitConfig := rbq.RBMQConfig{
		URL:          cfg.RabbitMQURL,
		ExchangeName: "user-service",
		RetryCount:   4,
		RetryDelay:   time.Second * 5,
		HeartBeat:    time.Second * 10,
	}

	producer, err := rbq.NewRabbitMQProducer(rabbitConfig)
	if err != nil {
		container.Logger.WithFields(logrus.Fields{
			"Operation :": "NewRabbitMQProducer",
			"Error":       err.Error(),
		}).Error("Failed to create RabbitMQProducer")
		return nil, err
	}

	err = producer.Connect()
	if err != nil {
		container.Logger.WithFields(logrus.Fields{
			"Operation :": "ProducerConnect",
			"Error":       err.Error(),
		}).Error("Failed to connect RabbitMQProducer")
		return nil, err
	}

	container.MessageProducer = producer
	eventHandler := evs.NewRabbitMQEventHandler(producer)
	container.EventHandler = *eventHandler
	dispatcher.Register("user.registered", eventHandler)
	dispatcher.Register("user.signed_in", eventHandler)

	authService := services.NewAuthService(
		container.UserRepo,
		container.RefreshTokenRepo,
		cfg.JWTSecret,
		cfg.RefreshTokenSecret,
		cfg.AccessTokenExpiry,
		cfg.RefreshTokenExpiry,
		dispatcher,
	)
	resetPassService := services.NewPasswordService(
		container.UserRepo,
		container.TokenRepo,
		dispatcher,
		cfg.RefreshTokenExpiry,
	)
	container.PasswordService = *resetPassService
	container.PasswordHandler = http.NewPasswordHandler(&container.PasswordService)

	container.AuthService = authService
	container.HttpAuthHandler = http.NewAuthHandler(container.AuthService)
	container.GRPCAuthHandler = grpc.NewAuthHandler(container.AuthService)
	return container, nil
}
