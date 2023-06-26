package main

import (
	"encoding/base64"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

type WSPostRequest struct {
	Title      string `json:"title"`
	Categories string `json:"categories"`
	Content    string `json:"content"`
	Privacy    string `json:"privacy"`
	Picture    string `json:"picture"`
	CreatedAt  string `json:"created_at"`
	AbleToSee  string `json:"able_to_see"`
}

type WSPosts struct {
	WSPosts []WSPost `json:"posts"`
}

type WSPost struct {
	ID         int    `json:"ID"`
	UserID     int    `json:"user_id"`
	Title      string `json:"title"`
	Categories string `json:"categories"`
	Content    string `json:"content"`
	Picture    string `json:"picture"`
	CreatedAt  string `json:"created_at"`
}

type WSPostDTOoutElement struct {
	ID              int    `json:"id"`
	Title           string `json:"title"`
	Content         string `json:"content"`
	Categories      string `json:"categories"`
	Picture         string `json:"picture"`
	CreatorFullName string `json:"creatorFullName"`
	CreatorEmail    string `json:"creatorEmail"`
	CreatedAt       string `json:"createdAt"`
}

func wsPostSubmitHandler(conn *websocket.Conn, messageData map[string]interface{}) {
	defer wsRecover()

	var data WSPostRequest

	data.Categories = sanitizeCategories(data.Categories)
	// process the picture
	postPicture := []byte{}
	if data.Picture != "" {
		avatarData, err := base64.StdEncoding.DecodeString(data.Picture)
		if err != nil {
			log.Println(err.Error())
			// jsonResponse(w, http.StatusUnprocessableEntity, "Invalid avatar")
			return
		}
		if !isImage(avatarData) {
			log.Println(err.Error())
			// jsonResponse(w, http.StatusUnsupportedMediaType, "avatar is not a valid image")
			return
		}
		postPicture = avatarData
	}

	userID := messageData["userID"].(int) // todo: FAKE code

	data.CreatedAt = time.Now().Format("2006-01-02 15:04:05")
	result, err := statements["addPost"].Exec(userID, data.Title, data.Categories, data.Content, data.Privacy, postPicture, data.CreatedAt)
	if err != nil {
		// jsonResponse(w, http.StatusInternalServerError, "addPost query failed")
		return
	}

	postId, err := result.LastInsertId()
	if err != nil {
		// jsonResponse(w, http.StatusInternalServerError, "LastInsertId of addPost query failed")
		return
	}
	// privacy check and able to see
	if data.Privacy == "almost private" {
		if data.AbleToSee != "" {
			listOfEmails := strings.Split(data.AbleToSee, " ")
			for _, email := range listOfEmails {
				// get the id of the user from the email
				userID, err := getIDbyEmail(email)
				if err != nil {
					// jsonResponse(w, http.StatusUnprocessableEntity, "Invalid email")
					return
				}
				// add the post to the almost_private table (user_id, post_id)
				_, err = statements["addAlmostPrivate"].Exec(userID, postId)
				if err != nil {
					// jsonResponse(w, http.StatusInternalServerError, "addAlmostPrivate query failed")
					return
				}
			}
		}
	}

	// jsonResponse(w, http.StatusOK, "Post created")

	rows, err := statements["getPosts"].Query(userID)
	if err != nil {
		log.Println(err.Error())
		// jsonResponse(w, http.StatusInternalServerError, "getPosts query failed")
		return
	}
	var post Post
	rows.Next()
	err = rows.Scan(&post.ID, &post.UserID, &post.Title, &post.Categories, &post.Content)
	if err != nil {
		// jsonResponse(w, http.StatusInternalServerError, "post scan failed")
		return
	}
	rows.Close()
	sendPost(post)
}

func wsPostsGetHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	defer recovery(w)
	rows, err := statements["getPosts"].Query()
	if err != nil {
		jsonResponse(w, http.StatusInternalServerError, "getPosts query failed")
		return
	}
	var wsPosts WSPosts
	for rows.Next() {
		var post WSPost
		pictureBlob := []byte{}
		err = rows.Scan(&post.ID, &post.UserID, &post.Title, &post.Categories, &post.Content, &pictureBlob)
		if err != nil {
			jsonResponse(w, http.StatusInternalServerError, "post scan failed")
			return
		}
		post.Picture = base64.StdEncoding.EncodeToString(pictureBlob)
		wsPosts.WSPosts = append(wsPosts.WSPosts, post)
	}
	rows.Close()
	w.WriteHeader(200)
	jsonResponseObj, err := json.Marshal(wsPosts)
	if err != nil {
		jsonResponse(w, http.StatusInternalServerError, "json.Marshal failed")
		return
	}
	_, err = w.Write(jsonResponseObj)
	if err != nil {
		jsonResponse(w, http.StatusInternalServerError, "w.Write(jsonResponseObj)<-posts failed")
		return
	}
}
