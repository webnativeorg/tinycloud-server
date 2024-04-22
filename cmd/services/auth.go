package services

import (
	"fmt"

	"github.com/golang-jwt/jwt/v5"
	"github.com/webnativeorg/tinycloud-server/cmd/environment"
	"golang.org/x/crypto/bcrypt"
)

func GenerateJWT(email string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": email,
	})

	return token.SignedString(environment.JWT_SECRET)
}

func ValidatePassword(password string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		fmt.Println("Password validation", err)
		return false
	}
	return true
}
