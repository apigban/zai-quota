#!/bin/bash
# Common UAT setup and configuration

set -e

UAT_LIB_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
UAT_ROOT="$(cd "${UAT_LIB_DIR}/.." && pwd)"
PROJECT_ROOT="$(cd "${UAT_ROOT}/.." && pwd)"

BINARY="${BINARY:-${PROJECT_ROOT}/zai-quota}"
export PROJECT_ROOT
MOCK_PORT="${MOCK_PORT:-19876}"
MOCK_URL="http://localhost:${MOCK_PORT}"
CONFIG_FILE="${HOME}/.zai-quota.yaml"
TEST_API_KEY="test-api-key-for-uat"
TEST_ENDPOINT="${MOCK_URL}/api/monitor/usage/quota/limit"

LOG_DIR="${PROJECT_ROOT}/tests/uat/logs"
LOG_FILE="${LOG_DIR}/uat_$(date +%Y%m%d_%H%M%S).log"

mkdir -p "$LOG_DIR"

ensure_binary() {
    if [ ! -f "$BINARY" ]; then
        echo "Error: Binary not found at $BINARY"
        echo "Run 'make build' first"
        return 1
    fi
    
    if [ ! -x "$BINARY" ]; then
        chmod +x "$BINARY"
    fi
    
    return 0
}

ensure_mock_server() {
    local mock_binary="${PROJECT_ROOT}/tests/uat/mock/mock-server"
    
    if [ ! -f "$mock_binary" ]; then
        echo "Error: Mock server not found at $mock_binary"
        echo "Run 'make mock' first"
        return 1
    fi
    
    if [ ! -x "$mock_binary" ]; then
        chmod +x "$mock_binary"
    fi
    
    return 0
}

backup_config() {
    if [ -f "$CONFIG_FILE" ]; then
        cp "$CONFIG_FILE" "${CONFIG_FILE}.uat-backup"
    fi
}

restore_config() {
    if [ -f "${CONFIG_FILE}.uat-backup" ]; then
        mv "${CONFIG_FILE}.uat-backup" "$CONFIG_FILE"
    else
        rm -f "$CONFIG_FILE"
    fi
}

clear_config() {
    rm -f "$CONFIG_FILE"
}

set_config() {
    local api_key="$1"
    local endpoint="${2:-$TEST_ENDPOINT}"
    cat > "$CONFIG_FILE" << EOF
api_key: "${api_key}"
endpoint: "${endpoint}"
timeout_seconds: 5
EOF
}

clear_env() {
    unset ZAI_API_KEY
    unset ZAI_QUOTA_CONFIG
}

set_env_key() {
    export ZAI_API_KEY="$1"
}

run_with_mock() {
    local args="$1"
    ZAI_API_KEY="$TEST_API_KEY" "$BINARY" --endpoint="$TEST_ENDPOINT" $args
}

run_with_mock_and_timeout() {
    local timeout_seconds="$1"
    local args="$2"
    timeout "$timeout_seconds" "$BINARY" --endpoint="$TEST_ENDPOINT" --api-key="$TEST_API_KEY" $args
}

cleanup() {
    stop_mock 2>/dev/null || true
    
    if [ -n "$EXPORTER_PID" ]; then
        kill $EXPORTER_PID 2>/dev/null || true
    fi
    
    restore_config
}

trap cleanup EXIT

log() {
    local message="$1"
    local timestamp=$(date '+%Y-%m-%d %H:%M:%S')
    echo "[$timestamp] $message" >> "$LOG_FILE"
}
