package database

import "fmt"

func (db *appdbimpl) DoSearch(pattern string) ([]User, DbError) {
	var dbErr DbError
	var users []User
	pattern = "%" + pattern + "%"
	query := fmt.Sprintf("SELECT id, name FROM %s WHERE name LIKE ?", UserTable)
	rows, err := db.c.Query(query, pattern)

	if err != nil {
		dbErr.InternalError = err
	} else {
		for rows.Next() {
			var user User
			err = rows.Scan(&user.Id, &user.Username)
			if err != nil {
				dbErr.InternalError = err
				return nil, dbErr
			}
			users = append(users, user)
		}

		if rows.Err() != nil {
			dbErr.InternalError = err
			return nil, dbErr
		}
	}

	defer rows.Close()

	return users, dbErr
}
