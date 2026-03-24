# Proposal: zai-quota-v2

## Summary

Fix critical API response parsing mismatch and add an interactive TUI mode to zai-quota using Bubble Tea v2, while preserving existing CLI output modes (`--json`, `--yaml`) for scripting.

## Problem Statement

The current implementation was built on incorrect assumptions about the API response structure:

**What the code expected:**
```json
{
  "data": {
    "limit": 10000,
    "remaining": 7500,
    "reset": "2024-03-08T00:00:00Z"
  }
}
```

**What the API actually returns:**
```json
{
  "code": 200,
  "msg": "Operation successful",
  "data": {
    "limits": [
      {
        "type": "TIME_LIMIT",
        "unit": 5,
        "number": 1,
        "usage": 1000,
        "currentValue": 8,
        "remaining": 992,
        "percentage": 1,
        "nextResetTime": 1775186469998,
        "usageDetails": [{"modelCode": "search-prime", "usage": 1}]
      },
      {
        "type": "TOKENS_LIMIT",
        "unit": 3,
        "number": 5,
        "percentage": 33,
        "nextResetTime": 1773052446291
      }
    ],
    "level": "pro"
  },
  "success": true
}
```

This mismatch causes the TUI to display "0 / 0" and "Unknown" reset times.

## Motivation

1. **Fix the core bug**: API response parsing is completely broken
2. **Improve user experience**: Interactive TUI for monitoring quota throughout a work session
3. **Support multiple limit types**: Display both 5-hour and weekly limits correctly
4. **Preserve scripting capability**: Keep `--json` and `--yaml` modes for automation

## Scope

### In Scope

- **API Response Parsing**: Rewrite models to match actual API structure
- **Multi-Limit Display**: Show both TIME_LIMIT and TOKENS_LIMIT
- **Usage Details**: Expandable breakdown of per-model usage for TIME_LIMIT
- **Level Display**: Show subscription level (lite/pro/max) in title bar
- **Interactive TUI mode** (default when no flags)
- **Manual refresh** with `r` key
- **Error overlay** with Dismiss/Retry options
- **Preserve last known data** on refresh failure
- **Empty state** with error on initial load failure
- **Keep existing CLI modes**: `--json`, `--yaml`, `--help`, `--version`

### Out of Scope

- Auto-refresh / polling
- Multiple accounts
- Usage history / graphs
- Configuration editing within TUI
- Write operations (reset session, etc.)

## Limit Types

Based on Z.ai documentation:

| Type | Description | Reset Period | Has Details |
|------|-------------|--------------|-------------|
| TIME_LIMIT | Prompt/request count | 5 hours | Yes (per-model breakdown) |
| TOKENS_LIMIT | Token consumption percentage | 7 days | No |

## Plan Levels

| Level | 5-Hour Prompts | Weekly Prompts |
|-------|----------------|----------------|
| lite | ~80 | ~400 |
| pro | ~400 | ~2,000 |
| max | ~1,600 | ~8,000 |

## Behavior Matrix

| Invocation | Behavior |
|------------|----------|
| `zai-quota` | Launch TUI (interactive) |
| `zai-quota --json` | Output JSON, exit immediately |
| `zai-quota --yaml` | Output YAML, exit immediately |
| `zai-quota --debug` | Launch TUI with debug logging |
| `zai-quota --help` | Show help, exit immediately |
| `zai-quota --version` | Show version, exit immediately |

## TUI Controls

| Key | Action |
|-----|--------|
| `r` / `R` | Refresh quota |
| `e` / `E` | Expand/collapse usage details |
| `q` / `Q` / `Ctrl+C` | Quit |

## Success Criteria

- [ ] API response parsing works correctly
- [ ] TUI displays both limit types with correct labels
- [ ] TIME_LIMIT shows expandable usage details
- [ ] Subscription level appears in title bar
- [ ] Progress bar colors match thresholds
- [ ] Manual refresh works without flicker
- [ ] Error overlay appears on failure
- [ ] CLI flags work unchanged for scripting
- [ ] All tests pass

## Dependencies

- `charm.land/bubbletea/v2`
- `charm.land/bubbles/v2` (progress component)
- `charm.land/lipgloss/v2`

## Risks

| Risk | Mitigation |
|------|------------|
| API structure changes | Use omitempty tags, validate required fields |
| Bubble Tea v2 stability | v2 is actively developed; monitor for breaking changes |
| Terminal compatibility | Test on common terminals (iTerm2, GNOME Terminal, Windows Terminal) |

## Timeline

Estimated 12-15 tasks covering:
1. Rewrite data models (2 tasks)
2. Update processor logic (1 task)
3. Update formatters (4 tasks)
4. Update TUI for multi-limit display (3 tasks)
5. Update tests (1 task)
6. Verification (1 task)
