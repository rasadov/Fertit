package services

import (
	"errors"
	"github.com/rasadov/MailManagerApp/internal/models"
	"github.com/rasadov/MailManagerApp/internal/repository"
	customErrors "github.com/rasadov/MailManagerApp/pkg/errors"
	"github.com/rasadov/MailManagerApp/pkg/utils"
	"gorm.io/gorm"
)

type SubscriberService interface {
	GetSubscriber(uuid string) (*models.Subscriber, error)
	GetSubscribers(page, elements int) ([]*models.Subscriber, int64, error)
	GetSubscribersByCategory(category string) ([]*utils.SubscriberEmail, int64, error)
	Subscribe(email string) error
	Update(uuid string, policyUpdates, incidents, newFeatures, news, other bool) error
}

type subscriberService struct {
	repo repository.SubscriberRepository
}

func NewSubscriberService(db *gorm.DB) SubscriberService {
	repo := repository.NewSubscriberRepository(db)
	return &subscriberService{
		repo: repo,
	}
}

func (s *subscriberService) GetSubscriber(uuid string) (*models.Subscriber, error) {
	user := &models.Subscriber{
		Uuid: uuid,
	}
	err := s.repo.GetSubscriber(user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *subscriberService) GetSubscribers(page, elements int) ([]*models.Subscriber, int64, error) {
	users, total, err := s.repo.ListSubscribers((page-1)*elements, elements)
	if err != nil {
		return nil, 0, err
	}
	return users, total, nil
}

func (s *subscriberService) GetSubscribersByCategory(category string) ([]*utils.SubscriberEmail, int64, error) {
	emails, total, err := s.repo.ListCategorySubscriberEmails(category)
	if err != nil {
		return nil, 0, err
	}
	return emails, total, nil
}

func (s *subscriberService) Subscribe(email string) error {
	user := &models.Subscriber{
		Email: email,
	}
	err := s.repo.GetSubscriber(user)
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return customErrors.ErrAlreadySubscribed
	}
	err = s.repo.CreateSubscriber(&models.Subscriber{
		Email:         email,
		PolicyUpdates: true,
		Incidents:     true,
		NewFeatures:   true,
		News:          true,
		Others:        true,
	})
	if err != nil {
		return customErrors.ErrSubscribing
	}
	return nil
}

func (s *subscriberService) Update(uuid string, policyUpdates, incidents, newFeatures, news, other bool) error {
	subscriber := &models.Subscriber{
		Uuid: uuid,
	}
	err := s.repo.GetSubscriber(subscriber)
	if err != nil {
		return err
	}

	subscriber.PolicyUpdates = policyUpdates
	subscriber.Incidents = incidents
	subscriber.NewFeatures = newFeatures
	subscriber.News = news
	subscriber.Others = other

	err = s.repo.UpdateSubscriber(subscriber)

	if err != nil {
		return err
	}

	return nil
}
