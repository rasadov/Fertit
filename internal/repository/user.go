package repository

import (
	"github.com/rasadov/MailManagerApp/internal/models"
	"gorm.io/gorm"
)

type UserRepository interface {
	GetUser(user *models.User) error
	CreateUser(user *models.User) error
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{
		db: db,
	}
}

func (r *userRepository) GetUser(user *models.User) error {
	return r.db.First(user).Error
}

func (r *userRepository) CreateUser(user *models.User) error {
	return r.db.Create(user).Error
}
