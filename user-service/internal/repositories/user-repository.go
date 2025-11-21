package repositories

import (
	"context"
	"user-service/internal/models"
)

type UserRepository interface {
	Create(ctx context.Context, user *models.User) error
	GetByID(ctx context.Context, id uint) (*models.User, error)
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	GetByPhoneNumber(ctx context.Context, phone uint) (*models.User, error)
	Update(ctx context.Context, user *models.User) error
	CreatePasswordResetToken(ctx context.Context, token *models.PassworResetToken) error
	GetValidPasswordResetToken(ctx context.Context, userID uint, token string) (*models.PassworResetToken, error)
	DeletePasswordResetToken(ctx context.Context, tokenID uint) error
}
