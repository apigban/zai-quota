#!/bin/bash
# Test: Help flag displays usage

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
source "${SCRIPT_DIR}/../../lib/assertions.sh"
source "${SCRIPT_DIR}/../../lib/common.sh"

test_header "Help flag works"

output=$("$BINARY" --help 2>&1)
exit_code=$?

assert_success "$exit_code" "Help should succeed"
assert_contains "$output" "zai-quota" "Help should mention program name"
assert_contains "$output" "--json" "Help should mention --json flag"
assert_contains "$output" "--yaml" "Help should mention --yaml flag"
assert_contains "$output" "--text" "Help should mention --text flag"

test_summary
