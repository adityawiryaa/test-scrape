package agent_test

import (
	"context"
	"errors"
	"testing"

	"github.com/adityawiryaa/api/domain/entity"
	"github.com/adityawiryaa/api/internal/repository/memory"
	agent "github.com/adityawiryaa/api/internal/usecases/agent"
	"github.com/adityawiryaa/api/pkg/backoff"
)

type mockWorkerClient struct {
	pushFunc func(ctx context.Context, cfg *entity.Config) error
}

func (m *mockWorkerClient) PushConfig(ctx context.Context, cfg *entity.Config) error {
	return m.pushFunc(ctx, cfg)
}

func TestForwardConfigToWorker(t *testing.T) {
	tests := []struct {
		name    string
		config  *entity.Config
		pushErr error
		wantErr bool
	}{
		{
			name:    "successful forward",
			config:  &entity.Config{Version: 3, Data: map[string]string{"k": "v"}},
			wantErr: false,
		},
		{
			name:    "no config to forward",
			config:  nil,
			wantErr: true,
		},
		{
			name:    "push fails",
			config:  &entity.Config{Version: 1},
			pushErr: errors.New("worker unreachable"),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := memory.NewConfigStore()
			if tt.config != nil {
				store.Set(tt.config)
			}

			workerClient := &mockWorkerClient{
				pushFunc: func(_ context.Context, _ *entity.Config) error {
					return tt.pushErr
				},
			}

			controllerClient := &mockControllerClient{
				registerFunc: func(_ context.Context, _ *entity.RegistrationRequest) (*entity.RegistrationResponse, error) {
					return nil, nil
				},
				fetchFunc: func(_ context.Context, _ int64) (*entity.Config, bool, error) {
					return nil, false, nil
				},
			}

			cfg := backoff.Config{
				InitialInterval: 0,
				MaxInterval:     0,
				Multiplier:      1,
				MaxRetries:      1,
			}

			uc := agent.NewCommandUsecase(controllerClient, workerClient, store, cfg)
			err := uc.ForwardConfigToWorker(context.Background())

			if tt.wantErr {
				if err == nil {
					t.Error("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}
