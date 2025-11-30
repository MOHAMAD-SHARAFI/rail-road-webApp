package repositories

import (
	"context"
	"user-service/internal/models"
)

type TokenRepository interface {
	CreatePasswordResetToken(ctx context.Context, token *models.PasswordResetToken) error
	FindValidPasswordResetToken(ctx context.Context, userid uint, token string) (*models.PasswordResetToken, error)
	DeletePasswordResetToken(ctx context.Context, tokenID uint) error
}
