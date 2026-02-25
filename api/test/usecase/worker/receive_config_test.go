package worker_test

import (
	"testing"

	"github.com/adityawiryaa/api/domain/entity"
	"github.com/adityawiryaa/api/internal/repository/memory"
	worker "github.com/adityawiryaa/api/internal/usecases/worker"
)

func TestReceiveConfig(t *testing.T) {
	tests := []struct {
		name        string
		configs     []*entity.Config
		wantVersion int64
	}{
		{
			name:        "store single config",
			configs:     []*entity.Config{{Version: 1, Data: map[string]string{"a": "b"}}},
			wantVersion: 1,
		},
		{
			name: "latest config wins",
			configs: []*entity.Config{
				{Version: 1, Data: map[string]string{"a": "b"}},
				{Version: 5, Data: map[string]string{"c": "d"}},
			},
			wantVersion: 5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := memory.NewConfigStore()
			cmdUc := worker.NewCommandUsecaseWithEnqueuer(nil, store, &mockQueueClient{})
			queryUc := worker.NewQueryUsecaseWithStore(store)

			for _, cfg := range tt.configs {
				cmdUc.ReceiveConfig(cfg)
			}

			current := queryUc.CurrentConfig()
			if current == nil {
				t.Fatal("expected non-nil config")
			}
			if current.Version != tt.wantVersion {
				t.Errorf("version = %d, want %d", current.Version, tt.wantVersion)
			}
		})
	}
}
