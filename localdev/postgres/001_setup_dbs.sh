#!/usr/bin/env bash

echo "Creating these databases: $EXTRA_DBS"

for DB in $(echo $EXTRA_DBS | tr ',' ' '); do
  echo "Creating database: $DB"
  psql -U $POSTGRES_USER template1 -c "CREATE DATABASE $DB"
done