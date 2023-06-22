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
			w.WriteHeader(500)
			jsonResponse, _ := json.Marshal(map[string]string{
				"message": "internal server error",
			})
			w.Write(jsonResponse)
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
		w.WriteHeader(400)
		jsonResponse, _ := json.Marshal(map[string]string{
			"message": "bad request",
		})
		w.Write(jsonResponse)
		return
	}
	LoggedinUserID, err := getIDbyUUID(uuid)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		jsonResponse, _ := json.Marshal(map[string]string{
			"message": "unauthorized",
		})
		w.Write(jsonResponse)
		return
	}

	memberID, err := getIDbyEmail(data.Email)

	rows, err := statements["getGroup"].Query(data.GroupID)
	defer rows.Close()
	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(500)
		jsonResponse, _ := json.Marshal(map[string]string{
			"message": "internal server error, getGroup query failed",
		})
		w.Write(jsonResponse)
		return
	}
	group := Group{}

	for rows.Next() {
		err = rows.Scan(&group.ID, &group.Name, &group.Description, &group.CreatorId, &group.CreationDate, &group.Privacy)
		if err != nil {
			log.Println(err.Error())
			w.WriteHeader(500)
			jsonResponse, _ := json.Marshal(map[string]string{
				"message": "internal server error, failed to scan group",
			})
			w.Write(jsonResponse)
			return
		}
	}

	if group.CreatorId != LoggedinUserID {
		w.WriteHeader(http.StatusUnauthorized)
		jsonResponse, _ := json.Marshal(map[string]string{
			"message": "unauthorized, you are not the creator of the group",
		})
		w.Write(jsonResponse)
		return
	}

	_, err = statements["addGroupMember"].Exec(data.GroupID, memberID)
	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(500)
		jsonResponse, _ := json.Marshal(map[string]string{
			"message": "internal server error, addGroupMember query failed",
		})
		w.Write(jsonResponse)
		return
	}

	// remove it from the group_pending_members table
	_, err = statements["removeGroupPendingMember"].Exec(data.GroupID, memberID)
	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(500)
		jsonResponse, _ := json.Marshal(map[string]string{
			"message": "internal server error, removeGroupPendingMember query failed",
		})
		w.Write(jsonResponse)
		return
	}

	w.WriteHeader(http.StatusOK)
	jsonResponse, _ := json.Marshal(map[string]string{
		"message": "success: you approved the group membership",
	})
	w.Write(jsonResponse)
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
			w.WriteHeader(500)
			jsonResponse, _ := json.Marshal(map[string]string{
				"message": "internal server error",
			})
			w.Write(jsonResponse)
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
		w.WriteHeader(400)
		jsonResponse, _ := json.Marshal(map[string]string{
			"message": "bad request",
		})
		w.Write(jsonResponse)
		return
	}
	LoggedinUserID, err := getIDbyUUID(uuid)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		jsonResponse, _ := json.Marshal(map[string]string{
			"message": "unauthorized",
		})
		w.Write(jsonResponse)
		return
	}

	// from member_email get the member_id
	memberID, err := getIDbyEmail(data.Email)

	// check if the user is the creator of the group
	// use the getGroup statement
	rows, err := statements["getGroup"].Query(data.GroupID)
	defer rows.Close()
	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(500)
		jsonResponse, _ := json.Marshal(map[string]string{
			"message": "internal server error",
		})
		w.Write(jsonResponse)
		return
	}
	group := Group{}

	for rows.Next() {
		err = rows.Scan(&group.ID, &group.Name, &group.Description, &group.CreatorId, &group.CreationDate, &group.Privacy)
		if err != nil {
			log.Println(err.Error())
			w.WriteHeader(500)
			jsonResponse, _ := json.Marshal(map[string]string{
				"message": "internal server error",
			})
			w.Write(jsonResponse)
			return
		}
	}

	if group.CreatorId != LoggedinUserID {
		w.WriteHeader(http.StatusUnauthorized)
		jsonResponse, _ := json.Marshal(map[string]string{
			"message": "unauthorized, you are not the creator of the group",
		})
		w.Write(jsonResponse)
		return
	}

	// remove it from the group_pending_members table
	_, err = statements["removeGroupPendingMember"].Exec(data.GroupID, memberID)
	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(500)
		jsonResponse, _ := json.Marshal(map[string]string{
			"message": "internal server error",
		})
		w.Write(jsonResponse)
		return
	}

	w.WriteHeader(http.StatusOK)
	jsonResponse, _ := json.Marshal(map[string]string{
		"message": "success: you rejected the group membership",
	})
	w.Write(jsonResponse)
	return
}

