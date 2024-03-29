package main

import (
	"fmt"
	"log"
	"net/http"
)

// === Route handlers ===

// Handle requests to root index
func indexHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: Debug middleware to log when routes are hit
	fmt.Println("Index hit")
	fmt.Fprintf(w, "Welcome to the server!")
}

// Respond to heartbeat requests to indicate server is reachable and up as expected
func heartbeatHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Heartbeat hit")
	fmt.Fprintf(w, "OK")
}

func main() {
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/heartbeat", heartbeatHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
