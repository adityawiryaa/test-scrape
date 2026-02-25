package backoff_test

import (
	"testing"
	"time"

	"github.com/adityawiryaa/api/pkg/backoff"
)

func TestNextInterval(t *testing.T) {
	tests := []struct {
		name    string
		cfg     backoff.Config
		attempt int
		wantMin time.Duration
		wantMax time.Duration
	}{
		{
			name:    "first attempt",
			cfg:     backoff.DefaultConfig(),
			attempt: 0,
			wantMin: 1 * time.Second,
			wantMax: 2 * time.Second,
		},
		{
			name:    "second attempt",
			cfg:     backoff.DefaultConfig(),
			attempt: 1,
			wantMin: 2 * time.Second,
			wantMax: 4 * time.Second,
		},
		{
			name:    "exceeds max retries",
			cfg:     backoff.DefaultConfig(),
			attempt: 15,
			wantMin: 30 * time.Second,
			wantMax: 45 * time.Second,
		},
		{
			name: "custom config",
			cfg: backoff.Config{
				InitialInterval: 500 * time.Millisecond,
				MaxInterval:     10 * time.Second,
				Multiplier:      3.0,
				MaxRetries:      5,
			},
			attempt: 2,
			wantMin: 4500 * time.Millisecond,
			wantMax: 7000 * time.Millisecond,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := backoff.NextInterval(tt.cfg, tt.attempt)
			if got < tt.wantMin || got > tt.wantMax {
				t.Errorf("NextInterval() = %v, want between %v and %v", got, tt.wantMin, tt.wantMax)
			}
		})
	}
}
