#!/bin/sh
set -e

echo "Waiting for PostgreSQL to start..."
until PGPASSWORD=$DB_PASSWORD pg_isready -h "$DB_HOST" -U "$DB_USER" -d "$DB_NAME"; do
  sleep 1
done

echo "Running migrations..."
PGPASSWORD=$DB_PASSWORD psql -v ON_ERROR_STOP=1 -h "$DB_HOST" -U "$DB_USER" -d "$DB_NAME" -f /app/migrations/001_init.up.sql

echo "Migrations completed!"