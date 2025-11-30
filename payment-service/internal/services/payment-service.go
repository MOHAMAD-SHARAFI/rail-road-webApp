package services

import (
	"context"
	"errors"
	"payment-service/internal/clients"
	"payment-service/internal/models"
	"payment-service/internal/repositories"
	"payment-service/pkg/logger"

	"github.com/sirupsen/logrus"
)

type PaymentService struct {
	paymentRepo repositories.PaymentRepository
	userClient  *clients.UserClients
	feeConfig   *FeeConfig
}

type FeeConfig struct {
	Percentage float64
	MinFee     float64
}

func NewPaymentService(paymentRepo repositories.PaymentRepository, userClient *clients.UserClients, feeConfig *FeeConfig) *PaymentService {
	return &PaymentService{
		paymentRepo: paymentRepo,
		userClient:  userClient,
		feeConfig:   feeConfig,
	}
}

func (s *PaymentService) CreatePayment(ctx context.Context, userToken string, amount float64) (*models.Payment, error) {
	logger.Log.WithFields(logrus.Fields{
		"Operation": "CreatePayment",
		"Amount":    amount,
	}).Info("creating new payment")

	//	validate user token
	valid, userID, err := s.userClient.ValidateToken(ctx, userToken)
	if err != nil {
		logger.Log.WithFields(logrus.Fields{
			"Operation": "CreatePayment",
			"Error":     err.Error(),
		}).Error("token validation failed")
		return nil, errors.New("token validation failed")
	}

	if !valid {
		logger.Log.WithFields(logrus.Fields{
			"Operation": "CreatePayment",
		}).Warn("invalid user token")
		return nil, errors.New("invalid user token")
	}

	//	calculate fee
	fee := amount * s.feeConfig.Percentage / 100
	if fee < s.feeConfig.MinFee {
		fee = s.feeConfig.MinFee
	}
	total := amount + fee
	// create payment record
	payment := &models.Payment{
		UserID:   userID,
		Amount:   amount,
		Fee:      fee,
		Total:    total,
		Currency: "IRR",
		Status:   "PENDING",
	}

	if err := s.paymentRepo.Create(ctx, payment); err != nil {
		logger.Log.WithFields(logrus.Fields{
			"Operation": "CreatePayment",
			"user_id :": userID,
			"amount :":  amount,
			"Error":     err.Error(),
		}).Error("failed to create payment record")
		return nil, err
	}

	logger.Log.WithFields(logrus.Fields{
		"Operation":    "CreatePayment",
		"payment_id :": payment.ID,
		"user_id :":    userID,
		"amount :":     amount,
		"fee :":        fee,
		"total :":      total,
	}).Info("payment created successfully")

	return payment, nil

}

func (s *PaymentService) ProcessPayment(ctx context.Context, paymentID uint) error {
	logger.Log.WithFields(logrus.Fields{
		"Operation":    "ProcessPayment",
		"payment_id :": paymentID,
	}).Info("starting payment processing")

	//	Get payment
	payment, err := s.paymentRepo.GetByID(ctx, paymentID)
	if err != nil {
		logger.Log.WithFields(logrus.Fields{
			"Operation":    "ProcessPayment",
			"payment_id :": paymentID,
			"error":        err.Error(),
		}).Error("payment not found")
		return errors.New("payment not found")
	}

	if payment.Status != "PENDING" {
		logger.Log.WithFields(logrus.Fields{
			"Operation":        "ProcessPayment",
			"payment_id :":     paymentID,
			"current_status :": payment.Status,
		}).Warn("payment is not in pending state")
		return errors.New("payment is not in pending status")
	}

	result, err := s.gateway.RequestPayment(
		payment.ID, payment.Amount,
		"http://localhost:8080/api/v1/payment/callback",
	)
	if err != nil {
		payment.Status = "FAILED"
		err := s.paymentRepo.UpdateStatus(ctx, payment.ID, "FAILED")
		if err != nil {
			return err
		}
		logger.Log.WithFields(logrus.Fields{
			"Operation":    "ProcessPayment",
			"payment_id :": payment.ID,
			"Error":        err.Error(),
		}).Error("Gateway request failed")
	}

	if result.Success {
		payment.Status = "PROCESSING"
		payment.GateWayTransactionID = result.RefID
		err := s.paymentRepo.UpdateStatus(ctx, payment.ID, "PROCESSING")
		if err != nil {
			return err
		}
	}

	logger.Log.WithFields(logrus.Fields{
		"Operation":    "ProcessPayment",
		"payment_id :": payment.ID,
	}).Info("payment process successfully")

	return nil
}
