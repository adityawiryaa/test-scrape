package usecases

import (
	"github.com/adityawiryaa/api/domain/repository"
	domainuc "github.com/adityawiryaa/api/domain/usecases"
)

type commandUsecase struct {
	agentRepoCommand  repository.AgentRepositoryCommand
	configRepoCommand repository.ConfigRepositoryCommand
	configRepoQuery   repository.ConfigRepositoryQuery
}

func NewCommandUsecase(
	agentRepoCommand repository.AgentRepositoryCommand,
	configRepoCommand repository.ConfigRepositoryCommand,
	configRepoQuery repository.ConfigRepositoryQuery,
) domainuc.UsecaseControllerCommand {
	return &commandUsecase{
		agentRepoCommand:  agentRepoCommand,
		configRepoCommand: configRepoCommand,
		configRepoQuery:   configRepoQuery,
	}
}
