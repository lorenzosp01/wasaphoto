package database

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/mattn/go-sqlite3"
)

func (db *appdbimpl) doesPhotoBelongToUser(photo int64, userId int64) bool {
	var count int
	query := fmt.Sprintf("SELECT count(*) FROM %s WHERE id=? AND owner=?", PhotoTable)
	err := db.c.QueryRow(query, photo, userId).Scan(&count)
	if err != nil {
		return false
	}

	return count > 0
}

func (db *appdbimpl) LikePhoto(authUser int64, photo int64, photoOwner int64) DbError {
	var dbErr DbError

	if db.doesPhotoBelongToUser(photo, photoOwner) {
		query := fmt.Sprintf("INSERT INTO %s (owner, photo) VALUES (?, ?)", LikeTable)
		_, err := db.c.Exec(query, authUser, photo)

		if err != nil {
			if err.(sqlite3.Error).ExtendedCode == sqlite3.ErrConstraintPrimaryKey {
				dbErr.Code = entityAlreadyExists
				dbErr.CustomMessage = "user already put like on this photo"
			}
			dbErr.InternalError = err
		}
	} else {
		dbErr.Code = genericConfilct
		dbErr.CustomMessage = "Photo doesn't belong to that user"
		dbErr.InternalError = errors.New("photo and photo owner don't match")
	}

	return dbErr
}

func (db *appdbimpl) UnlikePhoto(authUser int64, photo int64, photoOwner int64) DbError {
	var dbErr DbError

	if db.doesPhotoBelongToUser(photo, photoOwner) {
		query := fmt.Sprintf("DELETE FROM %s WHERE owner=? AND photo=?", LikeTable)
		res, err := db.c.Exec(query, authUser, photo)

		if err != nil {
			dbErr.InternalError = err
			dbErr.Code = genericError
		} else {
			affected, _ := res.RowsAffected()
			if affected == 0 {
				dbErr.Code = notFound
				dbErr.CustomMessage = "There is no like from that user in that photo"
				dbErr.InternalError = NoRowsDeleted
			}
		}
	} else {
		dbErr.Code = forbiddenAction
		dbErr.CustomMessage = "Photo doesn't belong to that user"
		dbErr.InternalError = errors.New("Photo and photo owner don't match")
	}

	return dbErr
}

func (db *appdbimpl) CommentPhoto(authUser int64, photo int64, photoOwner int64, commentText string) DbError {
	var dbErr DbError

	if db.doesPhotoBelongToUser(photo, photoOwner) {
		query := fmt.Sprintf("INSERT INTO %s (owner, photo, content) VALUES (?, ?, ?)", CommentTable)
		_, err := db.c.Exec(query, authUser, photo, commentText)

		if err != nil {
			dbErr.Code = genericError
			dbErr.InternalError = err
		}
	} else {
		dbErr.Code = genericConfilct
		dbErr.CustomMessage = "Photo doesn't belong to that user"
		dbErr.InternalError = errors.New("photo and photo owner don't match")
	}

	return dbErr
}

func (db *appdbimpl) DeleteComment(photo int64, photoOwner int64, commentOwner int64, comment int64) DbError {
	var dbErr DbError

	if db.doesPhotoBelongToUser(photo, photoOwner) {
		query := fmt.Sprintf("DELETE FROM %s WHERE id=? AND photo=? AND owner=?", CommentTable)
		res, err := db.c.Exec(query, comment, photo, commentOwner)

		if err != nil {
			dbErr.InternalError = err
			dbErr.Code = genericError
		} else {
			affected, _ := res.RowsAffected()
			if affected == 0 {
				dbErr.Code = genericConfilct
				dbErr.CustomMessage = "Comment doesn't belong to that photo or that user isn't the owner of the comment"
				dbErr.InternalError = NoRowsDeleted
			}
		}
	} else {
		dbErr.Code = genericConfilct
		dbErr.CustomMessage = "Photo doesn't belong to that user"
		dbErr.InternalError = errors.New("photo and photo owner don't match")
	}

	return dbErr
}

func (db *appdbimpl) GetPhotoComments(photo int64, photoOwner int64) ([]Comment, DbError) {
	var dbErr DbError
	var comments []Comment

	if db.doesPhotoBelongToUser(photo, photoOwner) {
		query := fmt.Sprintf("SELECT id, owner, content, created_at FROM %s WHERE photo=?", CommentTable)
		rows, err := db.c.Query(query, photo)

		if err != nil {
			dbErr.Code = genericError
			dbErr.InternalError = err
		} else {
			defer func(rows *sql.Rows) {
				err := rows.Close()
				if err != nil {
					dbErr.Code = genericError
					dbErr.InternalError = err

				}
			}(rows)
			for rows.Next() {
				var comment Comment
				err := rows.Scan(&comment.Id, &comment.Owner, &comment.Content, &comment.CreatedAt)
				if err != nil {
					dbErr.Code = genericError
					dbErr.InternalError = err
					break
				}
				comments = append(comments, comment)
			}
		}
	} else {
		dbErr.Code = genericConfilct
		dbErr.CustomMessage = "Photo doesn't belong to that user"
		dbErr.InternalError = errors.New("photo and photo owner don't match")
	}

	return comments, dbErr
}
