package queue

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

const resultTTL = 1 * time.Hour

type HitResult struct {
	TaskID     string `json:"task_id"`
	Status     string `json:"status"`
	StatusCode int    `json:"status_code,omitempty"`
	Body       string `json:"body,omitempty"`
	Error      string `json:"error,omitempty"`
}

type ResultStore struct {
	rdb *redis.Client
}

func NewResultStore(rdb *redis.Client) *ResultStore {
	return &ResultStore{rdb: rdb}
}

func (s *ResultStore) SaveResult(ctx context.Context, result *HitResult) error {
	data, err := json.Marshal(result)
	if err != nil {
		return fmt.Errorf("marshaling result: %w", err)
	}

	key := resultKey(result.TaskID)
	return s.rdb.Set(ctx, key, data, resultTTL).Err()
}

func (s *ResultStore) GetResult(ctx context.Context, taskID string) (*HitResult, error) {
	key := resultKey(taskID)
	data, err := s.rdb.Get(ctx, key).Bytes()
	if err == redis.Nil {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("getting result: %w", err)
	}

	var result HitResult
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("unmarshaling result: %w", err)
	}

	return &result, nil
}

func resultKey(taskID string) string {
	return "worker:hit:result:" + taskID
}
