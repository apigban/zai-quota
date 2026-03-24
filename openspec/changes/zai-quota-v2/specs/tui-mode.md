# Spec: TUI Mode

## Capability

Interactive terminal user interface for displaying Z.ai quota information with manual refresh capability and multi-limit display.

## Requirements

### REQ-TUI-001: Default TUI Mode

When `zai-quota` is invoked without `--json` or `--yaml` flags, the tool SHALL launch an interactive TUI.

### REQ-TUI-002: Multi-Limit Display

The TUI SHALL display both limit types returned by the API:

| Type | Label | Description |
|------|-------|-------------|
| TIME_LIMIT | "[5-Hour Prompt Limit]" | Request count, resets every 5 hours |
| TOKENS_LIMIT | "[Weekly Quota Limit]" | Token percentage, resets every 7 days |

### REQ-TUI-003: TIME_LIMIT Display

For TIME_LIMIT, the TUI SHALL display:
- Colored progress bar based on percentage
- Usage as "X / Y" (currentValue / usage)
- Remaining count
- Time until reset
- "[e] Expand details" hint if usageDetails present

### REQ-TUI-004: TOKENS_LIMIT Display

For TOKENS_LIMIT, the TUI SHALL display:
- Colored progress bar based on percentage
- Percentage value
- Time until reset

### REQ-TUI-005: Usage Details Expand/Collapse

When TIME_LIMIT has usageDetails:
- Collapsed by default, showing "[e] Expand details"
- Pressing 'e' expands to show per-model breakdown:
  ```
  Usage breakdown:
    └─ search-prime: 5
    └─ model-alpha: 3
  ```
- Expanded state shows "[e] Collapse"
- Pressing 'e' again collapses

### REQ-TUI-006: Level Display

The TUI title bar SHALL show the subscription level:
- "Z.AI Quota Monitor [Pro]" for pro plan
- "Z.AI Quota Monitor [Lite]" for lite plan
- "Z.AI Quota Monitor [Max]" for max plan
- Level is capitalized (first letter uppercase)

### REQ-TUI-007: Manual Refresh

The TUI SHALL support manual refresh via `r` or `R` key:
- While refreshing, header SHALL show "Refreshing..." indicator
- Previous quota data SHALL remain visible during refresh
- User SHALL NOT be able to trigger multiple concurrent refreshes

### REQ-TUI-008: Quit

The TUI SHALL quit on any of:
- `q` key
- `Q` key
- `Ctrl+C`

### REQ-TUI-009: Error Overlay

When an API error occurs:
- An error overlay SHALL appear at the top of the quota display
- The overlay SHALL show the error message
- The overlay SHALL provide `[Dismiss]` and `[Retry]` options
- Previous quota data (if any) SHALL remain visible below the overlay

### REQ-TUI-010: Empty State

If the initial API fetch fails (no prior data):
- The TUI SHALL show the error overlay
- The TUI SHALL display "No quota data available" message
- The TUI SHALL show `[r] Retry` in help instead of `[r] Refresh`

### REQ-TUI-011: Progress Bar Colors

The progress bar color SHALL indicate usage level:

| Percentage | Color | Level |
|------------|-------|-------|
| 0-79% | Green (#4CAF50) | Safe |
| 80-89% | Amber (#FFC107) | Warning |
| 90-94% | Orange (#FF5722) | Critical |
| 95-100% | Red (#F44336) | Emergency |

### REQ-TUI-012: Terminal Resize

The TUI SHALL handle terminal resize events gracefully.

## Non-Requirements

- Auto-refresh or polling
- Multiple accounts
- Usage history
- Configuration editing
- Write operations

## Dependencies

- `charm.land/bubbletea/v2`
- `charm.land/bubbles/v2` (progress component)
- `charm.land/lipgloss/v2`

## Acceptance Criteria

- [ ] `zai-quota` launches TUI
- [ ] `zai-quota --json` outputs JSON and exits (unchanged)
- [ ] `zai-quota --yaml` outputs YAML and exits (unchanged)
- [ ] Both TIME_LIMIT and TOKENS_LIMIT display correctly
- [ ] TIME_LIMIT shows "X / Y" format
- [ ] TOKENS_LIMIT shows percentage only
- [ ] Pressing `r` refreshes quota
- [ ] Pressing `e` expands/collapses usage details
- [ ] Pressing `q` quits
- [ ] Level appears in title bar
- [ ] Error overlay appears on API failure
- [ ] Progress bar colors match thresholds
