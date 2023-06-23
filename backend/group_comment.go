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
	defer func() {
		if err := recover(); err != nil {
			log.Println(err)
			w.WriteHeader(500)
			jsonResponse, _ := json.Marshal(map[string]string{
				"message": "internal server error",
			})
			w.Write(jsonResponse)
		}
	}()
	var data GroupCommentRequest
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&data)
	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(400)
		jsonResponse, _ := json.Marshal(map[string]string{
			"message": "Bad request",
		})
		w.Write(jsonResponse)
		return
	}
	// get user id from the cookie
	cookie, err := r.Cookie("user_uuid")
	if err != nil {
		w.WriteHeader(401)
		jsonResponse, _ := json.Marshal(map[string]string{
			"message": "Unauthorized",
		})
		w.Write(jsonResponse)
		return
	}
	myuuid := cookie.Value
	userID, err := getIDbyUUID(myuuid)
	if err != nil {
		w.WriteHeader(500)
		jsonResponse, _ := json.Marshal(map[string]string{
			"message": "internal server error",
		})
		w.Write(jsonResponse)
		return
	}
	_, err = statements["addGroupComment"].Exec(userID, data.GroupPostID, data.Content, data.Picture, time.Now().Format("2006-01-02 15:04:05"))
	if err != nil {
		w.WriteHeader(500)
		jsonResponse, _ := json.Marshal(map[string]string{
			"message": "internal server error",
		})
		w.Write(jsonResponse)
		return
	}
	w.WriteHeader(200)
	jsonResponse, _ := json.Marshal(map[string]string{
		"message": "Comment created",
	})
	w.Write(jsonResponse)
	rows, err := statements["getGroupComments"].Query(data.GroupPostID)
	if err != nil {
		w.WriteHeader(500)
		jsonResponse, _ := json.Marshal(map[string]string{
			"message": "internal server error",
		})
		w.Write(jsonResponse)
		return
	}
	// todo: 99% it must be remastered too, but 1% still exists :D . Check it precisely
	var comment Comment
	rows.Next()
	var firstName, lastName, nickname string
	var pictureBlob []byte
	rows.Scan(&comment.Email, &firstName, &lastName, &nickname, &comment.Content, &pictureBlob, &comment.CreatedAt)
	comment.Fullname = firstName + " " + lastName
	comment.Picture = base64.StdEncoding.EncodeToString(pictureBlob)
	rows.Close()
	sendComment(data.GroupPostID, comment)
}

// # groupCommentsGetHandler returns all comments for a post
//
// - @param group_post_id
func groupCommentsGetHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	defer func() {
		if err := recover(); err != nil {
			log.Println(err)
			w.WriteHeader(500)
			jsonResponse, _ := json.Marshal(map[string]string{
				"message": "internal server error",
			})
			w.Write(jsonResponse)
		}
	}()
	var data struct {
		GroupPostID int `json:"group_post_id"`
	}
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&data)
	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(400)
		jsonResponse, _ := json.Marshal(map[string]string{
			"message": "Bad request",
		})
		w.Write(jsonResponse)
		return
	}
	rows, err := statements["getGroupComments"].Query(data.GroupPostID)
	if err != nil {
		w.WriteHeader(500)
		jsonResponse, _ := json.Marshal(map[string]string{
			"message": "internal server error",
		})
		w.Write(jsonResponse)
		return
	}

	// todo: 99% it must be remastered too, but 1% still exists :D . Check it precisely
	var comments Comments
	for rows.Next() {
		var comment Comment
		var firstName, lastName string
		var pictureBlob []byte
		rows.Scan(&comment.Email, &firstName, &lastName, &comment.Content, &pictureBlob, &comment.CreatedAt)
		comment.Fullname = firstName + " " + lastName
		comment.Picture = base64.StdEncoding.EncodeToString(pictureBlob)
		comments.Comments = append(comments.Comments, comment)
	}
	rows.Close()
	w.WriteHeader(200)
	jsonResponse, _ := json.Marshal(comments)
	w.Write(jsonResponse)
}
