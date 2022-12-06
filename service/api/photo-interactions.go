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
	photoId, _ := params["photo_id"]

	isBanned, dbErr := rt.db.IsUserAlreadyTargeted(userId, authUserId, database.BanTable)
	if isBanned {
		rt.LoggerAndHttpErrorSender(w, errors.New("can't like a photo of a user that banned you"), utils.HttpError{StatusCode: 403})
		return
	} else {
		if dbErr.Err != nil {
			rt.LoggerAndHttpErrorSender(w, dbErr.Err, dbErr.ToHttp())
			return
		}
	}

	dbErr = rt.db.LikePhoto(authUserId, photoId)
	if dbErr.Err != nil {
		rt.LoggerAndHttpErrorSender(w, dbErr.Err, dbErr.ToHttp())
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	_, _ = w.Write([]byte("Photo liked successfully"))
}

func (rt *_router) unlikePhoto(w http.ResponseWriter, r *http.Request, params map[string]int64) {
	authUserId := params["token"]
	userId := params["user_id"]

	if authUserId != userId {
		rt.LoggerAndHttpErrorSender(w, errors.New("can't delete a like impersonating someone else"), utils.HttpError{StatusCode: 403})
		return
	}

	photoId, _ := params["photo_id"]

	dbErr := rt.db.UnlikePhoto(authUserId, photoId)
	if dbErr.Err != nil {
		rt.LoggerAndHttpErrorSender(w, dbErr.Err, dbErr.ToHttp())
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	_, _ = w.Write([]byte("Photo unliked successfully"))
}

func (rt *_router) commentPhoto(w http.ResponseWriter, r *http.Request, params map[string]int64) {
	authUserId := params["token"]
	userId := params["user_id"]
	photoId, _ := params["photo_id"]

	isBanned, dbErr := rt.db.IsUserAlreadyTargeted(userId, authUserId, database.BanTable)
	if isBanned {
		rt.LoggerAndHttpErrorSender(w, errors.New("can't comment a photo of a user that banned you"), utils.HttpError{StatusCode: 403})
		return
	} else {
		if dbErr.Err != nil {
			rt.LoggerAndHttpErrorSender(w, dbErr.Err, dbErr.ToHttp())
			return
		}
	}

	var comment Comment
	err := json.NewDecoder(r.Body).Decode(&comment)
	if err != nil {
		rt.LoggerAndHttpErrorSender(w, err, utils.HttpError{StatusCode: 400})
		return
	}

	dbErr = rt.db.CommentPhoto(authUserId, photoId, comment.Content)
	if dbErr.Err != nil {
		rt.LoggerAndHttpErrorSender(w, dbErr.Err, dbErr.ToHttp())
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	_, _ = w.Write([]byte("Comment added succesfully"))
}

func (rt *_router) getPhotoComments(w http.ResponseWriter, r *http.Request, params map[string]int64) {
	authUserId := params["token"]
	userId := params["user_id"]
	photoId, _ := params["photo_id"]

	isBanned, dbErr := rt.db.IsUserAlreadyTargeted(userId, authUserId, database.BanTable)
	if isBanned {
		rt.LoggerAndHttpErrorSender(w, errors.New("can't get the comments of a photo of a user that banned you"), utils.HttpError{StatusCode: 403})
		return
	} else {
		if dbErr.Err != nil {
			rt.LoggerAndHttpErrorSender(w, dbErr.Err, dbErr.ToHttp())
			return
		}
	}

	var commentsObject CommentsObject
	dbComments, dbErr := rt.db.GetPhotoComments(photoId)
	if dbErr.Err != nil {
		rt.LoggerAndHttpErrorSender(w, dbErr.Err, dbErr.ToHttp())
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

	dbErr := rt.db.DeleteComment(commentId)
	if dbErr.Err != nil {
		rt.LoggerAndHttpErrorSender(w, dbErr.Err, dbErr.ToHttp())
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	_, _ = w.Write([]byte("Comment deleted successfully"))
}
