package main

import (
	"context"
	"github.com/rasadov/MailManagerApp/internal/database/postgres"
	"github.com/rasadov/MailManagerApp/internal/database/redis"
	"github.com/rasadov/MailManagerApp/internal/handlers"
	"github.com/rasadov/MailManagerApp/internal/middleware"
	"github.com/rasadov/MailManagerApp/internal/services"
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize databases
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	db := postgres.InitDB()
	redis.InitRedis(ctx)

	// Initialize services
	rateLimiter := services.NewRedisRateLimiter(redis.Client, 3, 15*time.Minute)
	authService := services.NewJWTAuthService(os.Getenv("JWT_SECRET"), os.Getenv("JWT_ISSUER"), db)
	emailService := services.NewSMTPEmailService(
		os.Getenv("SMTP_HOST"),
		587,
		os.Getenv("SMTP_USER"),
		os.Getenv("SMTP_PASS"),
	)
	subscriberService := services.NewSubscriberService(db)

	// Initialize handlers
	adminHandler := handlers.NewAdminHandler(rateLimiter, authService, emailService)
	subscriberHandler := handlers.NewSubscriberHandler(subscriberService)
	staticHandler := handlers.NewStaticHandler()

	// Setup routes
	r := setupRoutes(adminHandler, subscriberHandler, staticHandler, rateLimiter, authService)

	log.Fatal(r.Run(":8080"))
}

func setupRoutes(
	adminHandler *handlers.AdminHandler,
	subscriberHandler *handlers.SubscriberHandler,
	staticHandler *handlers.StaticHandler,
	rateLimiter *services.RedisRateLimiter,
	authService services.AuthService,
) *gin.Engine {
	r := gin.Default()

	// Load HTML templates
	r.LoadHTMLGlob("web/templates/**/*.html")

	r.Static("/static", "./web/static")

	// ========== STATIC ROUTES ==========
	r.GET("/favicon.ico", staticHandler.Favicon)
	r.GET("/", staticHandler.Index)

	// ========== SUBSCRIBER ROUTES ==========
	r.POST("/subscribe", subscriberHandler.Subscribe)
	r.GET("/preferences/:uuid", subscriberHandler.ManagePreferencesPage)
	r.POST("/preferences/:uuid", subscriberHandler.UpdatePreferences)

	// ========== ADMIN ROUTES ==========
	// Admin login (with an obscured path)
	r.GET("/admin/xxxloginyyy", adminHandler.AdminLoginPage)
	r.POST("/admin/xxxloginyyy",
		middleware.RateLimitMiddleware(rateLimiter),
		adminHandler.LoginPost,
	)

	// Protected admin routes
	admin := r.Group("/admin")
	admin.Use(middleware.AuthRequired(authService))
	{
		admin.GET("", adminHandler.AdminDashboard)
		admin.POST("/send-email", adminHandler.SendEmail)
	}

	return r
}
