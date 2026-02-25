package usecases

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"time"
)

type DefaultHTTPExecutor struct {
	client *http.Client
}

func NewHTTPExecutor(timeout time.Duration) *DefaultHTTPExecutor {
	return &DefaultHTTPExecutor{
		client: &http.Client{Timeout: timeout},
	}
}

func (e *DefaultHTTPExecutor) Execute(ctx context.Context, method string, url string, headers map[string]string, body []byte) (int, []byte, error) {
	var reqBody io.Reader
	if body != nil {
		reqBody = bytes.NewReader(body)
	}

	req, err := http.NewRequestWithContext(ctx, method, url, reqBody)
	if err != nil {
		return 0, nil, fmt.Errorf("creating request: %w", err)
	}

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	resp, err := e.client.Do(req)
	if err != nil {
		return 0, nil, fmt.Errorf("executing request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return resp.StatusCode, nil, fmt.Errorf("reading response: %w", err)
	}

	return resp.StatusCode, respBody, nil
}
