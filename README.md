# amborosia-server

Backend server for Ambrosia.

## Setup

1. Run `go mod tidy` to get dependencies

### GraphQL Code Gen

1. Run `go generate ./..` within the graph/ directory to update based on changes to [graph/schema.graphqls](./graph/schema.graphqls)

### Server Build and Launching

1. Define the environment variables required for connecting to the server:

- `POSTGRES_USER`: Database user (e.g. `postgres`)
- `POSTGRES_PASSWORD`: Database user password (e.g. `postgres`)
- `POSTGRES_DB`: Database name (e.g. `ambrosia`)
- `POSTGRES_HOST`: Database hostname (e.g. `localhost`)
- `POSTGRES_PORT`: Database port (e.g. `5342`)

1. Run `go build .`
1. Run `./ambrosia-server`
1. Navigate to <http://localhost:8080> to see the server running

## Database

The server expects to connect to a Postgres database with connection information corresponding to the environment variables defined above.

### Migrations

#### Creating Migrations

Database migrations track changes over time to the database schema.
Any changes to the schema should be made via migrations.

1. Run `db/create_migration.sh <migration_name>` to generate migration files within db/migrations
1. Add appropriate SQL logic to the generated files to perform the migration

- The "up" file should include logic to apply the schema change (e.g. adding a column)
- The "down" file should include logic to revert the schema change (e.g. removing the column that was added)

#### Applying Migrations

Run the following script to apply any pending migrations to the database:

```sh
db/apply_migrations.sh
```

### Seeding

Populate initial data to the database using db/sql/seed.sql:

```bash
psql -U postgres -f db/sql/seed.sql
```
