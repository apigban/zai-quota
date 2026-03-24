#!/bin/bash
# Test: JSON output format

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
source "${SCRIPT_DIR}/../../lib/assertions.sh"
source "${SCRIPT_DIR}/../../lib/common.sh"
source "${SCRIPT_DIR}/../../lib/scenarios.sh"

test_header "JSON output format"

set_mock_scenario "success_full"

output=$(ZAI_API_KEY="$TEST_API_KEY" "$BINARY" --endpoint="$TEST_ENDPOINT" --json 2>&1)
exit_code=$?

assert_success "$exit_code" "JSON output should succeed"
assert_json_valid "$output" "Output should be valid JSON"
assert_contains "$output" '"level"' "Should include level field"
assert_contains "$output" '"limits"' "Should include limits array"
assert_contains "$output" '"TOKENS_LIMIT"' "Should include TOKENS_LIMIT"
assert_contains "$output" '"TIME_LIMIT"' "Should include TIME_LIMIT"
assert_contains "$output" '"percentage"' "Should include percentage"

test_summary
