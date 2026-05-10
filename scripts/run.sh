#!/bin/bash
set -e

# load .env
export $(grep -v '^#' .env | xargs)

# run migrations
goose -dir migrations sqlite3 "$DB_PATH" up

# start server
go run ./cmd/api