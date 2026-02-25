package config

import "os"

type ControllerConfig struct {
	Port   string
	DBPath string
	APIKey string
}

func LoadControllerConfig() *ControllerConfig {
	return &ControllerConfig{
		Port:   getEnv("CONTROLLER_PORT", "6001"),
		DBPath: getEnv("CONTROLLER_DB_PATH", "controller.db"),
		APIKey: getEnv("API_KEY", "default-api-key"),
	}
}

func getEnv(key string, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}
