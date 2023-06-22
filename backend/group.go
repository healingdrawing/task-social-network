package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"
)

type Group struct {
	ID           int       `json:"ID"`
	Name         string    `json:"name"`
	Description  string    `json:"description"`
	Creator      int       `json:"creator"`
	CreationDate time.Time `json:"creation_date"`
	Privacy      string    `json:"privacy"`
}

type GroupMember struct {
	GroupID  int `json:"group_id"`
	MemberID int `json:"member_id"`
}

type GroupCreationDTO struct {
	UUID        string `json:"UUID"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Privacy     string `json:"privacy"`
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
	data.Creator, err = getIDbyUUID(incomingData.UUID)
	if data.Name == "" || data.Description == "" || data.Privacy == "" {
		w.WriteHeader(400)
		jsonResponse, _ := json.Marshal(map[string]string{
			"message": "bad request, missing fields",
		})
		w.Write(jsonResponse)
		return
	}
	result, err := statements["addGroup"].Exec(data.Name, data.Description, data.Creator, time.Now().Format("2006-01-02 15:04:05"), data.Privacy)
	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(500)
		jsonResponse, _ := json.Marshal(map[string]string{
			"message": "internal server error",
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
	_, err = statements["addGroupMember"].Exec(groupID, data.Creator)
	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(500)
		jsonResponse, _ := json.Marshal(map[string]string{
			"message": "internal server error, failed to add group creator as group member",
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
	var data GroupMember
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
	// check ig group is private
	var group Group
	err = statements["getGroup"].QueryRow(data.GroupID).Scan(&group.ID, &group.Name, &group.Description, &group.Creator, &group.CreationDate, &group.Privacy)
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
		_, err = statements["addGroupPendingMember"].Exec(data.GroupID, data.MemberID)
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
	w.WriteHeader(200)
	jsonResponse, _ := json.Marshal(map[string]string{
		"message": "group joined",
	})
	w.Write(jsonResponse)
}

// Incoming JSON DTO for group creation over handler groupCreateHandler
// {
//	"UUID": 1,
// 	"name": "group name",
// 	"description": "group description",
// 	"privacy": "public"
// }

// Incoming JSON DTO for group joining over handler groupJoinHandler
// {
// 	"group_id": 1,
// 	"member_id": 1
// }
