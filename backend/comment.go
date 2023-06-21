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
	Email    string `json:"email"`
	Fullname string `json:"fullname"`
	Content  string `json:"content"`
	Picture  string `json:"picture"`
}

// # commentNewHandler creates a new comment on a post
//
// - @param postID
// - @param content
func commentNewHandler(w http.ResponseWriter, r *http.Request) {
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
	var data CommentRequest
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
	// get user id form the cookie
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
	ID, err := getIDbyUUID(myuuid)
	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(500)
		jsonResponse, _ := json.Marshal(map[string]string{
			"message": "internal server error, could not get user id",
		})
		w.Write(jsonResponse)
		return
	}
	// convert data.Picture to blob for sqlite
	pictureBlob := []byte{}
	// check the avatar validity
	if data.Picture != "" {
		avatarData, err := base64.StdEncoding.DecodeString(data.Picture)
		if err != nil {
			log.Println(err.Error())
			w.WriteHeader(http.StatusUnprocessableEntity)
			jsonResponse, _ := json.Marshal(map[string]string{
				"message": "Invalid avatar",
			})
			w.Write(jsonResponse)
			return
		}
		if !isImage(avatarData) {
			log.Println(err.Error())
			w.WriteHeader(http.StatusUnsupportedMediaType)
			jsonResponse, _ := json.Marshal(map[string]string{
				"message": "avatar is not a valid image",
			})
			w.Write(jsonResponse)
			return
		}
		pictureBlob = avatarData
	}
	_, err = statements["addComment"].Exec(ID, data.PostID, data.Content, pictureBlob, time.Now().Format("2006-01-02 15:04:05"))
	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(500)
		jsonResponse, _ := json.Marshal(map[string]string{
			"message": "internal server error, addComment query failed",
		})
		w.Write(jsonResponse)
		return
	}
	w.WriteHeader(200)
	jsonResponse, _ := json.Marshal(map[string]string{
		"message": "Comment created",
	})
	w.Write(jsonResponse)
	rows, err := statements["getComments"].Query(data.PostID)
	// TODO: superfulous error return to http client, need to fix by having error return to client overr websockets
	// if err != nil {
	// 	w.WriteHeader(500)
	// 	jsonResponse, _ := json.Marshal(map[string]string{
	// 		"message": "internal server error",
	// 	})
	// 	w.Write(jsonResponse)
	// 	return
	// }
	var comment Comment
	rows.Next()
	rows.Scan(&comment.Fullname, &comment.Content)
	rows.Close()
	sendComment(data.PostID, comment)
}

// # commentGetHandler returns all comments for a post
//
// - @param postID int
func commentGetHandler(w http.ResponseWriter, r *http.Request) {
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
		PostID int `json:"postID"`
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
	rows, err := statements["getComments"].Query(data.PostID)
	if err != nil {
		w.WriteHeader(500)
		jsonResponse, _ := json.Marshal(map[string]string{
			"message": "internal server error",
		})
		w.Write(jsonResponse)
		return
	}
	var comments Comments
	for rows.Next() {
		var comment Comment
		var firstName, lastName string
		var pictureBlob []byte
		rows.Scan(&firstName, &lastName, &comment.Content, &pictureBlob)
		comment.Fullname = firstName + " " + lastName
		comment.Picture = base64.StdEncoding.EncodeToString(pictureBlob)
		comments.Comments = append(comments.Comments, comment)
	}
	rows.Close()
	w.WriteHeader(200)
	jsonResponse, _ := json.Marshal(comments)
	w.Write(jsonResponse)
}
