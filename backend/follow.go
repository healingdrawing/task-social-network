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

// # FollowingHandler returns list of users that the target userID is following
//
// - If: the request has a JSON body, it will be decoded into a DTO, and the email of the user will be extracted
// and the user id will be extracted from the email
//
// - Else: the cookie of the current user will be extracted and the user id will be extracted from the uuid
func FollowingHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	defer func() {
		if err := recover(); err != nil {
			w.WriteHeader(500)
			jsonResponse, _ := json.Marshal(map[string]string{
				"message": "internal server error, something went wrong",
			})
			w.Write(jsonResponse)
		}
	}()
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
			w.WriteHeader(400)
			jsonResponse, _ := json.Marshal(map[string]string{
				"message": "bad request",
			})
			w.Write(jsonResponse)
			return
		}
		data.Email = incomingData["email"].(string)
		err = statements["getIDbyEmail"].QueryRow(data.Email).Scan(&userID)
		if err != nil {
			w.WriteHeader(500)
			jsonResponse, _ := json.Marshal(map[string]string{
				"message": "internal server error, getIDbyEmail query failed",
			})
			w.Write(jsonResponse)
			return
		}
	} else {
		cookie, err := r.Cookie("user_uuid")
		if err != nil {
			w.WriteHeader(400)
			jsonResponse, _ := json.Marshal(map[string]string{
				"message": "bad request, could not get cookie",
			})
			w.Write(jsonResponse)
			return
		}
		err = statements["getIDbyUUID"].QueryRow(cookie.Value).Scan(&userID)
		if err != nil {
			w.WriteHeader(500)
			jsonResponse, _ := json.Marshal(map[string]string{
				"message": "internal server error, getIDbyUUID query failed",
			})
			w.Write(jsonResponse)
			return
		}
	}
	rows, err := statements["getFollowing"].Query(userID)
	if err != nil {
		w.WriteHeader(500)
		jsonResponse, _ := json.Marshal(map[string]string{
			"message": "internal server error, getFollowing query failed",
		})
		w.Write(jsonResponse)
		return
	}
	allFollowingIds := []int{}
	for rows.Next() {
		var id int
		err := rows.Scan(&id)
		if err != nil {
			w.WriteHeader(500)
			jsonResponse, _ := json.Marshal(map[string]string{
				"message": "internal server error, scan failed",
			})
			w.Write(jsonResponse)
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
			w.WriteHeader(500)
			jsonResponse, _ := json.Marshal(map[string]string{
				"message": "internal server error, getUserbyID query failed",
			})
			w.Write(jsonResponse)
			return
		}
		followingProfile.FullName = followingProfile.firstName + " " + followingProfile.lastName
		followings = append(followings, followingProfile)
	}
	jsonResponse, _ := json.Marshal(followings)
	// write the response
	w.Write(jsonResponse)
}

func FollowersHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	defer func() {
		if err := recover(); err != nil {
			w.WriteHeader(500)
			jsonResponse, _ := json.Marshal(map[string]string{
				"message": "internal server error. Something went wrong",
			})
			w.Write(jsonResponse)
		}
	}()
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
			w.WriteHeader(400)
			jsonResponse, _ := json.Marshal(map[string]string{
				"message": "bad request",
			})
			w.Write(jsonResponse)
			return
		}
		data.Email = incomingData["email"].(string)
		// get the user id from the email

		err = statements["getIDbyEmail"].QueryRow(data.Email).Scan(&userID)
		if err != nil {
			w.WriteHeader(500)
			jsonResponse, _ := json.Marshal(map[string]string{
				"message": "internal server error. getIDbyEmail query failed",
			})
			w.Write(jsonResponse)
			return
		}
	} else {
		// get uuid of current user from cookies
		cookie, err := r.Cookie("user_uuid")
		if err != nil {
			w.WriteHeader(400)
			jsonResponse, _ := json.Marshal(map[string]string{
				"message": "bad request",
			})
			w.Write(jsonResponse)
			return
		}
		// get the user id from the uuid
		err = statements["getIDbyUUID"].QueryRow(cookie.Value).Scan(&userID)
		if err != nil {
			w.WriteHeader(500)
			jsonResponse, _ := json.Marshal(map[string]string{
				"message": "internal server error, getIDbyUUID query failed",
			})
			w.Write(jsonResponse)
			return
		}
	}

	// get the list of users that follow the target userID
	rows, err := statements["getFollowers"].Query(userID)
	if err != nil {
		w.WriteHeader(500)
		jsonResponse, _ := json.Marshal(map[string]string{
			"message": "internal server error, getFollowers query failed",
		})
		w.Write(jsonResponse)
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
			w.WriteHeader(500)
			jsonResponse, _ := json.Marshal(map[string]string{
				"message": "internal server error, scan failed",
			})
			w.Write(jsonResponse)
			return
		}
		// get the email, first_name, last_name of the follower
		var nickname string
		err = statements["getUserbyID"].QueryRow(follower.id).Scan(&follower.Email, &follower.firstName, &follower.lastName, &nickname)
		if err != nil {
			w.WriteHeader(500)
			jsonResponse, _ := json.Marshal(map[string]string{
				"message": "internal server error, getUserbyID query failed",
			})
			w.Write(jsonResponse)
			return
		}
		// set the FullName of the Follower
		follower.FullName = follower.firstName + " " + follower.lastName
		// append the Follower to the slice
		followers = append(followers, follower)
	}

	// encode the slice into json
	jsonResponse, _ := json.Marshal(followers)
	// write the response
	w.Write(jsonResponse)

}

