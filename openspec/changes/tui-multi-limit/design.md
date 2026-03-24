# Design: tui-multi-limit

## Context

The TUI (internal/tui/tui.go) was originally built assuming a single quota with limit/remaining/reset fields. The actual API returns multiple limit types (TIME_LIMIT and TOKENS_LIMIT) with different field structures. After updating the models (models/quota.go) and processor (processor/processor.go), the TUI needs a major rewrite to handle multi-limit display.

**Current State:**
- TUI Model stores `quota *models.QuotaResponse` (now has `Limits []Limit, Level string`)
- TUI processes single quota into `ProcessedData` struct
- TUI View renders single progress bar and usage display
- No expand/collapse functionality

**Constraints:**
- Must use Bubble Tea v2 framework (already integrated)
- Must preserve error handling, refresh, and quit functionality
- Must work with existing processor.ProcessedLimit structure
- Must maintain current progress bar color thresholds

## Goals / Non-Goals

**Goals:**
- Display both TIME_LIMIT and TOKENS_LIMIT with distinct visual presentation
- Add expand/collapse for TIME_LIMIT usage details
- Show subscription level in title bar
- Update labels to reflect actual limit types (5-hour vs weekly)
- Maintain all existing functionality (refresh, error handling, quit)

**Non-Goals:**
- Auto-refresh/polling (manual refresh only)
- Modifying API client or processor logic (already complete)
- Changing CLI modes (--json, --yaml)
- Adding configuration options for display preferences

## Decisions

### Decision 1: Model Structure Change

**Choice:** Replace `quota *models.QuotaResponse` with `limits []models.Limit` and `level string` in Model struct.

**Rationale:**
- Processor now returns `[]ProcessedLimit` directly from `quota.Limits`
- Level is top-level field in API response, needs separate storage
- Cleaner separation of concerns - Model stores raw data, processed data computed for display

**Alternative Considered:** Keep `quota *models.QuotaResponse` and extract limits/level from it
- **Why not:** Would require quota field to never be nil, complicates error handling

### Decision 2: Processed Data Structure

**Choice:** Replace `processed *ProcessedData` with `processed map[string]ProcessedLimitData` keyed by limit type.

**Rationale:**
- Easy lookup by limit type (TIME_LIMIT, TOKENS_LIMIT)
- Supports variable number of limits
- Natural fit for iterating in View

**Structure:**
```go
type ProcessedLimitData struct {
    Type          string
    Label         string
    Percentage    int
    Used          int
    Total         int
    Remaining     int
    ResetDisplay  string
    WarningLevel  string
    UsageDetails  []models.UsageDetail
}
```

**Alternative Considered:** Keep slice `[]ProcessedLimitData`
- **Why not:** Requires linear search for specific limit type, less efficient

### Decision 3: Expand/Collapse State Management

**Choice:** Add `expanded map[string]bool` to Model to track expand state per limit type.

**Rationale:**
- Supports multiple limits (future-proof)
- Simple toggle logic
- Persists across refreshes

**Alternative Considered:** Single boolean `expanded bool`
- **Why not:** Won't scale if more limit types added later

### Decision 4: Level Display in Title Bar

**Choice:** Format as "Z.AI Quota Monitor [Pro]" with gold color and bold styling.

**Rationale:**
- Provides clear visual indicator of subscription tier
- Gold color distinguishes from other UI elements
- Capitalization (strings.Title) normalizes API values

**Alternative Considered:** Show level in separate line below title
- **Why not:** Wastes vertical space, title bar is natural location

### Decision 5: View Rendering Strategy

**Choice:** Iterate over `processed` map and render each limit type with its own section.

**Rationale:**
- Clean separation between limit types
- Each limit gets full width for progress bar
- Natural vertical layout

**Rendering Order:**
1. Title bar with level
2. Error overlay (if error)
3. TIME_LIMIT section (if present)
4. TOKENS_LIMIT section (if present)
5. Help text
6. Debug log (if debug mode)

**Alternative Considered:** Side-by-side layout
- **Why not:** Progress bars need full width, vertical space more usable

### Decision 6: Usage Details Expansion

**Choice:** Display as hierarchical list with tree-style indentation.

**Format:**
```
Usage breakdown:
  └─ search-prime: 5
  └─ model-alpha: 3
```

**Rationale:**
- Clear visual hierarchy
- Consistent with common CLI tool patterns
- Compact yet readable

**Alternative Considered:** Table format
- **Why not:** Overkill for 2-4 items, wastes horizontal space

## Risks / Trade-offs

### Risk 1: Large TUI File
**Risk:** Single 380-line file may become difficult to maintain
**Mitigation:** 
- Keep functions small and focused
- Extract rendering helpers (renderTimeLimit, renderTokensLimit)
- Consider splitting in future if more features added

### Risk 2: State Management Complexity
**Risk:** More state fields (limits, level, processed, expanded) increase bug surface
**Mitigation:**
- Clear separation: raw data (limits/level) vs. computed data (processed)
- Use Update() for all state changes
- Comprehensive tests for state transitions

### Risk 3: Backward Compatibility
**Risk:** Changes break existing TUI behavior users expect
**Mitigation:**
- Preserve all keybindings (r, q, ctrl+c)
- Keep refresh behavior identical
- Maintain error overlay functionality
- No breaking changes to user-facing features, only additions

### Trade-off: Memory vs. Simplicity
**Trade-off:** Storing both raw limits and processed data duplicates information
**Decision:** Accept trade-off for cleaner separation and easier View rendering
**Rationale:** Data is small (2-4 limits), memory impact negligible

## Migration Plan

### Phase 1: Update Model Structure
1. Add new fields to Model struct (limits, level, processed, expanded)
2. Remove old fields (quota, processed *ProcessedData)
3. Update NewModel() constructor

### Phase 2: Update Processing Logic
1. Replace `processQuota()` with `processLimits()` 
2. Convert []models.Limit → map[string]ProcessedLimitData
3. Add label generation logic
4. Preserve reset time formatting using processor.FormatTimeUntil()

### Phase 3: Update View Rendering
1. Update `renderTitle()` to show level
2. Replace `renderQuotaDisplay()` with multi-limit rendering
3. Add expand/collapse rendering for usage details
4. Update `renderHelp()` to include 'e' key

### Phase 4: Update Key Handling
1. Add 'e'/'E' handler in `handleKeyPress()`
2. Toggle expanded state for TIME_LIMIT

### Phase 5: Update fetchQuotaCmd
1. No changes needed - already returns *models.QuotaResponse
2. handleQuotaFetched() extracts limits and level

### Rollback Strategy
If critical issues arise:
1. Revert Model struct to single quota approach
2. Revert processQuota() function
3. Revert View rendering to single limit display

Git tags will mark pre-change state for easy rollback.

## Open Questions

None - all technical decisions have been resolved through exploration and user clarification.
