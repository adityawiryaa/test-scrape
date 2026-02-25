package usecases

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/adityawiryaa/api/domain/dto"
	"github.com/adityawiryaa/api/domain/dto/mapper"
	"github.com/adityawiryaa/api/domain/entity"
	"github.com/adityawiryaa/api/domain/request"
	"github.com/adityawiryaa/api/domain/valueobject"
)

func (c *commandUsecase) RegisterAgent(ctx context.Context, req *request.RegisterAgentRequest) (*dto.RegistrationResponseDTO, error) {
	now := time.Now()
	agent := &entity.Agent{
		ID:        uuid.New().String(),
		Hostname:  req.Hostname,
		IPAddress: req.IPAddress,
		Port:      req.Port,
		Status:    valueobject.StatusActive,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := c.agentRepoCommand.Save(ctx, agent); err != nil {
		return nil, err
	}

	pollInterval := 30
	latestCfg, err := c.configRepoQuery.GetLatestConfig(ctx)
	if err == nil && latestCfg != nil && latestCfg.PollIntervalSeconds > 0 {
		pollInterval = latestCfg.PollIntervalSeconds
	}

	result := mapper.ToRegistrationResponseDTO(agent, pollInterval)
	return &result, nil
}
