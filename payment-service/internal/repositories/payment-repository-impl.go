package repositories

import (
	"context"

	"payment-service/internal/models"
	"payment-service/pkg/logger"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type paymentRepository struct {
	db *gorm.DB
}

func NewPaymentRepository(db *gorm.DB) PaymentRepository {
	return &paymentRepository{db: db}
}

func (p *paymentRepository) Create(ctx context.Context, payment *models.Payment) error {
	logger.Log.WithFields(logrus.Fields{
		"operation": "Create a payment in progress",
		"amount":    payment.Amount,
	}).Info("creating payment")

	result := p.db.WithContext(ctx).Create(payment)

	return result.Error
}

func (p *paymentRepository) UpdateStatus(ctx context.Context, id uint, status string) error {
	logger.Log.WithFields(logrus.Fields{
		"operation": "Update status of payment in progress",
		"status":    status,
	}).Info("updating payment status")
	result := p.db.WithContext(ctx).Model(&models.Payment{}).Where("id = ?", id).Update("status", status)
	if result.Error != nil {
		logger.Log.WithFields(logrus.Fields{
			"operation": "Update status",
			"status":    status,
		}).Error("failed to update payment status")
	}

	logger.Log.WithFields(logrus.Fields{
		"operation": "Update status of payment",
		"status":    status,
	}).Info("payment status updated successfully")
	return nil
}

func (p *paymentRepository) GetByID(ctx context.Context, id uint) (*models.Payment, error) {
	logger.Log.WithFields(logrus.Fields{
		"operation": "Get payment by id in progress",
		"id":        id,
	}).Info("getting payment by id")

	var payment models.Payment
	result := p.db.WithContext(ctx).Where("id = ?", id).First(&payment)
	if result.Error != nil {
		logger.Log.WithFields(logrus.Fields{
			"operation": "Get payment by id",
			"id":        id,
		}).Error("failed to get payment by id")
	}

	logger.Log.WithFields(logrus.Fields{
		"operation": "Get payment by id",
		"id":        id,
	}).Info("payment with this id is found")
	return &payment, nil
}

func (p *paymentRepository) UpdateGatewayTransactionID(ctx context.Context, id uint, gatewayID string) error {
	logger.Log.WithFields(logrus.Fields{
		"operation":  "Update gateway transaction id",
		"gateway_id": gatewayID,
	}).Info("updating gateway transaction id")

	result := p.db.WithContext(ctx).Where("id = ?", id).Update("gateway_transaction_id", gatewayID)
	if result.Error != nil {
		logger.Log.WithFields(logrus.Fields{
			"operation":  "Update gateway transaction id",
			"gateway_id": gatewayID,
			"Error :":    result.Error,
		}).Error("failed to update gateway transaction id")
	}

	logger.Log.WithFields(logrus.Fields{
		"operation":  "Update gateway transaction id",
		"gateway_id": gatewayID,
	}).Info("gateway transaction id updated successfully")

	return nil

}
func (p *paymentRepository) GetByUserID(ctx context.Context, userID uint) ([]models.Payment, error) {
	logger.Log.WithFields(logrus.Fields{
		"operation": "Get payment by user id in progress",
		"user_id":   userID,
	}).Info("getting payment by user id")

	var payments []models.Payment

	result := p.db.WithContext(ctx).Where("user_id = ?", userID).Find(&payments)
	if result.Error != nil {
		logger.Log.WithFields(logrus.Fields{
			"operation": "Get payment by user id",
			"user_id":   userID,
		}).Error("failed to get payment by user id")
	}

	logger.Log.WithFields(logrus.Fields{
		"operation": "Get payment by user id",
		"user_id":   userID,
	}).Info("payment with this user is found")
	return payments, nil

}
