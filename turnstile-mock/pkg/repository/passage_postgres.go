package repository

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"turnstile-mock/models"
	"turnstile-mock/pkg/logging"
)

type PassagePostgres struct {
	logger logging.Logger
	db     *sqlx.DB
}

func NewPassagePostgres(logger logging.Logger, db *sqlx.DB) *PassagePostgres {
	return &PassagePostgres{
		logger: logger,
		db:     db,
	}
}

func (r *PassagePostgres) GetEmployees() ([]models.Employee, error) {
	var employees []models.Employee

	query := fmt.Sprintf("SELECT * FROM %s;", employeeTable)
	if err := r.db.Select(&employees, query); err != nil {
		return nil, err
	}

	return employees, nil
}
