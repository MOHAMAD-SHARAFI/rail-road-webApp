package http

import (
	"net/http"
	"payment-service/internal/models"
	"payment-service/internal/services"
	"payment-service/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type PaymentHandler struct {
	paymentService *services.PaymentService
}

func NewPaymentHandler(paymentService *services.PaymentService) *PaymentHandler {
	return &PaymentHandler{
		paymentService: paymentService,
	}
}

func (h *PaymentHandler) CreatePayment(ctx *gin.Context) {
	logger.Log.WithFields(logrus.Fields{
		"Operation": "CreatePayment",
	}).Info("Payment is Creating")
	var req models.CreatePaymentRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		logger.Log.WithFields(logrus.Fields{
			"Operation": "CreatePayment",
			"Error":     err.Error(),
		}).Warn("Invalid request body")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userToken := ctx.GetHeader("Authorization")

	payment, err := h.paymentService.CreatePayment(ctx.Request.Context(), userToken, req.Amount)
	if err != nil {
		logger.Log.WithFields(logrus.Fields{
			"Operation": "CreatePayment",
			"Error":     err.Error(),
		}).Warn("Payment creation failed")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	logger.Log.WithFields(logrus.Fields{
		"Operation": "CreatePayment",
		"Amount":    req.Amount,
		"UserToken": userToken,
	}).Info("Payment created successfully")
	ctx.JSON(http.StatusOK, gin.H{"payment": payment})

}
