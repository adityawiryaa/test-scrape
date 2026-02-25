package controller

import (
	"github.com/adityawiryaa/api/domain/usecases"
)

type Handler struct {
	commandUC usecases.UsecaseControllerCommand
	queryUC   usecases.UsecaseControllerQuery
}

func NewHandler(commandUC usecases.UsecaseControllerCommand, queryUC usecases.UsecaseControllerQuery) *Handler {
	return &Handler{
		commandUC: commandUC,
		queryUC:   queryUC,
	}
}
