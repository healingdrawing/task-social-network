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
	defer func() {
		if err := recover(); err != nil {
			log.Println(err)
			w.WriteHeader(500)
			jsonResponse, _ := json.Marshal(map[string]string{
				"message": "internal server error",
			})
			w.Write(jsonResponse)
		}
	}()
	var data Group
	var incomingData GroupCreationDTO
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&incomingData)
	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(400)
		jsonResponse, _ := json.Marshal(map[string]string{
			"message": "bad request, malformed json",
		})
		w.Write(jsonResponse)
		return
	}

	// get user id from the cookie
	cookie, err := r.Cookie("user_uuid")
	if err != nil {
		w.WriteHeader(401)
		jsonResponse, _ := json.Marshal(map[string]string{
			"message": "unauthorized",
		})
		w.Write(jsonResponse)
		return
	}
	incomingData.UUID = cookie.Value

	data.Name = strings.TrimSpace(incomingData.Name)
	data.Description = strings.TrimSpace(incomingData.Description)
	data.Privacy = strings.TrimSpace(incomingData.Privacy)
	data.CreatorId, err = getIDbyUUID(incomingData.UUID)
	if data.Name == "" || data.Description == "" || data.Privacy == "" {
		w.WriteHeader(400)
		jsonResponse, _ := json.Marshal(map[string]string{
			"message": "bad request, missing fields",
		})
		w.Write(jsonResponse)
		return
	}
	result, err := statements["addGroup"].Exec(data.Name, data.Description, data.CreatorId, time.Now().Format("2006-01-02 15:04:05"), data.Privacy)
	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(500)
		jsonResponse, _ := json.Marshal(map[string]string{
			"message": "internal server error, addGroup query failed",
		})
		w.Write(jsonResponse)
		return
	}
	groupID, err := result.LastInsertId()
	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(500)
		jsonResponse, _ := json.Marshal(map[string]string{
			"message": "internal server error, failed to get last insert id of group creation",
		})
		w.Write(jsonResponse)
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
				w.WriteHeader(404)
				errorMsg := fmt.Sprintf("user with email %s not found", email)
				jsonResponse, _ := json.Marshal(map[string]string{
					"message": errorMsg,
				})
				w.Write(jsonResponse)
				return
			}
			// add the invited user to the groupPendingMembers table
			_, err = statements["addGroupInvitedUser"].Exec(invitedUserID, groupID, data.CreatorId, time.Now().Format("2006-01-02 15:04:05"))
			if err != nil {
				log.Println(err.Error())
				w.WriteHeader(500)
				jsonResponse, _ := json.Marshal(map[string]string{
					"message": "internal server error, failed to add user to group pending members",
				})
				w.Write(jsonResponse)
				return
			}
		}
	}
	// add group creator to group members
	_, err = statements["addGroupMember"].Exec(groupID, data.CreatorId)
	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(500)
		jsonResponse, _ := json.Marshal(map[string]string{
			"message": "addGroupMember query failed to add creator to group membership",
		})
		w.Write(jsonResponse)
		return
	}
	w.WriteHeader(200)
	jsonResponse, _ := json.Marshal(map[string]string{
		"message": "group created",
	})
	w.Write(jsonResponse)
}

