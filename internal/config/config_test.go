package config

import (
	"os"
	"testing"
	"time"
)

func TestLoadConfig(t *testing.T) {
	// Save original env vars
	originalEnv := map[string]string{
		"PORT":                os.Getenv("PORT"),
		"HOST":                os.Getenv("HOST"),
		"LOG_LEVEL":           os.Getenv("LOG_LEVEL"),
		"RATE_LIMIT_REQUESTS": os.Getenv("RATE_LIMIT_REQUESTS"),
		"RATE_LIMIT_ENABLED":  os.Getenv("RATE_LIMIT_ENABLED"),
	}

	// Restore env vars after test
	defer func() {
		for key, value := range originalEnv {
			if value == "" {
				os.Unsetenv(key)
			} else {
				os.Setenv(key, value)
			}
		}
	}()

	t.Run("default values", func(t *testing.T) {
		// Clear env vars
		os.Unsetenv("PORT")
		os.Unsetenv("HOST")
		os.Unsetenv("LOG_LEVEL")

		cfg, err := LoadConfig()
		if err != nil {
			t.Fatalf("LoadConfig() error = %v", err)
		}

		if cfg.Server.Port != "8080" {
			t.Errorf("Server.Port = %v, want 8080", cfg.Server.Port)
		}

		if cfg.Server.Host != "0.0.0.0" {
			t.Errorf("Server.Host = %v, want 0.0.0.0", cfg.Server.Host)
		}

		if cfg.Logging.Level != "info" {
			t.Errorf("Logging.Level = %v, want info", cfg.Logging.Level)
		}

		if cfg.RateLimit.Limit != 60 {
			t.Errorf("RateLimit.Limit = %v, want 60", cfg.RateLimit.Limit)
		}
	})

	t.Run("custom values", func(t *testing.T) {
		os.Setenv("PORT", "3000")
		os.Setenv("HOST", "127.0.0.1")
		os.Setenv("LOG_LEVEL", "debug")
		os.Setenv("RATE_LIMIT_REQUESTS", "100")

		cfg, err := LoadConfig()
		if err != nil {
			t.Fatalf("LoadConfig() error = %v", err)
		}

		if cfg.Server.Port != "3000" {
			t.Errorf("Server.Port = %v, want 3000", cfg.Server.Port)
		}

		if cfg.Server.Host != "127.0.0.1" {
			t.Errorf("Server.Host = %v, want 127.0.0.1", cfg.Server.Host)
		}

		if cfg.Logging.Level != "debug" {
			t.Errorf("Logging.Level = %v, want debug", cfg.Logging.Level)
		}

		if cfg.RateLimit.Limit != 100 {
			t.Errorf("RateLimit.Limit = %v, want 100", cfg.RateLimit.Limit)
		}
	})

	t.Run("rate limit disabled", func(t *testing.T) {
		os.Setenv("RATE_LIMIT_ENABLED", "false")

		cfg, err := LoadConfig()
		if err != nil {
			t.Fatalf("LoadConfig() error = %v", err)
		}

		if cfg.RateLimit.Enabled {
			t.Error("RateLimit.Enabled = true, want false")
		}
	})

	t.Run("address method", func(t *testing.T) {
		os.Setenv("HOST", "localhost")
		os.Setenv("PORT", "9000")

		cfg, err := LoadConfig()
		if err != nil {
			t.Fatalf("LoadConfig() error = %v", err)
		}

		expected := "localhost:9000"
		if addr := cfg.Address(); addr != expected {
			t.Errorf("Address() = %v, want %v", addr, expected)
		}
	})
}

func TestGetEnvHelpers(t *testing.T) {
	t.Run("getEnv", func(t *testing.T) {
		os.Setenv("TEST_STRING", "value")
		defer os.Unsetenv("TEST_STRING")

		if got := getEnv("TEST_STRING", "default"); got != "value" {
			t.Errorf("getEnv() = %v, want value", got)
		}

		if got := getEnv("NONEXISTENT", "default"); got != "default" {
			t.Errorf("getEnv() = %v, want default", got)
		}
	})

	t.Run("getEnvInt", func(t *testing.T) {
		os.Setenv("TEST_INT", "42")
		defer os.Unsetenv("TEST_INT")

		if got := getEnvInt("TEST_INT", 10); got != 42 {
			t.Errorf("getEnvInt() = %v, want 42", got)
		}

		if got := getEnvInt("NONEXISTENT", 10); got != 10 {
			t.Errorf("getEnvInt() = %v, want 10", got)
		}

		os.Setenv("TEST_INT", "invalid")
		if got := getEnvInt("TEST_INT", 10); got != 10 {
			t.Errorf("getEnvInt() with invalid = %v, want 10", got)
		}
	})

	t.Run("getEnvBool", func(t *testing.T) {
		os.Setenv("TEST_BOOL", "true")
		defer os.Unsetenv("TEST_BOOL")

		if got := getEnvBool("TEST_BOOL", false); !got {
			t.Error("getEnvBool() = false, want true")
		}

		if got := getEnvBool("NONEXISTENT", false); got {
			t.Error("getEnvBool() = true, want false")
		}

		os.Setenv("TEST_BOOL", "invalid")
		if got := getEnvBool("TEST_BOOL", false); got {
			t.Error("getEnvBool() with invalid = true, want false")
		}
	})

	t.Run("getEnvDuration", func(t *testing.T) {
		os.Setenv("TEST_DURATION", "5s")
		defer os.Unsetenv("TEST_DURATION")

		if got := getEnvDuration("TEST_DURATION", time.Second); got != 5*time.Second {
			t.Errorf("getEnvDuration() = %v, want 5s", got)
		}

		if got := getEnvDuration("NONEXISTENT", time.Second); got != time.Second {
			t.Errorf("getEnvDuration() = %v, want 1s", got)
		}

		os.Setenv("TEST_DURATION", "invalid")
		if got := getEnvDuration("TEST_DURATION", time.Second); got != time.Second {
			t.Errorf("getEnvDuration() with invalid = %v, want 1s", got)
		}
	})
}
