package handlers

import (
	"github.com/gin-gonic/gin"
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

func (s *SubscriberHandler) Subscribe(c *gin.Context) {}

func (s *SubscriberHandler) ManagePreferences(c *gin.Context) {}

func (s *SubscriberHandler) UpdatePreferences(c *gin.Context) {}
