package models

import "gorm.io/gorm"

type Payment struct {
	gorm.Model
	ID                   uint    `gorm:"primary_key"`
	UserId               uint    `gorm:"not null"`
	Amount               float64 `gorm:"not null"`
	Currency             string  `gorm:"not null;default:'IRR'"`
	Status               string  `gorm:"not null;default:'PENDING'"`
	GateWayTransactionID string
}
