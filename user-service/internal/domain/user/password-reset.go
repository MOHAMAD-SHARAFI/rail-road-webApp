package user

import (
	"crypto/rand"
	"encoding/hex"
	"time"
)

type ResetToken struct {
	Token    string
	ExpireAt time.Time
}

func (t ResetToken) IsExpired() bool {
	currenTime := time.Now()
	isExpired := currenTime.After(t.ExpireAt)
	return isExpired
}

func (t ResetToken) IsValid() bool {
	tokenNotEmpty := t.Token != ""
	tokenNotExpired := !t.IsExpired()
	return tokenNotEmpty && tokenNotExpired
}

func (s *PasswordResetService) GenerateResetToken(expiryTime time.Duration) (*ResetToken, error) {
	tokenByte := make([]byte, 32)
	_, err := rand.Read(tokenByte)
	if err != nil {
		return nil, err
	}

	tokenString := hex.EncodeToString(tokenByte)
	ExpirationTime := time.Now().Add(expiryTime)

	resetToken := &ResetToken{
		Token:    tokenString,
		ExpireAt: ExpirationTime,
	}

	return resetToken, nil

}

type PasswordResetService struct{}

func (s *PasswordResetService) ValidateResetToken(token *ResetToken) bool {
	isValid := token.IsValid()
	return isValid
}
