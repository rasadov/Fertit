package services

import (
	"errors"
	"github.com/rasadov/MailManagerApp/internal/models"
	"github.com/rasadov/MailManagerApp/internal/repository"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"log"
)

type AuthService interface {
	ValidateCredentials(username, password string) error
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
