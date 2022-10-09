package repo

import (
	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
	"turnstile/internal/models"
	"turnstile/pkg/logging"
)

type EmployeeRepo struct {
	logger    logging.Logger
	db        *sqlx.DB
	tableName string
}

func NewEmployeeRepo(logger logging.Logger, db *sqlx.DB, tableName string) *EmployeeRepo {
	return &EmployeeRepo{
		logger:    logger,
		db:        db,
		tableName: tableName,
	}
}

func (er *EmployeeRepo) SaveSlice(data []models.Employee) error {
	tx, err := er.db.Begin()
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
		err = er.save(tx, v)
		if err != nil {
			return err
		}
	}

	return err
}

func (er *EmployeeRepo) save(tx *sql.Tx, data models.Employee) error {
	query := fmt.Sprintf("INSERT INTO %s (card_number, employee_id, rv, isdeleted) values ($1, $2, $3, $4)",
		er.tableName)

	if _, err := tx.Exec(query, data.CardNumber, data.EmployeeID, data.Rv, data.IsDeleted); err != nil {
		return err
	}
	return nil
}

func (er *EmployeeRepo) Save(data models.Employee) error {
	tx, err := er.db.Begin()
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

	query := fmt.Sprintf("INSERT INTO %s (card_number, employee_id, rv, isdeleted) values ($1, $2, $3, $4)",
		er.tableName)

	if _, err = tx.Exec(query, data.CardNumber, data.EmployeeID, data.Rv, data.IsDeleted); err != nil {
		return err
	}

	return err
}

func (er *EmployeeRepo) GetEmployeeByCard(card uint64) (models.Employee, error) {
	query := fmt.Sprintf("SELECT * FROM %s WHERE card_number=$1 ", er.tableName)

	var employee models.Employee
	if err := er.db.Get(&employee, query, card); err != nil {
		return models.Employee{}, err
	}

	return employee, nil
}
