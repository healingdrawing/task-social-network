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

type WS_POST_RESPONSE_DTO struct {
	ID         int    `json:"id"` // post section. use to get comments
	Title      string `json:"title"`
	Content    string `json:"content"`
	Categories string `json:"categories"`
	Picture    string `json:"picture"`
	Privacy    string `json:"privacy"`
	Created_at string `json:"created_at"`
	Email      string `json:"email"` //creator section. use to show post credentials
	First_name string `json:"first_name"`
	Last_name  string `json:"last_name"`
}

type WS_POST_SUBMIT_DTO struct {
	User_uuid   string `json:"user_uuid"`
	Title       string `json:"title"`
	Categories  string `json:"categories"`
	Content     string `json:"content"`
	Privacy     string `json:"privacy"`
	Picture     string `json:"picture"`
	Created_at  string `json:"created_at"`
	Able_to_see string `json:"able_to_see"`
}

func wsPostSubmitHandler(conn *websocket.Conn, messageData map[string]interface{}) {
	defer wsRecover()

	for key, value := range messageData {
		log.Println("wsPostSubmitHandler,\nkey: ", key, "\nvalue: ", value)
	}

	var data WS_POST_SUBMIT_DTO
	data.User_uuid = messageData["user_uuid"].(string)
	data.Title = messageData["title"].(string)
	data.Categories = messageData["categories"].(string)
	data.Content = messageData["content"].(string)
	data.Privacy = messageData["privacy"].(string)
	data.Picture = messageData["picture"].(string)
	data.Able_to_see = (messageData["able_to_see"]).(string)

	data.Categories = sanitizeCategories(data.Categories)
	// process the picture
	postPicture := []byte{}
	if data.Picture != "null" { //todo: it was empty string ""
		avatarData, err := base64.StdEncoding.DecodeString(data.Picture)
		if err != nil {
			log.Println("Invalid avatar ", err.Error())
			// jsonResponse(w, http.StatusUnprocessableEntity, "Invalid avatar")
			return
		}
		if !isImage(avatarData) {
			log.Println("avatar is not a valid image", err.Error())
			// jsonResponse(w, http.StatusUnsupportedMediaType, "avatar is not a valid image")
			return
		}
		postPicture = avatarData
	}

	userID, err := getIDbyUUID(data.User_uuid)
	if err != nil {
		log.Println("failed to get ID of the request sender", err.Error()) // todo: some error common message sender needed
	}

	data.Created_at = time.Now().Format("2006-01-02 15:04:05")
	result, err := statements["addPost"].Exec(userID, data.Title, data.Categories, data.Content, data.Privacy, postPicture, data.Created_at)
	if err != nil {
		log.Println("addPost query failed", err.Error())
		// jsonResponse(w, http.StatusInternalServerError, "addPost query failed")
		return
	}

	postId, err := result.LastInsertId()
	if err != nil {
		log.Println("LastInsertId of addPost query failed", err.Error())
		// jsonResponse(w, http.StatusInternalServerError, "LastInsertId of addPost query failed")
		return
	}
	// privacy check and able to see
	if data.Privacy == "almost private" {
		if data.Able_to_see != "" {
			listOfEmails := strings.Split(data.Able_to_see, " ")
			for _, email := range listOfEmails {
				// get the id of the user from the email
				userID, err := getIDbyEmail(email)
				if err != nil {
					log.Println("Invalid email ", err.Error())
					// jsonResponse(w, http.StatusUnprocessableEntity, "Invalid email")
					return
				}
				// add the post to the almost_private table (user_id, post_id)
				_, err = statements["addAlmostPrivate"].Exec(userID, postId)
				if err != nil {
					log.Println("addAlmostPrivate query failed", err.Error())
					// jsonResponse(w, http.StatusInternalServerError, "addAlmostPrivate query failed")
					return
				}
			}
		}
	}

	//fix: fails with scan, because structure was changed, continue refactoring
	// jsonResponse(w, http.StatusOK, "Post created")

	rows, err := statements["getPosts"].Query(userID)
	if err != nil {
		log.Println("getPosts query failed", err.Error())
		// jsonResponse(w, http.StatusInternalServerError, "getPosts query failed")
		return
	}
	var post WS_POST_RESPONSE_DTO
	rows.Next()
	err = rows.Scan(&post.ID, &post.Title, &post.Content, &post.Categories, &post.Picture, &post.Privacy, &post.Created_at, &post.Email, &post.First_name, &post.Last_name)
	if err != nil {
		log.Println("post scan failed", err.Error())
		// jsonResponse(w, http.StatusInternalServerError, "post scan failed")
		return
	}
	rows.Close()
	wsSendPost(post)
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
