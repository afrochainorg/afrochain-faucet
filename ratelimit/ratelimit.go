package ratelimit

import (
	"sync"
	"time"
)

type RateLimiter struct {
	mu       sync.Mutex
	requests map[string]time.Time
}

// New creates a new in-memory rate limiter
func New() (*RateLimiter, error) {
	rl := &RateLimiter{
		requests: make(map[string]time.Time),
	}

	// Start a cleanup goroutine to remove expired entries
	go rl.cleanup()

	return rl, nil
}

func (r *RateLimiter) cleanup() {
	ticker := time.NewTicker(1 * time.Hour)
	for range ticker.C {
		r.mu.Lock()
		now := time.Now()
		for addr, timestamp := range r.requests {
			if now.Sub(timestamp) >= 24*time.Hour {
				delete(r.requests, addr)
			}
		}
		r.mu.Unlock()
	}
}

// IsAllowed checks if the address can make a request (not rate limited)
func (r *RateLimiter) IsAllowed(address string) (bool, time.Duration, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	lastRequest, exists := r.requests[address]
	if !exists {
		return true, 0, nil
	}

	timeElapsed := time.Since(lastRequest)
	rateLimitDuration := 24 * time.Hour

	if timeElapsed < rateLimitDuration {
		timeRemaining := rateLimitDuration - timeElapsed
		return false, timeRemaining, nil
	}

	return true, 0, nil
}

// RecordRequest records that the address made a request
func (r *RateLimiter) RecordRequest(address string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.requests[address] = time.Now()
	return nil
}

// Close the rate limiter
func (r *RateLimiter) Close() error {
	return nil
}
