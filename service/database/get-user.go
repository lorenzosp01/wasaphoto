package database

import (
	"database/sql"
	"errors"
)

type User struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

func (db *appdbimpl) GetUser(username string) (*User, error) {
	var user User
	err := db.c.QueryRow("SELECT id, name FROM User WHERE name=?", username).Scan(&user.Id, &user.Name)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, errors.New("user not found")
		// create user into db
	}

	return &user, nil
}

func (db *appdbimpl) CreateUser(name string) (*User, error) {
	var user User
	err := db.c.QueryRow("INSERT INTO User (name) VALUES (?)", name)
	return &user, nil
}
