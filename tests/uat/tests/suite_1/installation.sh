#!/bin/bash
# Suite 1: Installation & Configuration Tests
# Tests binary execution and API key configuration methods

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
source "${SCRIPT_DIR}/../../lib/assertions.sh"
source "${SCRIPT_DIR}/../../lib/scenarios.sh"
source "${SCRIPT_DIR}/../../lib/common.sh"

test_binary_exists() {
    test_header "Binary exists and is executable"
    
    assert_file_exists "$BINARY" "Main binary should exist"
    
    if [ -x "$BINARY" ]; then
        echo -e "  ${GREEN}✓${NC} Binary is executable"
        TESTS_PASSED=$((TESTS_PASSED + 1))
        TESTS_RUN=$((TESTS_RUN + 1))
    else
        echo -e "  ${RED}✗${NC} Binary is not executable"
        TESTS_FAILED=$((TESTS_FAILED + 1))
        TESTS_RUN=$((TESTS_RUN + 1))
    fi
}

test_help_flag() {
    test_header "Help flag works"
    
    local output
    output=$("$BINARY" --help 2>&1)
    local exit_code=$?
    
    assert_success "$exit_code" "Help should succeed"
    assert_contains "$output" "zai-quota" "Help should mention program name"
    assert_contains "$output" "--json" "Help should mention --json flag"
    assert_contains "$output" "--yaml" "Help should mention --yaml flag"
}

test_missing_api_key() {
    test_header "Missing API key shows error"
    
    local output
    output=$(ZAI_API_KEY="" ZAI_QUOTA_CONFIG="" "$BINARY" --json 2>&1)
    local exit_code=$?
    
    assert_eq "1" "$exit_code" "Should return exit code 1 for config error"
    assert_contains "$output" "ZAI_API_KEY" "Error should mention env var"
    assert_contains "$output" ".zai-quota.yaml" "Error should mention config file"
}

test_env_var_config() {
    test_header "API key from environment variable"
    
    set_mock_scenario "success_full"
    
    local output
    output=$(ZAI_API_KEY="$TEST_API_KEY" "$BINARY" --endpoint="${MOCK_URL}/api/monitor/usage/quota/limit" --json 2>&1)
    local exit_code=$?
    
    assert_success "$exit_code" "Should succeed with env var"
    assert_json_valid "$output" "Output should be valid JSON"
}

test_config_file() {
    test_header "API key from config file"
    
    local temp_config="/tmp/zai-quota-uat-config-$$.yaml"
    
    cat > "$temp_config" << EOF
api_key: "${TEST_API_KEY}"
endpoint: "${MOCK_URL}/api/monitor/usage/quota/limit"
timeout_seconds: 5
EOF
    
    set_mock_scenario "success_full"
    
    local output
    output=$(ZAI_QUOTA_CONFIG="$temp_config" "$BINARY" --json 2>&1)
    local exit_code=$?
    
    assert_success "$exit_code" "Should succeed with config file"
    
    rm -f "$temp_config"
}

run_tests() {
    init_tests
    
    test_binary_exists
    test_help_flag
    test_missing_api_key
    test_env_var_config
    test_config_file
    
    test_summary
}

if [[ "${BASH_SOURCE[0]}" == "${0}" ]]; then
    run_tests
fi
