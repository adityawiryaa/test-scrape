package controller_test

import (
	"context"
	"errors"
	"testing"

	"github.com/adityawiryaa/api/domain/entity"
	"github.com/adityawiryaa/api/domain/request"
	controller "github.com/adityawiryaa/api/internal/usecases/controller"
)

type mockConfigCommand struct {
	saveFunc func(ctx context.Context, cfg *entity.Config) error
}

func (m *mockConfigCommand) SaveConfig(ctx context.Context, cfg *entity.Config) error {
	return m.saveFunc(ctx, cfg)
}

type mockConfigQuery struct {
	getLatestFunc    func(ctx context.Context) (*entity.Config, error)
	getByVersionFunc func(ctx context.Context, version int64) (*entity.Config, error)
}

func (m *mockConfigQuery) GetLatestConfig(ctx context.Context) (*entity.Config, error) {
	return m.getLatestFunc(ctx)
}

func (m *mockConfigQuery) GetConfigByVersion(ctx context.Context, version int64) (*entity.Config, error) {
	return m.getByVersionFunc(ctx, version)
}

func TestUpdateConfig(t *testing.T) {
	tests := []struct {
		name        string
		req         *request.UpdateConfigRequest
		latest      *entity.Config
		latestErr   error
		saveErr     error
		wantErr     bool
		wantVersion int64
	}{
		{
			name:        "first config",
			req:         &request.UpdateConfigRequest{Data: map[string]string{"key": "value"}, PollIntervalSeconds: 15},
			latest:      nil,
			latestErr:   errors.New("no rows"),
			saveErr:     nil,
			wantErr:     false,
			wantVersion: 1,
		},
		{
			name:        "increments version",
			req:         &request.UpdateConfigRequest{Data: map[string]string{"key": "value2"}},
			latest:      &entity.Config{Version: 3},
			latestErr:   nil,
			saveErr:     nil,
			wantErr:     false,
			wantVersion: 4,
		},
		{
			name:      "save fails",
			req:       &request.UpdateConfigRequest{Data: map[string]string{"k": "v"}},
			latest:    nil,
			latestErr: errors.New("no rows"),
			saveErr:   errors.New("db error"),
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &mockConfigCommand{
				saveFunc: func(_ context.Context, _ *entity.Config) error {
					return tt.saveErr
				},
			}
			query := &mockConfigQuery{
				getLatestFunc: func(_ context.Context) (*entity.Config, error) {
					return tt.latest, tt.latestErr
				},
				getByVersionFunc: func(_ context.Context, _ int64) (*entity.Config, error) {
					return nil, nil
				},
			}

			uc := controller.NewCommandUsecase(nil, cmd, query)
			cfg, err := uc.UpdateConfig(context.Background(), tt.req)

			if tt.wantErr {
				if err == nil {
					t.Error("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if cfg.Version != tt.wantVersion {
				t.Errorf("version = %d, want %d", cfg.Version, tt.wantVersion)
			}
		})
	}
}
