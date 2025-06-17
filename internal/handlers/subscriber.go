package handlers

import (
	customErrors "github.com/rasadov/MailManagerApp/pkg/errors"

	"errors"
	"github.com/gin-gonic/gin"
	"github.com/rasadov/MailManagerApp/internal/services"
	"net/http"
)

type SubscriberHandler struct {
	service services.SubscriberService
}

func NewSubscriberHandler(
	service services.SubscriberService,
) *SubscriberHandler {
	return &SubscriberHandler{
		service: service,
	}
}

func (s *SubscriberHandler) Subscribe(c *gin.Context) {
	email := c.PostForm("email")

	err := s.service.Subscribe(email)

	if errors.Is(err, customErrors.ErrSubscribing) {
		c.AbortWithStatusJSON(http.StatusServiceUnavailable, customErrors.ErrSubscribing)
		return
	}

	if errors.Is(err, customErrors.ErrAlreadySubscribed) {
		c.AbortWithStatusJSON(http.StatusBadRequest, customErrors.ErrAlreadySubscribed)
		return
	}

	c.AbortWithStatusJSON(http.StatusNoContent, nil)
}

func (s *SubscriberHandler) ManagePreferencesPage(c *gin.Context) {
	uuid := c.Param("uuid")

	subscriber, err := s.service.GetSubscriber(uuid)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"message": err,
		})
		return
	}

	c.HTML(http.StatusOK, "preferences.tmpl", gin.H{
		"subscriber": subscriber,
		"Year":       "2025",
	})
}

func (s *SubscriberHandler) UpdatePreferences(c *gin.Context) {
	uuid := c.PostForm("uuid")
	policyUpdates := c.PostForm("policy") == "on"
	incident := c.PostForm("incident") == "on"
	newFeatures := c.PostForm("features") == "on"
	news := c.PostForm("news") == "on"
	other := c.PostForm("other") == "on"

	err := s.service.Update(uuid, policyUpdates, incident, newFeatures, news, other)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": err,
		})
	}

	c.AbortWithStatus(http.StatusOK)
}
