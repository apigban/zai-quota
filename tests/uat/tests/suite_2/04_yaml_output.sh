#!/bin/bash
# Test: YAML output format

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
source "${SCRIPT_DIR}/../../lib/assertions.sh"
source "${SCRIPT_DIR}/../../lib/common.sh"
source "${SCRIPT_DIR}/../../lib/scenarios.sh"

test_header "YAML output format"

set_mock_scenario "success_full"

output=$(ZAI_API_KEY="$TEST_API_KEY" "$BINARY" --endpoint="$TEST_ENDPOINT" --yaml 2>&1)
exit_code=$?

assert_success "$exit_code" "YAML output should succeed"
assert_yaml_valid "$output" "Output should be valid YAML"
assert_contains "$output" "level:" "Should include level field"
assert_contains "$output" "limits:" "Should include limits field"
assert_contains "$output" "type:" "Should include type field"

test_summary
