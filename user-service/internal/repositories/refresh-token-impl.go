package repositories

import (
	"context"

	"user-service/internal/models"
	"user-service/pkg/logger"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type refreshTokenRepository struct {
	db *gorm.DB
}

func NewRefreshTokenRepository(db *gorm.DB) RefreshTokenRepository {
	return &refreshTokenRepository{
		db: db,
	}
}

func (r *refreshTokenRepository) Create(ctx context.Context, token *models.RefreshToken) error {
	logger.Log.WithFields(logrus.Fields{
		"Operation :": "starting Create a refresh token",
	}).Info("Starting Create a refresh token")
	result := r.db.WithContext(ctx).Create(token)
	if result.Error != nil {
		logger.Log.WithFields(logrus.Fields{
			"Operation :": "starting Create a refresh token",
			"Error":       result.Error,
		}).Error("Error while creating a refresh token")
		return result.Error
	}
	logger.Log.WithFields(logrus.Fields{
		"Operation :": "finished Create a refresh token",
		"user_id :":   token.UserID,
	}).Info("Finished Create a refresh token")
	return nil
}

func (r *refreshTokenRepository) FindByToken(ctx context.Context, token string) (*models.RefreshToken, error) {
	logger.Log.WithFields(logrus.Fields{
		"Operation :": "starting FindByToken",
	}).Info("Starting FindByToken")
	refreshToken := &models.RefreshToken{
		Token: token,
	}
	result := r.db.WithContext(ctx).Where("token = ?", token).First(refreshToken)
	if result.Error != nil {
		logger.Log.WithFields(logrus.Fields{
			"Operation :": "starting FindByToken",
			"Error":       result.Error.Error(),
		}).Error("Error while finding a refresh token")
		return nil, result.Error
	}

	logger.Log.WithFields(logrus.Fields{
		"Operation :": "finished FindByToken",
		"user_id :":   refreshToken.UserID,
	}).Infoln("Finished FindByToken")
	return refreshToken, nil
}

func (r *refreshTokenRepository) DeleteByID(ctx context.Context, id uint) error {
	logger.Log.WithFields(logrus.Fields{
		"Operation :": "starting Delete token by id",
	}).Debug("Starting Delete token by id")
	result := r.db.WithContext(ctx).Where("id = ?", id).Delete(&models.RefreshToken{})
	if result.Error != nil {
		logger.Log.WithFields(logrus.Fields{
			"Operation :": "starting Delete token by id",
			"Error":       result.Error.Error(),
		}).Error("Error while deleting a refresh token")
		return result.Error
	}

	logger.Log.WithFields(logrus.Fields{
		"Operation :": "finished Delete token by id",
		"user_id :":   id,
	}).Info("Finished Delete token by id")
	return nil
}
