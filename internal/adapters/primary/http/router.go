package http

import (
	"logistics-api/internal/adapters/primary/health"
	"logistics-api/internal/adapters/primary/http/handlers"
	"logistics-api/internal/adapters/primary/http/middleware"
	"logistics-api/internal/pkg/logger"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

type Router struct {
	engine         *gin.Engine
	authHandler    *handlers.AuthHandler
	orderHandler   *handlers.OrderHandler
	healthHandler  *health.HealthHandler
	authMiddleware *middleware.AuthMiddleware
	logger         logger.Logger
}

type RouterConfig struct {
	AuthHandler    *handlers.AuthHandler
	OrderHandler   *handlers.OrderHandler
	HealthHandler  *health.HealthHandler
	AuthMiddleware *middleware.AuthMiddleware
	Logger         logger.Logger
	RateLimitRPS   float64
	RateLimitBurst int
}

func NewRouter(config RouterConfig) *Router {
	gin.SetMode(gin.ReleaseMode)

	engine := gin.New()

	engine.Use(gin.Recovery())

	engine.Use(middleware.CORS())

	engine.Use(middleware.LoggingMiddleware(config.Logger))

	rateLimiter := middleware.NewRateLimiter(
		rate.Limit(config.RateLimitRPS),
		config.RateLimitBurst,
		config.Logger,
	)
	engine.Use(rateLimiter.RateLimit())

	rateLimiter.StartCleanup(time.Minute * 5)

	return &Router{
		engine:         engine,
		authHandler:    config.AuthHandler,
		orderHandler:   config.OrderHandler,
		healthHandler:  config.HealthHandler,
		authMiddleware: config.AuthMiddleware,
		logger:         config.Logger,
	}
}

func (r *Router) SetupRoutes() {
	health := r.engine.Group("/health")
	{
		health.GET("/", r.healthHandler.HealthCheck)
		health.GET("/live", r.healthHandler.LivenessCheck)
		health.GET("/ready", r.healthHandler.ReadinessCheck)
	}

	api := r.engine.Group("/api/v1")

	auth := api.Group("/auth")
	{
		auth.POST("/register", r.authHandler.Register)
		auth.POST("/login", r.authHandler.Login)
	}

	protected := api.Group("/")
	protected.Use(r.authMiddleware.RequireAuth())
	{
		orders := protected.Group("/orders")
		{
			orders.POST("/", r.orderHandler.CreateOrder)
			orders.GET("/", r.orderHandler.GetOrders)
			orders.GET("/:id", r.orderHandler.GetOrderByID)
			orders.PUT("/:id/status", r.authMiddleware.RequireAdmin(), r.orderHandler.UpdateOrderStatus)
		}

		admin := protected.Group("/admin")
		admin.Use(r.authMiddleware.RequireAdmin())
		{
			// Future admin endpoints can go here
			// admin.GET("/users", r.userHandler.GetAllUsers)
			// admin.GET("/analytics", r.analyticsHandler.GetAnalytics)
		}
	}

	r.engine.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{
			"success": false,
			"error": map[string]interface{}{
				"code":    404,
				"type":    "not_found",
				"message": "Route not found",
			},
		})
	})

	r.logger.Info("Routes setup completed")
}

func (r *Router) GetEngine() *gin.Engine {
	return r.engine
}

func (r *Router) GetRoutes() []gin.RouteInfo {
	return r.engine.Routes()
}
