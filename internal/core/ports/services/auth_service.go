package services

import (
	"logistics-api/internal/core/domain"
)

type AuthService interface {
	GenerateToken(user *domain.User) (string, error)
	ValidateToken(token string) (*TokenClaims, error)
	HashPassword(password string) (string, error)
	ValidatePassword(hashedPassword, password string) bool
}

type TokenClaims struct {
	UserID string          `json:"user_id"`
	Email  string          `json:"email"`
	Role   domain.UserRole `json:"role"`
}
