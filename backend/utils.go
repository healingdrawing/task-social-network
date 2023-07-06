package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"runtime/debug"
	"strings"
	"time"
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

func get_user_id_by_email(email string) (user_id int, err error) {
	rows, err := statements["getUserIDByEmail"].Query(email)
	if err != nil {
		return 0, err
	}
	defer rows.Close()
	rows.Next()
	err = rows.Scan(&user_id)
	if err != nil {
		return 0, err
	}
	rows.Close()
	return user_id, nil
}

func get_email_by_user_id(user_id int) (email string, err error) {
	rows, err := statements["getEmailByID"].Query(user_id)
	if err != nil {
		return "", err
	}
	defer rows.Close()
	rows.Next()
	err = rows.Scan(&email)
	if err != nil {
		return "", err
	}
	rows.Close()
	return email, nil
}

// # getRequestSenderID gets the ID of the request sender from the cookie
//
// @params {r *http.Request}
func getRequestSenderID(r *http.Request) (int, error) {
	cookie, err := r.Cookie("user_uuid")
	if err != nil {
		return 0, errors.New("malformed cookie/cookie not found")
	}

	requestSenderID, err := get_user_id_by_uuid(cookie.Value)
	if err != nil {
		return 0, errors.New("failed to get ID of the request sender")
	}

	return requestSenderID, nil
}

// # get_user_id_by_uuid retrieves id of the user from uuid
//
// @params {uuid string}
// execute DB prepared statement get_user_id_by_uuid.query
func get_user_id_by_uuid(uuid string) (user_id int, err error) {
	rows, err := statements["getIDbyUUID"].Query(uuid)
	if err != nil {
		return 0, err
	}
	defer rows.Close()
	rows.Next()
	err = rows.Scan(&user_id)
	if err != nil {
		return 0, err
	}
	rows.Close()
	return user_id, nil
}

// extractImageData extracts image data from dataURI
//
// @params {dataURI string}
//
// it will split the dataURI into two parts, using "," as the delimiter, then return the second part
func extractImageData(dataURI string) (string, error) {
	parts := strings.SplitN(dataURI, ",", 2)
	if len(parts) != 2 {
		return "", errors.New("invalid dataURI")
	}
	return parts[1], nil
}

func isImage(data []byte) bool {
	if len(data) < 4 {
		log.Println("len(data) < 4")
		return false
	}

	switch {
	case bytes.HasPrefix(data, []byte{0xFF, 0xD8, 0xFF}): // JPEG
		return true
	case bytes.HasPrefix(data, []byte{0x89, 0x50, 0x4E, 0x47}): // PNG
		return true
	case bytes.HasPrefix(data, []byte{0x47, 0x49, 0x46, 0x38}): // GIF
		return true
	}
	return false

	// if data[0] == 0xFF && data[1] == 0xD8 && data[2] == 0xFF {
	// 	return true // JPEG
	// }

	// if data[0] == 0x89 && data[1] == 0x50 && data[2] == 0x4E && data[3] == 0x47 {
	// 	return true // PNG
	// }

	// if data[0] == 0x47 && data[1] == 0x49 && data[2] == 0x46 && data[3] == 0x38 {
	// 	return true // GIF
	// }

	// return
}

// randomNum returns a random number between min and max, both inclusive.
func randomNum(min, max int) int {
	rng := rand.New(rand.NewSource(time.Now().Unix()))
	rng.Seed(time.Now().Unix())
	return rng.Intn(max+1-min) + min
}
