package middleware

import (
	"strings"

	"logistics-api/internal/adapters/primary/http/dto"
	"logistics-api/internal/core/domain"
	"logistics-api/internal/core/ports/services"
	"logistics-api/internal/pkg/logger"

	"github.com/gin-gonic/gin"
)

type AuthMiddleware struct {
	authService services.AuthService
	logger      logger.Logger
}

func NewAuthMiddleware(authService services.AuthService, logger logger.Logger) *AuthMiddleware {
	return &AuthMiddleware{
		authService: authService,
		logger:      logger,
	}
}

func (m *AuthMiddleware) RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			m.logger.Warn("Missing authorization header")
			dto.UnauthorizedResponse(c)
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			m.logger.Warn("Invalid authorization header format")
			dto.UnauthorizedResponse(c)
			c.Abort()
			return
		}

		token := parts[1]
		claims, err := m.authService.ValidateToken(token)
		if err != nil {
			m.logger.Warn("Invalid token", logger.Error(err))
			dto.UnauthorizedResponse(c)
			c.Abort()
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("user_email", claims.Email)
		c.Set("user_role", claims.Role)

		c.Next()
	}
}

func (m *AuthMiddleware) RequireAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, exists := c.Get("user_role")
		if !exists {
			m.logger.Error("User role not found in context")
			dto.ForbiddenResponse(c)
			c.Abort()
			return
		}

		role, ok := userRole.(domain.UserRole)
		if !ok || role != domain.AdminRole {
			m.logger.Warn("Access denied - admin role required",
				logger.String("user_role", string(role)))
			dto.ForbiddenResponse(c)
			c.Abort()
			return
		}

		c.Next()
	}
}

func (m *AuthMiddleware) RequireClient() gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, exists := c.Get("user_role")
		if !exists {
			m.logger.Error("User role not found in context")
			dto.ForbiddenResponse(c)
			c.Abort()
			return
		}

		role, ok := userRole.(domain.UserRole)
		if !ok || role != domain.ClientRole {
			m.logger.Warn("Access denied - client role required",
				logger.String("user_role", string(role)))
			dto.ForbiddenResponse(c)
			c.Abort()
			return
		}

		c.Next()
	}
}
