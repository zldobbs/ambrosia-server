#!/bin/bash

# Create a new migration entry for the database
# This will create two new SQL files, one for up/down respectively
# The SQL for the migration itself will need to be manually written in the created files.

name=$1
if [ -z $name ]; then
    echo "No sequence name provided"
    exit 1
fi

migrations_dir=$(dirname "$0")/migrations

# Versioning via version numbers
latest_version=$(ls $migrations_dir | tail -n -1 | cut -d '_' -f1)
if [ -z $latest_version ]; then
    latest_version=0
fi
version=$((latest_version + 1))

# Versioning via timestamp
ts=$(date +%Y%m%d%H%M%S)

versioned_name="${ts}_${name}"
echo "-- ${versioned_name}_up" > $migrations_dir/${versioned_name}_up.sql
echo "-- ${versioned_name}_down" > $migrations_dir/${versioned_name}_down.sql

echo "Created migration $versioned_name"
