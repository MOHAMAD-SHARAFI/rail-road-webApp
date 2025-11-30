package http

import (
	"net/http"
	"user-service/internal/services"
	"user-service/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type PasswordHandler struct {
	passwordService *services.PasswordService
}

func NewPasswordHandler(passwordService *services.PasswordService) *PasswordHandler {
	return &PasswordHandler{
		passwordService: passwordService,
	}
}

func (h PasswordHandler) RequestPasswordReset(ctx *gin.Context) {
	logger.Log.WithFields(logrus.Fields{
		"Operation :": "http-RequestPasswordReset",
		"EndPoint":    ctx.Request.URL.Path,
		"Method":      ctx.Request.Method,
	}).Info("password reset request received")
	var req struct {
		Email string `json:"email" binding:"required,email"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		logger.Log.WithFields(logrus.Fields{
			"Operation :": "http-RequestPasswordReset",
			"Error":       err.Error(),
		}).Warn("invalid request body")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	err := h.passwordService.RequestPasswordReset(ctx.Request.Context(), req.Email)
	if err != nil {
		logger.Log.WithFields(logrus.Fields{
			"Operation :": "http-RequestPasswordReset",
			"Email":       req.Email,
			"Error":       err.Error(),
		}).Warn("password reset request failed")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	logger.Log.WithFields(logrus.Fields{
		"Operation :": "http-RequestPasswordReset",
		"Email":       req.Email,
	}).Info("password reset request completed successfully")
	ctx.JSON(http.StatusOK, gin.H{"Message": "password reset email sent successfully"})
}
