package main

import (
	"encoding/json"
	"net/http"
)

/*
jsonResponseWriterManager is
WriteHeader(statusCode) + Write(jsonResponse) + manage of "message marshlling to jsonResponse" error
*/
func jsonResponseWriterManager(w http.ResponseWriter, statusCode int, message string) {
	jsonResponse, err := json.Marshal(map[string]string{
		"message": http.StatusText(statusCode) + ": " + message,
	})
	if err == nil {
		w.WriteHeader(statusCode)
		_, err = w.Write(jsonResponse)
		if err != nil { // if error in time of writing jsonResponse
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("fuck off, golang, with this shitty infinity error nesting"))
			// fuck off golang, with this shit infinity error nesting
		}
	} else { // if error in time of marshalling error message
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("fuck off, golang, with this shitty infinity error nesting"))
	}
}
