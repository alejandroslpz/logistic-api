package http

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"logistics-api/internal/config"
	"logistics-api/internal/pkg/logger"
)

type Server struct {
	httpServer *http.Server
	logger     logger.Logger
}

func NewServer(router *Router, config *config.ServerConfig, logger logger.Logger) *Server {
	address := fmt.Sprintf("%s:%s", config.Host, config.Port)

	httpServer := &http.Server{
		Addr:         address,
		Handler:      router.GetEngine(),
		ReadTimeout:  time.Duration(config.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(config.WriteTimeout) * time.Second,
	}

	return &Server{
		httpServer: httpServer,
		logger:     logger,
	}
}

func (s *Server) Start() error {
	// Start server in a goroutine
	go func() {
		s.logger.Info("Starting HTTP server",
			logger.String("address", s.httpServer.Addr))

		if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.logger.Fatal("Failed to start server", logger.Error(err))
		}
	}()

	s.logger.Info("Server started successfully",
		logger.String("address", s.httpServer.Addr))

	// Wait for interrupt signal to gracefully shutdown
	return s.waitForShutdown()
}

func (s *Server) waitForShutdown() error {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit
	s.logger.Info("Shutdown signal received")

	// Create context with timeout for graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	s.logger.Info("Shutting down server...")

	if err := s.httpServer.Shutdown(ctx); err != nil {
		s.logger.Error("Server forced to shutdown", logger.Error(err))
		return err
	}

	s.logger.Info("Server shutdown completed")
	return nil
}

func (s *Server) Stop() error {
	return s.httpServer.Close()
}
