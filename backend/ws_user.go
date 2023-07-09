package main

import (
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/websocket"
	_ "github.com/mattn/go-sqlite3"
)

// used in two cases: wsUserVisitorStatusHandler and wsUserGroupVisitorStatusHandler
type WS_USER_VISITOR_STATUS_DTO struct {
	Status string `json:"status"`
}

func wsUserVisitorStatusHandler(conn *websocket.Conn, messageData map[string]interface{}) {
	defer wsRecover()

	uuid, ok := messageData["user_uuid"].(string)
	if !ok {
		log.Println("failed to get user_uuid from message data")
		wsSendError(WS_ERROR_RESPONSE_DTO{fmt.Sprint(http.StatusUnprocessableEntity) + " failed to get user_uuid from message data"})
		return
	}
	user_id, err := get_user_id_by_uuid(uuid)
	if err != nil {
		log.Println("failed to get ID of the message sender", err.Error())
		wsSendError(WS_ERROR_RESPONSE_DTO{fmt.Sprint(http.StatusUnprocessableEntity) + " failed to get ID of the message sender"})
		return
	}

	target_email, ok := messageData["target_email"].(string)
	if !ok {
		log.Println("failed to get target email", err.Error())
		wsSendError(WS_ERROR_RESPONSE_DTO{fmt.Sprint(http.StatusUnprocessableEntity) + " failed to get target email"})
		return
	}

	target_id, err := get_user_id_by_email(target_email)
	if err != nil {
		log.Println("failed to get ID of the target user", err.Error())
		wsSendError(WS_ERROR_RESPONSE_DTO{fmt.Sprint(http.StatusUnprocessableEntity) + " failed to get ID of the target user"})
		return
	}

	isFollower, err := isFollowing(user_id, target_id)
	if err != nil {
		log.Println("failed to check if user is a follower of the target user", err.Error())
		wsSendError(WS_ERROR_RESPONSE_DTO{fmt.Sprint(http.StatusUnprocessableEntity) + " failed to check if user is a follower of the target user"})
		return
	}

	log.Print("isFollower: ", isFollower, " target_id: ", target_id, " user_id: ", user_id)

	var profile WS_USER_PROFILE_DTO
	rows, err := statements["getUserProfile"].Query(target_id)
	if err != nil {
		log.Println("failed to get user profile", err.Error())
		wsSendError(WS_ERROR_RESPONSE_DTO{fmt.Sprint(http.StatusUnprocessableEntity) + " failed to get user profile"})
		return
	}
	defer rows.Close()
	rows.Next()
	err = rows.Scan(&profile.Email, &profile.First_name, &profile.Last_name, &profile.Dob,
		&profile.avatar_bytes, &profile.Nickname, &profile.About_me, &profile.Privacy)
	if err != nil {
		log.Println("failed to scan user profile", err.Error())
		wsSendError(WS_ERROR_RESPONSE_DTO{fmt.Sprint(http.StatusUnprocessableEntity) + " failed to scan user profile"})
		return
	}
	rows.Close()

	visitor := WS_USER_VISITOR_STATUS_DTO{}
	if target_id == user_id {
		visitor.Status = "owner"
	} else if isFollower {
		visitor.Status = "follower"
	} else {
		is_requester, err := isRequester(user_id, target_id)
		if err != nil {
			log.Println("failed to check if user is a requester of the target user", err.Error())
			wsSendError(WS_ERROR_RESPONSE_DTO{fmt.Sprint(http.StatusUnprocessableEntity) + " failed to check if user is a requester of the target user"})
			return
		}
		if is_requester {
			visitor.Status = "requester"
		} else {
			visitor.Status = "visitor"
		}
	}

	log.Println("sending user visitor status", visitor)
	wsSendUserVisitorStatus(visitor)
}

