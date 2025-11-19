package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	ID           uint   `gorm:"primary_key"`
	UserName     string `gorm:"uniqueIndex;not null"`
	Email        string `gorm:"uniqueIndex;not null"`
	PasswordHash string `gorm:"not null"`
	PhoneNumber  string `gorm:"uniqueIndex;not null"`
}
