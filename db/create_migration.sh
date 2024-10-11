#!/bin/bash

# Create a new migration entry for the database
# This will create two new SQL files, one for up/down respectively
# The SQL for the migration itself will need to be manually written in the created files.

name=$1
if [ -z $name ]; then
    echo "No sequence name provided"
    exit 1
fi

ts=$(date +%Y%m%d%H%M%S)
ts_name="${ts}_${name}"
echo "-- ${ts_name}_up" > $(dirname "$0")/migrations/${ts_name}_up.sql
echo "-- ${ts_name}_down" > $(dirname "$0")/migrations/${ts_name}_down.sql

echo "Created migration $ts_name"
