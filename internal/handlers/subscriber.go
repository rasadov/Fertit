package handlers

import (
	"github.com/rasadov/MailManagerApp/internal/services"
	"gorm.io/gorm"
)

type SubscriberHandler struct {
}

func NewSubscriberHandler(
	service services.SubscriberService,
	db *gorm.DB,
) *SubscriberHandler {
	return &SubscriberHandler{}
}
