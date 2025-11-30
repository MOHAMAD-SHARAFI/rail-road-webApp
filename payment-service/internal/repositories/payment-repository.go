package repositories

import (
	"context"
	"payment-service/internal/models"
)

type PaymentRepository interface {
	Create(ctx context.Context, payment *models.Payment) error
	GetByID(ctx context.Context, ID uint) (*models.Payment, error)
	GetByUserID(ctx context.Context, userID uint) ([]models.Payment, error)
	UpdateStatus(ctx context.Context, id uint, status string) error
	UpdateGatewayTransactionID(ctx context.Context, id uint, gatewayID string) error
}
