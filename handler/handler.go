package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/afrochainorg/afrochain-faucet/command"
	"github.com/afrochainorg/afrochain-faucet/config"
	"github.com/afrochainorg/afrochain-faucet/ratelimit"
)

type Handler struct {
	Config      config.Config
	RateLimiter *ratelimit.RateLimiter
}

// New creates a new Handler with rate limiting capabilities
func New(c config.Config) (*Handler, error) {
	rateLimiter, err := ratelimit.New(c.RedisURL)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize rate limiter: %w", err)
	}

	return &Handler{
		Config:      c,
		RateLimiter: rateLimiter,
	}, nil
}

type RequestPayload struct {
	Recipient string `json:"recipient"`
	Amount    string `json:"amount"`
}

type RequestResponse struct {
	TxHash string `json:"txhash"`
}

func (h *Handler) Request(c *gin.Context) {
	var f RequestPayload
	if err := c.BindJSON(&f); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	allowed, timeRemaining, err := h.RateLimiter.IsAllowed(f.Recipient)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to check rate limit: " + err.Error(),
		})
		return
	}

	if !allowed {
		hours := int(timeRemaining.Hours())
		minutes := int(timeRemaining.Minutes()) % 60

		nextRequestTime := time.Now().Add(timeRemaining)

		c.JSON(http.StatusTooManyRequests, gin.H{
			"error":           "Rate limit exceeded",
			"message":         "You can only request tokens once per 24 hours",
			"timeRemaining":   fmt.Sprintf("%dh %dm", hours, minutes),
			"nextRequestTime": nextRequestTime.Format(time.RFC3339),
		})
		return
	}

	output, err := command.ExecuteTransfer(h.Config, f.Recipient, f.Amount)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	var resp RequestResponse
	if err := json.Unmarshal(output, &resp); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	// Record successful request for rate limiting
	if err := h.RateLimiter.RecordRequest(f.Recipient); err != nil {
		// Non-fatal: transfer succeeded, just log the rate limit recording failure
		fmt.Printf("Warning: Failed to record rate limit: %v\n", err)
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "success",
		"txHash":  resp.TxHash,
	})
}
