#!/bin/bash
# Main UAT orchestrator

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "${SCRIPT_DIR}/../.." && pwd)"

source "${SCRIPT_DIR}/lib/assertions.sh"
source "${SCRIPT_DIR}/lib/scenarios.sh"
source "${SCRIPT_DIR}/lib/common.sh"

SUITE_PASSED=0
SUITE_FAILED=0
SUITE_TOTAL=0

run_test_file() {
    local test_file="$1"
    local test_name=$(basename "$test_file" .sh)
    
    echo ""
    echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
    echo -e "${BLUE}Running: $test_name${NC}"
    echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
    
    init_tests
    
    if bash "$test_file"; then
        local exit_code=$?
        if [ $exit_code -eq 0 ]; then
            SUITE_PASSED=$((SUITE_PASSED + 1))
            echo -e "${GREEN}✓ PASSED${NC}: $test_name"
        else
            SUITE_FAILED=$((SUITE_FAILED + 1))
            echo -e "${RED}✗ FAILED${NC}: $test_name"
        fi
        SUITE_TOTAL=$((SUITE_TOTAL + 1))
    else
        echo -e "${RED}✗ ERROR${NC}: $test_name (failed to execute)"
        SUITE_FAILED=$((SUITE_FAILED + 1))
        SUITE_TOTAL=$((SUITE_TOTAL + 1))
    fi
}

run_suite() {
    local suite_name="$1"
    local suite_dir="$2"
    
    echo ""
    echo -e "${BLUE}═══════════════════════════════════════════════════════════════${NC}"
    echo -e "${BLUE}SUITE: $suite_name${NC}"
    echo -e "${BLUE}═══════════════════════════════════════════════════════════════${NC}"
    
    local start_time=$(date +%s)
    
    for test_file in "$suite_dir"/*.sh; do
        if [ -f "$test_file" ]; then
            run_test_file "$test_file"
        fi
    done
    
    local end_time=$(date +%s)
    local duration=$((end_time - start_time))
    
    echo ""
    echo "Suite completed in ${duration}s"
}

print_final_report() {
    local end_time=$(date +%s)
    local total_duration=$((end_time - START_TIME))
    
    echo ""
    echo -e "${BLUE}═══════════════════════════════════════════════════════════════${NC}"
    echo -e "${BLUE}UAT FINAL REPORT${NC}"
    echo -e "${BLUE}═══════════════════════════════════════════════════════════════${NC}"
    echo ""
    echo "  Total Duration: ${total_duration}"
    echo "  Suites Run:     $SUITE_TOTAL"
    echo -e "  ${GREEN}Passed:${NC}        $SUITE_PASSED"
    echo -e "  ${RED}Failed:${NC}        $SUITE_FAILED"
    echo ""
    
    if [ $SUITE_FAILED -eq 0 ]; then
        echo -e "${GREEN}════════════════════════════════════════════${NC}"
        echo -e "${GREEN}ALL TESTS PASSED${NC}"
        echo -e "${GREEN}════════════════════════════════════════════${NC}"
        echo ""
        echo "Full log available at: $LOG_FILE"
        return 0
    else
        echo -e "${RED}══════════════════════════════════════════${NC}"
        echo -e "${RED}SOME TESTS FAILED${NC}"
        echo -e "${RED}══════════════════════════════════════════${NC}"
        echo ""
        echo "Full log available at: $LOG_FILE"
        return 1
    fi
}

check_prerequisites() {
    echo "Checking prerequisites..."
    
    if [ ! -f "$BINARY" ]; then
        echo -e "${RED}Error:${NC} Binary not found: $BINARY"
        echo "Run 'make build' first"
        return 1
    fi
    
    if ! command -v curl &>/dev/null; then
        echo -e "${RED}Error:${NC} curl not found"
        return 1
    fi
    
    local mock_binary="${SCRIPT_DIR}/mock/mock-server"
    if [ ! -f "$mock_binary" ]; then
        echo -e "${RED}Error:${NC} Mock server not found: $mock_binary"
        echo "Run 'make mock' first"
        return 1
    fi
    
    echo -e "${GREEN}Prerequisites OK${NC}"
    return 0
}

main() {
    START_TIME=$(date +%s)
    
    echo -e "${BLUE}═══════════════════════════════════════════════════════════════${NC}"
    echo -e "${BLUE}Z.ai Quota Monitor - User Acceptance Testing${NC}"
    echo -e "${BLUE}═══════════════════════════════════════════════════════════════${NC}"
    echo "Started: $(date)"
    echo ""
    
    check_prerequisites || exit 1
    
    backup_config
    
    echo "Starting mock server on port $MOCK_PORT..."
    start_mock "$MOCK_PORT" || fail "Failed to start mock server"
    
    reset_mock
    
    run_suite "Installation & Config" "${SCRIPT_DIR}/tests/suite_1"
    run_suite "CLI Output Formats" "${SCRIPT_DIR}/tests/suite_2"
    run_suite "Exit Codes" "${SCRIPT_DIR}/tests/suite_3"
    run_suite "Prometheus Exporter" "${SCRIPT_DIR}/tests/suite_4"
    run_suite "Error Handling" "${SCRIPT_DIR}/tests/suite_5"
    
    print_final_report
    local result=$?
    
    cleanup
    
    exit $result
}

if [[ "${BASH_SOURCE[0]}" == "${0}" ]]; then
    main "$@"
fi
