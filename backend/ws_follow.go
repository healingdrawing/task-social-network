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
	defer wsRecover()

	uuid := messageData["user_uuid"].(string)
	user_id, err := getIDbyUUID(uuid)
	if err != nil {
		log.Println("failed to get ID of the request sender", err.Error())
		wsSendError(WS_ERROR_RESPONSE_DTO{fmt.Sprint(http.StatusUnprocessableEntity) + " failed to get ID of the request sender"})
		return
	}

	target_email, ok := messageData["target_email"].(string)
	if !ok {
		log.Println("failed to get target email", err.Error())
		wsSendError(WS_ERROR_RESPONSE_DTO{fmt.Sprint(http.StatusUnprocessableEntity) + " failed to get target email"})
		return
	}

	target_id, err := getIDbyEmail(target_email)
	if err != nil {
		log.Println("failed to get ID of the target user", err.Error())
		wsSendError(WS_ERROR_RESPONSE_DTO{fmt.Sprint(http.StatusUnprocessableEntity) + " failed to get ID of the target user"})
		return
	}

	isFollower, err := isFollowing(target_id, user_id)
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

	if target_id != user_id && !isFollower && profile.Privacy == "private" {
		log.Println("user does not have permissions to see the target user profile")
		wsSendError(WS_ERROR_RESPONSE_DTO{fmt.Sprint(http.StatusForbidden) + " user does not have permissions to see the target user profile"})
		return
	}

	rows, err = statements["getFollowing"].Query(target_id)
	if err != nil {
		log.Println("failed to get following list", err.Error())
		wsSendError(WS_ERROR_RESPONSE_DTO{fmt.Sprint(http.StatusUnprocessableEntity) + " failed to get following list"})
		return
	}

	// todo: continue refactor here

	allFollowingIds := []int{}
	for rows.Next() {
		var id int
		err := rows.Scan(&id)
		if err != nil {
			log.Println("failed to scan following list", err.Error())
			wsSendError(WS_ERROR_RESPONSE_DTO{fmt.Sprint(http.StatusInternalServerError) + " failed to scan following list"})
			return
		}
		allFollowingIds = append(allFollowingIds, id)
	}

	var followings []UserForList
	for _, id := range allFollowingIds {
		var user UserForList
		err = statements["getUserbyID"].QueryRow(id).Scan(&user.Email, &user.First_name, &user.Last_name)
		if err != nil {
			log.Println("failed to scan into user for list", err.Error())
			wsSendError(WS_ERROR_RESPONSE_DTO{fmt.Sprint(http.StatusInternalServerError) + " failed to scan user for following list"})
			return
		}
		followings = append(followings, user)
	}

	wsSendFollowingList(followings)
}

