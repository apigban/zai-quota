#!/bin/bash
# Test: Prometheus exporter starts successfully

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
source "${SCRIPT_DIR}/../../lib/assertions.sh"
source "${SCRIPT_DIR}/../../lib/common.sh"
source "${SCRIPT_DIR}/../../lib/scenarios.sh"

test_header "Prometheus exporter startup"

set_mock_scenario "success_full"

port=19090
exporter_pid=""

start_exporter() {
    "$BINARY" exporter --endpoint="$TEST_ENDPOINT" --api-key="$TEST_API_KEY" --listen=":$port" --poll-interval=60 > /dev/null 2>&1 &
    exporter_pid=$!
    sleep 2
}

stop_exporter() {
    if [ -n "$exporter_pid" ]; then
        kill "$exporter_pid" 2>/dev/null || true
    fi
}

trap stop_exporter EXIT

start_exporter

assert_process_running "$exporter_pid" "Exporter process should be running"

health=$(curl -s "http://localhost:$port/health")
assert_eq "OK" "$health" "Health endpoint should return OK"

stop_exporter

test_summary
