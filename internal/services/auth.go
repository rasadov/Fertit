package services

import "context"

type AuthService interface {
	ValidateCredentials(ctx context.Context, username, password string) error
	GenerateToken(ctx context.Context, username string) (string, error)
	VerifyToken(ctx context.Context, token string) (string, error)
}
