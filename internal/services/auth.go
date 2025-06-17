package services

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/rasadov/MailManagerApp/internal/models"
	"github.com/rasadov/MailManagerApp/internal/repository"
	"github.com/rasadov/MailManagerApp/pkg/utils"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"log"
)

type AuthService interface {
	ValidateCredentials(username, password string) error
	GenerateToken(username string) (string, error)
	VerifyToken(token string) (*jwt.Token, error)
	EnsureAdminUser(username, password string) error
}

type authService struct {
	SecretKey      string
	Issuer         string
	userRepository repository.UserRepository
}

func NewJWTAuthService(SecretKey, Issuer string, db *gorm.DB) AuthService {
	userRepository := repository.NewUserRepository(db)
	return &authService{
		SecretKey:      SecretKey,
		Issuer:         Issuer,
		userRepository: userRepository,
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

func (s *authService) ValidateCredentials(username, password string) error {
	user, err := s.userRepository.GetUser(username)

	if err != nil {
		return err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))

	if err != nil {
		return err
	}

	return nil
}

func (s *authService) GenerateToken(username string) (string, error) {
	_, err := s.userRepository.GetUser(username)
	if err != nil {
		return "", err
	}

	tokenData := &utils.JWTCustomClaims{
		Username: username,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, tokenData)

	tokenString, err := token.SignedString([]byte(s.SecretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (s *authService) EnsureAdminUser(username, password string) error {
	_, err := s.userRepository.GetUser(username)
	if err == nil {
		log.Println("Admin user already exists")
		return nil
	}

	err = s.userRepository.CreateUser(models.User{
		Username: username,
		Password: password,
	})

	if err != nil {
		log.Println("Failed to create admin user")
		return err
	}

	log.Println("Admin user created")
	return nil
}
