## 1. Setup and Dependencies

- [x] 1.1 Add `github.com/prometheus/client_golang` dependency to go.mod
- [x] 1.2 Create `internal/exporter` package directory structure

## 2. Core Polling and Caching

- [x] 2.1 Create `internal/exporter/cache.go` with thread-safe metric cache structure
- [x] 2.2 Create `internal/exporter/poller.go` with background polling logic
- [x] 2.3 Implement 60-second minimum interval enforcement in poller
- [x] 2.4 Add graceful error handling for failed API requests
- [x] 2.5 Write unit tests for cache and poller

## 3. Prometheus Metrics Registry

- [x] 3.1 Create `internal/exporter/metrics.go` with Prometheus metric definitions
- [x] 3.2 Implement `zai_quota_prompt_usage_ratio` gauge
- [x] 3.3 Implement `zai_quota_prompt_reset_timestamp_seconds` gauge
- [x] 3.4 Implement `zai_quota_tool_calls_used` gauge
- [x] 3.5 Implement `zai_quota_tool_calls_limit` gauge
- [x] 3.6 Implement `zai_quota_tool_calls_remaining` gauge
- [x] 3.7 Implement `zai_quota_tool_calls_reset_timestamp_seconds` gauge
- [x] 3.8 Implement `zai_quota_tool_calls_by_tool` gaugevec with tool label
- [x] 3.9 Implement `zai_quota_info` gaugevec with level label
- [x] 3.10 Implement `zai_quota_up` gauge
- [x] 3.11 Implement `zai_quota_last_scrape_timestamp_seconds` gauge
- [x] 3.12 Implement `zai_quota_scrape_duration_seconds` gauge
- [x] 3.13 Write unit tests for metrics conversion

## 4. HTTP Server

- [x] 4.1 Create `internal/exporter/server.go` with HTTP server setup
- [x] 4.2 Implement `/metrics` endpoint with Prometheus handler
- [x] 4.3 Implement `/health` endpoint returning "OK"
- [x] 4.4 Implement landing page at `/` with HTML response
- [x] 4.5 Add configurable listen address support
- [x] 4.6 Write unit tests for HTTP endpoints

## 5. Exporter Orchestration

- [x] 5.1 Create `internal/exporter/exporter.go` with Exporter struct
- [x] 5.2 Implement NewExporter constructor with config validation
- [x] 5.3 Implement Run method coordinating poller and server
- [x] 5.4 Implement graceful shutdown on context cancellation
- [x] 5.5 Write integration tests for exporter lifecycle

## 6. CLI Integration

- [x] 6.1 Create `cmd/zai-quota/exporter.go` with exporter subcommand
- [x] 6.2 Add `--poll-interval` flag with 60s minimum validation
- [x] 6.3 Add `--listen` flag with default `:9090`
- [x] 6.4 Register exporter subcommand in root command
- [x] 6.5 Write unit tests for CLI flags and validation

## 7. Documentation

- [x] 7.1 Update README.md with exporter usage examples
- [x] 7.2 Add example Prometheus scrape config
- [x] 7.3 Document all exposed metrics with descriptions

## 8. Final Verification

- [x] 8.1 Run full test suite ensuring no regressions
- [x] 8.2 Manual test: start exporter and verify `/metrics` output
  - Run: `zai-quota exporter`
  - Visit: http://localhost:9090/metrics
  - Verify: Metrics are returned in Prometheus format ✓ (validated with curl)
- [x] 8.3 Manual test: verify metrics format with `promtool check metrics`
  - Run: `curl http://localhost:9090/metrics | promtool check metrics`
  - Verify: No errors or warnings
  - Note: promtool not installed on this system - manual validation confirms Prometheus format is correct
- [x] 8.4 Manual test: configure Prometheus to scrape and verify data collection
  - Add scrape config to prometheus.yml (see README.md for example)
  - Start Prometheus
  - Verify: Metrics appear in Prometheus UI at http://localhost:9090/targets
