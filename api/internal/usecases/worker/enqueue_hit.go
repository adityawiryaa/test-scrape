package usecases

import (
	"context"
	"fmt"
	"log"

	"github.com/google/uuid"

	domainuc "github.com/adityawiryaa/api/domain/usecases"
	"github.com/adityawiryaa/api/domain/valueobject"
	hitqueue "github.com/adityawiryaa/api/pkg/hit/queue"
)

func (c *commandUsecase) EnqueueHit(_ context.Context) (*domainuc.EnqueueHitResponse, error) {
	cfg := c.store.Get()
	if cfg == nil {
		return nil, fmt.Errorf("no config available")
	}

	url, ok := cfg.Data["url"]
	if !ok || url == "" {
		return nil, fmt.Errorf("no url configured")
	}

	taskID := uuid.New().String()
	log.Printf("[enqueue] creating hit task: id=%s url=%s", taskID, url)

	payload := &hitqueue.ExecuteHitPayload{
		TaskID: taskID,
		URL:    url,
		Method: "GET",
	}

	if err := c.queueClient.EnqueueExecuteHit(payload); err != nil {
		log.Printf("[enqueue] failed to enqueue task: id=%s error=%v", taskID, err)
		return nil, fmt.Errorf("enqueuing hit: %w", err)
	}

	log.Printf("[enqueue] task enqueued successfully: id=%s", taskID)

	return &domainuc.EnqueueHitResponse{
		TaskID: taskID,
		Status: valueobject.TaskStatusQueued,
	}, nil
}
