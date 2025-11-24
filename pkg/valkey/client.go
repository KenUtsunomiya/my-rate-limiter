package valkey

import (
	"cmp"
	"context"
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

func (vc *Client) Ping(ctx context.Context) error {
	resp := vc.Client.Do(ctx, vc.Client.B().Ping().Build())
	return resp.Error()
}

func (vc *Client) HGetAll(ctx context.Context, key string) (map[string]string, error) {
	resp := vc.Client.Do(ctx, vc.Client.B().Hgetall().Key(key).Build())
	if resp.Error() != nil {
		return nil, resp.Error()
	}

	v, err := resp.AsStrMap()
	if err != nil {
		return nil, err
	}

	return v, nil
}

func (vc *Client) HSet(ctx context.Context, key string, values map[string]string) error {
	cmd := vc.Client.B().Hset().Key(key).FieldValue()
	for k, v := range values {
		cmd = cmd.FieldValue(k, v)
	}
	resp := vc.Client.Do(ctx, cmd.Build())
	return resp.Error()
}
