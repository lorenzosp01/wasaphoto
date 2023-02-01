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
	GetImage(int64, int64) ([]byte, DbError)
	InsertPhoto([]byte, int64) DbError
	EntityExists(int64, string) (bool, DbError)
	ChangeUsername(int64, string) DbError
	DeletePhoto(int64, int64) (bool, DbError)
	GetUserProfile(int64, int64, int64) (UserProfile, DbError)
	GetMyStream(int64, int64, int64) ([]Photo, DbError)
	GetUserPhotos(int64, int64, int64) ([]Photo, DbError)
	getProfileCounters(int64) (ProfileCounters, DbError)
	TargetUser(int64, int64, string) (bool, DbError)
	IsUserTargeted(int64, int64, string) (bool, DbError)
	UntargetUser(int64, int64, string) (bool, DbError)
	GetUsersList(int64, string) ([]User, DbError)
	LikePhoto(int64, int64, int64) (bool, DbError)
	UnlikePhoto(int64, int64, int64) (bool, DbError)
	CommentPhoto(int64, int64, int64, string) (bool, DbError)
	GetPhotoComments(int64, int64) ([]Comment, DbError)
	DeleteComment(int64, int64, int64, int64) (bool, DbError)
	DoSearch(string) ([]User, DbError)
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
	Owner     User
	Content   string
	CreatedAt string
}

type Photo struct {
	Id         int64
	Owner      User
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
	InternalError error
	Code          int
}

// todo rivedere (i codici dovrebbe restituirli solo l'api)
const (
	StateConflict int = 4
	GenericError  int = 5
)

func (e DbError) ToHttp() utils.HttpError {
	var httpErr utils.HttpError

	switch e.Code {
	//case ForbiddenAction:
	//	httpErr.StatusCode = http.StatusForbidden
	//case BadInput:
	//	httpErr.StatusCode = http.StatusBadRequest
	case StateConflict:
		httpErr.StatusCode = http.StatusConflict
	case GenericError:
		httpErr.StatusCode = http.StatusInternalServerError
	}
	return httpErr
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

var ParamsNameToTable = map[string]string{
	"token":            UserTable,
	"user_id":          UserTable,
	"photo_id":         PhotoTable,
	"targeted_user_id": UserTable,
	"comment_id":       CommentTable,
}

// New returns a new instance of AppDatabase based on the SQLite connection `db`.
// `db` is required - an error will be returned if `db` is `nil`.
func New(db *sql.DB) (AppDatabase, error) {
	if db == nil {
		return nil, errors.New("database is required when building a AppDatabase")
	}

	_, _ = db.Exec("PRAGMA foreign_keys = ON")

	var tableName string
	err := db.QueryRow(`SELECT name FROM sqlite_master WHERE type='table' AND name='User';`).Scan(&tableName)
	if errors.Is(err, sql.ErrNoRows) {

		_, err := db.Exec(
			` create table User
					(
						id   integer primary key autoincrement,
						name text not null unique
					);

				create table Ban
				(
					banned  integer not null references User on delete cascade,
					banning integer not null references User on delete cascade,
					primary key (banned, banning)
				);

				create table Follow
				(
					follower  integer not null references User on delete cascade,
					following integer not null references User on delete cascade,
					primary key (follower, following)
				);

				create table Photo
				(
					id          integer
					primary key autoincrement,
					owner       integer not null
					references User
					on delete cascade,
					image       blob    not null,
					uploaded_at datetime default current_timestamp
				);

				create table Comment
				(
					id         integer
					primary key autoincrement,
					owner      integer                            not null
					references User
					on delete cascade,
					content    text                               not null,
					created_at datetime default current_timestamp not null,
					photo      integer                            not null
					references Photo
					on delete cascade
				);

				create table Like
				(
					owner integer not null
					references User
					on delete cascade,
					photo integer not null
					references Photo
					on delete cascade,
					primary key (owner, photo)
				);
`)

		if err != nil {
			return nil, err
		}
	}

	return &appdbimpl{
		c: db,
	}, nil
}

func (db *appdbimpl) Ping() error {
	return db.c.Ping()
}

func (db *appdbimpl) EntityExists(id int64, tableToUse string) (bool, DbError) {
	var entityCounter int
	query := fmt.Sprintf("SELECT count(*) FROM %s WHERE id=?", tableToUse)
	err := db.c.QueryRow(query, id).Scan(&entityCounter)
	var dbErr DbError
	if err != nil {
		dbErr.InternalError = err
		dbErr.Code = GenericError
		return false, dbErr
	}

	return entityCounter > 0, dbErr
}

func (db *appdbimpl) IsUserTargeted(targetingUserId int64, targetedUserId int64, tableName string) (bool, DbError) {
	var dbErr DbError
	var query string

	switch tableName {
	case BanTable:
		query = fmt.Sprintf("SELECT count(*) FROM %s WHERE banned=? AND banning=?", BanTable)
	case FollowTable:
		query = fmt.Sprintf("SELECT count(*) FROM %s WHERE following=? AND follower=?", FollowTable)
	default:
		return false, dbErr
	}

	var targetCount int
	err := db.c.QueryRow(query, targetedUserId, targetingUserId).Scan(&targetCount)

	if err != nil {
		dbErr.InternalError = err
		dbErr.Code = GenericError
		return false, dbErr
	}

	//if targetCount > 0 {
	//	dbErr.InternalError = errors.New("User is targeted")
	//	dbErr.Code = ForbiddenAction
	//}

	return targetCount > 0, dbErr
}