// wsUserGroupVisitorStatusHandler used to manage button on top of group page of frontend. Button will be: request to join group/waiting for decision/member of group.
// - group_id is the group that the user is visiting
//
// status will be one of the following:
//
// "requester", "member", "visitor"
func wsUserGroupVisitorStatusHandler(conn *websocket.Conn, messageData map[string]interface{}) {
	defer wsRecover()

	uuid, ok := messageData["user_uuid"].(string)
	if !ok {
		log.Println("failed to get user_uuid from message data")
		wsSendError(WS_ERROR_RESPONSE_DTO{fmt.Sprint(http.StatusUnprocessableEntity) + " failed to get user_uuid from message data"})
		return
	}
	user_id, err := get_user_id_by_uuid(uuid)
	if err != nil {
		log.Println("failed to get ID of the message sender", err.Error())
		wsSendError(WS_ERROR_RESPONSE_DTO{fmt.Sprint(http.StatusUnprocessableEntity) + " failed to get ID of the message sender"})
		return
	}

	_group_id, ok := messageData["group_id"].(float64)
	if !ok {
		log.Println("failed to get group_id from message data")
		wsSendError(WS_ERROR_RESPONSE_DTO{fmt.Sprint(http.StatusUnprocessableEntity) + " failed to get group_id from message data"})
		return
	}
	group_id := int(_group_id)

	// check if user already made a request to join the group
	is_requester, err := isGroupJoinRequester(user_id, group_id)
	if err != nil {
		log.Println("failed to check if user is a requester of the target group", err.Error())
		wsSendError(WS_ERROR_RESPONSE_DTO{fmt.Sprint(http.StatusUnprocessableEntity) + " failed to check if user is a requester of the target group"})
		return
	}

	if is_requester {
		wsSendUserGroupVisitorStatus(WS_USER_VISITOR_STATUS_DTO{Status: "requester"})
		return
	}

	// check if user is already a member of the group
	is_member, err := isGroupMember(user_id, group_id)
	if err != nil {
		log.Println("failed to check if user is a member of the target group", err.Error())
		wsSendError(WS_ERROR_RESPONSE_DTO{fmt.Sprint(http.StatusUnprocessableEntity) + " failed to check if user is a member of the target group"})
		return
	}

	if is_member {
		wsSendUserGroupVisitorStatus(WS_USER_VISITOR_STATUS_DTO{Status: "member"})
		return
	}
	wsSendUserGroupVisitorStatus(WS_USER_VISITOR_STATUS_DTO{Status: "visitor"})
}

func isGroupJoinRequester(user_id int, group_id int) (bool, error) {
	rows, err := statements["getGroupPendingMember"].Query(group_id, user_id)
	if err != nil {
		return false, err
	}
	defer rows.Close()
	// check length of rows
	// if length is 0, return false
	// else return true
	if rows.Next() {
		return true, nil
	}
	return false, nil
}

func isGroupMember(user_id int, group_id int) (bool, error) {
	rows, err := statements["getGroupMember"].Query(group_id, user_id)
	if err != nil {
		return false, err
	}
	defer rows.Close()
	// check length of rows
	// if length is 0, return false
	// else return true
	if rows.Next() {
		return true, nil
	}
	return false, nil
}

type WS_USER_PROFILE_DTO struct {
	Email        string `json:"email"`
	First_name   string `json:"first_name"`
	Last_name    string `json:"last_name"`
	Dob          string `json:"dob"`
	Avatar       string `json:"avatar"`
	avatar_bytes []byte `sqlite3:"avatar"`
	Nickname     string `json:"nickname"`
	About_me     string `json:"about_me"`
	Public       bool   `json:"public"`
	Privacy      string `sqlite3:"privacy"`
}

type WS_USER_PROFILE_RESPONSE_DTO struct {
	Email      string `json:"email"`
	First_name string `json:"first_name"`
	Last_name  string `json:"last_name"`
	Dob        string `json:"dob"`
	Avatar     string `json:"avatar"`
	Nickname   string `json:"nickname"`
	About_me   string `json:"about_me"`
	Public     bool   `json:"public"`
}

