package ratelimit

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

type RateLimiter struct {
	client *redis.Client
	ctx    context.Context
}

// New creates a new rate limiter with Redis connection
func New(redisURL string) (*RateLimiter, error) {
	opt, err := redis.ParseURL(redisURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse Redis URL: %w", err)
	}

	client := redis.NewClient(opt)
	ctx := context.Background()

	_, err = client.Ping(ctx).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	return &RateLimiter{
		client: client,
		ctx:    ctx,
	}, nil
}

// IsAllowed checks if the address can make a request (not rate limited)
func (r *RateLimiter) IsAllowed(address string) (bool, time.Duration, error) {
	key := fmt.Sprintf("faucet_request:%s", address)
	
	lastRequestTime, err := r.client.Get(r.ctx, key).Result()
	if err == redis.Nil {
		return true, 0, nil
	}
	if err != nil {
		return false, 0, fmt.Errorf("failed to check rate limit: %w", err)
	}

	timestamp, err := strconv.ParseInt(lastRequestTime, 10, 64)
	if err != nil {
		return true, 0, nil
	}

	lastRequest := time.Unix(timestamp, 0)
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
	key := fmt.Sprintf("faucet_request:%s", address)
	now := time.Now().Unix()
	
	err := r.client.Set(r.ctx, key, now, 24*time.Hour).Err()
	if err != nil {
		return fmt.Errorf("failed to record request: %w", err)
	}

	return nil
}

// Close the Redis connection
func (r *RateLimiter) Close() error {
	return r.client.Close()
}