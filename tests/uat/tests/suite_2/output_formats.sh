#!/bin/bash
# Suite 2: CLI Output Format Tests
# Tests JSON, YAML, and text output formats

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
source "${SCRIPT_DIR}/../../lib/assertions.sh"
source "${SCRIPT_DIR}/../../lib/scenarios.sh"
source "${SCRIPT_DIR}/../../lib/common.sh"

test_json_output_format() {
    test_header "JSON output format"
    
    set_mock_scenario "success_full"
    
    local output
    output=$(ZAI_API_KEY="$TEST_API_KEY" "$BINARY" --endpoint="${MOCK_URL}/api/monitor/usage/quota/limit" --json 2>&1)
    local exit_code=$?
    
    assert_success "$exit_code" "JSON output should succeed"
    assert_json_valid "$output" "Output should be valid JSON"
    assert_contains "$output" '"level"' "Should include level field"
    assert_contains "$output" '"limits"' "Should include limits array"
    assert_contains "$output" '"TOKENS_LIMIT"' "Should include TOKENS_LIMIT type"
    assert_contains "$output" '"TIME_LIMIT"' "Should include TIME_LIMIT type"
    assert_contains "$output" '"percentage"' "Should include percentage field"
    assert_contains "$output" '"usageDetails"' "Should include usageDetails"
}

test_yaml_output_format() {
    test_header "YAML output format"
    
    set_mock_scenario "success_full"
    
    local output
    output=$(ZAI_API_KEY="$TEST_API_KEY" "$BINARY" --endpoint="${MOCK_URL}/api/monitor/usage/quota/limit" --yaml 2>&1)
    local exit_code=$?
    
    assert_success "$exit_code" "YAML output should succeed"
    assert_yaml_valid "$output" "Output should be valid YAML"
    assert_contains "$output" "level:" "Should include level field"
    assert_contains "$output" "limits:" "Should include limits field"
}

test_text_output_format() {
    test_header "Plain text output format"
    
    set_mock_scenario "success_full"
    
    local output
    output=$(ZAI_API_KEY="$TEST_API_KEY" "$BINARY" --endpoint="${MOCK_URL}/api/monitor/usage/quota/limit" --text 2>&1)
    local exit_code=$?
    
    assert_success "$exit_code" "Text output should succeed"
    
    local ansi_count=$(printf '%s' "$output" | grep -c $'\x1b' || echo "0")
    assert_eq "0" "$ansi_count" "Should have zero ANSI codes"
    
    assert_contains "$output" "5-Hour Prompt Limit" "Should show prompt limit label"
    assert_contains "$output" "Tool Quota" "Should show tool quota label"
}

test_mutually_exclusive_flags() {
    test_header "Mutually exclusive format flags"
    
    local output
    output=$(ZAI_API_KEY="$TEST_API_KEY" "$BINARY" --json --yaml 2>&1)
    local exit_code=$?
    
    assert_failure "$exit_code" "Should fail with multiple format flags"
    assert_contains "$output" "mutually exclusive" "Should mention mutual exclusion"
}

test_tty_detection() {
    test_header "TTY auto-detection"
    
    set_mock_scenario "success_full"
    
    local output
    output=$(echo "" | ZAI_API_KEY="$TEST_API_KEY" "$BINARY" --endpoint="${MOCK_URL}/api/monitor/usage/quota/limit" 2>&1)
    local exit_code=$?
    
    assert_success "$exit_code" "Piped output should succeed"
    
    local ansi_count=$(printf '%s' "$output" | grep -c $'\x1b' || echo "0")
    assert_eq "0" "$ansi_count" "Piped output should have no ANSI"
}

test_json_partial_data() {
    test_header "JSON with partial data"
    
    set_mock_scenario "success_partial"
    
    local output
    output=$(ZAI_API_KEY="$TEST_API_KEY" "$BINARY" --endpoint="${MOCK_URL}/api/monitor/usage/quota/limit" --json 2>&1)
    local exit_code=$?
    
    assert_success "$exit_code" "Should handle partial data"
    assert_json_valid "$output" "Output should be valid JSON"
    assert_contains "$output" "TOKENS_LIMIT" "Should include TOKENS_LIMIT"
}

run_tests() {
    init_tests
    
    test_json_output_format
    test_yaml_output_format
    test_text_output_format
    test_mutually_exclusive_flags
    test_tty_detection
    test_json_partial_data
    
    test_summary
}

if [[ "${BASH_SOURCE[0]}" == "${0}" ]]; then
    run_tests
fi
