package aut

import (
	"time"
	"golang.org/x/crypto/bcrypt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {
	// Create the Claims
	now := time.Now().UTC()
	expire := now.Add(expiresIn)

	claims := jwt.RegisteredClaims{
		Issuer:    	"chirpy",
		IssuedAt:	jwt.NewNumericDate(now),
		ExpiresAt: 	jwt.NewNumericDate(expire),
		Subject:	userID.String(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	sign, err := token.SignedString(tokenSecret)
	if err != nil {
		return "", err
	}

	return sign, nil
}

//func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) 


func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hash), nil
}

func CheckPasswordHash(password, hash string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err
}