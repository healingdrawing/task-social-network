package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

type WS_GROUP_SUBMIT_DTO struct {
	User_uuid      string `json:"user_uuid"`
	Name           string `json:"name"`
	Description    string `json:"description"`
	Invited_emails string `json:"invited_emails"` // space separated emails
}

type WS_GROUP_RESPONSE_DTO struct {
	Id          int       `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Created_at  time.Time `json:"created_at"`
	Email       string    `json:"email"`
	First_name  string    `json:"first_name"`
	Last_name   string    `json:"last_name"`
}

type WS_GROUPS_LIST_DTO []WS_GROUP_RESPONSE_DTO

// AS IN DB
type WS_GROUP_CHECK_DTO struct {
	Id          int       `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Creator_id  int       `json:"creator_id"`
	Created_at  time.Time `json:"created_at"`
}

type WS_GROUP_REQUEST_RESPONSE_DTO struct {
	Group_id    int    `json:"group_id"`  // accept or reject in frontend
	Member_id   int    `json:"member_id"` // accept or reject in frontend
	Group_name  string `json:"group_name"`
	Description string `json:"description"`
	Email       string `json:"email"`
	First_name  string `json:"first_name"`
	Last_name   string `json:"last_name"`
}

type WS_GROUP_REQUESTS_LIST_DTO []WS_GROUP_REQUEST_RESPONSE_DTO

type WS_GROUP_INVITE_RESPONSE_DTO struct {
	Group_id           int       `json:"group_id"`
	Group_name         string    `json:"group_name"`
	Group_description  string    `json:"group_description"`
	Created_at         time.Time `json:"created_at"`
	Inviter_email      string    `json:"email"`
	Inviter_first_name string    `json:"first_name"`
	Inviter_last_name  string    `json:"last_name"`
}
type WS_GROUP_INVITES_LIST_DTO []WS_GROUP_INVITE_RESPONSE_DTO

//old refactored types bottom

type WS_GROUP_MEMBER struct {
	Email    int `json:"email"`
	Group_id int `json:"group_id"`
}

type WS_GROUP_INVITED_MEMBER struct {
	Email    int `json:"email"`
	Group_id int `json:"group_id"`
}

type WS_INVITED_USER_INFO_DTO struct {
	Email      string `json:"email"`
	First_name string `json:"first_name"`
	Last_name  string `json:"last_name"`
}

