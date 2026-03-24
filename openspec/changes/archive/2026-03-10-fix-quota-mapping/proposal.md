# Proposal: fix-quota-mapping

## Why

The quota type labels were incorrectly mapped in the TUI and processor. `TOKENS_LIMIT` was labeled as "[Weekly Quota Limit]" when it actually represents the 5-hour prompt limit, and `TIME_LIMIT` was labeled as "[5-Hour Prompt Limit]" when it represents the tool quota. This caused user confusion and misaligned reset cycle displays.

## What Changes

- **Correct quota mapping**: `TOKENS_LIMIT` → "[5-Hour Prompt Limit]" (short-term reset), `TIME_LIMIT` → "[Tool Quota]" (monthly reset)
- **Expandable details for both**: Both quota types now support expand/collapse for usage breakdown
- **Consolidated rendering**: Single `renderLimit()` function handles both quota types
- **Robust field handling**: Fallbacks for missing `total`/`remaining` fields using percentage-based display

## Capabilities

### New Capabilities

None - this is a correction to existing implementation.

### Modified Capabilities

- `multi-limit-display`: Label mapping corrected to match Z.AI website behavior
- `usage-details-expansion`: Now supports both `TOKENS_LIMIT` and `TIME_LIMIT` (previously only `TIME_LIMIT`)

## Impact

**Code**:
- `internal/processor/processor.go` - Corrected label assignment in `ProcessLimits()`
- `internal/tui/tui.go` - Consolidated rendering, added expand for both limits, fallback display logic

**Tests**:
- `internal/processor/processor_test.go` - Updated test expectations for correct labels
- `internal/tui/tui_test.go` - Updated tests for consolidated rendering and dual expand
