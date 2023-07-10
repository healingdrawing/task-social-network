package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

// type GroupMemberRequest struct {
// 	GroupID int    `json:"group_id"`
// 	Email   string `json:"member_email"`
// }

// type UserFollowerRequest struct {
// 	Email string `json:"email"` // the user who wants to follow you
// }

// wsGroupRequestAcceptHandler is the handler for accepting a group request
//
// @rparam {user_uuid string, group_id int, requester_email string}
func wsGroupRequestAcceptHandler(conn *websocket.Conn, messageData map[string]interface{}) {
	defer wsRecover()

	uuid, ok := messageData["user_uuid"].(string)
	if !ok {
		log.Println("failed to get user_uuid from message data")
		wsSendError(WS_ERROR_RESPONSE_DTO{fmt.Sprint(http.StatusUnprocessableEntity) + " failed to get user_uuid from message data"})
		return
	}

	group_creator_id, err := get_user_id_by_uuid(uuid)
	if err != nil {
		log.Println("failed to get ID of the accept request sender", err.Error())
		wsSendError(WS_ERROR_RESPONSE_DTO{fmt.Sprint(http.StatusUnprocessableEntity) + " failed to get ID of the accept request sender"})
		return
	}

	_group_id, ok := messageData["group_id"].(float64)
	if !ok {
		log.Println("failed to get group_id from message data")
		wsSendError(WS_ERROR_RESPONSE_DTO{fmt.Sprint(http.StatusUnprocessableEntity) + " failed to get group_id from message data"})
		return
	}
	group_id := int64(_group_id)

	requester_email, ok := messageData["requester_email"].(string)
	if !ok {
		log.Println("failed to get requester_email from message data")
		wsSendError(WS_ERROR_RESPONSE_DTO{fmt.Sprint(http.StatusUnprocessableEntity) + " failed to get requester_email from message data"})
		return
	}

	requester_id, err := get_user_id_by_email(requester_email)
	if err != nil {
		log.Println("failed to get ID of the message sender", err.Error())
		wsSendError(WS_ERROR_RESPONSE_DTO{fmt.Sprint(http.StatusUnprocessableEntity) + " failed to get ID of the message sender"})
		return
	}

	rows, err := statements["getGroup"].Query(group_id)
	if err != nil {
		log.Println("getGroup query failed", err.Error())
		wsSendError(WS_ERROR_RESPONSE_DTO{fmt.Sprint(http.StatusInternalServerError) + " getGroup query failed"})
		return
	}
	defer rows.Close()

	var group WS_GROUP_CHECK_DTO

	for rows.Next() {
		err = rows.Scan(&group.Id, &group.Name, &group.Description, &group.Creator_id, &group.Created_at)
		if err != nil {
			log.Println("failed to scan group", err.Error())
			wsSendError(WS_ERROR_RESPONSE_DTO{fmt.Sprint(http.StatusInternalServerError) + " failed to scan group"})
			return
		}
	}

	if group.Creator_id != group_creator_id {
		log.Println("request sender is not the creator of the group")
		wsSendError(WS_ERROR_RESPONSE_DTO{fmt.Sprint(http.StatusUnauthorized) + " request sender is not the creator of the group"})
		return
	}

	_, err = statements["addGroupMember"].Exec(group_id, requester_id)
	if err != nil {
		log.Println("addGroupMember query failed", err.Error())
		wsSendError(WS_ERROR_RESPONSE_DTO{fmt.Sprint(http.StatusInternalServerError) + " addGroupMember query failed"})
		return
	}

	_, err = statements["removeGroupPendingMember"].Exec(group_id, requester_id)
	if err != nil {
		log.Println("removeGroupPendingMember query failed", err.Error())
		wsSendError(WS_ERROR_RESPONSE_DTO{fmt.Sprint(http.StatusInternalServerError) + " removeGroupPendingMember query failed"})
		return
	}

	wsSendSuccess(WS_SUCCESS_RESPONSE_DTO{fmt.Sprint(http.StatusOK) + " Success: you approved the group membership"})
}

