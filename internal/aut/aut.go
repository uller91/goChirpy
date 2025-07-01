package aut

import (
	"time"
	"golang.org/x/crypto/bcrypt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"net/http"
	"strings"
	"errors"
	"crypto/rand"
	"encoding/hex"
	//"fmt"
)

func MakeRefreshToken() (string, error) {
	key := make([]byte, 32)
	rand.Read(key)
	hexKey := hex.EncodeToString(key)
	return hexKey, nil
}

func GetBearerToken(headers http.Header) (string, error) {
	line := headers.Get("Authorization")
	if line == "" {
		return "", errors.New("No token found in the Header")
	}

	splitToken := strings.Fields(line)
	if len(splitToken) < 2 || splitToken[0] != "Bearer" {
		return "", errors.New("Malformed authorization Header was sent")
	}

	return splitToken[1], nil
}


func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {
	// Create the Claims
	now := time.Now().UTC()
	expire := now.Add(expiresIn)

	claims := &jwt.RegisteredClaims{
		Issuer:    	"chirpy",
		IssuedAt:	jwt.NewNumericDate(now),
		ExpiresAt: 	jwt.NewNumericDate(expire),
		Subject:	userID.String(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	sign, err := token.SignedString([]byte(tokenSecret))
	if err != nil {
		return "", err
	}

	return sign, nil
}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	claims := &jwt.RegisteredClaims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return[]byte(tokenSecret), nil
	})
	if err != nil {
		return uuid.Nil, err
	}

	userId, err := token.Claims.GetSubject()
	if err != nil {
		return uuid.Nil, err
	}

	userUuid, err := uuid.Parse(userId)
	if err!= nil {
		return uuid.Nil, err
	} 

	return userUuid, nil
}


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