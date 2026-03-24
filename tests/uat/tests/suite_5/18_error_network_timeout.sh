#!/bin/bash
# Test: Network timeout error handling

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
source "${SCRIPT_DIR}/../../lib/assertions.sh"
source "${SCRIPT_DIR}/../../lib/common.sh"
source "${SCRIPT_DIR}/../../lib/scenarios.sh"

test_header "Network timeout error handling"

set_mock_scenario "timeout"

output=$(timeout 5 "$BINARY" --endpoint="$TEST_ENDPOINT" --api-key="$TEST_API_KEY" --json 2>&1)
exit_code=$?

assert_failure "$exit_code" "Should fail on network timeout"
assert_contains "$output" "timeout\|failed\|error" "Error should indicate timeout or failure"

test_summary
