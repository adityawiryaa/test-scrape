package worker_test

import (
	"context"
	"testing"

	"github.com/adityawiryaa/api/domain/entity"
	"github.com/adityawiryaa/api/domain/valueobject"
	"github.com/adityawiryaa/api/internal/repository/memory"
	worker "github.com/adityawiryaa/api/internal/usecases/worker"
	hitqueue "github.com/adityawiryaa/api/pkg/hit/queue"
)

type mockExecutor struct {
	executeFunc func(ctx context.Context, method string, url string, headers map[string]string, body []byte) (int, []byte, error)
}

func (m *mockExecutor) Execute(ctx context.Context, method string, url string, headers map[string]string, body []byte) (int, []byte, error) {
	return m.executeFunc(ctx, method, url, headers, body)
}

type mockQueueClient struct {
	enqueueFunc func(payload *hitqueue.ExecuteHitPayload) error
}

func (m *mockQueueClient) EnqueueExecuteHit(payload *hitqueue.ExecuteHitPayload) error {
	if m.enqueueFunc != nil {
		return m.enqueueFunc(payload)
	}
	return nil
}

func TestEnqueueHit(t *testing.T) {
	tests := []struct {
		name       string
		config     *entity.Config
		enqueueErr error
		wantErr    bool
		wantStatus string
	}{
		{
			name:       "successful enqueue",
			config:     &entity.Config{Version: 3, Data: map[string]string{"url": "https://example.com/api"}},
			wantErr:    false,
			wantStatus: valueobject.TaskStatusQueued,
		},
		{
			name:    "no config available",
			config:  nil,
			wantErr: true,
		},
		{
			name:    "no url in config",
			config:  &entity.Config{Version: 1, Data: map[string]string{"other": "value"}},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := memory.NewConfigStore()
			if tt.config != nil {
				store.Set(tt.config)
			}

			executor := &mockExecutor{
				executeFunc: func(_ context.Context, _ string, _ string, _ map[string]string, _ []byte) (int, []byte, error) {
					return 0, nil, nil
				},
			}

			uc := worker.NewCommandUsecaseWithEnqueuer(executor, store, &mockQueueClient{})
			resp, err := uc.EnqueueHit(context.Background())

			if tt.wantErr {
				if err == nil {
					t.Error("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if resp.Status != tt.wantStatus {
				t.Errorf("status = %s, want %s", resp.Status, tt.wantStatus)
			}
			if resp.TaskID == "" {
				t.Error("expected non-empty task ID")
			}
		})
	}
}
