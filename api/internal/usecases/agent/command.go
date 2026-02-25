package usecases

import (
	"github.com/adityawiryaa/api/domain/usecases"
	"github.com/adityawiryaa/api/internal/repository/memory"
	"github.com/adityawiryaa/api/pkg/backoff"
)

type commandUsecase struct {
	controllerClient usecases.ControllerClient
	workerClient     usecases.WorkerClient
	store            *memory.ConfigStore
	backoffCfg       backoff.Config
}

func NewCommandUsecase(
	controllerClient usecases.ControllerClient,
	workerClient usecases.WorkerClient,
	store *memory.ConfigStore,
	backoffCfg backoff.Config,
) usecases.UsecaseAgentCommand {
	return &commandUsecase{
		controllerClient: controllerClient,
		workerClient:     workerClient,
		store:            store,
		backoffCfg:       backoffCfg,
	}
}
