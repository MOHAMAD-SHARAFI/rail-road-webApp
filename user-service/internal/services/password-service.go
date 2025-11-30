package services

import (
	"context"
	"fmt"
	"time"
	"user-service/internal/models"

	usr "user-service/internal/domain/user"
	"user-service/internal/pkg/events"
	"user-service/internal/repositories"
	"user-service/pkg/logger"

	"github.com/sirupsen/logrus"
)

type PasswordService struct {
	userRepo         repositories.UserRepository
	tokenRepo        repositories.TokenRepository
	eventDispatcher  events.EventDispatcher
	resetTokenExpiry time.Duration
}

func NewPasswordService(
	userRepo repositories.UserRepository,
	tokenRepo repositories.TokenRepository,
	eventDispatcher events.EventDispatcher,
	resetTokenExpiry time.Duration,
) *PasswordService {
	return &PasswordService{
		userRepo:         userRepo,
		tokenRepo:        tokenRepo,
		eventDispatcher:  eventDispatcher,
		resetTokenExpiry: resetTokenExpiry,
	}
}

func (p *PasswordService) RequestPasswordReset(ctx context.Context, email string) error {
	logger.Log.WithFields(logrus.Fields{
		"Operation :": "RequestPasswordReset",
		"Email":       email,
	}).Infoln("starting password reset request")

	user, err := p.userRepo.FindByEmail(ctx, email)
	if err != nil {
		logger.Log.WithFields(logrus.Fields{
			"email":       email,
			"error":       err.Error(),
			"Operation :": "find user by email for password reset",
		}).Error("User not found for password reset")
		return fmt.Errorf("user not found")
	}

	passwordResetService := usr.PasswordResetService{}
	resetToken, err := passwordResetService.GenerateResetToken(p.resetTokenExpiry)
	if err != nil {
		logger.Log.WithFields(logrus.Fields{
			"Operation :": "RequestPasswordReset",
			"user_id :":   user.ID,
			"Error":       err.Error(),
		}).Error("Error generating reset token")
		return err
	}

	passwordResetModel := &models.PasswordResetToken{
		UserID:    user.ID,
		Token:     resetToken.Token,
		ExpiresAt: resetToken.ExpireAt,
	}

	err = p.tokenRepo.CreatePasswordResetToken(ctx, passwordResetModel)
	if err != nil {
		logger.Log.WithFields(logrus.Fields{
			"Operation :": "RequestPasswordReset",
			"user_id :":   user.ID,
			"Error":       err.Error(),
		}).Error("failed to save password reset token")
		return err
	}
	//	sending event
	event := usr.NewPasswordResetRequestEvent(user.ID, user.Email, resetToken.Token)
	err = p.eventDispatcher.Dispatch(ctx, event)
	if err != nil {
		logger.Log.WithFields(logrus.Fields{
			"Operation :": "RequestPasswordReset",
			"user_id :":   user.ID,
			"Error":       err.Error(),
		}).Warn("failed to dispatch password reset event")
	}

	logger.Log.WithFields(logrus.Fields{
		"Operation :": "RequestPasswordReset",
		"user_id :":   user.ID,
		"Email":       user.Email,
	}).Info("password reset request completed successfully")

	return nil
}
