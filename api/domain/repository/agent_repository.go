package repository

import (
	"context"
	"github.com/adityawiryaa/api/domain/entity"
)

type AgentRepositoryCommand interface {
	Save(ctx context.Context, agent *entity.Agent) error
}

type AgentRepositoryQuery interface {
	FindByID(ctx context.Context, id string) (*entity.Agent, error)
}
