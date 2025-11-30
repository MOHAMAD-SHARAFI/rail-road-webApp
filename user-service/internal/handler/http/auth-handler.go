package http

import (
	"net/http"
	"user-service/internal/models"
	"user-service/internal/services"
	"user-service/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type AuthHandler struct {
	authService *services.AuthService
}

func NewAuthHandler(authService *services.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

func (h AuthHandler) SignUp(ctx *gin.Context) {
	logger.Log.WithFields(logrus.Fields{
		"Operation :": "http-SignUp",
		"endpoint :":  ctx.Request.URL.Path,
		"method":      ctx.Request.Method,
	}).Info("sign-Up request received")

	var req models.SignUpRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		logger.Log.WithFields(logrus.Fields{
			"Operation :": "http-SignUp",
			"error :":     err.Error(),
		}).Warn("invalid request body")
		ctx.JSON(http.StatusBadRequest, gin.H{"Error": "invalid request body"})
		return
	}

	logger.Log.WithFields(logrus.Fields{
		"Operation :": "http-SignUp",
		"username :":  req.Username,
		"Email :":     req.Email,
	}).Debug("processing sign-Up request")

	response, err := h.authService.SignUp(ctx.Request.Context(), req.Username, req.Email, req.Password)
	if err != nil {
		logger.Log.WithFields(logrus.Fields{
			"Operation :": "http-SignUp",
			"Username :":  req.Username,
			"Email :":     req.Email,
			"error :":     err.Error(),
		}).Warn("failed to sign-Up")
		ctx.JSON(http.StatusBadRequest, gin.H{"Error :": err.Error()})
		return
	}

	logger.Log.WithFields(logrus.Fields{
		"Operation :": "http-SignUp",
		"User_id :":   response.ID,
		"Email :":     req.Email,
		"Username :":  req.Username,
	}).Info("SignUp completed successfully")
	ctx.JSON(http.StatusOK, response)
}

func (h AuthHandler) SignIn(ctx *gin.Context) {
	logger.Log.WithFields(logrus.Fields{
		"Operation :": "http-SignIn",
		"endpoint :":  ctx.Request.URL.Path,
		"method":      ctx.Request.Method,
	}).Info("sign-In request received")

	var req models.SignInRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		logger.Log.WithFields(logrus.Fields{
			"Operation :": "http-SignIn",
			"error :":     err.Error(),
		}).Warn("invalid request body")
		ctx.JSON(http.StatusBadRequest, gin.H{"Error": "invalid request body"})
		return
	}

	logger.Log.WithFields(logrus.Fields{
		"Operation :": "http-SignIn",
		"Username :":  req.UserName,
	}).Debug("processing sign-In request")

	response, err := h.authService.SignIn(ctx.Request.Context(), req.UserName, req.Password)
	if err != nil {
		logger.Log.WithFields(logrus.Fields{
			"Operation :": "http-SignIn",
			"Username :":  req.UserName,
			"error :":     err.Error(),
		}).Warn("SingIn failed")
		ctx.JSON(http.StatusUnauthorized, gin.H{"Error :": err.Error()})
		return
	}

	logger.Log.WithFields(logrus.Fields{
		"Operation :": "http-SignIn",
		"User_Id :":   response.UserID,
		"Username :":  req.UserName,
	}).Info("SignIn completed successfully")

	ctx.JSON(http.StatusOK, response)
}

func (h AuthHandler) ValidateToken(ctx *gin.Context) {
	logger.Log.WithFields(logrus.Fields{
		"Operation :": "http-ValidateToken",
		"endpoint :":  ctx.Request.URL.Path,
		"method":      ctx.Request.Method,
	}).Debug("Validate request received")

	var req models.ValidateTokenRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		logger.Log.WithFields(logrus.Fields{
			"Operation :": "http-ValidateToken",
			"error :":     err.Error(),
		}).Warn("invalid request body")
		ctx.JSON(http.StatusBadRequest, gin.H{"Error :": "Invalid request body"})
		return
	}

	response, err := h.authService.ValidateToken(ctx.Request.Context(), req.Token)
	if err != nil {
		logger.Log.WithFields(logrus.Fields{
			"Operation :": "http-ValidateToken",
			"error :":     err.Error(),
		}).Warn("failed to ValidateToken")

		ctx.JSON(http.StatusBadRequest, gin.H{"Error :": err.Error()})
		return
	}

	logger.Log.WithFields(logrus.Fields{
		"Operation :": "http-ValidateToken",
		"User_Id :":   response.UserID,
		"valid :":     response.Valid,
	}).Info("ValidateToken completed successfully")
	ctx.JSON(http.StatusOK, response)
}

func (h AuthHandler) RefreshToken(ctx *gin.Context) {
	logger.Log.WithFields(logrus.Fields{
		"Operation :": "http-RefreshToken",
		"endpoint :":  ctx.Request.URL.Path,
		"method":      ctx.Request.Method,
	}).Debug("refresh-token request received")

	var req models.RefreshTokenRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		logger.Log.WithFields(logrus.Fields{
			"Operation :": "http-RefreshToken",
			"error :":     err.Error(),
		}).Warn("invalid request body")
		ctx.JSON(http.StatusBadRequest, gin.H{"Error :": "Invalid request body"})
		return
	}

	response, err := h.authService.RefreshToken(ctx.Request.Context(), req.RefreshToken)
	if err != nil {
		logger.Log.WithFields(logrus.Fields{
			"Operation :": "http-RefreshToken",
			"error :":     err.Error(),
		}).Warn("failed to RefreshToken")
		ctx.JSON(http.StatusUnauthorized, gin.H{"Error :": err.Error()})
		return
	}

	logger.Log.WithFields(logrus.Fields{
		"Operation :": "http-RefreshToken",
	}).Info("Token Refresh completed successfully")

	ctx.JSON(http.StatusOK, response)
}