// groupGetHandler makes the user join the group
//
// @params: group_id
func groupJoinHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	defer func() {
		if err := recover(); err != nil {
			log.Println(err)
			w.WriteHeader(500)
			jsonResponse, _ := json.Marshal(map[string]string{
				"message": "internal server error",
			})
			w.Write(jsonResponse)
		}
	}()

	// extract requestor id from cookie
	cookie, err := r.Cookie("user_uuid")
	if err != nil {
		w.WriteHeader(401)
		jsonResponse, _ := json.Marshal(map[string]string{
			"message": "unauthorized",
		})
		w.Write(jsonResponse)
		return
	}
	requestorID, err := getIDbyUUID(cookie.Value)
	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(500)
		jsonResponse, _ := json.Marshal(map[string]string{
			"message": "internal server error, failed to get id of request sender",
		})
		w.Write(jsonResponse)
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
		w.WriteHeader(400)
		jsonResponse, _ := json.Marshal(map[string]string{
			"message": "bad request",
		})
		w.Write(jsonResponse)
		return
	}
	// check if group is private
	var group Group
	err = statements["getGroup"].QueryRow(data.GroupID).Scan(&group.ID, &group.Name, &group.Description, &group.CreatorId, &group.CreationDate, &group.Privacy)
	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(404)
		jsonResponse, _ := json.Marshal(map[string]string{
			"message": "group not found",
		})
		w.Write(jsonResponse)
		return
	}
	if group.Privacy == "private" {
		// add the member to the groupPendingMembers table
		_, err = statements["addGroupPendingMember"].Exec(data.GroupID, requestorID)
		if err != nil {
			log.Println(err.Error())
			w.WriteHeader(500)
			jsonResponse, _ := json.Marshal(map[string]string{
				"message": "internal server error",
			})
			w.Write(jsonResponse)
			return
		}
		w.WriteHeader(200)
		jsonResponse, _ := json.Marshal(map[string]string{
			"message": "group joining request sent to group creator, waiting for approval",
		})
		w.Write(jsonResponse)
		return
	}

	_, err = statements["addGroupMember"].Exec(data.GroupID, requestorID)
	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(500)
		jsonResponse, _ := json.Marshal(map[string]string{
			"message": "internal server error",
		})
		w.Write(jsonResponse)
		return
	}
	w.WriteHeader(200)
	jsonResponse, _ := json.Marshal(map[string]string{
		"message": "group joined",
	})
	w.Write(jsonResponse)
}

// groupGetHandler makes the user leave the group
//
// @params: group_id
func groupLeaveHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	defer func() {
		if err := recover(); err != nil {
			log.Println(err)
			w.WriteHeader(500)
			jsonResponse, _ := json.Marshal(map[string]string{
				"message": "internal server error",
			})
			w.Write(jsonResponse)
		}
	}()

	// get the id of the request sender
	cookie, err := r.Cookie("user_uuid")
	if err != nil {
		w.WriteHeader(401)
		jsonResponse, _ := json.Marshal(map[string]string{
			"message": "unauthorized, malformed cookie",
		})
		w.Write(jsonResponse)
		return
	}
	requestSenderID, err := getIDbyUUID(cookie.Value)
	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(500)
		jsonResponse, _ := json.Marshal(map[string]string{
			"message": "internal server error, failed to get id of request sender",
		})
		w.Write(jsonResponse)
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
		w.WriteHeader(400)
		jsonResponse, _ := json.Marshal(map[string]string{
			"message": "bad request",
		})
		w.Write(jsonResponse)
		return
	}

	// check if the request sender is the group creator
	var group Group
	err = statements["getGroup"].QueryRow(data.GroupID).Scan(&group.ID, &group.Name, &group.Description, &group.CreatorId, &group.CreationDate, &group.Privacy)
	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(404)
		jsonResponse, _ := json.Marshal(map[string]string{
			"message": "group not found",
		})
		w.Write(jsonResponse)
		return
	}
	if group.CreatorId == requestSenderID {
		// tell them that creator can't leave the group
		w.WriteHeader(403)
		jsonResponse, _ := json.Marshal(map[string]string{
			"message": "group creator can't leave the group",
		})
		w.Write(jsonResponse)
		return
	}
	// remove the member from the groupMembers table
	_, err = statements["removeGroupMember"].Exec(data.GroupID, requestSenderID)
	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(500)
		jsonResponse, _ := json.Marshal(map[string]string{
			"message": "internal server error",
		})
		w.Write(jsonResponse)
		return
	}
	w.WriteHeader(200)
	returnMessage := fmt.Sprintf("user with id %d left the group %s", requestSenderID, group.Name)
	jsonResponse, _ := json.Marshal(map[string]string{
		"message": returnMessage,
	})
	w.Write(jsonResponse)
}

