package database

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/mattn/go-sqlite3"
)

func (db *appdbimpl) TargetUser(authUserId int64, userId int64, tableName string) DbError {
	var dbErr DbError
	var query string

	switch tableName {
	case BanTable:
		query = fmt.Sprintf("INSERT INTO %s (banning, banned) VALUES (?, ?)", BanTable)
	case FollowTable:
		query = fmt.Sprintf("INSERT INTO %s (follower, following) VALUES (?, ?)", FollowTable)
	default:
		dbErr.Code = genericError
		return dbErr
	}

	_, err := db.c.Exec(query, authUserId, userId)
	if err != nil {
		if errors.As(err, &sqlite3.Error{}) {
			specificErr, ok := err.(sqlite3.Error)
			if ok {
				if errors.Is(specificErr.ExtendedCode, sqlite3.ErrConstraintPrimaryKey) {
					dbErr.InternalError = err
					dbErr.Code = entityAlreadyExists
					dbErr.CustomMessage = "user already targeted by a " + tableName
				}
			} else {
				dbErr.InternalError = errors.New("error casting error to sqlite3.Error")
				dbErr.Code = genericError
			}
		}
	}

	return dbErr
}

func (db *appdbimpl) UntargetUser(authUserId int64, userId int64, tableName string) DbError {

	var query string
	var dbErr DbError
	switch tableName {
	case BanTable:
		query = fmt.Sprintf("DELETE FROM %s WHERE banning=? AND banned=?", BanTable)
	case FollowTable:
		query = fmt.Sprintf("DELETE FROM %s WHERE follower=? AND following=?", FollowTable)
	default:
		dbErr.Code = genericError
	}

	res, err := db.c.Exec(query, authUserId, userId)
	if err == nil {
		affected, _ := res.RowsAffected()
		if affected == 0 {
			dbErr.Code = notFound
			dbErr.CustomMessage = tableName + " not found"
			dbErr.InternalError = ErrNoRowsDeleted
		}
	} else {
		dbErr.InternalError = err
		dbErr.Code = genericError
	}

	return dbErr
}

func (db *appdbimpl) GetUsersList(authUserId int64, tableName string) ([]User, DbError) {
	var query string
	var dbErr DbError
	dbErr.Code = genericError

	switch tableName {
	case BanTable:
		query = fmt.Sprintf("SELECT banned FROM %s WHERE banning=?", BanTable)
	case FollowTable:
		query = fmt.Sprintf("SELECT following FROM %s WHERE follower=?", FollowTable)
	default:
		return []User{}, dbErr
	}

	rows, err := db.c.Query(query, authUserId)
	if err != nil {
		dbErr.InternalError = err
	}

	var users []User
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			dbErr.Code = genericError
			dbErr.InternalError = err
		}
	}(rows)
	for rows.Next() {
		var user User
		err = rows.Scan(&user.Id)
		if err != nil {
			dbErr.InternalError = err
		}
		query = fmt.Sprintf("SELECT name FROM %s WHERE id=?", UserTable)
		err = db.c.QueryRow(query, user.Id).Scan(&user.Username)
		if err != nil {
			dbErr.InternalError = err
		}
		users = append(users, user)
	}

	if users == nil {
		dbErr.Code = notFound
		dbErr.CustomMessage = "no users targeted by " + tableName
		dbErr.InternalError = errors.New("no users found in db")
	}

	return users, dbErr
}
