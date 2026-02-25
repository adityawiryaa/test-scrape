package usecases

import (
	"log"

	"github.com/adityawiryaa/api/domain/entity"
)

func (c *commandUsecase) ReceiveConfig(cfg *entity.Config) {
	log.Printf("[config] received config version=%d", cfg.Version)
	c.store.Set(cfg)
}
