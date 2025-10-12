package ratelimit

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/KenUtsunomiya/my-rate-limiter/pkg/valkey"
)

type VkRateLimiter struct {
	vkClient   valkey.Client
	maxTokens  float64
	refillRate float64
	keyPrefix  string
}

func NewRateLimiter(vkClient valkey.Client) *VkRateLimiter {
	return &VkRateLimiter{
		vkClient:   vkClient,
		maxTokens:  10,
		refillRate: 1,
		keyPrefix:  "ratelimit",
	}
}

// Allow checks if a request is allowed under the rate limit policy.
func (rl *VkRateLimiter) Allow(ctx context.Context, userID string, method string, resource string) (bool, error) {
	key := fmt.Sprintf("%s:%s:%s:%s", rl.keyPrefix, userID, method, resource)
	now := time.Now().UnixMilli()

	data, err := rl.vkClient.HGetAll(ctx, key)
	if err != nil {
		return false, err
	}

	var tokens float64
	var lastRefill int64
	if len(data) == 0 {
		tokens = rl.maxTokens
		lastRefill = now
	} else {
		tokens, _ = strconv.ParseFloat(data["tokens"], 64)
		lastRefill, _ = strconv.ParseInt(data["lastRefillTs"], 10, 64)
	}

	elapsed := float64(now - lastRefill)
	if elapsed > 0 {
		tokens = min(rl.maxTokens, tokens+(elapsed/1000.0)*rl.refillRate)
		lastRefill = now
	}

	allowed := false
	if tokens >= 1 {
		tokens -= 1
		allowed = true
	}

	err = rl.vkClient.HSet(ctx, key, map[string]string{
		"tokens":       fmt.Sprintf("%f", tokens),
		"lastRefillTs": fmt.Sprintf("%d", lastRefill),
	})
	if err != nil {
		return false, err
	}

	_ = rl.vkClient.Expire(ctx, key, int64(2*time.Minute))

	return allowed, nil
}
