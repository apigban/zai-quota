#!/bin/bash
# Test: Prometheus metrics format

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
source "${SCRIPT_DIR}/../../lib/assertions.sh"
source "${SCRIPT_DIR}/../../lib/common.sh"
source "${SCRIPT_DIR}/../../lib/scenarios.sh"

test_header "Prometheus metrics format"

set_mock_scenario "success_full"

port=19091
"$BINARY" exporter --endpoint="$TEST_ENDPOINT" --api-key="$TEST_API_KEY" --listen=":$port" --poll-interval=60 > /dev/null 2>&1 &
exporter_pid=$!

trap "kill $exporter_pid 2>/dev/null" EXIT

sleep 2

metrics=$(curl -s "http://localhost:$port/metrics")

assert_contains "$metrics" "# TYPE zai_quota_prompt_usage_ratio gauge" "Should have prompt usage ratio metric"
assert_contains "$metrics" "# TYPE zai_quota_prompt_reset_timestamp_seconds gauge" "Should have prompt reset metric"
assert_contains "$metrics" "# TYPE zai_quota_tool_calls_used gauge" "Should have tool calls used metric"
assert_contains "$metrics" "# TYPE zai_quota_tool_calls_limit gauge" "Should have tool calls limit metric"
assert_contains "$metrics" "# TYPE zai_quota_tool_calls_remaining gauge" "Should have tool calls remaining metric"
assert_contains "$metrics" "# TYPE zai_quota_info gauge" "Should have info metric"
assert_contains "$metrics" "# TYPE zai_quota_up gauge" "Should have up metric"

assert_metric_value "$metrics" "zai_quota_prompt_usage_ratio" "0.28"
assert_metric_value "$metrics" "zai_quota_tool_calls_used" "16"
assert_metric_value "$metrics" "zai_quota_tool_calls_limit" "1000"
assert_metric_value "$metrics" "zai_quota_tool_calls_remaining" "984"
assert_metric_value "$metrics" "zai_quota_up" "1"

assert_metric_label "$metrics" "zai_quota_info" "level" "pro"

kill $exporter_pid 2>/dev/null

test_summary
