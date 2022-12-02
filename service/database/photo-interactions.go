package database

import (
	"database/sql"
	"fmt"
	"github.com/mattn/go-sqlite3"
)

func (db *appdbimpl) LikePhoto(authUserId int64, photoId int64) DbError {
	var dbErr DbError

	query := fmt.Sprintf("INSERT INTO %s (owner, photo) VALUES (?, ?)", LikeTable)
	_, dbErr.Err = db.c.Exec(query, authUserId, photoId)

	if dbErr.Err != nil {
		if dbErr.Err.(sqlite3.Error).Code == sqlite3.ErrConstraint {
			return DbError{EntityAlreadyExists}
		}
	}

	return dbErr
}

func (db *appdbimpl) UnlikePhoto(authUserId int64, photoId int64) DbError {
	var dbErr DbError
	var res sql.Result

	query := fmt.Sprintf("DELETE FROM %s WHERE owner=? AND photo=?", LikeTable)
	res, dbErr.Err = db.c.Exec(query, authUserId, photoId)

	if dbErr.Err == nil {
		affected, _ := res.RowsAffected()
		if affected == 0 {
			return DbError{sql.ErrNoRows}
		}
	}

	return dbErr
}
