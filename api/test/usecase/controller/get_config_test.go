package controller_test

import (
	"context"
	"errors"
	"testing"

	"github.com/adityawiryaa/api/domain/entity"
	controller "github.com/adityawiryaa/api/internal/usecases/controller"
)

func TestGetConfig(t *testing.T) {
	tests := []struct {
		name      string
		latest    *entity.Config
		latestErr error
		wantErr   bool
	}{
		{
			name:    "returns latest config",
			latest:  &entity.Config{Version: 5, Data: map[string]string{"env": "prod"}},
			wantErr: false,
		},
		{
			name:      "no config found",
			latestErr: errors.New("no rows"),
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			query := &mockConfigQuery{
				getLatestFunc: func(_ context.Context) (*entity.Config, error) {
					return tt.latest, tt.latestErr
				},
				getByVersionFunc: func(_ context.Context, _ int64) (*entity.Config, error) {
					return nil, nil
				},
			}

			uc := controller.NewQueryUsecase(query)
			cfg, err := uc.GetLatestConfig(context.Background())

			if tt.wantErr {
				if err == nil {
					t.Error("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if cfg.Version != tt.latest.Version {
				t.Errorf("version = %d, want %d", cfg.Version, tt.latest.Version)
			}
		})
	}
}
