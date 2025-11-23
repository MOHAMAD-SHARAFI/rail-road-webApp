package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	UserName     string `gorm:"uniqueIndex"`
	Email        string `gorm:"uniqueIndex;not null"`
	PasswordHash string
	PhoneNumber  string `gorm:"uniqueIndex;not null"`
}

type PasswordResetToken struct {
	ID        uint   `gorm:"primary_key"`
	UserID    uint   `gorm:"not null"`
	Token     string `gorm:"not null"`
	ExpiresAt time.Time
	CreatedAt time.Time
}
