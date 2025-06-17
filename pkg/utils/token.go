package utils

import "github.com/golang-jwt/jwt/v5"

type JWTCustomClaims struct {
	// Storing Token data
	Username string `json:"username"`
	jwt.RegisteredClaims
}
