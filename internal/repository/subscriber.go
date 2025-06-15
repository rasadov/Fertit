package repository

import (
	"github.com/rasadov/MailManagerApp/internal/models"
	"gorm.io/gorm"
)

type SubscriberRepository interface {
	GetSubscriber(uuid string) (models.Subscriber, error)
	ListSubscribers(offset, limit int) ([]models.Subscriber, error)
	GetSubscriberByEmail(email string) (models.Subscriber, error)
	UpdateSubscriber(subscriber models.Subscriber) error
	CreateSubscriber(subscriber models.Subscriber) error
}

type subscriberRepository struct {
	db *gorm.DB
}

func NewSubscriberRepository(db *gorm.DB) SubscriberRepository {
	return &subscriberRepository{
		db: db,
	}
}

func (r subscriberRepository) GetSubscriber(uuid string) (models.Subscriber, error) {
	var sub models.Subscriber
	res := r.db.First(&sub, "uuid = ?", uuid)

	return sub, res.Error
}

func (r subscriberRepository) GetSubscriberByEmail(email string) (models.Subscriber, error) {
	var sub models.Subscriber
	res := r.db.First(&sub, "email = ?", email)

	return sub, res.Error
}

func (r subscriberRepository) ListSubscribers(offset, limit int) ([]models.Subscriber, error) {
	var subs []models.Subscriber
	res := r.db.Offset(offset).Limit(limit).Find(&subs)
	return subs, res.Error
}

func (r subscriberRepository) UpdateSubscriber(subscriber models.Subscriber) error {
	res := r.db.Save(&subscriber)

	if res.Error != nil {
		return res.Error
	}

	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

func (r subscriberRepository) CreateSubscriber(subscriber models.Subscriber) error {
	res := r.db.Create(&subscriber)

	if res.Error != nil {
		return res.Error
	}

	return nil
}