/*
wsFollowersHandler returns list of users that are following the target user,

if user has permission to view the target user's followers list.

Otherwise, it returns an error:

"403 user does not have permissions to see the target user profile"
*/
func wsFollowersListHandler(conn *websocket.Conn, messageData map[string]interface{}) {
	defer wsRecover()

	uuid := messageData["user_uuid"].(string)
	user_id, err := getIDbyUUID(uuid)
	if err != nil {
		log.Println("failed to get ID of the request sender", err.Error())
		wsSendError(WS_ERROR_RESPONSE_DTO{fmt.Sprint(http.StatusUnprocessableEntity) + " failed to get ID of the request sender"})
		return
	}

	target_email, ok := messageData["target_email"].(string)
	if !ok {
		log.Println("failed to get target email", err.Error())
		wsSendError(WS_ERROR_RESPONSE_DTO{fmt.Sprint(http.StatusUnprocessableEntity) + " failed to get target email"})
		return
	}

	target_id, err := getIDbyEmail(target_email)
	if err != nil {
		log.Println("failed to get ID of the target user", err.Error())
		wsSendError(WS_ERROR_RESPONSE_DTO{fmt.Sprint(http.StatusUnprocessableEntity) + " failed to get ID of the target user"})
		return
	}

	isFollower, err := isFollowing(target_id, user_id)
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

	if target_id != user_id && !isFollower && profile.Privacy == "private" {
		log.Println("user does not have permissions to see the target user profile")
		wsSendError(WS_ERROR_RESPONSE_DTO{fmt.Sprint(http.StatusForbidden) + " user does not have permissions to see the target user profile"})
		return
	}

	rows, err = statements["getFollowers"].Query(target_id)
	if err != nil {
		log.Println("failed to get followers list", err.Error())
		wsSendError(WS_ERROR_RESPONSE_DTO{fmt.Sprint(http.StatusUnprocessableEntity) + " failed to get followers list"})
		return
	}

	allFollowersIds := []int{}
	for rows.Next() {
		var id int
		err := rows.Scan(&id)
		if err != nil {
			log.Println("failed to scan followers list", err.Error())
			wsSendError(WS_ERROR_RESPONSE_DTO{fmt.Sprint(http.StatusInternalServerError) + " failed to scan followers list"})
			return
		}
		allFollowersIds = append(allFollowersIds, id)
	}

	var followers []UserForList
	for _, id := range allFollowersIds {
		var user UserForList
		err = statements["getUserbyID"].QueryRow(id).Scan(&user.Email, &user.First_name, &user.Last_name)
		if err != nil {
			log.Println("failed to scan into user for list", err.Error())
			wsSendError(WS_ERROR_RESPONSE_DTO{fmt.Sprint(http.StatusInternalServerError) + " failed to scan user for following list"})
			return
		}
		followers = append(followers, user)
	}

	wsSendFollowersList(followers)
}

// followHandler manages the current user follow target user request.
//
// - If: the target user is private, the current user will be added to the followers_pending table
//
// - Else: the current user will be added to the followers table
func wsFollowHandler(conn *websocket.Conn, messageData map[string]interface{}) {
	defer wsRecover()
	uuid := messageData["user_uuid"].(string)
	user_id, err := getIDbyUUID(uuid)
	if err != nil {
		log.Println("failed to get ID of the request sender", err.Error())
		wsSendError(WS_ERROR_RESPONSE_DTO{fmt.Sprint(http.StatusUnprocessableEntity) + " failed to get ID of the request sender"})
		return
	}

	target_email, ok := messageData["target_email"].(string)
	if !ok {
		log.Println("failed to get target email", err.Error())
		wsSendError(WS_ERROR_RESPONSE_DTO{fmt.Sprint(http.StatusUnprocessableEntity) + " failed to get target email"})
		return
	}

	target_id, err := getIDbyEmail(target_email)
	if err != nil {
		log.Println("failed to get ID of the target user", err.Error())
		wsSendError(WS_ERROR_RESPONSE_DTO{fmt.Sprint(http.StatusUnprocessableEntity) + " failed to get ID of the target user"})
		return
	}

	// check if our target whom we want to follow is private or public
	var privacy string
	err = statements["getUserPrivacy"].QueryRow(target_id).Scan(&privacy)
	if err != nil {
		log.Println("failed to get target user privacy", err.Error())
		wsSendError(WS_ERROR_RESPONSE_DTO{fmt.Sprint(http.StatusUnprocessableEntity) + " failed to get target user privacy"})
		return
	}

	if privacy == "private" {
		// if the target is private, add the follower to the followers_pending table
		_, err = statements["addFollowerPending"].Exec(target_id, user_id)
		if err != nil {
			log.Println("failed to add follower to followers_pending table", err.Error())
			wsSendError(WS_ERROR_RESPONSE_DTO{fmt.Sprint(http.StatusInternalServerError) + " failed to add follower to followers_pending table"})
			return
		}
		wsSendSuccess(WS_SUCCESS_RESPONSE_DTO{"Request to become a follower was added"})
		return
	} else {
		// if the target is public, add the follower to the followers table
		_, err = statements["addFollower"].Exec(target_id, user_id)
		if err != nil {
			log.Println("failed to add follower to followers table", err.Error())
			wsSendError(WS_ERROR_RESPONSE_DTO{fmt.Sprint(http.StatusInternalServerError) + " failed to add follower to followers table"})
			return
		}
		wsSendSuccess(WS_SUCCESS_RESPONSE_DTO{"Follower was added"})
		return
	}
}

