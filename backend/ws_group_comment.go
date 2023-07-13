package main

import (
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

// the same implementation as "ws_comment.go", but different SQL queries

// type WS_COMMENT_SUBMIT_DTO struct {
// 	User_uuid string `json:"user_uuid"`
// 	Post_id   int    `json:"post_id"`
// 	Content   string `json:"content"`
// 	Picture   string `json:"picture"`
// }

// type WS_COMMENT_RESPONSE_DTO struct {
// 	Content    string `json:"content"`
// 	Picture    string `json:"picture"`
// 	Created_at string `json:"created_at"`
// 	Email      string `json:"email"`
// 	First_name string `json:"first_name"`
// 	Last_name  string `json:"last_name"`
// }

// type WS_COMMENTS_LIST_DTO []WS_COMMENT_RESPONSE_DTO

// wsGroupPostCommentSubmitHandler creates a new comment on a group post, then return all comments on group post.
func wsGroupPostCommentSubmitHandler(conn *websocket.Conn, messageData map[string]interface{}) {
	defer wsRecover()

	var data WS_COMMENT_SUBMIT_DTO

	_post_id, ok := messageData["post_id"].(float64)
	if !ok {
		log.Println("failed to get post_id from messageData")
		wsSendError(WS_ERROR_RESPONSE_DTO{fmt.Sprint(http.StatusUnprocessableEntity, " failed to get post_id from messageData")})
		return
	}
	data.Post_id = int(_post_id)

	fields := map[string]*string{
		"user_uuid": &data.User_uuid,
		"content":   &data.Content,
		"picture":   &data.Picture,
	}

	for key, ptr := range fields {
		value, ok := messageData[key].(string)
		if !ok {
			log.Printf("failed to get %s from messageData\n", key)
			wsSendError(WS_ERROR_RESPONSE_DTO{fmt.Sprintf("%d failed to get %s from messageData", http.StatusUnprocessableEntity, key)})
			return
		}
		*ptr = value
	}

	user_id, err := get_user_id_by_uuid(data.User_uuid)
	if err != nil {
		log.Println("failed to get ID of the message sender", err.Error())
		wsSendError(WS_ERROR_RESPONSE_DTO{fmt.Sprint(http.StatusUnprocessableEntity) + " failed to get ID of the message sender"})
		return
	}

	// process the picture
	commentPicture := []byte{}
	if data.Picture != "" {
		imageData, err := extractImageData(data.Picture)
		if err != nil {
			log.Println("=FAIL extractPictureData: ", err.Error())
			wsSendError(WS_ERROR_RESPONSE_DTO{fmt.Sprint(http.StatusUnprocessableEntity) + " =FAIL extractPictureData:" + err.Error()})
			return
		}
		pictureData, err := base64.StdEncoding.DecodeString(imageData)
		if err != nil {
			log.Println("Invalid picture ", err.Error())
			wsSendError(WS_ERROR_RESPONSE_DTO{fmt.Sprint(http.StatusUnprocessableEntity) + " Invalid picture"})
			return
		}
		if !isImage(pictureData) {
			log.Println("picture is not a valid image", err.Error())
			wsSendError(WS_ERROR_RESPONSE_DTO{fmt.Sprint(http.StatusUnsupportedMediaType) + " picture is not a valid image"})
			return
		}
		commentPicture = pictureData
	}

	created_at := time.Now().Format("2006-01-02 15:04:05")

	_, err = statements["addGroupComment"].Exec(user_id, data.Post_id, data.Content, commentPicture, created_at)
	if err != nil {
		log.Println("addGroupComment query failed", err.Error())
		wsSendError(WS_ERROR_RESPONSE_DTO{fmt.Sprint(http.StatusInternalServerError) + " addGroupComment query failed"})
		return
	}

	wsSendSuccess(WS_SUCCESS_RESPONSE_DTO{fmt.Sprint(http.StatusOK) + " Group Comment created"})

	// send all comments on group post
	// duplicate defer wsRecover(), do not want to crap the code using bool + if
	wsGroupPostCommentsListHandler(conn, messageData)
}

// wsGroupPostCommentsListHandler returns all comments for a group post
//
// - @param post_id int
func wsGroupPostCommentsListHandler(conn *websocket.Conn, messageData map[string]interface{}) {
	defer wsRecover()

	uuid, ok := messageData["user_uuid"].(string)
	if !ok {
		log.Println("failed to get user_uuid from messageData")
		wsSendError(WS_ERROR_RESPONSE_DTO{fmt.Sprint(http.StatusUnprocessableEntity, " failed to get user_uuid from messageData")})
		return
	}
	_, err := get_user_id_by_uuid(uuid)
	if err != nil {
		log.Println("failed to get ID of the message sender", err.Error())
		wsSendError(WS_ERROR_RESPONSE_DTO{fmt.Sprint(http.StatusUnprocessableEntity) + " failed to get ID of the message sender"})
		return
	}

	_post_id, ok := messageData["post_id"].(float64)
	if !ok {
		log.Println("failed to get post_id from messageData")
		wsSendError(WS_ERROR_RESPONSE_DTO{fmt.Sprint(http.StatusUnprocessableEntity, " failed to get post_id from messageData")})
		return
	}
	post_id := int(_post_id)

	rows, err := statements["getGroupComments"].Query(post_id)
	if err != nil {
		log.Println("getGroupComments query failed", err.Error())
		wsSendError(WS_ERROR_RESPONSE_DTO{fmt.Sprint(http.StatusInternalServerError) + " getGroupComments query failed"})
		return
	}
	defer rows.Close()

	var commentsList WS_COMMENTS_LIST_DTO
	for rows.Next() {
		var comment WS_COMMENT_RESPONSE_DTO
		var pictureBlob []byte
		err = rows.Scan(&comment.Email, &comment.First_name, &comment.Last_name, &comment.Content, &pictureBlob, &comment.Created_at)
		if err != nil {
			log.Println("getGroupComments scan failed", err.Error())
			wsSendError(WS_ERROR_RESPONSE_DTO{fmt.Sprint(http.StatusInternalServerError) + " getGroupComments scan failed"})
			return
		}
		comment.Picture = base64.StdEncoding.EncodeToString(pictureBlob)
		commentsList = append(commentsList, comment)
	}
	wsSendCommentsList(commentsList)
}
