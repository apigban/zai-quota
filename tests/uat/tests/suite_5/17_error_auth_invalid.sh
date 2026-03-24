#!/bin/bash
# Test: Auth invalid error handling

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
source "${SCRIPT_DIR}/../../lib/assertions.sh"
source "${SCRIPT_DIR}/../../lib/common.sh"
source "${SCRIPT_DIR}/../../lib/scenarios.sh"

test_header "Auth invalid error handling"

set_mock_scenario "auth_invalid"

output=$(ZAI_API_KEY="$TEST_API_KEY" "$BINARY" --endpoint="$TEST_ENDPOINT" --json 2>&1)
exit_code=$?

assert_failure "$exit_code" "Should fail with invalid auth"
assert_contains "$output" "unauthorized\|401\|invalid" "Error should indicate auth failure"

test_summary
