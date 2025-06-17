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
	"time"
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
		log.Println("There was a problem verifying the token")
		log.Println(err)
		return nil, err
	}

	if claims, ok := token.Claims.(*utils.JWTCustomClaims); ok && token.Valid {
		if claims.Issuer != s.Issuer {
			return nil, errors.New("invalid Issuer in token")
		}
		return token, nil
	}

	log.Println("Invalid token")
	return nil, errors.New("invalid token claims")
}

func (s *authService) ValidateCredentials(username, password string) error {
	user, err := s.userRepository.GetUser(username)

	log.Println("Credentials:")
	log.Println(username, password)

	if err != nil {
		return err
	}

	if len(password) < 8 || len(password) > 128 {
		return errors.New("password length must be between 8 and 128")
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
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    s.Issuer,
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, tokenData)

	tokenString, err := token.SignedString([]byte(s.SecretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (s *authService) EnsureAdminUser(username, password string) error {
	if len(password) < 8 || len(password) > 128 {
		return errors.New("password length must be between 8 and 128")
	}

	_, err := s.userRepository.GetUser(username)
	if err == nil {
		log.Println("Admin user already exists")
		return nil
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		return err
	}

	err = s.userRepository.CreateUser(models.User{
		Username: username,
		Password: string(hashedPassword),
	})

	if err != nil {
		log.Println("Failed to create admin user")
		return err
	}

	log.Println("Admin user created")
	return nil
}
