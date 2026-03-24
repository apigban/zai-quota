## Why

Users currently see reset times as relative durations only (e.g., "Reset: 4h 6m" or "Reset: 20 days"). This requires mental calculation to determine the exact reset time. For long durations like "20 days", users cannot easily plan around when their quota will reset without manually calculating the date.

## What Changes

- Display exact date/time alongside duration in reset display
- Include timezone information (offset + city name from IANA timezone)
- Graceful fallback to UTC when timezone cannot be detected
- Apply to both TOKENS_LIMIT and TIME_LIMIT quota displays

**New display format:**
```
Reset: 4h 6m (Mar 13, 20:00 +04:00 Dubai)
Reset: 20 days (Apr 2, 00:00 +04:00 Dubai)
Reset: 4h 6m (Mar 13, 16:00 UTC)    ← fallback when timezone unavailable
```

## Capabilities

### New Capabilities

- `timezone-display`: Automatic detection and display of local timezone with UTC fallback

### Modified Capabilities

- `multi-limit-display`: Reset display format changes from duration-only to duration + exact datetime with timezone

## Impact

- **Files modified:**
  - `internal/processor/time.go` - Add new formatting function
  - `internal/processor/timezone.go` - New file for timezone detection
  - `internal/tui/tui.go` - Update ProcessedLimitData and renderLimit()
- **No API changes** - Internal display formatting only
- **No breaking changes** - Existing functionality preserved, display enhanced
