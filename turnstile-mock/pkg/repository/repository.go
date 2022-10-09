package repository

import (
	"github.com/jmoiron/sqlx"
	"turnstile-mock/models"
	"turnstile-mock/pkg/logging"
)

type Passage interface {
	GetEmployees() ([]models.Employee, error)
}

type Repository struct {
	Passage
}

func NewRepository(logger logging.Logger, db *sqlx.DB) *Repository {
	return &Repository{
		Passage: NewPassagePostgres(logger, db),
	}
}
