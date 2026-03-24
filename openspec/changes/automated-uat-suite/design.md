## Architecture

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                        UAT AUTOMATION ARCHITECTURE                            │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│   ┌──────────────────────────────────────────────────────────────────────┐  │
│   │                         TEST ORCHESTRATOR                              │  │
│   │                         run_uat.sh                                    │  │
│   │                                                                        │  │
│   │   • Initializes mock server                                          │  │
│   │   • Runs test suites in order                                        │  │
│   │   • Collects results                                                 │  │
│   │   • Generates report                                                 │  │
│   └──────────────────────────────────────────────────────────────────────┘  │
│                                    │                                         │
│                                    ▼                                         │
│   ┌──────────────────────────────────────────────────────────────────────┐  │
│   │                         TEST LIBRARY                                  │  │
│   │                         lib/                                          │  │
│   │                                                                        │  │
│   │   assertions.sh    scenarios.sh    mock_control.sh                   │  │
│   │   ─────────────    ────────────    ───────────────                   │  │
│   │   • assert_eq()    • set_scenario()                                   │  │
│   │   • assert_exit()  • mock_status()                                    │  │
│   │   • assert_contains()                                                 │  │
│   └──────────────────────────────────────────────────────────────────────┘  │
│                                    │                                         │
│                                    ▼                                         │
│   ┌──────────────────────────────────────────────────────────────────────┐  │
│   │                         TEST SUITES                                   │  │
│   │                         tests/                                        │  │
│   │                                                                        │  │
│   │   Suite 1: Installation & Config                                      │  │
│   │   ├── 01_binary_exists.sh                                             │  │
│   │   └── 02_api_key_config.sh                                            │  │
│   │                                                                        │  │
│   │   Suite 2: CLI Output Formats                                         │  │
│   │   ├── 03_json_output.sh                                               │  │
│   │   ├── 04_yaml_output.sh                                               │  │
│   │   ├── 05_text_output.sh                                               │  │
│   │   └── 06_tty_detection.sh                                             │  │
│   │                                                                        │  │
│   │   Suite 3: Exit Codes                                                 │  │
│   │   ├── 07_exit_success.sh                                              │  │
│   │   ├── 08_exit_config_error.sh                                         │  │
│   │   ├── 09_exit_network_error.sh                                        │  │
│   │   └── 10_exit_auth_error.sh                                           │  │
│   │                                                                        │  │
│   │   Suite 4: Prometheus Exporter                                        │  │
│   │   ├── 11_exporter_startup.sh                                          │  │
│   │   ├── 12_exporter_metrics.sh                                          │  │
│   │   ├── 13_exporter_health.sh                                           │  │
│   │   ├── 14_exporter_landing.sh                                          │  │
│   │   ├── 15_exporter_polling.sh                                          │  │
│   │   └── 16_exporter_caching.sh                                          │  │
│   │                                                                        │  │
│   │   Suite 5: Error Handling                                             │  │
│   │   ├── 17_error_auth_invalid.sh                                        │  │
│   │   ├── 18_error_network_timeout.sh                                     │  │
│   │   └── 19_error_server_5xx.sh                                          │  │
│   │                                                                        │  │
│   └──────────────────────────────────────────────────────────────────────┘  │
│                                    │                                         │
│                                    ▼                                         │
│   ┌──────────────────────────────────────────────────────────────────────┐  │
│   │                         MOCK SERVER                                   │  │
│   │                         mock/main.go                                  │  │
│   │                                                                        │  │
│   │   Endpoints:                                                           │  │
│   │   POST /api/monitor/usage/quota/limit                                 │  │
│   │   GET  /control/scenario/{name}   <- Switch response scenario         │  │
│   │   GET  /control/status            <- Check mock status                │  │
│   │                                                                        │  │
│   │   Scenarios:                                                           │  │
│   │   • success_full          - Both limits, full details                 │  │
│   │   • success_partial       - Only one limit                            │  │
│   │   • success_empty         - No limits                                 │  │
│   │   • auth_invalid          - 401 Unauthorized                          │  │
│   │   • auth_forbidden        - 403 Forbidden                             │  │
│   │   • server_error          - 500 Internal Server Error                 │  │
│   │   • server_unavailable    - 503 Service Unavailable                   │  │
│   │   • timeout               - Delays response (simulates timeout)       │  │
│   │   • malformed             - Invalid JSON response                     │  │
│   │   • rate_limited          - 429 Too Many Requests                     │  │
│   │                                                                        │  │
│   └──────────────────────────────────────────────────────────────────────┘  │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

