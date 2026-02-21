#!/usr/bin/env bash
set -euo pipefail

# MoCaCo Backend Start Script

# Determine runtime directory
RUNTIME_DIR="${XDG_RUNTIME_DIR:-/tmp}"
PID_FILE="${RUNTIME_DIR}/macaco.pid"
LOG_FILE="${RUNTIME_DIR}/macaco.log"
PORT="${1:-8080}"

# Check if already running
if [ -f "$PID_FILE" ]; then
    PID=$(cat "$PID_FILE")
    if kill -0 "$PID" 2>/dev/null; then
        echo "MoCaCo backend is already running (PID: $PID)"
        exit 0
    else
        echo "Removing stale PID file"
        rm -f "$PID_FILE"
    fi
fi

# Find the macaco binary
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_DIR="$(dirname "$SCRIPT_DIR")"

MACACO_BIN=""
if [ -x "${PROJECT_DIR}/bin/macaco" ]; then
    MACACO_BIN="${PROJECT_DIR}/bin/macaco"
elif [ -x "${PROJECT_DIR}/result/bin/macaco" ]; then
    MACACO_BIN="${PROJECT_DIR}/result/bin/macaco"
elif command -v macaco &>/dev/null; then
    MACACO_BIN="$(command -v macaco)"
else
    echo "Error: Could not find macaco binary"
    echo "Please run 'make build' or 'nix build' first"
    exit 1
fi

echo "Starting MoCaCo backend on port $PORT..."
echo "Binary: $MACACO_BIN"
echo "Log file: $LOG_FILE"

# Start the server in the background
nohup "$MACACO_BIN" --server --addr "localhost:$PORT" > "$LOG_FILE" 2>&1 &
BACKEND_PID=$!

# Write PID file
echo "$BACKEND_PID" > "$PID_FILE"

# Wait a moment and check if it started
sleep 1

if kill -0 "$BACKEND_PID" 2>/dev/null; then
    echo "MoCaCo backend started successfully (PID: $BACKEND_PID)"
    echo "API available at: http://localhost:$PORT/api/v1"
else
    echo "Error: Backend failed to start"
    echo "Check the log file: $LOG_FILE"
    rm -f "$PID_FILE"
    exit 1
fi
