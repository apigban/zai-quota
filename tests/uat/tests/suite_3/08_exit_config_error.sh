#!/bin/bash
# Test: Exit code 1 on config error

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
source "${SCRIPT_DIR}/../../lib/assertions.sh"
source "${SCRIPT_DIR}/../../lib/common.sh"

test_header "Exit code 1 on config error"

unset ZAI_API_KEY
unset ZAI_API_ENDPOINT

output=$("$BINARY" --json 2>&1)
exit_code=$?

assert_eq "1" "$exit_code" "Should return exit code 1 for config error"
assert_contains "$output" "API key" "Error should mention API key"

test_summary
