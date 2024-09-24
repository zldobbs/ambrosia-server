package main

import (
	"log"
	"net/http"
	"time"
)

// Middleware function that logs when a URL is requested.
//
// Parameters:
// 	- next: Next HTTP handler to call after this
//
// Returns:
// 	- This handler function as middleware
func logMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		log.Printf("Started request: %s %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
		log.Printf("Completed in: %s", time.Since(start))
	})
}

// Middleware function for handling authentication.
// If authentication fails, reject proceeding to the next handler.
//
// Parameters:
// 	- next: Next HTTP handler to call after this
//
// Returns:
// 	- This handler function as middleware
func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// TODO: Implement actual authorization checks
		if false {
			http.Error(w, "Forbidden: Access Denied", http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// Main entrypoint; will handle launching the HTTP server.
func main() {
	// TODO: Consider using gorilla/mux
	mux := http.NewServeMux()

	// Public routes
	mux.HandleFunc("/", indexHandler)
	mux.HandleFunc("/heartbeat", heartbeatHandler)

	// Protected routes
	mux.Handle("/graphql", authMiddleware(http.HandlerFunc(graphQLHandler)))

	// Wrap all handlers with logging middleware
	loggedMux := logMiddleware(mux)

	log.Fatal(http.ListenAndServe(":8080", loggedMux))
}
