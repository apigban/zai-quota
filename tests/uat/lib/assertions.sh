#!/bin/bash
# Assertion library for UAT tests

set -e

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

TESTS_RUN=0
TESTS_PASSED=0
TESTS_FAILED=0

init_tests() {
    TESTS_RUN=0
    TESTS_PASSED=0
    TESTS_FAILED=0
}

test_header() {
    local test_name="$1"
    echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
    echo -e "${BLUE}TEST: $test_name${NC}"
    echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
}

test_summary() {
    echo ""
    echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
    echo -e "${BLUE}TEST SUMMARY${NC}"
    echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
    echo -e "  Tests Run:    $TESTS_RUN"
    echo -e "  ${GREEN}Passed:${NC}       $TESTS_PASSED"
    echo -e "  ${RED}Failed:${NC}       $TESTS_FAILED"
    echo ""
    
    if [ $TESTS_FAILED -eq 0 ]; then
        echo -e "${GREEN}✓ ALL TESTS PASSED${NC}"
        return 0
    else
        echo -e "${RED}✗ SOME TESTS FAILED${NC}"
        return 1
    fi
}

assert_eq() {
    local expected="$1"
    local actual="$2"
    local message="${3:-Values should be equal}"
    
    TESTS_RUN=$((TESTS_RUN + 1))
    
    if [ "$expected" == "$actual" ]; then
        echo -e "  ${GREEN}✓${NC} $message"
        TESTS_PASSED=$((TESTS_PASSED + 1))
        return 0
    else
        echo -e "  ${RED}✗${NC} $message"
        echo -e "    ${YELLOW}Expected:${NC} $expected"
        echo -e "    ${YELLOW}Actual:${NC}   $actual"
        TESTS_FAILED=$((TESTS_FAILED + 1))
        return 1
    fi
}

assert_contains() {
    local haystack="$1"
    local needle="$2"
    local message="${3:-String should contain substring}"
    
    TESTS_RUN=$((TESTS_RUN + 1))
    
    if [[ "$haystack" == *"$needle"* ]]; then
        echo -e "  ${GREEN}✓${NC} $message"
        TESTS_PASSED=$((TESTS_PASSED + 1))
        return 0
    else
        echo -e "  ${RED}✗${NC} $message"
        echo -e "    ${YELLOW}Substring not found:${NC} $needle"
        TESTS_FAILED=$((TESTS_FAILED + 1))
        return 1
    fi
}

assert_not_contains() {
    local haystack="$1"
    local needle="$2"
    local message="${3:-String should not contain substring}"
    
    TESTS_RUN=$((TESTS_RUN + 1))
    
    if [[ "$haystack" != *"$needle"* ]]; then
        echo -e "  ${GREEN}✓${NC} $message"
        TESTS_PASSED=$((TESTS_PASSED + 1))
        return 0
    else
        echo -e "  ${RED}✗${NC} $message"
        echo -e "    ${YELLOW}Unexpected substring found:${NC} $needle"
        TESTS_FAILED=$((TESTS_FAILED + 1))
        return 1
    fi
}

assert_success() {
    local exit_code="$1"
    local message="${2:-Command should succeed}"
    
    assert_eq "0" "$exit_code" "$message"
}

assert_failure() {
    local exit_code="$1"
    local message="${2:-Command should fail}"
    
    TESTS_RUN=$((TESTS_RUN + 1))
    
    if [ "$exit_code" -ne 0 ]; then
        echo -e "  ${GREEN}✓${NC} $message (exit code: $exit_code)"
        TESTS_PASSED=$((TESTS_PASSED + 1))
        return 0
    else
        echo -e "  ${RED}✗${NC} $message (expected non-zero, got 0)"
        TESTS_FAILED=$((TESTS_FAILED + 1))
        return 1
    fi
}

assert_file_exists() {
    local file="$1"
    local message="${2:-File should exist}"
    
    TESTS_RUN=$((TESTS_RUN + 1))
    
    if [ -f "$file" ]; then
        echo -e "  ${GREEN}✓${NC} $message"
        TESTS_PASSED=$((TESTS_PASSED + 1))
        return 0
    else
        echo -e "  ${RED}✗${NC} $message"
        echo -e "    ${YELLOW}File not found:${NC} $file"
        TESTS_FAILED=$((TESTS_FAILED + 1))
        return 1
    fi
}

assert_json_valid() {
    local json="$1"
    local message="${2:-JSON should be valid}"
    
    TESTS_RUN=$((TESTS_RUN + 1))
    
    if echo "$json" | python3 -c "import json,sys; json.load(sys.stdin)" 2>/dev/null; then
        echo -e "  ${GREEN}✓${NC} $message"
        TESTS_PASSED=$((TESTS_PASSED + 1))
        return 0
    else
        echo -e "  ${RED}✗${NC} $message"
        TESTS_FAILED=$((TESTS_FAILED + 1))
        return 1
    fi
}

assert_yaml_valid() {
    local yaml="$1"
    local message="${2:-YAML should be valid}"
    
    TESTS_RUN=$((TESTS_RUN + 1))
    
    if echo "$yaml" | python3 -c "import yaml,sys; yaml.safe_load(sys.stdin)" 2>/dev/null; then
        echo -e "  ${GREEN}✓${NC} $message"
        TESTS_PASSED=$((TESTS_PASSED + 1))
        return 0
    else
        echo -e "  ${RED}✗${NC} $message"
        TESTS_FAILED=$((TESTS_FAILED + 1))
        return 1
    fi
}

assert_http_status() {
    local expected="$1"
    local actual="$2"
    local message="${3:-HTTP status should be correct}"
    
    assert_eq "$expected" "$actual" "$message"
}

assert_metric_exists() {
    local metrics="$1"
    local metric_name="$2"
    local message="${3:-Metric should exist}"
    
    assert_contains "$metrics" "$metric_name" "$message"
}

assert_metric_label() {
    local metrics="$1"
    local metric_name="$2"
    local label_name="$3"
    local label_value="$4"
    local message="${5:-Metric should have correct label}"
    
    TESTS_RUN=$((TESTS_RUN + 1))
    
    if echo "$metrics" | grep -q "${metric_name}{.*${label_name}=\"${label_value}\".*}"; then
        echo -e "  ${GREEN}✓${NC} $message"
        TESTS_PASSED=$((TESTS_PASSED + 1))
        return 0
    else
        echo -e "  ${RED}✗${NC} $message"
        echo -e "    ${YELLOW}Expected label:${NC} ${label_name}=\"${label_value}\""
        TESTS_FAILED=$((TESTS_FAILED + 1))
        return 1
    fi
}

skip_test() {
    local reason="$1"
    echo -e "  ${YELLOW}⊘ SKIPPED:${NC} $reason"
}

fail() {
    local message="$1"
    echo -e "  ${RED}✗ FATAL:${NC} $message"
    exit 1
}
