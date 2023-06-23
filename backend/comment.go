package main

import (
	"encoding/base64"
	"encoding/json"
	"log"
	"net/http"
	"time"
)

type CommentRequest struct {
	PostID  int    `json:"postID"`
	Content string `json:"content"`
	Picture string `json:"picture"`
}

type Comments struct {
	Comments []Comment `json:"comments"`
}

type Comment struct {
	Email     string `json:"email"`
	Fullname  string `json:"fullname"`
	Content   string `json:"content"`
	Picture   string `json:"picture"`
	CreatedAt string `json:"created_at"`
}

// # commentNewHandler creates a new comment on a post
//
// - @param postID
// - @param content
func commentNewHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	defer recovery(w)
	var data CommentRequest
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&data)
	if err != nil {
		log.Println(err.Error())
		jsonResponse(w, http.StatusBadRequest, "")
		return
	}
	// get user id form the cookie
	cookie, err := r.Cookie("user_uuid")
	if err != nil {
		jsonResponse(w, http.StatusUnauthorized, "")
		return
	}
	myuuid := cookie.Value
	ID, err := getIDbyUUID(myuuid)
	if err != nil {
		log.Println(err.Error())
		jsonResponse(w, http.StatusUnauthorized, "")
		return
	}
	// convert data.Picture to blob for sqlite
	pictureBlob := []byte{}
	// check the avatar validity
	if data.Picture != "" {
		avatarData, err := base64.StdEncoding.DecodeString(data.Picture)
		if err != nil {
			log.Println(err.Error())
			jsonResponse(w, http.StatusUnprocessableEntity, "Invalid avatar")
			return
		}
		if !isImage(avatarData) {
			log.Println(err.Error())
			jsonResponse(w, http.StatusUnsupportedMediaType, "Avatar is not a valid image")
			return
		}
		pictureBlob = avatarData
	}
	_, err = statements["addComment"].Exec(ID, data.PostID, data.Content, pictureBlob, time.Now().Format("2006-01-02 15:04:05"))
	if err != nil {
		log.Println(err.Error())
		jsonResponse(w, http.StatusInternalServerError, "addComment query failed")
		return
	}
	jsonResponse(w, http.StatusOK, "Comment created")

	// TODO: UNCOMMENT+REFACTOR THIS LATER!!! Old code is commented by / * * / because chat is not ready now for new version
	/*
		rows, err := statements["getComments"].Query(data.PostID)
		// TODO: superfulous error return to http client, need to fix by having error return to client overr websockets
		// if err != nil {
		// 	w.WriteHeader(500)
		// 	jsonResponseObj, _ := json.Marshal(map[string]string{
		// 		"message": "internal server error",
		// 	})
		// 	w.Write(jsonResponseObj)
		// 	return
		// }
		var comment Comment
		rows.Next()
		rows.Scan(&comment.Fullname, &comment.Content)
		rows.Close()
		sendComment(data.PostID, comment)
	*/
}

// # commentsGetHandler returns all comments for a post
//
// - @param postID int
func commentsGetHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	defer recovery(w)
	var data struct {
		PostID int `json:"postID"`
	}
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&data)
	if err != nil {
		log.Println(err.Error())
		jsonResponse(w, http.StatusBadRequest, "")
		return
	}
	rows, err := statements["getComments"].Query(data.PostID)
	if err != nil {
		jsonResponse(w, http.StatusInternalServerError, "")
		return
	}
	var comments Comments
	for rows.Next() {
		var comment Comment
		var firstName, lastName string
		var pictureBlob []byte
		err = rows.Scan(&firstName, &lastName, &comment.Content, &pictureBlob)
		if err != nil {
			jsonResponse(w, http.StatusInternalServerError, "rows.Scan comments for loop failed")
			return
		}
		comment.Fullname = firstName + " " + lastName
		comment.Picture = base64.StdEncoding.EncodeToString(pictureBlob)
		comments.Comments = append(comments.Comments, comment)
	}
	rows.Close()
	w.WriteHeader(200)
	jsonResponseObject, err := json.Marshal(comments)
	if err != nil {
		jsonResponse(w, http.StatusInternalServerError, "json.Marshal comments failed")
		return
	}
	_, err = w.Write(jsonResponseObject)
	if err != nil {
		jsonResponse(w, http.StatusInternalServerError, "w.Write comments failed") // todo: CHECK IT! handler was added, in old code it was not
		return
	}
}
