package repo

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"turnstile/internal/models"
	"turnstile/pkg/logging"
)

type LogRepo struct {
	logger    logging.Logger
	db        *sqlx.DB
	tableName string
}

func NewLogRepo(logger logging.Logger, db *sqlx.DB, tableName string) *LogRepo {
	return &LogRepo{
		logger:    logger,
		db:        db,
		tableName: tableName,
	}
}

func (lr *LogRepo) Save(log models.PassageLogForApi) error {
	tx, err := lr.db.Begin()
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

	query := fmt.Sprintf("INSERT INTO %s (turnstile_id, employee_id, card, direction, dt) "+
		"values ($1, $2, $3, $4, $5)", lr.tableName)

	if _, err = tx.Exec(query, log.TurnstileID, log.EmployeeID, log.CardID, log.Direction, log.DateTime); err != nil {
		return err
	}

	return err
}

func (lr *LogRepo) GetAll() (models.PassageLogsForApi, error) {
	var logs models.PassageLogsForApi

	query := fmt.Sprintf("SELECT * FROM %s", lr.tableName)
	if err := lr.db.Select(&logs.Logs, query); err != nil {
		return models.PassageLogsForApi{}, err
	}

	return logs, nil
}

func (lr *LogRepo) DeleteAll() error {
	query := fmt.Sprintf("DELETE FROM %s", lr.tableName)

	_, err := lr.db.Exec(query)

	return err
}
