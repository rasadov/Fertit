package services

import "gorm.io/gorm"

type SubscriberService interface {
}

func NewSubscriberService(db *gorm.DB) SubscriberService {
	return nil
}
