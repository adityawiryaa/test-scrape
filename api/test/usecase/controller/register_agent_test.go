package controller_test

import (
	"context"
	"errors"
	"testing"

	"github.com/adityawiryaa/api/domain/entity"
	"github.com/adityawiryaa/api/domain/request"
	controller "github.com/adityawiryaa/api/internal/usecases/controller"
)

type mockAgentCommand struct {
	saveFunc func(ctx context.Context, agent *entity.Agent) error
}

func (m *mockAgentCommand) Save(ctx context.Context, agent *entity.Agent) error {
	return m.saveFunc(ctx, agent)
}

func TestRegisterAgent(t *testing.T) {
	tests := []struct {
		name    string
		req     *request.RegisterAgentRequest
		saveErr error
		wantErr bool
	}{
		{
			name: "successful registration",
			req: &request.RegisterAgentRequest{
				Hostname:  "agent-01",
				IPAddress: "192.168.1.10",
				Port:      8081,
			},
			saveErr: nil,
			wantErr: false,
		},
		{
			name: "save fails",
			req: &request.RegisterAgentRequest{
				Hostname:  "agent-02",
				IPAddress: "192.168.1.11",
				Port:      8082,
			},
			saveErr: errors.New("db error"),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &mockAgentCommand{
				saveFunc: func(_ context.Context, _ *entity.Agent) error {
					return tt.saveErr
				},
			}

			query := &mockConfigQuery{
				getLatestFunc: func(_ context.Context) (*entity.Config, error) {
					return &entity.Config{PollIntervalSeconds: 30}, nil
				},
				getByVersionFunc: func(_ context.Context, _ int64) (*entity.Config, error) {
					return nil, nil
				},
			}

			uc := controller.NewCommandUsecase(cmd, nil, query)
			resp, err := uc.RegisterAgent(context.Background(), tt.req)

			if tt.wantErr {
				if err == nil {
					t.Error("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if resp.AgentID == "" {
				t.Error("expected non-empty agent ID")
			}
			if resp.Status != "active" {
				t.Errorf("status = %s, want active", resp.Status)
			}
			if resp.PollURL != "/config" {
				t.Errorf("poll_url = %s, want /config", resp.PollURL)
			}
			if resp.PollIntervalSeconds != 30 {
				t.Errorf("poll_interval_seconds = %d, want 30", resp.PollIntervalSeconds)
			}
		})
	}
}
