package utils

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
)

type JWTClaims struct {
	*jwt.StandardClaims
}

func GenerateToken(secret string, userID string, exp time.Duration) (tokenKey string, expTime time.Time, err error) {
	expTime = time.Now().Local().Add(exp)

	claims := &JWTClaims{
		&jwt.StandardClaims{
			ExpiresAt: expTime.Unix(),
			Id:        userID,
		},
	}

	//	Create Token With Claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	if tokenKey, err = token.SignedString([]byte(secret)); err != nil {
		return
	}

	return
}

func ValidateToken(tokenString string, secret string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		//	Ensure The Singing Method Is As Expected
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})

	if err != nil {
		return nil, fmt.Errorf("error parsing token: %w", err)
	}

	//	Validate Token And Extract Claims
	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}