// # FollowHandler adds the current user to the followers list of the target user
//
// - If: the target user is private, the current user will be added to the followers_pending table
//
// - Else: the current user will be added to the followers table
func FollowHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	defer func() {
		if err := recover(); err != nil {
			w.WriteHeader(500)
			jsonResponse, _ := json.Marshal(map[string]string{
				"message": "internal server error. Something went wrong",
			})
			w.Write(jsonResponse)
		}
	}()
	// get the uuid of the current user from the cookies
	cookie, err := r.Cookie("user_uuid")
	if err != nil {
		w.WriteHeader(400)
		jsonResponse, _ := json.Marshal(map[string]string{
			"message": "bad request",
		})
		w.Write(jsonResponse)
		return
	}
	// get the user id from the uuid
	var userID int
	err = statements["getIDbyUUID"].QueryRow(cookie.Value).Scan(&userID)
	if err != nil {
		w.WriteHeader(500)
		jsonResponse, _ := json.Marshal(map[string]string{
			"message": "internal server error. getIDbyUUID query failed",
		})
		w.Write(jsonResponse)
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
		w.WriteHeader(400)
		jsonResponse, _ := json.Marshal(map[string]string{
			"message": "bad request",
		})
		w.Write(jsonResponse)
		return
	}
	// get the user id from the email
	var targetUserID int
	err = statements["getIDbyEmail"].QueryRow(data.Email).Scan(&targetUserID)
	if err != nil {
		w.WriteHeader(500)
		jsonResponse, _ := json.Marshal(map[string]string{
			"message": "internal server error. getIDbyEmail query failed",
		})
		w.Write(jsonResponse)
		return
	}

	// check if our target whom we want to follow is private or public
	var privacy string
	err = statements["getUserPrivacy"].QueryRow(targetUserID).Scan(&privacy)
	if err != nil {
		w.WriteHeader(500)
		jsonResponse, _ := json.Marshal(map[string]string{
			"message": "internal server error. getUserPrivacy query failed",
		})
		w.Write(jsonResponse)
		return
	}

	if privacy == "private" {
		// if the target is private, add the follower to the followers_pending table
		_, err = statements["addFollowerPending"].Exec(targetUserID, userID)
		if err != nil {
			w.WriteHeader(500)
			jsonResponse, _ := json.Marshal(map[string]string{
				"message": "internal server error. could not add you to the pending followers list",
			})
			w.Write(jsonResponse)
			return
		}
		w.WriteHeader(200)
		jsonResponse, _ := json.Marshal(map[string]string{
			"message": "request sent to follow the user",
		})
		w.Write(jsonResponse)
		return
	} else {
		// if the target is public, add the follower to the followers table
		_, err = statements["addFollower"].Exec(targetUserID, userID)
		if err != nil {
			w.WriteHeader(500)
			jsonResponse, _ := json.Marshal(map[string]string{
				"message": "internal server error. could not add you to the followers list",
			})
			w.Write(jsonResponse)
			return
		}
		w.WriteHeader(200)
		jsonResponse, _ := json.Marshal(map[string]string{
			"message": "user followed",
		})
		w.Write(jsonResponse)
		return
	}
}

// # UnfollowHandler removes the current user from the followers list of the target user
//
// - If the followers table has a row with the current loggedin user as follower_id and the target user as user_id,
// the row will be deleted
func UnfollowHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	defer func() {
		if err := recover(); err != nil {
			w.WriteHeader(500)
			jsonResponse, _ := json.Marshal(map[string]string{
				"message": "internal server error",
			})
			w.Write(jsonResponse)
		}
	}()
	// get the uuid of the current user from the cookies
	cookie, err := r.Cookie("user_uuid")
	if err != nil {
		w.WriteHeader(400)
		jsonResponse, _ := json.Marshal(map[string]string{
			"message": "bad request",
		})
		w.Write(jsonResponse)
		return
	}
	// get the user id from the uuid
	var userID int
	err = statements["getIDbyUUID"].QueryRow(cookie.Value).Scan(&userID)
	if err != nil {
		w.WriteHeader(500)
		jsonResponse, _ := json.Marshal(map[string]string{
			"message": "internal server error",
		})
		w.Write(jsonResponse)
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
		w.WriteHeader(400)
		jsonResponse, _ := json.Marshal(map[string]string{
			"message": "bad request",
		})
		w.Write(jsonResponse)
		return
	}
	// get the user id from the email
	var targetUserID int
	err = statements["getIDbyEmail"].QueryRow(data.Email).Scan(&targetUserID)
	if err != nil {
		w.WriteHeader(500)
		jsonResponse, _ := json.Marshal(map[string]string{
			"message": "internal server error. could not get target user id",
		})
		w.Write(jsonResponse)
		return
	}
	_, err = statements["removeFollower"].Exec(targetUserID, userID)
	if err != nil {
		w.WriteHeader(500)
		jsonResponse, _ := json.Marshal(map[string]string{
			"message": "internal server error. could not remove you from the followers list",
		})
		w.Write(jsonResponse)
		return
	}
	w.WriteHeader(200)
	jsonResponse, _ := json.Marshal(map[string]string{
		"message": "user unfollowed",
	})
	w.Write(jsonResponse)
	return
}
