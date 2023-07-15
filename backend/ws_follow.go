package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type UserForList struct {
	Email      string `json:"email"`
	First_name string `json:"first_name"`
	Last_name  string `json:"last_name"`
	Nickname   string `json:"nickname"`
}

type WS_FOLLOWING_LIST_DTO []UserForList

type WS_FOLLOWERS_LIST_DTO []UserForList

type WS_FOLLOW_REQUESTS_LIST_DTO []UserForList

/*
wsFollowingHandler returns list of users that the target user is following,

if user has permission to view the target user's following list.

Otherwise, it returns an error:

"403 user does not have permissions to see the target user profile"
*/
func wsFollowingListHandler(conn *websocket.Conn, messageData map[string]interface{}) {
	defer wsRecover(messageData)

	uuid, ok := messageData["user_uuid"].(string)
	if !ok {
		log.Println("failed to get user_uuid from message data")
		return
	}

	user_id, err := get_user_id_by_uuid(uuid)
	if err != nil {
		log.Println("failed to get ID of the message sender", err.Error())
		wsSend(WS_ERROR_RESPONSE, WS_ERROR_RESPONSE_DTO{fmt.Sprint(http.StatusUnprocessableEntity) + " failed to get ID of the message sender"}, []string{uuid})
		return
	}

	target_email, ok := messageData["target_email"].(string)
	if !ok {
		log.Println("failed to get target email", err.Error())
		wsSend(WS_ERROR_RESPONSE, WS_ERROR_RESPONSE_DTO{fmt.Sprint(http.StatusUnprocessableEntity) + " failed to get target email"}, []string{uuid})
		return
	}

	target_id, err := get_user_id_by_email(target_email)
	if err != nil {
		log.Println("failed to get ID of the target user", err.Error())
		wsSend(WS_ERROR_RESPONSE, WS_ERROR_RESPONSE_DTO{fmt.Sprint(http.StatusUnprocessableEntity) + " failed to get ID of the target user"}, []string{uuid})
		return
	}

	isFollower, err := isFollowing(user_id, target_id)
	if err != nil {
		log.Println("failed to check if user is a follower of the target user", err.Error())
		wsSend(WS_ERROR_RESPONSE, WS_ERROR_RESPONSE_DTO{fmt.Sprint(http.StatusUnprocessableEntity) + " failed to check if user is a follower of the target user"}, []string{uuid})
		return
	}

	var profile WS_USER_PROFILE_DTO
	rows, err := statements["getUserProfile"].Query(target_id)
	if err != nil {
		log.Println("failed to get user profile", err.Error())
		wsSend(WS_ERROR_RESPONSE, WS_ERROR_RESPONSE_DTO{fmt.Sprint(http.StatusUnprocessableEntity) + " failed to get user profile"}, []string{uuid})
		return
	}
	defer rows.Close()
	rows.Next()
	err = rows.Scan(&profile.Email, &profile.First_name, &profile.Last_name, &profile.Dob,
		&profile.avatar_bytes, &profile.Nickname, &profile.About_me, &profile.Privacy)
	if err != nil {
		log.Println("failed to scan user profile", err.Error())
		wsSend(WS_ERROR_RESPONSE, WS_ERROR_RESPONSE_DTO{fmt.Sprint(http.StatusUnprocessableEntity) + " failed to scan user profile"}, []string{uuid})
		return
	}
	rows.Close()

	if target_id != user_id && !isFollower && profile.Privacy == "private" {
		log.Println("user does not have permissions to see the target user profile")
		wsSend(WS_ERROR_RESPONSE, WS_ERROR_RESPONSE_DTO{fmt.Sprint(http.StatusForbidden) + " user does not have permissions to see the target user profile"}, []string{uuid})
		return
	}

	rows, err = statements["getFollowing"].Query(target_id)
	if err != nil {
		log.Println("failed to get following list", err.Error())
		wsSend(WS_ERROR_RESPONSE, WS_ERROR_RESPONSE_DTO{fmt.Sprint(http.StatusUnprocessableEntity) + " failed to get following list"}, []string{uuid})
		return
	}

	// todo: continue refactor here

	allFollowingIds := []int{}
	for rows.Next() {
		var id int
		err := rows.Scan(&id)
		if err != nil {
			log.Println("failed to scan following list", err.Error())
			wsSend(WS_ERROR_RESPONSE, WS_ERROR_RESPONSE_DTO{fmt.Sprint(http.StatusInternalServerError) + " failed to scan following list"}, []string{uuid})
			return
		}
		allFollowingIds = append(allFollowingIds, id)
	}

	var followings []UserForList
	for _, id := range allFollowingIds {
		var user UserForList
		err = statements["getUserbyID"].QueryRow(id).Scan(&user.Email, &user.First_name, &user.Last_name, &user.Nickname)
		if err != nil {
			log.Println("failed to scan into user for list", err.Error())
			wsSend(WS_ERROR_RESPONSE, WS_ERROR_RESPONSE_DTO{fmt.Sprint(http.StatusInternalServerError) + " failed to scan user for following list"}, []string{uuid})
			return
		}
		followings = append(followings, user)
	}

	wsSend(WS_USER_FOLLOWING_LIST, followings, []string{uuid})
}

