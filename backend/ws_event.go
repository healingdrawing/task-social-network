package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

type WS_GROUP_EVENT_SUBMIT_DTO struct {
	User_uuid   string `json:"user_uuid"`
	Group_id    int    `json:"group_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Date        string `json:"date"`
	Decision    string `json:"decision"` // "going", "not going"
}

// wsGroupEventSubmitHandler create a new event
//
// @rparam {group_id int, name string, description string, date string}
func wsGroupEventSubmitHandler(conn *websocket.Conn, messageData map[string]interface{}) {
	defer wsRecover()

	uuid, ok := messageData["user_uuid"].(string)
	if !ok {
		log.Println("failed to get user_uuid from messageData")
		wsSendError(WS_ERROR_RESPONSE_DTO{fmt.Sprint(http.StatusUnprocessableEntity) + " failed to get user_uuid from messageData"})
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
		log.Println("failed to get group_id from messageData")
		wsSendError(WS_ERROR_RESPONSE_DTO{fmt.Sprint(http.StatusUnprocessableEntity) + " failed to get group_id from messageData"})
		return
	}
	group_id := int(_group_id)

	is_member, err := isGroupMember(user_id, group_id)
	if err != nil {
		log.Println("failed to check if user is member of the group", err.Error())
		wsSendError(WS_ERROR_RESPONSE_DTO{fmt.Sprint(http.StatusUnprocessableEntity) + " failed to check if user is member of the group"})
		return
	}
	if !is_member {
		log.Println("user is not member of the group")
		wsSendError(WS_ERROR_RESPONSE_DTO{fmt.Sprint(http.StatusUnprocessableEntity) + " user is not member of the group"})
		return
	}

	var data WS_GROUP_EVENT_SUBMIT_DTO
	fields := map[string]*string{
		"title":       &data.Title,
		"description": &data.Description,
		"date":        &data.Date,
		"decision":    &data.Decision,
	}

	for key, ptr := range fields {
		value, ok := messageData[key].(string)
		if !ok {
			log.Printf("failed to get %s from messageData\n", key)
			wsSendError(WS_ERROR_RESPONSE_DTO{fmt.Sprintf("%d failed to get %s from messageData", http.StatusUnprocessableEntity, key)})
			return
		}
		*ptr = value
	}

	data.Title = strings.TrimSpace(data.Title)
	data.Description = strings.TrimSpace(data.Description)
	data.Date = strings.TrimSpace(data.Date)
	data.Decision = strings.TrimSpace(data.Decision)

	if data.Decision != "going" && data.Decision != "not going" {
		log.Println("invalid decision")
		wsSendError(WS_ERROR_RESPONSE_DTO{fmt.Sprint(http.StatusUnprocessableEntity) + " invalid decision"})
		return
	}

	if data.Title == "" || data.Description == "" || data.Date == "" {
		log.Println("empty fields")
		wsSendError(WS_ERROR_RESPONSE_DTO{fmt.Sprint(http.StatusUnprocessableEntity) + " empty fields"})
		return
	}
	created_at := time.Now().Format("2006-01-02 15:04:05")

	// add the event to the table
	result, err := statements["addEvent"].Exec(group_id, data.Title, data.Description, data.Date, created_at)
	if err != nil {
		log.Println("addEvent query failed", err.Error())
		wsSendError(WS_ERROR_RESPONSE_DTO{fmt.Sprint(http.StatusInternalServerError) + " addEvent query failed"})
		return
	}

	event_id, err := result.LastInsertId()
	if err != nil {
		log.Println("failed to get last insert id", err.Error())
		wsSendError(WS_ERROR_RESPONSE_DTO{fmt.Sprint(http.StatusInternalServerError) + " failed to get last insert id"})
		return
	}

	_, err = statements["addEventParticipant"].Exec(event_id, user_id, data.Decision)
	if err != nil {
		log.Println("addEventParticipant query failed", err.Error())
		wsSendError(WS_ERROR_RESPONSE_DTO{fmt.Sprint(http.StatusInternalServerError) + " addEventParticipant query failed"})
		return
	}

	wsGroupEventsListHandler(conn, messageData)
}

type WS_GROUP_EVENT_RESPONSE_DTO struct {
	Id          int    `json:"id"`
	Group_id    int    `json:"group_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Date        string `json:"date"`
	Created_at  string `json:"created_at"`
	Decision    string `json:"decision"`
}

