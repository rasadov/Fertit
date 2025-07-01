package utils

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/rasadov/MailManagerApp/internal/config"
	"log"
	"time"
)

func VerifyToken(encodedToken string) (*jwt.Token, error) {
	token, err := jwt.ParseWithClaims(
		encodedToken,
		&JWTCustomClaims{},
		func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}
			return []byte(config.JWTSecret), nil
		},
	)
	if err != nil {
		log.Println("There was a problem verifying the token")
		log.Println(err)
		return nil, err
	}

	if claims, ok := token.Claims.(*JWTCustomClaims); ok && token.Valid {
		if claims.Issuer != config.JWTIssuer {
			return nil, errors.New("invalid Issuer in token")
		}
		return token, nil
	}

	log.Println("Invalid token")
	return nil, errors.New("invalid token claims")
}

func GenerateToken(username string) (string, error) {
	tokenData := &JWTCustomClaims{
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    config.JWTIssuer,
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, tokenData)

	tokenString, err := token.SignedString([]byte(config.JWTSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
