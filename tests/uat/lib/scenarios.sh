#!/bin/bash
# Mock scenario control utilities

MOCK_PORT="${MOCK_PORT:-19876}"
MOCK_URL="http://localhost:${MOCK_PORT}"
MOCK_CONTROL_URL="${MOCK_URL}/control"
MOCK_PID_FILE="/tmp/zai-quota-mock-${MOCK_PORT}.pid"

set_mock_scenario() {
    local scenario="$1"
    local response
    response=$(curl -s "${MOCK_CONTROL_URL}/scenario/${scenario}" 2>&1)
    if echo "$response" | grep -q "OK"; then
        return 0
    else
        echo "Failed to set scenario: $response" >&2
        return 1
    fi
}

get_mock_status() {
    curl -s "${MOCK_CONTROL_URL}/status" 2>&1
}

get_mock_scenarios() {
    curl -s "${MOCK_CONTROL_URL}/list" 2>&1
}

reset_mock() {
    set_mock_scenario "success_full"
}

wait_for_mock() {
    local max_attempts="${1:-30}"
    local attempt=0
    
    while [ $attempt -lt $max_attempts ]; do
        if curl -s "${MOCK_CONTROL_URL}/status" > /dev/null 2>&1; then
            return 0
        fi
        attempt=$((attempt + 1))
        sleep 0.1
    done
    
    echo "Mock server not available after ${max_attempts} attempts" >&2
    return 1
}

start_mock() {
    local port="${1:-$MOCK_PORT}"
    local mock_binary="${PROJECT_ROOT}/tests/uat/mock/mock-server"
    
    if [ ! -f "$mock_binary" ]; then
        echo "Mock server binary not found at $mock_binary" >&2
        echo "Run 'make mock' first" >&2
        return 1
    fi
    
    if [ -f "$MOCK_PID_FILE" ]; then
        local old_pid=$(cat "$MOCK_PID_FILE" 2>/dev/null)
        if kill -0 "$old_pid" 2>/dev/null; then
            echo "Mock server already running (PID: $old_pid)" >&2
            return 0
        fi
    fi
    
    "$mock_binary" --port="$port" > /dev/null 2>&1 &
    local pid=$!
    echo $pid > "$MOCK_PID_FILE"
    
    if wait_for_mock "$max_attempts"; then
        return 0
    else
        return 1
    fi
}

stop_mock() {
    if [ -f "$MOCK_PID_FILE" ]; then
        local pid=$(cat "$MOCK_PID_FILE" 2>/dev/null)
        if [ -n "$pid" ] && kill -0 "$pid" 2>/dev/null; then
            kill "$pid" 2>/dev/null
            sleep 0.5
        fi
        rm -f "$MOCK_PID_FILE"
    fi
}

is_mock_running() {
    if [ -f "$MOCK_PID_FILE" ]; then
        local pid=$(cat "$MOCK_PID_FILE" 2>/dev/null)
        if [ -n "$pid" ] && kill -0 "$pid" 2>/dev/null; then
            return 0
        fi
    fi
    return 1
}
