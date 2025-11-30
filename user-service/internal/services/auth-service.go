package services

import (
	"context"
	"errors"
	"strconv"
	"time"

	usr "user-service/internal/domain/user"
	"user-service/internal/models"
	"user-service/internal/pkg/events"
	"user-service/internal/repositories"
	"user-service/pkg/logger"
	"user-service/pkg/utils"

	"github.com/sirupsen/logrus"
)

type AuthService struct {
	userRepo           repositories.UserRepository
	refreshTokenRepo   repositories.RefreshTokenRepository
	jwtSecret          string
	refreshTokenSecret string
	accessTokenExpiry  time.Duration
	refreshTokenExpiry time.Duration
	eventDispatcher    events.EventDispatcher
}

func NewAuthService(
	userRepo repositories.UserRepository,
	refreshTokenRepo repositories.RefreshTokenRepository,
	jwtSecret string,
	refreshTokenSecret string,
	accessTokenExpiry time.Duration,
	refreshTokenExpiry time.Duration,
	eventDispatcher events.EventDispatcher,
) *AuthService {
	return &AuthService{
		userRepo:           userRepo,
		refreshTokenRepo:   refreshTokenRepo,
		jwtSecret:          jwtSecret,
		refreshTokenSecret: refreshTokenSecret,
		accessTokenExpiry:  accessTokenExpiry,
		refreshTokenExpiry: refreshTokenExpiry,
		eventDispatcher:    eventDispatcher,
	}
}

func (a *AuthService) SignUp(ctx context.Context, username, email, password string) (*models.UserResponse, error) {
	logger.Log.WithFields(logrus.Fields{
		"Operation :": "SignUp",
		"Username":    username,
		"Email":       email,
	}).Info("starting User registration")

	//check if user with this email already exists
	if _, err := a.userRepo.FindByEmail(ctx, email); err == nil {
		logger.Log.WithFields(logrus.Fields{
			"Operation :": "SignUp",
			"email":       email,
		}).Warn("registration Failed -> email already exists")
		return nil, errors.New("user with this email already exist")
	}

	//	check if user with username already exists
	if _, err := a.userRepo.FindByUserName(ctx, username); err == nil {
		logger.Log.WithFields(logrus.Fields{
			"Operation :": "SignUp",
			"username":    username,
		}).Warn("Registration Failed -> user with this username already exists")
		return nil, errors.New("user with this username already exist")
	}

	//hash password
	hashedPassword := utils.HashPassword(password)

	//	create new user
	user := &models.User{
		UserName:     username,
		Email:        email,
		PasswordHash: hashedPassword,
	}

	// Save user to database
	if err := a.userRepo.Create(ctx, user); err != nil {
		logger.Log.WithFields(logrus.Fields{
			"Operation :": "SignUp",
			"username":    username,
			"email":       email,
			"Error":       err.Error(),
		}).Error("Failed to create in database")
		return nil, err
	}

	//Convert To Response
	response := &models.UserResponse{
		ID:        user.ID,
		UserName:  user.UserName,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
	}

	logger.Log.WithFields(logrus.Fields{
		"Operation :": "SignUp",
		"username":    username,
		"email":       email,
		"user_id":     user.ID,
	}).Info("user registered successfully")
	event := usr.NewUserRegisteredEvent(user.ID, username, email)
	err := a.eventDispatcher.Dispatch(ctx, event)
	if err != nil {
		logger.Log.WithFields(logrus.Fields{
			"Operation :": "SignUp",
			"user-ID":     user.ID,
			"Error":       err.Error(),
		}).Warn("Failed to dispatch user registered event")
	}
	return response, nil

}

