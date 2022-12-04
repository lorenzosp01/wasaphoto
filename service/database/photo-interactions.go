package database

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/mattn/go-sqlite3"
)

func (db *appdbimpl) LikePhoto(authUserId int64, photoId int64) DbError {
	var dbErr DbError
	query := fmt.Sprintf("INSERT INTO %s (owner, photo) VALUES (?, ?)", LikeTable)
	_, dbErr.Err = db.c.Exec(query, authUserId, photoId)

	if dbErr.Err != nil {
		if dbErr.Err.(sqlite3.Error).Code == sqlite3.ErrConstraint {
			return DbError{EntityAlreadyExists}
		}
	}

	return dbErr
}

func (db *appdbimpl) UnlikePhoto(authUserId int64, photoId int64) DbError {
	var res sql.Result
	var dbErr DbError
	query := fmt.Sprintf("DELETE FROM %s WHERE owner=? AND photo=?", LikeTable)
	res, dbErr.Err = db.c.Exec(query, authUserId, photoId)

	if dbErr.Err == nil {
		affected, _ := res.RowsAffected()
		if affected == 0 {
			return DbError{sql.ErrNoRows}
		}
	}

	return dbErr
}

func (db *appdbimpl) CommentPhoto(authUserId int64, photoId int64, commentText string) DbError {
	var dbErr DbError

	query := fmt.Sprintf("INSERT INTO %s (owner, photo, content) VALUES (?, ?, ?)", CommentTable)
	_, dbErr.Err = db.c.Exec(query, authUserId, photoId, commentText)

	return dbErr
}

func (db *appdbimpl) DeleteComment(commentId int64) DbError {
	var dbErr DbError
	var res sql.Result

	query := fmt.Sprintf("DELETE FROM %s WHERE id=?", CommentTable)
	res, dbErr.Err = db.c.Exec(query, commentId)

	if dbErr.Err == nil {
		affected, _ := res.RowsAffected()
		if affected == 0 {
			return DbError{sql.ErrNoRows}
		}
	}

	return dbErr
}

func (db *appdbimpl) GetPhotoComments(photoId int64) ([]Comment, DbError) {
	var dbErr DbError
	var comments []Comment
	var rows *sql.Rows
	query := fmt.Sprintf("SELECT id, owner, created_at, content FROM %s WHERE photo=? ", CommentTable)
	rows, dbErr.Err = db.c.Query(query, photoId)

	if dbErr.Err == nil {
		for rows.Next() {
			var comment Comment
			dbErr.Err = rows.Scan(&comment.Id, &comment.Owner, &comment.CreatedAt, &comment.Content)
			if dbErr.Err != nil {
				return nil, dbErr
			}
			comments = append(comments, comment)
		}
	}

	if comments == nil {
		dbErr.Err = errors.New("no comments found")
	}

	return comments, dbErr
}
