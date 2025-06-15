package services

import (
	"context"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/rasadov/MailManagerApp/pkg/utils"
	"gorm.io/gorm"
)

type AuthService interface {
	ValidateCredentials(ctx context.Context, username, password string) error
	GenerateToken(ctx context.Context, username string) (string, error)
	VerifyToken(token string) (*jwt.Token, error)
}

type authService struct {
	SecretKey string
	Issuer    string
	db        *gorm.DB
}

func NewJWTAuthService(SecretKey, Issuer string, db *gorm.DB) AuthService {
	return &authService{
		SecretKey: SecretKey,
		Issuer:    Issuer,
		db:        db,
	}
}

func (s *authService) VerifyToken(encodedToken string) (*jwt.Token, error) {
	token, err := jwt.ParseWithClaims(
		encodedToken,
		&utils.JWTCustomClaims{},
		func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}
			return []byte(s.SecretKey), nil
		},
	)
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*utils.JWTCustomClaims); ok && token.Valid {
		if claims.Issuer != s.Issuer {
			return nil, errors.New("invalid Issuer in token")
		}
		return token, nil
	}

	return nil, errors.New("invalid token claims")
}

func (s *authService) ValidateCredentials(ctx context.Context, username, password string) error {
	return nil
}

func (s *authService) GenerateToken(ctx context.Context, username string) (string, error) {
	return "", nil
}
