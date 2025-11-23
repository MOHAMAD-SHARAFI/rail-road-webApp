package repositories

import (
	"context"
	"log"
	"payment-service/internal/models"

	"gorm.io/gorm"
)

type paymentRepository struct {
	db *gorm.DB
}

func NewPaymentRepository(db *gorm.DB) PaymentRepository {
	return paymentRepository{db: db}
}

func (p paymentRepository) Create(ctx context.Context, payment *models.Payment) error {
	err := p.db.WithContext(ctx).Create(payment)
	if err != nil {
		log.Fatal(models.ErrFailedCreatePayment)
	}

	return nil
}

func (p paymentRepository) UpdateStatus(ctx context.Context, id uint, status string) error {
	err := p.db.WithContext(ctx).Model(&models.Payment{}).Where("id = ?", models.Payment{Status: status}).Update("status", status).Error
	if err != nil {
		log.Fatal(models.ErrUpdateFailed)
	}
	return nil
}

func (p paymentRepository) GetByID(ctx context.Context, id uint) (*models.Payment, error) {
	var payment models.Payment
	err := p.db.WithContext(ctx).Where("id = ?", models.Payment{ID: id}).First(&payment)
	if err != nil {
		log.Fatal(models.ErrPaymentNotFound)
	}

	return &payment, nil
}
