package usecases

import (
	"github.com/adityawiryaa/api/domain/usecases"
	"github.com/adityawiryaa/api/internal/repository/memory"
)

type queryUsecase struct {
	client usecases.ControllerClient
	store  *memory.ConfigStore
}

func NewQueryUsecase(
	client usecases.ControllerClient,
	store *memory.ConfigStore,
) usecases.UsecaseAgentQuery {
	return &queryUsecase{
		client: client,
		store:  store,
	}
}
