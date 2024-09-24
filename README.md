# amborosia-server

Backend server for Ambrosia.

## Setup

1. Run `go mod tidy` to get dependencies

### GraphQL Code Gen

1. Run

### Server Build and Launching

> TODO: Put the server in a container, launch with compose.yaml

1. Run `go build .`
1. Run `./ambrosia-server`
1. Navigate to <http://localhost:8080> to see the server running

## Database Setup

1. Run db/initialize.sql within a `psql` terminal

  > NOTE: While testing may need to recreate the database from scratch.
  > If this is the case, first drop the existing database with: `DROP DATABASE ambrosia;` before running [db/initialize.sql](./db/initialize.sql).
