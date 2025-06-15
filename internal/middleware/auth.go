package middleware

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/rasadov/MailManagerApp/internal/services"
	"github.com/rasadov/MailManagerApp/pkg/utils"
	"log"
)

func AuthRequired(jwtService services.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authCookie, err := c.Cookie("token")
		if err != nil || authCookie == "" {
			c.Set("userID", "")
			c.Next()
			return
		}

		token, err := jwtService.VerifyToken(authCookie)
		if err != nil {
			c.Set("userID", "")
			c.Next()
			return
		}

		if claims, ok := token.Claims.(*utils.JWTCustomClaims); ok && token.Valid {
			log.Println("Successfully validated token")
			log.Println("User ID from token:", claims.UserID)

			c.Set("userID", claims.UserID)
			ctxWithVal := context.WithValue(c.Request.Context(), "userID", claims.UserID)
			c.Request = c.Request.WithContext(ctxWithVal)
		} else {
			c.Set("userID", "")
		}

		c.Next()
	}
}
