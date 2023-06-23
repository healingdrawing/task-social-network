package main

import (
	"encoding/base64"
	"encoding/json"
	"log"
	"net/http"
	"time"
)

type GroupCommentRequest struct {
	GroupPostID int    `json:"group_post_id"`
	Content     string `json:"content"`
	Picture     string `json:"picture"`
}

type GroupComments struct {
	GroupComments []GroupComment `json:"group_comments"`
}

// todo: Finally on frontend side, the "Username" data should be cummulative of first name and last name plus email in round brackets, cause nickname is not unique. And for "comment.go" too. To provide clickable link on user profile. Or you can not make following request etc.
type GroupComment struct {
	Fullname string `json:"username"`
	Content  string `json:"content"`
	Picture  string `json:"picture"`
}

// # groupCommentNewHandler creates a new comment on a group post
//
// @r.param {group_post_id int, content string, picture string}
func groupCommentNewHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	defer recovery(w)
	var data GroupCommentRequest
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&data)
	if err != nil {
		log.Println(err.Error())
		jsonResponse(w, 400, "Bad request")
		return
	}
	// get user id from the cookie
	cookie, err := r.Cookie("user_uuid")
	if err != nil {
		jsonResponse(w, 401, "cannot get cookie")
		return
	}
	myuuid := cookie.Value
	userID, err := getIDbyUUID(myuuid)
	if err != nil {
		jsonResponse(w, 401, " getIDbyUUID failed")
		return
	}
	_, err = statements["addGroupComment"].Exec(userID, data.GroupPostID, data.Content, data.Picture, time.Now().Format("2006-01-02 15:04:05"))
	if err != nil {
		jsonResponse(w, 500, "addGroupComment query failed")
		return
	}
	jsonResponse(w, 200, "Comment created")

	// rows, err := statements["getGroupComments"].Query(data.GroupPostID)
	// if err != nil {
	// 	jsonResponseWriterManager(w, 500, "getGroupComments query failed")
	// 	return
	// }
	// // todo: 99% it must be remastered too, but 1% still exists :D . Check it precisely
	// var comment Comment
	// rows.Next()
	// var firstName, lastName, nickname string
	// var pictureBlob []byte
	// rows.Scan(&comment.Email, &firstName, &lastName, &nickname, &comment.Content, &pictureBlob, &comment.CreatedAt)
	// comment.Fullname = firstName + " " + lastName
	// comment.Picture = base64.StdEncoding.EncodeToString(pictureBlob)
	// rows.Close()
	// sendComment(data.GroupPostID, comment)
}

// # groupCommentsGetHandler returns all comments for a post
//
// - @param group_post_id
func groupCommentsGetHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	defer recovery(w)
	var data struct {
		GroupPostID int `json:"group_post_id"`
	}
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&data)
	if err != nil {
		log.Println(err.Error())
		jsonResponse(w, 400, "Bad request")
		return
	}
	rows, err := statements["getGroupComments"].Query(data.GroupPostID)
	if err != nil {
		jsonResponse(w, 500, "getGroupComments query failed")
		return
	}

	// todo: 99% it must be remastered too, but 1% still exists :D . Check it precisely
	var comments Comments
	for rows.Next() {
		var comment Comment
		var firstName, lastName string
		var pictureBlob []byte
		err = rows.Scan(&comment.Email, &firstName, &lastName, &comment.Content, &pictureBlob, &comment.CreatedAt)
		if err != nil {
			jsonResponse(w, 500, " scan getGroupComments query failed")
			return
		}
		comment.Fullname = firstName + " " + lastName
		comment.Picture = base64.StdEncoding.EncodeToString(pictureBlob)
		comments.Comments = append(comments.Comments, comment)
	}
	rows.Close()
	w.WriteHeader(200)
	jsonResponseObj, err := json.Marshal(comments)
	if err != nil {
		jsonResponse(w, 500, "json.Marshal(comments) failed")
		return
	}
	_, err = w.Write(jsonResponseObj)
	if err != nil {
		jsonResponse(w, 500, "w.Write(jsonResponseObj)<-comments failed")
		return
	}

}
