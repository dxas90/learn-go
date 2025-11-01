#!/bin/bash
set -e

echo "🚀 Starting learn-go application..."

# Export environment variables
export HOST=${HOST:-127.0.0.1}
export PORT=${PORT:-8080}
export GO_ENV=${GO_ENV:-development}
export APP_VERSION=${APP_VERSION:-0.0.1}

# Build and run
go build -o bin/learn-go ./cmd/api
./bin/learn-go
