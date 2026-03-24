## Context

The zai-quota tool currently supports:
- Interactive TUI mode (default)
- One-shot JSON/YAML/text output for scripting

Users need Prometheus integration for:
- Historical trend tracking in Grafana
- Alerting when quota approaches limits
- Integration with existing observability stacks

**Critical Constraint**: Z.ai API enforces minimum 60-second interval between requests. This precludes on-demand fetching per Prometheus scrape and requires a polling/caching architecture.

## Goals / Non-Goals

**Goals:**
- Expose quota metrics in Prometheus exposition format
- Support configurable polling interval (minimum 60s)
- Provide per-tool usage breakdown
- Include exporter health metrics
- Follow Prometheus naming conventions and best practices

**Non-Goals:**
- Per-model prompt usage breakdown (API doesn't provide this for TOKENS_LIMIT)
- Authentication/TLS for the metrics endpoint (out of scope for initial version)
- Multi-target exporter pattern (single Z.ai account per exporter instance)

## Decisions

### D1: Polling Architecture (Cached Metrics)

**Decision**: Background poller with in-memory cache, not scrape-on-demand.

**Rationale**: 
- API requires 60s minimum between requests
- Prometheus may scrape at any interval; on-demand would violate rate limit
- Cache ensures API is called at most once per poll interval

**Alternatives Considered**:
- On-demand fetch: Rejected - would violate 60s API constraint
- File-based cache: Rejected - unnecessary complexity for single-process exporter

```
┌──────────────┐    60s min    ┌─────────────┐    on-demand    ┌────────┐
│   Z.ai API   │◀─────────────│   Poller    │────────────────▶│ Cache  │
└──────────────┘               └─────────────┘                 └───┬────┘
                                                                    │
                                    ┌───────────────────────────────┘
                                    ▼
                              ┌───────────┐
                              │  /metrics │◀── Prometheus
                              └───────────┘
```

### D2: Prometheus Client Library

**Decision**: Use `github.com/prometheus/client_golang` for metric exposition.

**Rationale**:
- Standard library for Go Prometheus exporters
- Handles text exposition format correctly
- Built-in support for gauges, counters, labels
- Active maintenance and wide adoption

### D3: Metric Naming Convention

**Decision**: Use `zai_quota_` namespace with Prometheus best practices.

**Rationale**: Per Prometheus naming guidelines:
- Single-word application prefix: `zai_quota_`
- Base units: `_seconds` for timestamps, no unit for counts
- Ratio not percentage: `_ratio` (0-1) instead of `_percent` (0-100)
- Gauge semantics: Values can increase and decrease

**Metric Schema**:

```promql
# Prompt usage (TOKENS_LIMIT) - only percentage available
zai_quota_prompt_usage_ratio                              gauge
zai_quota_prompt_reset_timestamp_seconds                  gauge

# Tool calls (TIME_LIMIT) - full details available  
zai_quota_tool_calls_used                                 gauge
zai_quota_tool_calls_limit                                gauge
zai_quota_tool_calls_remaining                            gauge
zai_quota_tool_calls_reset_timestamp_seconds              gauge
zai_quota_tool_calls_by_tool{tool="..."}                  gauge

# Metadata
zai_quota_info{level="..."}                               gauge (always 1)

# Exporter health
zai_quota_up                                              gauge (1=ok, 0=error)
zai_quota_last_scrape_timestamp_seconds                   gauge
zai_quota_scrape_duration_seconds                         gauge
```

### D4: Subcommand Integration

**Decision**: Add as `exporter` subcommand, not separate binary.

**Rationale**:
- Shares existing config loading, API client, models
- Single binary distribution remains simple
- Consistent with existing `--json`, `--yaml` output modes

**Usage**:
```bash
zai-quota exporter [--poll-interval=60] [--listen=:9090]
```

### D5: Minimum Poll Interval Enforcement

**Decision**: Hard minimum of 60s, reject lower values at startup.

**Rationale**:
- API constraint is non-negotiable
- Clear error message better than silent rate limiting
- Default of 60s matches minimum (conservative)

## Risks / Trade-offs

**R1: Stale Metrics** → Mitigation: Document that metrics may be up to `poll-interval` seconds old. Include `zai_quota_last_scrape_timestamp_seconds` so users can detect staleness.

**R2: API Downtime** → Mitigation: Continue serving last cached metrics. Set `zai_quota_up=0` to indicate scrape failure. Don't crash the exporter on API errors.

**R3: Memory Growth** → Mitigation: Cache is bounded (single QuotaResponse). No unbounded growth possible.

**R4: Clock Drift on Reset Timestamps** → Mitigation: Use API-provided timestamps directly (milliseconds → seconds conversion only). No local time calculations.
