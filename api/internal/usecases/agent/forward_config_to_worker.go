package usecases

import (
	"context"
	"fmt"
)

func (c *commandUsecase) ForwardConfigToWorker(ctx context.Context) error {
	cfg := c.store.Get()
	if cfg == nil {
		return fmt.Errorf("no config to forward")
	}
	return c.workerClient.PushConfig(ctx, cfg)
}
