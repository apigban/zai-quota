# Z.ai Quota Monitor - User Acceptance Testing (UAT)

## Overview

This directory contains automated end-to-end tests for the Z.ai Quota Monitor. Tests use a mock API server to simulate various scenarios without requiring a real Z.ai account.

## Quick Start

```bash
# Build both binaries
make build mock

# Run all UAT tests
make uat

# Or run directly
./tests/uat/run_uat.sh
```

## Architecture

```
tests/uat/
├── run_uat.sh              # Main entry point
├── mock/
│   └── main.go             # Mock Z.ai API server
├── lib/
│   ├── assertions.sh       # Test assertion helpers
│   ├── common.sh           # Shared variables and setup
│   └── scenarios.sh        # Mock scenario control
├── tests/
│   ├── suite_1/             # Installation & Config
│   ├── suite_2/             # CLI Output Formats
│   ├── suite_3/             # Exit Codes
│   ├── suite_4/             # Prometheus Exporter
│   └── suite_5/             # Error Handling
└── logs/                    # Test execution logs
```

## Test Suites

| Suite | Focus | Tests |
|-------|-------|-------|
| 1 | Installation & Config | Binary exists, help flags, config loading |
| 2 | CLI Output Formats | JSON, YAML, text, TTY detection |
| 3 | Exit Codes | 0 (success), 1 (config), 2 (network), 3 (auth) |
| 4 | Prometheus Exporter | Startup, metrics format, health, landing, polling, caching |
| 5 | Error Handling | Auth failures, network errors, server errors |

## Mock Server Scenarios

The mock server simulates these API responses:

| Scenario | Description |
|----------|-------------|
| `success_full` | Both TOKENS_LIMIT and TIME_LIMIT with full details |
| `success_partial` | Only TOKENS_LIMIT present |
| `success_empty` | No limits returned |
| `success_warning` | 82% usage (warning level) |
| `success_critical` | 98% usage (critical level) |
| `auth_invalid` | 401 Unauthorized |
| `auth_forbidden` | 403 Forbidden |
| `server_error` | 500 Internal Server Error |
| `server_unavailable` | 503 Service Unavailable |
| `timeout` | Delayed response (simulates timeout) |
| `malformed` | Invalid JSON response |
| `rate_limited` | 429 Too Many Requests |

### Switching Scenarios

```bash
# Set scenario
curl http://localhost:19876/control/scenario/auth_invalid

# Check current scenario
curl http://localhost:19876/control/status
```

## Writing New Tests

### Test Structure

```bash
#!/bin/bash
set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
source "${SCRIPT_DIR}/../../lib/assertions.sh"
source "${SCRIPT_DIR}/../../lib/common.sh"
source "${SCRIPT_DIR}/../../lib/scenarios.sh"

test_header "My test description"

# Set mock scenario
set_mock_scenario "success_full"

# Run binary
output=$("$BINARY" --endpoint="$TEST_ENDPOINT" --api-key="$TEST_API_KEY" --json 2>&1)
exit_code=$?

# Assertions
assert_success "$exit_code" "Should succeed"
assert_json_valid "$output" "Should be valid JSON"

test_summary
```

### Available Assertions

```bash
# Exit codes
assert_success $exit_code "message"
assert_failure $exit_code "message"
assert_eq "expected" "actual" "message"

# Content
assert_contains "$output" "substring" "message"
assert_not_contains "$output" "substring" "message"

# Files
assert_file_exists "/path/to/file" "message"
assert_file_not_exists "/path/to/file" "message"

# Formats
assert_json_valid "$output" "message"
assert_yaml_valid "$output" "message"

# Prometheus metrics
assert_metric_exists "$metrics" "metric_name" "message"
assert_metric_label "$metrics" "metric" "label" "value" "message"
```

## Running Individual Tests

```bash
# Run single test
./tests/uat/tests/suite_2/03_json_output.sh

# Run suite
for test in tests/uat/tests/suite_2/*.sh; do
    bash "$test"
done
```

## Debugging

### View Test Logs

```bash
# Most recent log
cat tests/uat/logs/uat_*.log | tail -100

# All logs
ls -la tests/uat/logs/
```

### Manual Mock Server

```bash
# Start mock server on custom port
./tests/uat/mock/mock-server --port=19877

# Test endpoint
curl -X POST http://localhost:19877/api/monitor/usage/quota/limit \
  -H "Authorization: Bearer test-key" \
  -H "Content-Type: application/json"
```

## CI/CD Integration

The test suite returns proper exit codes:
- `0` - All tests passed
- `1` - One or more tests failed

Example GitHub Actions:

```yaml
- name: Run UAT
  run: |
    make build mock
    make uat
```

## Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `MOCK_PORT` | 19876 | Port for mock server |
| `BINARY` | `./zai-quota` | Path to binary |
| `TEST_API_KEY` | `test-api-key` | API key for tests |
| `TEST_ENDPOINT` | Mock URL | API endpoint |
| `LOG_DIR` | `tests/uat/logs` | Log output directory |
