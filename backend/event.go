package main

import (
	"encoding/base64"
	"encoding/json"
	"log"
	"net/http"
	"time"
)

type EventDTOin struct {
	GroupID     int    `json:"group_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Date        string `json:"date"`
	Location    string `json:"location"`
	CreatorId   int    `json:"creator_id"`
	CreatedAt   string `json:"created_at"`
}

// # eventNewHandler create a new event
//
// @rparam {group_id int, name string, description string, date string, location string }
func eventNewHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	defer recovery(w)
	err := error(nil)
	var inEvent EventDTOin
	inEvent.CreatorId, err = getRequestSenderID(r)
	if err != nil {
		jsonResponse(w, 401, err.Error())
		return
	}
	inEvent.CreatedAt = time.Now().Format("2006-01-02 15:04:05")
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err = decoder.Decode(&inEvent)
	if err != nil {
		log.Println(err.Error())
		jsonResponse(w, http.StatusUnprocessableEntity, " malformed json")
		return
	}
	// add the event to the table
	_, err = statements["addEvent"].Exec(inEvent.GroupID, inEvent.Name, inEvent.Description, inEvent.Date, inEvent.Location, inEvent.CreatorId, inEvent.CreatedAt)
	if err != nil {
		jsonResponse(w, http.StatusInternalServerError, " addEvent query failed")
		return
	}
	// add creator as going to the event
	// event_praticipants table has event_id, user_id, status, status_updated_at
	_, err = statements["addEventParticipant"].Exec(inEvent.GroupID, inEvent.CreatorId, "going", inEvent.CreatedAt)

	jsonResponse(w, 200, "Event created")
}

type EventDTOout struct {
	ID          int    `json:"id"`
	GroupID     int    `json:"group_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Date        string `json:"date"`
	Location    string `json:"location"`
	CreatorId   int    `json:"creator_id"`
	CreatedAt   string `json:"created_at"`
}

// # eventsGetHandler returns all events for a group
//
// @rparam {group_id int}
func eventsGetHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	defer recovery(w)
	_, err := getRequestSenderID(r)
	if err != nil {
		jsonResponse(w, 401, err.Error()+" you are not logged in")
		return
	}
	var data struct {
		GroupID int `json:"group_id"`
	}
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err = decoder.Decode(&data)
	if err != nil {
		log.Println(err.Error())
		jsonResponse(w, http.StatusUnprocessableEntity, " malformed json")
		return
	}
	rows, err := statements["getEvents"].Query(data.GroupID)
	if err != nil {
		log.Println(err.Error())
		jsonResponse(w, http.StatusInternalServerError, " getEvents query failed")
		return
	}
	var events []EventDTOout
	for rows.Next() {
		var event EventDTOout
		err = rows.Scan(&event.ID, &event.GroupID, &event.Name, &event.Description, &event.Date, &event.Location, &event.CreatorId, &event.CreatedAt)
		if err != nil {
			log.Println(err.Error())
			jsonResponse(w, http.StatusInternalServerError, " getEvents query failed")
			return
		}
		events = append(events, event)
	}
	jsonResponse(w, 200, events)
}

type ParticipantDTOout struct {
	ID              int    `json:"id"`
	FullName        string `json:"full_name"`
	Email           string `json:"email"`
	Picture         string `json:"picture"`
	Status          string `json:"status"`
	StatusUpdatedAt string `json:"status_updated_at"`
}

// # eventParticipantsGetHandler returns all participants for an event
//
// @rparam {event_id int}
func eventParticipantsGetHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	defer recovery(w)
	var data struct {
		EventID int `json:"event_id"`
	}
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&data)
	if err != nil {
		log.Println(err.Error())
		jsonResponse(w, http.StatusUnprocessableEntity, " malformed json")
		return
	}
	rows, err := statements["getEventParticipants"].Query(data.EventID)
	if err != nil {
		jsonResponse(w, http.StatusInternalServerError, " getEventParticipants query failed")
		return
	}
	var participants []ParticipantDTOout
	for rows.Next() {
		var participant ParticipantDTOout
		var firstName, lastName string
		picBlob := []byte{}
		rows.Scan(&participant.ID, &firstName, &lastName, &participant.Email, &picBlob, &participant.Status, &participant.StatusUpdatedAt)
		participant.FullName = firstName + " " + lastName
		participant.Picture = base64.StdEncoding.EncodeToString(picBlob)
		participants = append(participants, participant)
	}
	jsonResponse(w, 200, participants)
}

// # eventAttendHandler marks a user as going to an event
//
// @rparam {event_id int}
func eventAttendHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	defer recovery(w)
	userID, err := getRequestSenderID(r)
	if err != nil {
		jsonResponse(w, 401, err.Error())
		return
	}
	var data struct {
		EventID int `json:"event_id"`
	}
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err = decoder.Decode(&data)
	if err != nil {
		log.Println(err.Error())
		jsonResponse(w, http.StatusUnprocessableEntity, " malformed json")
		return
	}
	// delete the row where user_id = userID and event_id = data.EventID
	_, err = statements["updateEventParticipant"].Exec("going", time.Now().Format("2006-01-02 15:04:05"), data.EventID, userID)
	if err != nil {
		jsonResponse(w, http.StatusInternalServerError, " updateEventParticipant query failed")
		return
	}
	jsonResponse(w, 200, "User marked as going")
	return
}

// # eventNotAttendHandler marks a user as not going to an event
//
// @rparam {event_id int}
func eventNotAttendHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	defer recovery(w)
	userID, err := getRequestSenderID(r)
	if err != nil {
		jsonResponse(w, 401, err.Error())
		return
	}
	var data struct {
		EventID int `json:"event_id"`
	}
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err = decoder.Decode(&data)
	if err != nil {
		log.Println(err.Error())
		jsonResponse(w, http.StatusUnprocessableEntity, " malformed json")
		return
	}
	result, err := statements["updateEventParticipant"].Exec("notgoing", time.Now().Format("2006-01-02 15:04:05"), data.EventID, userID)
	if err != nil {
		log.Println("Error updating event participant status:", err)
		jsonResponse(w, http.StatusInternalServerError, " updateEventParticipant query failed")
		return
	} else {
		rowsAffected, _ := result.RowsAffected()
		log.Println("Rows affected:", rowsAffected)
	}
	jsonResponse(w, 200, "User marked as not going")
	return
}
