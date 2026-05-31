#!/bin/sh
set -e

DB_URL="mysql://${DB_USER}:${DB_PASSWORD}@tcp(${DB_HOST}:${DB_PORT})/${DB_NAME}"

echo "[entrypoint] Running database migrations..."
migrate -path /app/db/migrations -database "$DB_URL" up
echo "[entrypoint] Migrations done."

exec "$@"
