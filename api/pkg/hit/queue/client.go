package queue

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/hibiken/asynq"
)

type Client struct {
	client *asynq.Client
}

func NewClient(redisAddr string, db int) *Client {
	log.Printf("[queue] client connected to redis: addr=%s db=%d", redisAddr, db)
	return &Client{
		client: asynq.NewClient(asynq.RedisClientOpt{
			Addr: redisAddr,
			DB:   db,
		}),
	}
}

func (c *Client) Close() error {
	return c.client.Close()
}

func (c *Client) EnqueueExecuteHit(payload *ExecuteHitPayload) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshaling payload: %w", err)
	}

	task := asynq.NewTask(TypeHitExecute, data)
	_, err = c.client.Enqueue(task, getTaskOptions()...)
	if err != nil {
		return fmt.Errorf("enqueuing task: %w", err)
	}

	return nil
}

func getTaskOptions() []asynq.Option {
	maxRetry := 3
	if v := os.Getenv("WORKER_RETRY_MAX"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			maxRetry = n
		}
	}

	return []asynq.Option{
		asynq.MaxRetry(maxRetry),
		asynq.Queue("default"),
	}
}