// wsGroupRequestRejectHandler is the handler for rejecting a group request
//
// @rparam {user_uuid string, group_id int, requester_email string}
func wsGroupRequestRejectHandler(conn *websocket.Conn, messageData map[string]interface{}) {
	defer wsRecover()

	uuid, ok := messageData["user_uuid"].(string)
	if !ok {
		log.Println("failed to get user_uuid from message data")
		wsSendError(WS_ERROR_RESPONSE_DTO{fmt.Sprint(http.StatusUnprocessableEntity) + " failed to get user_uuid from message data"})
		return
	}

	group_creator_id, err := get_user_id_by_uuid(uuid)
	if err != nil {
		log.Println("failed to get ID of the accept request sender", err.Error())
		wsSendError(WS_ERROR_RESPONSE_DTO{fmt.Sprint(http.StatusUnprocessableEntity) + " failed to get ID of the accept request sender"})
		return
	}

	_group_id, ok := messageData["group_id"].(float64)
	if !ok {
		log.Println("failed to get group_id from message data")
		wsSendError(WS_ERROR_RESPONSE_DTO{fmt.Sprint(http.StatusUnprocessableEntity) + " failed to get group_id from message data"})
		return
	}
	group_id := int64(_group_id)

	requester_email, ok := messageData["requester_email"].(string)
	if !ok {
		log.Println("failed to get requester_email from message data")
		wsSendError(WS_ERROR_RESPONSE_DTO{fmt.Sprint(http.StatusUnprocessableEntity) + " failed to get requester_email from message data"})
		return
	}

	requester_id, err := get_user_id_by_email(requester_email)
	if err != nil {
		log.Println("failed to get ID of the message sender", err.Error())
		wsSendError(WS_ERROR_RESPONSE_DTO{fmt.Sprint(http.StatusUnprocessableEntity) + " failed to get ID of the message sender"})
		return
	}

	rows, err := statements["getGroup"].Query(group_id)
	if err != nil {
		log.Println("getGroup query failed", err.Error())
		wsSendError(WS_ERROR_RESPONSE_DTO{fmt.Sprint(http.StatusInternalServerError) + " getGroup query failed"})
		return
	}
	defer rows.Close()

	var group WS_GROUP_CHECK_DTO

	for rows.Next() {
		err = rows.Scan(&group.Id, &group.Name, &group.Description, &group.Creator_id, &group.Created_at)
		if err != nil {
			log.Println("failed to scan group", err.Error())
			wsSendError(WS_ERROR_RESPONSE_DTO{fmt.Sprint(http.StatusInternalServerError) + " failed to scan group"})
			return
		}
	}

	if group.Creator_id != group_creator_id {
		log.Println("request sender is not the creator of the group")
		wsSendError(WS_ERROR_RESPONSE_DTO{fmt.Sprint(http.StatusUnauthorized) + " request sender is not the creator of the group"})
		return
	}

	_, err = statements["removeGroupPendingMember"].Exec(group_id, requester_id)
	if err != nil {
		log.Println("removeGroupPendingMember query failed", err.Error())
		wsSendError(WS_ERROR_RESPONSE_DTO{fmt.Sprint(http.StatusInternalServerError) + " removeGroupPendingMember query failed"})
		return
	}

	wsSendSuccess(WS_SUCCESS_RESPONSE_DTO{fmt.Sprint(http.StatusOK) + " Success: you rejected the group membership"})
}

