package service

import (
	"data-generator-mock/models"
	"data-generator-mock/pkg/logging"
	"data-generator-mock/pkg/repository"
)

type DataGenerator interface {
	GetByRv(rv int64) ([]models.Employee, error)
	GenNewEmployee() error
	GenSlice(n int) error
}

type Service struct {
	DataGenerator
}

func NewService(logger logging.Logger, repos *repository.Repository) *Service {
	return &Service{
		DataGenerator: NewDataGeneratorService(logger, repos.DataGenerator),
	}
}
