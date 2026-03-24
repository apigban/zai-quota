## Why

Users need to monitor Z.ai quota usage through Prometheus-compatible tooling for integration with existing observability stacks (Grafana, Alertmanager, etc.). Currently, the tool only supports interactive TUI viewing or one-shot JSON/YAML output, making it impossible to track usage trends, set up alerts, or integrate with monitoring dashboards.

The Z.ai API enforces a minimum 60-second interval between requests, requiring a polling/caching architecture rather than on-demand fetches.

## What Changes

- New `exporter` subcommand that runs a Prometheus-compatible metrics server
- Background polling of Z.ai API with configurable interval (minimum 60s, default 60s)
- HTTP endpoints: `/metrics` (Prometheus format), `/health`, and landing page
- Metrics for both prompt usage (TOKENS_LIMIT) and tool call quota (TIME_LIMIT)
- Per-tool breakdown metrics for tool call usage
- Exporter health metrics (up status, last scrape timestamp, scrape duration)

## Capabilities

### New Capabilities

- `prometheus-exporter`: Prometheus-compatible metrics server that polls Z.ai API and exposes quota metrics for monitoring systems

### Modified Capabilities

(None - this is a new feature with no requirement changes to existing capabilities)

## Impact

- **New packages**: `internal/exporter` (HTTP server, polling, caching), `internal/prometheus` (metric formatting)
- **Command changes**: New subcommand registration in `cmd/zai-quota/`
- **Config extension**: New flags `--poll-interval` and `--listen` for exporter mode
- **Dependencies**: Will need Prometheus client library or custom exposition formatter
- **No breaking changes**: Existing TUI and CLI output modes unchanged
