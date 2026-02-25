package worker

import (
	"context"
	"fmt"
	"time"

	"github.com/adityawiryaa/api/domain/entity"
	"github.com/adityawiryaa/api/pkg/httpclient"
	"github.com/adityawiryaa/api/pkg/response"
)

type Client struct {
	httpClient *httpclient.Client
	baseURL    string
}

func NewClient(baseURL string, timeout time.Duration) *Client {
	return &Client{
		httpClient: httpclient.New(timeout),
		baseURL:    baseURL,
	}
}

func (c *Client) PushConfig(ctx context.Context, cfg *entity.Config) error {
	resp, err := c.httpClient.Post(ctx, c.baseURL+"/config", cfg, nil)
	if err != nil {
		return fmt.Errorf("pushing config to worker: %w", err)
	}

	var apiResp response.APIResponse
	if err := httpclient.DecodeResponse(resp, &apiResp); err != nil {
		return err
	}

	if !apiResp.Success {
		return fmt.Errorf("push config failed: %s", apiResp.Error.Message)
	}

	return nil
}
