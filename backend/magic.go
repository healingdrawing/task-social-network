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
	if err != nil { // if error in time of marshalling error message
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("\"" + http.StatusText(http.StatusInternalServerError) + "\" in time of json marshalling of message: \"" + message + "\"\nMarshalling error.Error() is \"" + err.Error() + "\""))
	} else {
		w.WriteHeader(statusCode)
		w.Write(jsonResponse)
	}
}
