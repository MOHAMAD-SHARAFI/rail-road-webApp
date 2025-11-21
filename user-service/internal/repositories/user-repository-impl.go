package repositories

import (
	"context"
	"errors"
	"log"
	"time"
	"user-service/internal/models"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"user-service/internal/logger"
)

type userRepository struct {
	db *gorm.DB
}

func NewUSerRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r userRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	if email == "" {
		return nil, errors.New("email must be filled")
	}

	err := r.db.WithContext(ctx).Where("email= ?", models.User{Email: email}).First(&user)
	if err != nil {
		log.Println("cannot find user with this Email", err)
	}

	return &user, nil
}

func (r userRepository) Create(ctx context.Context, user *models.User) error {
	err := r.db.WithContext(ctx).Create(user)
	if err != nil {
		logger.Log.WithFields(logrus.Fields{
			"email":     user.Email,
			"firstName": user.FirstName,
			"lastName":  user.LastName,
			"ID":        user.ID,
		})
	}

	return nil
}

func (r userRepository) Update(ctx context.Context, user *models.User) error {
	err := r.db.WithContext(ctx).Save(&user)
	if err != nil {
		logger.Log.WithFields(logrus.Fields{
			"Operation":   "Update",
			"Error":       "cannot update user",
			"ErrorDetail": err,
		})
	} else {
		logger.Log.WithFields(logrus.Fields{
			"Operation": "Update",
			"UpdatedFields": logrus.Fields{
				"email":        user.Email,
				"firstName":    user.FirstName,
				"lastName":     user.LastName,
				"ID":           user.ID,
				"PhoneNumber ": user.PhoneNumber,
			},
		})
	}

	return nil
}

func (r userRepository) GetByID(ctx context.Context, id *uint) (*models.User, error) {
	if *id == nil || id == 0 {
		err := errors.New("id cannot be 0")
		logger.Log.WithFields(logrus.Fields{
			"Operation":   "GetByID",
			"Error":       "cannot find user because id is zero or negative",
			"ErrorDetail": err,
		})

	}

	var user models.User
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&user)
	if err != nil {
		logger.Log.WithFields(logrus.Fields{
			"Operation":   "GetByID",
			"ID":          id,
			"Error":       "cannot find user",
			"ErrorDetail": err,
		})
	} else {
		logger.Log.WithFields(logrus.Fields{
			"Operation":   "GetByID",
			"ID":          id,
			"FirstName":   user.FirstName,
			"LastName":    user.LastName,
			"PhoneNumber": user.PhoneNumber,
		})
	}

	return &user, nil
}

func (r userRepository) GetByPhoneNumber(ctx context.Context, phone string) (*models.User, error) {
	var user models.User
	if phone == "" || phone == "0" {
		return nil, logger.Log.WithFields(logrus.Fields{
			"Operation":   "GetByPhoneNumber",
			"PhoneNumber": phone,
			"Error":       "cannot find user Because it's phone number is zero or empty",
		})
	}

	err := r.db.WithContext(ctx).Where("phone = ?", models.User{PhoneNumber: phone}).First(&user)
	if err != nil {
		logger.Log.WithFields(logrus.Fields{
			"Operation":   "GetByPhoneNumber",
			"PhoneNumber": phone,
			"Error":       "cannot find user",
			"ErrorDetail": err,
		})
	} else {
		logger.Log.WithFields(logrus.Fields{
			"Operation":   "GetByPhoneNumber",
			"PhoneNumber": phone,
			"FirstName":   user.FirstName,
			"LastName":    user.LastName,
		})
	}

	return &user, nil
}

func (r userRepository) GetValidPasswordResetToken(ctx context.Context, userID uint, token string) (*models.PassworResetToken, error) {
	err := r.db.WithContext(ctx).Where("user_id = ? AND token=? AND expires_at > ?", userID, token, time.Now()).First(&token)
	if err != nil {
		log.Fatalf("your token,claims is not valid %v \n", models.ErrTokenNotValid)
	}

	return &models.PassworResetToken{Token: token}, nil
}

func (r userRepository) CreatePasswordResetToken(ctx context.Context, token *models.PassworResetToken) error {
	err := r.db.WithContext(ctx).Create(token)
	if err != nil {
		logger.Log.WithFields(logrus.Fields{
			"user_id":   token.UserID,
			"operation": "CreatePasswordResetToken",
			"Error":     "failed to create token",
		}).Error(err)
	} else {
		logger.Log.WithFields(logrus.Fields{
			"user_id":   token.UserID,
			"Token":     token.Token,
			"Operation": "CreatePasswordResetToken",
			"CreatedAt": token.CreatedAt,
		})
	}

	return nil
}

func (r userRepository) DeletePasswordResetToken(ctx context.Context, tokenId uint) error {
	err := r.db.WithContext(ctx).Where("token_id = ?", tokenId).Delete(&models.PassworResetToken{})
	if err != nil {
		logger.Log.WithFields(logrus.Fields{
			"Operation": "DeletePasswordResetToken",
			"Error":     "failed to delete token",
		})
	} else {
		logger.Log.WithFields(logrus.Fields{
			"Operation": "DeletePasswordResetToken",
			"Log":       "token succesfully deleted",
		})
	}

	return nil
}
