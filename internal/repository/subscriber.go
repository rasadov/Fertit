package repository

import (
	"github.com/rasadov/MailManagerApp/internal/models"
	"github.com/rasadov/MailManagerApp/pkg/utils"
	"gorm.io/gorm"
)

// SubscriberRepository - interface with repository level functions
type SubscriberRepository interface {
	GetSubscriber(sub *models.Subscriber) error
	ListSubscribers(offset, limit int) ([]*models.Subscriber, int64, error)
	ListCategorySubscriberEmails(category string) ([]*utils.SubscriberEmail, int64, error)
	UpdateSubscriber(sub *models.Subscriber) error
	CreateSubscriber(sub *models.Subscriber) error
}

// subscriberRepository - implementation of the subscriberRepository with gorm
type subscriberRepository struct {
	db *gorm.DB
}

func NewSubscriberRepository(db *gorm.DB) SubscriberRepository {
	return &subscriberRepository{
		db: db,
	}
}

func (r subscriberRepository) GetSubscriber(sub *models.Subscriber) error {
	return r.db.First(sub).Error
}

func (r subscriberRepository) ListSubscribers(offset, limit int) ([]*models.Subscriber, int64, error) {
	var subs []*models.Subscriber
	res := r.db.Offset(offset).Limit(limit).Find(&subs)
	return subs, res.RowsAffected, res.Error
}

func (r subscriberRepository) ListCategorySubscriberEmails(category string) ([]*utils.SubscriberEmail, int64, error) {
	var subs []*utils.SubscriberEmail

	query := r.db.Model(&models.Subscriber{}).
		Select("email, uuid")

	if category != "all" {
		query = query.Where(category+" = ?", true)
	}

	res := query.Find(subs)
	return subs, res.RowsAffected, res.Error
}

func (r subscriberRepository) UpdateSubscriber(subscriber *models.Subscriber) error {
	res := r.db.Save(subscriber)

	if res.Error != nil {
		return res.Error
	}

	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

func (r subscriberRepository) CreateSubscriber(subscriber *models.Subscriber) error {
	res := r.db.Create(subscriber)

	if res.Error != nil {
		return res.Error
	}

	return nil
}
