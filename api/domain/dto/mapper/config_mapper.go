package mapper

import (
	"github.com/adityawiryaa/api/domain/dto"
	"github.com/adityawiryaa/api/domain/entity"
)

func ToConfigDTO(cfg *entity.Config) dto.ConfigDTO {
	return dto.ConfigDTO{
		ID:                  cfg.ID,
		Version:             cfg.Version,
		Data:                cfg.Data,
		PollIntervalSeconds: cfg.PollIntervalSeconds,
	}
}
