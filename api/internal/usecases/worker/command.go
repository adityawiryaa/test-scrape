package usecases

import (
	"context"

	"github.com/adityawiryaa/api/internal/repository/memory"
	hitqueue "github.com/adityawiryaa/api/pkg/hit/queue"
	domainuc "github.com/adityawiryaa/api/domain/usecases"
)

type HTTPExecutor interface {
	Execute(ctx context.Context, method string, url string, headers map[string]string, body []byte) (int, []byte, error)
}

type HitEnqueuer interface {
	EnqueueExecuteHit(payload *hitqueue.ExecuteHitPayload) error
}

type commandUsecase struct {
	executor    HTTPExecutor
	store       *memory.ConfigStore
	queueClient HitEnqueuer
}

func NewCommandUsecase(executor HTTPExecutor, store *memory.ConfigStore, queueClient *hitqueue.Client) domainuc.UsecaseWorkerCommand {
	return &commandUsecase{
		executor:    executor,
		store:       store,
		queueClient: queueClient,
	}
}

func NewCommandUsecaseWithEnqueuer(executor HTTPExecutor, store *memory.ConfigStore, queueClient HitEnqueuer) domainuc.UsecaseWorkerCommand {
	return &commandUsecase{
		executor:    executor,
		store:       store,
		queueClient: queueClient,
	}
}
