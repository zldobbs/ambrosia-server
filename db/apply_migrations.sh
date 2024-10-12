#!/bin/bash

# Apply pending database migrations
# NOTE: This script will apply all "up" migrations in order
#       To be more granular or apply "down" migrations, run the SQL scripts manually.

function die() {
    echo "$1" && exit 1
}

[ ! -z $POSTGRES_DB ] || die "POSTGRES_DB not set"
[ ! -z $POSTGRES_USER ] || die "POSTGRES_USER not set"
[ ! -z $POSTGRES_PASSWORD ] || die "POSTGRES_PASSWORD not set"
[ ! -z $POSTGRES_HOST ] || die "POSTGRES_HOST not set"
[ ! -z $POSTGRES_PORT ] || die "POSTGRES_PORT not set"

function execute_sql_default() {
    # Execute SQL command using default Postgres database, instead of the one defined by POSTGRES_DB
    PGPASSWORD=$POSTGRES_PASSWORD psql -U $POSTGRES_USER -h $POSTGRES_HOST -p $POSTGRES_PORT -tAc "$1"
}

function execute_sql() {
    # Execute SQL command using connection variables defined within environment
    PGPASSWORD=$POSTGRES_PASSWORD psql -U $POSTGRES_USER -h $POSTGRES_HOST -p $POSTGRES_PORT -d $POSTGRES_DB -tAc "$1"
}

function execute_sql_file() {
    # Execute SQL script file using connection variables defined within environment
    PGPASSWORD=$POSTGRES_PASSWORD psql -U $POSTGRES_USER -h $POSTGRES_HOST -p $POSTGRES_PORT -d $POSTGRES_DB -f "$1"
}

# If POSTGRES_DB does not exist, all commands will fail.
db_exists=$(execute_sql_default "SELECT 1 FROM pg_database WHERE datname='$POSTGRES_DB';")
if [ "$db_exists" != "1" ]; then
    echo "Creating $POSTGRES_DB database..."
    execute_sql_default "CREATE DATABASE $POSTGRES_DB;" || die "Failed to create $POSTGRES_DB"
fi

# Create a table to track which migrations have been applied so far
execute_sql "CREATE TABLE IF NOT EXISTS schema_migrations (version VARCHAR(255) PRIMARY KEY, applied_at TIMESTAMPTZ NOT NULL DEFAULT NOW());"

# Loop through migration "up" files in order and apply any pending migrations
migrations_dir=$(dirname "$0")/migrations
for file in $migrations_dir/*_up.sql; do
    filename=$(basename -- "$file")
    version="${filename%.*}"

    # Check if this version has already been applied within the database
    is_applied=$(execute_sql "SELECT 1 FROM schema_migrations WHERE version='$version';")
    [ "$?" != "0" ] && die "Failed to query schema_migrations"
    if [ "$is_applied" != "1" ]; then
        echo "Applying $version..."
        execute_sql_file "$file" || die "Failed to apply migration $version ($file)"
        execute_sql "INSERT INTO schema_migrations (version) VALUES ('$version');" || die "Failed to update schema_migrations table for $version"
        echo "Applied $version!"
    else
        echo "Skipping $version; already applied..."
    fi
done
