package db

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"wctest/app/config"
)

type DB struct {
	*sql.DB
}

func InitDB(cfg *config.Config) (*DB, error) {
	db, err := sql.Open("sqlite3", cfg.DBPath)
	if err != nil {
		return nil, err
	}

	createTableSQL := `
	CREATE TABLE IF NOT EXISTS employees (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		first_name TEXT NOT NULL,
		last_name TEXT NOT NULL,
		title TEXT NOT NULL,
		reports_to_id INTEGER,
		FOREIGN KEY (reports_to_id) REFERENCES employees(id)
	);
	`

	_, err = db.Exec(createTableSQL)
	if err != nil {
		return nil, err
	}

	return &DB{db}, nil
} 