func wsGroupSubmitHandler(conn *websocket.Conn, messageData map[string]interface{}) {
	defer wsRecover(messageData)

	uuid, ok := messageData["user_uuid"].(string)
	if !ok {
		log.Println("failed to get user_uuid from message data")
		return
	}

	var data WS_GROUP_SUBMIT_DTO
	data.User_uuid = uuid

	fields := map[string]*string{
		"name":           &data.Name,
		"description":    &data.Description,
		"invited_emails": &data.Invited_emails,
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

	// sanitize data
	data.Name = strings.TrimSpace(data.Name)
	data.Description = strings.TrimSpace(data.Description)
	if data.Name == "" || data.Description == "" {
		log.Println("empty fields")
		wsSend(WS_ERROR_RESPONSE, WS_ERROR_RESPONSE_DTO{fmt.Sprint(http.StatusUnprocessableEntity) + " empty fields"}, []string{uuid})
		return
	}
	created_at := time.Now().Format("2006-01-02 15:04:05")
	result, err := statements["addGroup"].Exec(data.Name, data.Description, user_id, created_at)
	if err != nil {
		log.Println("addGroup query failed", err.Error())
		wsSend(WS_ERROR_RESPONSE, WS_ERROR_RESPONSE_DTO{fmt.Sprint(http.StatusInternalServerError) + " addGroup query failed"}, []string{uuid})
		return
	}
	group_id, err := result.LastInsertId()
	if err != nil {
		log.Println("failed to get last insert id", err.Error())
		wsSend(WS_ERROR_RESPONSE, WS_ERROR_RESPONSE_DTO{fmt.Sprint(http.StatusInternalServerError) + " failed to get last insert id"}, []string{uuid})
		return
	}
	invited_emails_string := strings.TrimSpace(data.Invited_emails)
	if invited_emails_string != "" {
		invited_emails := strings.Split(invited_emails_string, " ")
		for _, email := range invited_emails {
			// get the user id from the email
			invited_user_id, err := get_user_id_by_email(email)
			if err != nil {
				log.Println("failed to get user id by email", err.Error())
				wsSend(WS_ERROR_RESPONSE, WS_ERROR_RESPONSE_DTO{fmt.Sprint(http.StatusInternalServerError) + " failed to get user id by email"}, []string{uuid})
				return
			}
			// add the invited user to the groupPendingMembers table
			created_at := time.Now().Format("2006-01-02 15:04:05")
			_, err = statements["addGroupInvitedUser"].Exec(invited_user_id, group_id, user_id, created_at)
			if err != nil {
				log.Println("addGroupInvitedUser query failed", err.Error())
				wsSend(WS_ERROR_RESPONSE, WS_ERROR_RESPONSE_DTO{fmt.Sprint(http.StatusInternalServerError) + " addGroupInvitedUser query failed"}, []string{uuid})
				return
			}
		}
	}
	// add group creator to group members
	_, err = statements["addGroupMember"].Exec(group_id, user_id)
	if err != nil {
		log.Println("addGroupMember query failed", err.Error())
		wsSend(WS_ERROR_RESPONSE, WS_ERROR_RESPONSE_DTO{fmt.Sprint(http.StatusInternalServerError) + " addGroupMember query failed"}, []string{uuid})
		return
	}

	// send response
	wsSend(WS_SUCCESS_RESPONSE, WS_SUCCESS_RESPONSE_DTO{fmt.Sprint(http.StatusOK) + " group created"}, []string{uuid})

	// send the list of all groups where user is a member, for GroupsView.vue
	wsGroupsListHandler(conn, messageData)

}

// wsGroupsListHandler returns list of groups where user is a member, for GroupsView.vue
func wsGroupsListHandler(conn *websocket.Conn, messageData map[string]interface{}) {
	uuid, ok := messageData["user_uuid"].(string)
	if !ok {
		log.Println("failed to get user_uuid from messageData")
		wsSend(WS_ERROR_RESPONSE, WS_ERROR_RESPONSE_DTO{fmt.Sprint(http.StatusUnprocessableEntity) + " failed to get user_uuid from messageData"}, []string{uuid})
		return
	}
	user_id, err := get_user_id_by_uuid(uuid)
	if err != nil {
		log.Println("failed to get ID of the message sender", err.Error())
		wsSend(WS_ERROR_RESPONSE, WS_ERROR_RESPONSE_DTO{fmt.Sprint(http.StatusUnprocessableEntity) + " failed to get ID of the message sender"}, []string{uuid})
		return
	}
	rows, err := statements["getGroups"].Query(user_id)
	if err != nil {
		log.Println("getGroups query failed", err.Error())
		wsSend(WS_ERROR_RESPONSE, WS_ERROR_RESPONSE_DTO{fmt.Sprint(http.StatusInternalServerError) + " getGroups query failed"}, []string{uuid})
		return
	}
	defer rows.Close()
	var groups WS_GROUPS_LIST_DTO
	for rows.Next() {
		var group WS_GROUP_RESPONSE_DTO
		err = rows.Scan(&group.Id, &group.Name, &group.Description, &group.Created_at, &group.Email, &group.First_name, &group.Last_name)
		if err != nil {
			log.Println("getGroups query failed to scan", err.Error())
			wsSend(WS_ERROR_RESPONSE, WS_ERROR_RESPONSE_DTO{fmt.Sprint(http.StatusInternalServerError) + " getGroups query failed to scan row"}, []string{uuid})
			return
		}
		groups = append(groups, group)
	}

	wsSend(WS_GROUPS_LIST, groups, []string{uuid})
}

// wsGroupsAllListHandler returns list of all groups, for GroupsAllView.vue
func wsGroupsAllListHandler(conn *websocket.Conn, messageData map[string]interface{}) {
	uuid, ok := messageData["user_uuid"].(string)
	if !ok {
		log.Println("failed to get user_uuid from messageData")
		wsSend(WS_ERROR_RESPONSE, WS_ERROR_RESPONSE_DTO{fmt.Sprint(http.StatusUnprocessableEntity) + " failed to get user_uuid from messageData"}, []string{uuid})
		return
	}
	_, err := get_user_id_by_uuid(uuid)
	if err != nil {
		log.Println("failed to get ID of the message sender", err.Error())
		wsSend(WS_ERROR_RESPONSE, WS_ERROR_RESPONSE_DTO{fmt.Sprint(http.StatusUnprocessableEntity) + " failed to get ID of the message sender"}, []string{uuid})
		return
	}
	rows, err := statements["getAllGroups"].Query()
	if err != nil {
		log.Println("getAllGroups query failed", err.Error())
		wsSend(WS_ERROR_RESPONSE, WS_ERROR_RESPONSE_DTO{fmt.Sprint(http.StatusInternalServerError) + " getAllGroups query failed"}, []string{uuid})
		return
	}
	defer rows.Close()
	var groups WS_GROUPS_LIST_DTO
	for rows.Next() {
		var group WS_GROUP_RESPONSE_DTO
		err = rows.Scan(&group.Id, &group.Name, &group.Description, &group.Created_at, &group.Email, &group.First_name, &group.Last_name)
		if err != nil {
			log.Println("getAllGroups query failed to scan", err.Error())
			wsSend(WS_ERROR_RESPONSE, WS_ERROR_RESPONSE_DTO{fmt.Sprint(http.StatusInternalServerError) + " getAllGroups query failed to scan row"}, []string{uuid})
			return
		}
		groups = append(groups, group)
	}

	wsSend(WS_GROUPS_LIST, groups, []string{uuid})
}

// wsGroupRequestSubmitHandler add user to groupPendingMembers table
//
// @params: group_id
func wsGroupRequestSubmitHandler(conn *websocket.Conn, messageData map[string]interface{}) {
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

	_group_id, ok := messageData["group_id"].(float64)
	if !ok {
		log.Println("failed to get group_id from messageData")
		wsSend(WS_ERROR_RESPONSE, WS_ERROR_RESPONSE_DTO{fmt.Sprint(http.StatusUnprocessableEntity) + " failed to get group_id from messageData"}, []string{uuid})
		return
	}
	group_id := int64(_group_id)

	// add the member to the groupPendingMembers table
	_, err = statements["addGroupPendingMember"].Exec(group_id, user_id)
	if err != nil {
		log.Println("addGroupPendingMember query failed", err.Error())
		wsSend(WS_ERROR_RESPONSE, WS_ERROR_RESPONSE_DTO{fmt.Sprint(http.StatusInternalServerError) + " addGroupPendingMember query failed"}, []string{uuid})
		return
	}

	//todo: perhaps remove later, at the moment not used
	wsSend(WS_SUCCESS_RESPONSE, WS_SUCCESS_RESPONSE_DTO{fmt.Sprint(http.StatusOK) + " group joining request sent to group creator, waiting for approval"}, []string{uuid})

	// send back the updated visitor status
	wsUserGroupVisitorStatusHandler(conn, messageData)
}

// wsGroupRequestsListHandler gets a list of all the requests to join the group from group_pending_members table
//
// @params: {group_id : int}
func wsGroupRequestsListHandler(conn *websocket.Conn, messageData map[string]interface{}) {
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

	// use getCreatorAllGroupsPendings and user_id, to get all pendings for the groups, where user_id is creator of the group

	rows, err := statements["getCreatorAllGroupsPendings"].Query(user_id)
	if err != nil {
		log.Println("getCreatorAllGroupsPendings query failed", err.Error())
		wsSend(WS_ERROR_RESPONSE, WS_ERROR_RESPONSE_DTO{fmt.Sprint(http.StatusInternalServerError) + " getCreatorAllGroupsPendings query failed"}, []string{uuid})
		return
	}
	defer rows.Close()

	var requests_list WS_GROUP_REQUESTS_LIST_DTO
	for rows.Next() {
		var request WS_GROUP_REQUEST_RESPONSE_DTO
		err = rows.Scan(
			&request.Group_id,
			&request.Member_id,
			&request.Group_name,
			&request.Description,
			&request.Email,
			&request.First_name,
			&request.Last_name,
		)
		if err != nil {
			log.Println("request scan failed", err.Error())
			wsSend(WS_ERROR_RESPONSE, WS_ERROR_RESPONSE_DTO{fmt.Sprint(http.StatusInternalServerError) + " request scan failed"}, []string{uuid})
			return
		}
		requests_list = append(requests_list, request)
	}

	wsSend(WS_GROUP_REQUESTS_LIST, requests_list, []string{uuid})
}

// wsGroupInvitesSubmitHandler send invites to the group for emails(users) space separated
//
// @params: group_id, invited_emails
func wsGroupInvitesSubmitHandler(conn *websocket.Conn, messageData map[string]interface{}) {
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

	_group_id, ok := messageData["group_id"].(float64)
	if !ok {
		log.Println("failed to get group_id from messageData")
		wsSend(WS_ERROR_RESPONSE, WS_ERROR_RESPONSE_DTO{fmt.Sprint(http.StatusUnprocessableEntity) + " failed to get group_id from messageData"}, []string{uuid})
		return
	}
	group_id := int64(_group_id)

	invited_emails, ok := messageData["invited_emails"].(string) // space separated emails
	if !ok {
		log.Println("failed to get invited_emails from messageData")
		wsSend(WS_ERROR_RESPONSE, WS_ERROR_RESPONSE_DTO{fmt.Sprint(http.StatusUnprocessableEntity) + " failed to get invited_emails from messageData"}, []string{uuid})
		return
	}
	//remove all not single spaces and trailing spaces
	invited_emails = strings.Join(strings.Fields(invited_emails), " ")
	//check the string is enough long, at least "a@" or "a b"(should not happens, cause filtered by frontend html tag built-in checker for email)
	if len(invited_emails) < 2 {
		log.Println("no invited_emails in messageData")
		wsSend(WS_ERROR_RESPONSE, WS_ERROR_RESPONSE_DTO{fmt.Sprint(http.StatusUnprocessableEntity) + " no invited_emails in messageData"}, []string{uuid})
		return
	}
	emails := strings.Split(invited_emails, " ")

	err_counter := 0

	var invited_user_ids []int
	for _, email := range emails {
		invited_user_id, err := get_user_id_by_email(email)
		if err != nil {
			log.Printf("failed to get ID of the invited user with email [%s] %s", email, err.Error())
			wsSend(WS_ERROR_RESPONSE, WS_ERROR_RESPONSE_DTO{fmt.Sprint(http.StatusUnprocessableEntity) + " failed to get ID of the invited user"}, []string{uuid})
			err_counter++
			continue
		}
		invited_user_ids = append(invited_user_ids, invited_user_id)
	}

	rows, err := statements["getGroupMembers"].Query(group_id)
	if err != nil {
		log.Println("getGroupMembers query failed", err.Error())
		wsSend(WS_ERROR_RESPONSE, WS_ERROR_RESPONSE_DTO{fmt.Sprint(http.StatusInternalServerError) + " getGroupMembers query failed"}, []string{uuid})
		return
	}
	defer rows.Close()

	group_member_ids := map[int]int{}
	for rows.Next() {
		var group_member_id int
		err = rows.Scan(&group_member_id)
		if err != nil {
			log.Println("group_member_id scan failed", err.Error())
			wsSend(WS_ERROR_RESPONSE, WS_ERROR_RESPONSE_DTO{fmt.Sprint(http.StatusInternalServerError) + " scan failed"}, []string{uuid})
			return
		}
		group_member_ids[group_member_id] = group_member_id
	}

	// check if the user is already a member of the group
	var not_member_user_ids []int
	for _, id := range invited_user_ids {
		_, ok := group_member_ids[id]
		if ok {
			log.Printf("user with ID [%d] is already a member of the group with ID [%d]", id, group_id)
			wsSend(WS_ERROR_RESPONSE, WS_ERROR_RESPONSE_DTO{fmt.Sprint(http.StatusUnprocessableEntity) + " user is already a member of the group"}, []string{uuid})
			err_counter++
			continue
		}
		not_member_user_ids = append(not_member_user_ids, id)
	}

	// check if the user is already invited to the group
	rows, err = statements["getGroupInvitedUsers"].Query(group_id)
	if err != nil {
		log.Println("getGroupInvitedUsers query failed", err.Error())
		wsSend(WS_ERROR_RESPONSE, WS_ERROR_RESPONSE_DTO{fmt.Sprint(http.StatusInternalServerError) + " getGroupInvitedUsers query failed"}, []string{uuid})
		return
	}
	defer rows.Close()

	group_invited_user_ids := map[int]int{}
	for rows.Next() {
		var group_invited_user_id int
		err = rows.Scan(&group_invited_user_id)
		if err != nil {
			log.Println("group_invited_user_id scan failed", err.Error())
			wsSend(WS_ERROR_RESPONSE, WS_ERROR_RESPONSE_DTO{fmt.Sprint(http.StatusInternalServerError) + " scan failed"}, []string{uuid})
			return
		}
		group_invited_user_ids[group_invited_user_id] = group_invited_user_id
	}

	var user_ids []int
	// check if the user is already invited to the group
	for _, id := range not_member_user_ids {
		_, ok := group_invited_user_ids[id]
		if ok {
			log.Printf("user with ID [%d] is already invited to the group with ID [%d]", id, group_id)
			wsSend(WS_ERROR_RESPONSE, WS_ERROR_RESPONSE_DTO{fmt.Sprint(http.StatusUnprocessableEntity) + " user is already invited to the group"}, []string{uuid})
			err_counter++
			continue
		}
		user_ids = append(user_ids, id)
	}

	// todo: CHECK THIS BULK APPROACH. if it fails, remove and use for loop, which is slower
	// Prepare the SQL statement for bulk insertion
	bulkInsertStmt, err := db.Prepare("INSERT INTO group_invited_users (user_id, group_id, inviter_id, created_at) VALUES (?, ?, ?, ?)")
	if err != nil {
		log.Println("Failed to prepare bulk insert statement", err.Error())
		wsSend(WS_ERROR_RESPONSE, WS_ERROR_RESPONSE_DTO{fmt.Sprint(http.StatusInternalServerError) + " failed to prepare bulk insert statement"}, []string{uuid})
		return
	}

	// Create a slice to hold the values for bulk insertion
	var values []interface{}

	// Populate the values slice with the data for bulk insertion
	for _, invited_user_id := range user_ids {
		values = append(values, invited_user_id, group_id, user_id, time.Now().Format("2006-01-02 15:04:05"))
	}

	// Execute the bulk insert statement
	_, err = bulkInsertStmt.Exec(values...)
	if err != nil {
		log.Println("Bulk insert failed", err.Error())
		wsSend(WS_ERROR_RESPONSE, WS_ERROR_RESPONSE_DTO{fmt.Sprint(http.StatusInternalServerError) + " bulk insert failed"}, []string{uuid})
		return
	}

	// for _, invited_user_id := range user_ids {
	// 	// add the invited user to the group_invited_users table
	// 	_, err = statements["addGroupInvitedUser"].Exec(invited_user_id, group_id, user_id, time.Now().Format("2006-01-02 15:04:05"))
	// 	if err != nil {
	// 		log.Println("addGroupInvitedUser query failed", err.Error())
	// 		wsSend(WS_ERROR_RESPONSE,WS_ERROR_RESPONSE_DTO{fmt.Sprint(http.StatusInternalServerError) + " addGroupInvitedUser query failed"},[]string{uuid})
	// 		err_counter++
	// 		continue
	// 	}
	// }

	wsSend(WS_SUCCESS_RESPONSE, WS_SUCCESS_RESPONSE_DTO{fmt.Sprintf("%d Number of errors:%d. User/s invited to the group", http.StatusOK, err_counter)}, []string{uuid})
}

// wsGroupInvitesListHandler returns invites into the groups for the user
func wsGroupInvitesListHandler(conn *websocket.Conn, messageData map[string]interface{}) {
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

	// get the invites into the groups for the user
	var invites_list WS_GROUP_INVITES_LIST_DTO
	rows, err := statements["getUserInvites"].Query(user_id)
	if err != nil {
		log.Println("getUserInvites query failed", err.Error())
		wsSend(WS_ERROR_RESPONSE, WS_ERROR_RESPONSE_DTO{fmt.Sprint(http.StatusInternalServerError) + " getUserInvites query failed"}, []string{uuid})
		return
	}
	defer rows.Close()
	for rows.Next() {
		var invite WS_GROUP_INVITE_RESPONSE_DTO
		err = rows.Scan(
			&invite.Group_id,
			&invite.Group_name,
			&invite.Group_description,
			&invite.Created_at,
			&invite.Inviter_email,
			&invite.Inviter_first_name,
			&invite.Inviter_last_name)
		if err != nil {
			log.Println("invite scan failed", err.Error())
			wsSend(WS_ERROR_RESPONSE, WS_ERROR_RESPONSE_DTO{fmt.Sprint(http.StatusInternalServerError) + " invite scan failed"}, []string{uuid})
			return
		}
		invites_list = append(invites_list, invite)
	}

	wsSend(WS_GROUP_INVITES_LIST, invites_list, []string{uuid})
}
