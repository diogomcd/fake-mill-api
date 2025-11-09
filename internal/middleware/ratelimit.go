package middleware

import (
	"net"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

// Clock interface for time operations (allows mocking in tests)
type Clock interface {
	Now() time.Time
}

// realClock implements Clock using real time
type realClock struct{}

func (realClock) Now() time.Time {
	return time.Now()
}

// RateLimiter implements rate limiting by IP
type RateLimiter struct {
	requests map[string][]time.Time
	mu       sync.RWMutex
	limit    int
	window   time.Duration
	clock    Clock
}

// NewRateLimiter creates a new instance of RateLimiter
func NewRateLimiter(limit int, window time.Duration) *RateLimiter {
	rl := &RateLimiter{
		requests: make(map[string][]time.Time),
		limit:    limit,
		window:   window,
		clock:    realClock{},
	}

	go func() {
		ticker := time.NewTicker(5 * time.Minute)
		defer ticker.Stop()

		for range ticker.C {
			rl.cleanup()
		}
	}()

	return rl
}

// NewRateLimiterWithClock creates a RateLimiter with a custom clock (for testing)
func NewRateLimiterWithClock(limit int, window time.Duration, clock Clock) *RateLimiter {
	rl := &RateLimiter{
		requests: make(map[string][]time.Time),
		limit:    limit,
		window:   window,
		clock:    clock,
	}

	return rl
}

// Middleware returns the Fiber middleware
func (rl *RateLimiter) Middleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		ip := getClientIP(c)

		if !rl.isAllowed(ip) {
			log.Warn().
				Str("ip", ip).
				Str("path", c.Path()).
				Str("method", c.Method()).
				Msg("Rate limit exceeded")

			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"error":   "rate limit exceeded",
				"message": "You have reached the request limit. Please try again later.",
			})
		}

		return c.Next()
	}
}

// isAllowed verifies if the request is allowed
func (rl *RateLimiter) isAllowed(ip string) bool {
	now := rl.clock.Now()
	windowStart := now.Add(-rl.window)

	rl.mu.RLock()
	requests := rl.requests[ip]

	validCount := 0
	for _, req := range requests {
		if req.After(windowStart) {
			validCount++
		}
	}
	rl.mu.RUnlock()

	if validCount >= rl.limit {
		return false
	}

	rl.mu.Lock()
	defer rl.mu.Unlock()

	// Re-verify after acquiring exclusive lock
	requests = rl.requests[ip]
	validRequests := []time.Time{}
	for _, req := range requests {
		if req.After(windowStart) {
			validRequests = append(validRequests, req)
		}
	}

	// Verify again after cleanup
	if len(validRequests) >= rl.limit {
		rl.requests[ip] = validRequests
		return false
	}

	// Add the new request
	validRequests = append(validRequests, now)
	rl.requests[ip] = validRequests

	return true
}

// cleanup remove old entries from the map
func (rl *RateLimiter) cleanup() {
	now := rl.clock.Now()
	windowStart := now.Add(-rl.window * 2)

	// First verify if there is something to clean with RLock
	rl.mu.RLock()
	hasWork := false
	for _, requests := range rl.requests {
		for _, req := range requests {
			if req.Before(windowStart) {
				hasWork = true
				break
			}
		}
		if hasWork {
			break
		}
	}
	rl.mu.RUnlock()

	// If there is no work, return immediately
	if !hasWork {
		return
	}

	// Now clean with exclusive Lock
	rl.mu.Lock()
	defer rl.mu.Unlock()

	removedIPs := 0
	for ip, requests := range rl.requests {
		validRequests := []time.Time{}
		for _, req := range requests {
			if req.After(windowStart) {
				validRequests = append(validRequests, req)
			}
		}

		if len(validRequests) == 0 {
			delete(rl.requests, ip)
			removedIPs++
		} else {
			rl.requests[ip] = validRequests
		}
	}

	if removedIPs > 0 {
		log.Debug().
			Int("removed_ips", removedIPs).
			Int("active_ips", len(rl.requests)).
			Msg("Rate limiter cleanup completed")
	}
}

// getClientIP gets the client IP
func getClientIP(c *fiber.Ctx) string {
	// Verify X-Forwarded-For (proxy)
	if ip := c.Get("X-Forwarded-For"); ip != "" {
		if idx := len(ip) - 1; idx >= 0 {
			for i := 0; i <= idx; i++ {
				if ip[i] == ',' {
					return ip[:i]
				}
			}
		}
		return ip
	}

	// Verify X-Real-IP
	if ip := c.Get("X-Real-IP"); ip != "" {
		return ip
	}

	// Get the IP from the Fiber context
	ip := c.IP()
	if ip != "" {
		return ip
	}

	// Fallback
	if addr := c.Context().RemoteAddr().String(); addr != "" {
		host, _, err := net.SplitHostPort(addr)
		if err == nil {
			return host
		}
		return addr
	}

	return "unknown"
}
