#!/bin/bash
# Test: Health endpoint returns OK

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
source "${SCRIPT_DIR}/../../lib/assertions.sh"
source "${SCRIPT_DIR}/../../lib/common.sh"

test_header "Health endpoint"

port=19092
"$BINARY" exporter --endpoint="$TEST_ENDPOINT" --api-key="$TEST_API_KEY" --listen=":$port" --poll-interval=60 > /dev/null 2>&1 &
pid=$!
trap "kill $pid 2>/dev/null" EXIT
sleep 2

response=$(curl -s "http://localhost:$port/health")
exit_code=$?

assert_success "$exit_code" "Health check request should succeed"
assert_eq "OK" "$response" "Health endpoint should return OK"

test_summary
