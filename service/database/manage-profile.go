package database

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/mattn/go-sqlite3"
)

func (db *appdbimpl) InsertPhoto(image []byte, ownerId int64) DbError {
	// Upload the photo to the database
	query := fmt.Sprintf("INSERT INTO %s (owner, image) VALUES (?, ?)", PhotoTable)
	_, err := db.c.Exec(query, ownerId, image)
	var dbErr DbError
	// If the insert was unsuccessful, return an error
	if err != nil {
		dbErr.InternalError = err
	}

	return dbErr
}

// Photo has to belong to the user in path
func (db *appdbimpl) GetImage(photo int64, user int64) ([]byte, DbError) {
	var image []byte
	query := fmt.Sprintf("SELECT image FROM %s WHERE id=? AND owner=?", PhotoTable)
	err := db.c.QueryRow(query, photo, user).Scan(&image)
	var dbErr DbError

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			dbErr.Code = StateConflict
		}
		dbErr.InternalError = err
	}

	return image, dbErr
}

func (db *appdbimpl) ChangeUsername(id int64, newUsername string) DbError {
	query := fmt.Sprintf("UPDATE %s SET name=? WHERE ID=?", UserTable)
	_, err := db.c.Exec(query, newUsername, id)
	var dbErr DbError

	if err != nil {
		var sqlErr sqlite3.Error
		dbErr.InternalError = err
		if errors.As(err, &sqlErr) {
			if errors.Is(sqlErr.ExtendedCode, sqlite3.ErrConstraintUnique) {
				dbErr.Code = StateConflict
			}
		}
	}

	return dbErr
}

// Photo has to belong to the authenticated user
func (db *appdbimpl) DeletePhoto(photo int64, user int64) (bool, DbError) {
	var dbErr DbError
	var affected int64

	query := fmt.Sprintf("DELETE FROM %s WHERE id=? AND owner=?", PhotoTable)
	res, err := db.c.Exec(query, photo, user)
	if err != nil {
		dbErr.InternalError = err
		return false, dbErr
	} else {
		affected, _ = res.RowsAffected()
	}

	return affected > 0, dbErr
}

func (db *appdbimpl) GetUserProfile(id int64, photosAmount int64, photosOffset int64) (UserProfile, DbError) {
	var up UserProfile
	var dbErr DbError

	up.UserInfo.Id = id
	query := fmt.Sprintf("SELECT name FROM %s WHERE id=?", UserTable)
	err := db.c.QueryRow(query, id).Scan(&up.UserInfo.Username)

	if err != nil {
		dbErr.InternalError = err
		return up, dbErr
	}

	up.Photos, dbErr = db.GetUserPhotos(id, photosAmount, photosOffset)
	if dbErr.InternalError != nil {
		return up, dbErr
	}

	up.ProfileInfo, dbErr = db.getProfileCounters(id)

	return up, dbErr
}

func (db *appdbimpl) GetUserPhotos(id int64, amount int64, offset int64) ([]Photo, DbError) {
	var dbErr DbError
	joinParam := UserTable + ".id"
	userColumn := "name"
	photoColumn := PhotoTable + ".id"

	query := fmt.Sprintf("SELECT %s, %s, owner, uploaded_at FROM %s, %s WHERE owner=%s AND owner=? "+
		"ORDER BY uploaded_at DESC LIMIT ? OFFSET ?", photoColumn, userColumn, PhotoTable, UserTable, joinParam)
	rows, err := db.c.Query(query, id, amount, offset)

	if err != nil {
		dbErr.InternalError = err
		return nil, dbErr
	}

	var photos []Photo
	var photo Photo

	for rows.Next() {
		err = rows.Scan(&photo.Id, &photo.Owner.Username, &photo.Owner.Id, &photo.UploadedAt)
		if err != nil {
			dbErr.InternalError = err
			return nil, dbErr
		}

		photo.PhotoInfo, dbErr = db.getPhotoCounters(photo.Id)
		if dbErr.InternalError != nil {
			return nil, dbErr
		}

		photos = append(photos, photo)
	}

	err = rows.Err()
	if err != nil {
		dbErr.InternalError = err
	}

	defer rows.Close()

	return photos, dbErr
}

func (db *appdbimpl) getPhotoCounters(photoId int64) (PhotoCounters, DbError) {
	var photoCounters PhotoCounters
	var dbErr DbError

	query := fmt.Sprintf("SELECT count(*) FROM %s WHERE photo=?", LikeTable)
	err := db.c.QueryRow(query, photoId).Scan(&photoCounters.LikesCounter)
	if err != nil {
		dbErr.InternalError = err
		return photoCounters, dbErr
	}

	query = fmt.Sprintf("SELECT count(*) FROM %s WHERE photo=?", CommentTable)
	err = db.c.QueryRow(query, photoId).Scan(&photoCounters.CommentsCounter)
	if err != nil {
		dbErr.InternalError = err
	}

	return photoCounters, dbErr
}

func (db *appdbimpl) getProfileCounters(id int64) (ProfileCounters, DbError) {
	var dbErr DbError
	var profileCounters ProfileCounters

	query := fmt.Sprintf("SELECT count(*) FROM %s WHERE follower=?", FollowTable)
	err := db.c.QueryRow(query, id).Scan(&profileCounters.FollowingCounter)
	if err != nil {
		dbErr.InternalError = err
		return profileCounters, dbErr
	}

	query = fmt.Sprintf("SELECT count(*) FROM %s WHERE following=?", FollowTable)
	err = db.c.QueryRow(query, id).Scan(&profileCounters.FollowersCounter)
	if err != nil {
		dbErr.InternalError = err
		return profileCounters, dbErr
	}

	query = fmt.Sprintf("SELECT count(*) FROM %s WHERE owner=?", PhotoTable)
	err = db.c.QueryRow(query, id).Scan(&profileCounters.PhotosCounter)
	if err != nil {
		dbErr.InternalError = err
	}

	return profileCounters, dbErr
}
