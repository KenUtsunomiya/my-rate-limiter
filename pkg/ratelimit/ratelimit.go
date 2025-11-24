package ratelimit

import "context"

type RateLimiter struct {
}

func NewRateLimiter() *RateLimiter {
	return &RateLimiter{}
}

func (rl *RateLimiter) Allow(ctx context.Context, userId string, method string, resource string) (bool, error) {
	return true, nil
}
