#!/bin/bash
# start-backend.sh — Start the Go backend with optional .env loading

set -e
SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
cd "$SCRIPT_DIR/backend"

echo "╔══════════════════════════════════════╗"
echo "║       ShopApp Backend (Go + MySQL)   ║"
echo "╚══════════════════════════════════════╝"

# Load .env if it exists
if [ -f ".env" ]; then
  echo "📄 Loading .env..."
  export $(grep -v '^#' .env | grep -v '^$' | xargs)
else
  echo "⚠️  No .env found — using defaults (DB_PASSWORD=empty)"
  echo "   Copy backend/.env.example to backend/.env and set your MySQL password."
fi

echo "📦 Downloading Go dependencies..."
go mod tidy

echo ""
echo "🚀 Starting server on http://localhost:${SERVER_PORT:-8080}"
echo "🗄️  MySQL: ${DB_USER:-root}@${DB_HOST:-127.0.0.1}:${DB_PORT:-3306}/${DB_NAME:-shopapp}"
echo ""
go run main.go
