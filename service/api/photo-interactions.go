package api

import (
	"encoding/json"
	"errors"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strconv"
	"wasaphoto/service/database"
	"wasaphoto/service/utils"
)

// todo controllare che la foto appartenga all'utente nel path
func (rt *_router) likePhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, token utils.Token) {
	authUserId, _ := strconv.ParseInt(token.Value, 10, 64)
	paramAuthUserId, _ := strconv.ParseInt(ps.ByName("auth_user_id"), 10, 64)
	photoOwnerId, _ := strconv.ParseInt(ps.ByName("user_id"), 10, 64)

	if authUserId != paramAuthUserId {
		rt.LoggerAndHttpErrorSender(w, errors.New("can't put a like impersonating someone else"), utils.HttpError{StatusCode: 403})
		return
	}

	photoId, _ := strconv.ParseInt(ps.ByName("photo_id"), 10, 64)

	isBanned, dbErr := rt.db.IsUserAlreadyTargeted(photoOwnerId, authUserId, database.BanTable)
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

// todo controllare che la foto appartieni all'utente nel path
func (rt *_router) unlikePhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, token utils.Token) {
	authUserId, _ := strconv.ParseInt(token.Value, 10, 64)
	paramAuthUserId, _ := strconv.ParseInt(ps.ByName("auth_user_id"), 10, 64)
	photoOwnerId, _ := strconv.ParseInt(ps.ByName("user_id"), 10, 64)

	if authUserId != paramAuthUserId {
		rt.LoggerAndHttpErrorSender(w, errors.New("can't delete a like impersonating someone else"), utils.HttpError{StatusCode: 403})
		return
	}

	photoId, _ := strconv.ParseInt(ps.ByName("photo_id"), 10, 64)

	isBanned, dbErr := rt.db.IsUserAlreadyTargeted(photoOwnerId, authUserId, database.BanTable)

	if isBanned {
		rt.LoggerAndHttpErrorSender(w, errors.New("can't unlike a photo of a user that banned you"), utils.HttpError{StatusCode: 403})
		return
	} else {
		if dbErr.Err != nil {
			rt.LoggerAndHttpErrorSender(w, dbErr.Err, dbErr.ToHttp())
			return
		}
	}

	dbErr = rt.db.UnlikePhoto(authUserId, photoId)
	if dbErr.Err != nil {
		rt.LoggerAndHttpErrorSender(w, dbErr.Err, dbErr.ToHttp())
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	_, _ = w.Write([]byte("Photo unliked successfully"))
}

// todo controllare che la foto appartieni all'utente nel path e commento appartenga alla foto
func (rt *_router) commentPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, token utils.Token) {
	authUserId, _ := strconv.ParseInt(token.Value, 10, 64)

	photoId, _ := strconv.ParseInt(ps.ByName("photo_id"), 10, 64)
	photoOwnerId, _ := strconv.ParseInt(ps.ByName("user_id"), 10, 64)

	isBanned, dbErr := rt.db.IsUserAlreadyTargeted(photoOwnerId, authUserId, database.BanTable)
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

// todo controllare che la foto appartieni all'utente nel path
func (rt *_router) getPhotoComments(w http.ResponseWriter, r *http.Request, ps httprouter.Params, token utils.Token) {
	authUserId, _ := strconv.ParseInt(token.Value, 10, 64)
	photoId, _ := strconv.ParseInt(ps.ByName("photo_id"), 10, 64)
	photoOwnerId, _ := strconv.ParseInt(ps.ByName("user_id"), 10, 64)
	// se owner ha bloccato authUserId, errore
	isBanned, dbErr := rt.db.IsUserAlreadyTargeted(photoOwnerId, authUserId, database.BanTable)
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

// todo controllare che la foto appartieni all'utente nel path, che il commento appartenga alla foto e che l'utente che vuole cancellare il commento sia il proprietario del commento
func (rt *_router) deleteComment(w http.ResponseWriter, r *http.Request, ps httprouter.Params, token utils.Token) {
	authUserId, _ := strconv.ParseInt(token.Value, 10, 64)
	photoOwnerId, _ := strconv.ParseInt(ps.ByName("user_id"), 10, 64)
	commentId, _ := strconv.ParseInt(ps.ByName("comment_id"), 10, 64)

	isBanned, dbErr := rt.db.IsUserAlreadyTargeted(photoOwnerId, authUserId, database.BanTable)
	if isBanned {
		rt.LoggerAndHttpErrorSender(w, errors.New("can't get the comments of a photo of a user that banned you"), utils.HttpError{StatusCode: 403})
		return
	} else {
		if dbErr.Err != nil {
			rt.LoggerAndHttpErrorSender(w, dbErr.Err, dbErr.ToHttp())
			return
		}
	}
	dbErr = rt.db.DeleteComment(commentId)
	if dbErr.Err != nil {
		rt.LoggerAndHttpErrorSender(w, dbErr.Err, dbErr.ToHttp())
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	_, _ = w.Write([]byte("Comment deleted successfully"))
}
