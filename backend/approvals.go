package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type GroupMemberRequest struct {
	GroupID int    `json:"group_id"`
	Email   string `json:"member_email"`
}

type UserFollowerRequest struct {
	Email string `json:"email"` // the user who wants to follow you
}

// # groupRequestAcceptHandler is the handler for accepting a group request
//
// @rparam {group_id int, member_email string}
func groupRequestAcceptHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	defer recovery(w)

	LoggedinUserID, err := getRequestSenderID(r)
	if err != nil {
		jsonResponse(w, http.StatusUnauthorized, err.Error())
		return
	}

	data := GroupMemberRequest{}
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err = decoder.Decode(&data)
	if err != nil {
		log.Println(err.Error())
		jsonResponse(w, http.StatusBadRequest, "")
		return
	}

	memberID, err := getIDbyEmail(data.Email)
	if err != nil {
		jsonResponse(w, http.StatusBadRequest, "") // todo: fresh handling, in old version was just skipped
		return
	}

	rows, err := statements["getGroup"].Query(data.GroupID)
	defer rows.Close() // todo: it says this should be after error checking, but it is only warning
	if err != nil {
		log.Println(err.Error())
		jsonResponse(w, http.StatusInternalServerError, "getGroup query failed")
		return
	}
	group := Group{}

	for rows.Next() {
		err = rows.Scan(&group.ID, &group.Name, &group.Description, &group.CreatorId, &group.CreationDate, &group.Privacy)
		if err != nil {
			log.Println(err.Error())
			jsonResponse(w, http.StatusInternalServerError, "failed to scan group")
			return
		}
	}

	if group.CreatorId != LoggedinUserID {
		jsonResponse(w, http.StatusUnauthorized, "you are not the creator of the group")
		return
	}

	_, err = statements["addGroupMember"].Exec(data.GroupID, memberID)
	if err != nil {
		log.Println(err.Error())
		jsonResponse(w, http.StatusInternalServerError, "addGroupMember query failed")
		return
	}

	_, err = statements["removeGroupPendingMember"].Exec(data.GroupID, memberID)
	if err != nil {
		log.Println(err.Error())
		jsonResponse(w, http.StatusInternalServerError, "removeGroupPendingMember query failed")
		return
	}

	jsonResponse(w, http.StatusOK, "success: you approved the group membership")
	return
}

// # groupRequestRejectHandler is the handler for rejecting a group request
//
// @rparam {group_id int, member_email string}
func groupRequestRejectHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	defer recovery(w)

	LoggedinUserID, err := getRequestSenderID(r)
	if err != nil {
		jsonResponse(w, http.StatusUnauthorized, err.Error())
		return
	}

	// incoming DTO GroupMember
	data := GroupMemberRequest{}
	// decode the request body into the DTO
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err = decoder.Decode(&data)
	if err != nil {
		log.Println(err.Error())
		jsonResponse(w, http.StatusBadRequest, "")
		return
	}

	memberID, err := getIDbyEmail(data.Email)
	if err != nil {
		jsonResponse(w, http.StatusBadRequest, "") // todo: fresh handling, in old version was just skipped
		return
	}

	rows, err := statements["getGroup"].Query(data.GroupID)
	defer rows.Close() // todo: it says this should be after error checking, but it is only warning
	if err != nil {
		log.Println(err.Error())
		jsonResponse(w, http.StatusInternalServerError, "getGroup query failed")
		return
	}
	group := Group{}

	for rows.Next() {
		err = rows.Scan(&group.ID, &group.Name, &group.Description, &group.CreatorId, &group.CreationDate, &group.Privacy)
		if err != nil {
			log.Println(err.Error())
			jsonResponse(w, http.StatusInternalServerError, "failed to scan group")
			return
		}
	}

	if group.CreatorId != LoggedinUserID {
		jsonResponse(w, http.StatusUnauthorized, "you are not the creator of the group")
		return
	}

	_, err = statements["removeGroupPendingMember"].Exec(data.GroupID, memberID)
	if err != nil {
		log.Println(err.Error())
		jsonResponse(w, http.StatusInternalServerError, "removeGroupPendingMember query failed")
		return
	}

	jsonResponse(w, http.StatusOK, "success: you rejected the group membership")
	return
}

// # approveFollowerHandler is the handler for accepting a follower request
//
// @rparam {email string}
func approveFollowerHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	defer recovery(w)

	LoggedinUserID, err := getRequestSenderID(r)
	if err != nil {
		jsonResponse(w, http.StatusUnauthorized, err.Error())
		return
	}

	// incoming DTO UserFollowerRequest
	data := UserFollowerRequest{}
	// decode the request body into the DTO
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err = decoder.Decode(&data)
	if err != nil {
		log.Println(err.Error())
		jsonResponse(w, http.StatusBadRequest, "")
		return
	}

	fanID, err := getIDbyEmail(data.Email)
	if err != nil {
		jsonResponse(w, http.StatusBadRequest, "") // todo: fresh handling, in old version was just skipped
		return
	}

	// check if the user has the follower in the followers_pending table
	// use the getFollowersPending statement
	rows, err := statements["getFollowersPending"].Query(LoggedinUserID)
	defer rows.Close()
	if err != nil {
		log.Println(err.Error())
		jsonResponse(w, http.StatusInternalServerError, "getFollowersPending query failed")
		return
	}
	var followerID int
	var followersPending []int
	for rows.Next() {
		err = rows.Scan(&followerID)
		if err != nil {
			log.Println(err.Error())
			jsonResponse(w, http.StatusInternalServerError, "failed to scan followers pending")
			return
		}
		followersPending = append(followersPending, followerID)
	}
	// check if the user is in the list of followersPending request
	found := false
	for _, follower := range followersPending {
		if follower == fanID {
			found = true
			break
		}
	}

	if !found {
		jsonResponse(w, http.StatusBadRequest, "the user is not in the list of followers pending")
		return
	}
	// add the follower to the followers table
	_, err = statements["addFollower"].Exec(LoggedinUserID, fanID)
	if err != nil {
		log.Println(err.Error())
		jsonResponse(w, http.StatusInternalServerError, "addFollower query failed")
		return
	}
	// remove it from the followers_pending table
	_, err = statements["removeFollowerPending"].Exec(LoggedinUserID, fanID)
	if err != nil {
		log.Println(err.Error())
		jsonResponse(w, http.StatusInternalServerError, "removeFollowerPending query failed")
		return // todo: CHECK IT! Was added , before it was just continue of the flow, to next line
	}
	jsonResponse(w, http.StatusOK, "success: you accepted the follow request")
	return
}

