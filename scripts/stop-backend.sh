#!/usr/bin/env bash
set -euo pipefail

# MoCaCo Backend Stop Script

RUNTIME_DIR="${XDG_RUNTIME_DIR:-/tmp}"
PID_FILE="${RUNTIME_DIR}/macaco.pid"
FORCE="${1:-}"

if [ ! -f "$PID_FILE" ]; then
    echo "MoCaCo backend is not running (no PID file found)"
    exit 0
fi

PID=$(cat "$PID_FILE")

if ! kill -0 "$PID" 2>/dev/null; then
    echo "MoCaCo backend is not running (stale PID file)"
    rm -f "$PID_FILE"
    exit 0
fi

echo "Stopping MoCaCo backend (PID: $PID)..."

if [ "$FORCE" = "-f" ] || [ "$FORCE" = "--force" ]; then
    kill -9 "$PID" 2>/dev/null || true
    echo "Force killed backend"
else
    # Graceful shutdown
    kill -TERM "$PID" 2>/dev/null || true

    # Wait for graceful shutdown
    TIMEOUT=5
    while [ $TIMEOUT -gt 0 ]; do
        if ! kill -0 "$PID" 2>/dev/null; then
            echo "Backend stopped gracefully"
            rm -f "$PID_FILE"
            exit 0
        fi
        sleep 1
        TIMEOUT=$((TIMEOUT - 1))
    done

    # Force kill if still running
    if kill -0 "$PID" 2>/dev/null; then
        echo "Graceful shutdown timed out, force killing..."
        kill -9 "$PID" 2>/dev/null || true
    fi
fi

rm -f "$PID_FILE"
echo "Backend stopped"
