package models

import (
	"time"
)

type PasswordResetToken struct {
	ID        uint   `gorm:"primary_key"`
	UserID    uint   `gorm:"not null"`
	Token     string `gorm:"not null"`
	ExpiresAt time.Time
	CreatedAt time.Time
}
