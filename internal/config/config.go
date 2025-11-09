package config

import (
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
)

// Config contains all application configurations
type Config struct {
	Server    ServerConfig
	RateLimit RateLimitConfig
	CORS      CORSConfig
	Logging   LoggingConfig
}

// ServerConfig HTTP server configurations
type ServerConfig struct {
	Host         string
	Port         string
	Prefork      bool
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
	AppName      string
}

// RateLimitConfig rate limiting configurations
type RateLimitConfig struct {
	Enabled bool
	Limit   int
	Window  time.Duration
}

// CORSConfig CORS configurations
type CORSConfig struct {
	AllowOrigins string
	AllowMethods string
	AllowHeaders string
	MaxAge       string
}

// LoggingConfig logging configurations
type LoggingConfig struct {
	Level string
	Env   string
}

// LoadConfig loads configurations from environment variables
func LoadConfig() (*Config, error) {
	// Load .env if it exists
	_ = godotenv.Load()

	cfg := &Config{
		Server: ServerConfig{
			Host:         getEnv("HOST", "0.0.0.0"),
			Port:         getEnv("PORT", "8080"),
			Prefork:      getEnvBool("PREFORK", false),
			ReadTimeout:  getEnvDuration("READ_TIMEOUT", 5*time.Second),
			WriteTimeout: getEnvDuration("WRITE_TIMEOUT", 10*time.Second),
			IdleTimeout:  getEnvDuration("IDLE_TIMEOUT", 120*time.Second),
			AppName:      getEnv("APP_NAME", "Fake Mill API"),
		},
		RateLimit: RateLimitConfig{
			Enabled: getEnvBool("RATE_LIMIT_ENABLED", true),
			Limit:   getEnvInt("RATE_LIMIT_REQUESTS", 60),
			Window:  getEnvDuration("RATE_LIMIT_WINDOW", time.Minute),
		},
		CORS: CORSConfig{
			AllowOrigins: getEnv("CORS_ALLOW_ORIGINS", "*"),
			AllowMethods: getEnv("CORS_ALLOW_METHODS", "GET, OPTIONS, POST"),
			AllowHeaders: getEnv("CORS_ALLOW_HEADERS", "Content-Type, Authorization"),
			MaxAge:       getEnv("CORS_MAX_AGE", "86400"),
		},
		Logging: LoggingConfig{
			Level: getEnv("LOG_LEVEL", "info"),
			Env:   getEnv("ENV", "production"),
		},
	}

	log.Debug().
		Str("host", cfg.Server.Host).
		Str("port", cfg.Server.Port).
		Int("rate_limit", cfg.RateLimit.Limit).
		Str("log_level", cfg.Logging.Level).
		Msg("Configuration loaded")

	return cfg, nil
}

// getEnv returns environment variable or default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvInt returns environment variable as int or default value
func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

// getEnvBool returns environment variable as bool or default value
func getEnvBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return defaultValue
}

// getEnvDuration returns environment variable as duration or default value
func getEnvDuration(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}

// Address returns the full server address
func (c *Config) Address() string {
	return c.Server.Host + ":" + c.Server.Port
}
