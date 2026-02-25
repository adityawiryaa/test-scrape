package config

import "time"

type WorkerConfig struct {
	Port           string
	APIKey         string
	RequestTimeout time.Duration
	Redis          *RedisConfig
}

func LoadWorkerConfig() *WorkerConfig {
	return &WorkerConfig{
		Port:           getEnv("WORKER_PORT", "6002"),
		APIKey:         getEnv("API_KEY", "default-api-key"),
		RequestTimeout: 30 * time.Second,
		Redis:          LoadRedisConfig(),
	}
}
