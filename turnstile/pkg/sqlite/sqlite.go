package sqlite

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

const (
	schema = `
	CREATE TABLE IF NOT EXISTS employee (
		card_number UNSIGNED BIG INT NOT NULL,
		employee_id UNSIGNED BIG INT NOT NULL,
		rv UNSIGNED BIG INT NOT NULL, 
		isdeleted BOOLEAN NOT NULL CHECK (isdeleted IN (0, 1))
	  );
	CREATE TABLE IF NOT EXISTS log (
		turnstile_id UNSIGNED BIG INT NOT NULL,
		employee_id UNSIGNED BIG INT NOT NULL,
		card UNSIGNED BIG INT NOT NULL, 
		direction INTEGER NOT NULL,
	    dt VARCHAR(50) NOT NULL
	  );
`
)

type Config struct {
	FileName string
}

func New(cfg Config) (*sqlx.DB, error) {
	db, err := sqlx.Open("sqlite3", cfg.FileName)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	db.MustExec(schema)

	return db, nil
}
