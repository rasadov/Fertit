package middleware

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/rasadov/MailManagerApp/internal/services"
	"github.com/rasadov/MailManagerApp/pkg/utils"
	"log"
	"net/http"
)

func AuthRequired(jwtService services.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authCookie, err := c.Cookie("token")
		log.Println("authCookie", authCookie)
		if err != nil || authCookie == "" {
			log.Println("We aborted here in if err != nil block 1")
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		token, err := jwtService.VerifyToken(authCookie)
		if err != nil {
			log.Println("We aborted here in if err != nil block")
			log.Println(err)
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		if claims, ok := token.Claims.(*utils.JWTCustomClaims); ok && token.Valid {
			log.Println("Successfully validated token")
			log.Println("User ID from token:", claims.Username)

			c.Set("username", claims.Username)
			ctxWithVal := context.WithValue(c.Request.Context(), "username", claims.Username)
			c.Request = c.Request.WithContext(ctxWithVal)
		} else {
			log.Println("We aborted here in else block")
			log.Println("User ID from token:", claims.Username)
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		c.Next()
	}
}
