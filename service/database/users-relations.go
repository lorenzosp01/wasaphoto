package database

import (
	"errors"
	"fmt"
)

func (db *appdbimpl) TargetUser(authUserId int64, userId int64, tableName string) DbError {

	var query string
	switch tableName {
	case BanTable:
		query = fmt.Sprintf("INSERT INTO %s (banning, banned) VALUES (?, ?)", BanTable)
	case FollowTable:
		query = fmt.Sprintf("INSERT INTO %s (follower, following) VALUES (?, ?)", FollowTable)
	default:
		return DbError{errors.New("invalid table name")}
	}

	var dbErr DbError
	_, dbErr.Err = db.c.Exec(query, authUserId, userId)

	return dbErr
}
