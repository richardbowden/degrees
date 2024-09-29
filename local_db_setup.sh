#!/bin/bash

# PostgreSQL connection details
PG_HOST="localhost"
PG_PORT="5432"
PG_SUPERUSER="postgres"  # an existing superuser to run these commands
PG_DATABASE="postgres"  # the database you want to grant connect privileges to

# New user details
NEW_USERNAME="p402"
NEW_USER_PASSWORD="letmein"
DB_NAME="p402"

psql -c "CREATE DATABASE $DB_NAME;"

# Create the new role without superuser privileges
psql -h $PG_HOST -p $PG_PORT -U $PG_SUPERUSER -d $DB_NAME -c "CREATE ROLE $NEW_USERNAME WITH LOGIN PASSWORD '$NEW_USER_PASSWORD';"

# Grant superuser privileges to the new role
psql -h $PG_HOST -p $PG_PORT -U $PG_SUPERUSER -d $DB_NAME -c "ALTER ROLE $NEW_USERNAME WITH SUPERUSER;"

# Grant connect privileges to the database
psql -h $PG_HOST -p $PG_PORT -U $PG_SUPERUSER -d $DB_NAME -c "GRANT CONNECT ON DATABASE $DB_NAME TO $NEW_USERNAME;"

echo "User $NEW_USERNAME has been created with superuser privileges and can now connect to $DB_NAME."


