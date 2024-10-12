#!/bin/bash

function die() {
    echo "$1" && exit 1
}

[ ! -z $POSTGRES_DB ] || die "POSTGRES_DB not set"
[ ! -z $POSTGRES_USER ] || die "POSTGRES_USER not set"
[ ! -z $POSTGRES_PASSWORD ] || die "POSTGRES_PASSWORD not set"
[ ! -z $POSTGRES_HOST ] || die "POSTGRES_HOST not set"
[ ! -z $POSTGRES_PORT ] || die "POSTGRES_PORT not set"

function execute_sql() {
    PGPASSWORD=$POSTGRES_PASSWORD psql -U $POSTGRES_USER -h $POSTGRES_HOST -p $POSTGRES_PORT -d $POSTGRES_DB -tAc "$1"
}

function execute_sql_file() {
    PGPASSWORD=$POSTGRES_PASSWORD psql -U $POSTGRES_USER -h $POSTGRES_HOST -p $POSTGRES_PORT -d $POSTGRES_DB -f "$1"
}

# FIXME: If POSTGRES_DB does not exist, all commands will fail.

# Create a table to track which migrations have been applied so far
execute_sql "CREATE TABLE IF NOT EXISTS schema_migrations (version VARCHAR(255) PRIMARY KEY, applied_at TIMESTAMPTZ NOT NULL DEFAULT NOW());"

# Loop through migration "up" files in order and apply any pending migrations
# NOTE: This script will apply all "up" migrations in order
#       To be more granular or apply "down" migrations, run the SQL scripts manually.
migrations_dir=$(dirname "$0")/migrations
for file in $migrations_dir/*_up.sql; do
    filename=$(basename -- "$file")
    version="${filename%.*}"

    # Check if this version has already been applied within the database
    is_applied=$(execute_sql "SELECT 1 FROM schema_migrations WHERE version='$version';")
    if [ "$is_applied" != "1" ]; then
        echo "Applying $version..."
        # FIXME: Check if these commands are successful, don't insert into schema_migrations if not
        execute_sql_file "$file"
        execute_sql "INSERT INTO schema_migrations (version) VALUES ('$version');"
        echo "Applied $version!"
    else
        echo "Skipping $version; already applied..."
    fi
done
