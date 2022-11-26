package database

import (
	"database/sql"
	"errors"
)

func (db *appdbimpl) InsertPhoto(image []byte, owner_id int64) error {

	// Upload the photo to the database
	_, err := db.c.Exec("INSERT INTO Photo (owner, image) VALUES (?, ?)", owner_id, image)
	// If the insert was unsuccessful, return an error
	if err != nil {
		return err
	}

	return nil
}

func (db *appdbimpl) GetImage(photoId int64) ([]byte, error) {

	var image []byte
	err := db.c.QueryRow("SELECT image FROM Photo WHERE id=?", photoId).Scan(&image)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}

	return image, nil
}
