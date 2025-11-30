package clients

import (
	"context"
	"payment-service/pkg/logger"

	"user-service/proto/user-service/gen/auth"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type UserClients struct {
	conn   *grpc.ClientConn
	client auth.AuthServiceClient
}

func NewUserClient(serverAddr string) (*UserClients, error) {
	conn, err := grpc.Dial(serverAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	client := auth.NewAuthServiceClient(conn)
	return &UserClients{
		conn:   conn,
		client: client,
	}, nil
}

func (c *UserClients) ValidateToken(ctx context.Context, token string) (bool, uint, error) {
	logger.Log.WithFields(logrus.Fields{
		"Operation :": "gRPC validate token",
	}).Debug("sending gRPC token validation request")

	response, err := c.client.ValidateToken(ctx, &auth.ValidateTokenRequest{Token: token})
	if err != nil {
		logger.Log.WithFields(logrus.Fields{
			"Operation :": "gRPC_validate_token",
			"Error":       err.Error(),
		}).Error("gRPC validate token failed")
		return false, 0, err
	}

	logger.Log.WithFields(logrus.Fields{
		"Operation :": "gRPC_validate_token",
		"user_id :":   response.UserId,
		"valid :":     response.Valid,
	}).Info("gRPC token validation response received")

	return response.Valid, uint(response.UserId), nil
}

func (c UserClients) Close() error {
	return c.conn.Close()
}
