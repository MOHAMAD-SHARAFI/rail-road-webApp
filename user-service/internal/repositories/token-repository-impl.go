package repositories

import (
	"context"
	"time"

	"user-service/internal/models"
	"user-service/pkg/logger"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type tokenRepository struct {
	db *gorm.DB
}

func NewTokenRepository(db *gorm.DB) TokenRepository {
	return &tokenRepository{db: db}
}

func (t *tokenRepository) CreatePasswordResetToken(ctx context.Context, token *models.PassworResetToken) error {
	result := t.db.WithContext(ctx).Create(token)
	if result.Error != nil {
		logger.Log.WithFields(logrus.Fields{
			"operation":     "CreatePasswordResetToken",
			"error-details": result.Error,
		}).Error("cannot create password reset token")
	}

	logger.Log.WithFields(logrus.Fields{
		"operation": "CreatePasswordResetToken",
		"token":     token,
	}).Info("password reset token created successfully")

	return result.Error
}

func (t *tokenRepository) FindValidPasswordResetToken(ctx context.Context, userid uint, token string) (*models.PassworResetToken, error) {
	var passwordResetToken models.PassworResetToken
	validToken := t.db.WithContext(ctx).Where("user_id = ? AND token = ? AND expires_at > ?", userid, token, time.Now()).First(&passwordResetToken)
	if validToken.Error != nil {
		logger.Log.WithFields(logrus.Fields{
			"operation":     "FindValidPasswordResetToken",
			"error-details": validToken.Error,
		}).Error("cannot find valid password reset token")
		return nil, validToken.Error
	}

	logger.Log.WithFields(logrus.Fields{
		"operation": "FindValidPasswordResetToken",
		"token":     token,
	}).Info("valid password reset token found successfully")
	return &passwordResetToken, nil
}

func (t *tokenRepository) DeletePasswordResetToken(ctx context.Context, tokenID uint) error {
	var token models.PassworResetToken
	result := t.db.WithContext(ctx).Delete(&token, tokenID)
	if result.Error != nil {
		logger.Log.WithFields(logrus.Fields{
			"operation":     "DeletePasswordResetToken",
			"error-details": result.Error,
		}).Error("cannot delete password reset token")
		return result.Error
	}

	logger.Log.WithFields(logrus.Fields{
		"operation": "DeletePasswordResetToken",
		"token":     token,
	}).Info("password reset token deleted successfully")
	return nil
}
