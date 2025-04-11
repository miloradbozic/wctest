package db

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"wctest/app/config"
)

type DB struct {
	*sql.DB
}

func InitDB(cfg *config.Config) (*DB, error) {
	db, err := sql.Open("sqlite3", cfg.DBPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %v", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %v", err)
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS employees (
			id INTEGER PRIMARY KEY,
			name TEXT NOT NULL,
			title TEXT NOT NULL,
			manager_id INTEGER,
			FOREIGN KEY (manager_id) REFERENCES employees(id)
		)
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to create employees table: %v", err)
	}

	return &DB{db}, nil
} 