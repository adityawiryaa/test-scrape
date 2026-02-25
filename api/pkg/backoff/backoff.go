package backoff

import (
	"math"
	"math/rand"
	"time"
)

type Config struct {
	InitialInterval time.Duration
	MaxInterval     time.Duration
	Multiplier      float64
	MaxRetries      int
}

func DefaultConfig() Config {
	return Config{
		InitialInterval: 1 * time.Second,
		MaxInterval:     30 * time.Second,
		Multiplier:      2.0,
		MaxRetries:      10,
	}
}

func NextInterval(cfg Config, attempt int) time.Duration {
	if attempt >= cfg.MaxRetries {
		return cfg.MaxInterval
	}
	interval := float64(cfg.InitialInterval) * math.Pow(cfg.Multiplier, float64(attempt))
	if interval > float64(cfg.MaxInterval) {
		interval = float64(cfg.MaxInterval)
	}
	jitter := rand.Float64() * 0.5 * interval
	return time.Duration(interval + jitter)
}
