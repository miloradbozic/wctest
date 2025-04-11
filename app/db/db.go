package db

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func InitDB(isDevelopment bool) error {
	var err error
	DB, err = sql.Open("sqlite3", "./employees.db")
	if err != nil {
		return err
	}

	// Only create table in development mode
	if isDevelopment {
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

		_, err = DB.Exec(createTableSQL)
		if err != nil {
			return err
		}
	}

	return nil
} 