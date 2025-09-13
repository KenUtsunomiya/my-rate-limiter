package ratelimit

import (
	"context"

	"github.com/KenUtsunomiya/my-rate-limiter/pkg/valkey"
)

type VkRateLimiter struct {
	vkClient valkey.Client
}

func NewRateLimiter(vkClient valkey.Client) *VkRateLimiter {
	return &VkRateLimiter{
		vkClient: vkClient,
	}
}

func (rl *VkRateLimiter) Allow(ctx context.Context, userID string, method string, resource string) (bool, error) {
	return true, nil
}
