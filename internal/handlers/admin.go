package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/rasadov/MailManagerApp/internal/config"
	"github.com/rasadov/MailManagerApp/internal/services"
	customErrors "github.com/rasadov/MailManagerApp/pkg/errors"
	"github.com/rasadov/MailManagerApp/pkg/utils"
	"log"
	"net/http"
	"strconv"
)

type AdminHandler struct {
	rateLimiter       services.RateLimiter
	authService       services.AuthService
	emailService      services.EmailService
	subscriberService services.SubscriberService
}

func NewAdminHandler(
	rateLimiter services.RateLimiter,
	authService services.AuthService,
	emailService services.EmailService,
	subscriberService services.SubscriberService,
) *AdminHandler {
	return &AdminHandler{
		rateLimiter:       rateLimiter,
		authService:       authService,
		emailService:      emailService,
		subscriberService: subscriberService,
	}
}

func (h *AdminHandler) AdminLoginPage(c *gin.Context) {
	c.HTML(http.StatusOK, "login.tmpl", gin.H{
		"Year": "2025",
	})
}

func (h *AdminHandler) LoginPost(c *gin.Context) {
	ip := c.ClientIP()
	username := c.PostForm("username")
	password := c.PostForm("password")

	isBlocked, err := h.rateLimiter.IsBlocked(c, ip)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, err)
	}

	if isBlocked {
		c.AbortWithStatusJSON(http.StatusTooManyRequests, customErrors.ErrTooManyAttempts)
		return
	}

	if err = h.authService.ValidateCredentials(username, password); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "Invalid username or password",
		})

		err = h.rateLimiter.RecordFailedAttempt(c, ip)
		if err != nil {
			log.Printf("Error recording failed attempt: %v", err)
		}
		return
	}

	token, err := utils.GenerateToken(username)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "Error generating token",
		})
		return
	}

	c.SetSameSite(http.SameSiteStrictMode)
	c.SetCookie(
		"token",
		token,
		config.AppConfig.TokenExpirationSeconds,
		"/",
		"",
		!config.AppConfig.Debug,
		true)

	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "Successfully logged in",
	})
}

func (h *AdminHandler) AdminDashboard(c *gin.Context) {
	pageStr := c.DefaultQuery("page", "1")
	elementsStr := c.DefaultQuery("elements", "10")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 || page > 100 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "Invalid page number",
		})
		return
	}

	elements, err := strconv.Atoi(elementsStr)
	if err != nil || elements < 1 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "Invalid elements number",
		})
		return
	}

	if elements > 50 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "Invalid elements number",
		})
		return
	}

	subscribers, err := h.subscriberService.GetSubscribers(page, elements)
	if err != nil {
		log.Println("Error: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve subscribers",
		})
		return
	}

	c.HTML(http.StatusOK, "dashboard.tmpl", gin.H{
		"Page":        page,
		"Elements":    elements,
		"Subscribers": subscribers,
	})
}

func (h *AdminHandler) SendEmail(c *gin.Context) {
	category := c.PostForm("category")
	subject := c.PostForm("subject")
	body := c.PostForm("body")
	isHtml := c.PostForm("html") == "1"

	emails, err := h.subscriberService.GetSubscribersByCategory(category)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve subscribers",
		})
	}

	err = h.emailService.SendNewsletter(emails, subject, body, isHtml)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to send email",
		})
		return
	}

	c.HTML(http.StatusOK, "success.tmpl", gin.H{})
}