// wsAcceptFollowerHandler is the handler for approving a follower request
//
// @rparam {email string}
func wsAcceptFollowerHandler(conn *websocket.Conn, messageData map[string]interface{}) {
	defer wsRecover()

	uuid, ok := messageData["user_uuid"].(string)
	if !ok {
		log.Println("failed to get user_uuid from messageData")
		wsSendError(WS_ERROR_RESPONSE_DTO{fmt.Sprint(http.StatusUnprocessableEntity) + " failed to get user_uuid from messageData"})
		return
	}

	idol_user_id, err := get_user_id_by_uuid(uuid)
	if err != nil {
		log.Println("failed to get ID of the message sender", err.Error())
		wsSendError(WS_ERROR_RESPONSE_DTO{fmt.Sprint(http.StatusUnprocessableEntity) + " failed to get ID of the message sender"})
		return
	}

	target_email, ok := messageData["target_email"].(string)
	if !ok {
		log.Println("failed to get target_email from messageData")
		wsSendError(WS_ERROR_RESPONSE_DTO{fmt.Sprint(http.StatusUnprocessableEntity) + " failed to get target_email from messageData"})
		return
	}

	fan_user_id, err := get_user_id_by_email(target_email)
	if err != nil {
		log.Println("failed to get ID of the follow request creator", err.Error())
		wsSendError(WS_ERROR_RESPONSE_DTO{fmt.Sprint(http.StatusUnprocessableEntity) + " failed to get ID of the follow request creator"})
		return
	}

	// check if the user has the follower in the followers_pending table
	// use the getFollowersPending statement
	rows, err := statements["getFollowersPending"].Query(idol_user_id)
	defer rows.Close()
	if err != nil {
		log.Println("getFollowersPending query failed", err.Error())
		wsSendError(WS_ERROR_RESPONSE_DTO{fmt.Sprint(http.StatusInternalServerError) + " getFollowersPending query failed"})
		return
	}
	var follower_id int
	followers_pending := map[int]int{}
	for rows.Next() {
		err = rows.Scan(&follower_id)
		if err != nil {
			log.Println("failed to scan followers pending", err.Error())
			wsSendError(WS_ERROR_RESPONSE_DTO{fmt.Sprint(http.StatusInternalServerError) + " failed to scan followers pending"})
			return
		}
		followers_pending[follower_id] = follower_id
	}
	// check if the user is in the list of followersPending request
	_, ok = followers_pending[fan_user_id]

	if !ok {
		log.Println("the user is not in the list of followers pending")
		wsSendError(WS_ERROR_RESPONSE_DTO{fmt.Sprint(http.StatusBadRequest) + " the user is not in the list of followers pending"})
		return
	}

	// add the follower to the followers table
	_, err = statements["addFollower"].Exec(idol_user_id, fan_user_id)
	if err != nil {
		log.Println("addFollower query failed", err.Error())
		wsSendError(WS_ERROR_RESPONSE_DTO{fmt.Sprint(http.StatusInternalServerError) + " addFollower query failed"})
		return
	}

	// remove it from the followers_pending table
	_, err = statements["removeFollowerPending"].Exec(idol_user_id, fan_user_id)
	if err != nil {
		log.Println("removeFollowerPending query failed", err.Error())
		wsSendError(WS_ERROR_RESPONSE_DTO{fmt.Sprint(http.StatusInternalServerError) + " removeFollowerPending query failed"})
		return
	}
	wsSendSuccess(WS_SUCCESS_RESPONSE_DTO{fmt.Sprint(http.StatusOK) + " Success: you approved the follow request"})
}