func (a *AuthService) SignIn(ctx context.Context, username, password string) (*models.SignInResponse, error) {
	logger.Log.WithFields(logrus.Fields{
		"Operation :": "SignIn",
		"Username":    username,
	}).Info("starting User login")

	// Find user by username
	user, err := a.userRepo.FindByUserName(ctx, username)
	if err != nil {
		logger.Log.WithFields(logrus.Fields{
			"Operation :": "SignIn",
			"username :":  username,
			"Error :":     "user not found",
		}).Warn("SignIn failed -> user not found")
		return nil, errors.New("invalid credentials")
	}

	//	Check password
	if !utils.ComparePassword(user.PasswordHash, password) {
		logger.Log.WithFields(logrus.Fields{
			"Operation :": "SignIn",
			"username :":  username,
			"user-id":     user.ID,
			"error":       "invalid password",
		}).Warn("SignIn failed -> invalid password")
		return nil, errors.New("invalid credentials")
	}

	//	Generate access token
	accessToken, accessExp, err := utils.GenerateToken(a.jwtSecret, strconv.Itoa(int(user.ID)), a.accessTokenExpiry)
	if err != nil {
		logger.Log.WithFields(logrus.Fields{
			"Operation :": "SignIn",
			"User-Id":     user.ID,
			"Error":       err.Error(),
		}).Error("Failed to generate access token")
		return nil, err
	}

	//	Generate refresh token
	refreshToken, refreshExp, err := utils.GenerateToken(a.refreshTokenSecret, strconv.Itoa(int(user.ID)), a.refreshTokenExpiry)
	if err != nil {
		logger.Log.WithFields(logrus.Fields{
			"Operation :": "SignIn",
			"User-Id":     user.ID,
			"Error":       err.Error(),
		}).Error("Failed to generate refresh token")
		return nil, err
	}

	// Save refresh token to database
	refreshTokenModel := models.RefreshToken{
		UserID:    user.ID,
		Token:     refreshToken,
		ExpiresAt: refreshExp,
		CreatedAt: time.Now(),
	}

	err = a.refreshTokenRepo.Create(ctx, &refreshTokenModel)
	if err != nil {
		logger.Log.WithFields(logrus.Fields{
			"Operation :": "SignIn",
			"User-Id":     user.ID,
			"Error":       err.Error(),
		}).Error("failed to save refresh token")
		return nil, err
	}

	//  Create Response
	response := &models.SignInResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    accessExp,
		UserID:       user.ID,
	}

	usersignedInEvent := usr.NewUserSignedInEvent(user.ID, username)
	err = a.eventDispatcher.Dispatch(ctx, usersignedInEvent)
	if err != nil {
		logger.Log.WithFields(logrus.Fields{
			"Operation :": "SignIn",
			"User-Id":     user.ID,
			"Error":       err.Error(),
		}).Warn("Failed to dispatch user signed in event")
	}

	logger.Log.WithFields(logrus.Fields{
		"Operation :": "SignIn",
		"user-id":     user.ID,
		"user-name":   user.UserName,
	}).Info("user signed in successfully")

	return response, nil
}

func (a *AuthService) ValidateToken(ctx context.Context, token string) (*models.ValidateResponse, error) {
	logger.Log.WithFields(logrus.Fields{
		"Operation :": "ValidateToken",
	}).Debug("Starting token validation")

	claims, err := utils.ValidateToken(token, a.jwtSecret)
	if err != nil {
		logger.Log.WithFields(logrus.Fields{
			"Operation :": "ValidateToken",
			"Error :":     err.Error(),
		}).Warn("Failed to validate token")
		return &models.ValidateResponse{Valid: false, Message: "invalid token"}, nil
	}

	// Find user by ID
	claimsuint, err := strconv.ParseUint(claims.Id, 10, 32)
	if err != nil {
		return nil, err
	}
	user, err := a.userRepo.FindByID(ctx, claimsuint)
	if err != nil {
		logger.Log.WithFields(logrus.Fields{
			"Operation :": "ValidateToken",
			"user-id":     user.ID,
			"error :":     "user not found",
		}).Warn("token validation failed -> user not found")
		return &models.ValidateResponse{Valid: false, Message: "user not found"}, nil
	}

	logger.Log.WithFields(logrus.Fields{
		"Operation :": "ValidateToken",
		"user-id":     user.ID,
		"user-name":   user.UserName,
	}).Info("Token validation successfully")

	return &models.ValidateResponse{
		Valid:   true,
		UserID:  user.ID,
		Message: "token is valid",
	}, nil
}

func (a *AuthService) RefreshToken(ctx context.Context, refreshToken string) (*models.RefreshTokenResponse, error) {
	logger.Log.WithFields(logrus.Fields{
		"Operation :": "RefreshToken",
	}).Info("Starting token refresh")

	claims, err := utils.ValidateToken(refreshToken, a.refreshTokenSecret)
	if err != nil {
		logger.Log.WithFields(logrus.Fields{
			"Operation :": "RefreshToken",
			"Error :":     err.Error(),
		}).Warn("refresh token validation failed")
		return nil, errors.New("invalid refresh token")
	}

	//Find the refresh Token in database
	token, err := a.refreshTokenRepo.FindByToken(ctx, refreshToken)
	if err != nil {
		logger.Log.WithFields(logrus.Fields{
			"Operation :": "RefreshToken",
			"user-id":     claims.Id,
			"error :":     "token not found in database",
		}).Warn("refresh token not found in database")
		return nil, errors.New("invalid refresh token")
	}

	//Check if token is expired
	if time.Now().After(token.ExpiresAt) {
		logger.Log.WithFields(logrus.Fields{
			"Operation :": "RefreshToken",
			"user-id":     claims.Id,
			"token-id":    token.ID,
		}).Warn("refresh token expired")
		return nil, errors.New("refresh token expired")
	}

	// Generate new access token
	accessToken, accessExp, err := utils.GenerateToken(a.jwtSecret, claims.Id, a.accessTokenExpiry)
	if err != nil {
		logger.Log.WithFields(logrus.Fields{
			"Operation :": "RefreshToken",
			"user-id":     claims.Id,
			"error :":     err.Error(),
		}).Error("failed to generate new access token")
		return nil, err
	}

	logger.Log.WithFields(logrus.Fields{
		"Operation :": "RefreshToken",
		"user-id":     claims.Id,
	}).Info("Token refreshed successfully")

	return &models.RefreshTokenResponse{
		AccessToken: accessToken,
		ExpiresAt:   accessExp,
	}, nil

}
