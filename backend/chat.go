package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
)

type Users struct {
	Users []User `json:"users"`
}

type MessagesRequest struct {
	Username  string `json:"username"`
	OtherUser string `json:"otheruser"`
}

type UsernameData struct {
	Username string `json:"username"`
}

type MessagesResponse struct {
	Messages []Message `json:"messages"`
}

type Message struct {
	UsernameFrom string    `json:"usernameFrom"`
	UsernameTo   string    `json:"usernameTo"`
	Content      string    `json:"content"`
	Time         time.Time `json:"time"`
}

type User struct {
	ID          int       `json:"ID"`
	Username    string    `json:"username"`
	LastMessage time.Time `json:"time"`
	Online      bool      `json:"online"`
}

type Typing struct {
	UsernameFrom string `json:"usernameFrom"`
	UsernameTo   string `json:"usernameTo"`
	Typing       bool   `json:"typing"`
}

func chatTypingHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	defer func() {
		if err := recover(); err != nil {
			log.Println(err)
			jsonResponseWriterManager(w, http.StatusInternalServerError, "recover - chatTypingHandler")
		}
	}()
	var data Typing
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&data)
	if err != nil {
		log.Println(err.Error())
		jsonResponseWriterManager(w, http.StatusBadRequest, "")
		return
	}
	jsonResponseWriterManager(w, http.StatusOK, "")
	sendTyping(data)
}

func chatMessagesHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	defer func() {
		if err := recover(); err != nil {
			log.Println(err)
			jsonResponseWriterManager(w, http.StatusInternalServerError, "recover - chatMessagesHandler")
		}
	}()
	var data MessagesRequest
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&data)
	if err != nil {
		log.Println(err.Error())
		jsonResponseWriterManager(w, http.StatusBadRequest, "")
		return
	}
	ID1, err := getIDbyEmail(data.Username)
	if err != nil {
		jsonResponseWriterManager(w, http.StatusInternalServerError, "ID1")
		return
	}
	ID2, err := getIDbyEmail(data.OtherUser)
	if err != nil {
		jsonResponseWriterManager(w, http.StatusInternalServerError, "ID2")
		return
	}
	rows, err := statements["getMessages"].Query(ID1, ID2, ID2, ID1)
	if err != nil {
		jsonResponseWriterManager(w, http.StatusInternalServerError, "getMessages query failed")
		return
	}
	var messages MessagesResponse
	var IDpairs [][2]int
	for rows.Next() {
		var message Message
		var (
			IDFrom int
			IDTo   int
		)
		err = rows.Scan(&IDFrom, &IDTo, &message.Content, &message.Time)
		if err != nil {
			jsonResponseWriterManager(w, http.StatusInternalServerError, "getMessages rows.Scan for loop failed")
			return
		}
		IDpairs = append(IDpairs, [2]int{IDFrom, IDTo})
		messages.Messages = append(messages.Messages, message)
	}
	rows.Close()
	for i := 0; i < len(messages.Messages); i++ {
		messages.Messages[i].UsernameFrom, err = getUserEmailbyID(IDpairs[i][0])
		if err != nil {
			jsonResponseWriterManager(w, http.StatusInternalServerError, "getUserEmailbyID IDpairs[i][0]")
			return
		}
		messages.Messages[i].UsernameTo, err = getUserEmailbyID(IDpairs[i][1])
		if err != nil {
			jsonResponseWriterManager(w, http.StatusInternalServerError, "getUserEmailbyID IDpairs[i][1]")
			return
		}
	}
	jsonResponseWriterManager(w, http.StatusOK, "")
}

func chatUsersHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	defer func() {
		if err := recover(); err != nil {
			log.Println(err)
			jsonResponseWriterManager(w, http.StatusInternalServerError, "recover - chatUsersHandler")
		}
	}()
	var data UsernameData
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&data)
	if err != nil {
		log.Println(err.Error())
		jsonResponseWriterManager(w, http.StatusBadRequest, "")
		return
	}
	rows, err := statements["getAllUsers"].Query()
	if err != nil {
		jsonResponseWriterManager(w, http.StatusInternalServerError, "getAllUsers query failed")
		return
	}
	var users Users
	for rows.Next() {
		var user User
		rows.Scan(&user.ID, &user.Username)
		user.Online = isOnline(user.Username)
		users.Users = append(users.Users, user)
	}
	rows.Close()
	ID1, err := getIDbyEmail(data.Username)
	if err != nil {
		jsonResponseWriterManager(w, http.StatusInternalServerError, "getIDbyEmail ID1 failed")
		return
	}
	for i := 0; i < len(users.Users); i++ {
		ID2, err := getIDbyEmail(users.Users[i].Username)
		if err != nil {
			jsonResponseWriterManager(w, http.StatusInternalServerError, "getIDbyEmail ID2 failed")
			return
		}
		rows, err := statements["getMessages"].Query(ID1, ID2, ID2, ID1)
		if err != nil {
			jsonResponseWriterManager(w, http.StatusInternalServerError, "getMessages query failed")
			return
		}
		defer rows.Close()
		rows.Next()
		var (
			IDFrom int
			IDTo   int
			Text   string
		)
		err = rows.Scan(&IDFrom, &IDTo, &Text, &users.Users[i].LastMessage)
		rows.Close()
		if err != nil {
			jsonResponseWriterManager(w, http.StatusInternalServerError, "getMessages scan failed")
			return
		}
	}
	jsonResponse, _ := json.Marshal(users)
	if err != nil {
		jsonResponseWriterManager(w, http.StatusInternalServerError, "json.Marshal users failed")
		return
	}

	w.WriteHeader(200) // todo: CHECK IT!!! NECESSARILY. very muddy place of code, compare with old version of code! Also WriteHeader docs says if it was not called before w.Write the status 200 will be used by default. So, probably, it is not necessary to call it at all.
	_, err = w.Write(jsonResponse)
	if err != nil {
		jsonResponseWriterManager(w, http.StatusInternalServerError, "w.Write jsonResponse failed")
		return
	}

}

func isOnline(username string) (found bool) {
	clients.Range(func(_, value interface{}) bool {
		if c, ok := value.(string); ok && c == username {
			found = true
		}
		return true
	})
	return
}

func chatNewHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	defer func() {
		if err := recover(); err != nil {
			log.Println(err)
			jsonResponseWriterManager(w, http.StatusInternalServerError, "recover - chatNewHandler")
		}
	}()
	var data Message
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&data)
	if err != nil {
		log.Println(err.Error())
		jsonResponseWriterManager(w, http.StatusBadRequest, "")
		return
	}
	fromID, err := getIDbyEmail(data.UsernameFrom)
	if err != nil {
		log.Println(err.Error())
		jsonResponseWriterManager(w, http.StatusInternalServerError, "getIDbyEmail fromID failed")
		return
	}
	toID, err := getIDbyEmail(data.UsernameTo)
	if err != nil {
		log.Println(err.Error())
		jsonResponseWriterManager(w, http.StatusInternalServerError, "getIDbyEmail toID failed")
		return
	}
	data.Time = time.Now()
	_, err = statements["addMessage"].Exec(fromID, toID, data.Content, data.Time)
	if err != nil {
		jsonResponseWriterManager(w, http.StatusInternalServerError, "addMessage query failed")
		return
	}
	jsonResponseWriterManager(w, http.StatusOK, "")
	sendMessage(data)
}
