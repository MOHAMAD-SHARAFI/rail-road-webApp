package grpc

import (
	"context"

	"user-service/internal/services"
	"user-service/pkg/logger"
	"user-service/proto/user-service/gen/auth"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AuthHandler struct {
	auth.UnimplementedAuthServiceServer
	authService *services.AuthService
}

func NewAuthHandler(authService *services.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

func (h AuthHandler) ValidateToken(ctx context.Context, req *auth.ValidateTokenRequest) (*auth.ValidateTokenResponse, error) {
	logger.Log.WithFields(logrus.Fields{
		"Operation :": "gRPC_validate_token",
	}).Info("gRPC validate token request received")

	response, err := h.authService.ValidateToken(ctx, req.Token)
	if err != nil {
		logger.Log.WithFields(logrus.Fields{
			"Operation :": "gRPC_validate_token",
			"Error":       err.Error(),
		}).Error("gRPC token validation failed")
		return nil, status.Error(codes.Internal, "internal server error")
	}

	grpcResponse := &auth.ValidateTokenResponse{
		Valid:   response.Valid,
		UserId:  uint64(response.UserID),
		Message: response.Message,
	}

	logger.Log.WithFields(logrus.Fields{
		"Operation :": "gRPC_validate_token",
		"user_id :":   response.UserID,
		"valid :":     response.Valid,
	}).Info("gRPC token validation succeeded")

	return grpcResponse, nil
}
