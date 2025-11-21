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
		log.Printf("failed to create payment: %v", err)
	}

	return nil
}

func (p paymentRepository) UpdateStatus(ctx context.Context, id uint, status string) error {
	err := p.db.WithContext(ctx).Model(&models.Payment{}).Where("id = ?", id).Update("status", status).Error
	if err != nil {
		log.Printf("failed to update status: %v", err)
	}
	return nil
}

func (p paymentRepository) GetByID(ctx context.Context, id uint) (*models.Payment, error) {
	var payment models.Payment
	err := p.db.WithContext(ctx).Where("id = ?", id).First(&payment)
	if err != nil {
		log.Printf("failed to fetch payment: %v", err)
	}

	return &payment, nil
}

func (p paymentRepository) GetFeeStructure(ctx context.Context) (*models.FeeStructure, error) {
	var feeStructure models.FeeStructure
	err := p.db.WithContext(ctx).First(&feeStructure)
	if err != nil {
		log.Printf("failed to fetch feeStructure: %v", err)
	}

	return &feeStructure, nil
}
