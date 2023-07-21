package main

import (
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

type WS_GROUP_POST_SUBMIT_DTO struct {
	User_uuid  string `json:"user_uuid"`
	Group_id   int    `json:"group_id"`
	Title      string `json:"title"`
	Categories string `json:"categories"`
	Content    string `json:"content"`
	Picture    string `json:"picture"`
	Created_at string `json:"created_at"`
}

type WS_GROUP_POST_RESPONSE_DTO struct {
	Group_id          int    `json:"group_id"` // use to open group from profile
	Group_name        string `json:"group_name"`
	Group_description string `json:"group_description"`
	Id                int    `json:"id"` // post section. use to get comments
	Title             string `json:"title"`
	Content           string `json:"content"`
	Categories        string `json:"categories"`
	Picture           string `json:"picture"`
	Created_at        string `json:"created_at"`
	Email             string `json:"email"` //creator section. use to show post credentials
	First_name        string `json:"first_name"`
	Last_name         string `json:"last_name"`
}

type WS_GROUP_POSTS_LIST_DTO []WS_GROUP_POST_RESPONSE_DTO

// wsGroupPostSubmitHandler creates a new group post, then return all posts of group
func wsGroupPostSubmitHandler(conn *websocket.Conn, messageData map[string]interface{}) {
	defer wsRecover(messageData)

	uuid, ok := messageData["user_uuid"].(string)
	if !ok {
		log.Println("failed to get user_uuid from message data")
		return
	}

	_group_id, ok := messageData["group_id"].(float64)
	if !ok {
		log.Println("failed to get group_id from messageData")
		wsSend(WS_ERROR_RESPONSE, WS_ERROR_RESPONSE_DTO{fmt.Sprint(http.StatusUnprocessableEntity) + " failed to get group_id from messageData"}, []string{uuid})
		return
	}
	group_id := int(_group_id)

	var data WS_GROUP_POST_SUBMIT_DTO
	data.User_uuid = uuid

	fields := map[string]*string{
		"title":      &data.Title,
		"categories": &data.Categories,
		"content":    &data.Content,
		"picture":    &data.Picture,
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

	data.Created_at = time.Now().Format("2006-01-02 15:04:05")
	result, err := statements["addGroupPost"].Exec(user_id, data.Title, data.Categories, data.Content, postPicture, data.Created_at)
	if err != nil {
		log.Println("addGroupPost query failed", err.Error())
		wsSend(WS_ERROR_RESPONSE, WS_ERROR_RESPONSE_DTO{fmt.Sprint(http.StatusInternalServerError) + " addGroupPost query failed"}, []string{uuid})
		return
	}

	group_post_id, err := result.LastInsertId()
	if err != nil {
		log.Println("LastInsertId of addGroupPost query failed", err.Error())
		wsSend(WS_ERROR_RESPONSE, WS_ERROR_RESPONSE_DTO{fmt.Sprint(http.StatusInternalServerError) + " LastInsertId of addGroupPost query failed"}, []string{uuid})
		return
	}

	// add post to group
	_, err = statements["addGroupPostMembership"].Exec(group_id, group_post_id)
	if err != nil {
		log.Println("addGroupPostMembership query failed", err.Error())
		wsSend(WS_ERROR_RESPONSE, WS_ERROR_RESPONSE_DTO{fmt.Sprint(http.StatusInternalServerError) + " addGroupPostMembership query failed"}, []string{uuid})
		return
	}

	wsSend(WS_SUCCESS_RESPONSE, WS_SUCCESS_RESPONSE_DTO{fmt.Sprint(http.StatusOK) + " group post created"}, []string{uuid})

	// return all group posts
	wsGroupPostsListHandler(conn, messageData)
}

func wsGroupPostsListHandler(conn *websocket.Conn, messageData map[string]interface{}) {
	defer wsRecover(messageData)

	uuid, ok := messageData["user_uuid"].(string)
	if !ok {
		log.Println("failed to get user_uuid from message data")
		return
	}

	_, err := get_user_id_by_uuid(uuid)
	if err != nil {
		log.Println("failed to get ID of the message sender", err.Error())
		wsSend(WS_ERROR_RESPONSE, WS_ERROR_RESPONSE_DTO{fmt.Sprint(http.StatusUnprocessableEntity) + " failed to get ID of the message sender"}, []string{uuid})
		return
	}

	_group_id, ok := messageData["group_id"].(float64)
	if !ok {
		log.Println("failed to get group_id from messageData")
		wsSend(WS_ERROR_RESPONSE, WS_ERROR_RESPONSE_DTO{fmt.Sprint(http.StatusUnprocessableEntity) + " failed to get group_id from messageData"}, []string{uuid})
		return
	}
	group_id := int(_group_id)

	rows, err := statements["getGroupPosts"].Query(group_id)
	if err != nil {
		log.Println("getGroupPosts query failed", err.Error())
		wsSend(WS_ERROR_RESPONSE, WS_ERROR_RESPONSE_DTO{fmt.Sprint(http.StatusInternalServerError) + " getGroupPosts query failed"}, []string{uuid})
		return
	}
	defer rows.Close()

	var group_posts_list WS_GROUP_POSTS_LIST_DTO
	for rows.Next() {
		var group_post WS_GROUP_POST_RESPONSE_DTO
		pictureBytes := []byte{}
		err = rows.Scan(
			&group_post.Group_id,
			&group_post.Group_name,
			&group_post.Group_description,
			&group_post.Id,
			&group_post.Title,
			&group_post.Content,
			&group_post.Categories,
			&pictureBytes,
			&group_post.Created_at,
			&group_post.Email,
			&group_post.First_name,
			&group_post.Last_name,
		)
		if err != nil {
			log.Println("group post scan failed", err.Error())
			wsSend(WS_ERROR_RESPONSE, WS_ERROR_RESPONSE_DTO{fmt.Sprint(http.StatusInternalServerError) + " group post scan failed"}, []string{uuid})
			return
		}
		group_post.Picture = base64.StdEncoding.EncodeToString(pictureBytes)
		group_posts_list = append(group_posts_list, group_post)
	}

	wsSend(WS_GROUP_POSTS_LIST, group_posts_list, []string{uuid})
}

/** includes the posts created by user, BUT NOT GROUP POSTS, created by user */
func wsUserGroupPostsListHandler(conn *websocket.Conn, messageData map[string]interface{}) {
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

	target_email := messageData["target_email"].(string)
	target_id, err := get_user_id_by_email(target_email)
	if err != nil {
		log.Println("failed to get ID of the target", err.Error())
		wsSend(WS_ERROR_RESPONSE, WS_ERROR_RESPONSE_DTO{fmt.Sprint(http.StatusUnprocessableEntity) + " failed to get ID of the target"}, []string{uuid})
		return
	}

	rows, err := statements["getUserAllGroupPosts"].Query(target_id)
	if err != nil {
		log.Println("getUserAllGroupPosts query failed", err.Error())
		wsSend(WS_ERROR_RESPONSE, WS_ERROR_RESPONSE_DTO{fmt.Sprint(http.StatusInternalServerError) + " getUserAllGroupPosts query failed"}, []string{uuid})
		return
	}
	defer rows.Close()

	var group_posts_list WS_GROUP_POSTS_LIST_DTO
	for rows.Next() {
		var group_post WS_GROUP_POST_RESPONSE_DTO
		pictureBytes := []byte{}
		err = rows.Scan(
			&group_post.Group_id,
			&group_post.Group_name,
			&group_post.Group_description,
			&group_post.Id,
			&group_post.Title,
			&group_post.Content,
			&group_post.Categories,
			&pictureBytes,
			&group_post.Created_at,
			&group_post.Email,
			&group_post.First_name,
			&group_post.Last_name,
		)
		if err != nil {
			log.Println("group post scan failed", err.Error())
			wsSend(WS_ERROR_RESPONSE, WS_ERROR_RESPONSE_DTO{fmt.Sprint(http.StatusInternalServerError) + " group post scan failed"}, []string{uuid})
			return
		}
		group_post.Picture = base64.StdEncoding.EncodeToString(pictureBytes)
		group_posts_list = append(group_posts_list, group_post)
	}

	wsSend(WS_GROUP_POSTS_LIST, group_posts_list, []string{uuid})
}
