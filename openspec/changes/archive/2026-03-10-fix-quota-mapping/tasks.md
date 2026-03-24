# Tasks: fix-quota-mapping

## 1. Processor Label Updates

- [x] 1.1 Update `TOKENS_LIMIT` label in `internal/processor/processor.go`
  - Change from `"[Weekly Quota Limit]"` to `"[5-Hour Prompt Limit]"`
  - Add full data processing (used, total, remaining, usageDetails)

- [x] 1.2 Update `TIME_LIMIT` label in `internal/processor/processor.go`
  - Change from `"[5-Hour Prompt Limit]"` to `"[Tool Quota]"`
  - Add full data processing (used, total, remaining, usageDetails)

## 2. TUI Label and Processing Updates

- [x] 2.1 Update `processLimits()` in `internal/tui/tui.go`
  - `TOKENS_LIMIT` → `"[5-Hour Prompt Limit]"`
  - `TIME_LIMIT` → `"[Tool Quota]"`
  - Populate Total, Used, Remaining for both types

- [x] 2.2 Update `renderQuotaDisplay()` order
  - Render `TOKENS_LIMIT` first (prompts)
  - Render `TIME_LIMIT` second (tools)

## 3. Consolidated Rendering

- [x] 3.1 Replace `renderTimeLimit()` and `renderTokensLimit()` with `renderLimit()`
  - Single function handles both quota types
  - Use `data.Type` for expand/collapse state key

- [x] 3.2 Add fallback display logic
  - Show `"X% used"` when `total <= 0`
  - Hide `"Remaining: X"` when `remaining <= 0`

## 4. Expand/Collapse Updates

- [x] 4.1 Update 'e' key handler for both quota types
  - Toggle `TOKENS_LIMIT` expansion when it has `UsageDetails`
  - Toggle `TIME_LIMIT` expansion when it has `UsageDetails`

- [x] 4.2 Update `renderHelp()` to check both types
  - Check `TOKENS_LIMIT` for usage details
  - Check `TIME_LIMIT` for usage details

## 5. Test Updates

- [x] 5.1 Update `internal/processor/processor_test.go`
  - Verify `TOKENS_LIMIT` label is `"[5-Hour Prompt Limit]"`
  - Verify `TIME_LIMIT` label is `"[Tool Quota]"`
  - Update test fixtures with correct field mappings

- [x] 5.2 Update `internal/tui/tui_test.go`
  - Update label expectations
  - Test consolidated rendering
  - Test expand/collapse for both quota types

## Dependencies

```
1.1 → 1.2 → 2.1 → 2.2 → 3.1 → 3.2 → 4.1 → 4.2 → 5.1 → 5.2
```

## Estimated Effort

| Phase | Tasks | Time |
|-------|-------|------|
| Processor Labels | 2 tasks | ~15 min |
| TUI Labels | 2 tasks | ~10 min |
| Consolidated Rendering | 2 tasks | ~20 min |
| Expand/Collapse | 2 tasks | ~10 min |
| Tests | 2 tasks | ~20 min |

**Total: ~1.25 hours** (Already completed)
