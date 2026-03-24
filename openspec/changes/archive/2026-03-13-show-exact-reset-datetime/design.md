## Context

The TUI currently displays reset times as relative durations only (`Reset: 4h 6m`). The underlying data flow:

```
API Response (NextResetTime: int64 ms)
         ↓
ConvertTimestamp() → time.Time
         ↓
FormatTimeUntil() → "4h 6m" or "20 days"
         ↓
Displayed as: "Reset: 4h 6m"
```

The absolute `time.Time` is available but not exposed to the display layer.

## Goals / Non-Goals

**Goals:**
- Show exact date/time alongside duration in reset display
- Detect and display local timezone (offset + city name)
- Graceful UTC fallback when timezone cannot be detected
- Cross-platform support (Linux, macOS)

**Non-Goals:**
- Windows timezone detection (can be added later)
- User-configurable timezone override
- Persisting timezone preference

## Decisions

### D1: Timezone Detection Strategy

**Decision:** Check multiple sources in priority order:
1. `TZ` environment variable
2. `/etc/timezone` file (Debian/Ubuntu)
3. `/etc/localtime` symlink resolution (RHEL/CentOS/macOS)
4. Fallback to UTC

**Rationale:** No external dependencies, covers most Unix-like systems.

**Alternatives considered:**
- Third-party library (adds dependency weight)
- Only use numeric offset (less descriptive)

### D2: Display Format

**Decision:** `Reset: {duration} ({date}, {time} {offset} {city})`

Examples:
- `Reset: 4h 6m (Mar 13, 20:00 +04:00 Dubai)`
- `Reset: 20 days (Apr 2, 00:00 +04:00 Dubai)`
- `Reset: 4h 6m (Mar 13, 16:00 UTC)` (fallback)

**Rationale:** Provides complete context - how long until reset AND exactly when.

**Alternatives considered:**
- Time only, no date (ambiguous for multi-day resets)
- Date only for long durations (lose precision)
- IANA name instead of city (too verbose: "Asia/Dubai")

### D3: City Name Extraction

**Decision:** Extract city from IANA name, replace underscores with spaces.

**Examples:**
- `Asia/Dubai` → `Dubai`
- `America/New_York` → `New York`
- `Europe/London` → `London`
- `UTC` → `UTC` (no transformation)

**Rationale:** City names are more recognizable than full IANA identifiers.

## Risks / Trade-offs

| Risk | Mitigation |
|------|------------|
| Timezone detection fails on some systems | Graceful UTC fallback with clear indication |
| Symlink resolution fails in containers | Try multiple detection methods |
| User in different TZ than server | Uses local system TZ, which is correct for terminal users |
