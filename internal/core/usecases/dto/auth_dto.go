package dto

import "logistics-api/internal/core/domain"

type RegisterRequest struct {
	Email    string          `json:"email" validate:"required,email"`
	Password string          `json:"password" validate:"required,min=6"`
	Role     domain.UserRole `json:"role" validate:"required,oneof=client admin"`
}

type RegisterResponse struct {
	User  *UserResponse `json:"user"`
	Token string        `json:"token"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type LoginResponse struct {
	User  *UserResponse `json:"user"`
	Token string        `json:"token"`
}

type UserResponse struct {
	ID        string          `json:"id"`
	Email     string          `json:"email"`
	Role      domain.UserRole `json:"role"`
	CreatedAt string          `json:"created_at"`
}

func ToUserResponse(user *domain.User) *UserResponse {
	return &UserResponse{
		ID:        user.ID,
		Email:     user.Email,
		Role:      user.Role,
		CreatedAt: user.CreatedAt.Format("2006-01-02T15:04:05Z"),
	}
}
