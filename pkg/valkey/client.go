package valkey

import (
	"cmp"
	"context"
	"fmt"
	"os"

	vk "github.com/valkey-io/valkey-go"
)

type Client struct {
	vk.Client
}

func NewClient() (*Client, error) {
	addr := cmp.Or(os.Getenv("VALKEY_BASE_URL"), "localhost:6379")
	userName := os.Getenv("VALKEY_USER_NAME")
	password := os.Getenv("VALKEY_PASSWORD")

	vkClient, err := vk.NewClient(vk.ClientOption{
		InitAddress: []string{addr},
		Username:    userName,
		Password:    password,
	})
	if err != nil {
		return nil, err
	}
	return &Client{vkClient}, nil
}

func (vc *Client) Close() {
	vc.Client.Close()
}

func (vc *Client) Hello(ctx context.Context) error {
	resp := vc.Client.Do(ctx, vc.Client.B().Ping().Build())
	return resp.Error()
}

func (vc *Client) CheckRateLimit(ctx context.Context, userId string, method string, resource string) (bool, error) {
	key := fmt.Sprintf("%s:%s:%s", userId, method, resource)
	cmd := vc.Client.B().Incr().Key(key).Build()
	resp := vc.Client.Do(ctx, cmd)
	if resp.Error() != nil {
		return false, resp.Error()
	}
	v, err := resp.AsInt64()
	if err != nil {
		return false, err
	}
	return v <= 10, nil
}
