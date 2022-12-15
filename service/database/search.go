package database

import "fmt"

func (db *appdbimpl) DoSearch(pattern string) ([]User, DbError) {
	var dbErr DbError
	var users []User
	pattern = "%" + pattern + "%"
	query := fmt.Sprintf("SELECT id, name FROM %s WHERE name LIKE ?", UserTable)
	rows, err := db.c.Query(query, pattern)

	defer rows.Close()

	if err != nil {
		dbErr.Code = GenericError
		dbErr.InternalError = err
	} else {
		for rows.Next() {
			var user User
			err = rows.Scan(&user.Id, &user.Username)
			if err != nil {
				dbErr.Code = GenericError
				dbErr.InternalError = err
			}
			users = append(users, user)
		}

		if rows.Err() != nil {
			dbErr.Code = GenericError
			dbErr.InternalError = err
		}
	}

	return users, dbErr
}
