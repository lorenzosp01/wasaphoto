package database

import (
	"errors"
	"fmt"
	"github.com/mattn/go-sqlite3"
)

//func (db *appdbimpl) doesPhotoBelongToUser(photo int64, userId int64) bool {
//	var count int
//	query := fmt.Sprintf("SELECT count(*) FROM %s WHERE id=? AND owner=?", PhotoTable)
//	err := db.c.QueryRow(query, photo, userId).Scan(&count)
//	if err != nil {
//		return false
//	}
//
//	return count > 0
//}

// Ogni funzione restituisce vero se l'operazione è andata a buon fine, falso altrimenti e l'errore generato
// Per i delete devo restituire state conflict se non esiste il record da cancellare
func (db *appdbimpl) LikePhoto(authUser int64, photo int64, photoOwner int64) (bool, DbError) {
	var dbErr DbError
	var affected int64
	query := fmt.Sprintf("INSERT INTO %s (owner, photo) VALUES (?, ?) WHERE EXISTS (SELECT * FROM %s WHERE id=? AND owner=?)", LikeTable, PhotoTable)
	res, err := db.c.Exec(query, authUser, photo, photo, photoOwner)

	if err != nil {
		var sqlErr sqlite3.Error
		if errors.As(err, &sqlErr) {
			if errors.Is(sqlErr.ExtendedCode, sqlite3.ErrConstraintPrimaryKey) {
				dbErr.InternalError = err
				dbErr.Code = StateConflict
			} else {
				dbErr.InternalError = err
				dbErr.Code = GenericError
			}
		}
		return false, dbErr
	} else {
		affected, _ = res.RowsAffected()
	}

	return affected > 0, dbErr
}

func (db *appdbimpl) UnlikePhoto(authUser int64, photo int64, photoOwner int64) (bool, DbError) {
	var dbErr DbError
	var affected int64
	query := fmt.Sprintf("DELETE FROM %s WHERE owner=? AND photo=? WHERE EXISTS (SELECT * FROM %s WHERE id=? AND owner=?  )", LikeTable, PhotoTable)
	res, err := db.c.Exec(query, authUser, photo, photo, photoOwner)

	if err != nil {
		dbErr.InternalError = err
		dbErr.Code = GenericError
		return false, dbErr
	} else {
		affected, _ = res.RowsAffected()
	}

	return affected > 0, dbErr
}

func (db *appdbimpl) CommentPhoto(authUser int64, photo int64, photoOwner int64, commentText string) (bool, DbError) {
	var dbErr DbError
	var affected int64

	query := fmt.Sprintf("INSERT INTO %s (owner, photo, content) VALUES (?, ?, ?) WHERE EXISTS (SELECT * FROM %s WHERE id=? AND owner=?)", CommentTable, PhotoTable)
	res, err := db.c.Exec(query, authUser, photo, commentText, photo, photoOwner)

	if err != nil {
		dbErr.Code = GenericError
		dbErr.InternalError = err
		return false, dbErr
	} else {
		affected, _ = res.RowsAffected()
	}

	return affected > 0, dbErr
}

func (db *appdbimpl) DeleteComment(photo int64, photoOwner int64, commentOwner int64, comment int64) (bool, DbError) {
	var dbErr DbError
	var affected int64

	query := fmt.Sprintf("DELETE FROM %s WHERE id=? AND photo=? AND owner=? WHERE EXISTS (SELECT * FROM %s WHERE id=? AND owner=?)", CommentTable, PhotoTable)
	res, err := db.c.Exec(query, comment, photo, commentOwner, photo, photoOwner)

	if err != nil {
		dbErr.InternalError = err
		dbErr.Code = GenericError
		return false, dbErr
	} else {
		affected, _ = res.RowsAffected()
	}

	return affected > 0, dbErr
}

func (db *appdbimpl) GetPhotoComments(photo int64, photoOwner int64) ([]Comment, DbError) {
	var dbErr DbError
	var comments []Comment

	joinParam := UserTable + ".id"
	userColumn := "name"
	commentColumn := CommentTable + ".id"
	query := fmt.Sprintf("SELECT %s, owner, %s, content, created_at FROM %s, %s WHERE owner=%s AND photo=?"+
		"ORDER BY created_at DESC WHERE EXISTS (SELECT * FROM %s WHERE id=? AND owner=?)", commentColumn, userColumn, CommentTable, UserTable, joinParam, PhotoTable)
	rows, err := db.c.Query(query, photo, photo, photoOwner)

	if err != nil {
		dbErr.Code = GenericError
		dbErr.InternalError = err
	} else {
		for rows.Next() {
			var comment Comment
			err := rows.Scan(&comment.Id, &comment.Owner.Id, &comment.Owner.Username, &comment.Content, &comment.CreatedAt)
			if err != nil {
				dbErr.Code = GenericError
				dbErr.InternalError = err
				break
			}
			comments = append(comments, comment)
		}

		err = rows.Err()
		if err != nil {
			dbErr.Code = GenericError
			dbErr.InternalError = err
		}

		defer rows.Close()
	}

	return comments, dbErr
}
