#!/bin/bash
# Test: Server 5xx error handling

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
source "${SCRIPT_DIR}/../../lib/assertions.sh"
source "${SCRIPT_DIR}/../../lib/common.sh"
source "${SCRIPT_DIR}/../../lib/scenarios.sh"

test_header "Server 5xx error handling"

set_mock_scenario "server_error"

output=$("$BINARY" --endpoint="$TEST_ENDPOINT" --api-key="$TEST_API_KEY" --json 2>&1)
exit_code=$?

assert_failure "$exit_code" "Should fail on server error"
assert_contains "$output" "500\|server\|error\|unavailable" "Error should indicate server error"

test_summary