## Test Data

### Mock API Response: success_full

```json
{
  "success": true,
  "code": 200,
  "msg": "OK",
  "data": {
    "limits": [
      {
        "type": "TOKENS_LIMIT",
        "percentage": 28,
        "usage": 1400,
        "currentValue": 1400,
        "total": 5000,
        "remaining": 3600,
        "nextResetTime": 1773234696431,
        "usageDetails": [
          {"modelCode": "claude-3-opus", "usage": 800},
          {"modelCode": "claude-3-sonnet", "usage": 600}
        ]
      },
      {
        "type": "TIME_LIMIT",
        "percentage": 2,
        "usage": 1000,
        "currentValue": 16,
        "total": 1000,
        "remaining": 984,
        "nextResetTime": 1775186469998,
        "usageDetails": [
          {"modelCode": "search-prime", "usage": 10},
          {"modelCode": "web-reader", "usage": 4},
          {"modelCode": "ref", "usage": 2}
        ]
      }
    ],
    "level": "pro"
  }
}
```

### Mock API Response: success_partial

```json
{
  "success": true,
  "code": 200,
  "msg": "OK",
  "data": {
    "limits": [
      {
        "type": "TOKENS_LIMIT",
        "percentage": 75,
        "usage": 3750,
        "currentValue": 3750,
        "total": 5000,
        "remaining": 1250,
        "nextResetTime": 1773234696431,
        "usageDetails": []
      }
    ],
    "level": "free"
  }
}
```

### Mock API Response: auth_invalid

```json
HTTP/1.1 401 Unauthorized
{"success": false, "code": 401, "msg": "Invalid API key", "data": null}
```

### Mock API Response: server_error

```json
HTTP/1.1 500 Internal Server Error
{"success": false, "code": 500, "msg": "Internal server error", "data": null}
```

## Exit Codes

| Code | Meaning | Test Coverage |
|------|---------|---------------|
| 0 | Success | 07_exit_success.sh |
| 1 | Configuration error | 08_exit_config_error.sh |
| 2 | Network error | 09_exit_network_error.sh |
| 3 | Authentication error | 10_exit_auth_error.sh |

## Prometheus Metrics Format

Expected `/metrics` output structure:

