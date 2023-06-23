package main

import (
	"encoding/json"
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
