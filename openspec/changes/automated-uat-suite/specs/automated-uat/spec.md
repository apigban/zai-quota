## ADDED Requirements

### Requirement: UAT runner script executes all test suites

The system SHALL provide a `run_uat.sh` script that orchestrates test execution.

#### Scenario: Execute full UAT suite
- **GIVEN** the project has been built (`make build`)
- **AND** the mock server has been built (`make mock`)
- **WHEN** user runs `./tests/uat/run_uat.sh`
- **THEN** system starts mock server on port 19876
- **AND** system executes all test suites in order
- **AND** system generates report in `tests/uat/logs/`
- **AND** system exits with code 0 if all tests pass
- **AND** system exits with code 1 if any tests fail

#### Scenario: UAT suite handles missing mock server
- **GIVEN** mock server binary does not exist
- **WHEN** user runs `./tests/uat/run_uat.sh`
- **THEN** system prints error message
- **AND** system exits with code 1

---

### Requirement: Mock server simulates Z.ai API responses

The system SHALL provide a mock server that simulates the Z.ai API.

#### Scenario: Mock server returns success response
- **GIVEN** mock server is running on port 19876
- **WHEN** client requests `POST /api/monitor/usage/quota/limit` with valid auth
- **THEN** server returns HTTP 200
- **AND** response body contains valid JSON with `success: true`
- **AND** response includes `data.limits` array

#### Scenario: Mock server switches scenarios via control endpoint
- **GIVEN** mock server is running
- **WHEN** client requests `GET /control/scenario/auth_invalid`
- **THEN** server switches to auth_invalid scenario
- **AND** subsequent API requests return 401 Unauthorized

---

### Requirement: Assertion library provides test helpers

The system SHALL provide assertion helper functions in `lib/assertions.sh`.

#### Scenario: assert_eq passes for equal values
- **GIVEN** test script sources `lib/assertions.sh`
- **WHEN** `assert_eq "expected" "expected" "values match"` is called
- **THEN** function returns 0
- **AND** output shows green checkmark

#### Scenario: assert_eq fails for unequal values
- **GIVEN** test script sources `lib/assertions.sh`
- **WHEN** `assert_eq "expected" "actual" "values differ"` is called
- **THEN** function returns 1
- **AND** output shows red X with both values

#### Scenario: assert_exit verifies exit codes
- **GIVEN** test script sources `lib/assertions.sh`
- **WHEN** `assert_exit 0 0 "success"` is called
- **THEN** function returns 0

---

### Requirement: JSON output format is valid

The system SHALL output valid JSON when `--json` flag is provided.

#### Scenario: JSON output is parseable
- **GIVEN** API returns success response
- **WHEN** user runs `zai-quota --json`
- **THEN** output is valid JSON
- **AND** output contains `level` field
- **AND** output contains `limits` array

---

### Requirement: YAML output format is valid

The system SHALL output valid YAML when `--yaml` flag is provided.

#### Scenario: YAML output is parseable
- **GIVEN** API returns success response
- **WHEN** user runs `zai-quota --yaml`
- **THEN** output is valid YAML
- **AND** output contains `level:` key

---

### Requirement: Text output has no ANSI codes

The system SHALL output plain text without ANSI escape codes when `--text` flag is provided.

#### Scenario: Text output contains no escape sequences
- **GIVEN** API returns success response
- **WHEN** user runs `zai-quota --text`
- **THEN** output contains no `\x1b` characters
- **AND** output is suitable for piping to files

---

### Requirement: Exit code 0 on success

The system SHALL exit with code 0 when quota fetch succeeds.

#### Scenario: Successful quota fetch
- **GIVEN** valid API key is configured
- **AND** API returns success response
- **WHEN** user runs `zai-quota --json`
- **THEN** process exits with code 0

---

### Requirement: Exit code 1 on config error

The system SHALL exit with code 1 when API key is not configured.

#### Scenario: Missing API key
- **GIVEN** no API key is configured
- **WHEN** user runs `zai-quota --json`
- **THEN** process exits with code 1
- **AND** error message mentions configuration options

---

### Requirement: Exit code 2 on network error

The system SHALL exit with code 2 when network request fails.

#### Scenario: Network timeout
- **GIVEN** API server does not respond
- **WHEN** user runs `zai-quota --json`
- **THEN** process exits with code 2

---

### Requirement: Exit code 3 on auth error

The system SHALL exit with code 3 when authentication fails.

#### Scenario: Invalid API key
- **GIVEN** API returns 401 Unauthorized
- **WHEN** user runs `zai-quota --json`
- **THEN** process exits with code 3

---

### Requirement: Prometheus exporter starts successfully

The system SHALL start the Prometheus exporter when `exporter` subcommand is used.

#### Scenario: Exporter startup with defaults
- **GIVEN** valid API key is configured
- **WHEN** user runs `zai-quota exporter`
- **THEN** HTTP server starts on port 9090
- **AND** server polls API every 60 seconds

---

### Requirement: Metrics endpoint returns Prometheus format

The system SHALL expose `/metrics` endpoint returning Prometheus text format.

#### Scenario: Prometheus scrapes metrics
- **GIVEN** exporter is running
- **WHEN** client requests `GET /metrics`
- **THEN** response has content-type `text/plain; version=0.0.4`
- **AND** response body contains `zai_quota_` prefixed metrics

---

### Requirement: Health endpoint returns OK

The system SHALL expose `/health` endpoint returning "OK".

#### Scenario: Health check request
- **GIVEN** exporter is running
- **WHEN** client requests `GET /health`
- **THEN** response status is 200
- **AND** response body is "OK"

---

### Requirement: Exporter caches metrics between polls

The system SHALL serve cached metrics without triggering new API calls.

#### Scenario: Scrape between polls returns cached data
- **GIVEN** exporter polled API at T0
- **WHEN** client scrapes at T0+30s
- **THEN** response contains same data as T0
- **AND** no new API request is made