// # unfollowHandler removes the current user from the followers list of the target user
//
// - If the followers table has a row with the current loggedin user as follower_id and the target user as user_id,
// the row will be deleted
func wsUnfollowHandler(conn *websocket.Conn, messageData map[string]interface{}) {
	defer wsRecover()

	uuid := messageData["user_uuid"].(string)
	user_id, err := getIDbyUUID(uuid)
	if err != nil {
		log.Println("failed to get ID of the request sender", err.Error())
		wsSendError(WS_ERROR_RESPONSE_DTO{fmt.Sprint(http.StatusUnprocessableEntity) + " failed to get ID of the request sender"})
		return
	}

	target_email, ok := messageData["target_email"].(string)
	if !ok {
		log.Println("failed to get target email", err.Error())
		wsSendError(WS_ERROR_RESPONSE_DTO{fmt.Sprint(http.StatusUnprocessableEntity) + " failed to get target email"})
		return
	}

	target_id, err := getIDbyEmail(target_email)
	if err != nil {
		log.Println("failed to get ID of the target user", err.Error())
		wsSendError(WS_ERROR_RESPONSE_DTO{fmt.Sprint(http.StatusUnprocessableEntity) + " failed to get ID of the target user"})
		return
	}

	_, err = statements["removeFollower"].Exec(target_id, user_id)
	if err != nil {
		log.Println("failed to remove follower from followers table", err.Error())
		wsSendError(WS_ERROR_RESPONSE_DTO{fmt.Sprint(http.StatusInternalServerError) + " failed to remove follower from followers table"})
		return
	}
	wsSendSuccess(WS_SUCCESS_RESPONSE_DTO{fmt.Sprint(http.StatusOK) + " user unfollowed"})
	return
}

// FollowRequestListHandler returns a list of users that have sent a follow request to the current user
func wsFollowRequestsListHandler(conn *websocket.Conn, messageData map[string]interface{}) {
	defer wsRecover()

	uuid := messageData["user_uuid"].(string)
	user_id, err := getIDbyUUID(uuid)
	if err != nil {
		log.Println("failed to get ID of the request sender", err.Error())
		wsSendError(WS_ERROR_RESPONSE_DTO{fmt.Sprint(http.StatusUnprocessableEntity) + " failed to get ID of the request sender"})
		return
	}

	rows, err := statements["getFollowersPending"].Query(user_id)
	if err != nil {
		log.Println("failed to get followers pending", err.Error())
		wsSendError(WS_ERROR_RESPONSE_DTO{fmt.Sprint(http.StatusInternalServerError) + " failed to get followers pending"})
		return
	}
	var follow_requests_list []UserForList
	for rows.Next() {
		var follow_requester_id int
		var follow_requester UserForList
		err = rows.Scan(&follow_requester_id)
		if err != nil {
			log.Println("failed to scan followers pending", err.Error())
			wsSendError(WS_ERROR_RESPONSE_DTO{fmt.Sprint(http.StatusInternalServerError) + " failed to scan followers pending"})
			return
		}

		err = statements["getUserbyID"].QueryRow(follow_requester_id).Scan(&follow_requester.Email, &follow_requester.First_name, &follow_requester.Last_name)
		if err != nil {
			log.Println("failed to get user by ID", err.Error())
			wsSendError(WS_ERROR_RESPONSE_DTO{fmt.Sprint(http.StatusInternalServerError) + " failed to get user by ID"})
			return
		}
		follow_requests_list = append(follow_requests_list, follow_requester)
	}

	wsSendFollowRequestsList(follow_requests_list)
}