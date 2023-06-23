package main

import (
	"encoding/json"
	"net/http"
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
