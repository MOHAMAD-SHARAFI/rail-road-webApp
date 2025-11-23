package utils

import (
	"fmt"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt"
)

type JWTClaims struct {
	jwt.StandardClaims
	UserID uint
	ID     uint
}

func GenerateToken(secret string, userID uint, exp time.Duration) (string, time.Time, error) {
	expTime := time.Now().Add(exp)

	claims := &JWTClaims{
		UserID: userID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expTime.Unix(),
			IssuedAt:  time.Now().Unix(),
			Subject:   strconv.FormatUint(uint64(userID), 10), // convert userID to string
		},
	}

	//	Create Token With Claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", time.Time{}, err
	}

	return tokenString, expTime, nil
}

func ValidateToken(tokenString string, secret string) (uint, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		//	Ensure The Singing Method Is As Expected
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})

	if err != nil {
		return 0, fmt.Errorf("error parsing token: %w", err)
	}

	//	Validate Token And Extract Claims
	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims.ID, nil
	}

	return 0, fmt.Errorf("invalid token")
}
