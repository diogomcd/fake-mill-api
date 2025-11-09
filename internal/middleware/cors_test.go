package middleware

import (
	"net/http/httptest"
	"testing"

	"github.com/diogomcd/fake-mill-api/internal/config"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func TestCORS_Headers(t *testing.T) {
	cfg := config.CORSConfig{
		AllowOrigins: "*",
		AllowMethods: "GET, OPTIONS, POST",
		AllowHeaders: "Content-Type, Authorization",
		MaxAge:       "86400",
	}

	app := fiber.New()
	app.Use(CORSWithConfig(cfg))
	app.Get("/test", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})

	req := httptest.NewRequest("GET", "/test", nil)
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, "*", resp.Header.Get("Access-Control-Allow-Origin"))
	assert.Equal(t, "GET, OPTIONS, POST", resp.Header.Get("Access-Control-Allow-Methods"))
	assert.Equal(t, "Content-Type, Authorization", resp.Header.Get("Access-Control-Allow-Headers"))
	assert.Equal(t, "86400", resp.Header.Get("Access-Control-Max-Age"))
}

func TestCORS_PreflightRequest(t *testing.T) {
	cfg := config.CORSConfig{
		AllowOrigins: "*",
		AllowMethods: "GET, OPTIONS, POST",
		AllowHeaders: "Content-Type, Authorization",
		MaxAge:       "86400",
	}

	app := fiber.New()
	app.Use(CORSWithConfig(cfg))
	app.Get("/test", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})

	req := httptest.NewRequest("OPTIONS", "/test", nil)
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, 204, resp.StatusCode, "Preflight request should return 204")
}
