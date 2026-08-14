package main

import (
	"database/sql"

	_ "modernc.org/sqlite"
)

var db *sql.DB

func initDB() error {
	var err error

	db, err = sql.Open("sqlite", "characters.db")

	if err != nil {
		return err
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS characters (
			id INTEGER PRIMARY KEY,
			name TEXT NOT NULL,
			description TEXT NOT NULL,
			alignment TEXT NOT NULL
		)
	`)

	return err
}