package graph

//go:generate go run github.com/99designs/gqlgen generate

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

import (
	"github.com/jackc/pgx/v5/pgxpool"
)

type Resolver struct {
	DB_POOL *pgxpool.Pool
}

// Create a new Resolver with the SQL database connection
//
// Parameters:
// 	- db: Connection to SQL database
//
// Returns:
// 	A GraphQL Resolver object with a connection to the SQL DB.
func NewResolver(pool *pgxpool.Pool) *Resolver {
	return &Resolver{DB_POOL: pool}
}
