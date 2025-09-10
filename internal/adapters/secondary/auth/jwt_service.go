package auth

import (
	"errors"
	"time"

	"logistics-api/internal/core/domain"
	"logistics-api/internal/core/ports/services"

	"github.com/golang-jwt/jwt/v5"
)

type JWTService struct {
	secretKey  []byte
	expiryHour int
}

func NewJWTService(secretKey string, expiryHour int) *JWTService {
	return &JWTService{
		secretKey:  []byte(secretKey),
		expiryHour: expiryHour,
	}
}

func (j *JWTService) GenerateToken(user *domain.User) (string, error) {
	claims := jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"role":    user.Role,
		"exp":     time.Now().Add(time.Hour * time.Duration(j.expiryHour)).Unix(),
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.secretKey)
}

func (j *JWTService) ValidateToken(tokenString string) (*services.TokenClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return j.secretKey, nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid token claims")
	}

	userID, ok := claims["user_id"].(string)
	if !ok {
		return nil, errors.New("invalid user_id in token")
	}

	email, ok := claims["email"].(string)
	if !ok {
		return nil, errors.New("invalid email in token")
	}

	roleStr, ok := claims["role"].(string)
	if !ok {
		return nil, errors.New("invalid role in token")
	}

	role := domain.UserRole(roleStr)
	if role != domain.ClientRole && role != domain.AdminRole {
		return nil, errors.New("invalid role value")
	}

	return &services.TokenClaims{
		UserID: userID,
		Email:  email,
		Role:   role,
	}, nil
}

func (j *JWTService) HashPassword(password string) (string, error) {
	return password, nil
}

func (j *JWTService) ValidatePassword(hashedPassword, password string) bool {
	return hashedPassword == password
}
