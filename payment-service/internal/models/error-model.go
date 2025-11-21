package models

import (
	"errors"
)

var (
	ErrFailedConnectDB     = errors.New("cannot connect to database")
	ErrFailedMigrateDB     = errors.New("cannot migrate database")
	ErrFailedCreatePayment = errors.New("could not create payment")
	ErrUpdateFailed        = errors.New("could not update payment")
	ErrPaymentNotFound     = errors.New("payment not found")
	ErrCantGetFeeStructure = errors.New("can't get fee structure")
)
