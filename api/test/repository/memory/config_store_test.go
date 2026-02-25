package memory_test

import (
	"sync"
	"testing"

	"github.com/adityawiryaa/api/domain/entity"
	"github.com/adityawiryaa/api/internal/repository/memory"
)

func TestConfigStore(t *testing.T) {
	tests := []struct {
		name        string
		setup       func(s *memory.ConfigStore)
		wantVersion int64
		wantNil     bool
	}{
		{
			name:        "empty store returns nil",
			setup:       func(s *memory.ConfigStore) {},
			wantVersion: 0,
			wantNil:     true,
		},
		{
			name: "set and get config",
			setup: func(s *memory.ConfigStore) {
				s.Set(&entity.Config{Version: 5, Data: map[string]string{"key": "value"}})
			},
			wantVersion: 5,
			wantNil:     false,
		},
		{
			name: "overwrite config",
			setup: func(s *memory.ConfigStore) {
				s.Set(&entity.Config{Version: 1})
				s.Set(&entity.Config{Version: 2})
			},
			wantVersion: 2,
			wantNil:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := memory.NewConfigStore()
			tt.setup(store)

			got := store.Get()
			if tt.wantNil && got != nil {
				t.Error("expected nil config")
			}
			if !tt.wantNil && got == nil {
				t.Error("expected non-nil config")
			}
			if store.Version() != tt.wantVersion {
				t.Errorf("Version() = %d, want %d", store.Version(), tt.wantVersion)
			}
		})
	}
}

func TestConfigStoreConcurrency(t *testing.T) {
	store := memory.NewConfigStore()
	var wg sync.WaitGroup

	for i := range 100 {
		wg.Add(2)
		go func(v int64) {
			defer wg.Done()
			store.Set(&entity.Config{Version: v})
		}(int64(i))
		go func() {
			defer wg.Done()
			store.Get()
			store.Version()
		}()
	}
	wg.Wait()
}
