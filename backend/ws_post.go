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
	Id         int    `json:"id"` // post section. use to get comments
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

// wsPostSubmitHandler creates a new post, then return all posts, which user can see.
func wsPostSubmitHandler(conn *websocket.Conn, messageData map[string]interface{}) {
	defer wsRecover(messageData)

	uuid, ok := messageData["user_uuid"].(string)
	if !ok {
		log.Println("failed to get user_uuid from message data")
		return
	}

	var data WS_POST_SUBMIT_DTO
	data.User_uuid = uuid

	fields := map[string]*string{
		"title":       &data.Title,
		"categories":  &data.Categories,
		"content":     &data.Content,
		"privacy":     &data.Privacy,
		"picture":     &data.Picture,
		"able_to_see": &data.Able_to_see,
	}

	for key, ptr := range fields {
		value, ok := messageData[key].(string)
		if !ok {
			log.Printf("failed to get %s from messageData\n", key)
			wsSend(WS_ERROR_RESPONSE, WS_ERROR_RESPONSE_DTO{fmt.Sprintf("%d failed to get %s from messageData", http.StatusUnprocessableEntity, key)}, []string{uuid})
			return
		}
		*ptr = value
	}

	data.Categories = sanitizeCategories(data.Categories)
	// process the picture
	postPicture := []byte{}
	if data.Picture != "" {
		imageData, err := extractImageData(data.Picture)
		if err != nil {
			log.Println("=FAIL extractPictureData: ", err.Error())
			wsSend(WS_ERROR_RESPONSE, WS_ERROR_RESPONSE_DTO{fmt.Sprint(http.StatusUnprocessableEntity) + " =FAIL extractPictureData:" + err.Error()}, []string{uuid})
			return
		}
		pictureData, err := base64.StdEncoding.DecodeString(imageData)
		if err != nil {
			log.Println("Invalid picture ", err.Error())
			wsSend(WS_ERROR_RESPONSE, WS_ERROR_RESPONSE_DTO{fmt.Sprint(http.StatusUnprocessableEntity) + " Invalid picture"}, []string{uuid})
			return
		}
		if !isImage(pictureData) {
			log.Println("picture is not a valid image", err.Error())
			wsSend(WS_ERROR_RESPONSE, WS_ERROR_RESPONSE_DTO{fmt.Sprint(http.StatusUnsupportedMediaType) + " picture is not a valid image"}, []string{uuid})
			return
		}
		postPicture = pictureData
	}

	user_id, err := get_user_id_by_uuid(data.User_uuid)
	if err != nil {
		log.Println("failed to get ID of the message sender", err.Error())
		wsSend(WS_ERROR_RESPONSE, WS_ERROR_RESPONSE_DTO{fmt.Sprint(http.StatusUnprocessableEntity) + " failed to get ID of the message sender"}, []string{uuid})
		return
	}

	// privacy check
	if _, ok := map[string]int{"public": 0, "private": 0, "almost private": 0}[data.Privacy]; !ok {
		log.Println("Invalid privacy ", err.Error())
		wsSend(WS_ERROR_RESPONSE, WS_ERROR_RESPONSE_DTO{fmt.Sprint(http.StatusUnprocessableEntity) + " Invalid privacy"}, []string{uuid})
		return
	}

	data.Created_at = time.Now().Format("2006-01-02 15:04:05")
	result, err := statements["addPost"].Exec(user_id, data.Title, data.Categories, data.Content, data.Privacy, postPicture, data.Created_at)
	if err != nil {
		log.Println("addPost query failed", err.Error())
		wsSend(WS_ERROR_RESPONSE, WS_ERROR_RESPONSE_DTO{fmt.Sprint(http.StatusInternalServerError) + " addPost query failed"}, []string{uuid})
		return
	}

	postId, err := result.LastInsertId()
	if err != nil {
		log.Println("LastInsertId of addPost query failed", err.Error())
		wsSend(WS_ERROR_RESPONSE, WS_ERROR_RESPONSE_DTO{fmt.Sprint(http.StatusInternalServerError) + " LastInsertId of addPost query failed"}, []string{uuid})
		return
	}
	// privacy check and able to see
	if data.Privacy == "almost private" {
		if data.Able_to_see != "" {
			listOfEmails := strings.Split(data.Able_to_see, " ")
			for _, email := range listOfEmails {
				// get the id of the user from the email
				userID, err := get_user_id_by_email(email)
				if err != nil {
					log.Println("Invalid email ", err.Error())
					wsSend(WS_ERROR_RESPONSE, WS_ERROR_RESPONSE_DTO{fmt.Sprint(http.StatusUnprocessableEntity) + " Invalid email"}, []string{uuid})
					return
				}
				// add the post to the almost_private table (user_id, post_id)
				_, err = statements["addAlmostPrivate"].Exec(userID, postId)
				if err != nil {
					log.Println("addAlmostPrivate query failed", err.Error())
					wsSend(WS_ERROR_RESPONSE, WS_ERROR_RESPONSE_DTO{fmt.Sprint(http.StatusInternalServerError) + " addAlmostPrivate query failed"}, []string{uuid})
					return
				}
			}
		}
	}

	wsSend(WS_SUCCESS_RESPONSE, WS_SUCCESS_RESPONSE_DTO{fmt.Sprint(http.StatusOK) + " post created"}, []string{uuid})

	// return all posts, which user can see.
	wsPostsListHandler(conn, messageData)
}

