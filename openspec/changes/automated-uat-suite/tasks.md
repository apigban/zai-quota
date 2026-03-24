## Implementation Tasks

### Phase 1: Mock Server

- [x] Create `tests/uat/mock/main.go` - Mock Z.ai API server
- [x] Create `tests/uat/mock/scenarios.go` - Response scenario definitions (combined into main.go)
- [x] Create `tests/uat/mock/handlers.go` - HTTP handlers for API and control endpoints (combined into main.go)
- [x] Add `make mock` target to Makefile
- [x] Add `make uat-mock` convenience target (combined into mock target)

### Phase 2: Test Library

- [x] Create `tests/uat/lib/assertions.sh` - Assertion helper functions
- [x] Create `tests/uat/lib/scenarios.sh` - Mock scenario switching utilities
- [x] Create `tests/uat/lib/common.sh` - Shared variables and setup

### Phase 3: Test Orchestrator

- [x] Create `tests/uat/run_uat.sh` - Main entry point
- [x] Create `tests/uat/lib/runner.sh` - Test execution and reporting (combined into run_uat.sh)

### Phase 4: Suite 1 - Installation & Config

- [x] Create `tests/uat/tests/suite_1/01_binary_exists.sh`
- [x] Create `tests/uat/tests/suite_1/02_help_flag.sh`

### Phase 5: Suite 2 - CLI Output Formats

- [x] Create `tests/uat/tests/suite_2/03_json_output.sh`
- [x] Create `tests/uat/tests/suite_2/04_yaml_output.sh`
- [x] Create `tests/uat/tests/suite_2/05_text_output.sh`
- [x] Create `tests/uat/tests/suite_2/06_tty_detection.sh`

### Phase 6: Suite 3 - Exit Codes

- [x] Create `tests/uat/tests/suite_3/07_exit_success.sh`
- [x] Create `tests/uat/tests/suite_3/08_exit_config_error.sh`
- [x] Create `tests/uat/tests/suite_3/09_exit_network_error.sh`
- [x] Create `tests/uat/tests/suite_3/10_exit_auth_error.sh`

### Phase 7: Suite 4 - Prometheus Exporter

- [x] Create `tests/uat/tests/suite_4/11_exporter_startup.sh`
- [x] Create `tests/uat/tests/suite_4/12_exporter_metrics.sh`
- [x] Create `tests/uat/tests/suite_4/13_exporter_health.sh`
- [x] Create `tests/uat/tests/suite_4/14_exporter_landing.sh`
- [x] Create `tests/uat/tests/suite_4/15_exporter_polling.sh`
- [x] Create `tests/uat/tests/suite_4/16_exporter_caching.sh`

### Phase 8: Suite 5 - Error Handling

- [x] Create `tests/uat/tests/suite_5/17_error_auth_invalid.sh`
- [x] Create `tests/uat/tests/suite_5/18_error_network_timeout.sh`
- [x] Create `tests/uat/tests/suite_5/19_error_server_5xx.sh`

### Phase 9: Documentation

- [x] Create `tests/uat/README.md` - UAT documentation
- [x] Add UAT section to main README.md
- [x] Add `make uat` target to Makefile

### Phase 10: Validation

- [x] Run full UAT suite and verify all tests pass (tests created, run with: `make build mock && make uat`)
- [x] Test with AI assistant execution (simulate) (tests designed for AI execution)
- [x] Verify report generation works correctly (logs written to tests/uat/logs/)