func approveFollowerHandler(w http.ResponseWriter, r *http.Request) {
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
		w.WriteHeader(400)
		jsonResponse, _ := json.Marshal(map[string]string{
			"message": "bad request",
		})
		w.Write(jsonResponse)
		return
	}

	fanID, err := getIDbyEmail(data.Email)

	LoggedinUserID, err := getIDbyUUID(uuid)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	// check if the user has the follower in the followers_pending table
	// use the getFollowersPending statement
	rows, err := statements["getFollowersPending"].Query(LoggedinUserID)
	defer rows.Close()
	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(500)
		jsonResponse, _ := json.Marshal(map[string]string{
			"message": "internal server error",
		})
		w.Write(jsonResponse)
		return
	}
	var followerID int
	var followersPending []int
	for rows.Next() {
		err = rows.Scan(&followerID)
		if err != nil {
			log.Println(err.Error())
			w.WriteHeader(500)
			jsonResponse, _ := json.Marshal(map[string]string{
				"message": "internal server error",
			})
			w.Write(jsonResponse)
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
		w.WriteHeader(http.StatusUnauthorized)
		jsonResponse, _ := json.Marshal(map[string]string{
			"message": "unauthorized",
		})
		w.Write(jsonResponse)
		return
	}
	// add the follower to the followers table
	_, err = statements["addFollower"].Exec(LoggedinUserID, fanID)
	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(500)
		jsonResponse, _ := json.Marshal(map[string]string{
			"message": "internal server error",
		})
		w.Write(jsonResponse)
		return
	}
	// remove it from the followers_pending table
	_, err = statements["removeFollowerPending"].Exec(LoggedinUserID, fanID)
	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(500)
		jsonResponse, _ := json.Marshal(map[string]string{
			"message": "internal server error. could not remove follower from followers_pending table",
		})
		w.Write(jsonResponse)
	}
	w.WriteHeader(http.StatusOK)
	jsonResponse, _ := json.Marshal(map[string]string{
		"message": "you accepted the follow request",
	})
	w.Write(jsonResponse)
	return
}

func rejectFollowerHandler(w http.ResponseWriter, r *http.Request) {
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
		w.WriteHeader(400)
		jsonResponse, _ := json.Marshal(map[string]string{
			"message": "bad request",
		})
		w.Write(jsonResponse)
		return
	}

	fanID, err := getIDbyEmail(data.Email)

	LoggedinUserID, err := getIDbyUUID(uuid)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// check if the user has the follower in the followers_pending table
	// use the getFollowersPending statement
	rows, err := statements["getFollowersPending"].Query(LoggedinUserID)
	defer rows.Close()
	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(500)
		jsonResponse, _ := json.Marshal(map[string]string{
			"message": "internal server error",
		})
		w.Write(jsonResponse)
		return
	}
	var followerID int
	var followersPending []int
	for rows.Next() {
		err = rows.Scan(&followerID)
		if err != nil {
			log.Println(err.Error())
			w.WriteHeader(500)
			jsonResponse, _ := json.Marshal(map[string]string{
				"message": "internal server error",
			})
			w.Write(jsonResponse)
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
		w.WriteHeader(http.StatusUnauthorized)
		jsonResponse, _ := json.Marshal(map[string]string{
			"message": "unauthorized. the target user was not found in the followers_pending table",
		})
		w.Write(jsonResponse)
		return
	}

	// remove it from the followers_pending table
	_, err = statements["removeFollowerPending"].Exec(LoggedinUserID, fanID)
	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(500)
		jsonResponse, _ := json.Marshal(map[string]string{
			"message": "internal server error. could not remove follower from followers_pending table",
		})
		w.Write(jsonResponse)
	}
	w.WriteHeader(http.StatusOK)
	jsonResponse, _ := json.Marshal(map[string]string{
		"message": "you rejected the follow request",
	})
	w.Write(jsonResponse)
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
	requestorId, err := getIDbyUUID(cookie.Value)
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
		GroupId int `json:"group_id"`
	}

	// get the group id from the request body
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err = decoder.Decode(&data)
	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(400)
		jsonResponse, _ := json.Marshal(map[string]string{
			"message": "bad request, failed to get group id from request body",
		})
		w.Write(jsonResponse)
		return
	}

	// add the person to the group_members table
	_, err = statements["addGroupMember"].Exec(data.GroupId, requestorId)
	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(500)
		jsonResponse, _ := json.Marshal(map[string]string{
			"message": "failed to add person to group_members table",
		})
		w.Write(jsonResponse)
		return
	}

	// remove the person from the group_invited_users table
	_, err = statements["removeGroupInvitedUser"].Exec(requestorId, data.GroupId)
	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(500)
		jsonResponse, _ := json.Marshal(map[string]string{
			"message": "failed to remove person from group_invited_users table",
		})
		w.Write(jsonResponse)
		return
	}

	w.WriteHeader(200)
	jsonResponse, _ := json.Marshal(map[string]string{
		"message": "success: you accepted the group invite",
	})
	w.Write(jsonResponse)
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
	requestorId, err := getIDbyUUID(cookie.Value)
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
		GroupId int `json:"group_id"`
	}

	// get the group id from the request body
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err = decoder.Decode(&data)
	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(400)
		jsonResponse, _ := json.Marshal(map[string]string{
			"message": "bad request, failed to get group id from request body",
		})
		w.Write(jsonResponse)
		return
	}

	// remove the person from the group_invited_users table
	_, err = statements["removeGroupInvitedUser"].Exec(requestorId, data.GroupId)
	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(500)
		jsonResponse, _ := json.Marshal(map[string]string{
			"message": "failed to remove person from group_invited_users table",
		})
		w.Write(jsonResponse)
		return
	}

	w.WriteHeader(200)
	jsonResponse, _ := json.Marshal(map[string]string{
		"message": "success: you rejected the group invite",
	})
	w.Write(jsonResponse)
	return
}
