#!/usr/bin/env bash

# -----------------------------

# Backend Alpha Development Environment

# -----------------------------

export PORT=8080

export DB_USER=trading_user
export DB_PASSWORD=password
export DB_HOST=localhost
export DB_PORT=5432
export DB_NAME=trading_db

# Convenience connection string

export DATABASE_URL="postgres://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}"

# Integration test database

export TEST_DATABASE_URL="postgres://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/backend_alpha_test"

echo "Backend Alpha environment loaded."
echo ""
echo "DATABASE_URL=${DATABASE_URL}"
echo "TEST_DATABASE_URL=${TEST_DATABASE_URL}"

