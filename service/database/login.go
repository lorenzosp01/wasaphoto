package database

import (
	"database/sql"
	"errors"
)

// GetUserId returns the user id for the given username.
func (db *appdbimpl) GetUserId(username string) (int64, error) {
	var id int64

	// Get user identifier if a user with the given username exists
	err := db.c.QueryRow("SELECT id FROM User WHERE name=?", username).Scan(&id)

	if errors.Is(err, sql.ErrNoRows) {
		// If no user has been found, create a new one
		id, err = db.createUser(username)
		if err != nil {
			return 0, err
		}
	}

	return id, nil
}

// createUser creates a new user with the given username and returns the user id.
func (db *appdbimpl) createUser(username string) (int64, error) {
	var id int64
	// Insert the new user into the database
	row, err := db.c.Exec("INSERT INTO User (name) VALUES (?)", username)
	// If the insert was unsuccessful, return a no one identifier and an error
	if err != nil {
		return 0, err
	}

	// da modificare in bass alla concorrenza
	id, err = row.LastInsertId()
	return id, err
}
