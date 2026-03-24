#!/bin/bash
# Test: Exit code 2 on network error

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
source "${SCRIPT_DIR}/../../lib/assertions.sh"
source "${SCRIPT_DIR}/../../lib/common.sh"
source "${SCRIPT_DIR}/../../lib/scenarios.sh"

test_header "Exit code 2 on network error"

set_mock_scenario "timeout"

output=$(timeout 5 "$BINARY" --endpoint="$TEST_ENDPOINT" --api-key="$TEST_API_KEY" --json 2>&1)
exit_code=$?

assert_failure "$exit_code" "Should return non-zero exit code"
assert_contains "$output" "timeout" "Error should mention timeout or network"

test_summary
