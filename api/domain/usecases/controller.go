package usecases

import (
	"context"
	"github.com/adityawiryaa/api/domain/dto"
	"github.com/adityawiryaa/api/domain/request"
)

type UsecaseControllerCommand interface {
	RegisterAgent(ctx context.Context, req *request.RegisterAgentRequest) (*dto.RegistrationResponseDTO, error)
	UpdateConfig(ctx context.Context, req *request.UpdateConfigRequest) (*dto.ConfigDTO, error)
}

type UsecaseControllerQuery interface {
	GetLatestConfig(ctx context.Context) (*dto.ConfigDTO, error)
	GetConfigByVersion(ctx context.Context, version int64) (*dto.ConfigDTO, error)
}
