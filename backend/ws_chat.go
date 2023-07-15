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
	Content    string
	Email      string
	First_name string
	Last_name  string
	Created_at string
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

	//get all uuids of users who is connected/logged in
	var connected_user_uuids []string
	clients.Range(func(key, value interface{}) bool {
		connected_user_uuids = append(connected_user_uuids, key.(string))
		return true
	})

	// get all ids of connected users
	query := fmt.Sprintf("SELECT user_id FROM session WHERE uuid IN (%s)", strings.Join(connected_user_uuids, ","))
	rows, err = db.Query(query)
	if err != nil {
		log.Println("failed to get connected users", err.Error())
		wsSend(WS_ERROR_RESPONSE, WS_ERROR_RESPONSE_DTO{fmt.Sprint(http.StatusInternalServerError, " failed to get connected users")}, []string{uuid})
		return
	}
	defer rows.Close()
	connected_user_ids := map[string]int{}
	for rows.Next() {
		var user_id int
		err := rows.Scan(&user_id)
		if err != nil {
			log.Println("failed to scan connected users", err.Error())
			wsSend(WS_ERROR_RESPONSE, WS_ERROR_RESPONSE_DTO{fmt.Sprint(http.StatusInternalServerError, " failed to scan connected users")}, []string{uuid})
			return
		}
		connected_user_ids = append(connected_user_ids, user_id)
	}

	// keep only user ids who is the member and connected
	var user_ids []int
	for _, member_id := range group_member_ids {
		for _, connected_user_id := range connected_user_ids {
			if member_id == connected_user_id {
				user_ids = append(user_ids, member_id)
			}
		}
	}

	// get uuids of users who is the member and connected
	var user_uuids []string
	clients.Range(func(key, value interface{}) bool {
		return true
	})

	wsSend(WS_GROUP_CHAT_MESSAGE, message, user_uuids)

}
