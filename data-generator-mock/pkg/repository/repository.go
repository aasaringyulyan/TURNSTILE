package repository

import (
	"data-generator-mock/models"
	"data-generator-mock/pkg/logging"
	"github.com/jmoiron/sqlx"
)

type DataGenerator interface {
	GetEmployees() ([]models.Employee, error)
	AddEmployee(employee models.Employee) error
	SaveSlice(data []models.Employee) error
}

type Repository struct {
	DataGenerator
}

func NewRepository(logger logging.Logger, db *sqlx.DB) *Repository {
	return &Repository{
		DataGenerator: NewDataGeneratorPostgres(logger, db),
	}
}
