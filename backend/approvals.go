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
	defer func() {
		if err := recover(); err != nil {
			log.Println(err)
			jsonResponseWriterManager(w, http.StatusInternalServerError, "")
		}
	}()
	cookie, err := r.Cookie("user_uuid")
	if err != nil {
		if err == http.ErrNoCookie {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	uuid := cookie.Value
	data := GroupMemberRequest{}
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err = decoder.Decode(&data)
	if err != nil {
		log.Println(err.Error())
		jsonResponseWriterManager(w, http.StatusBadRequest, "")
		return
	}
	LoggedinUserID, err := getIDbyUUID(uuid)
	if err != nil {
		jsonResponseWriterManager(w, http.StatusUnauthorized, "")
		return
	}

	memberID, err := getIDbyEmail(data.Email)
	if err != nil {
		jsonResponseWriterManager(w, http.StatusBadRequest, "") // todo: fresh handling, in old version was just skipped
		return
	}

	rows, err := statements["getGroup"].Query(data.GroupID)
	defer rows.Close() // todo: it says this should be after error checking, but it is only warning
	if err != nil {
		log.Println(err.Error())
		jsonResponseWriterManager(w, http.StatusInternalServerError, "getGroup query failed")
		return
	}
	group := Group{}

	for rows.Next() {
		err = rows.Scan(&group.ID, &group.Name, &group.Description, &group.CreatorId, &group.CreationDate, &group.Privacy)
		if err != nil {
			log.Println(err.Error())
			jsonResponseWriterManager(w, http.StatusInternalServerError, "failed to scan group")
			return
		}
	}

	if group.CreatorId != LoggedinUserID {
		jsonResponseWriterManager(w, http.StatusUnauthorized, "you are not the creator of the group")
		return
	}

	_, err = statements["addGroupMember"].Exec(data.GroupID, memberID)
	if err != nil {
		log.Println(err.Error())
		jsonResponseWriterManager(w, http.StatusInternalServerError, "addGroupMember query failed")
		return
	}

	// remove it from the group_pending_members table //todo: "remove it ..." remove what? :D
	_, err = statements["removeGroupPendingMember"].Exec(data.GroupID, memberID)
	if err != nil {
		log.Println(err.Error())
		jsonResponseWriterManager(w, http.StatusInternalServerError, "removeGroupPendingMember query failed")
		return
	}

	jsonResponseWriterManager(w, http.StatusOK, "success: you approved the group membership")
	return
}

// # groupRequestRejectHandler is the handler for rejecting a group request
//
// @rparam {group_id int, member_email string}
func groupRequestRejectHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	defer func() {
		if err := recover(); err != nil {
			log.Println(err)
			jsonResponseWriterManager(w, http.StatusInternalServerError, "")
		}
	}()
	// get the logged in user id from the uuid in cookies
	cookie, err := r.Cookie("user_uuid")
	if err != nil {
		if err == http.ErrNoCookie {
			// If the cookie is not set, return an unauthorized status
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		// For any other type of error, return a bad request status
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Get the UUID value from the cookie
	uuid := cookie.Value
	// incoming DTO GroupMember
	data := GroupMemberRequest{}
	// decode the request body into the DTO
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err = decoder.Decode(&data)
	if err != nil {
		log.Println(err.Error())
		jsonResponseWriterManager(w, http.StatusBadRequest, "")
		return
	}
	LoggedinUserID, err := getIDbyUUID(uuid)
	if err != nil {
		jsonResponseWriterManager(w, http.StatusUnauthorized, "")
		return
	}

	// from member_email get the member_id
	memberID, err := getIDbyEmail(data.Email)
	if err != nil {
		jsonResponseWriterManager(w, http.StatusBadRequest, "") // todo: fresh handling, in old version was just skipped
		return
	}

	// check if the user is the creator of the group
	// use the getGroup statement
	rows, err := statements["getGroup"].Query(data.GroupID)
	defer rows.Close() // todo: it says this should be after error checking, but it is only warning
	if err != nil {
		log.Println(err.Error())
		jsonResponseWriterManager(w, http.StatusInternalServerError, "getGroup query failed")
		return
	}
	group := Group{}

	for rows.Next() {
		err = rows.Scan(&group.ID, &group.Name, &group.Description, &group.CreatorId, &group.CreationDate, &group.Privacy)
		if err != nil {
			log.Println(err.Error())
			jsonResponseWriterManager(w, http.StatusInternalServerError, "failed to scan group")
			return
		}
	}

	if group.CreatorId != LoggedinUserID {
		jsonResponseWriterManager(w, http.StatusUnauthorized, "you are not the creator of the group")
		return
	}

	// remove it from the group_pending_members table
	_, err = statements["removeGroupPendingMember"].Exec(data.GroupID, memberID)
	if err != nil {
		log.Println(err.Error())
		jsonResponseWriterManager(w, http.StatusInternalServerError, "removeGroupPendingMember query failed")
		return
	}

	jsonResponseWriterManager(w, http.StatusOK, "success: you rejected the group membership")
	return
}

func approveFollowerHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	defer func() {
		if err := recover(); err != nil {
			log.Println(err)
			jsonResponseWriterManager(w, http.StatusInternalServerError, "")
		}
	}()
	// get the logged in user id from the uuid in cookies
	cookie, err := r.Cookie("user_uuid")
	if err != nil {
		if err == http.ErrNoCookie {
			// If the cookie is not set, return an unauthorized status
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		// For any other type of error, return a bad request status
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Get the UUID value from the cookie
	uuid := cookie.Value

	// incoming DTO UserFollowerRequest
	data := UserFollowerRequest{}
	// decode the request body into the DTO
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err = decoder.Decode(&data)
	if err != nil {
		log.Println(err.Error())
		jsonResponseWriterManager(w, http.StatusBadRequest, "")
		return
	}

	fanID, err := getIDbyEmail(data.Email)
	if err != nil {
		jsonResponseWriterManager(w, http.StatusBadRequest, "") // todo: fresh handling, in old version was just skipped
		return
	}

	LoggedinUserID, err := getIDbyUUID(uuid)
	if err != nil {
		jsonResponseWriterManager(w, http.StatusUnauthorized, "")
		// w.WriteHeader(http.StatusUnauthorized) // todo: CHECK IT, WAS REPLACED BY THE LINE ABOVE
		return
	}
	// check if the user has the follower in the followers_pending table
	// use the getFollowersPending statement
	rows, err := statements["getFollowersPending"].Query(LoggedinUserID)
	defer rows.Close()
	if err != nil {
		log.Println(err.Error())
		jsonResponseWriterManager(w, http.StatusInternalServerError, "getFollowersPending query failed")
		return
	}
	var followerID int
	var followersPending []int
	for rows.Next() {
		err = rows.Scan(&followerID)
		if err != nil {
			log.Println(err.Error())
			jsonResponseWriterManager(w, http.StatusInternalServerError, "failed to scan followers pending")
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
		jsonResponseWriterManager(w, http.StatusBadRequest, "the user is not in the list of followers pending")
		return
	}
	// add the follower to the followers table
	_, err = statements["addFollower"].Exec(LoggedinUserID, fanID)
	if err != nil {
		log.Println(err.Error())
		jsonResponseWriterManager(w, http.StatusInternalServerError, "addFollower query failed")
		return
	}
	// remove it from the followers_pending table
	_, err = statements["removeFollowerPending"].Exec(LoggedinUserID, fanID)
	if err != nil {
		log.Println(err.Error())
		jsonResponseWriterManager(w, http.StatusInternalServerError, "removeFollowerPending query failed")
		return // todo: CHECK IT! Was added , before it was just continue of the flow, to next line
	}
	jsonResponseWriterManager(w, http.StatusOK, "success: you accepted the follow request")
	return
}

func rejectFollowerHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	defer func() {
		if err := recover(); err != nil {
			log.Println(err)
			jsonResponseWriterManager(w, http.StatusInternalServerError, "")
		}
	}()
	// get the logged in user id from the uuid in cookies
	cookie, err := r.Cookie("user_uuid")
	if err != nil {
		if err == http.ErrNoCookie {
			// If the cookie is not set, return an unauthorized status
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		// For any other type of error, return a bad request status
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Get the UUID value from the cookie
	uuid := cookie.Value

	// incoming DTO UserFollowerRequest
	data := UserFollowerRequest{}
	// decode the request body into the DTO
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err = decoder.Decode(&data)
	if err != nil {
		log.Println(err.Error())
		jsonResponseWriterManager(w, http.StatusBadRequest, "")
		return
	}

	fanID, err := getIDbyEmail(data.Email)
	if err != nil {
		jsonResponseWriterManager(w, http.StatusBadRequest, "") // todo: fresh handling, in old version was just skipped
		return
	}

	LoggedinUserID, err := getIDbyUUID(uuid)
	if err != nil {
		jsonResponseWriterManager(w, http.StatusUnauthorized, "")
		// w.WriteHeader(http.StatusUnauthorized) // todo: CHECK IT, WAS REPLACED BY THE LINE ABOVE
		return
	}

	// check if the user has the follower in the followers_pending table
	// use the getFollowersPending statement
	rows, err := statements["getFollowersPending"].Query(LoggedinUserID)
	defer rows.Close()
	if err != nil {
		log.Println(err.Error())
		jsonResponseWriterManager(w, http.StatusInternalServerError, "getFollowersPending query failed")
		return
	}
	var followerID int
	var followersPending []int
	for rows.Next() {
		err = rows.Scan(&followerID)
		if err != nil {
			log.Println(err.Error())
			jsonResponseWriterManager(w, http.StatusInternalServerError, "failed to scan followers pending")
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
		jsonResponseWriterManager(w, http.StatusBadRequest, "the user is not in the list of followers pending")
		return
	}

	// remove it from the followers_pending table
	_, err = statements["removeFollowerPending"].Exec(LoggedinUserID, fanID)
	if err != nil {
		log.Println(err.Error())
		jsonResponseWriterManager(w, http.StatusInternalServerError, "removeFollowerPending query failed")
		return // todo: CHECK IT! Was added , before it was just continue of the flow, to next line
	}
	jsonResponseWriterManager(w, http.StatusOK, "success: you rejected the follow request")
	return
}

// groupInviteAcceptHandler is the handler for accepting a group invite
//
// @r.param {group_id int}
func groupInviteAcceptHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	defer func() {
		if err := recover(); err != nil {
			log.Println(err)
			jsonResponseWriterManager(w, http.StatusInternalServerError, "")
		}
	}()

	// get the id of the request sender
	cookie, err := r.Cookie("user_uuid")
	if err != nil {
		jsonResponseWriterManager(w, http.StatusUnauthorized, "")
		return
	}
	requestorId, err := getIDbyUUID(cookie.Value)
	if err != nil {
		log.Println(err.Error())
		jsonResponseWriterManager(w, http.StatusInternalServerError, "failed to get id of request sender")
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
		jsonResponseWriterManager(w, http.StatusBadRequest, "failed to get group id from request body")
		return
	}

	// add the person to the group_members table
	_, err = statements["addGroupMember"].Exec(data.GroupId, requestorId)
	if err != nil {
		log.Println(err.Error())
		jsonResponseWriterManager(w, http.StatusInternalServerError, "failed to add person to group_members table")
		return
	}

	// remove the person from the group_invited_users table
	_, err = statements["removeGroupInvitedUser"].Exec(requestorId, data.GroupId)
	if err != nil {
		log.Println(err.Error())
		jsonResponseWriterManager(w, http.StatusInternalServerError, "failed to remove person from group_invited_users table")
		return
	}
	jsonResponseWriterManager(w, http.StatusOK, "success: you accepted the group invite")
	return
}

// groupInviteRejectHandler is the handler for rejecting a group invite
//
// @r.param {group_id int}
func groupInviteRejectHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	defer func() {
		if err := recover(); err != nil {
			log.Println(err)
			jsonResponseWriterManager(w, http.StatusInternalServerError, "")
		}
	}()

	// get the id of the request sender
	cookie, err := r.Cookie("user_uuid")
	if err != nil {
		jsonResponseWriterManager(w, http.StatusUnauthorized, "")
		return
	}
	requestorId, err := getIDbyUUID(cookie.Value)
	if err != nil {
		log.Println(err.Error())
		jsonResponseWriterManager(w, http.StatusInternalServerError, "failed to get id of request sender")
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
		jsonResponseWriterManager(w, http.StatusBadRequest, "failed to get group id from request body")
		return
	}

	// remove the person from the group_invited_users table
	_, err = statements["removeGroupInvitedUser"].Exec(requestorId, data.GroupId)
	if err != nil {
		log.Println(err.Error())
		jsonResponseWriterManager(w, http.StatusInternalServerError, "failed to remove person from group_invited_users table")
		return
	}
	jsonResponseWriterManager(w, http.StatusOK, "success: you rejected the group invite")
	return
}
