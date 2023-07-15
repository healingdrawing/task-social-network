package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

type WS_GROUP_CHAT_MESSAGE_DTO struct {
	Content    string `json:"content"`
	Email      string `json:"email"`
	First_name string `json:"first_name"`
	Last_name  string `json:"last_name"`
	Created_at string `json:"created_at"`
	Group_id   int    `json:"group_id"`
}

func wsGroupChatMessageHandler(conn *websocket.Conn, messageData map[string]interface{}) {
	defer wsRecover(messageData)

	uuid, ok := messageData["user_uuid"].(string)
	if !ok {
		log.Println("failed to get user_uuid from messageData")
		return
	}

	_group_id, ok := messageData["group_id"].(float64)
	if !ok {
		log.Println("failed to get group_id from messageData")
		wsSend(WS_ERROR_RESPONSE, WS_ERROR_RESPONSE_DTO{fmt.Sprint(http.StatusUnprocessableEntity, " failed to get group_id from messageData")}, []string{uuid})
		return
	}
	group_id := int(_group_id)

	content, ok := messageData["content"].(string)
	if !ok {
		log.Println("failed to get content from messageData")
		wsSend(WS_ERROR_RESPONSE, WS_ERROR_RESPONSE_DTO{fmt.Sprint(http.StatusUnprocessableEntity, " failed to get content from messageData")}, []string{uuid})
		return
	}

	if strings.TrimSpace(content) == "" {
		log.Println("content is empty")
		wsSend(WS_ERROR_RESPONSE, WS_ERROR_RESPONSE_DTO{fmt.Sprint(http.StatusUnprocessableEntity, " content is empty")}, []string{uuid})
		return
	}

	created_at := time.Now().Format("2006-01-02 15:04:05")

	var message WS_GROUP_CHAT_MESSAGE_DTO

	message.Content = content
	message.Created_at = created_at
	message.Group_id = group_id

	// Get user info
	user_id, err := get_user_id_by_uuid(uuid)
	if err != nil {
		log.Println("failed to get user_id", err.Error())
		wsSend(WS_ERROR_RESPONSE, WS_ERROR_RESPONSE_DTO{fmt.Sprint(http.StatusInternalServerError, " failed to get user_id")}, []string{uuid})
		return
	}

	var gap string
	rows, err := statements["getUserbyID"].Query(user_id)
	if err != nil {
		log.Println("failed to get user info", err.Error())
		wsSend(WS_ERROR_RESPONSE, WS_ERROR_RESPONSE_DTO{fmt.Sprint(http.StatusInternalServerError, " failed to get user info")}, []string{uuid})
		return
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&message.Email, &message.First_name, &message.Last_name, &gap)
		if err != nil {
			log.Println("failed to scan user info", err.Error())
			wsSend(WS_ERROR_RESPONSE, WS_ERROR_RESPONSE_DTO{fmt.Sprint(http.StatusInternalServerError, " failed to scan user info")}, []string{uuid})
			return
		}
	}

	//get all id's of users who is the member of the group of this group chat
	var group_member_ids []int
	rows, err = statements["getGroupMembers"].Query(group_id)
	if err != nil {
		log.Println("failed to get group members", err.Error())
		wsSend(WS_ERROR_RESPONSE, WS_ERROR_RESPONSE_DTO{fmt.Sprint(http.StatusInternalServerError, " failed to get group members")}, []string{uuid})
		return
	}
	defer rows.Close()
	for rows.Next() {
		var member_id int
		err := rows.Scan(&member_id)
		if err != nil {
			log.Println("failed to scan group members", err.Error())
			wsSend(WS_ERROR_RESPONSE, WS_ERROR_RESPONSE_DTO{fmt.Sprint(http.StatusInternalServerError, " failed to scan group members")}, []string{uuid})
			return
		}
		group_member_ids = append(group_member_ids, member_id)
	}

	// get connected user ids

	// get ids of users who is member and connected
	connected_group_member_ids := map[int]int{}
	clients.Range(func(key, value interface{}) bool {
		client := value.(*Client)
		for _, member_id := range group_member_ids {
			if client.USER_ID == member_id {
				connected_group_member_ids[member_id] = member_id
				break
			}
		}
		return true
	})

	//get all uuids of users who is connected/logged in and member of the group
	var user_uuids []string
	clients.Range(func(key, value interface{}) bool {
		client := value.(*Client)
		_, ok := connected_group_member_ids[client.USER_ID]
		if ok {
			user_uuids = append(user_uuids, key.(string))
		} else {
			log.Println("user is not connected/logged in anymore")
		}
		return true
	})

	wsSend(WS_GROUP_CHAT_MESSAGE, message, user_uuids)

}

