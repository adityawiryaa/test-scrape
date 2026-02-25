package usecases

import (
	"context"

	"github.com/adityawiryaa/api/domain/dto"
	"github.com/adityawiryaa/api/domain/dto/mapper"
)

func (q *queryUsecase) GetLatestConfig(ctx context.Context) (*dto.ConfigDTO, error) {
	cfg, err := q.configRepoQuery.GetLatestConfig(ctx)
	if err != nil {
		return nil, err
	}
	result := mapper.ToConfigDTO(cfg)
	return &result, nil
}

func (q *queryUsecase) GetConfigByVersion(ctx context.Context, version int64) (*dto.ConfigDTO, error) {
	cfg, err := q.configRepoQuery.GetConfigByVersion(ctx, version)
	if err != nil {
		return nil, err
	}
	result := mapper.ToConfigDTO(cfg)
	return &result, nil
}
