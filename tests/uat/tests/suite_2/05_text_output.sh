#!/bin/bash
# Test: Text output format (no ANSI codes)

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
source "${SCRIPT_DIR}/../../lib/assertions.sh"
source "${SCRIPT_DIR}/../../lib/common.sh"
source "${SCRIPT_DIR}/../../lib/scenarios.sh"

test_header "Text output format (no ANSI)"

set_mock_scenario "success_full"

output=$(ZAI_API_KEY="$TEST_API_KEY" "$BINARY" --endpoint="$TEST_ENDPOINT" --text 2>&1)
exit_code=$?

assert_success "$exit_code" "Text output should succeed"

ansi_count=$(printf '%s' "$output" | grep -c $'\x1b' || echo "0")
assert_eq "0" "$ansi_count" "Should have zero ANSI codes"

assert_contains "$output" "5-Hour Prompt Limit" "Should show prompt limit label"
assert_contains "$output" "Tool Quota" "Should show tool quota label"
assert_contains "$output" "28%" "Should show percentage"

test_summary
