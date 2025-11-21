package models

import (
	"errors"
)

var (
	ErrFailedConnectDB  = errors.New("cannot connect to database")
	ErrFailedMigrateDB  = errors.New("cannot migrate database")
	ErrUserNotFound     = errors.New("user not found")
	ErrCannotCreateUser = errors.New("cannot create user")
	ErrTokenNotValid    = errors.New("token is not valid")
)
