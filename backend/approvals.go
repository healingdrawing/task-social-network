package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type GroupMemberRequest struct {
	UUID     string `json:"UUID"`
	GroupID  int    `json:"group_id"`
	MemberID int    `json:"member_id"`
}

// place to receive requests and approve requests
// to join groups
// to follow users

func approveGroupMembershipHandler(w http.ResponseWriter, r *http.Request) {
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
	// incoming DTO GroupMember
	data := GroupMemberRequest{}
	// decode the request body into the DTO
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
	LoggedinUserID, err := getIDbyUUID(data.UUID)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		jsonResponse, _ := json.Marshal(map[string]string{
			"message": "unauthorized",
		})
		w.Write(jsonResponse)
		return
	}
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
		err = rows.Scan(&group.ID, &group.Name, &group.Description, &group.Creator, &group.CreationDate, &group.Privacy)
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

	if group.Creator != LoggedinUserID {
		w.WriteHeader(http.StatusUnauthorized)
		jsonResponse, _ := json.Marshal(map[string]string{
			"message": "unauthorized",
		})
		w.Write(jsonResponse)
		return
	}

	_, err = statements["addGroupMember"].Exec(data.GroupID, data.MemberID)
	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(500)
		jsonResponse, _ := json.Marshal(map[string]string{
			"message": "internal server error",
		})
		w.Write(jsonResponse)
		return
	}

	// remove it from the group_pending_members table
	_, err = statements["removeGroupPendingMember"].Exec(data.GroupID, data.MemberID)
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
		"message": "success",
	})
	w.Write(jsonResponse)
	return
}

func rejectGroupMembershipHandler(w http.ResponseWriter, r *http.Request) {
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
	// incoming DTO GroupMember
	data := GroupMemberRequest{}
	// decode the request body into the DTO
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
	LoggedinUserID, err := getIDbyUUID(data.UUID)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		jsonResponse, _ := json.Marshal(map[string]string{
			"message": "unauthorized",
		})
		w.Write(jsonResponse)
		return
	}
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
		err = rows.Scan(&group.ID, &group.Name, &group.Description, &group.Creator, &group.CreationDate, &group.Privacy)
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

	if group.Creator != LoggedinUserID {
		w.WriteHeader(http.StatusUnauthorized)
		jsonResponse, _ := json.Marshal(map[string]string{
			"message": "unauthorized",
		})
		w.Write(jsonResponse)
		return
	}

	// remove it from the group_pending_members table
	_, err = statements["removeGroupPendingMember"].Exec(data.GroupID, data.MemberID)
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
