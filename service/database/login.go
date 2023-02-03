package database

import (
	"database/sql"
	"errors"
	"fmt"
)

// GetUserId returns the user id for the given username.
func (db *appdbimpl) GetUserId(username string) (bool, int64, DbError) {
	query := fmt.Sprintf("SELECT id FROM %s WHERE name=?", UserTable)
	var id int64
	// Get user identifier if a user with the given username exists
	err := db.c.QueryRow(query, username).Scan(&id)
	var dbErr DbError
	created := false
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// If no user has been found, create a new one
			id, err = db.createUser(username)
			if id != -1 {
				created = true
			}
		}
		dbErr.InternalError = err

		return created, id, dbErr
	}

	return created, id, dbErr
}

// createUser creates a new user with the given username and returns the user id.
func (db *appdbimpl) createUser(username string) (int64, error) {
	var id int64
	// Insert the new user into the database
	row, err := db.c.Exec("INSERT INTO User (name) VALUES (?)", username)
	// If the insert was unsuccessful, return a no one identifier and an error
	if err != nil {
		return -1, err
	}

	id, err = row.LastInsertId()
	return id, err
}
