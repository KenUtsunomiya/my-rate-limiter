package server

import "context"

type RateLimiter interface {
	allow(key string) (bool, error)
}

type InMemoryRateLimiter struct{}

func newInMemoryRateLimiter() *InMemoryRateLimiter {
	return &InMemoryRateLimiter{}
}

func (rl *InMemoryRateLimiter) allow(ctx context.Context, userId string, method string, resource string) (bool, error) {
	return true, nil
}
