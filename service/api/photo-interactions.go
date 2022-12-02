package api

import (
	"errors"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strconv"
	"wasaphoto/service/utils"
)

func (rt *_router) likePhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, token utils.Token) {
	authUserId, _ := strconv.ParseInt(token.Value, 10, 64)
	paramAuthUserId, _ := strconv.ParseInt(ps.ByName("auth_user_id"), 10, 64)

	if authUserId != paramAuthUserId {
		rt.LoggerAndHttpErrorSender(w, errors.New("can't put a like impersonating someone else"), utils.HttpError{StatusCode: 403})
		return
	}

	photoId, _ := strconv.ParseInt(ps.ByName("photo_id"), 10, 64)
	//	ownerPhotoId, _ := strconv.ParseInt(ps.ByName("user_id"), 10, 64)
	// se owner ha bloccato authUserId, errore

	dbErr := rt.db.LikePhoto(authUserId, photoId)
	if dbErr.Err != nil {
		rt.LoggerAndHttpErrorSender(w, dbErr.Err, dbErr.ToHttp())
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	_, _ = w.Write([]byte("Photo liked successfully"))
}


func (rt *_router) unlikePhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, token utils.Token){
	authUserId, _ := strconv.ParseInt(token.Value, 10, 64)
	paramAuthUserId, _ := strconv.ParseInt(ps.ByName("auth_user_id"), 10, 64)

	if authUserId != paramAuthUserId {
		rt.LoggerAndHttpErrorSender(w, errors.New("can't delete a like impersonating someone else"), utils.HttpError{StatusCode: 403})
		return
	}

	photoId, _ := strconv.ParseInt(ps.ByName("photo_id"), 10, 64)
	//	ownerPhotoId, _ := strconv.ParseInt(ps.ByName("user_id"), 10, 64)
	// se owner ha bloccato authUserId, errore

	dbErr := rt.db.UnlikePhoto(authUserId, photoId)
	if dbErr.Err != nil {
		rt.LoggerAndHttpErrorSender(w, dbErr.Err, dbErr.ToHttp())
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	_, _ = w.Write([]byte("Photo unliked successfully"))
}

func (rt *_router) commentPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, token utils.Token) {

}

func (rt *_router) getPhotoComments(w http.ResponseWriter, r *http.Request, ps httprouter.Params, token utils.Token) {

}

func (rt *_router) deleteComment (w http.ResponseWriter, r *http.Request, ps httprouter.Params, token utils.Token) {

}