type WS_GROUP_EVENTS_LIST_DTO []WS_GROUP_EVENT_RESPONSE_DTO

// wsGroupEventsListHandler returns all events for a group, with the user's decision
//
// @rparam {group_id int}
func wsGroupEventsListHandler(conn *websocket.Conn, messageData map[string]interface{}) {
	defer wsRecover()

	uuid, ok := messageData["user_uuid"].(string)
	if !ok {
		log.Println("failed to get user_uuid from messageData")
		wsSendError(WS_ERROR_RESPONSE_DTO{fmt.Sprint(http.StatusUnprocessableEntity) + " failed to get user_uuid from messageData"})
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
		log.Println("failed to get group_id from messageData\n", messageData)
		wsSendError(WS_ERROR_RESPONSE_DTO{fmt.Sprint(http.StatusUnprocessableEntity) + " failed to get group_id from messageData"})
		return
	}
	group_id := int(_group_id)

	rows, err := statements["getEvents"].Query(group_id)
	if err != nil {
		log.Println("getEvents query failed", err.Error())
		wsSendError(WS_ERROR_RESPONSE_DTO{fmt.Sprint(http.StatusInternalServerError) + " getEvents query failed"})
		return
	}
	var events_list WS_GROUP_EVENTS_LIST_DTO
	for rows.Next() {
		var event WS_GROUP_EVENT_RESPONSE_DTO
		err = rows.Scan(
			&event.Id,
			&event.Group_id,
			&event.Title,
			&event.Description,
			&event.Date,
			&event.Created_at,
		)
		if err != nil {
			log.Println("getEvents query failed to scan", err.Error())
			wsSendError(WS_ERROR_RESPONSE_DTO{fmt.Sprint(http.StatusInternalServerError) + " getEvents query failed to scan"})
			return
		}
		events_list = append(events_list, event)
	}

	// get the user's decision for each event
	for i, event := range events_list {
		rows, err := statements["getEventParticipantDecision"].Query(event.Id, user_id)
		if err != nil {
			log.Println("getEventParticipantDecision query failed", err.Error())
			wsSendError(WS_ERROR_RESPONSE_DTO{fmt.Sprint(http.StatusInternalServerError) + " getEventParticipantDecision query failed"})
			return
		}
		defer rows.Close()
		if rows.Next() {
			var decision string
			err = rows.Scan(&decision)
			if err != nil {
				log.Println("getEventParticipantDecision query failed to scan", err.Error())
				wsSendError(WS_ERROR_RESPONSE_DTO{fmt.Sprint(http.StatusInternalServerError) + " getEventParticipantDecision query failed to scan"})
				return
			}
			events_list[i].Decision = decision
		} else {
			events_list[i].Decision = "waiting" // this case raise two buttons
		}
		rows.Close() // yes, it is ugly solution inside loop, but not a time for adventures now
	}

	wsSendGroupEventsList(events_list)
}

type WS_USER_GROUPS_FRESH_EVENT_RESPOSE_DTO struct {
	Event_id          int    `json:"event_id"`
	Event_title       string `json:"event_title"`
	Event_description string `json:"event_description"`
	Event_date        string `json:"event_date"`

	Group_id          int    `json:"group_id"`
	Group_name        string `json:"group_name"`
	Group_description string `json:"group_description"`
}

type WS_USER_GROUPS_FRESH_EVENTS_LIST_DTO []WS_USER_GROUPS_FRESH_EVENT_RESPOSE_DTO