// groupInviteHandler makes the user join the group_invited_users table
//
// @params: group_id, member_email
// only one person is invited at a time via this handler
func groupInviteHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	defer func() {
		if err := recover(); err != nil {
			log.Println(err)
			w.WriteHeader(500)
			jsonResponse, _ := json.Marshal(map[string]string{
				"message": "internal server error",
			})
			w.Write(jsonResponse)
		}
	}()
	var data struct {
		GroupID      int    `json:"group_id"`
		InvitedEmail string `json:"member_email"`
	}
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&data)
	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(400)
		jsonResponse, _ := json.Marshal(map[string]string{
			"message": "bad request",
		})
		w.Write(jsonResponse)
		return
	}
	// get the id of the request sender
	cookie, err := r.Cookie("user_uuid")
	if err != nil {
		w.WriteHeader(401)
		jsonResponse, _ := json.Marshal(map[string]string{
			"message": "unauthorized",
		})
		w.Write(jsonResponse)
		return
	}
	requestSenderID, err := getIDbyUUID(cookie.Value)
	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(500)
		jsonResponse, _ := json.Marshal(map[string]string{
			"message": "internal server error, failed to get id of request sender",
		})
		w.Write(jsonResponse)
		return
	}
	// get the id of the invited user
	var invitedUserID int
	err = statements["getUserIDByEmail"].QueryRow(data.InvitedEmail).Scan(&invitedUserID)
	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(404)
		errorMsg := fmt.Sprintf("user with email %s not found", data.InvitedEmail)
		jsonResponse, _ := json.Marshal(map[string]string{
			"message": errorMsg,
		})
		w.Write(jsonResponse)
		return
	}
	// add the invited user to the groupInvitedUsers table
	_, err = statements["addGroupInvitedUser"].Exec(invitedUserID, data.GroupID, requestSenderID, time.Now().Format("2006-01-02 15:04:05"))
	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(500)
		jsonResponse, _ := json.Marshal(map[string]string{
			"message": "addGroupInvitedUser query failed",
		})
		w.Write(jsonResponse)
		return
	}
	w.WriteHeader(200)
	jsonResponse, _ := json.Marshal(map[string]string{
		"message": "user invited to the group",
	})
	w.Write(jsonResponse)
}

// groupInvitedHandler returns the list of users invited to the group
//
// @params: group_id
func groupInvitedHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	defer func() {
		if err := recover(); err != nil {
			log.Println(err)
			w.WriteHeader(500)
			jsonResponse, _ := json.Marshal(map[string]string{
				"message": "internal server error",
			})
			w.Write(jsonResponse)
		}
	}()
	// get the id of the request sender
	cookie, err := r.Cookie("user_uuid")
	if err != nil {
		w.WriteHeader(401)
		jsonResponse, _ := json.Marshal(map[string]string{
			"message": "unauthorized",
		})
		w.Write(jsonResponse)
		return
	}
	_, err = getIDbyUUID(cookie.Value)
	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(500)
		jsonResponse, _ := json.Marshal(map[string]string{
			"message": "internal server error, failed to get id of request sender",
		})
		w.Write(jsonResponse)
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
		w.WriteHeader(400)
		jsonResponse, _ := json.Marshal(map[string]string{
			"message": "bad request",
		})
		w.Write(jsonResponse)
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
		w.WriteHeader(500)
		jsonResponse, _ := json.Marshal(map[string]string{
			"message": "getGroupInvitedUsers query failed",
		})
		w.Write(jsonResponse)
		return
	}
	for rows.Next() {
		var invitedUser invitedUserData
		err = rows.Scan(&invitedUser.InvitedUserID, &invitedUser.GroupID, &invitedUser.InviterID, &invitedUser.InvitaionTime)
		if err != nil {
			log.Println(err.Error())
			w.WriteHeader(500)
			jsonResponse, _ := json.Marshal(map[string]string{
				"message": "failed to scan invited user id",
			})
			w.Write(jsonResponse)
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
			w.WriteHeader(500)
			jsonResponse, _ := json.Marshal(map[string]string{
				"message": "getUserbyID query failed",
			})
			w.Write(jsonResponse)
			return
		}
		invitedUserInfo.FullName = invitedFirstName + " " + invitedLastName
		err = statements["getUserbyID"].QueryRow(invitedUserID.InviterID).Scan(&invitedUserInfo.InviterEmail, &inviterFirstName, &inviterLastName, &invitorNick)
		if err != nil {
			log.Println(err.Error())
			w.WriteHeader(500)
			jsonResponse, _ := json.Marshal(map[string]string{
				"message": "getUserbyID query failed",
			})
			w.Write(jsonResponse)
			return
		}
		invitedUserInfo.InviterFullName = inviterFirstName + " " + inviterLastName
		invitedUserInfo.InvitaionTime = invitedUserID.InvitaionTime
		invitedUsersInfo = append(invitedUsersInfo, invitedUserInfo)
	}

	w.WriteHeader(200)
	jsonResponse, _ := json.Marshal(invitedUsersInfo)
	w.Write(jsonResponse)
}