// wsRejectFollowerHandler is the handler for rejecting a follower request
//
// @rparam invited_emails (string space separated}
func wsRejectFollowerHandler(conn *websocket.Conn, messageData map[string]interface{}) {
	defer wsRecover()

	uuid, ok := messageData["user_uuid"].(string)
	if !ok {
		log.Println("failed to get user_uuid from messageData")
		wsSendError(WS_ERROR_RESPONSE_DTO{fmt.Sprint(http.StatusUnprocessableEntity) + " failed to get user_uuid from messageData"})
		return
	}

	idol_user_id, err := get_user_id_by_uuid(uuid)
	if err != nil {
		log.Println("failed to get ID of the message sender", err.Error())
		wsSendError(WS_ERROR_RESPONSE_DTO{fmt.Sprint(http.StatusUnprocessableEntity) + " failed to get ID of the message sender"})
		return
	}

	target_email, ok := messageData["target_email"].(string)
	if !ok {
		log.Println("failed to get target_email from messageData")
		wsSendError(WS_ERROR_RESPONSE_DTO{fmt.Sprint(http.StatusUnprocessableEntity) + " failed to get target_email from messageData"})
		return
	}

	fan_user_id, err := get_user_id_by_email(target_email)
	if err != nil {
		log.Println("failed to get ID of the follow request creator", err.Error())
		wsSendError(WS_ERROR_RESPONSE_DTO{fmt.Sprint(http.StatusUnprocessableEntity) + " failed to get ID of the follow request creator"})
		return
	}

	// check if the user has the follower in the followers_pending table
	// use the getFollowersPending statement
	rows, err := statements["getFollowersPending"].Query(idol_user_id)
	defer rows.Close()
	if err != nil {
		log.Println("getFollowersPending query failed", err.Error())
		wsSendError(WS_ERROR_RESPONSE_DTO{fmt.Sprint(http.StatusInternalServerError) + " getFollowersPending query failed"})
		return
	}
	var follower_id int
	followers_pending := map[int]int{}
	for rows.Next() {
		err = rows.Scan(&follower_id)
		if err != nil {
			log.Println("failed to scan followers pending", err.Error())
			wsSendError(WS_ERROR_RESPONSE_DTO{fmt.Sprint(http.StatusInternalServerError) + " failed to scan followers pending"})
			return
		}
		followers_pending[follower_id] = follower_id
	}
	// check if the user is in the list of followersPending request
	_, ok = followers_pending[fan_user_id]

	if !ok {
		log.Println("the user is not in the list of followers pending")
		wsSendError(WS_ERROR_RESPONSE_DTO{fmt.Sprint(http.StatusBadRequest) + " the user is not in the list of followers pending"})
		return
	}

	// remove it from the followers_pending table
	_, err = statements["removeFollowerPending"].Exec(idol_user_id, fan_user_id)
	if err != nil {
		log.Println("removeFollowerPending query failed", err.Error())
		wsSendError(WS_ERROR_RESPONSE_DTO{fmt.Sprint(http.StatusInternalServerError) + " removeFollowerPending query failed"})
		return
	}

	wsSendSuccess(WS_SUCCESS_RESPONSE_DTO{fmt.Sprint(http.StatusOK) + " Success: you rejected the follow request"})
}

// wsGroupInviteAcceptHandler is the handler for accepting a group invite
//
// @r.param {group_id int}
func wsGroupInviteAcceptHandler(conn *websocket.Conn, messageData map[string]interface{}) {
	defer wsRecover()

	log.Println("=== wsGroupInviteAcceptHandler")

	uuid, ok := messageData["user_uuid"].(string)
	if !ok {
		log.Println("failed to get user_uuid from message data")
		wsSendError(WS_ERROR_RESPONSE_DTO{fmt.Sprint(http.StatusUnprocessableEntity) + " failed to get user_uuid from message data"})
		return
	}

	user_id, err := get_user_id_by_uuid(uuid)
	if err != nil {
		log.Println("failed to get ID of the accept request sender", err.Error())
		wsSendError(WS_ERROR_RESPONSE_DTO{fmt.Sprint(http.StatusUnprocessableEntity) + " failed to get ID of the accept request sender"})
		return
	}

	_group_id, ok := messageData["group_id"].(float64)
	if !ok {
		log.Println("failed to get group_id from message data")
		wsSendError(WS_ERROR_RESPONSE_DTO{fmt.Sprint(http.StatusUnprocessableEntity) + " failed to get group_id from message data"})
		return
	}
	group_id := int64(_group_id)

	// check if the user was invited to the group, using getGroupInvitedUsers statement
	rows, err := statements["getGroupInvitedUsers"].Query(group_id)
	defer rows.Close()
	if err != nil {
		log.Println("getGroupInvitedUsers query failed", err.Error())
		wsSendError(WS_ERROR_RESPONSE_DTO{fmt.Sprint(http.StatusInternalServerError) + " getGroupInvitedUsers query failed"})
		return
	}
	var invited_user_id int
	var gapint int
	var gapstring string
	invited_user_ids := map[int]int{}
	for rows.Next() {
		err = rows.Scan(&invited_user_id, &gapint, &gapint, &gapstring)
		if err != nil {
			log.Println("failed to scan invited users", err.Error())
			wsSendError(WS_ERROR_RESPONSE_DTO{fmt.Sprint(http.StatusInternalServerError) + " failed to scan invited users"})
			return
		}
		invited_user_ids[invited_user_id] = invited_user_id
	}
	// check if the user is in the list of invited users
	_, ok = invited_user_ids[user_id]

	if !ok {
		log.Println("the user is not in the list of invited users")
		wsSendError(WS_ERROR_RESPONSE_DTO{fmt.Sprint(http.StatusBadRequest) + " the user is not in the list of invited users"})
		return
	}

	_, err = statements["addGroupMember"].Exec(group_id, user_id)
	if err != nil {
		log.Println("addGroupMember query failed", err.Error())
		wsSendError(WS_ERROR_RESPONSE_DTO{fmt.Sprint(http.StatusInternalServerError) + " addGroupMember query failed"})
		return
	}

	_, err = statements["removeGroupInvitedUser"].Exec(user_id, group_id)
	if err != nil {
		log.Println("removeGroupInvitedUser query failed", err.Error())
		wsSendError(WS_ERROR_RESPONSE_DTO{fmt.Sprint(http.StatusInternalServerError) + " removeGroupInvitedUser query failed"})
		return
	}

	wsSendSuccess(WS_SUCCESS_RESPONSE_DTO{fmt.Sprint(http.StatusOK) + " Success: you accepted the group invite"})
}

