# Tasks: tui-multi-limit

## 1. Model Structure Updates

- [x] 1.1 Update Model struct in internal/tui/tui.go
  - Replace `quota *models.QuotaResponse` with `limits []models.Limit`
  - Add `level string` field
  - Replace `processed *ProcessedData` with `processed map[string]ProcessedLimitData`
  - Add `expanded map[string]bool` field
  
- [x] 1.2 Create ProcessedLimitData struct
  - Add Type, Label, Percentage fields
  - Add Used, Total, Remaining fields
  - Add ResetDisplay, WarningLevel fields
  - Add UsageDetails []models.UsageDetail field

- [x] 1.3 Update NewModel() constructor
  - Initialize processed as empty map
  - Initialize expanded as empty map
  - Initialize limits as empty slice

## 2. Processing Logic Updates

- [x] 2.1 Replace processQuota() with processLimits()
  - Accept []models.Limit as input
  - Return map[string]ProcessedLimitData
  - Iterate over limits and process each one
  
- [x] 2.2 Implement TIME_LIMIT processing
  - Map Usage field → Total
  - Map CurrentValue field → Used
  - Map Remaining field → Remaining
  - Copy UsageDetails array
  - Set Label to "[5-Hour Prompt Limit]"
  
- [x] 2.3 Implement TOKENS_LIMIT processing
  - Only populate Percentage and reset time
  - Set Used, Total, Remaining to 0
  - Set UsageDetails to nil
  - Set Label to "[Weekly Quota Limit]"

- [x] 2.4 Implement warning level calculation
  - Calculate based on percentage thresholds
  - <80% = "safe", 80-89% = "warning", 90-94% = "critical", ≥95% = "emergency"

- [x] 2.5 Update handleQuotaFetched()
  - Extract limits from msg.quota.Limits
  - Extract level from msg.quota.Level
  - Call processLimits() instead of processQuota()

## 3. View Rendering Updates

- [x] 3.1 Update renderTitle() to show level
  - Capitalize level using strings.Title
  - Format as "Z.AI Quota Monitor [Level]"
  - Style level with gold color (#FFD700) and bold
  
- [x] 3.2 Replace renderQuotaDisplay() with multi-limit rendering
  - Iterate over processed map
  - Render each limit type in its own section
  - Add spacing between sections
  
- [x] 3.3 Implement TIME_LIMIT rendering
  - Display label: "[5-Hour Prompt Limit]"
  - Show progress bar with percentage
  - Show usage as "X / Y"
  - Show remaining count
  - Show reset time
  - Show expand/collapse hint if usageDetails present
  
- [x] 3.4 Implement TOKENS_LIMIT rendering
  - Display label: "[Weekly Quota Limit]"
  - Show progress bar with percentage
  - Show percentage value
  - Show reset time
  - NO usage/remaining display
  
- [x] 3.5 Implement usage details expansion rendering
  - When expanded: show "Usage breakdown:" header
  - Show each model as "└─ <modelCode>: <usage>"
  - Show "[e] Collapse" hint
  - When collapsed: show "[e] Expand details" hint

- [x] 3.6 Update renderHelp() to include 'e' key
  - Add 'e' for expand/collapse when usageDetails available
  - Format: "[r] Refresh  [e] Expand  [q] Quit"

## 4. Key Handling Updates

- [x] 4.1 Add 'e'/'E' handler in handleKeyPress()
  - Check if TIME_LIMIT has usageDetails
  - Toggle expanded["TIME_LIMIT"] state
  - No-op if no usageDetails available

## 5. Style Updates

- [x] 5.1 Add levelStyle to internal/tui/styles.go
  - Set foreground color to gold (#FFD700)
  - Set bold to true
  - Use in renderTitle()

## 6. Test Updates

- [x] 6.1 Update internal/tui/tui_test.go
  - Replace QuotaResponse{Limit:..., Remaining:..., Reset:...} with QuotaResponse{Limits:[]Limit{...}, Level:"pro"}
  - Update all test fixtures with new structure
  - Add tests for multi-limit display
  - Add tests for expand/collapse functionality
  - Add tests for level display

- [x] 6.2 Update internal/api/client_test.go
  - Update mock responses to use new API structure
  - Update assertions to check limits array
  - Update assertions to check level field

- [x] 6.3 Update internal/formatter/json_test.go
  - Update ProcessedLimit fixtures to use new fields
  - Add level parameter to FormatJSON calls
  - Update test expectations

- [x] 6.4 Update internal/formatter/yaml_test.go
  - Update ProcessedLimit fixtures to use new fields
  - Add level parameter to FormatYAML calls
  - Update test expectations

- [x] 6.5 Update internal/formatter/human_test.go
  - Update ProcessedLimit fixtures to use new fields
  - Update test expectations

- [x] 6.6 Update internal/formatter/colored_test.go
  - Update ProcessedLimit fixtures to use new fields
  - Rename FormatProcessedLimitWithColor to FormatColored
  - Update test expectations

## 7. Verification

- [x] 7.1 Run full test suite
  - Execute `go test ./...`
  - Verify all tests pass
  - Fix any remaining compilation errors

- [ ] 7.2 Manual TUI testing
  - Test with real API (if API key available)
  - Verify both limits display correctly
  - Test expand/collapse with 'e' key
  - Verify level shows in title bar
  - Test refresh with 'r' key
  - Verify error handling still works
  - **BUG FIX**: Fixed timestamp conversion in `internal/tui/tui.go:210`
    - Changed `time.Unix(limit.NextResetTime, 0)` to `processor.ConvertTimestamp(limit.NextResetTime)`
    - API returns milliseconds, but `time.Unix()` expects seconds

- [x] 7.3 Update README.md
  - Document multi-limit display feature
  - Add 'e' key to keyboard shortcuts
  - Update TUI screenshot (if applicable)
  - Document level indicator in title bar

## Dependencies

```
1.1 → 1.2 → 1.3
         ↓
      2.1 → 2.2 → 2.3 → 2.4 → 2.5
         ↓
      3.1 → 3.2 → 3.3 → 3.4 → 3.5 → 3.6
         ↓               ↓
      4.1             5.1
         ↓
      6.1 → 6.2 → 6.3 → 6.4 → 6.5 → 6.6
         ↓
      7.1 → 7.2 → 7.3
```

## Estimated Effort

| Phase | Tasks | Time |
|-------|-------|------|
| Phase 1: Model Structure | 3 tasks | ~20 min |
| Phase 2: Processing Logic | 5 tasks | ~30 min |
| Phase 3: View Rendering | 6 tasks | ~45 min |
| Phase 4: Key Handling | 1 task | ~10 min |
| Phase 5: Styles | 1 task | ~5 min |
| Phase 6: Tests | 6 tasks | ~40 min |
| Phase 7: Verification | 3 tasks | ~20 min |

**Total: ~3 hours**
