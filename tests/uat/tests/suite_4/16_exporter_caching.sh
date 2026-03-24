#!/bin/bash
# Test: Exporter caching behavior

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
source "${SCRIPT_DIR}/../../lib/assertions.sh"
source "${SCRIPT_DIR}/../../lib/common.sh"
source "${SCRIPT_DIR}/../../lib/scenarios.sh"

test_header "Exporter caching behavior"

set_mock_scenario "success_full"

port=19096
poll_interval=60

"$BINARY" exporter --endpoint="$TEST_ENDPOINT" --api-key="$TEST_API_KEY" --listen=":$port" --poll-interval=$poll_interval > /dev/null 2>&1 &
pid=$!
trap "kill $pid 2>/dev/null" EXIT

sleep 2

metrics1=$(curl -s "http://localhost:$port/metrics")
ts1=$(echo "$metrics1" | grep "zai_quota_last_scrape_timestamp_seconds" | awk '{print $2}')

sleep 3

metrics2=$(curl -s "http://localhost:$port/metrics")
ts2=$(echo "$metrics2" | grep "zai_quota_last_scrape_timestamp_seconds" | awk '{print $2}')

assert_eq "$ts1" "$ts2" "Scrape timestamp should not change (cached)"

kill $pid 2>/dev/null

test_summary
