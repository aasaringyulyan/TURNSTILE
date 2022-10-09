package repository

const (
	schema = `
	CREATE TABLE IF NOT EXISTS employee (
		card_number VARCHAR(255) NOT NULL,
		employee_id VARCHAR(255) NOT NULL,
		rv VARCHAR(255) NOT NULL, 
		isdeleted BOOLEAN NOT NULL
	  );`
)
