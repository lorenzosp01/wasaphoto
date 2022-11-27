package database

import "fmt"

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
