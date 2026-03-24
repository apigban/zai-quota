#!/bin/bash
# Test: TTY auto-detection

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
source "${SCRIPT_DIR}/../../lib/assertions.sh"
source "${SCRIPT_DIR}/../../lib/common.sh"
source "${SCRIPT_DIR}/../../lib/scenarios.sh"

test_header "TTY auto-detection"

set_mock_scenario "success_full"

output=$(echo "" | ZAI_API_KEY="$TEST_API_KEY" "$BINARY" --endpoint="$TEST_ENDPOINT" 2>&1)
exit_code=$?

assert_success "$exit_code" "Piped output should succeed"

ansi_count=$(printf '%s' "$output" | grep -c $'\x1b' || echo "0")
assert_eq "0" "$ansi_count" "Piped output should have no ANSI codes"

test_summary
