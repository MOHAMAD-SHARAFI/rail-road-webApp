package repositories

import (
	"context"
	"errors"
	"user-service/internal/models"
	"user-service/pkg/logger"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

func NewUSerRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r userRepository) Create(ctx context.Context, user *models.User) error {
	result := r.db.WithContext(ctx).Create(user)
	if result.Error != nil {
		logger.Log.WithFields(logrus.Fields{
			"Error :":       "Cannot Create User",
			"User Info :":   user,
			"Operation :":   "CreateUser",
			"ErrorDetail :": result.Error,
		}).Error("failed to create user")
		return result.Error
	}

	logger.Log.WithFields(logrus.Fields{
		"Info :":      "User Created Successfully",
		"Operation :": "CreateUser",
		"email :":     user.Email,
		"firstName :": user.UserName,
		"ID :":        user.ID,
	}).Info("User Created Successfully")
	return nil
}

func (r userRepository) FindByUserName(ctx context.Context, username string) (*models.User, error) {
	var user models.User
	result := r.db.WithContext(ctx).Where("username = ?", username).First(&user)
	if result.Error != nil {
		logger.Log.WithFields(logrus.Fields{
			"Error :":       "Cannot Find User",
			"User Info :":   username,
			"Operation :":   "FindByUserName",
			"ErrorDetail :": result.Error,
		}).Error("failed to find user by username")
		return nil, result.Error
	}

	logger.Log.WithFields(logrus.Fields{
		"Info :":      "User Find Successfully",
		"Operation :": "FindByUserName",
		"username :":  models.User{UserName: username},
	}).Info("User Find Successfully")

	return &user, nil
}

func (r userRepository) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	result := r.db.WithContext(ctx).Where("email = ?", email).First(&user)
	if result.Error != nil {
		logger.Log.WithFields(logrus.Fields{
			"Error :":       "Cannot Find User",
			"User Email :":  email,
			"Operation :":   "FindByEmail",
			"ErrorDetail :": result.Error,
		}).Error("failed to find user by email")
	}

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		logger.Log.WithFields(logrus.Fields{
			"Error :":      "Cannot Find User",
			"User Email :": email,
			"Operation :":  "FindByEmail",
		}).Error("failed to find user by email")
		return nil, result.Error
	}
	logger.Log.WithFields(logrus.Fields{
		"Info :":      "User Find Successfully",
		"Operation :": "FindByEmail",
		"email :":     email,
		"firstName :": user.UserName,
	}).Info("User Find Successfully")
	return &user, nil
}

func (r userRepository) FindByID(ctx context.Context, id uint) (*models.User, error) {
	var user models.User
	result := r.db.WithContext(ctx).First(&user, id)
	if result.Error != nil {
		logger.Log.WithFields(logrus.Fields{
			"Error :":       "Cannot Find User",
			"User ID :":     id,
			"Operation :":   "FindByID",
			"ErrorDetail :": result.Error,
		}).Error("failed to find user by id")
	}

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		logger.Log.WithFields(logrus.Fields{
			"Error :":     "Cannot Find User",
			"User ID :":   id,
			"Operation :": "FindByID",
		}).Error("failed to find user by id")
	}

	logger.Log.WithFields(logrus.Fields{
		"Info :":      "User Find Successfully",
		"Operation :": "FindByID",
		"ID :":        id,
		"UserName :":  user.UserName,
	}).Info("User Find Successfully")
	return &user, nil
}

func (r userRepository) Update(ctx context.Context, user *models.User) error {
	result := r.db.WithContext(ctx).Save(user)
	if result.Error != nil {
		logger.Log.WithFields(logrus.Fields{
			"Error :":       "Cannot Update User",
			"User Info :":   user,
			"Operation :":   "UpdateUser",
			"ErrorDetail :": result.Error,
		}).Error("failed to update user")
	}

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		logger.Log.WithFields(logrus.Fields{
			"Error :":     "Cannot Update User",
			"User Info :": user,
			"Operation :": "UpdateUser",
		}).Error("failed to update user")
	}

	logger.Log.WithFields(logrus.Fields{
		"Info :":      "User Updated Successfully",
		"Operation :": "UpdateUser",
		"email :":     user.Email,
		"UserName :":  user.UserName,
	}).Info("User Updated Successfully")
	return result.Error
}

func (r userRepository) Delete(ctx context.Context, id uint) error {
	var user models.User
	result := r.db.WithContext(ctx).Delete(&user, id)
	if result.Error != nil {
		logger.Log.WithFields(logrus.Fields{
			"Error :":       "Cannot Delete User",
			"User ID :":     id,
			"Operation :":   "Delete",
			"ErrorDetail :": result.Error,
		}).Error("failed to delete user")
		return result.Error
	}

	logger.Log.WithFields(logrus.Fields{
		"Info :":      "User Deleted Successfully",
		"Operation :": "Delete",
		"user-id :":   id,
	}).Info("User Deleted Successfully")
	return nil
}

var _ UserRepository = (*userRepository)(nil)