/*
wsFollowersHandler returns list of users that are following the target user,

if user has permission to view the target user's followers list.

Otherwise, it returns an error:

"403 user does not have permissions to see the target user profile"
*/
func wsFollowersListHandler(conn *websocket.Conn, messageData map[string]interface{}) {
	defer wsRecover(messageData)

	uuid, ok := messageData["user_uuid"].(string)
	if !ok {
		log.Println("failed to get user_uuid from message data")
		return
	}
	user_id, err := get_user_id_by_uuid(uuid)
	if err != nil {
		log.Println("failed to get ID of the message sender", err.Error())
		wsSend(WS_ERROR_RESPONSE, WS_ERROR_RESPONSE_DTO{fmt.Sprint(http.StatusUnprocessableEntity) + " failed to get ID of the message sender"}, []string{uuid})
		return
	}

	target_email, ok := messageData["target_email"].(string)
	if !ok {
		log.Println("failed to get target email", err.Error())
		wsSend(WS_ERROR_RESPONSE, WS_ERROR_RESPONSE_DTO{fmt.Sprint(http.StatusUnprocessableEntity) + " failed to get target email"}, []string{uuid})
		return
	}

	target_id, err := get_user_id_by_email(target_email)
	if err != nil {
		log.Println("failed to get ID of the target user", err.Error())
		wsSend(WS_ERROR_RESPONSE, WS_ERROR_RESPONSE_DTO{fmt.Sprint(http.StatusUnprocessableEntity) + " failed to get ID of the target user"}, []string{uuid})
		return
	}

	isFollower, err := isFollowing(user_id, target_id)
	if err != nil {
		log.Println("failed to check if user is a follower of the target user", err.Error())
		wsSend(WS_ERROR_RESPONSE, WS_ERROR_RESPONSE_DTO{fmt.Sprint(http.StatusUnprocessableEntity) + " failed to check if user is a follower of the target user"}, []string{uuid})
		return
	}

	var profile WS_USER_PROFILE_DTO
	rows, err := statements["getUserProfile"].Query(target_id)
	if err != nil {
		log.Println("failed to get user profile", err.Error())
		wsSend(WS_ERROR_RESPONSE, WS_ERROR_RESPONSE_DTO{fmt.Sprint(http.StatusUnprocessableEntity) + " failed to get user profile"}, []string{uuid})
		return
	}
	defer rows.Close()
	rows.Next()
	err = rows.Scan(&profile.Email, &profile.First_name, &profile.Last_name, &profile.Dob,
		&profile.avatar_bytes, &profile.Nickname, &profile.About_me, &profile.Privacy)
	if err != nil {
		log.Println("failed to scan user profile", err.Error())
		wsSend(WS_ERROR_RESPONSE, WS_ERROR_RESPONSE_DTO{fmt.Sprint(http.StatusUnprocessableEntity) + " failed to scan user profile"}, []string{uuid})
		return
	}
	rows.Close()

	if target_id != user_id && !isFollower && profile.Privacy == "private" {
		log.Println("user does not have permissions to see the target user profile")
		wsSend(WS_ERROR_RESPONSE, WS_ERROR_RESPONSE_DTO{fmt.Sprint(http.StatusForbidden) + " user does not have permissions to see the target user profile"}, []string{uuid})
		return
	}

	rows, err = statements["getFollowers"].Query(target_id)
	if err != nil {
		log.Println("failed to get followers list", err.Error())
		wsSend(WS_ERROR_RESPONSE, WS_ERROR_RESPONSE_DTO{fmt.Sprint(http.StatusUnprocessableEntity) + " failed to get followers list"}, []string{uuid})
		return
	}

	allFollowersIds := []int{}
	for rows.Next() {
		var id int
		err := rows.Scan(&id)
		if err != nil {
			log.Println("failed to scan followers list", err.Error())
			wsSend(WS_ERROR_RESPONSE, WS_ERROR_RESPONSE_DTO{fmt.Sprint(http.StatusInternalServerError) + " failed to scan followers list"}, []string{uuid})
			return
		}
		allFollowersIds = append(allFollowersIds, id)
	}

	var followers []UserForList
	for _, id := range allFollowersIds {
		var user UserForList
		err = statements["getUserbyID"].QueryRow(id).Scan(&user.Email, &user.First_name, &user.Last_name, &user.Nickname)
		if err != nil {
			log.Println("failed to scan into user for list", err.Error())
			wsSend(WS_ERROR_RESPONSE, WS_ERROR_RESPONSE_DTO{fmt.Sprint(http.StatusInternalServerError) + " failed to scan user for following list"}, []string{uuid})
			return
		}
		followers = append(followers, user)
	}

	wsSend(WS_USER_FOLLOWERS_LIST, followers, []string{uuid})
}

