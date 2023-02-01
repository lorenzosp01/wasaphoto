package database

import (
	"fmt"
)

func (db *appdbimpl) GetMyStream(userId int64, offset int64, amount int64) ([]Photo, DbError) {
	var dbErr DbError

	query := fmt.Sprintf("SELECT Photo.id, User.name, owner, uploaded_at FROM %s, %s WHERE owner=User.id AND"+
		" owner IN (SELECT following FROM %s WHERE follower=?) ORDER BY uploaded_at DESC LIMIT ? OFFSET ?", PhotoTable, UserTable, FollowTable)
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
				dbErr.Code = GenericError
				dbErr.InternalError = err
				return nil, dbErr
			}

			photo.PhotoInfo, dbErr = db.getPhotoCounters(photo.Id)
			if dbErr.InternalError != nil {
				return nil, dbErr
			}

			photos = append(photos, photo)
		}

		if rows.Err() != nil {
			dbErr.Code = GenericError
			dbErr.InternalError = err
			return nil, dbErr
		}
	}

	defer rows.Close()

	return photos, dbErr
}
