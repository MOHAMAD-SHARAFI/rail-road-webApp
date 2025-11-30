package models

import (
	"gorm.io/gorm"
)

type Payment struct {
	gorm.Model
	UserID               uint    `gorm:"not null"`
	Amount               float64 `gorm:"not null"`
	Fee                  float64 `gorm:"not null"`
	Total                float64 `gorm:"not null"`
	Currency             string  `gorm:"not null;default:'IRR'"`
	Status               string  `gorm:"not null;default:'PENDING'"`
	GateWayTransactionID string
}

type CreatePaymentRequest struct {
	Amount      float64 `json:"amount"`
	Description string  `json:"description"`
}
type FeeStructure struct {
	ID         uint `gorm:"primary_key"`
	Percentage float64
	MinFee     float64
}
