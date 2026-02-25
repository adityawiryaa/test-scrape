package agent_test

import (
	"context"
	"errors"
	"testing"

	"github.com/adityawiryaa/api/domain/entity"
	"github.com/adityawiryaa/api/internal/repository/memory"
	agent "github.com/adityawiryaa/api/internal/usecases/agent"
)

func TestPollConfig(t *testing.T) {
	tests := []struct {
		name         string
		storeVersion int64
		fetchCfg     *entity.Config
		fetchChanged bool
		fetchErr     error
		wantErr      bool
		wantInterval int
	}{
		{
			name:         "no change",
			storeVersion: 5,
			fetchChanged: false,
			wantErr:      false,
			wantInterval: 0,
		},
		{
			name:         "config updated",
			storeVersion: 1,
			fetchCfg:     &entity.Config{Version: 2, PollIntervalSeconds: 15},
			fetchChanged: true,
			wantErr:      false,
			wantInterval: 15,
		},
		{
			name:     "fetch error",
			fetchErr: errors.New("network error"),
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := memory.NewConfigStore()
			if tt.storeVersion > 0 {
				store.Set(&entity.Config{Version: tt.storeVersion})
			}

			client := &mockControllerClient{
				registerFunc: func(_ context.Context, _ *entity.RegistrationRequest) (*entity.RegistrationResponse, error) {
					return nil, nil
				},
				fetchFunc: func(_ context.Context, _ int64) (*entity.Config, bool, error) {
					return tt.fetchCfg, tt.fetchChanged, tt.fetchErr
				},
			}

			uc := agent.NewQueryUsecase(client, store)
			interval, err := uc.PollConfig(context.Background())

			if tt.wantErr {
				if err == nil {
					t.Error("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if interval != tt.wantInterval {
				t.Errorf("interval = %d, want %d", interval, tt.wantInterval)
			}
		})
	}
}
