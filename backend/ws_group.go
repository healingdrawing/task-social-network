package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

type WS_GROUP_SUBMIT_DTO struct {
	User_uuid   string `json:"user_uuid"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Privacy     string `json:"privacy"`
	Invited     string `json:"invited"` // space separated emails
}

type WS_GROUP_RESPONSE_DTO struct {
	Id            int       `json:"id"`
	Name          string    `json:"name"`
	Description   string    `json:"description"`
	Creation_date time.Time `json:"creation_date"`
	Privacy       string    `json:"privacy"`
	Email         string    `json:"email"`
	First_name    string    `json:"first_name"`
	Last_name     string    `json:"last_name"`
}

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
	defer wsRecover()

	var data WS_GROUP_SUBMIT_DTO

	fields := map[string]*string{
		"user_uuid":   &data.User_uuid,
		"name":        &data.Name,
		"description": &data.Description,
		"privacy":     &data.Privacy,
		"invited":     &data.Invited,
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

	user_id, err := getIDbyUUID(data.User_uuid)
	if err != nil {
		log.Println("failed to get ID of the request sender", err.Error())
		wsSendError(WS_ERROR_RESPONSE_DTO{fmt.Sprint(http.StatusUnprocessableEntity) + " failed to get ID of the request sender"})
		return
	}

	// sanitize data
	data.Name = strings.TrimSpace(data.Name)
	data.Description = strings.TrimSpace(data.Description)
	data.Privacy = strings.TrimSpace(data.Privacy)
	if data.Name == "" || data.Description == "" || data.Privacy == "" {
		log.Println("empty fields")
		wsSendError(WS_ERROR_RESPONSE_DTO{fmt.Sprint(http.StatusUnprocessableEntity) + " empty fields"})
		return
	}
	result, err := statements["addGroup"].Exec(data.Name, data.Description, user_id, time.Now().Format("2006-01-02 15:04:05"), data.Privacy)
	if err != nil {
		log.Println("addGroup query failed", err.Error())
		wsSendError(WS_ERROR_RESPONSE_DTO{fmt.Sprint(http.StatusInternalServerError) + " addGroup query failed"})
		return
	}
	group_id, err := result.LastInsertId()
	if err != nil {
		log.Println("failed to get last insert id", err.Error())
		wsSendError(WS_ERROR_RESPONSE_DTO{fmt.Sprint(http.StatusInternalServerError) + " failed to get last insert id"})
		return
	}
	invited_users_emails_string := strings.TrimSpace(data.Invited)
	if invited_users_emails_string != "" {
		invited_users_emails := strings.Split(invited_users_emails_string, " ")
		for _, email := range invited_users_emails {
			// get the user id from the email
			var invited_user_id int
			err = statements["getUserIDByEmail"].QueryRow(email).Scan(&invited_user_id)
			if err != nil {
				log.Println("getUserIDByEmail query failed", err.Error())
				wsSendError(WS_ERROR_RESPONSE_DTO{fmt.Sprint(http.StatusInternalServerError) + " getUserIDByEmail query failed"})
				return
			}
			// add the invited user to the groupPendingMembers table
			_, err = statements["addGroupInvitedUser"].Exec(invited_user_id, group_id, user_id, time.Now().Format("2006-01-02 15:04:05"))
			if err != nil {
				log.Println("addGroupInvitedUser query failed", err.Error())
				wsSendError(WS_ERROR_RESPONSE_DTO{fmt.Sprint(http.StatusInternalServerError) + " addGroupInvitedUser query failed"})
				return
			}
		}
	}
	// add group creator to group members
	_, err = statements["addGroupMember"].Exec(group_id, user_id)
	if err != nil {
		log.Println("addGroupMember query failed", err.Error())
		wsSendError(WS_ERROR_RESPONSE_DTO{fmt.Sprint(http.StatusInternalServerError) + " addGroupMember query failed"})
		return
	}
	// send response
	wsSendSuccess(WS_SUCCESS_RESPONSE_DTO{fmt.Sprint(http.StatusOK) + " group created"})

}

// groupGetHandler makes the user join the group
//
// @params: group_id
func groupJoinHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	defer recovery(w)
	requestorID, err := getRequestSenderID(r)
	if err != nil {
		jsonResponse(w, 401, err.Error())
		return
	}

	var data struct {
		GroupID int `json:"group_id"`
	}

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err = decoder.Decode(&data)
	if err != nil {
		log.Println(err.Error())
		jsonResponse(w, http.StatusUnprocessableEntity, " malformed json, bad request")
		return
	}
	// check if group is private
	var group Group
	err = statements["getGroup"].QueryRow(data.GroupID).Scan(&group.ID, &group.Name, &group.Description, &group.CreatorId, &group.CreationDate, &group.Privacy)
	if err != nil {
		log.Println(err.Error())
		jsonResponse(w, http.StatusInternalServerError, "getGroup query failed, group not found")
		return
	}
	if group.Privacy == "private" {
		// add the member to the groupPendingMembers table
		_, err = statements["addGroupPendingMember"].Exec(data.GroupID, requestorID)
		if err != nil {
			log.Println(err.Error())
			jsonResponse(w, http.StatusInternalServerError, "addGroupPendingMember query failed, failed to add user to group pending members")
			return
		}
		jsonResponse(w, http.StatusOK, "group joining request sent to group creator, waiting for approval")
		return
	}

	_, err = statements["addGroupMember"].Exec(data.GroupID, requestorID)
	if err != nil {
		log.Println(err.Error())
		jsonResponse(w, http.StatusInternalServerError, "addGroupMember query failed to add user to group membership")
		return
	}
	jsonResponse(w, http.StatusOK, "group joined")
}

// groupGetHandler makes the user leave the group
//
// @params: group_id
func groupLeaveHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	defer recovery(w)
	requestSenderID, err := getRequestSenderID(r)
	if err != nil {
		jsonResponse(w, 401, err.Error())
		return
	}

	var data struct {
		GroupID int `json:"group_id"`
	}

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err = decoder.Decode(&data)
	if err != nil {
		log.Println(err.Error())
		jsonResponse(w, http.StatusBadRequest, "")
		return
	}

	// check if the request sender is the group creator
	var group Group
	err = statements["getGroup"].QueryRow(data.GroupID).Scan(&group.ID, &group.Name, &group.Description, &group.CreatorId, &group.CreationDate, &group.Privacy)
	if err != nil {
		log.Println(err.Error())
		jsonResponse(w, http.StatusInternalServerError, "getGroup query failed, group not found")
		return
	}
	if group.CreatorId == requestSenderID {
		// tell them that creator can't leave the group
		jsonResponse(w, http.StatusForbidden, "group creator can't leave the group")
		return
	}
	// remove the member from the groupMembers table
	_, err = statements["removeGroupMember"].Exec(data.GroupID, requestSenderID)
	if err != nil {
		log.Println(err.Error())
		jsonResponse(w, http.StatusInternalServerError, "removeGroupMember query failed to remove user from group membership")
		return
	}
	jsonResponse(w, http.StatusOK, fmt.Sprintf("user with id %d left the group %s", requestSenderID, group.Name))
}

// groupInviteHandler makes the user join the group_invited_users table
//
// @params: group_id, member_email
// only one person is invited at a time via this handler
func groupInviteHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	defer recovery(w)
	requestSenderID, err := getRequestSenderID(r)
	if err != nil {
		jsonResponse(w, 401, err.Error())
		return
	}
	var data struct {
		GroupID      int    `json:"group_id"`
		InvitedEmail string `json:"member_email"`
	}
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err = decoder.Decode(&data)
	if err != nil {
		log.Println(err.Error())
		jsonResponse(w, http.StatusBadRequest, "")
		return
	}

	// get the id of the invited user
	var invitedUserID int
	err = statements["getUserIDByEmail"].QueryRow(data.InvitedEmail).Scan(&invitedUserID)
	if err != nil {
		log.Println(err.Error())
		jsonResponse(w, http.StatusInternalServerError, "getUserIDByEmail query failed. "+fmt.Sprintf("user with email %s not found", data.InvitedEmail))
		return
	}
	// add the invited user to the groupInvitedUsers table
	_, err = statements["addGroupInvitedUser"].Exec(invitedUserID, data.GroupID, requestSenderID, time.Now().Format("2006-01-02 15:04:05"))
	if err != nil {
		log.Println(err.Error())
		jsonResponse(w, http.StatusInternalServerError, "addGroupInvitedUser query failed")
		return
	}
	jsonResponse(w, http.StatusOK, "user invited to the group")
}

// groupInvitedHandler returns the list of users invited to the group
//
// @params: group_id
func groupInvitedHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	defer recovery(w)
	// just to check session
	_, err := getRequestSenderID(r)
	if err != nil {
		jsonResponse(w, 401, err.Error())
		return
	}
	var data struct {
		GroupID int `json:"group_id"`
	}

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err = decoder.Decode(&data)
	if err != nil {
		log.Println(err.Error())
		jsonResponse(w, http.StatusBadRequest, "")
		return
	}
	// get the list of users invited to the group
	type invitedUserData struct {
		InvitedUserID int
		GroupID       int
		InviterID     int
		InvitaionTime time.Time
	}
	var invitedUsers []invitedUserData
	rows, err := statements["getGroupInvitedUsers"].Query(data.GroupID)
	if err != nil {
		log.Println(err.Error())
		jsonResponse(w, http.StatusInternalServerError, "getGroupInvitedUsers query failed")
		return
	}
	for rows.Next() {
		var invitedUser invitedUserData
		err = rows.Scan(&invitedUser.InvitedUserID, &invitedUser.GroupID, &invitedUser.InviterID, &invitedUser.InvitaionTime)
		if err != nil {
			log.Println(err.Error())
			jsonResponse(w, http.StatusInternalServerError, "failed to scan invited user id")
			return
		}
		invitedUsers = append(invitedUsers, invitedUser)
	}

	type outgoingData struct {
		FullName        string    `json:"full_name"`
		Email           string    `json:"email"`
		InviterFullName string    `json:"inviter_full_name"`
		InviterEmail    string    `json:"inviter_email"`
		InvitaionTime   time.Time `json:"invitation_time"`
	}

	invitedUsersInfo := []outgoingData{}

	for _, invitedUserID := range invitedUsers {
		var invitedUserInfo outgoingData
		var invitedFirstName, invitedLastName, invitedNick, inviterFirstName, inviterLastName, invitorNick string
		err = statements["getUserbyID"].QueryRow(invitedUserID.InvitedUserID).Scan(&invitedUserInfo.Email, &invitedFirstName, &invitedLastName, &invitedNick)
		if err != nil {
			log.Println(err.Error())
			jsonResponse(w, http.StatusInternalServerError, "getUserbyID query failed")
			return
		}
		invitedUserInfo.FullName = invitedFirstName + " " + invitedLastName
		err = statements["getUserbyID"].QueryRow(invitedUserID.InviterID).Scan(&invitedUserInfo.InviterEmail, &inviterFirstName, &inviterLastName, &invitorNick)
		if err != nil {
			log.Println(err.Error())
			jsonResponse(w, http.StatusInternalServerError, "getUserbyID query failed")
			return
		}
		invitedUserInfo.InviterFullName = inviterFirstName + " " + inviterLastName
		invitedUserInfo.InvitaionTime = invitedUserID.InvitaionTime
		invitedUsersInfo = append(invitedUsersInfo, invitedUserInfo)
	}

	w.WriteHeader(200)
	jsonResponseObj, _ := json.Marshal(invitedUsersInfo)
	_, err = w.Write(jsonResponseObj)
	if err != nil {
		jsonResponse(w, http.StatusInternalServerError, "w.Write(jsonResponseObj)<-invitedUsersInfo failed")
		return
	}

}

// # groupRequestsHandler gets a list of all the requests to join the group from group_pending_members table
//
// @r.params: {group_id : int}
func groupRequestsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	defer recovery(w)
	requestSenderID, err := getRequestSenderID(r)
	if err != nil {
		jsonResponse(w, 401, err.Error())
		return
	}
	// check if request sender is the group creator
	var data struct {
		GroupID int `json:"group_id"`
	}

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err = decoder.Decode(&data)
	if err != nil {
		log.Println(err.Error())
		jsonResponse(w, http.StatusBadRequest, "")
		return
	}

	var group Group
	err = statements["getGroup"].QueryRow(data.GroupID).Scan(&group.ID, &group.Name, &group.Description, &group.CreatorId, &group.CreationDate, &group.Privacy)
	if err != nil {
		log.Println(err.Error())
		jsonResponse(w, http.StatusNotFound, "getGroup query failed, group not found")
		return
	}

	if group.CreatorId != requestSenderID {
		jsonResponse(w, http.StatusForbidden, "only group creator can view requests")
		return
	}

	// get the list of users invited to the group
	type pendingUserData struct {
		PendingUserID int
	}
	var pendingUsers []pendingUserData
	rows, err := statements["getGroupPendingMembers"].Query(data.GroupID)
	if err != nil {
		log.Println(err.Error())
		jsonResponse(w, http.StatusInternalServerError, "getGroupPendingMembers query failed")
		return
	}
	for rows.Next() {
		var pendingUser pendingUserData
		err = rows.Scan(&pendingUser.PendingUserID)
		if err != nil {
			log.Println(err.Error())
			jsonResponse(w, http.StatusInternalServerError, "failed to scan pending user id")
			return
		}
		pendingUsers = append(pendingUsers, pendingUser)
	}
	// send the array of name and email of the pending users as response
	type outgoingData struct {
		FullName string `json:"full_name"`
		Email    string `json:"email"`
	}
	// get the full name and email of the pending users
	pendingUsersInfo := []outgoingData{}
	for _, pendingUserID := range pendingUsers {
		var pendingUserInfo outgoingData
		var pendingFirstName, pendingLastName, pendingNick string
		err = statements["getUserbyID"].QueryRow(pendingUserID.PendingUserID).Scan(&pendingUserInfo.Email, &pendingFirstName, &pendingLastName, &pendingNick)
		if err != nil {
			log.Println(err.Error())
			jsonResponse(w, http.StatusInternalServerError, "getUserbyID query failed")
			return
		}
		pendingUserInfo.FullName = pendingFirstName + " " + pendingLastName
		pendingUsersInfo = append(pendingUsersInfo, pendingUserInfo)
	}

	w.WriteHeader(200)
	jsonResponseObj, _ := json.Marshal(pendingUsersInfo)
	_, err = w.Write(jsonResponseObj)
	if err != nil {
		jsonResponse(w, http.StatusInternalServerError, "w.Write(jsonResponseObj)<-pendingUsersInfo failed")
		return
	}

}
