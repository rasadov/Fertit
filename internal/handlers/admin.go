package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/rasadov/MailManagerApp/internal/config"
	"github.com/rasadov/MailManagerApp/internal/services"
	"log"
	"net/http"
)

type AdminHandler struct {
	rateLimiter  services.RateLimiter
	authService  services.AuthService
	emailService services.EmailService
}

func NewAdminHandler(
	rateLimiter services.RateLimiter,
	authService services.AuthService,
	emailService services.EmailService,
) *AdminHandler {
	return &AdminHandler{
		rateLimiter:  rateLimiter,
		authService:  authService,
		emailService: emailService,
	}
}

func (h *AdminHandler) AdminLoginPage(c *gin.Context) {}

func (h *AdminHandler) LoginPost(c *gin.Context) {
	ip := c.ClientIP()
	username := c.PostForm("username")
	password := c.PostForm("password")

	if err := h.authService.ValidateCredentials(c, username, password); err != nil {
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

	token, err := h.authService.GenerateToken(c, username)
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

func (h *AdminHandler) AdminDashboard(c *gin.Context) {}

func (h *AdminHandler) SendEmail(c *gin.Context) {}
