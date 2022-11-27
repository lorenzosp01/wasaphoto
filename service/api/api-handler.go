package api

import (
	"net/http"
)

// Handler returns an instance of httprouter.Router that handle APIs registered here
func (rt *_router) Handler() http.Handler {
	// Register routes
	rt.router.POST("/session", rt.doLogin)
	rt.router.POST("/profiles/:user_id/photos", rt.uploadPhoto)
	rt.router.GET("/profiles/:user_id/photos/:photo_id", rt.getImage)
	rt.router.PUT("/profiles/:user_id/name", rt.setMyUsername)
	rt.router.DELETE("/profiles/:user_id/photos/:photo_id", rt.deletePhoto)
	// Special routes
	rt.router.GET("/liveness", rt.liveness)

	return rt.router
}
