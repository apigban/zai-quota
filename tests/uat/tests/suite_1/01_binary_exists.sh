#!/bin/bash
# Test: Binary exists and is executable

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
source "${SCRIPT_DIR}/../../lib/assertions.sh"
source "${SCRIPT_DIR}/../../lib/common.sh"

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

test_summary
