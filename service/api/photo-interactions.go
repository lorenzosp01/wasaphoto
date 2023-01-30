package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"wasaphoto/service/database"
	"wasaphoto/service/utils"
)

func (rt *_router) likePhoto(w http.ResponseWriter, r *http.Request, params map[string]int64) {
	authUserId := params["token"]
	userId := params["user_id"]
	photoId := params["photo_id"]

	if authUserId != params["targeted_user_id"] {
		rt.LoggerAndHttpErrorSender(w, errors.New("who puts like and authenticated user id are different"), utils.HttpError{StatusCode: http.StatusForbidden, Message: "You can't like a photo impersonating someone else"})
		return
	}

	isBanned, dbErr := rt.db.IsUserTargeted(userId, authUserId, database.BanTable)
	if dbErr.InternalError != nil {
		if isBanned {
			dbErr.CustomMessage = utils.BannedMessage
		}
		rt.LoggerAndHttpErrorSender(w, dbErr.InternalError, dbErr.ToHttp())
		return
	}

	dbErr = rt.db.LikePhoto(authUserId, photoId, userId)
	if dbErr.InternalError != nil {
		rt.LoggerAndHttpErrorSender(w, dbErr.InternalError, dbErr.ToHttp())
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	_, _ = w.Write([]byte("Photo liked successfully"))
}

func (rt *_router) unlikePhoto(w http.ResponseWriter, r *http.Request, params map[string]int64) {
	authUserId := params["token"]
	userId := params["user_id"]

	if authUserId != params["targeted_user_id"] {
		rt.LoggerAndHttpErrorSender(w, errors.New("who deletes like and authenticated user id are different"), utils.HttpError{StatusCode: http.StatusForbidden, Message: "You can't unlike a photo impersonating someone else"})
		return
	}

	isBanned, dbErr := rt.db.IsUserTargeted(userId, authUserId, database.BanTable)
	if dbErr.InternalError != nil {
		if isBanned {
			dbErr.CustomMessage = utils.BannedMessage
		}
		rt.LoggerAndHttpErrorSender(w, dbErr.InternalError, dbErr.ToHttp())
		return
	}

	photoId := params["photo_id"]

	dbErr = rt.db.UnlikePhoto(authUserId, photoId, userId)
	if dbErr.InternalError != nil {
		rt.LoggerAndHttpErrorSender(w, dbErr.InternalError, dbErr.ToHttp())
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	_, _ = w.Write([]byte("Photo unliked successfully"))
}

func (rt *_router) commentPhoto(w http.ResponseWriter, r *http.Request, params map[string]int64) {
	authUserId := params["token"]
	userId := params["user_id"]
	photoId := params["photo_id"]

	isBanned, dbErr := rt.db.IsUserTargeted(userId, authUserId, database.BanTable)
	if dbErr.InternalError != nil {
		if isBanned {
			dbErr.CustomMessage = utils.BannedMessage
		}
		rt.LoggerAndHttpErrorSender(w, dbErr.InternalError, dbErr.ToHttp())
		return
	}

	var comment Comment
	err := json.NewDecoder(r.Body).Decode(&comment)
	if err != nil {
		rt.LoggerAndHttpErrorSender(w, err, utils.HttpError{StatusCode: http.StatusBadRequest})
		return
	}

	dbErr = rt.db.CommentPhoto(authUserId, photoId, userId, comment.Content)
	if dbErr.InternalError != nil {
		rt.LoggerAndHttpErrorSender(w, dbErr.InternalError, dbErr.ToHttp())
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	_, _ = w.Write([]byte("Comment added succesfully"))
}

func (rt *_router) getPhotoComments(w http.ResponseWriter, r *http.Request, params map[string]int64) {
	authUserId := params["token"]
	userId := params["user_id"]
	photoId := params["photo_id"]

	isBanned, dbErr := rt.db.IsUserTargeted(userId, authUserId, database.BanTable)
	if dbErr.InternalError != nil {
		if isBanned {
			dbErr.CustomMessage = utils.BannedMessage
		}
		rt.LoggerAndHttpErrorSender(w, dbErr.InternalError, dbErr.ToHttp())
		return
	}

	var commentsObject CommentsObject
	dbComments, dbErr := rt.db.GetPhotoComments(photoId, userId)
	if dbErr.InternalError != nil {
		rt.LoggerAndHttpErrorSender(w, dbErr.InternalError, dbErr.ToHttp())
		return
	}

	for _, dbComment := range dbComments {
		var comment Comment
		comment.fromDatabase(dbComment)
		commentsObject.Comments = append(commentsObject.Comments, comment)
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(commentsObject)
}

func (rt *_router) deleteComment(w http.ResponseWriter, r *http.Request, params map[string]int64) {
	commentId := params["comment_id"]
	userId := params["user_id"]
	photoId := params["photo_id"]
	authUserId := params["token"]

	dbErr := rt.db.DeleteComment(photoId, userId, authUserId, commentId)
	if dbErr.InternalError != nil {
		rt.LoggerAndHttpErrorSender(w, dbErr.InternalError, dbErr.ToHttp())
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	_, _ = w.Write([]byte("Comment deleted successfully"))
}
