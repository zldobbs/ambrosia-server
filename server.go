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
//   - next: Next HTTP handler to call after this
//
// Returns:
//   - This handler function as middleware
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
//   - next: Next HTTP handler to call after this
//
// Returns:
//   - This handler function as middleware
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

// Middleware function for enabling CORS on target routes
//
// Parameters:
//   - next: Next HTTP handler to call after this
//
// Returns:
//   - This handler function as middleware
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// TODO: CORS should be restricted in production!
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// Handle preflight requests
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
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
	defer pool.Close()

	// GraphQL Server (using gqlgen)
	gql_server := handler.NewDefaultServer(
		graph.NewExecutableSchema(
			graph.Config{Resolvers: &graph.Resolver{DB_POOL: pool}},
		),
	)

	// Public routes
	mux.HandleFunc("/", indexHandler)
	mux.HandleFunc("/heartbeat", heartbeatHandler)

	// Protected routes
	mux.Handle("/graphql", authMiddleware(gql_server))

	// Wrap all handlers with logging and cors middleware
	loggedMux := logMiddleware(mux)
	corsLoggedMux := corsMiddleware(loggedMux)

	// TODO: Consider running the server in a goroutine for better logging here
	log.Println("ambrosia server starting now, no output means OK")
	log.Fatal(http.ListenAndServe(":"+port, corsLoggedMux))
}
