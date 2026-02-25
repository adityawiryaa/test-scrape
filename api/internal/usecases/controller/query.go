package usecases

import (
	"github.com/adityawiryaa/api/domain/repository"
	domainuc "github.com/adityawiryaa/api/domain/usecases"
)

type queryUsecase struct {
	configRepoQuery repository.ConfigRepositoryQuery
}

func NewQueryUsecase(configRepoQuery repository.ConfigRepositoryQuery) domainuc.UsecaseControllerQuery {
	return &queryUsecase{
		configRepoQuery: configRepoQuery,
	}
}
