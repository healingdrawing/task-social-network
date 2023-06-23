package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"runtime/debug"
	"strings"
)

// # jsonResponse marshals and forwards json response writing to http.ResponseWriter
//
// @params {w http.ResponseWriter, statusCode int, message string}
// @sideEffect {jsonResponse -> w}
func jsonResponse(w http.ResponseWriter, statusCode int, message string) {
	jsonResponseObj, _ := json.Marshal(map[string]string{
		"message": http.StatusText(statusCode) + ": " + message,
	})
	w.WriteHeader(statusCode)
	w.Write(jsonResponseObj)
}

// # recovery is a utility function to recover from panic and send a json err response over http
//
// @sideEffect {log, debug}
//
// - for further debugging uncomment {print stack trace}
func recovery(w http.ResponseWriter) {
	if r := recover(); r != nil {
		fmt.Println("=====================================")
		stackTrace := debug.Stack()
		lines := strings.Split(string(stackTrace), "\n")
		relevantPanicLine := ""
		if len(lines) > 2 {
			relevantPanicLine = fmt.Sprintf("if in the same func, the panic occurred at below line (else, print full stack trace by going into utils func recovery)\n%s", lines[len(lines)-2])
		}
		log.Println(relevantPanicLine)
		jsonResponse(w, http.StatusInternalServerError, "internal server error"+"\n"+relevantPanicLine)
		fmt.Println("=====================================")
		// to print the full stack trace
		// log.Println(string(stackTrace))
	}
}

// # getRequestSenderID gets the ID of the request sender from the cookie
//
// @params {r *http.Request}
func getRequestSenderID(r *http.Request) (int, error) {
	cookie, err := r.Cookie("user_uuid")
	if err != nil {
		return 0, errors.New("malformed cookie/cookie not found")
	}

	requestSenderID, err := getIDbyUUID(cookie.Value)
	if err != nil {
		return 0, errors.New("failed to get ID of the request sender")
	}

	return requestSenderID, nil
}

// # getIDbyUUID retrieves ID of the user from uuid
//
// @params {UUID string}
// execute DB prepared statement getIDbyUUID.query
func getIDbyUUID(UUID string) (ID int, err error) {
	rows, err := statements["getIDbyUUID"].Query(UUID)
	if err != nil {
		return 0, err
	}
	defer rows.Close()
	rows.Next()
	err = rows.Scan(&ID)
	if err != nil {
		return 0, err
	}
	rows.Close()
	return ID, nil
}
