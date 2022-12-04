package database

import (
	"database/sql"
	"fmt"
)

func (db *appdbimpl) InsertPhoto(image []byte, ownerId int64) DbError {
	// Upload the photo to the database
	query := fmt.Sprintf("INSERT INTO %s (owner, image) VALUES (?, ?)", PhotoTable)
	_, err := db.c.Exec(query, ownerId, image)
	// If the insert was unsuccessful, return an error
	if err != nil {
		return DbError{
			Err: err,
		}
	}

	return DbError{}
}

func (db *appdbimpl) GetImage(photoId int64) ([]byte, DbError) {
	var image []byte
	query := fmt.Sprintf("SELECT image %s FROM Photo WHERE id=?", PhotoTable)
	err := db.c.QueryRow(query, photoId).Scan(&image)

	if err != nil {
		return nil, DbError{
			Err: err,
		}
	}

	return image, DbError{}
}

func (db *appdbimpl) ChangeUsername(id int64, newUsername string) DbError {
	query := fmt.Sprintf("UPDATE %s SET name=? WHERE ID=?", UserTable)
	_, err := db.c.Exec(query, newUsername, id)
	return DbError{
		Err: err,
	}
}

func (db *appdbimpl) DeletePhoto(id int64) DbError {
	query := fmt.Sprintf("DELETE FROM %s WHERE ID=?", PhotoTable)
	_, err := db.c.Exec(query, id)
	return DbError{
		Err: err,
	}
}

func (db *appdbimpl) GetUserProfile(id int64, photosAmount int64, photosOffset int64) (UserProfile, DbError) {
	var up UserProfile
	var dbErr DbError

	up.UserInfo.Id = id
	query := fmt.Sprintf("SELECT name FROM %s WHERE id=?", UserTable)
	dbErr.Err = db.c.QueryRow(query, id).Scan(&up.UserInfo.Username)

	if dbErr.Err != nil {
		return up, dbErr
	}

	up.Photos, dbErr = db.GetUserPhotos(id, photosAmount, photosOffset)
	if dbErr.Err != nil {
		return up, dbErr
	}

	up.ProfileInfo, dbErr = db.getProfileCounters(id)

	return up, dbErr
}

func (db *appdbimpl) GetUserPhotos(id int64, amount int64, offset int64) ([]Photo, DbError) {
	// Per ogni foto
	var dbErr DbError
	var rows *sql.Rows
	query := fmt.Sprintf("SELECT id, uploaded_at FROM %s WHERE owner=? LIMIT ? OFFSET ?", PhotoTable)
	rows, dbErr.Err = db.c.Query(query, id, amount, offset)

	if dbErr.Err != nil {
		return nil, dbErr
	}

	var photos []Photo
	var photo Photo

	for rows.Next() {
		photo.Owner = id
		dbErr.Err = rows.Scan(&photo.Id, &photo.UploadedAt)
		if dbErr.Err != nil {
			return nil, dbErr
		}

		photo.PhotoInfo, dbErr = db.getPhotoCounters(photo.Id)
		if dbErr.Err != nil {
			return nil, dbErr
		}

		photos = append(photos, photo)
		amount--
	}

	return photos, dbErr
}

func (db *appdbimpl) getPhotoCounters(photoId int64) (PhotoCounters, DbError) {
	var photoCounters PhotoCounters
	var dbErr DbError

	query := fmt.Sprintf("SELECT count(*) FROM %s WHERE photo=?", LikeTable)
	dbErr.Err = db.c.QueryRow(query, photoId).Scan(&photoCounters.LikesCounter)
	if dbErr.Err != nil {
		return photoCounters, dbErr
	}

	query = fmt.Sprintf("SELECT count(*) FROM %s WHERE photo=?", CommentTable)
	dbErr.Err = db.c.QueryRow(query, photoId).Scan(&photoCounters.CommentsCounter)

	return photoCounters, dbErr
}

func (db *appdbimpl) getProfileCounters(id int64) (ProfileCounters, DbError) {
	var dbErr DbError
	var profileCounters ProfileCounters

	query := fmt.Sprintf("SELECT count(*) FROM %s WHERE follower=?", FollowTable)
	dbErr.Err = db.c.QueryRow(query, id).Scan(&profileCounters.FollowingCounter)
	if dbErr.Err != nil {
		return profileCounters, dbErr
	}

	query = fmt.Sprintf("SELECT count(*) FROM %s WHERE following=?", FollowTable)
	dbErr.Err = db.c.QueryRow(query, id).Scan(&profileCounters.FollowersCounter)
	if dbErr.Err != nil {
		return profileCounters, dbErr
	}

	query = fmt.Sprintf("SELECT count(*) FROM %s WHERE owner=?", PhotoTable)
	dbErr.Err = db.c.QueryRow(query, id).Scan(&profileCounters.PhotosCounter)
	return profileCounters, dbErr
}
