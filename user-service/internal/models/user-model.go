package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	FirstName    string `gorm:"uniqueIndex;not null"`
	LastName     string `gorm:"uniqueIndex;not null"`
	Email        string `gorm:"uniqueIndex;not null"`
	PasswordHash string `gorm:"not null"`
	PhoneNumber  string `gorm:"uniqueIndex;not null"`
}

type PassworResetToken struct {
	ID        uint   `gorm:"primary_key"`
	UserID    uint   `gorm:"not null"`
	Token     string `gorm:"not null"`
	ExpiresAt time.Time
	CreatedAt time.Time
}
