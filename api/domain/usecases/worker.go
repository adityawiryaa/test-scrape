package usecases

import (
	"context"

	"github.com/adityawiryaa/api/domain/dto"
	"github.com/adityawiryaa/api/domain/entity"
	hitqueue "github.com/adityawiryaa/api/pkg/hit/queue"
)

type EnqueueHitResponse struct {
	TaskID string `json:"task_id"`
	Status string `json:"status"`
}

type HitResultResponse struct {
	TaskID     string `json:"task_id"`
	Status     string `json:"status"`
	StatusCode int    `json:"status_code,omitempty"`
	Body       string `json:"body,omitempty"`
	Error      string `json:"error,omitempty"`
}

type UsecaseWorkerCommand interface {
	ReceiveConfig(cfg *entity.Config)
	EnqueueHit(ctx context.Context) (*EnqueueHitResponse, error)
}

type UsecaseWorkerQuery interface {
	CurrentConfig() *dto.ConfigDTO
	GetHitResult(ctx context.Context, taskID string) (*HitResultResponse, error)
}

func ToHitResultResponse(r *hitqueue.HitResult) *HitResultResponse {
	return &HitResultResponse{
		TaskID:     r.TaskID,
		Status:     r.Status,
		StatusCode: r.StatusCode,
		Body:       r.Body,
		Error:      r.Error,
	}
}
