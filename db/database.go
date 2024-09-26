package db

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Initialize a connection pool for the Ambrosia database.
// Using pgx to connect to Postgres.
// Expects to find connection information in environment:
// 	- POSTGRES_USER: Database user username
//	- POSTGRES_PASSWORD: Database user password
// 	- POSTGRES_DB: Database name
// 	- POSTGRES_HOST: Database hostname
// 	- POSTGRES_PORT: Database port
//
// Returns:
// 	- Connection pool for configured database
func InitDB() *pgxpool.Pool {
	// Load configuration from the environment
	user := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	dbname := os.Getenv("POSTGRES_DB")
	host := os.Getenv("POSTGRES_HOST")
	port := os.Getenv("POSTGRES_PORT")

	// Connection string
	connString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		user,
		password,
		host,
		port,
		dbname,
	)

	// Open a postgres connection pool
	pool, err := pgxpool.New(context.Background(), connString)
	if err != nil {
		log.Fatalf("Failed to connect to database, error: %s", err)
	}

	// Test connection
	err = pool.Ping(context.Background())
	if err != nil {
		log.Fatalf("Failed to ping database, error: %s", err)
	}

	return pool
}
