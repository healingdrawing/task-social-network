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
// @params {w http.ResponseWriter, statusCode int, data any}
// @sideEffect {jsonResponse -> w}
func jsonResponse(w http.ResponseWriter, statusCode int, data any) {
	jsonResponseObj := []byte{}
	// if data type is string
	if message, ok := data.(string); ok {
		jsonResponseObj, _ = json.Marshal(map[string]string{
			"message": http.StatusText(statusCode) + ": " + message,
		})
	}
	// if data type is int
	if message, ok := data.(int); ok {
		jsonResponseObj, _ = json.Marshal(map[string]int{
			"message": message,
		})
	}
	// if data type is bool
	if message, ok := data.(bool); ok {
		jsonResponseObj, _ = json.Marshal(map[string]bool{
			"message": message,
		})
	}
	// if data type is slice
	if _, ok := data.([]any); ok {
		jsonResponseObj, _ = json.Marshal(map[string][]any{
			"data": data.([]any),
		})
	}
	// if data type is object
	if _, ok := data.(map[string]any); ok {
		jsonResponseObj, _ = json.Marshal(map[string]any{
			"data": data.(map[string]any),
		})
	}
	// if unhandled by above custom conversion
	if len(jsonResponseObj) == 0 {
		w.WriteHeader(statusCode)
		err := json.NewEncoder(w).Encode(data)
		if err != nil {
			log.Println(err.Error())
		}
		return
	}
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
		relevantPanicLines := []string{}
		for _, line := range lines {
			if strings.Contains(line, "backend/") {
				relevantPanicLines = append(relevantPanicLines, line)
			}
		}
		if len(relevantPanicLines) > 1 {
			for i, line := range relevantPanicLines {
				if strings.Contains(line, "utils.go") {
					relevantPanicLines = append(relevantPanicLines[:i], relevantPanicLines[i+1:]...)
				}
			}
		}
		relevantPanicLine := strings.Join(relevantPanicLines, "\n")
		log.Println(relevantPanicLines)
		jsonResponse(w, http.StatusInternalServerError, relevantPanicLine)
		fmt.Println("=====================================")
		// to print the full stack trace
		log.Println(string(stackTrace))
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
