package api

import (
	"regexp"
	"wasaphoto/service/database"
	"wasaphoto/service/utils"
)

type CommentsObject struct {
	Comments []Comment `json:"comments"`
}

type PathParams struct {
	AuthUserId   int64
	TargetUserId int64
	PhotoId      int64
	Token        utils.Token
}

type User struct {
	Id       int64  `json:"identifier"`
	Username string `json:"username"`
}

type ProfileCounters struct {
	PhotosCounter    int `json:"photosCounter"`
	FollowingCounter int `json:"followingCounter"`
	FollowersCounter int `json:"followersCounter"`
}

type UserStream struct {
	Photos []Photo `json:"photos"`
}

type Photo struct {
	Id         int64         `json:"id"`
	Owner      int64         `json:"owner"`
	UploadedAt string        `json:"uploadedAt"`
	PhotoInfo  PhotoCounters `json:"photoInfo"`
}

type Comment struct {
	Id        int64  `json:"id"`
	Owner     int64  `json:"owner"`
	Content   string `json:"content"`
	CreatedAt string `json:"uploadedAt"`
}

func (c *Comment) fromDatabase(dbComment database.Comment) {
	c.Id = dbComment.Id
	c.Owner = dbComment.Owner
	c.Content = dbComment.Content
	c.CreatedAt = dbComment.CreatedAt
}

type UserIdentifier struct {
	Id int64 `json:"identifier"`
}

type Username struct {
	Username string `json:"username"`
}

func (u Username) IsValid() bool {
	var valid = false
	if len(u.Username) > 0 && len(u.Username) < 16 {
		valid, _ = regexp.Match(`^[a-zA-Z0-9]+$`, []byte(u.Username))
	}
	return valid
}

type PhotoCounters struct {
	LikesCounter    int `json:"likesCounter"`
	CommentsCounter int `json:"commentsCounter"`
}

type UserProfile struct {
	UserInfo    User            `json:"user_info"`
	Photos      []Photo         `json:"photos"`
	ProfileInfo ProfileCounters `json:"profileInfo"`
}

func (u *User) fromDatabase(dbUser database.User) {
	u.Id = dbUser.Id
	u.Username = dbUser.Username
}

func (up *UserProfile) fromDatabase(upDb database.UserProfile) {
	up.UserInfo.fromDatabase(upDb.UserInfo)
	for _, photo := range upDb.Photos {
		newPhoto := Photo{}
		newPhoto.fromDatabase(photo)
		up.Photos = append(up.Photos, newPhoto)
	}
	up.ProfileInfo.FollowingCounter = upDb.ProfileInfo.FollowingCounter
	up.ProfileInfo.FollowersCounter = upDb.ProfileInfo.FollowersCounter
	up.ProfileInfo.PhotosCounter = upDb.ProfileInfo.PhotosCounter
}

func (p *Photo) fromDatabase(dbPhoto database.Photo) {
	p.Id = dbPhoto.Id
	p.Owner = dbPhoto.Owner
	p.UploadedAt = dbPhoto.UploadedAt
	p.PhotoInfo.LikesCounter = dbPhoto.PhotoInfo.LikesCounter
	p.PhotoInfo.CommentsCounter = dbPhoto.PhotoInfo.CommentsCounter
}
