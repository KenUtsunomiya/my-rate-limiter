package ratelimit

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/KenUtsunomiya/my-rate-limiter/pkg/valkey"
)

type RateLimiter struct {
	client     valkey.Client
	keyPrefix  string
	maxTokens  float64
	refillRate float64
}

func NewRateLimiter(client valkey.Client, keyPrefix string, maxTokens float64, refillRate float64) *RateLimiter {
	return &RateLimiter{
		client:     client,
		keyPrefix:  keyPrefix,
		maxTokens:  maxTokens,
		refillRate: refillRate,
	}
}

func (rl *RateLimiter) Allow(ctx context.Context, userId string, method string, resource string) (bool, error) {
	now := time.Now().UnixMilli()

	key := rl.keyPrefix + userId + ":" + method + ":" + resource
	data, err := rl.client.HGetAll(ctx, key)
	if err != nil {
		log.Printf("failed to get data from valkey: %v", err)
		return true, nil
	}

	var tokens float64
	var lastRefill int64
	if len(data) == 0 {
		tokens = rl.maxTokens
		lastRefill = now
	} else {
		tokens, err = strconv.ParseFloat(data["tokens"], 64)
		if err != nil {
			log.Printf("failed to parse tokens: %v", err)
			tokens = rl.maxTokens
		}
		lastRefill, err = strconv.ParseInt(data["lastRefill"], 10, 64)
		if err != nil {
			log.Printf("failed to parse lastRefill: %v", err)
			lastRefill = now
		}
	}

	elapsed := float64(now - lastRefill)
	if elapsed > 0 {
		tokens = tokens + (elapsed/1000.0)*rl.refillRate
		lastRefill = now
	}

	allowed := false
	if tokens >= 1 {
		tokens -= 1
		allowed = true
	}

	if err = rl.client.HSet(ctx, key, map[string]string{
		"tokens":     fmt.Sprintf("%f", tokens),
		"lastRefill": fmt.Sprintf("%d", lastRefill),
	}); err != nil {
		log.Printf("failed to set data to valkey: %v", err)
		return true, nil
	}

	return allowed, nil
}
