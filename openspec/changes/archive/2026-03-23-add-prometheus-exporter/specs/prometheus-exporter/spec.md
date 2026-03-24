## ADDED Requirements

### Requirement: Exporter subcommand starts metrics server

The system SHALL provide an `exporter` subcommand that starts an HTTP server exposing Prometheus-compatible metrics.

#### Scenario: Start exporter with defaults
- **WHEN** user runs `zai-quota exporter`
- **THEN** system starts HTTP server on default port `:9090`
- **AND** server polls Z.ai API every 60 seconds

#### Scenario: Start exporter with custom configuration
- **WHEN** user runs `zai-quota exporter --poll-interval=120 --listen=:8080`
- **THEN** system starts HTTP server on port `:8080`
- **AND** server polls Z.ai API every 120 seconds

#### Scenario: Reject invalid poll interval
- **WHEN** user runs `zai-quota exporter --poll-interval=30`
- **THEN** system exits with error message indicating minimum interval is 60 seconds

---

### Requirement: Metrics endpoint exposes Prometheus format

The system SHALL expose a `/metrics` endpoint returning metrics in Prometheus text exposition format.

#### Scenario: Prometheus scrapes metrics
- **WHEN** Prometheus performs HTTP GET on `/metrics`
- **THEN** system returns HTTP 200 with `text/plain; version=0.0.4` content type
- **AND** response body contains valid Prometheus metric format

---

### Requirement: Prompt usage metrics exposed

The system SHALL expose metrics for prompt usage (TOKENS_LIMIT) as gauges.

#### Scenario: Prompt usage ratio available
- **WHEN** Z.ai API reports TOKENS_LIMIT with `percentage=28`
- **THEN** `/metrics` exposes `zai_quota_prompt_usage_ratio 0.28`

#### Scenario: Prompt reset timestamp available
- **WHEN** Z.ai API reports TOKENS_LIMIT with `nextResetTime=1773234696431`
- **THEN** `/metrics` exposes `zai_quota_prompt_reset_timestamp_seconds 1773234696`

---

### Requirement: Tool call metrics exposed

The system SHALL expose metrics for tool call quota (TIME_LIMIT) as gauges.

#### Scenario: Tool call usage available
- **WHEN** Z.ai API reports TIME_LIMIT with `currentValue=16`
- **THEN** `/metrics` exposes `zai_quota_tool_calls_used 16`

#### Scenario: Tool call limit available
- **WHEN** Z.ai API reports TIME_LIMIT with `usage=1000`
- **THEN** `/metrics` exposes `zai_quota_tool_calls_limit 1000`

#### Scenario: Tool call remaining available
- **WHEN** Z.ai API reports TIME_LIMIT with `remaining=984`
- **THEN** `/metrics` exposes `zai_quota_tool_calls_remaining 984`

#### Scenario: Tool call reset timestamp available
- **WHEN** Z.ai API reports TIME_LIMIT with `nextResetTime=1775186469998`
- **THEN** `/metrics` exposes `zai_quota_tool_calls_reset_timestamp_seconds 1775186469`

---

### Requirement: Per-tool usage breakdown exposed

The system SHALL expose per-tool usage metrics from usageDetails.

#### Scenario: Multiple tools reported
- **WHEN** Z.ai API reports usageDetails with `[{"modelCode":"search-prime","usage":4},{"modelCode":"web-reader","usage":3}]`
- **THEN** `/metrics` exposes:
  - `zai_quota_tool_calls_by_tool{tool="search-prime"} 4`
  - `zai_quota_tool_calls_by_tool{tool="web-reader"} 3`

#### Scenario: No tools used
- **WHEN** Z.ai API reports empty or missing usageDetails
- **THEN** no `zai_quota_tool_calls_by_tool` metrics are exposed

---

### Requirement: Subscription level exposed as info metric

The system SHALL expose subscription level as an info metric.

#### Scenario: Pro subscription
- **WHEN** Z.ai API reports `level="pro"`
- **THEN** `/metrics` exposes `zai_quota_info{level="pro"} 1`

---

### Requirement: Exporter health metrics exposed

The system SHALL expose health metrics for the exporter itself.

#### Scenario: Successful API poll
- **WHEN** last API poll succeeded
- **THEN** `/metrics` exposes `zai_quota_up 1`

#### Scenario: Failed API poll
- **WHEN** last API poll failed
- **THEN** `/metrics` exposes `zai_quota_up 0`

#### Scenario: Last scrape timestamp
- **WHEN** API poll completes at Unix timestamp 1710123400
- **THEN** `/metrics` exposes `zai_quota_last_scrape_timestamp_seconds 1710123400`

#### Scenario: Scrape duration tracked
- **WHEN** API poll takes 0.234 seconds
- **THEN** `/metrics` exposes `zai_quota_scrape_duration_seconds 0.234`

---

### Requirement: Polling respects minimum interval

The system SHALL enforce a minimum 60-second interval between API requests.

#### Scenario: Poll at configured interval
- **WHEN** poll interval is set to 120 seconds
- **THEN** system makes API request exactly once every 120 seconds

#### Scenario: Consecutive polls respect minimum
- **WHEN** exporter runs for 5 minutes with default 60s interval
- **THEN** system makes exactly 5 API requests (at 0s, 60s, 120s, 180s, 240s)

---

### Requirement: Cached metrics served on demand

The system SHALL serve cached metrics when Prometheus scrapes, without triggering new API calls.

#### Scenario: Scrape between polls
- **WHEN** Prometheus scrapes `/metrics` 30 seconds after last API poll
- **THEN** system returns cached metrics from 30 seconds ago
- **AND** no new API request is made

#### Scenario: Metrics available immediately
- **WHEN** exporter first starts
- **THEN** system makes initial API request before serving `/metrics`
- **AND** `/metrics` returns data from initial request

---

### Requirement: Graceful handling of API failures

The system SHALL continue operating when API requests fail.

#### Scenario: API returns error
- **WHEN** Z.ai API returns 5xx error
- **THEN** exporter continues running
- **AND** `zai_quota_up` is set to 0
- **AND** previous cached metrics remain available

#### Scenario: Network timeout
- **WHEN** API request times out
- **THEN** exporter continues running
- **AND** `zai_quota_up` is set to 0
- **AND** next poll attempt occurs at scheduled interval

---

### Requirement: Landing page provided

The system SHALL provide a simple HTML landing page at the root path.

#### Scenario: Access root path
- **WHEN** user navigates to `http://localhost:9090/`
- **THEN** system returns HTML page with exporter name
- **AND** page includes link to `/metrics`

---

### Requirement: Health endpoint provided

The system SHALL provide a `/health` endpoint for health checks.

#### Scenario: Health check
- **WHEN** client performs HTTP GET on `/health`
- **THEN** system returns HTTP 200 with body `OK`
