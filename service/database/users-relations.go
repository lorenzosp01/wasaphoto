package database

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/mattn/go-sqlite3"
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
	if dbErr.Err != nil {
		if dbErr.Err.(sqlite3.Error).Code == sqlite3.ErrConstraint {
			return DbError{EntityAlreadyExists}
		}
	}

	return dbErr
}

func (db *appdbimpl) UntargetUser(authUserId int64, userId int64, tableName string) DbError {

	var query string
	switch tableName {
	case BanTable:
		query = fmt.Sprintf("DELETE FROM %s WHERE banning=? AND banned=?", BanTable)
	case FollowTable:
		query = fmt.Sprintf("DELETE FROM %s WHERE follower=? AND following=?", FollowTable)
	default:
		return DbError{errors.New("invalid table name")}
	}

	var dbErr DbError
	var res sql.Result
	res, dbErr.Err = db.c.Exec(query, authUserId, userId)
	if dbErr.Err == nil {
		affected, _ := res.RowsAffected()
		if affected == 0 {
			return DbError{sql.ErrNoRows}
		}
	}
	return dbErr
}

func (db *appdbimpl) GetUsersList(authUserId int64, tableName string) ([]User, DbError) {
	var query string
	switch tableName {
	case BanTable:
		query = fmt.Sprintf("SELECT banned FROM %s WHERE banning=?", BanTable)
	case FollowTable:
		query = fmt.Sprintf("SELECT following FROM %s WHERE follower=?", FollowTable)
	default:
		return []User{}, DbError{errors.New("invalid table name")}
	}

	var dbErr DbError
	var rows *sql.Rows
	rows, dbErr.Err = db.c.Query(query, authUserId)
	var users []User

	for rows.Next() {
		var user User
		dbErr.Err = rows.Scan(&user.Id)
		if dbErr.Err != nil {
			return users, dbErr
		}
		query = fmt.Sprintf("SELECT name FROM %s WHERE id=?", UserTable)
		dbErr.Err = db.c.QueryRow(query, user.Id).Scan(&user.Username)
		if dbErr.Err != nil {
			return users, dbErr
		}
		users = append(users, user)
	}

	if users == nil {
		dbErr.Err = errors.New("no users found")
	}

	return users, dbErr
}
