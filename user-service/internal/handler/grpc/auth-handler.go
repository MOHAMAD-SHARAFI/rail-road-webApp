package grpc

import (
	"context"
	"user-service/internal/repositories"
	"user-service/proto/user-service/gen/user"
)

type AuthHandler struct {
	repo repositories.UserRepository
}

func NewAuthHandler(repo repositories.UserRepository) *AuthHandler {
	return &AuthHandler{repo: repo}
}

func (h AuthHandler) nameValidateToken(ctx context.Context, req *user.ValidateTokenRequest) (*user.ValidateTokenResponse, error) {
	Token := req.GetToken()
}
