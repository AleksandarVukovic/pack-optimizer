#!/usr/bin/env bash

set -euo pipefail

BINARY=${BINARY:-./bin/pack-optimizer}
API_URL=${API_URL:-http://localhost:8080}
MAX_RETRIES=${MAX_RETRIES:-3}
RETRY_INTERVAL=${RETRY_INTERVAL:-2}
SERVER_PID=""

log()  { echo "[integration-test] $*"; }
error() { echo "[integration-test] ERROR: $*" >&2; }

start_server() {
  log "Starting server: $BINARY"
  $BINARY &
  SERVER_PID=$!
  log "Server PID: $SERVER_PID"
}

wait_for_server() {
  log "Waiting for server to be ready at $API_URL..."
  for i in $(seq 1 "$MAX_RETRIES"); do
    if curl -sf "$API_URL/health" > /dev/null; then
      log "Server is up"
      return 0
    fi
    log "Attempt $i/$MAX_RETRIES failed, retrying in ${RETRY_INTERVAL}s..."
    sleep "$RETRY_INTERVAL"
  done
  error "Server did not start in time"
  return 1
}

stop_server() {
  if [ -n "$SERVER_PID" ] && kill -0 "$SERVER_PID" 2>/dev/null; then
    log "Stopping server (PID: $SERVER_PID)"
    kill "$SERVER_PID"
    wait "$SERVER_PID" 2>/dev/null || true
    log "Server stopped"
  fi
}

run_tests() {
  log "Running integration tests against $API_URL"
  API_URL="$API_URL" go test -v ./test/integration/...
}

# always stop the server on exit
trap stop_server EXIT

start_server
wait_for_server
run_tests