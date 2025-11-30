package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	UserName     string `gorm:"uniqueIndex"`
	Email        string `gorm:"uniqueIndex;not null"`
	PasswordHash string
	PhoneNumber  string `gorm:"uniqueIndex;not null"`
}