// # rejectFollowerHandler is the handler for rejecting a follower request
//
// @rparam {email string}
func rejectFollowerHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	defer recovery(w)

	LoggedinUserID, err := getRequestSenderID(r)
	if err != nil {
		jsonResponse(w, http.StatusUnauthorized, err.Error())
		return
	}

	// incoming DTO UserFollowerRequest
	data := UserFollowerRequest{}
	// decode the request body into the DTO
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err = decoder.Decode(&data)
	if err != nil {
		log.Println(err.Error())
		jsonResponse(w, http.StatusBadRequest, "")
		return
	}

	fanID, err := getIDbyEmail(data.Email)
	if err != nil {
		jsonResponse(w, http.StatusBadRequest, "") // todo: fresh handling, in old version was just skipped
		return
	}

	// check if the user has the follower in the followers_pending table
	// use the getFollowersPending statement
	rows, err := statements["getFollowersPending"].Query(LoggedinUserID)
	defer rows.Close()
	if err != nil {
		log.Println(err.Error())
		jsonResponse(w, http.StatusInternalServerError, "getFollowersPending query failed")
		return
	}
	var followerID int
	var followersPending []int
	for rows.Next() {
		err = rows.Scan(&followerID)
		if err != nil {
			log.Println(err.Error())
			jsonResponse(w, http.StatusInternalServerError, "failed to scan followers pending")
			return
		}
		followersPending = append(followersPending, followerID)
	}
	// check if the user is in the list of followersPending request
	found := false
	for _, follower := range followersPending {
		if follower == fanID {
			found = true
			break
		}
	}

	if !found {
		jsonResponse(w, http.StatusBadRequest, "the user is not in the list of followers pending")
		return
	}

	// remove it from the followers_pending table
	_, err = statements["removeFollowerPending"].Exec(LoggedinUserID, fanID)
	if err != nil {
		log.Println(err.Error())
		jsonResponse(w, http.StatusInternalServerError, "removeFollowerPending query failed")
		return // todo: CHECK IT! Was added , before it was just continue of the flow, to next line
	}
	jsonResponse(w, http.StatusOK, "success: you rejected the follow request")
	return
}

// groupInviteAcceptHandler is the handler for accepting a group invite
//
// @r.param {group_id int}
func groupInviteAcceptHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	defer recovery(w)

	requestorId, err := getRequestSenderID(r)
	if err != nil {
		jsonResponse(w, http.StatusUnauthorized, err.Error())
		return
	}

	var data struct {
		GroupId int `json:"group_id"`
	}

	// get the group id from the request body
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err = decoder.Decode(&data)
	if err != nil {
		log.Println(err.Error())
		jsonResponse(w, http.StatusBadRequest, "failed to get group id from request body")
		return
	}

	// add the person to the group_members table
	_, err = statements["addGroupMember"].Exec(data.GroupId, requestorId)
	if err != nil {
		log.Println(err.Error())
		jsonResponse(w, http.StatusInternalServerError, "failed to add person to group_members table")
		return
	}

	// remove the person from the group_invited_users table
	_, err = statements["removeGroupInvitedUser"].Exec(requestorId, data.GroupId)
	if err != nil {
		log.Println(err.Error())
		jsonResponse(w, http.StatusInternalServerError, "failed to remove person from group_invited_users table")
		return
	}
	jsonResponse(w, http.StatusOK, "success: you accepted the group invite")
	return
}

// groupInviteRejectHandler is the handler for rejecting a group invite
//
// @r.param {group_id int}
func groupInviteRejectHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	defer recovery(w)

	requestorId, err := getRequestSenderID(r)
	if err != nil {
		jsonResponse(w, http.StatusUnauthorized, err.Error())
		return
	}

	var data struct {
		GroupId int `json:"group_id"`
	}

	// get the group id from the request body
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err = decoder.Decode(&data)
	if err != nil {
		log.Println(err.Error())
		jsonResponse(w, http.StatusBadRequest, "failed to get group id from request body")
		return
	}

	// remove the person from the group_invited_users table
	_, err = statements["removeGroupInvitedUser"].Exec(requestorId, data.GroupId)
	if err != nil {
		log.Println(err.Error())
		jsonResponse(w, http.StatusInternalServerError, "failed to remove person from group_invited_users table")
		return
	}
	jsonResponse(w, http.StatusOK, "success: you rejected the group invite")
	return
}
