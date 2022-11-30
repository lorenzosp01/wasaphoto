package database

import "fmt"

func (db *appdbimpl) BanUser(authUserId, userId int64) DbError {
	var dbErr DbError
	query := fmt.Sprintf("INSERT INTO %s (banned, banning) VALUES (?, ?)", BanTable)
	_, dbErr.Err = db.c.Exec(query, userId, authUserId)

	return dbErr
}