/*
wsUserProfileHandler returns the profile of the target user,

if user has permission to view the target user's profile.

Otherwise, it returns an error:

"403 user does not have permissions to see the target user profile"
*/
func wsUserProfileHandler(conn *websocket.Conn, messageData map[string]interface{}) {
	defer wsRecover()

	uuid, ok := messageData["user_uuid"].(string)
	if !ok {
		log.Println("failed to get user_uuid from message data")
		wsSendError(WS_ERROR_RESPONSE_DTO{fmt.Sprint(http.StatusUnprocessableEntity) + " failed to get user_uuid from message data"})
		return
	}
	user_id, err := get_user_id_by_uuid(uuid)
	if err != nil {
		log.Println("failed to get ID of the message sender", err.Error())
		wsSendError(WS_ERROR_RESPONSE_DTO{fmt.Sprint(http.StatusUnprocessableEntity) + " failed to get ID of the message sender"})
		return
	}

	target_email, ok := messageData["target_email"].(string)
	if !ok {
		log.Println("failed to get target email", err.Error())
		wsSendError(WS_ERROR_RESPONSE_DTO{fmt.Sprint(http.StatusUnprocessableEntity) + " failed to get target email"})
		return
	}

	target_id, err := get_user_id_by_email(target_email)
	if err != nil {
		log.Println("failed to get ID of the target user", err.Error())
		wsSendError(WS_ERROR_RESPONSE_DTO{fmt.Sprint(http.StatusUnprocessableEntity) + " failed to get ID of the target user"})
		return
	}

	isFollower, err := isFollowing(user_id, target_id)
	if err != nil {
		log.Println("failed to check if user is a follower of the target user", err.Error())
		wsSendError(WS_ERROR_RESPONSE_DTO{fmt.Sprint(http.StatusUnprocessableEntity) + " failed to check if user is a follower of the target user"})
		return
	}

	var profile WS_USER_PROFILE_DTO
	rows, err := statements["getUserProfile"].Query(target_id)
	if err != nil {
		log.Println("failed to get user profile", err.Error())
		wsSendError(WS_ERROR_RESPONSE_DTO{fmt.Sprint(http.StatusUnprocessableEntity) + " failed to get user profile"})
		return
	}
	defer rows.Close()
	rows.Next()
	err = rows.Scan(&profile.Email, &profile.First_name, &profile.Last_name, &profile.Dob,
		&profile.avatar_bytes, &profile.Nickname, &profile.About_me, &profile.Privacy)
	if err != nil {
		log.Println("failed to scan user profile", err.Error())
		wsSendError(WS_ERROR_RESPONSE_DTO{fmt.Sprint(http.StatusUnprocessableEntity) + " failed to scan user profile"})
		return
	}
	rows.Close()

	full_profile := true
	if target_id != user_id && !isFollower && profile.Privacy == "private" {
		log.Println("user does not have permissions to see the target user profile")
		wsSendError(WS_ERROR_RESPONSE_DTO{fmt.Sprint(http.StatusForbidden) + " user does not have permissions to see the target user profile"})
		full_profile = false
	}

	profileDTO := WS_USER_PROFILE_RESPONSE_DTO{}

	profileDTO.Email = profile.Email
	profileDTO.First_name = profile.First_name
	profileDTO.Last_name = profile.Last_name
	if full_profile {
		profileDTO.Dob = profile.Dob
		profileDTO.Avatar = base64.StdEncoding.EncodeToString(profile.avatar_bytes)
		profileDTO.Nickname = profile.Nickname
		profileDTO.About_me = profile.About_me
		if profile.Privacy == "public" {
			profileDTO.Public = true
		} else {
			profileDTO.Public = false
		}
	}
	log.Println("sending user profile", profileDTO)
	wsSendUserProfile(profileDTO)

}

func wsChangePrivacyHandler(conn *websocket.Conn, messageData map[string]interface{}) {
	defer wsRecover()

	uuid, ok := messageData["user_uuid"].(string)
	if !ok {
		log.Println("failed to get user_uuid from message data")
		wsSendError(WS_ERROR_RESPONSE_DTO{fmt.Sprint(http.StatusUnprocessableEntity) + " failed to get user_uuid from message data"})
		return
	}
	user_id, err := get_user_id_by_uuid(uuid)
	if err != nil {
		log.Println("failed to get ID of the message sender", err.Error())
		wsSendError(WS_ERROR_RESPONSE_DTO{fmt.Sprint(http.StatusUnprocessableEntity) + " failed to get ID of the message sender"})
		return
	}

	_make_public, ok := messageData["make_public"].(string)
	if !ok {
		log.Println("failed to get make_public", err.Error())
		wsSendError(WS_ERROR_RESPONSE_DTO{fmt.Sprint(http.StatusUnprocessableEntity) + " failed to get make_public"})
		return
	}
	make_public, err := strconv.ParseBool(_make_public)
	if err != nil {
		log.Println("failed to parse make_public", err.Error())
		wsSendError(WS_ERROR_RESPONSE_DTO{fmt.Sprint(http.StatusUnprocessableEntity) + " failed to parse make_public"})
		return
	}

	privacy_value := map[bool]string{true: "public", false: "private"}[make_public]

	_, err = statements["updateUserPrivacy"].Exec(privacy_value, user_id)
	if err != nil {
		log.Println("failed to update user privacy", err.Error())
		wsSendError(WS_ERROR_RESPONSE_DTO{fmt.Sprint(http.StatusUnprocessableEntity) + " failed to update user privacy"})
		return
	}

	wsSendSuccess(WS_SUCCESS_RESPONSE_DTO{fmt.Sprint(http.StatusOK) + " privacy updated"})
}

func isRequester(user_id int, target_id int) (bool, error) {
	rows, err := statements["doesSecondRequesterFollowFirst"].Query(target_id, user_id)
	if err != nil {
		return false, err
	}
	defer rows.Close()
	// check length of rows
	// if length is 0, return false
	// else return true
	if rows.Next() {
		return true, nil
	}
	return false, nil
}