func wsPostsListHandler(conn *websocket.Conn, messageData map[string]interface{}) {
	defer wsRecover(messageData)

	uuid, ok := messageData["user_uuid"].(string)
	if !ok {
		log.Println("failed to get user_uuid from messageData")
		return
	}
	user_id, err := get_user_id_by_uuid(uuid)
	if err != nil {
		log.Println("failed to get ID of the message sender", err.Error())
		wsSend(WS_ERROR_RESPONSE, WS_ERROR_RESPONSE_DTO{fmt.Sprint(http.StatusUnprocessableEntity) + " failed to get ID of the message sender"}, []string{uuid})
		return
	}

	rows, err := statements["getPostsAbleToSee"].Query(user_id, user_id, user_id)
	if err != nil {
		log.Println("getPostsAbleToSee query failed", err.Error())
		wsSend(WS_ERROR_RESPONSE, WS_ERROR_RESPONSE_DTO{fmt.Sprint(http.StatusInternalServerError) + " getPostsAbleToSee query failed"}, []string{uuid})
		return
	}
	defer rows.Close()

	var postsList WS_POSTS_LIST_DTO
	for rows.Next() {
		var post WS_POST_RESPONSE_DTO
		pictureBytes := []byte{}
		err = rows.Scan(&post.Id, &post.Title, &post.Content, &post.Categories, &pictureBytes, &post.Privacy, &post.Created_at, &post.Email, &post.First_name, &post.Last_name)
		if err != nil {
			log.Println("post scan failed", err.Error())
			wsSend(WS_ERROR_RESPONSE, WS_ERROR_RESPONSE_DTO{fmt.Sprint(http.StatusInternalServerError) + " post scan failed"}, []string{uuid})
			return
		}
		post.Picture = base64.StdEncoding.EncodeToString(pictureBytes)
		postsList = append(postsList, post)
	}

	wsSend(WS_POSTS_LIST, postsList, []string{uuid})
}

/** includes the posts created by user, BUT NOT GROUP POSTS, created by user */
func wsUserPostsListHandler(conn *websocket.Conn, messageData map[string]interface{}) {
	defer wsRecover(messageData)

	uuid, ok := messageData["user_uuid"].(string)
	if !ok {
		log.Println("failed to get user_uuid from messageData")
		return
	}
	user_id, err := get_user_id_by_uuid(uuid)
	if err != nil {
		log.Println("failed to get ID of the message sender", err.Error())
		wsSend(WS_ERROR_RESPONSE, WS_ERROR_RESPONSE_DTO{fmt.Sprint(http.StatusUnprocessableEntity) + " failed to get ID of the message sender"}, []string{uuid})
		return
	}

	target_email := messageData["target_email"].(string)
	target_id, err := get_user_id_by_email(target_email)
	if err != nil {
		log.Println("failed to get ID of the target", err.Error())
		wsSend(WS_ERROR_RESPONSE, WS_ERROR_RESPONSE_DTO{fmt.Sprint(http.StatusUnprocessableEntity) + " failed to get ID of the target"}, []string{uuid})
		return
	}

	rows, err := statements["getPostsAbleToSeeToVisitor"].Query(target_id, user_id, target_id, target_id, target_id, target_id, user_id, target_id, user_id)
	if err != nil {
		log.Println("getPostsAbleToSeeToVisitor query failed", err.Error())
		wsSend(WS_ERROR_RESPONSE, WS_ERROR_RESPONSE_DTO{fmt.Sprint(http.StatusInternalServerError) + " getPostsAbleToSeeToVisitor query failed"}, []string{uuid})
		return
	}
	defer rows.Close()

	var postsList WS_POSTS_LIST_DTO
	for rows.Next() {
		var post WS_POST_RESPONSE_DTO
		pictureBytes := []byte{}
		err = rows.Scan(&post.Id, &post.Title, &post.Content, &post.Categories, &pictureBytes, &post.Privacy, &post.Created_at, &post.Email, &post.First_name, &post.Last_name)
		if err != nil {
			log.Println("post scan failed", err.Error())
			wsSend(WS_ERROR_RESPONSE, WS_ERROR_RESPONSE_DTO{fmt.Sprint(http.StatusInternalServerError) + " post scan failed"}, []string{uuid})
			return
		}
		post.Picture = base64.StdEncoding.EncodeToString(pictureBytes)
		postsList = append(postsList, post)
	}

	wsSend(WS_POSTS_LIST, postsList, []string{uuid})
}
