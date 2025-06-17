package repository

import (
	"github.com/rasadov/MailManagerApp/internal/models"
	"gorm.io/gorm"
)

type UserRepository interface {
	GetUser(username string) (models.User, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{
		db: db,
	}
}

func (r *userRepository) GetUser(username string) (models.User, error) {
	var admin models.User

	err := r.db.First(&admin, "username = ?", username).Error
	if err != nil {
		return models.User{}, err
	}
	return admin, nil
}
