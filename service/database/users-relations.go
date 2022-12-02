package database

import (
	"database/sql"
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

// todo gestire qua dentro l'errore di non presenza del ban/follow come conflict
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
	_, dbErr.Err = db.c.Exec(query, authUserId, userId)

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

	var DbErr DbError
	var rows *sql.Rows
	rows, DbErr.Err = db.c.Query(query, authUserId)
	var users []User

	for rows.Next() {
		var user User
		DbErr.Err = rows.Scan(&user.Id)
		if DbErr.Err != nil {
			return users, DbErr
		}
		query = fmt.Sprintf("SELECT name FROM %s WHERE id=?", UserTable)
		DbErr.Err = db.c.QueryRow(query, user.Id).Scan(&user.Username)
		if DbErr.Err != nil {
			return users, DbErr
		}
		users = append(users, user)
	}

	if users == nil {
		DbErr.Err = errors.New("no users found")
	}

	return users, DbErr
}
