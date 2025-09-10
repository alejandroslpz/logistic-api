package health

import (
	"context"
	"net/http"
	"time"

	"logistics-api/internal/adapters/primary/http/dto"
	"logistics-api/internal/pkg/logger"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type HealthHandler struct {
	db     *gorm.DB
	logger logger.Logger
}

type HealthResponse struct {
	Status    string           `json:"status"`
	Timestamp string           `json:"timestamp"`
	Version   string           `json:"version"`
	Checks    map[string]Check `json:"checks"`
}

type Check struct {
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`
}

func NewHealthHandler(db *gorm.DB, logger logger.Logger) *HealthHandler {
	return &HealthHandler{
		db:     db,
		logger: logger,
	}
}

func (h *HealthHandler) HealthCheck(c *gin.Context) {
	response := HealthResponse{
		Status:    "healthy",
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Version:   "1.0.0",
		Checks:    make(map[string]Check),
	}

	// Check database connection
	dbCheck := h.checkDatabase()
	response.Checks["database"] = dbCheck

	// If any check fails, mark overall status as unhealthy
	if dbCheck.Status != "healthy" {
		response.Status = "unhealthy"
	}

	statusCode := http.StatusOK
	if response.Status == "unhealthy" {
		statusCode = http.StatusServiceUnavailable
	}

	c.JSON(statusCode, response)
}

func (h *HealthHandler) ReadinessCheck(c *gin.Context) {
	// Simple readiness check
	dbCheck := h.checkDatabase()

	if dbCheck.Status == "healthy" {
		dto.SuccessResponse(c, http.StatusOK, "Service is ready", map[string]string{
			"status":    "ready",
			"timestamp": time.Now().UTC().Format(time.RFC3339),
		})
	} else {
		dto.ErrorResponse(c, http.StatusServiceUnavailable, "service_unavailable", "Service is not ready")
	}
}

func (h *HealthHandler) LivenessCheck(c *gin.Context) {
	// Simple liveness check - just return OK
	dto.SuccessResponse(c, http.StatusOK, "Service is alive", map[string]string{
		"status":    "alive",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

func (h *HealthHandler) checkDatabase() Check {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	sqlDB, err := h.db.DB()
	if err != nil {
		h.logger.Error("Failed to get database instance", logger.Error(err))
		return Check{
			Status:  "unhealthy",
			Message: "Failed to get database instance",
		}
	}

	if err := sqlDB.PingContext(ctx); err != nil {
		h.logger.Error("Database ping failed", logger.Error(err))
		return Check{
			Status:  "unhealthy",
			Message: "Database connection failed",
		}
	}

	return Check{
		Status: "healthy",
	}
}
