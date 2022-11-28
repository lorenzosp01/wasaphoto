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
	UserExists(int64) DbError
	ChangeUsername(int64, string) DbError
	DeletePhoto(int64) DbError
	IsPhotoOwner(int64, int64) (bool, DbError)
	IsUserBannedBy(int64, int64) (bool, DbError)
	GetUserProfile(int64) (UserProfile, DbError)
	getUserPhotos(int64) ([]Photo, DbError)
	getProfileCounters(int64) (ProfileCounters, DbError)
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

func (e DbError) ToHttp() utils.HttpError {
	switch e.Err {
	case sql.ErrNoRows:
		return utils.HttpError{
			StatusCode: http.StatusNotFound,
			Message:    "Not found",
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

	// Check if table exists. If not, the database is empty, and we need to create the structure
	var tableName string
	err := db.QueryRow(`SELECT name FROM sqlite_master WHERE type='table' AND name='example_table';`).Scan(&tableName)
	if errors.Is(err, sql.ErrNoRows) {
		sqlStmt := `CREATE TABLE example_table (id INTEGER NOT NULL PRIMARY KEY, name TEXT);`
		_, err = db.Exec(sqlStmt)
		if err != nil {
			return nil, fmt.Errorf("error creating database structure: %w", err)
		}
	}

	return &appdbimpl{
		c: db,
	}, nil
}

func (db *appdbimpl) Ping() error {
	return db.c.Ping()
}

// UserExists returns an error if the user does not exist or if there is an error.
// during query execution.
func (db *appdbimpl) UserExists(id int64) DbError {
	query := fmt.Sprintf("SELECT id FROM %s WHERE id=?", UserTable)
	return DbError{db.c.QueryRow(query, id).Scan(&id)}
}

func (db *appdbimpl) IsPhotoOwner(id int64, owner_id int64) (bool, DbError) {
	var realOwner int64
	var dbErr DbError
	query := fmt.Sprintf("SELECT owner FROM %s WHERE id=?", PhotoTable)
	dbErr.Err = db.c.QueryRow(query, id, owner_id).Scan(&realOwner)

	return realOwner == owner_id, dbErr
}

func (db *appdbimpl) IsUserBannedBy(banned_id int64, banning_id int64) (bool, DbError) {
	var dbErr DbError
	query := fmt.Sprintf("SELECT banned FROM %s WHERE banned=? AND banning=?", BanTable)
	dbErr.Err = db.c.QueryRow(query, banned_id, banning_id).Scan(&banned_id)

	return !errors.Is(dbErr.Err, sql.ErrNoRows), dbErr
}
