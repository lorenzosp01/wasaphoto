/*
Package database is the middleware between the app database and the code. All data (de)serialization (save/load) from a
persistent database are handled here. Database specific logic should never escape this package.

To use this package you need to apply migrations to the database if needed/wanted, connect to it (using the database
data source name from config), and then initialize an instance of AppDatabase from the DB connection.

For example, this code adds a parameter in `webapi` executable for the database data source name (add it to the
main.WebAPIConfiguration structure):

	DB struct {
		Filename string `conf:""`
	}

This is an example on how to migrate the DB and connect to it:

	// Start Database
	logger.Println("initializing database support")
	db, err := sql.Open("sqlite3", "./foo.db")
	if err != nil {
		logger.WithError(err).Error("error opening SQLite DB")
		return fmt.Errorf("opening SQLite: %w", err)
	}
	defer func() {
		logger.Debug("database stopping")
		_ = db.Close()
	}()

Then you can initialize the AppDatabase and pass it to the api package.
*/
package database

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"wasaphoto/service/utils"
)

// AppDatabase is the high level interface for the DB
type AppDatabase interface {
	Ping() error
	GetUserId(string) (int64, DbError)
	GetImage(int64) ([]byte, DbError)
	InsertPhoto([]byte, int64) DbError
	EntityExists(int64, string) DbError
	ChangeUsername(int64, string) DbError
	DeletePhoto(int64) DbError
	GetUserProfile(int64, int64, int64) (UserProfile, DbError)
	GetUserPhotos(int64, int64, int64) ([]Photo, DbError)
	getProfileCounters(int64) (ProfileCounters, DbError)
	TargetUser(int64, int64, string) DbError
	IsUserAlreadyTargeted(int64, int64, string) (bool, DbError)
	UntargetUser(int64, int64, string) DbError
	GetUsersList(int64, string) ([]User, DbError)
	LikePhoto(int64, int64) DbError
	UnlikePhoto(int64, int64) DbError
	CommentPhoto(int64, int64, string) DbError
	GetPhotoComments(int64) ([]Comment, DbError)
	DeleteComment(int64) DbError
	DoesEntityBelongsTo(int64, int64, string) (bool, DbError)
}

type UserProfile struct {
	UserInfo    User
	Photos      []Photo
	ProfileInfo ProfileCounters
}

type ProfileCounters struct {
	PhotosCounter    int
	FollowingCounter int
	FollowersCounter int
}

type Comment struct {
	Id        int64
	Owner     int64
	Content   string
	CreatedAt string
}

type Photo struct {
	Id         int64
	Owner      int64
	UploadedAt string
	PhotoInfo  PhotoCounters
}

type PhotoCounters struct {
	LikesCounter    int
	CommentsCounter int
}

type User struct {
	Id       int64
	Username string
}

type DbError struct {
	Err error
}

var EntityAlreadyExists = errors.New("entity already exists")
var ForbiddenAction = errors.New("forbidden action")
var BadInput = errors.New("bad input")

func (e DbError) ToHttp() utils.HttpError {
	switch e.Err {
	case sql.ErrNoRows:
		return utils.HttpError{
			StatusCode: http.StatusNotFound,
			Message:    "Not found",
		}
	case nil:
		return utils.HttpError{}
	case BadInput:
		return utils.HttpError{
			StatusCode: http.StatusBadRequest,
			Message:    "Bad parameters",
		}
	case EntityAlreadyExists:
		return utils.HttpError{
			StatusCode: http.StatusConflict,
			Message:    "Conflict with the server state",
		}
	case ForbiddenAction:
		return utils.HttpError{
			StatusCode: http.StatusForbidden,
			Message:    "Forbidden action",
		}
	default:
		return utils.HttpError{
			StatusCode: http.StatusInternalServerError,
			Message:    "Internal server error",
		}
	}
}

const (
	UserTable    string = "User"
	PhotoTable   string = "Photo"
	BanTable     string = "Ban"
	LikeTable    string = "Like"
	FollowTable  string = "Follow"
	CommentTable string = "Comment"
)

type appdbimpl struct {
	c *sql.DB
}

// New returns a new instance of AppDatabase based on the SQLite connection `db`.
// `db` is required - an error will be returned if `db` is `nil`.
func New(db *sql.DB) (AppDatabase, error) {
	if db == nil {
		return nil, errors.New("database is required when building a AppDatabase")
	}

	_, _ = db.Exec("PRAGMA foreign_keys = ON")

	return &appdbimpl{
		c: db,
	}, nil
}

func (db *appdbimpl) Ping() error {
	return db.c.Ping()
}

func (db *appdbimpl) EntityExists(id int64, tableToUse string) DbError {
	query := fmt.Sprintf("SELECT id FROM %s WHERE id=?", tableToUse)
	return DbError{db.c.QueryRow(query, id).Scan(&id)}
}

func (db *appdbimpl) DoesEntityBelongsTo(entityId int64, ownerId int64, entityTable string) (bool, DbError) {
	var dbErr DbError
	var count int64
	var query string

	switch entityTable {
	case CommentTable:
		query = fmt.Sprintf("SELECT count(*) FROM %s WHERE id=? AND photo=?", entityTable)
	default:
		query = fmt.Sprintf("SELECT count(*) FROM %s WHERE id=? AND owner=?", entityTable)
	}

	dbErr.Err = db.c.QueryRow(query, entityId, ownerId).Scan(&count)

	return count > 0, dbErr
}

func (db *appdbimpl) IsUserAlreadyTargeted(targetingUserId int64, targetedUserId int64, tableName string) (bool, DbError) {
	var dbErr DbError
	var query string
	switch tableName {
	case BanTable:
		query = fmt.Sprintf("SELECT count(*) FROM %s WHERE banned=? AND banning=?", BanTable)
	case FollowTable:
		query = fmt.Sprintf("SELECT count(*) FROM %s WHERE following=? AND follower=?", FollowTable)
	default:
		return false, DbError{errors.New("invalid table name")}
	}

	var targetCount int
	dbErr.Err = db.c.QueryRow(query, targetedUserId, targetingUserId).Scan(&targetCount)

	return targetCount > 0, dbErr
}