// wsUserGroupsFreshEventsListHandler returns all fresh(user decision needed) events for a user
func wsUserGroupsFreshEventsListHandler(conn *websocket.Conn, messageData map[string]interface{}) {
	defer wsRecover()

	uuid, ok := messageData["user_uuid"].(string)
	if !ok {
		log.Println("failed to get user_uuid from messageData")
		wsSendError(WS_ERROR_RESPONSE_DTO{fmt.Sprint(http.StatusUnprocessableEntity) + " failed to get user_uuid from messageData"})
		return
	}

	user_id, err := get_user_id_by_uuid(uuid)
	if err != nil {
		log.Println("failed to get ID of the message sender", err.Error())
		wsSendError(WS_ERROR_RESPONSE_DTO{fmt.Sprint(http.StatusUnprocessableEntity) + " failed to get ID of the message sender"})
		return
	}

	rows, err := statements["getFreshEvents"].Query(user_id, user_id)
	if err != nil {
		log.Println("getFreshEvents query failed", err.Error())
		wsSendError(WS_ERROR_RESPONSE_DTO{fmt.Sprint(http.StatusInternalServerError) + " getFreshEvents query failed"})
		return
	}
	var fresh_events_list WS_USER_GROUPS_FRESH_EVENTS_LIST_DTO
	for rows.Next() {
		var fresh_event WS_USER_GROUPS_FRESH_EVENT_RESPOSE_DTO
		err = rows.Scan(
			&fresh_event.Event_id,
			&fresh_event.Event_title,
			&fresh_event.Event_description,
			&fresh_event.Event_date,
			&fresh_event.Group_id,
			&fresh_event.Group_name,
			&fresh_event.Group_description,
		)
		if err != nil {
			log.Println("getFreshEvents query failed to scan", err.Error())
			wsSendError(WS_ERROR_RESPONSE_DTO{fmt.Sprint(http.StatusInternalServerError) + " getFreshEvents query failed to scan"})
			return
		}
		fresh_events_list = append(fresh_events_list, fresh_event)
	}

	wsSendUserGroupsFreshEventsList(fresh_events_list)
}

// wsEventGoingHandler marks a user as going to an event
//
// @rparam {event_id int}
func wsGroupEventGoingHandler(conn *websocket.Conn, messageData map[string]interface{}) {
	defer wsRecover()

	uuid, ok := messageData["user_uuid"].(string)
	if !ok {
		log.Println("failed to get user_uuid from messageData")
		wsSendError(WS_ERROR_RESPONSE_DTO{fmt.Sprint(http.StatusUnprocessableEntity) + " failed to get user_uuid from messageData"})
		return
	}

	user_id, err := get_user_id_by_uuid(uuid)
	if err != nil {
		log.Println("failed to get ID of the message sender", err.Error())
		wsSendError(WS_ERROR_RESPONSE_DTO{fmt.Sprint(http.StatusUnprocessableEntity) + " failed to get ID of the message sender"})
		return
	}

	_event_id, ok := messageData["event_id"].(float64)
	if !ok {
		log.Println("failed to get event_id from messageData")
		wsSendError(WS_ERROR_RESPONSE_DTO{fmt.Sprint(http.StatusUnprocessableEntity) + " failed to get event_id from messageData"})
		return
	}
	event_id := int(_event_id)

	_, err = statements["addEventParticipant"].Exec(event_id, user_id, "going")
	if err != nil {
		log.Println("addEventParticipant query failed", err.Error())
		wsSendError(WS_ERROR_RESPONSE_DTO{fmt.Sprint(http.StatusInternalServerError) + " addEventParticipant query failed"})
		return
	}

	wsGroupEventsListHandler(conn, messageData)

}

// wsGroupEventNotGoingHandler marks a user as not going to an event
//
// @rparam {event_id int}
func wsGroupEventNotGoingHandler(conn *websocket.Conn, messageData map[string]interface{}) {
	defer wsRecover()

	uuid, ok := messageData["user_uuid"].(string)
	if !ok {
		log.Println("failed to get user_uuid from messageData")
		wsSendError(WS_ERROR_RESPONSE_DTO{fmt.Sprint(http.StatusUnprocessableEntity) + " failed to get user_uuid from messageData"})
		return
	}

	user_id, err := get_user_id_by_uuid(uuid)
	if err != nil {
		log.Println("failed to get ID of the message sender", err.Error())
		wsSendError(WS_ERROR_RESPONSE_DTO{fmt.Sprint(http.StatusUnprocessableEntity) + " failed to get ID of the message sender"})
		return
	}

	_event_id, ok := messageData["event_id"].(float64)
	if !ok {
		log.Println("failed to get event_id from messageData")
		wsSendError(WS_ERROR_RESPONSE_DTO{fmt.Sprint(http.StatusUnprocessableEntity) + " failed to get event_id from messageData"})
		return
	}
	event_id := int(_event_id)

	_, err = statements["addEventParticipant"].Exec(event_id, user_id, "not going")
	if err != nil {
		log.Println("addEventParticipant query failed", err.Error())
		wsSendError(WS_ERROR_RESPONSE_DTO{fmt.Sprint(http.StatusInternalServerError) + " addEventParticipant query failed"})
		return
	}

	wsGroupEventsListHandler(conn, messageData)

}
