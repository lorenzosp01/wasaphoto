package api

import (
	"net/http"
)

// todo add wrapper
// Handler returns an instance of httprouter.Router that handle APIs registered here
func (rt *_router) Handler() http.Handler {
	// Register routes
	rt.router.POST("/session", rt.doLogin)
	rt.router.POST("/profiles/:user_id/photos", rt.uploadPhoto)
	rt.router.GET("/profiles/:user_id/photos/:photo_id", rt.getImage)
	rt.router.PUT("/profiles/:user_id/name", rt.setMyUsername)
	rt.router.DELETE("/profiles/:user_id/photos/:photo_id", rt.deletePhoto)
	rt.router.GET("/profiles/:user_id", rt.getUserProfile)
	rt.router.PUT("/profiles/:user_id/ban/:targeted_user_id", rt.banUser)
	//rt.router.DELETE("/profiles/:auth_user_id/ban/:user_id", rt.unbanUser)
	rt.router.PUT("/profiles/:user_id/following/:targeted_user_id", rt.followUser)
	//rt.router.DELETE("/profiles/:auth_user_id/following/:user_id", rt.unfollowUser)
	// Special routes
	rt.router.GET("/liveness", rt.liveness)

	return rt.router
}
