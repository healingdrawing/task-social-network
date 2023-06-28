package main

import (
	"encoding/base64"
	"fmt"
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

type WS_POSTS_LIST_DTO []WS_POST_RESPONSE_DTO

func wsPostSubmitHandler(conn *websocket.Conn, messageData map[string]interface{}) {
	defer wsRecover()

	log.Println("wsPostSubmitHandler-------NEW POST") //todo: remove this
	for key, value := range messageData {
		log.Println("KEY: ", key, "\nVALUE: ", value)
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
			wsSendError(WS_ERROR_RESPONSE_DTO{fmt.Sprint(http.StatusUnprocessableEntity) + " Invalid avatar"})
			return
		}
		if !isImage(avatarData) {
			log.Println("avatar is not a valid image", err.Error())
			wsSendError(WS_ERROR_RESPONSE_DTO{fmt.Sprint(http.StatusUnsupportedMediaType) + " avatar is not a valid image"})
			return
		}
		postPicture = avatarData
	}

	userID, err := getIDbyUUID(data.User_uuid)
	if err != nil {
		log.Println("failed to get ID of the request sender", err.Error())
		wsSendError(WS_ERROR_RESPONSE_DTO{fmt.Sprint(http.StatusUnprocessableEntity) + " failed to get ID of the request sender"})
	}

	data.Created_at = time.Now().Format("2006-01-02 15:04:05")
	result, err := statements["addPost"].Exec(userID, data.Title, data.Categories, data.Content, data.Privacy, postPicture, data.Created_at)
	if err != nil {
		log.Println("addPost query failed", err.Error())
		wsSendError(WS_ERROR_RESPONSE_DTO{fmt.Sprint(http.StatusInternalServerError) + " addPost query failed"})
		return
	}

	postId, err := result.LastInsertId()
	if err != nil {
		log.Println("LastInsertId of addPost query failed", err.Error())
		wsSendError(WS_ERROR_RESPONSE_DTO{fmt.Sprint(http.StatusInternalServerError) + " LastInsertId of addPost query failed"})
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
					wsSendError(WS_ERROR_RESPONSE_DTO{fmt.Sprint(http.StatusUnprocessableEntity) + " Invalid email"})
					return
				}
				// add the post to the almost_private table (user_id, post_id)
				_, err = statements["addAlmostPrivate"].Exec(userID, postId)
				if err != nil {
					log.Println("addAlmostPrivate query failed", err.Error())
					wsSendError(WS_ERROR_RESPONSE_DTO{fmt.Sprint(http.StatusInternalServerError) + " addAlmostPrivate query failed"})
					return
				}
			}
		}
	}

	rows, err := statements["getPosts"].Query(userID)
	if err != nil {
		log.Println("getPosts query failed", err.Error())
		wsSendError(WS_ERROR_RESPONSE_DTO{fmt.Sprint(http.StatusInternalServerError) + " getPosts query failed"})
		return
	}
	defer rows.Close()

	var post WS_POST_RESPONSE_DTO
	pictureBytes := []byte{}
	rows.Next()
	err = rows.Scan(&post.ID, &post.Title, &post.Content, &post.Categories, &pictureBytes, &post.Privacy, &post.Created_at, &post.Email, &post.First_name, &post.Last_name)
	if err != nil {
		log.Println("post scan failed", err.Error())
		wsSendError(WS_ERROR_RESPONSE_DTO{fmt.Sprint(http.StatusInternalServerError) + " post scan failed"})
		return
	}
	post.Picture = base64.StdEncoding.EncodeToString(pictureBytes)

	wsSendPost(post)
}

func wsPostsListHandler(conn *websocket.Conn, messageData map[string]interface{}) {
	defer wsRecover()

	user_uuid := messageData["user_uuid"].(string)
	userID, err := getIDbyUUID(user_uuid)
	if err != nil {
		log.Println("failed to get ID of the request sender", err.Error())
		wsSendError(WS_ERROR_RESPONSE_DTO{fmt.Sprint(http.StatusUnprocessableEntity) + " failed to get ID of the request sender"})
	}

	rows, err := statements["getPostsAbleToSee"].Query(userID, userID, userID)
	if err != nil {
		log.Println("getPosts query failed", err.Error())
		wsSendError(WS_ERROR_RESPONSE_DTO{fmt.Sprint(http.StatusInternalServerError) + " getPostsAbleToSee query failed"})
		return
	}
	defer rows.Close()

	var postsList WS_POSTS_LIST_DTO
	for rows.Next() {
		var post WS_POST_RESPONSE_DTO
		pictureBytes := []byte{}
		err = rows.Scan(&post.ID, &post.Title, &post.Content, &post.Categories, &pictureBytes, &post.Privacy, &post.Created_at, &post.Email, &post.First_name, &post.Last_name)
		if err != nil {
			log.Println("post scan failed", err.Error())
			wsSendError(WS_ERROR_RESPONSE_DTO{fmt.Sprint(http.StatusInternalServerError) + " post scan failed"})
			return
		}
		post.Picture = base64.StdEncoding.EncodeToString(pictureBytes)
		postsList = append(postsList, post)
	}

	wsSendPostsList(postsList)

}

// todo: global
// Picture not managed properly.
// Need:
// - link to profile of post author
// - some sorting of posts on backend side by able_to_see emails
// - convert to picture for screen after server respond to frontend
// - clear the storage when picture uplaoaded, or fail checks.

// Now the last successed picture will be uploaded on the next post automatically.
