package database

import (
	"fmt"
	"github.com/mattn/go-sqlite3"
)

func (db *appdbimpl) LikePhoto(authUserId int64, photoId int64) DbError {
	var dbErr DbError
	query := fmt.Sprintf("INSERT INTO %s (owner, photo) VALUES (?, ?)", LikeTable)
	_, err := db.c.Exec(query, authUserId, photoId)

	dbErr.InternalError = err
	dbErr.Code = genericError
	if err != nil {
		if err.(sqlite3.Error).Code == sqlite3.ErrConstraint {
			dbErr.Code = entityAlreadyExists
			dbErr.CustomMessage = "user already put like on this photo"
		}
	}

	return dbErr
}

func (db *appdbimpl) UnlikePhoto(authUserId int64, photoId int64) DbError {
	var dbErr DbError
	query := fmt.Sprintf("DELETE FROM %s WHERE owner=? AND photo=?", LikeTable)
	res, err := db.c.Exec(query, authUserId, photoId)

	if err != nil {
		dbErr.InternalError = err
		dbErr.Code = genericError
	} else {
		affected, _ := res.RowsAffected()
		if affected == 0 {
			dbErr.Code = notFound
			dbErr.CustomMessage = "Like not found"
			dbErr.InternalError = NoRowsDeleted
		}
	}

	return dbErr
}

func (db *appdbimpl) CommentPhoto(authUserId int64, photoId int64, commentText string) DbError {
	var dbErr DbError

	query := fmt.Sprintf("INSERT INTO %s (owner, photo, content) VALUES (?, ?, ?)", CommentTable)
	_, err := db.c.Exec(query, authUserId, photoId, commentText)

	if err != nil {
		dbErr.InternalError = err
		dbErr.Code = genericError
	}

	return dbErr
}

func (db *appdbimpl) DeleteComment(commentId int64) DbError {
	var dbErr DbError

	query := fmt.Sprintf("DELETE FROM %s WHERE id=?", CommentTable)
	res, err := db.c.Exec(query, commentId)

	if err != nil {
		dbErr.InternalError = err
		dbErr.Code = genericError
	} else {
		affected, _ := res.RowsAffected()
		if affected == 0 {
			dbErr.Code = notFound
			dbErr.CustomMessage = "Comment not found"
			dbErr.InternalError = NoRowsDeleted
		}
	}

	return dbErr
}

func (db *appdbimpl) GetPhotoComments(photoId int64) ([]Comment, DbError) {
	var dbErr DbError
	var comments []Comment
	query := fmt.Sprintf("SELECT id, owner, created_at, content FROM %s WHERE photo=? ", CommentTable)
	rows, err := db.c.Query(query, photoId)

	if err != nil {
		dbErr.Code = genericError
		dbErr.InternalError = err
		return nil, dbErr
	}

	for rows.Next() {
		var comment Comment
		err = rows.Scan(&comment.Id, &comment.Owner, &comment.CreatedAt, &comment.Content)
		if err != nil {
			dbErr.Code = genericError
			dbErr.InternalError = err
			return nil, dbErr
		}
		comments = append(comments, comment)
	}

	if comments == nil {
		dbErr.Code = 1
		dbErr.CustomMessage = "Comments not found"
	}

	return comments, dbErr
}
