package app

import (
	"logistics-api/internal/adapters/primary/health"
	"logistics-api/internal/adapters/primary/http"
	"logistics-api/internal/adapters/primary/http/handlers"
	"logistics-api/internal/adapters/primary/http/middleware"
	authService "logistics-api/internal/adapters/secondary/auth"
	"logistics-api/internal/adapters/secondary/database/postgres"
	"logistics-api/internal/adapters/secondary/external"
	loggerAdapter "logistics-api/internal/adapters/secondary/logger"
	"logistics-api/internal/config"
	authUseCase "logistics-api/internal/core/usecases/auth"
	"logistics-api/internal/core/usecases/order"
	"logistics-api/internal/pkg/logger"
	"logistics-api/internal/pkg/validator"

	"gorm.io/gorm"
)

type Container struct {
	// Config
	Config *config.Config

	// Infrastructure
	DB     *gorm.DB
	Logger logger.Logger

	// Services
	AuthService       *authService.JWTService
	CoordinateService *external.CoordinateService

	// Repositories
	UserRepository  *postgres.UserRepository
	OrderRepository *postgres.OrderRepository

	// Use Cases
	RegisterUC     *authUseCase.RegisterUseCase
	LoginUC        *authUseCase.LoginUseCase
	CreateOrderUC  *order.CreateOrderUseCase
	GetOrdersUC    *order.GetOrdersUseCase
	UpdateStatusUC *order.UpdateOrderStatusUseCase

	// HTTP Layer
	Validator      *validator.Validator
	AuthMiddleware *middleware.AuthMiddleware
	AuthHandler    *handlers.AuthHandler
	OrderHandler   *handlers.OrderHandler
	HealthHandler  *health.HealthHandler
	Router         *http.Router
	Server         *http.Server
}

func NewContainer() (*Container, error) {
	container := &Container{}

	if err := container.initConfig(); err != nil {
		return nil, err
	}

	if err := container.initInfrastructure(); err != nil {
		return nil, err
	}

	if err := container.initServices(); err != nil {
		return nil, err
	}

	if err := container.initRepositories(); err != nil {
		return nil, err
	}

	if err := container.initUseCases(); err != nil {
		return nil, err
	}

	if err := container.initHTTPLayer(); err != nil {
		return nil, err
	}

	return container, nil
}

func (c *Container) initConfig() error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}
	c.Config = cfg
	return nil
}

func (c *Container) initInfrastructure() error {
	// Initialize logger
	c.Logger = loggerAdapter.NewLogrusAdapter(c.Config.Logger.Level, c.Config.Logger.Format)

	// Initialize database
	db, err := postgres.NewConnection(&c.Config.Database)
	if err != nil {
		return err
	}
	c.DB = db

	c.Logger.Info("Infrastructure initialized successfully")
	return nil
}

func (c *Container) initServices() error {
	// Auth service
	c.AuthService = authService.NewJWTService(c.Config.JWT.Secret, c.Config.JWT.ExpiryHour)

	// Coordinate service
	c.CoordinateService = external.NewCoordinateService()

	c.Logger.Info("Services initialized successfully")
	return nil
}

func (c *Container) initRepositories() error {
	// User repository
	c.UserRepository = postgres.NewUserRepository(c.DB)

	// Order repository
	c.OrderRepository = postgres.NewOrderRepository(c.DB)

	c.Logger.Info("Repositories initialized successfully")
	return nil
}

func (c *Container) initUseCases() error {
	// Auth use cases
	c.RegisterUC = authUseCase.NewRegisterUseCase(c.UserRepository, c.AuthService, c.Logger)
	c.LoginUC = authUseCase.NewLoginUseCase(c.UserRepository, c.AuthService, c.Logger)

	// Order use cases
	c.CreateOrderUC = order.NewCreateOrderUseCase(c.OrderRepository, c.UserRepository, c.CoordinateService, c.Logger)
	c.GetOrdersUC = order.NewGetOrdersUseCase(c.OrderRepository, c.UserRepository, c.Logger)
	c.UpdateStatusUC = order.NewUpdateOrderStatusUseCase(c.OrderRepository, c.Logger)

	c.Logger.Info("Use cases initialized successfully")
	return nil
}

func (c *Container) initHTTPLayer() error {
	// Validator
	c.Validator = validator.New()

	// Middleware
	c.AuthMiddleware = middleware.NewAuthMiddleware(c.AuthService, c.Logger)

	// Handlers
	c.AuthHandler = handlers.NewAuthHandler(c.RegisterUC, c.LoginUC, c.Validator, c.Logger)
	c.OrderHandler = handlers.NewOrderHandler(c.CreateOrderUC, c.GetOrdersUC, c.UpdateStatusUC, c.Validator, c.Logger)
	c.HealthHandler = health.NewHealthHandler(c.DB, c.Logger)

	// Router
	routerConfig := http.RouterConfig{
		AuthHandler:    c.AuthHandler,
		OrderHandler:   c.OrderHandler,
		HealthHandler:  c.HealthHandler,
		AuthMiddleware: c.AuthMiddleware,
		Logger:         c.Logger,
		RateLimitRPS:   10.0, // 10 requests per second
		RateLimitBurst: 20,   // Burst of 20 requests
	}
	c.Router = http.NewRouter(routerConfig)
	c.Router.SetupRoutes()

	// Server
	c.Server = http.NewServer(c.Router, &c.Config.Server, c.Logger)

	c.Logger.Info("HTTP layer initialized successfully")
	return nil
}

func (c *Container) Start() error {
	c.Logger.Info("Starting application...")
	return c.Server.Start()
}

func (c *Container) Stop() error {
	c.Logger.Info("Stopping application...")
	return c.Server.Stop()
}
