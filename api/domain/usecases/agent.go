package usecases

import (
	"context"
	"time"
	"github.com/adityawiryaa/api/domain/dto"
	"github.com/adityawiryaa/api/domain/entity"
)

type ControllerClient interface {
	Register(ctx context.Context, req *entity.RegistrationRequest) (*entity.RegistrationResponse, error)
	FetchConfig(ctx context.Context, currentVersion int64) (*entity.Config, bool, error)
}

type WorkerClient interface {
	PushConfig(ctx context.Context, cfg *entity.Config) error
}

type UsecaseAgentCommand interface {
	RegisterWithController(ctx context.Context, req *entity.RegistrationRequest) (*dto.RegistrationResponseDTO, error)
	ForwardConfigToWorker(ctx context.Context) error
}

type UsecaseAgentQuery interface {
	PollConfig(ctx context.Context) (int, error)
	StartPolling(ctx context.Context, initialInterval time.Duration, forwardFunc func(context.Context) error)
}
