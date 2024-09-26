package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/zldobbs/ambrosia-server/db"
	"github.com/zldobbs/ambrosia-server/graph"
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

const defaultPort = "8080"

// Main entrypoint; will handle launching the HTTP server.
func main() {
	// TODO: Consider using gorilla/mux
	mux := http.NewServeMux()

	port := os.Getenv("AMBROSIA_BACKEND_PORT")
	if port == "" {
		port = defaultPort
	}

	// Get database connection pool
	pool := db.InitDB()

	// GraphQL Server (using gqlgen)
	// TODO: Setup "playground" for testing (or use postman)
	gql_server := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{DB_POOL: pool}}))

	// Public routes
	mux.HandleFunc("/", indexHandler)
	mux.HandleFunc("/heartbeat", heartbeatHandler)

	// Protected routes
	mux.Handle("/graphql", authMiddleware(gql_server))

	// Wrap all handlers with logging middleware
	loggedMux := logMiddleware(mux)

	log.Fatal(http.ListenAndServe(":"+port, loggedMux))
}

func InitDB() {
	panic("unimplemented")
}
