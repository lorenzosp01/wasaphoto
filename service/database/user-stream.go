package database

import (
	"database/sql"
	"errors"
	"fmt"
)

func (db *appdbimpl) GetMyStream(userId int64, offset int64, amount int64) ([]Photo, DbError) {
	var dbErr DbError

	query := fmt.Sprintf("SELECT Photo.id, User.name, owner, uploaded_at FROM %s, %s WHERE owner=User.id AND"+
		" owner IN (SELECT following FROM %s WHERE follower=%d) ORDER BY uploaded_at DESC LIMIT ? OFFSET ?", PhotoTable, UserTable, FollowTable, userId)
	rows, err := db.c.Query(query, userId, amount, offset)

	var photos []Photo

	if err != nil {
		dbErr.Code = GenericError
		dbErr.InternalError = err
	} else {
		for rows.Next() {
			var photo Photo
			err = rows.Scan(&photo.Id, &photo.Owner.Username, &photo.Owner.Id, &photo.UploadedAt)
			if err != nil {
				if errors.Is(err, sql.ErrNoRows) {
					dbErr.Code = NotFound
					dbErr.CustomMessage = "User not found"
				} else {
					dbErr.Code = GenericError
				}
				dbErr.InternalError = err
				return nil, dbErr
			}

			photo.PhotoInfo, dbErr = db.getPhotoCounters(photo.Id)
			if dbErr.InternalError != nil {
				return nil, dbErr
			}

			photos = append(photos, photo)
		}
	}

	defer rows.Close()

	return photos, dbErr
}
