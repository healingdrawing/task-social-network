package main

import (
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

type WS_COMMENT_SUBMIT_DTO struct {
	User_uuid string `json:"user_uuid"`
	Post_id   int    `json:"post_id"`
	Content   string `json:"content"`
	Picture   string `json:"picture"`
}

type WS_COMMENT_RESPONSE_DTO struct {
	Content    string `json:"content"`
	Picture    string `json:"picture"`
	Created_at string `json:"created_at"`
	Email      string `json:"email"`
	First_name string `json:"first_name"`
	Last_name  string `json:"last_name"`
}

type WS_COMMENTS_LIST_DTO []WS_COMMENT_RESPONSE_DTO

// wsCommentSubmitHandler creates a new comment on a post, then return all comments on that post.
func wsCommentSubmitHandler(conn *websocket.Conn, messageData map[string]interface{}) {
	defer wsRecover(messageData)

	uuid, ok := messageData["user_uuid"].(string)
	if !ok {
		log.Println("failed to get user_uuid from messageData")
		return
	}

	var data WS_COMMENT_SUBMIT_DTO
	data.User_uuid = uuid

	_post_id, ok := messageData["post_id"].(float64)
	if !ok {
		log.Println("failed to get post_id from messageData")
		wsSend(WS_ERROR_RESPONSE, WS_ERROR_RESPONSE_DTO{fmt.Sprint(http.StatusUnprocessableEntity, " failed to get post_id from messageData")}, []string{uuid})
		return
	}
	data.Post_id = int(_post_id)

	fields := map[string]*string{
		"content": &data.Content,
		"picture": &data.Picture,
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

	user_id, err := get_user_id_by_uuid(data.User_uuid)
	if err != nil {
		log.Println("failed to get ID of the message sender", err.Error())
		wsSend(WS_ERROR_RESPONSE, WS_ERROR_RESPONSE_DTO{fmt.Sprint(http.StatusUnprocessableEntity) + " failed to get ID of the message sender"}, []string{uuid})
		return
	}

	// process the picture
	commentPicture := []byte{}
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
		commentPicture = pictureData
	}

	created_at := time.Now().Format("2006-01-02 15:04:05")

	_, err = statements["addComment"].Exec(user_id, data.Post_id, data.Content, commentPicture, created_at)
	if err != nil {
		log.Println("addComment query failed", err.Error())
		wsSend(WS_ERROR_RESPONSE, WS_ERROR_RESPONSE_DTO{fmt.Sprint(http.StatusInternalServerError) + " addComment query failed"}, []string{uuid})
		return
	}

	wsSend(WS_SUCCESS_RESPONSE, WS_SUCCESS_RESPONSE_DTO{fmt.Sprint(http.StatusOK) + " Comment created"}, []string{uuid})

	// send all comments on that post
	wsCommentsListHandler(conn, messageData)
}

// wsCommentsListHandler returns all comments for a post
//
// - @param post_id int
func wsCommentsListHandler(conn *websocket.Conn, messageData map[string]interface{}) {
	defer wsRecover(messageData)

	uuid, ok := messageData["user_uuid"].(string)
	if !ok {
		log.Println("failed to get user_uuid from messageData")
		return
	}
	_, err := get_user_id_by_uuid(uuid)
	if err != nil {
		log.Println("failed to get ID of the message sender", err.Error())
		wsSend(WS_ERROR_RESPONSE, WS_ERROR_RESPONSE_DTO{fmt.Sprint(http.StatusUnprocessableEntity) + " failed to get ID of the message sender"}, []string{uuid})
		return
	}

	_post_id, ok := messageData["post_id"].(float64)
	if !ok {
		log.Println("failed to get post_id from messageData")
		wsSend(WS_ERROR_RESPONSE, WS_ERROR_RESPONSE_DTO{fmt.Sprint(http.StatusUnprocessableEntity, " failed to get post_id from messageData")}, []string{uuid})
		return
	}
	post_id := int(_post_id)

	rows, err := statements["getComments"].Query(post_id)
	if err != nil {
		log.Println("getComments query failed", err.Error())
		wsSend(WS_ERROR_RESPONSE, WS_ERROR_RESPONSE_DTO{fmt.Sprint(http.StatusInternalServerError) + " getComments query failed"}, []string{uuid})
		return
	}
	defer rows.Close()

	var commentsList WS_COMMENTS_LIST_DTO
	for rows.Next() {
		var comment WS_COMMENT_RESPONSE_DTO
		var pictureBlob []byte
		err = rows.Scan(&comment.Email, &comment.First_name, &comment.Last_name, &comment.Content, &pictureBlob, &comment.Created_at)
		if err != nil {
			log.Println("getComments scan failed", err.Error())
			wsSend(WS_ERROR_RESPONSE, WS_ERROR_RESPONSE_DTO{fmt.Sprint(http.StatusInternalServerError) + " getComments scan failed"}, []string{uuid})
			return
		}
		comment.Picture = base64.StdEncoding.EncodeToString(pictureBlob)
		commentsList = append(commentsList, comment)
	}

	wsSend(WS_COMMENTS_LIST, commentsList, []string{uuid})
}
