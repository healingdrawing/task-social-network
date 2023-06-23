package main

import (
	"encoding/json"
	"net/http"
)

type Following struct {
	id        int
	Email     string `json:"email"`
	firstName string
	lastName  string
	FullName  string `json:"full_name"`
}

type Follower struct {
	id        int
	Email     string `json:"email"`
	firstName string
	lastName  string
	FullName  string `json:"full_name"`
}

// # followingHandler returns list of users that the target userID is following
//
// - If: the request has a JSON body, it will be decoded into a DTO, and the email of the user will be extracted
// and the user id will be extracted from the email
//
// - Else: the cookie of the current user will be extracted and the user id will be extracted from the uuid
func followingHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	defer recovery(w)
	var userID int

	// if the incoming request has a JSON body, decode it into the DTO, else get the cookie of current user
	// and get the list of users that the current user is following
	if r.Header.Get("Content-Type") == "application/json" {
		data := struct {
			Email string `json:"email"`
		}{}
		incomingData := map[string]any{}
		// decode the request body into the DTO
		decoder := json.NewDecoder(r.Body)
		decoder.DisallowUnknownFields()
		err := decoder.Decode(&incomingData)
		if err != nil {
			jsonResponse(w, http.StatusBadRequest, "")
			return
		}
		data.Email = incomingData["email"].(string)
		err = statements["getIDbyEmail"].QueryRow(data.Email).Scan(&userID)
		if err != nil {
			jsonResponse(w, http.StatusInternalServerError, " getIDbyEmail query failed")
			return
		}
	} else {
		cookie, err := r.Cookie("user_uuid")
		if err != nil {
			jsonResponse(w, http.StatusUnauthorized, "could not get cookie")
			return
		}
		err = statements["getIDbyUUID"].QueryRow(cookie.Value).Scan(&userID)
		if err != nil {
			jsonResponse(w, http.StatusInternalServerError, " getIDbyUUID query failed")
			return
		}
	}
	rows, err := statements["getFollowing"].Query(userID)
	if err != nil {
		jsonResponse(w, http.StatusInternalServerError, " getFollowing query failed")
		return
	}
	allFollowingIds := []int{}
	for rows.Next() {
		var id int
		err := rows.Scan(&id)
		if err != nil {
			jsonResponse(w, http.StatusInternalServerError, " scan allFollowingIds for loop failed")
			return
		}
		allFollowingIds = append(allFollowingIds, id)
	}

	var followings []Following

	for _, id := range allFollowingIds {
		var nickname string
		var followingProfile Following
		err = statements["getUserbyID"].QueryRow(id).Scan(&followingProfile.Email, &followingProfile.firstName, &followingProfile.lastName, &nickname)
		if err != nil {
			jsonResponse(w, http.StatusInternalServerError, " getUserbyID query failed")
			return
		}
		followingProfile.FullName = followingProfile.firstName + " " + followingProfile.lastName
		followings = append(followings, followingProfile)
	}
	jsonResponseObj, _ := json.Marshal(followings)
	// write the response
	_, err = w.Write(jsonResponseObj)
	if err != nil {
		jsonResponse(w, http.StatusInternalServerError, " w.Write(jsonResponseObj)<-followings failed")
		return
	}
}

func followersHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	defer recovery(w)
	var userID int

	// if the incoming request has a JSON body, decode it into the DTO, else get the cookie of current user
	// and get the list of users that the current user is following
	if r.Header.Get("Content-Type") == "application/json" {
		data := struct {
			Email string `json:"email"`
		}{}
		incomingData := map[string]any{}
		// decode the request body into the DTO
		decoder := json.NewDecoder(r.Body)
		decoder.DisallowUnknownFields()
		err := decoder.Decode(&incomingData)
		if err != nil {
			jsonResponse(w, http.StatusBadRequest, "")
			return
		}
		data.Email = incomingData["email"].(string)
		// get the user id from the email

		err = statements["getIDbyEmail"].QueryRow(data.Email).Scan(&userID)
		if err != nil {
			jsonResponse(w, http.StatusInternalServerError, " getIDbyEmail query failed")
			return
		}
	} else {
		// get uuid of current user from cookies
		cookie, err := r.Cookie("user_uuid")
		if err != nil {
			jsonResponse(w, http.StatusUnauthorized, "could not get cookie")
			return
		}
		// get the user id from the uuid
		err = statements["getIDbyUUID"].QueryRow(cookie.Value).Scan(&userID)
		if err != nil {
			jsonResponse(w, http.StatusInternalServerError, " getIDbyUUID query failed")
			return
		}
	}

	// get the list of users that follow the target userID
	rows, err := statements["getFollowers"].Query(userID)
	if err != nil {
		jsonResponse(w, http.StatusInternalServerError, " getFollowers query failed")
		return
	}
	// create a slice of Follower
	var followers []Following
	// iterate over the rows
	for rows.Next() {
		// create a new Follower
		var follower Following
		// scan the row into the Follower
		err := rows.Scan(&follower.id)
		if err != nil {
			jsonResponse(w, http.StatusInternalServerError, " scan getFollowers &follower.id for loop failed")
			return
		}
		// get the email, first_name, last_name of the follower
		var nickname string
		err = statements["getUserbyID"].QueryRow(follower.id).Scan(&follower.Email, &follower.firstName, &follower.lastName, &nickname)
		if err != nil {
			jsonResponse(w, http.StatusInternalServerError, " getUserbyID query failed")
			return
		}
		follower.FullName = follower.firstName + " " + follower.lastName
		followers = append(followers, follower)
	}

	jsonResponseObj, err := json.Marshal(followers)
	if err != nil {
		jsonResponse(w, http.StatusInternalServerError, " json.Marshal(followers) failed")
		return
	}
	_, err = w.Write(jsonResponseObj)
	if err != nil {
		jsonResponse(w, http.StatusInternalServerError, " w.Write(jsonResponseObj)<-followers failed")
		return
	}

}

