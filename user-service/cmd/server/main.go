package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
	"user-service/internal/config"
	"user-service/internal/container"
	"user-service/internal/pkg/messaging"

	"user-service/proto/user-service/gen/auth"

	"user-service/pkg/logger"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
)

func main() {
	// loading logger
	logger.InitLogger()
	//loading Config
	cfg := config.LoadConfig()
	//create container
	cntainer, err := container.NewContainer(cfg)
	if err != nil {
		log.Fatal("failed to init container", err)
	}
	//defer close producer
	defer func(MessageProducer messaging.MessageProducer) {
		err := MessageProducer.Close()
		if err != nil {
			log.Fatal("failed to close message producer", err)
		}
	}(cntainer.MessageProducer)
	// ctx for shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	//for signal
	sigChan := make(chan os.Signal, 1)
	//notify signal                   listen and interrupt
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	//for cancel
	go func() {
		<-sigChan //wait for signal
		cancel()  //cancel context
	}() //goroutine for shutdown

	//	create router and implement routes
	r := gin.Default()

	api := r.Group("/api/v1")
	{
		api.POST("/sign_up", cntainer.HttpAuthHandler.SignUp)
		api.POST("/sign_in", cntainer.HttpAuthHandler.SignIn)
		api.POST("/validate_token", cntainer.HttpAuthHandler.ValidateToken)
		api.POST("/refresh_token", cntainer.HttpAuthHandler.RefreshToken)
		api.POST("/password/reset", cntainer.PasswordHandler.RequestPasswordReset)
	}
	//go for gRPC
	go func() {
		lis, err := net.Listen("tcp", "50051")
		if err != nil {
			log.Fatalf("failed to listen for gRPC: %v", err)
		}
		//build new gRPC server
		grpcServer := grpc.NewServer()
		auth.RegisterAuthServiceServer(grpcServer, cntainer.GRPCAuthHandler)
		log.Println("gRPC server starting on port 50051")
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("failed to serve gRPC: %v", err)
		}
	}() //goroutine for gRPC
	//run HTTP server gracefully
	go func() {
		if err := r.Run(":" + cfg.Port); err != nil {
			log.Fatalf("failed to start server: %v", err)
		}
	}()

	//wait for shutdown
	<-ctx.Done() //wait for cancel
	log.Println("shutting down server...")
}
