# amborosia-server

Backend server for Ambrosia.

## Setup

1. Run `go mod tidy` to get dependencies

### GraphQL Code Gen

1. Run `go generate ./..` within the graph/ directory to update based on changes to [graph/schema.graphqls](./graph/schema.graphqls)

### Server Build and Launching

> TODO: Put the server in a container, launch with compose.yaml

1. Define the environment variables required for connecting to the server:

- `POSTGRES_USER`: Database user (e.g. `postgres`)
- `POSTGRES_PASSWORD`: Database user password (e.g. `postgres`)
- `POSTGRES_DB`: Database name (e.g. `ambrosia`)
- `POSTGRES_HOST`: Database hostname (e.g. `localhost`)
- `POSTGRES_PORT`: Database port (e.g. `5342`)

1. Run `go build .`
1. Run `./ambrosia-server`
1. Navigate to <http://localhost:8080> to see the server running

## Database Setup

Scripts are provided to help setup the expected tables and seed data.
These can be found in [/db/sql](./db/sql)

1. Run initialize.sql within a `psql` terminal

  > NOTE: While testing may need to recreate the database from scratch.
  > If this is the case, first drop the existing database with: `DROP DATABASE ambrosia;` before running [initialize.sql](./db/sql/initialize.sql).

1. Additionally run the seed.sql script to populate the ambrosia database with some sample data.
