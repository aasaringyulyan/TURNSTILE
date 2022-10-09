package repository

import (
	"data-generator-mock/models"
	"data-generator-mock/pkg/logging"
	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
)

type DataGeneratorPostgres struct {
	logger logging.Logger
	db     *sqlx.DB
}

func NewDataGeneratorPostgres(logger logging.Logger, db *sqlx.DB) *DataGeneratorPostgres {
	return &DataGeneratorPostgres{
		logger: logger,
		db:     db,
	}
}

func (r *DataGeneratorPostgres) GetEmployees() ([]models.Employee, error) {
	var employees []models.Employee

	query := fmt.Sprintf("SELECT * FROM %s;", employeeTable)
	if err := r.db.Select(&employees, query); err != nil {
		return nil, err
	}

	return employees, nil
}

func (r *DataGeneratorPostgres) AddEmployee(employee models.Employee) error {
	query := fmt.Sprintf("INSERT INTO %s (card_number, employee_id, rv, isdeleted) values ($1, $2, $3, $4);",
		employeeTable)
	_, err := r.db.Exec(query, employee.CardNumber, employee.EmployeeID, employee.Rv, employee.IsDeleted)
	if err != nil {
		return err
	}
	return nil
}

func (r *DataGeneratorPostgres) SaveSlice(data []models.Employee) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	defer func() {
		switch err {
		case nil:
			err = tx.Commit()
		default:
			tx.Rollback()
		}
	}()

	for _, v := range data {
		err = r.save(tx, v)
		if err != nil {
			return err
		}
	}

	return err
}

func (r *DataGeneratorPostgres) save(tx *sql.Tx, employee models.Employee) error {
	query := fmt.Sprintf("INSERT INTO %s (card_number, employee_id, rv, isdeleted) values ($1, $2, $3, $4);",
		employeeTable)

	if _, err := tx.Exec(query, employee.CardNumber, employee.EmployeeID, employee.Rv, employee.IsDeleted); err != nil {
		return err
	}
	return nil
}
