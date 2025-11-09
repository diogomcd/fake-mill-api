package middleware

import (
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

// mockClock implements Clock interface for testing
type mockClock struct {
	currentTime time.Time
	mu          sync.Mutex
}

func (m *mockClock) Now() time.Time {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.currentTime
}

func (m *mockClock) Advance(duration time.Duration) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.currentTime = m.currentTime.Add(duration)
}

func (m *mockClock) Set(t time.Time) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.currentTime = t
}

func TestRateLimiter_WithinLimit(t *testing.T) {
	clock := &mockClock{currentTime: time.Now()}
	rl := NewRateLimiterWithClock(60, time.Minute, clock)

	app := fiber.New()
	app.Use(rl.Middleware())
	app.Get("/test", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})

	for i := 0; i < 60; i++ {
		req := httptest.NewRequest("GET", "/test", nil)
		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, 200, resp.StatusCode, "Request %d should be allowed", i+1)
	}
}

func TestRateLimiter_ExceedsLimit(t *testing.T) {
	clock := &mockClock{currentTime: time.Now()}
	rl := NewRateLimiterWithClock(60, time.Minute, clock)

	app := fiber.New()
	app.Use(rl.Middleware())
	app.Get("/test", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})

	for i := 0; i < 60; i++ {
		req := httptest.NewRequest("GET", "/test", nil)
		resp, _ := app.Test(req)
		assert.Equal(t, 200, resp.StatusCode)
	}

	req := httptest.NewRequest("GET", "/test", nil)
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, 429, resp.StatusCode, "Request 61 should be rate limited")
}

func TestRateLimiter_Cleanup(t *testing.T) {
	clock := &mockClock{currentTime: time.Now()}
	rl := NewRateLimiterWithClock(60, time.Minute, clock)

	app := fiber.New()
	app.Use(rl.Middleware())
	app.Get("/test", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})

	for i := 0; i < 10; i++ {
		req := httptest.NewRequest("GET", "/test", nil)
		req.Header.Set("X-Forwarded-For", "192.168.1.1")
		app.Test(req)
	}

	clock.Advance(2 * time.Hour)
	rl.cleanup()

	rl.mu.RLock()
	_, exists := rl.requests["192.168.1.1"]
	rl.mu.RUnlock()

	assert.False(t, exists, "Old IP should be cleaned up")
}

func TestRateLimiter_MultipleIPs(t *testing.T) {
	clock := &mockClock{currentTime: time.Now()}
	rl := NewRateLimiterWithClock(60, time.Minute, clock)

	app := fiber.New()
	app.Use(rl.Middleware())
	app.Get("/test", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})

	var wg sync.WaitGroup
	ips := []string{"192.168.1.1", "192.168.1.2", "192.168.1.3", "192.168.1.4", "192.168.1.5"}

	for _, ip := range ips {
		wg.Add(1)
		go func(ipAddr string) {
			defer wg.Done()
			for i := 0; i < 10; i++ {
				req := httptest.NewRequest("GET", "/test", nil)
				req.Header.Set("X-Forwarded-For", ipAddr)
				resp, _ := app.Test(req)
				assert.Equal(t, 200, resp.StatusCode)
			}
		}(ip)
	}

	wg.Wait()
}

func TestRateLimiter_ResetAfterWindow(t *testing.T) {
	clock := &mockClock{currentTime: time.Now()}
	rl := NewRateLimiterWithClock(60, time.Minute, clock)

	app := fiber.New()
	app.Use(rl.Middleware())
	app.Get("/test", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})

	for i := 0; i < 60; i++ {
		req := httptest.NewRequest("GET", "/test", nil)
		app.Test(req)
	}

	req := httptest.NewRequest("GET", "/test", nil)
	resp, _ := app.Test(req)
	assert.Equal(t, 429, resp.StatusCode, "Should be rate limited")

	clock.Advance(2 * time.Minute)

	req = httptest.NewRequest("GET", "/test", nil)
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode, "Should allow requests after window expires")
}

