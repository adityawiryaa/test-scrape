package worker

import (
	"github.com/adityawiryaa/api/domain/usecases"
)

type Handler struct {
	commandUC usecases.UsecaseWorkerCommand
	queryUC   usecases.UsecaseWorkerQuery
}

func NewHandler(commandUC usecases.UsecaseWorkerCommand, queryUC usecases.UsecaseWorkerQuery) *Handler {
	return &Handler{
		commandUC: commandUC,
		queryUC:   queryUC,
	}
}
