# Design: fix-quota-mapping

## Context

The TUI and processor layers incorrectly mapped API quota types to user-facing labels. The Z.AI API returns two limit types:
- `TOKENS_LIMIT`: Short-term prompt quota (resets every ~5 hours)
- `TIME_LIMIT`: Tool usage quota (resets monthly)

**Previous (Incorrect) Mapping:**
- `TOKENS_LIMIT` → "[Weekly Quota Limit]" (wrong - no usage details)
- `TIME_LIMIT` → "[5-Hour Prompt Limit]" (wrong label, correct functionality)

**Correct Mapping (matching Z.AI website):**
- `TOKENS_LIMIT` → "[5-Hour Prompt Limit]" (has usage details per model)
- `TIME_LIMIT` → "[Tool Quota]" (has usage details per tool)

**Constraints:**
- Must not change API response structure
- Must maintain backward compatibility with all existing tests
- Must preserve progress bar color thresholds and warning levels

## Goals / Non-Goals

**Goals:**
- Correct label mapping to match Z.AI website behavior
- Enable expandable usage details for both quota types
- Consolidate rendering logic to reduce code duplication
- Add fallback display when `total`/`remaining` are not available

**Non-Goals:**
- Changing API client or data models
- Adding new quota types
- Modifying progress bar styling or thresholds

## Decisions

### Decision 1: Label Assignment Location

**Choice:** Update labels in both `processor.go` and `tui.go` to ensure consistency.

**Rationale:**
- Processor provides labels for CLI output (JSON, YAML, human)
- TUI has its own processing for display-specific formatting
- Both need correct labels to match Z.AI website

**Implementation:**
```go
// processor.go
if limit.Type == "TOKENS_LIMIT" {
    label = "[5-Hour Prompt Limit]"
} else if limit.Type == "TIME_LIMIT" {
    label = "[Tool Quota]"
}

// tui.go processLimits()
switch limit.Type {
case "TOKENS_LIMIT":
    data.Label = "[5-Hour Prompt Limit]"
case "TIME_LIMIT":
    data.Label = "[Tool Quota]"
}
```

### Decision 2: Consolidated Rendering Function

**Choice:** Replace separate `renderTimeLimit()` and `renderTokensLimit()` with single `renderLimit()` function.

**Rationale:**
- Both quota types now have same data structure (usage details, counts)
- Eliminates ~30 lines of duplicate code
- Easier to maintain and extend

**Implementation:**
- Single function handles all limit types
- Conditional display of `total`/`remaining` based on availability
- Uses `data.Type` for expand/collapse state key

### Decision 3: Fallback Display Logic

**Choice:** Show percentage-based display when `total` is 0 or unavailable.

**Rationale:**
- Some quota responses may not include `total` count
- Prevents showing confusing "0 / 0" display
- Falls back to "X% used" format

**Implementation:**
```go
if data.Total > 0 {
    // Show "X / Y" format
} else {
    // Show "X% used" format
}
```

### Decision 4: Dual Expand/Collapse Support

**Choice:** Enable expand/collapse for both `TOKENS_LIMIT` and `TIME_LIMIT`.

**Rationale:**
- Both quota types now have `UsageDetails` array
- Users want to see per-model and per-tool breakdown
- Toggle each independently using `expanded[data.Type]`

**Implementation:**
```go
// In handleKeyPress for 'e' key
if data, exists := m.processed["TOKENS_LIMIT"]; exists && len(data.UsageDetails) > 0 {
    m.expanded["TOKENS_LIMIT"] = !m.expanded["TOKENS_LIMIT"]
}
if data, exists := m.processed["TIME_LIMIT"]; exists && len(data.UsageDetails) > 0 {
    m.expanded["TIME_LIMIT"] = !m.expanded["TIME_LIMIT"]
}
```

## Risks / Trade-offs

### Risk 1: User Confusion During Transition
**Risk:** Users familiar with old labels may be confused
**Mitigation:** Labels now match Z.AI website, which is authoritative source

### Risk 2: Test Coverage Gaps
**Risk:** Existing tests may not cover all edge cases
**Mitigation:** All tests updated to verify correct label assignment

### Trade-off: Consolidation vs. Flexibility
**Trade-off:** Single render function reduces flexibility for quota-type-specific formatting
**Decision:** Accept trade-off - both quota types have identical display needs after fix
