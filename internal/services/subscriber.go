package services

import (
	"errors"
	"github.com/rasadov/MailManagerApp/internal/models"
	"github.com/rasadov/MailManagerApp/internal/repository"
	customErrors "github.com/rasadov/MailManagerApp/pkg/errors"
	"gorm.io/gorm"
)

type SubscriberService interface {
	GetSubscriber(uuid string) (models.Subscriber, error)
	GetSubscribers(page, elements int) ([]models.Subscriber, error)
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

func (s *subscriberService) GetSubscriber(uuid string) (models.Subscriber, error) {
	res, err := s.repo.GetSubscriber(uuid)
	if err != nil {
		return models.Subscriber{}, err
	}
	return res, nil
}

func (s *subscriberService) GetSubscribers(page, elements int) ([]models.Subscriber, error) {
	res, err := s.repo.ListSubscribers((page-1)*elements, elements)
	if err != nil {
		return []models.Subscriber{}, err
	}
	return res, nil
}

func (s *subscriberService) Subscribe(email string) error {
	_, err := s.repo.GetSubscriberByEmail(email)
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return customErrors.ErrAlreadySubscribed
	}
	err = s.repo.CreateSubscriber(models.Subscriber{
		Email:         email,
		PolicyUpdates: true,
		Incidents:     true,
		NewFeatures:   true,
		News:          true,
		Other:         true,
	})
	if err != nil {
		return customErrors.ErrSubscribing
	}
	return nil
}

func (s *subscriberService) Update(uuid string, policyUpdates, incidents, newFeatures, news, other bool) error {
	subscriber, err := s.repo.GetSubscriber(uuid)
	if err != nil {
		return err
	}

	err = s.repo.UpdateSubscriber(models.Subscriber{
		Uuid:          subscriber.Uuid,
		Email:         subscriber.Email,
		PolicyUpdates: policyUpdates,
		Incidents:     incidents,
		NewFeatures:   newFeatures,
		News:          news,
		Other:         other,
	})

	if err != nil {
		return err
	}

	return nil
}
