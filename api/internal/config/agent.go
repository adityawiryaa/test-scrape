package config

import (
	"strconv"
	"time"
)

type AgentConfig struct {
	Hostname       string
	IPAddress      string
	Port           int
	ControllerURL  string
	WorkerURL      string
	APIKey         string
	PollInterval   time.Duration
	RequestTimeout time.Duration
}

func LoadAgentConfig() *AgentConfig {
	port, _ := strconv.Atoi(getEnv("AGENT_PORT", "8081"))
	pollSec, _ := strconv.Atoi(getEnv("POLL_INTERVAL_SECONDS", "30"))
	timeoutSec, _ := strconv.Atoi(getEnv("REQUEST_TIMEOUT_SECONDS", "10"))

	return &AgentConfig{
		Hostname:       getEnv("AGENT_HOSTNAME", "agent-01"),
		IPAddress:      getEnv("AGENT_IP", "127.0.0.1"),
		Port:           port,
		ControllerURL:  getEnv("CONTROLLER_URL", "http://localhost:6001"),
		WorkerURL:      getEnv("WORKER_URL", "http://localhost:6002"),
		APIKey:         getEnv("API_KEY", "default-api-key"),
		PollInterval:   time.Duration(pollSec) * time.Second,
		RequestTimeout: time.Duration(timeoutSec) * time.Second,
	}
}
