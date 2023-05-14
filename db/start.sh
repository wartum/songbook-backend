#!/bin/sh

password="Songs"

docker run -d -e POSTGRES_PASSWORD="$password" -p 5432:5432 postgres:15.2-alpine
sleep 2
PGPASSWORD="$password" psql -h localhost -p 5432 -U postgres -f create_db.sql
PGPASSWORD="$password" psql -h localhost -p 5432 -U postgres -d songs -f create_tables.sql
PGPASSWORD="$password" psql -h localhost -p 5432 -U postgres -d songs -f insert_users.sql
