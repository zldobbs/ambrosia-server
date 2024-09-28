// Store the ambrosia handlers
package main

import (
	"fmt"
	"log"
	"net/http"
)

// Helper function to write an HTTP response with error handling.
//
// Parameters:
// 	- w: ResponseWriter object to write to
//	- res: String to write to the response
func safeResponseWrite(w http.ResponseWriter, responseString string) {
	if _, err := fmt.Fprint(w, responseString); err != nil {
		log.Println("Error writing response:", err)
	}
}

// Handle requests to root index.
//
// Parameters:
// 	- w: ResponseWriter object
//	- r: Incoming request
func indexHandler(w http.ResponseWriter, r *http.Request) {
	safeResponseWrite(w, "Welcome to the server!")
}

// Respond to heartbeat requests to indicate server is reachable and up as expected
// NOTE: This is likely redundant to IndexHandler...
//
// Parameters:
// 	- w: ResponseWriter object
//	- r: Incoming request
func heartbeatHandler(w http.ResponseWriter, r *http.Request) {
	safeResponseWrite(w, "OK")
}

// GraphQL endpoint.
// NOTE: This should be moved behind an authentication middleware.
//
// Parameters:
// 	- w: ResponseWriter object
//	- r: Incoming request
func graphQLHandler(w http.ResponseWriter, r *http.Request) {
	safeResponseWrite(w, "GraphQL is not in yet!")
}
