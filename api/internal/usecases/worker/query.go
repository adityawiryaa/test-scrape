package usecases

import (
	"context"
	"fmt"

	"github.com/adityawiryaa/api/domain/dto"
	"github.com/adityawiryaa/api/domain/dto/mapper"
	domainuc "github.com/adityawiryaa/api/domain/usecases"
	"github.com/adityawiryaa/api/domain/valueobject"
	"github.com/adityawiryaa/api/internal/repository/memory"
	hitqueue "github.com/adityawiryaa/api/pkg/hit/queue"
)

type queryUsecase struct {
	store       *memory.ConfigStore
	resultStore *hitqueue.ResultStore
}

func NewQueryUsecase(store *memory.ConfigStore, resultStore *hitqueue.ResultStore) domainuc.UsecaseWorkerQuery {
	return &queryUsecase{
		store:       store,
		resultStore: resultStore,
	}
}

func NewQueryUsecaseWithStore(store *memory.ConfigStore) domainuc.UsecaseWorkerQuery {
	return &queryUsecase{
		store: store,
	}
}

func (q *queryUsecase) CurrentConfig() *dto.ConfigDTO {
	cfg := q.store.Get()
	if cfg == nil {
		return nil
	}
	result := mapper.ToConfigDTO(cfg)
	return &result
}

func (q *queryUsecase) GetHitResult(ctx context.Context, taskID string) (*domainuc.HitResultResponse, error) {
	if q.resultStore == nil {
		return nil, fmt.Errorf("result store not configured")
	}

	result, err := q.resultStore.GetResult(ctx, taskID)
	if err != nil {
		return nil, fmt.Errorf("getting hit result: %w", err)
	}

	if result == nil {
		return &domainuc.HitResultResponse{
			TaskID: taskID,
			Status: valueobject.TaskStatusPending,
		}, nil
	}

	return domainuc.ToHitResultResponse(result), nil
}
