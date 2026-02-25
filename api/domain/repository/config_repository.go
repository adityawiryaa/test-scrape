package repository

import (
	"context"
	"github.com/adityawiryaa/api/domain/entity"
)

type ConfigRepositoryCommand interface {
	SaveConfig(ctx context.Context, cfg *entity.Config) error
}

type ConfigRepositoryQuery interface {
	GetLatestConfig(ctx context.Context) (*entity.Config, error)
	GetConfigByVersion(ctx context.Context, version int64) (*entity.Config, error)
}
