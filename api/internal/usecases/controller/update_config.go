package usecases

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/adityawiryaa/api/domain/dto"
	"github.com/adityawiryaa/api/domain/dto/mapper"
	"github.com/adityawiryaa/api/domain/entity"
	"github.com/adityawiryaa/api/domain/request"
)

func (c *commandUsecase) UpdateConfig(ctx context.Context, req *request.UpdateConfigRequest) (*dto.ConfigDTO, error) {
	var nextVersion int64 = 1

	latest, err := c.configRepoQuery.GetLatestConfig(ctx)
	if err == nil && latest != nil {
		nextVersion = latest.Version + 1
	}

	pollInterval := 30
	if req.PollIntervalSeconds > 0 {
		pollInterval = req.PollIntervalSeconds
	}

	cfg := &entity.Config{
		ID:                  uuid.New().String(),
		Version:             nextVersion,
		Data:                req.Data,
		PollIntervalSeconds: pollInterval,
		CreatedAt:           time.Now(),
	}

	if err := c.configRepoCommand.SaveConfig(ctx, cfg); err != nil {
		return nil, err
	}

	result := mapper.ToConfigDTO(cfg)
	return &result, nil
}
