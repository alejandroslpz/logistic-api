package config

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	JWT      JWTConfig
	Logger   LoggerConfig
}

type ServerConfig struct {
	Host         string
	Port         string
	ReadTimeout  int
	WriteTimeout int
}

type DatabaseConfig struct {
	DatabaseURL  string
	MaxOpenConns int
	MaxIdleConns int
}

type JWTConfig struct {
	Secret     string
	ExpiryHour int
}

type LoggerConfig struct {
	Level  string
	Format string
}

func Load() (*Config, error) {
	config := &Config{
		Server:   loadServerConfig(),
		Database: loadDatabaseConfig(),
		JWT:      loadJWTConfig(),
		Logger:   loadLoggerConfig(),
	}

	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("configuration validation failed: %w", err)
	}

	return config, nil
}

func loadServerConfig() ServerConfig {
	return ServerConfig{
		Host:         getEnv("SERVER_HOST", "0.0.0.0"),
		Port:         getEnv("SERVER_PORT", getEnv("PORT", "8080")),
		ReadTimeout:  getEnvInt("SERVER_READ_TIMEOUT", 30),
		WriteTimeout: getEnvInt("SERVER_WRITE_TIMEOUT", 30),
	}
}

func loadDatabaseConfig() DatabaseConfig {
	return DatabaseConfig{
		DatabaseURL:  getEnv("DATABASE_URL", ""),
		MaxOpenConns: getEnvInt("DB_MAX_OPEN_CONNS", 10),
		MaxIdleConns: getEnvInt("DB_MAX_IDLE_CONNS", 2),
	}
}

func loadJWTConfig() JWTConfig {
	return JWTConfig{
		Secret:     getEnv("JWT_SECRET", ""),
		ExpiryHour: getEnvInt("JWT_EXPIRY_HOUR", 24),
	}
}

func loadLoggerConfig() LoggerConfig {
	return LoggerConfig{
		Level:  getEnv("LOG_LEVEL", "info"),
		Format: getEnv("LOG_FORMAT", "json"),
	}
}

func (c *Config) Validate() error {
	if c.Database.DatabaseURL == "" {
		return fmt.Errorf("DATABASE_URL is required")
	}

	if c.JWT.Secret == "" {
		return fmt.Errorf("JWT_SECRET is required")
	}

	if len(c.JWT.Secret) < 32 {
		return fmt.Errorf("JWT_SECRET must be at least 32 characters long")
	}

	if c.Logger.Level != "debug" && c.Logger.Level != "info" && c.Logger.Level != "warn" && c.Logger.Level != "error" {
		return fmt.Errorf("LOG_LEVEL must be one of: debug, info, warn, error")
	}

	if c.Logger.Format != "json" && c.Logger.Format != "text" {
		return fmt.Errorf("LOG_FORMAT must be either 'json' or 'text'")
	}

	return nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func (d *DatabaseConfig) GetDSN() string {
	return d.DatabaseURL
}