```
# HELP zai_quota_prompt_usage_ratio Current prompt usage as ratio (0-1)
# TYPE zai_quota_prompt_usage_ratio gauge
zai_quota_prompt_usage_ratio 0.28

# HELP zai_quota_prompt_reset_timestamp_seconds Unix timestamp when prompt limit resets
# TYPE zai_quota_prompt_reset_timestamp_seconds gauge
zai_quota_prompt_reset_timestamp_seconds 1773234696

# HELP zai_quota_tool_calls_used Number of tool calls used
# TYPE zai_quota_tool_calls_used gauge
zai_quota_tool_calls_used 16

# HELP zai_quota_tool_calls_limit Maximum allowed tool calls
# TYPE zai_quota_tool_calls_limit gauge
zai_quota_tool_calls_limit 1000

# HELP zai_quota_tool_calls_remaining Remaining tool calls
# TYPE zai_quota_tool_calls_remaining gauge
zai_quota_tool_calls_remaining 984

# HELP zai_quota_tool_calls_reset_timestamp_seconds Unix timestamp when tool limit resets
# TYPE zai_quota_tool_calls_reset_timestamp_seconds gauge
zai_quota_tool_calls_reset_timestamp_seconds 1775186469

# HELP zai_quota_tool_calls_by_tool Per-tool usage breakdown
# TYPE zai_quota_tool_calls_by_tool gauge
zai_quota_tool_calls_by_tool{tool="search-prime"} 10
zai_quota_tool_calls_by_tool{tool="web-reader"} 4
zai_quota_tool_calls_by_tool{tool="ref"} 2

# HELP zai_quota_info Subscription level info
# TYPE zai_quota_info gauge
zai_quota_info{level="pro"} 1

# HELP zai_quota_up Whether last scrape succeeded
# TYPE zai_quota_up gauge
zai_quota_up 1

# HELP zai_quota_last_scrape_timestamp_seconds Timestamp of last API poll
# TYPE zai_quota_last_scrape_timestamp_seconds gauge
zai_quota_last_scrape_timestamp_seconds 1710123400

# HELP zai_quota_scrape_duration_seconds Duration of last API poll
# TYPE zai_quota_scrape_duration_seconds gauge
zai_quota_scrape_duration_seconds 0.234
```

## Test Execution Flow

```
┌─────────────────────────────────────────────────────────────────┐
│                    TEST EXECUTION ORDER                          │
├─────────────────────────────────────────────────────────────────┤
│                                                                  │
│  1. PREPARATION                                                  │
│     ├── Build binary: make build                                │
│     ├── Build mock server: make mock                            │
│     └── Start mock server on port 19876                         │
│                                                                  │
│  2. SUITE 1: Installation & Config                              │
│     ├── Verify binary exists and is executable                  │
│     └── Test API key configuration methods                      │
│                                                                  │
│  3. SUITE 2: CLI Output Formats                                 │
│     ├── Test --json output structure                            │
│     ├── Test --yaml output structure                            │
│     ├── Test --text output (no ANSI codes)                      │
│     └── Test TTY auto-detection                                 │
│                                                                  │
│  4. SUITE 3: Exit Codes                                         │
│     ├── Verify exit code 0 on success                           │
│     ├── Verify exit code 1 on config error                      │
│     ├── Verify exit code 2 on network error                     │
│     └── Verify exit code 3 on auth error                        │
│                                                                  │
│  5. SUITE 4: Prometheus Exporter                                │
│     ├── Start exporter in background                            │
│     ├── Verify /metrics format and content                      │
│     ├── Verify /health returns "OK"                             │
│     ├── Verify / (landing page) has content                     │
│     ├── Verify polling interval is respected                    │
│     └── Verify caching between polls                            │
│                                                                  │
│  6. SUITE 5: Error Handling                                     │
│     ├── Test invalid API key error handling                     │
│     ├── Test network timeout handling                           │
│     └── Test server 5xx error handling                          │
│                                                                  │
│  7. CLEANUP                                                      │
│     ├── Stop exporter (if running)                              │
│     ├── Stop mock server                                        │
│     └── Generate report                                         │
│                                                                  │
└─────────────────────────────────────────────────────────────────┘
```

## AI Execution Protocol

When an AI assistant executes this UAT suite:

1. **Environment Setup**
   - Ensure working directory is project root
   - Verify Go is installed
   - Verify bash is available

2. **Execution Command**
   ```bash
   ./tests/uat/run_uat.sh
   ```

3. **Expected Duration**
   - Suite 1-3: ~10 seconds
   - Suite 4: ~30 seconds (includes polling wait)
   - Suite 5: ~15 seconds
   - Total: ~1 minute

4. **Success Criteria**
   - All tests pass (exit code 0)
   - Report shows 0 failures
   - No ERROR lines in output

5. **Failure Handling**
   - Report shows which test failed
   - AI can read individual test logs in `tests/uat/logs/`
   - AI can re-run single test: `./tests/uat/tests/03_json_output.sh`
