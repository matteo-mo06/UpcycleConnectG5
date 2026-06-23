#!/bin/sh
set -e

DB_URL="mysql://${DB_USER}:${DB_PASSWORD}@tcp(${DB_HOST}:${DB_PORT})/${DB_NAME}"

echo "[entrypoint] Running database migrations..."
TRIES=0
until migrate -path /app/db/migrations -database "$DB_URL" up; do
  TRIES=$((TRIES + 1))
  if [ "$TRIES" -ge 15 ]; then
    echo "[entrypoint] Migrations failed after 15 attempts, aborting."
    exit 1
  fi
  echo "[entrypoint] Retry $TRIES/15 — database not ready yet, waiting 4s..."
  sleep 4
done
echo "[entrypoint] Migrations done."

echo "[entrypoint] Running seed..."
./seed || echo "[entrypoint] Seed skipped (already applied or error ignored)."

exec "$@"
