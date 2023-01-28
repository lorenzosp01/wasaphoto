package api

import (
	"net/http"
)

// Handler returns an instance of httprouter.Router that handle APIs registered here
func (rt *_router) Handler() http.Handler {
	// Register routes
	// Login
	rt.router.POST("/session", rt.doLogin)
	// Manage profile
	rt.router.POST("/profiles/:user_id/photos/", rt.wrap(rt.authWrap(rt.uploadPhoto)))
	rt.router.GET("/profiles/:user_id/photos/:photo_id", rt.wrap(rt.getImage))
	rt.router.PUT("/profiles/:user_id/name", rt.wrap(rt.authWrap(rt.setMyUsername)))
	rt.router.DELETE("/profiles/:user_id/photos/:photo_id", rt.wrap(rt.authWrap(rt.deletePhoto)))
	rt.router.GET("/profiles/:user_id", rt.wrap(rt.getUserProfile))
	// Users relations
	rt.router.PUT("/profiles/:user_id/ban/:targeted_user_id", rt.wrap(rt.authWrap(rt.banUser)))
	rt.router.DELETE("/profiles/:user_id/ban/:targeted_user_id", rt.wrap(rt.authWrap(rt.unbanUser)))
	rt.router.PUT("/profiles/:user_id/following/:targeted_user_id", rt.wrap(rt.authWrap(rt.followUser)))
	rt.router.DELETE("/profiles/:user_id/following/:targeted_user_id", rt.wrap(rt.authWrap(rt.unfollowUser)))
	rt.router.GET("/profiles/:user_id/following/", rt.wrap(rt.authWrap(rt.getFollowedUsers)))
	rt.router.GET("/profiles/:user_id/ban/", rt.wrap(rt.authWrap(rt.getBannedUsers)))
	// Photo interactions
	rt.router.PUT("/profiles/:user_id/photos/:photo_id/likes/:auth_user_id", rt.wrap(rt.likePhoto))
	rt.router.DELETE("/profiles/:user_id/photos/:photo_id/likes/:auth_user_id", rt.wrap(rt.unlikePhoto))
	rt.router.POST("/profiles/:user_id/photos/:photo_id/comments/", rt.wrap(rt.commentPhoto))
	rt.router.DELETE("/profiles/:user_id/photos/:photo_id/comments/:comment_id", rt.wrap(rt.deleteComment))
	rt.router.GET("/profiles/:user_id/photos/:photo_id/comments/", rt.wrap(rt.getPhotoComments))
	rt.router.GET("/search", rt.wrap(rt.doSearch))
	// Stream
	rt.router.GET("/stream/:user_id", rt.wrap(rt.authWrap(rt.getMyStream)))
	// Special routes
	rt.router.GET("/liveness", rt.liveness)

	return rt.router
}
