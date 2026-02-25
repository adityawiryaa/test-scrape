package usecases

import (
	"context"
	"log"
	"time"

	"github.com/adityawiryaa/api/domain/dto"
	"github.com/adityawiryaa/api/domain/entity"
	"github.com/adityawiryaa/api/pkg/backoff"
)

func (c *commandUsecase) RegisterWithController(ctx context.Context, req *entity.RegistrationRequest) (*dto.RegistrationResponseDTO, error) {
	var lastErr error
	for attempt := range c.backoffCfg.MaxRetries + 1 {
		resp, err := c.controllerClient.Register(ctx, req)
		if err == nil {
			return &dto.RegistrationResponseDTO{
				AgentID: resp.AgentID,
				Status:  resp.Status,
			}, nil
		}
		lastErr = err
		log.Printf("registration attempt %d failed: %v", attempt+1, err)

		if attempt < c.backoffCfg.MaxRetries {
			interval := backoff.NextInterval(c.backoffCfg, attempt)
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-time.After(interval):
			}
		}
	}
	return nil, lastErr
}
