# Proposal: tui-multi-limit

## Why

The current TUI was built on incorrect API assumptions. It expects a single quota (limit/remaining/reset) but the actual API returns multiple limit types (TIME_LIMIT and TOKENS_LIMIT) with different field structures. This causes the TUI to display "0 / 0" and "Unknown" reset times.

After fixing the API response parsing in the models and processor layers, the TUI needs to be rewritten to:
1. Display both TIME_LIMIT (5-hour prompt quota) and TOKENS_LIMIT (weekly token quota)
2. Show subscription level (lite/pro/max) in the title bar
3. Provide expandable usage details for TIME_LIMIT
4. Update labels to accurately reflect each limit type

## What Changes

- **Multi-limit display**: Show both TIME_LIMIT and TOKENS_LIMIT simultaneously with distinct visual presentation
- **Level display**: Show subscription level in title bar (e.g., "Z.AI Quota Monitor [Pro]")
- **Expand/collapse interaction**: Add 'e' key to expand/collapse usage details for TIME_LIMIT
- **Correct labels**: TIME_LIMIT → "[5-Hour Prompt Limit]", TOKENS_LIMIT → "[Weekly Quota Limit]"
- **Per-limit progress bars**: Each limit type gets its own colored progress bar based on percentage
- **Usage details display**: TIME_LIMIT shows per-model breakdown when expanded (e.g., "search-prime: 5")

**BREAKING**: TUI Model structure changes significantly - `quota *models.QuotaResponse` replaced with `limits []models.Limit` and new `processed map[string]ProcessedLimitData`

## Capabilities

### New Capabilities

- `multi-limit-display`: Display multiple quota limits (TIME_LIMIT and TOKENS_LIMIT) with distinct visual presentation
- `usage-details-expansion`: Expandable/collapsible view of per-model usage breakdown for TIME_LIMIT
- `level-indicator`: Display subscription level (lite/pro/max) in TUI title bar

### Modified Capabilities

None - this is a new TUI implementation, not a modification of existing specs.

## Impact

**Code**:
- `internal/tui/tui.go` - Major rewrite (~380 lines): Model struct, processing logic, view rendering, key handling
- `internal/tui/styles.go` - Add level style for title bar
- `internal/tui/tui_test.go` - Update test fixtures for new model structure

**APIs**:
- Consumes `models.QuotaResponse` with new structure (limits array, level field)
- Uses `processor.ProcessedLimit` with new fields (Label, Used, Total, UsageDetails)

**Dependencies**:
- Requires completed API response parsing fixes (models/quota.go, processor/processor.go)
- Requires updated formatters (already completed)

**Systems**:
- TUI mode (default when running `zai-quota` without flags)
- CLI modes (--json, --yaml) unaffected