// followHandler manages the current user follow target user request.
//
// - If: the target user is private, the current user will be added to the followers_pending table
//
// - Else: the current user will be added to the followers table
func wsFollowHandler(conn *websocket.Conn, messageData map[string]interface{}) {
	defer wsRecover(messageData)
	uuid, ok := messageData["user_uuid"].(string)
	if !ok {
		log.Println("failed to get user_uuid from message data")
		return
	}
	user_id, err := get_user_id_by_uuid(uuid)
	if err != nil {
		log.Println("failed to get ID of the message sender", err.Error())
		wsSend(WS_ERROR_RESPONSE, WS_ERROR_RESPONSE_DTO{fmt.Sprint(http.StatusUnprocessableEntity) + " failed to get ID of the message sender"}, []string{uuid})
		return
	}

	target_email, ok := messageData["target_email"].(string)
	if !ok {
		log.Println("failed to get target email", err.Error())
		wsSend(WS_ERROR_RESPONSE, WS_ERROR_RESPONSE_DTO{fmt.Sprint(http.StatusUnprocessableEntity) + " failed to get target email"}, []string{uuid})
		return
	}

	target_id, err := get_user_id_by_email(target_email)
	if err != nil {
		log.Println("failed to get ID of the target user", err.Error())
		wsSend(WS_ERROR_RESPONSE, WS_ERROR_RESPONSE_DTO{fmt.Sprint(http.StatusUnprocessableEntity) + " failed to get ID of the target user"}, []string{uuid})
		return
	}

	// check if our target whom we want to follow is private or public
	var privacy string
	err = statements["getUserPrivacy"].QueryRow(target_id).Scan(&privacy)
	if err != nil {
		log.Println("failed to get target user privacy", err.Error())
		wsSend(WS_ERROR_RESPONSE, WS_ERROR_RESPONSE_DTO{fmt.Sprint(http.StatusUnprocessableEntity) + " failed to get target user privacy"}, []string{uuid})
		return
	}

	// bug: add works even if the user is already a follower, perhaps can be fixed when profile is loaded
	if privacy == "private" {
		// if the target is private, add the follower to the followers_pending table
		_, err = statements["addFollowerPending"].Exec(target_id, user_id)
		if err != nil {
			log.Println("failed to add follower to followers_pending table", err.Error())
			wsSend(WS_ERROR_RESPONSE, WS_ERROR_RESPONSE_DTO{fmt.Sprint(http.StatusInternalServerError) + " failed to add follower to followers_pending table"}, []string{uuid})
			return
		}
		wsSend(WS_SUCCESS_RESPONSE, WS_SUCCESS_RESPONSE_DTO{fmt.Sprint(http.StatusOK) + "Request to become a follower was added"}, []string{uuid})
		return
	} else {
		// if the target is public, add the follower to the followers table
		_, err = statements["addFollower"].Exec(target_id, user_id)
		if err != nil {
			log.Println("failed to add follower to followers table", err.Error())
			wsSend(WS_ERROR_RESPONSE, WS_ERROR_RESPONSE_DTO{fmt.Sprint(http.StatusInternalServerError) + " failed to add follower to followers table"}, []string{uuid})
			return
		}
		wsSend(WS_SUCCESS_RESPONSE, WS_SUCCESS_RESPONSE_DTO{fmt.Sprint(http.StatusOK) + "Follower was added"}, []string{uuid})
		return
	}
}

