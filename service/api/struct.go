package api

import "wasaphoto/service/database"

type User struct {
	Id       int64  `json:"identifier"`
	Username string `json:"username"`
}

type UserProfile struct {
	UserInfo    User            `json:"user_info"`
	Photos      []Photo         `json:"photos"`
	ProfileInfo ProfileCounters `json:"profileInfo"`
}

type ProfileCounters struct {
	PhotosCounter    int `json:"photosCounter"`
	FollowingCounter int `json:"followingCounter"`
	FollowersCounter int `json:"followersCounter"`
}

type Photo struct {
	Id         int64         `json:"id"`
	Owner      int64         `json:"owner"`
	UploadedAt string        `json:"uploadedAt"`
	PhotoInfo  PhotoCounters `json:"photoInfo"`
}

type PhotoCounters struct {
	LikesCounter    int `json:"likesCounter"`
	CommentsCounter int `json:"commentsCounter"`
}

func (up *UserProfile) fromDatabase(upDb database.UserProfile) {
	up.UserInfo.Id = upDb.UserInfo.Id
	up.UserInfo.Username = upDb.UserInfo.Username
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
