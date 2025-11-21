package repositories

import (
	"context"
	"payment-service/internal/models"
)

type PaymentRepository interface {
	Create(ctx context.Context, payment *models.Payment) error
	GetByID(ctx context.Context, ID uint) (*models.Payment, error)
	UpdateStatus(ctx context.Context, ID uint, status string) error
	GetFeeStructure(ctx context.Context) (*models.FeeStructure, error)
}