// # followHandler adds the current user to the followers list of the target user
//
// - If: the target user is private, the current user will be added to the followers_pending table
//
// - Else: the current user will be added to the followers table
func followHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	defer recovery(w)
	// get the uuid of the current user from the cookies
	cookie, err := r.Cookie("user_uuid")
	if err != nil {
		jsonResponse(w, http.StatusUnauthorized, "could not get cookie")
		return
	}
	// get the user id from the uuid
	var userID int
	err = statements["getIDbyUUID"].QueryRow(cookie.Value).Scan(&userID)
	if err != nil {
		jsonResponse(w, http.StatusInternalServerError, " getIDbyUUID query failed")
		return
	}
	// incoming DTO UserFollowerRequest
	data := struct {
		Email string `json:"email"`
	}{}
	// decode the request body into the DTO
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err = decoder.Decode(&data)
	if err != nil {
		jsonResponse(w, http.StatusBadRequest, "")
		return
	}
	// get the user id from the email
	var targetUserID int
	err = statements["getIDbyEmail"].QueryRow(data.Email).Scan(&targetUserID)
	if err != nil {
		jsonResponse(w, http.StatusInternalServerError, " getIDbyEmail query failed")
		return
	}

	// check if our target whom we want to follow is private or public
	var privacy string
	err = statements["getUserPrivacy"].QueryRow(targetUserID).Scan(&privacy)
	if err != nil {
		jsonResponse(w, http.StatusInternalServerError, " getUserPrivacy query failed")
		return
	}

	if privacy == "private" {
		// if the target is private, add the follower to the followers_pending table
		_, err = statements["addFollowerPending"].Exec(targetUserID, userID)
		if err != nil {
			jsonResponse(w, http.StatusInternalServerError, " addFollowerPending query failed")
			return
		}
		jsonResponse(w, http.StatusOK, "request sent to follow the user")
		return
	} else {
		// if the target is public, add the follower to the followers table
		_, err = statements["addFollower"].Exec(targetUserID, userID)
		if err != nil {
			jsonResponse(w, http.StatusInternalServerError, " addFollower query failed")
			return
		}
		jsonResponse(w, http.StatusOK, "user followed")
		return
	}
}

// # unfollowHandler removes the current user from the followers list of the target user
//
// - If the followers table has a row with the current loggedin user as follower_id and the target user as user_id,
// the row will be deleted
func unfollowHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	defer recovery(w)
	// get the uuid of the current user from the cookies
	cookie, err := r.Cookie("user_uuid")
	if err != nil {
		jsonResponse(w, http.StatusUnauthorized, "could not get cookie")
		return
	}
	// get the user id from the uuid
	var userID int
	err = statements["getIDbyUUID"].QueryRow(cookie.Value).Scan(&userID)
	if err != nil {
		jsonResponse(w, http.StatusInternalServerError, " getIDbyUUID query failed")
		return
	}
	// incoming DTO UserFollowerRequest
	data := struct {
		Email string `json:"email"`
	}{}
	// decode the request body into the DTO
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err = decoder.Decode(&data)
	if err != nil {
		jsonResponse(w, http.StatusBadRequest, "")
		return
	}
	// get the user id from the email
	var targetUserID int
	err = statements["getIDbyEmail"].QueryRow(data.Email).Scan(&targetUserID)
	if err != nil {
		jsonResponse(w, http.StatusInternalServerError, " getIDbyEmail->targetUserID query failed")
		return
	}
	_, err = statements["removeFollower"].Exec(targetUserID, userID)
	if err != nil {
		jsonResponse(w, http.StatusInternalServerError, " removeFollower query failed")
		return
	}
	jsonResponse(w, http.StatusOK, "user unfollowed")
	return
}

// # FollowRequestListHandler returns a list of users that have sent a follow request to the current user
//
// - No JSON body is required
func followRequestListHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	defer recovery(w)
	// get the uuid of the current user from the cookies
	cookie, err := r.Cookie("user_uuid")
	if err != nil {
		jsonResponse(w, http.StatusUnauthorized, "could not get cookie")
		return
	}
	// get the user id from the uuid
	var userID int
	err = statements["getIDbyUUID"].QueryRow(cookie.Value).Scan(&userID)
	if err != nil {
		jsonResponse(w, http.StatusInternalServerError, " getIDbyUUID query failed")
		return
	}
	rows, err := statements["getFollowersPending"].Query(userID)
	if err != nil {
		jsonResponse(w, http.StatusInternalServerError, " getFollowersPending query failed")
		return
	}
	var followers []Follower
	for rows.Next() {
		var follower Follower
		err = rows.Scan(&follower.id)
		if err != nil {
			jsonResponse(w, http.StatusInternalServerError, " &follower.id scan for loop failed")
			return
		}
		var nickname string
		err = statements["getUserbyID"].QueryRow(follower.id).Scan(&follower.Email, &follower.firstName, &follower.lastName, &nickname)
		if err != nil {
			jsonResponse(w, http.StatusInternalServerError, " getUserbyID query failed")
			return
		}
		follower.FullName = follower.firstName + " " + follower.lastName
		followers = append(followers, follower)
	}
	jsonResponseObj, _ := json.Marshal(followers)
	_, err = w.Write(jsonResponseObj)
	if err != nil {
		jsonResponse(w, http.StatusInternalServerError, "w.Write(jsonResponseObj)<-followers failed")
		return
	}
}
