#!/bin/bash
# Test: Exit code 3 on auth error

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
source "${SCRIPT_DIR}/../../lib/assertions.sh"
source "${SCRIPT_DIR}/../../lib/common.sh"
source "${SCRIPT_DIR}/../../lib/scenarios.sh"

test_header "Exit code 3 on auth error"

set_mock_scenario "auth_invalid"

output=$("$BINARY" --endpoint="$TEST_ENDPOINT" --api-key="$TEST_API_KEY" --json 2>&1)
exit_code=$?

assert_eq "3" "$exit_code" "Should return exit code 3 on auth error"
assert_contains "$output" "401\|unauthorized\|invalid" "Error should indicate auth failure"

test_summary
