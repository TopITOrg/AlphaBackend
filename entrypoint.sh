#!/bin/sh
set -e

goose -dir ./migrations postgres "host=${Config_HOST} port=${Config_PORT} database=${Config_DB_NAME} user=${Config_USERNAME} password=${Config_PASSWORD} sslmode=disable" up

exec /app/go-server