type WS_PRIVATE_CHAT_MESSAGE_DTO struct {
	Content        string `json:"content"`
	Email          string `json:"email"`
	First_name     string `json:"first_name"`
	Last_name      string `json:"last_name"`
	Created_at     string `json:"created_at"`
	Target_user_id int    `json:"target_user_id"`
}

func wsPrivateChatMessageHandler(conn *websocket.Conn, messageData map[string]interface{}) {
	defer wsRecover(messageData)

	uuid, ok := messageData["user_uuid"].(string)
	if !ok {
		log.Println("failed to get user_uuid from messageData")
		return
	}

	_target_user_id, ok := messageData["target_user_id"].(float64)
	if !ok {
		log.Println("failed to get target_user_id from messageData")
		wsSend(WS_ERROR_RESPONSE, WS_ERROR_RESPONSE_DTO{fmt.Sprint(http.StatusUnprocessableEntity, " failed to get target_user_id from messageData")}, []string{uuid})
		return
	}
	target_user_id := int(_target_user_id)

	content, ok := messageData["content"].(string)
	if !ok {
		log.Println("failed to get content from messageData")
		wsSend(WS_ERROR_RESPONSE, WS_ERROR_RESPONSE_DTO{fmt.Sprint(http.StatusUnprocessableEntity, " failed to get content from messageData")}, []string{uuid})
		return
	}

	if strings.TrimSpace(content) == "" {
		log.Println("content is empty")
		wsSend(WS_ERROR_RESPONSE, WS_ERROR_RESPONSE_DTO{fmt.Sprint(http.StatusUnprocessableEntity, " content is empty")}, []string{uuid})
		return
	}

	created_at := time.Now().Format("2006-01-02 15:04:05")

	var message WS_PRIVATE_CHAT_MESSAGE_DTO

	message.Content = content
	message.Created_at = created_at
	message.Target_user_id = target_user_id

	// Get user info
	user_id, err := get_user_id_by_uuid(uuid)
	if err != nil {
		log.Println("failed to get user_id", err.Error())
		wsSend(WS_ERROR_RESPONSE, WS_ERROR_RESPONSE_DTO{fmt.Sprint(http.StatusInternalServerError, " failed to get user_id")}, []string{uuid})
		return
	}

	var gap string
	rows, err := statements["getUserbyID"].Query(user_id)
	if err != nil {
		log.Println("failed to get user info", err.Error())
		wsSend(WS_ERROR_RESPONSE, WS_ERROR_RESPONSE_DTO{fmt.Sprint(http.StatusInternalServerError, " failed to get user info")}, []string{uuid})
		return
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&message.Email, &message.First_name, &message.Last_name, &gap)
		if err != nil {
			log.Println("failed to scan user info", err.Error())
			wsSend(WS_ERROR_RESPONSE, WS_ERROR_RESPONSE_DTO{fmt.Sprint(http.StatusInternalServerError, " failed to scan user info")}, []string{uuid})
			return
		}
	}

	// get ids of users (EXPECTED TWO USERS OR LESS) who is connected and target_user_id or user_id
	connected_ids := map[int]int{}
	clients.Range(func(key, value interface{}) bool {
		client := value.(*Client)

		if client.USER_ID == target_user_id ||
			client.USER_ID == user_id {
			connected_ids[client.USER_ID] = client.USER_ID
		}

		return true
	})

	//get all uuids of users who is connected/logged in and who is target_user_id
	var user_uuids []string
	clients.Range(func(key, value interface{}) bool {
		client := value.(*Client)
		_, ok := connected_ids[client.USER_ID]
		if ok {
			user_uuids = append(user_uuids, key.(string))
		} else {
			log.Println("user is not connected/logged in anymore")
		}
		return true
	})

	if len(user_uuids) > 2 {
		log.Println("there are more than two receivers, which is not expected")
		wsSend(WS_ERROR_RESPONSE, WS_ERROR_RESPONSE_DTO{fmt.Sprint(http.StatusInternalServerError, " there are more than two receivers, which is not expected")}, []string{uuid})
	}

	wsSend(WS_PRIVATE_CHAT_MESSAGE, message, user_uuids)

}
