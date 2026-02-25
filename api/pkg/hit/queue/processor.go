package queue

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/adityawiryaa/api/domain/valueobject"
	"github.com/hibiken/asynq"
)

type HTTPExecutor interface {
	Execute(ctx context.Context, method string, url string, headers map[string]string, body []byte) (int, []byte, error)
}

type Processor struct {
	executor HTTPExecutor
	store    *ResultStore
}

func NewProcessor(executor HTTPExecutor, store *ResultStore) *Processor {
	return &Processor{
		executor: executor,
		store:    store,
	}
}

func (p *Processor) HandleExecuteHit(ctx context.Context, t *asynq.Task) error {
	var payload ExecuteHitPayload
	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		return fmt.Errorf("unmarshaling payload: %w", err)
	}

	log.Printf("[processor] picked up task: id=%s type=%s url=%s method=%s", payload.TaskID, TypeHitExecute, payload.URL, payload.Method)

	statusCode, body, err := p.executor.Execute(ctx, payload.Method, payload.URL, nil, nil)
	if err != nil {
		log.Printf("[processor] execution failed: id=%s error=%v", payload.TaskID, err)
		saveErr := p.store.SaveResult(ctx, &HitResult{
			TaskID: payload.TaskID,
			Status: valueobject.TaskStatusFailed,
			Error:  err.Error(),
		})
		if saveErr != nil {
			log.Printf("[processor] failed to save error result: id=%s error=%v", payload.TaskID, saveErr)
		}
		return fmt.Errorf("executing hit: %w", err)
	}

	result := &HitResult{
		TaskID:     payload.TaskID,
		Status:     valueobject.TaskStatusCompleted,
		StatusCode: statusCode,
		Body:       string(body),
	}

	if err := p.store.SaveResult(ctx, result); err != nil {
		log.Printf("[processor] failed to save result: id=%s error=%v", payload.TaskID, err)
		return fmt.Errorf("saving result: %w", err)
	}

	log.Printf("[processor] task completed: id=%s status_code=%d body_size=%d", payload.TaskID, statusCode, len(body))
	return nil
}

func (p *Processor) RegisterHandlers(mux *asynq.ServeMux) {
	mux.HandleFunc(TypeHitExecute, p.HandleExecuteHit)
}