// # unfollowHandler removes the current user from the followers list of the target user
//
// - If the followers table has a row with the current loggedin user as follower_id and the target user as user_id,
// the row will be deleted
func wsUnfollowHandler(conn *websocket.Conn, messageData map[string]interface{}) {
	defer wsRecover(messageData)

	uuid, ok := messageData["user_uuid"].(string)
	if !ok {
		log.Println("failed to get user_uuid from message data")
		return
	}
	user_id, err := get_user_id_by_uuid(uuid)
	if err != nil {
		log.Println("failed to get ID of the message sender", err.Error())
		wsSend(WS_ERROR_RESPONSE, WS_ERROR_RESPONSE_DTO{fmt.Sprint(http.StatusUnprocessableEntity) + " failed to get ID of the message sender"}, []string{uuid})
		return
	}

	target_email, ok := messageData["target_email"].(string)
	if !ok {
		log.Println("failed to get target email", err.Error())
		wsSend(WS_ERROR_RESPONSE, WS_ERROR_RESPONSE_DTO{fmt.Sprint(http.StatusUnprocessableEntity) + " failed to get target email"}, []string{uuid})
		return
	}

	target_id, err := get_user_id_by_email(target_email)
	if err != nil {
		log.Println("failed to get ID of the target user", err.Error())
		wsSend(WS_ERROR_RESPONSE, WS_ERROR_RESPONSE_DTO{fmt.Sprint(http.StatusUnprocessableEntity) + " failed to get ID of the target user"}, []string{uuid})
		return
	}

	_, err = statements["removeFollower"].Exec(target_id, user_id)
	if err != nil {
		log.Println("failed to remove follower from followers table", err.Error())
		wsSend(WS_ERROR_RESPONSE, WS_ERROR_RESPONSE_DTO{fmt.Sprint(http.StatusInternalServerError) + " failed to remove follower from followers table"}, []string{uuid})
		return
	}

	wsSend(WS_SUCCESS_RESPONSE, WS_SUCCESS_RESPONSE_DTO{fmt.Sprint(http.StatusOK) + " user unfollowed"}, []string{uuid})
}

// FollowRequestListHandler returns a list of users that have sent a follow request to the current user
func wsFollowRequestsListHandler(conn *websocket.Conn, messageData map[string]interface{}) {
	defer wsRecover(messageData)

	uuid, ok := messageData["user_uuid"].(string)
	if !ok {
		log.Println("failed to get user_uuid from message data")
		return
	}
	user_id, err := get_user_id_by_uuid(uuid)
	if err != nil {
		log.Println("failed to get ID of the message sender", err.Error())
		wsSend(WS_ERROR_RESPONSE, WS_ERROR_RESPONSE_DTO{fmt.Sprint(http.StatusUnprocessableEntity) + " failed to get ID of the message sender"}, []string{uuid})
		return
	}

	rows, err := statements["getFollowersPending"].Query(user_id)
	if err != nil {
		log.Println("failed to get followers pending", err.Error())
		wsSend(WS_ERROR_RESPONSE, WS_ERROR_RESPONSE_DTO{fmt.Sprint(http.StatusInternalServerError) + " failed to get followers pending"}, []string{uuid})
		return
	}
	var follow_requests_list []UserForList
	for rows.Next() {
		var follow_requester_id int
		var follow_requester UserForList
		err = rows.Scan(&follow_requester_id)
		if err != nil {
			log.Println("failed to scan followers pending", err.Error())
			wsSend(WS_ERROR_RESPONSE, WS_ERROR_RESPONSE_DTO{fmt.Sprint(http.StatusInternalServerError) + " failed to scan followers pending"}, []string{uuid})
			return
		}

		err = statements["getUserbyID"].QueryRow(follow_requester_id).Scan(&follow_requester.Email, &follow_requester.First_name, &follow_requester.Last_name, &follow_requester.Nickname)
		if err != nil {
			log.Println("failed to get user by ID", err.Error())
			wsSend(WS_ERROR_RESPONSE, WS_ERROR_RESPONSE_DTO{fmt.Sprint(http.StatusInternalServerError) + " failed to get user by ID"}, []string{uuid})
			return
		}
		follow_requests_list = append(follow_requests_list, follow_requester)
	}

	wsSend(WS_FOLLOW_REQUESTS_LIST, follow_requests_list, []string{uuid})
}
