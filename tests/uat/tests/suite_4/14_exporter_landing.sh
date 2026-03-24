#!/bin/bash
# Test: Landing page at root path

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
source "${SCRIPT_DIR}/../../lib/assertions.sh"
source "${SCRIPT_DIR}/../../lib/common.sh"
source "${SCRIPT_DIR}/../../lib/scenarios.sh"

test_header "Landing page at root path"

set_mock_scenario "success_full"

port=19093
"$BINARY" exporter --endpoint="$TEST_ENDPOINT" --api-key="$TEST_API_KEY" --listen=":$port" --poll-interval=60 > /dev/null 2>&1 &
pid=$!
trap "kill $pid 2>/dev/null" EXIT
sleep 2

response=$(curl -s "http://localhost:$port/")
exit_code=$?

assert_success "$exit_code" "Landing page request should succeed"
assert_contains "$response" "Z.ai" "Landing page should mention Z.ai"
assert_contains "$response" "/metrics" "Landing page should have link to /metrics"
assert_contains "$response" "Prometheus" "Landing page should mention Prometheus"

kill $pid 2>/dev/null

test_summary
