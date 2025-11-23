package utils

import "golang.org/x/crypto/bcrypt"

func HashPassword(pass string) string {
	bytePass := []byte(pass)
	hPass, _ := bcrypt.GenerateFromPassword(bytePass, bcrypt.DefaultCost)
	return string(hPass)
}

func ComparePassword(dbpass, pass string) bool {
	return bcrypt.CompareHashAndPassword([]byte(dbpass), []byte(pass)) == nil
}
