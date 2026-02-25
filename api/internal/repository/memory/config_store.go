package memory

import (
	"sync"

	"github.com/adityawiryaa/api/domain/entity"
)

type ConfigStore struct {
	mu     sync.RWMutex
	config *entity.Config
}

func NewConfigStore() *ConfigStore {
	return &ConfigStore{}
}

func (s *ConfigStore) Get() *entity.Config {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.config
}

func (s *ConfigStore) Set(cfg *entity.Config) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.config = cfg
}

func (s *ConfigStore) Version() int64 {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.config == nil {
		return 0
	}
	return s.config.Version
}
