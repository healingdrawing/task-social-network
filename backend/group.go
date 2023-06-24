package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
)

type Group struct {
	ID           int       `json:"ID"`
	Name         string    `json:"name"`
	Description  string    `json:"description"`
	CreatorId    int       `json:"creator"`
	CreationDate time.Time `json:"creation_date"`
	Privacy      string    `json:"privacy"`
}

type GroupMember struct {
	GroupID     int `json:"group_id"`
	MemberEmail int `json:"member_email"`
}

type GroupInvitedMember struct {
	GroupID      int `json:"group_id"`
	MemberEmails int `json:"member_emails"`
}

type GroupCreationDTO struct {
	UUID        string `json:"UUID"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Privacy     string `json:"privacy"`
	Invited     string `json:"invited"`
}

type InvitedUserInfoDTO struct {
	FullName string `json:"full_name"`
	Email    string `json:"email"`
}

func groupNewHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	defer recovery(w)
	var data Group
	err := error(nil)
	data.CreatorId, err = getRequestSenderID(r)
	if err != nil {
		jsonResponse(w, 401, err.Error())
		return
	}
	var incomingData GroupCreationDTO
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err = decoder.Decode(&incomingData)
	if err != nil {
		log.Println(err.Error())
		jsonResponse(w, http.StatusUnprocessableEntity, " malformed json")
		return
	}

	data.Name = strings.TrimSpace(incomingData.Name)
	data.Description = strings.TrimSpace(incomingData.Description)
	data.Privacy = strings.TrimSpace(incomingData.Privacy)
	if data.Name == "" || data.Description == "" || data.Privacy == "" {
		jsonResponse(w, http.StatusBadRequest, "missing fields")
		return
	}
	result, err := statements["addGroup"].Exec(data.Name, data.Description, data.CreatorId, time.Now().Format("2006-01-02 15:04:05"), data.Privacy)
	if err != nil {
		log.Println(err.Error())
		jsonResponse(w, http.StatusInternalServerError, " addGroup query failed")
		return
	}
	groupID, err := result.LastInsertId()
	if err != nil {
		log.Println(err.Error())
		jsonResponse(w, http.StatusInternalServerError, " failed to get last insert id of group creation")
		return
	}
	invitedUsersEmailsString := strings.TrimSpace(incomingData.Invited)
	if invitedUsersEmailsString != "" {
		invitedUsersEmails := strings.Split(invitedUsersEmailsString, " ")
		for _, email := range invitedUsersEmails {
			// get the user id from the email
			var invitedUserID int
			err = statements["getUserIDByEmail"].QueryRow(email).Scan(&invitedUserID)
			if err != nil {
				log.Println(err.Error())
				jsonResponse(w, http.StatusInternalServerError, fmt.Sprintf("getUserIDByEmail query failed, user with email %s not found", email))
				return
			}
			// add the invited user to the groupPendingMembers table
			_, err = statements["addGroupInvitedUser"].Exec(invitedUserID, groupID, data.CreatorId, time.Now().Format("2006-01-02 15:04:05"))
			if err != nil {
				log.Println(err.Error())
				jsonResponse(w, http.StatusInternalServerError, "addGroupInvitedUser query failed, failed to add user to group pending members")
				return
			}
		}
	}
	// add group creator to group members
	_, err = statements["addGroupMember"].Exec(groupID, data.CreatorId)
	if err != nil {
		log.Println(err.Error())
		jsonResponse(w, http.StatusInternalServerError, "addGroupMember query failed to add creator to group membership")
		return
	}
	jsonResponse(w, http.StatusOK, "group created")
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
