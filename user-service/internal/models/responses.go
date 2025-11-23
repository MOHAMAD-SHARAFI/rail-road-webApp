package models

import "time"

type UserResponse struct {
	ID          uint      `json:"id"`
	UserName    string    `json:"username"`
	Email       string    `json:"email"`
	PhoneNumber string    `json:"phone_number"`
	CreatedAt   time.Time `json:"created_at"`
}

type SignInResponse struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	ExpiresAt    time.Time `json:"expires_at"`
	UserID       uint      `json:"user_id"`
}

type ValidateResponse struct {
	Valid     bool      `json:"valid"`
	UserID    uint      `json:"user_id"`
	Message   string    `json:"message"`
	ExpiresAt time.Time `json:"expires_at"`
}

type RefreshTokenResponse struct {
	AccessToken string    `json:"access_token"`
	ExpiresAt   time.Time `json:"expires_at"`
}

type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Details any    `json:"details"`
}
