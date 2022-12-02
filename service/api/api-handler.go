package api

import (
	"net/http"
	"wasaphoto/service/database"
)

// Handler returns an instance of httprouter.Router that handle APIs registered here
func (rt *_router) Handler() http.Handler {
	// Register routes
	rt.router.POST("/session", rt.doLogin)
	rt.router.POST("/profiles/:user_id/photos", rt.wrap(rt.authWrap(rt.uploadPhoto), []string{database.UserTable}))
	rt.router.GET("/profiles/:user_id/photos/:photo_id", rt.wrap(rt.getImage, []string{database.UserTable, database.PhotoTable}))
	rt.router.PUT("/profiles/:user_id/name", rt.wrap(rt.authWrap(rt.setMyUsername), []string{database.UserTable}))
	rt.router.DELETE("/profiles/:user_id/photos/:photo_id", rt.wrap(rt.authWrap(rt.deletePhoto), []string{database.UserTable, database.PhotoTable}))
	rt.router.GET("/profiles/:user_id", rt.wrap(rt.getUserProfile, []string{database.UserTable}))
	rt.router.PUT("/profiles/:user_id/ban/:targeted_user_id", rt.wrap(rt.authWrap(rt.banUser), []string{database.UserTable, database.UserTable}))
	rt.router.DELETE("/profiles/:user_id/ban/:targeted_user_id", rt.wrap(rt.authWrap(rt.unbanUser), []string{database.UserTable, database.UserTable}))
	rt.router.PUT("/profiles/:user_id/following/:targeted_user_id", rt.wrap(rt.authWrap(rt.followUser), []string{database.UserTable, database.UserTable}))
	rt.router.DELETE("/profiles/:user_id/following/:targeted_user_id", rt.wrap(rt.authWrap(rt.unfollowUser), []string{database.UserTable, database.UserTable}))
	rt.router.GET("/profiles/:user_id/following/", rt.wrap(rt.authWrap(rt.getFollowedUsers), []string{database.UserTable}))
	rt.router.GET("/profiles/:user_id/ban/", rt.wrap(rt.authWrap(rt.getBannedUsers), []string{database.UserTable}))
	rt.router.PUT("/profiles/:user_id/photos/:photo_id/likes/:auth_user_id", rt.wrap(rt.likePhoto, []string{database.UserTable, database.PhotoTable, database.UserTable}))
	rt.router.DELETE("/profiles/:user_id/photos/:photo_id/likes/:auth_user_id", rt.wrap(rt.unlikePhoto, []string{database.UserTable, database.PhotoTable, database.UserTable}))
	rt.router.POST("/profiles/:user_id/photos/:photo_id/comments", rt.wrap(rt.commentPhoto, []string{database.UserTable, database.PhotoTable}))
	rt.router.DELETE("/profiles/:user_id/photos/:photo_id/comments/:comment_id", rt.wrap(rt.deleteComment, []string{database.UserTable, database.PhotoTable, database.CommentTable}))
	rt.router.GET("/profiles/:user_id/photos/:photo_id/comments", rt.wrap(rt.getPhotoComments, []string{database.UserTable, database.PhotoTable}))
	// Special routes
	rt.router.GET("/liveness", rt.liveness)

	return rt.router
}
