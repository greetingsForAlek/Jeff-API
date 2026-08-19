package main

import (
	"database/sql"

	_ "modernc.org/sqlite"
)

var db *sql.DB // The databse

// the function that initialises the db
func initDB() error { 
	var err error // An error that will be thrown

	db, err = sql.Open("sqlite", "characters.db") // Open the Datebase

	if err != nil { // Handle any errors
		return err
	}

	// Execute some SQL that basically initialises the Database if it is not already initialised.
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS characters (
			id INTEGER PRIMARY KEY,
			name TEXT NOT NULL,
			description TEXT NOT NULL,
			alignment TEXT NOT NULL,
			image TEXT NOT NULL,
			canon INTEGER NOT NULL CHECK (canon IN (0, 1))
		)
	`)

	return err // Return any errors
}