package api

import (
	"net/http"
	"wasaphoto/service/database"
)

// Handler returns an instance of httprouter.Router that handle APIs registered here
func (rt *_router) Handler() http.Handler {
	// Register routes
	// Login
	rt.router.POST("/session", rt.doLogin)
	// Manage profile
	rt.router.POST("/profiles/:user_id/photos", rt.wrap(rt.authWrap(rt.uploadPhoto), []string{database.UserTable}, false))
	rt.router.GET("/profiles/:user_id/photos/:photo_id", rt.wrap(rt.getImage, []string{database.UserTable, database.PhotoTable}, true))
	rt.router.PUT("/profiles/:user_id/name", rt.wrap(rt.authWrap(rt.setMyUsername), []string{database.UserTable}, false))
	rt.router.DELETE("/profiles/:user_id/photos/:photo_id", rt.wrap(rt.authWrap(rt.deletePhoto), []string{database.UserTable, database.PhotoTable}, true))
	rt.router.GET("/profiles/:user_id", rt.wrap(rt.getUserProfile, []string{database.UserTable},false))
	// Users relations
	rt.router.PUT("/profiles/:user_id/ban/:targeted_user_id", rt.wrap(rt.authWrap(rt.banUser), []string{database.UserTable, database.UserTable}, false))
	rt.router.DELETE("/profiles/:user_id/ban/:targeted_user_id", rt.wrap(rt.authWrap(rt.unbanUser), []string{database.UserTable, database.UserTable}, false))
	rt.router.PUT("/profiles/:user_id/following/:targeted_user_id", rt.wrap(rt.authWrap(rt.followUser), []string{database.UserTable, database.UserTable}, false))
	rt.router.DELETE("/profiles/:user_id/following/:targeted_user_id", rt.wrap(rt.authWrap(rt.unfollowUser), []string{database.UserTable, database.UserTable}, false))
	rt.router.GET("/profiles/:user_id/following/", rt.wrap(rt.authWrap(rt.getFollowedUsers), []string{database.UserTable}, false))
	rt.router.GET("/profiles/:user_id/ban/", rt.wrap(rt.authWrap(rt.getBannedUsers), []string{database.UserTable}, false))
	// Photo interactions
	rt.router.PUT("/profiles/:user_id/photos/:photo_id/likes", rt.wrap(rt.likePhoto, []string{database.UserTable, database.PhotoTable, database.UserTable}, true))
	rt.router.DELETE("/profiles/:user_id/photos/:photo_id/likes/:auth_user_id", rt.wrap(rt.unlikePhoto, []string{database.UserTable, database.PhotoTable, database.UserTable}, true))
	rt.router.POST("/profiles/:user_id/photos/:photo_id/comments", rt.wrap(rt.commentPhoto, []string{database.UserTable, database.PhotoTable},true))
	rt.router.DELETE("/profiles/:user_id/photos/:photo_id/comments/:comment_id", rt.wrap(rt.deleteComment, []string{database.UserTable, database.PhotoTable, database.CommentTable},true))
	rt.router.GET("/profiles/:user_id/photos/:photo_id/comments", rt.wrap(rt.getPhotoComments, []string{database.UserTable, database.PhotoTable}, true))
	// Stream
	rt.router.GET("/stream/:user_id", rt.wrap(rt.authWrap(rt.getMyStream), []string{database.UserTable}, false))
	// Special routes
	rt.router.GET("/liveness", rt.liveness)

	return rt.router
}
