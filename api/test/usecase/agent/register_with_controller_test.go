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

type mockControllerClient struct {
	registerFunc func(ctx context.Context, req *entity.RegistrationRequest) (*entity.RegistrationResponse, error)
	fetchFunc    func(ctx context.Context, currentVersion int64) (*entity.Config, bool, error)
}

func (m *mockControllerClient) Register(ctx context.Context, req *entity.RegistrationRequest) (*entity.RegistrationResponse, error) {
	return m.registerFunc(ctx, req)
}

func (m *mockControllerClient) FetchConfig(ctx context.Context, currentVersion int64) (*entity.Config, bool, error) {
	return m.fetchFunc(ctx, currentVersion)
}

func TestRegisterWithController(t *testing.T) {
	tests := []struct {
		name    string
		regResp *entity.RegistrationResponse
		regErr  error
		wantErr bool
	}{
		{
			name:    "successful registration",
			regResp: &entity.RegistrationResponse{AgentID: "agent-123", Status: "active"},
			wantErr: false,
		},
		{
			name:    "all retries fail",
			regErr:  errors.New("connection refused"),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := &mockControllerClient{
				registerFunc: func(_ context.Context, _ *entity.RegistrationRequest) (*entity.RegistrationResponse, error) {
					return tt.regResp, tt.regErr
				},
				fetchFunc: func(_ context.Context, _ int64) (*entity.Config, bool, error) {
					return nil, false, nil
				},
			}

			workerClient := &mockWorkerClient{
				pushFunc: func(_ context.Context, _ *entity.Config) error {
					return nil
				},
			}

			cfg := backoff.Config{
				InitialInterval: 0,
				MaxInterval:     0,
				Multiplier:      1,
				MaxRetries:      1,
			}

			store := memory.NewConfigStore()
			uc := agent.NewCommandUsecase(client, workerClient, store, cfg)
			resp, err := uc.RegisterWithController(context.Background(), &entity.RegistrationRequest{
				Hostname:  "test-host",
				IPAddress: "127.0.0.1",
				Port:      8080,
			})

			if tt.wantErr {
				if err == nil {
					t.Error("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if resp.AgentID != tt.regResp.AgentID {
				t.Errorf("agentID = %s, want %s", resp.AgentID, tt.regResp.AgentID)
			}
		})
	}
}
