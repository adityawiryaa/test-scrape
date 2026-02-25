package controller

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/adityawiryaa/api/domain/entity"
	"github.com/adityawiryaa/api/pkg/httpclient"
	"github.com/adityawiryaa/api/pkg/response"
)

type Client struct {
	httpClient *httpclient.Client
	baseURL    string
	apiKey     string
}

func NewClient(baseURL string, apiKey string, timeout time.Duration) *Client {
	return &Client{
		httpClient: httpclient.New(timeout),
		baseURL:    baseURL,
		apiKey:     apiKey,
	}
}

func (c *Client) Register(ctx context.Context, req *entity.RegistrationRequest) (*entity.RegistrationResponse, error) {
	headers := map[string]string{
		"X-API-Key": c.apiKey,
	}

	resp, err := c.httpClient.Post(ctx, c.baseURL+"/register", req, headers)
	if err != nil {
		return nil, fmt.Errorf("registering with controller: %w", err)
	}

	var apiResp response.APIResponse
	if err := httpclient.DecodeResponse(resp, &apiResp); err != nil {
		return nil, err
	}

	if !apiResp.Success {
		return nil, fmt.Errorf("registration failed: %s", apiResp.Error.Message)
	}

	data, ok := apiResp.Data.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("unexpected response format")
	}

	return &entity.RegistrationResponse{
		AgentID: fmt.Sprintf("%v", data["agent_id"]),
		Status:  fmt.Sprintf("%v", data["status"]),
	}, nil
}

func (c *Client) FetchConfig(ctx context.Context, currentVersion int64) (*entity.Config, bool, error) {
	headers := map[string]string{
		"If-None-Match": strconv.FormatInt(currentVersion, 10),
		"X-API-Key":     c.apiKey,
	}

	resp, err := c.httpClient.Get(ctx, c.baseURL+"/config", headers)
	if err != nil {
		return nil, false, fmt.Errorf("fetching config: %w", err)
	}

	if resp.StatusCode == http.StatusNotModified {
		resp.Body.Close()
		return nil, false, nil
	}

	var apiResp response.APIResponse
	if err := httpclient.DecodeResponse(resp, &apiResp); err != nil {
		return nil, false, err
	}

	if !apiResp.Success {
		return nil, false, fmt.Errorf("fetch config failed: %s", apiResp.Error.Message)
	}

	data, ok := apiResp.Data.(map[string]any)
	if !ok {
		return nil, false, fmt.Errorf("unexpected response format")
	}

	version, _ := data["version"].(float64)
	pollInterval, _ := data["poll_interval_seconds"].(float64)

	configData := make(map[string]string)
	if rawData, ok := data["data"].(map[string]any); ok {
		for k, v := range rawData {
			configData[k] = fmt.Sprintf("%v", v)
		}
	}

	cfg := &entity.Config{
		ID:                  fmt.Sprintf("%v", data["id"]),
		Version:             int64(version),
		Data:                configData,
		PollIntervalSeconds: int(pollInterval),
	}

	return cfg, true, nil
}
