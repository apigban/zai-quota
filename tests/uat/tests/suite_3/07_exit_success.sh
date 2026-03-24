#!/bin/bash
# Test: Exit code 0 on success

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
source "${SCRIPT_DIR}/../../lib/assertions.sh"
source "${SCRIPT_DIR}/../../lib/common.sh"
source "${SCRIPT_DIR}/../../lib/scenarios.sh"

test_header "Exit code 0 on success"

set_mock_scenario "success_full"

ZAI_API_KEY="$TEST_API_KEY" "$BINARY" --endpoint="$TEST_ENDPOINT" --json >/dev/null 2>&1
exit_code=$?

assert_eq "0" "$exit_code" "Should return exit code 0 on success"

test_summary