// wsGroupInviteRejectHandler is the handler for rejecting a group invite
//
// @param {group_id int}
func wsGroupInviteRejectHandler(conn *websocket.Conn, messageData map[string]interface{}) {
	defer wsRecover()

	log.Println("=== wsGroupInviteRejectHandler")

	uuid, ok := messageData["user_uuid"].(string)
	if !ok {
		log.Println("failed to get user_uuid from message data")
		wsSendError(WS_ERROR_RESPONSE_DTO{fmt.Sprint(http.StatusUnprocessableEntity) + " failed to get user_uuid from message data"})
		return
	}

	user_id, err := get_user_id_by_uuid(uuid)
	if err != nil {
		log.Println("failed to get ID of the accept request sender", err.Error())
		wsSendError(WS_ERROR_RESPONSE_DTO{fmt.Sprint(http.StatusUnprocessableEntity) + " failed to get ID of the accept request sender"})
		return
	}

	_group_id, ok := messageData["group_id"].(float64)
	if !ok {
		log.Println("failed to get group_id from message data")
		wsSendError(WS_ERROR_RESPONSE_DTO{fmt.Sprint(http.StatusUnprocessableEntity) + " failed to get group_id from message data"})
		return
	}
	group_id := int64(_group_id)

	// check if the user was invited to the group, using getGroupInvitedUsers statement
	rows, err := statements["getGroupInvitedUsers"].Query(group_id)
	defer rows.Close()
	if err != nil {
		log.Println("getGroupInvitedUsers query failed", err.Error())
		wsSendError(WS_ERROR_RESPONSE_DTO{fmt.Sprint(http.StatusInternalServerError) + " getGroupInvitedUsers query failed"})
		return
	}
	var invited_user_id int
	var gapint int
	var gapstring string
	invited_user_ids := map[int]int{}
	for rows.Next() {
		err = rows.Scan(&invited_user_id, &gapint, &gapint, &gapstring)
		if err != nil {
			log.Println("failed to scan invited users", err.Error())
			wsSendError(WS_ERROR_RESPONSE_DTO{fmt.Sprint(http.StatusInternalServerError) + " failed to scan invited users"})
			return
		}
		invited_user_ids[invited_user_id] = invited_user_id
	}
	// check if the user is in the list of invited users
	_, ok = invited_user_ids[user_id]

	if !ok {
		log.Println("the user is not in the list of invited users")
		wsSendError(WS_ERROR_RESPONSE_DTO{fmt.Sprint(http.StatusBadRequest) + " the user is not in the list of invited users"})
		return
	}

	_, err = statements["removeGroupInvitedUser"].Exec(user_id, group_id)
	if err != nil {
		log.Println("removeGroupInvitedUser query failed", err.Error())
		wsSendError(WS_ERROR_RESPONSE_DTO{fmt.Sprint(http.StatusInternalServerError) + " removeGroupInvitedUser query failed"})
		return
	}

	wsSendSuccess(WS_SUCCESS_RESPONSE_DTO{fmt.Sprint(http.StatusOK) + " Success: you rejected the group invite"})
}
