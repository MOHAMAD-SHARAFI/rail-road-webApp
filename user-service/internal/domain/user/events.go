package user

import "user-service/internal/pkg/events"

type RegisteredEvent struct {
	events.BaseEvent
	UserID   uint
	UserName string
	Email    string
}

type SignedInEvent struct {
	events.BaseEvent
	UserID   uint
	UserName string
}

type PasswordResetRequestEvent struct {
	events.BaseEvent
	UserID     uint
	Email      string
	ResetToken string
}
type PasswordResetCompletedEvent struct {
	events.BaseEvent
	UserID uint
	Email  string
}

func NewUserRegisteredEvent(userID uint, userName, email string) RegisteredEvent {
	payLoad := map[string]interface{}{
		"user_id":   userID,
		"user_name": userName,
		"email":     email,
	}
	baseEvent := events.NewBaseEvent("user.registered", payLoad)

	return RegisteredEvent{
		BaseEvent: baseEvent,
		UserID:    userID,
		UserName:  userName,
		Email:     email,
	}
}

func NewUserSignedInEvent(userID uint, userName string) SignedInEvent {
	payLoad := map[string]interface{}{
		"user_id":   userID,
		"user_name": userName,
	}

	baseEvent := events.NewBaseEvent("user.signed-in", payLoad)

	return SignedInEvent{
		BaseEvent: baseEvent,
		UserID:    userID,
		UserName:  userName,
	}
}

func NewPasswordResetRequestEvent(userId uint, email string, resetToken string) PasswordResetRequestEvent {
	payLoad := map[string]interface{}{
		"user_id":     userId,
		"email":       email,
		"reset_Token": resetToken,
	}

	baseEvent := events.NewBaseEvent("user.password_reset_request", payLoad)
	return PasswordResetRequestEvent{
		BaseEvent:  baseEvent,
		UserID:     userId,
		Email:      email,
		ResetToken: resetToken,
	}
}

func NewPasswordResetCompletedEvent(userID uint, email string) PasswordResetCompletedEvent {
	payLoad := map[string]interface{}{
		"user_id": userID,
		"email":   email,
	}
	baseEvent := events.NewBaseEvent("user.password_reset_completed", payLoad)
	return PasswordResetCompletedEvent{
		BaseEvent: baseEvent,
		UserID:    userID,
		Email:     email,
	}
}
