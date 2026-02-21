#!/usr/bin/env bash
set -euo pipefail

# MoCaCo Backend Status Script

RUNTIME_DIR="${XDG_RUNTIME_DIR:-/tmp}"
PID_FILE="${RUNTIME_DIR}/macaco.pid"
LOG_FILE="${RUNTIME_DIR}/macaco.log"
PORT="${1:-8080}"

echo "MoCaCo Backend Status"
echo "====================="
echo ""

# Check PID file
if [ ! -f "$PID_FILE" ]; then
    echo "Status: NOT RUNNING (no PID file)"
    exit 1
fi

PID=$(cat "$PID_FILE")

if ! kill -0 "$PID" 2>/dev/null; then
    echo "Status: NOT RUNNING (stale PID file)"
    rm -f "$PID_FILE"
    exit 1
fi

echo "Status: RUNNING"
echo "PID: $PID"
echo "Port: $PORT"
echo ""

# Check health endpoint
echo "Health Check:"
if command -v curl &>/dev/null; then
    if curl -s "http://localhost:$PORT/api/v1/health" 2>/dev/null | head -c 200; then
        echo ""
        echo ""
        echo "API: HEALTHY"
    else
        echo "API: UNHEALTHY (could not connect)"
    fi
else
    echo "curl not available, skipping health check"
fi

# Show recent logs
if [ -f "$LOG_FILE" ]; then
    echo ""
    echo "Recent Logs:"
    echo "------------"
    tail -n 10 "$LOG_FILE"
fi
