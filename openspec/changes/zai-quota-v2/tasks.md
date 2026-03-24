# Tasks: zai-quota-v2

## Overview

This change fixes a critical API response parsing mismatch and adds multi-limit TUI display. The original TUI implementation was completed but built on incorrect API assumptions.

## Implementation Tasks

### Phase 1: Data Models

- [x] **1.1. Rewrite models/quota.go**
  - Replace `QuotaResponse` with correct structure (Limits []Limit, Level string)
  - Update `Limit` struct with all API fields (Type, Unit, Number, Usage, CurrentValue, Remaining, Percentage, NextResetTime, UsageDetails)
  - Add `UsageDetail` struct (ModelCode, Usage)
  - Add `omitempty` tags for optional fields
  - Delete old flat structure (limit/remaining/reset fields)

- [ ] **1.2. Update models/quota_test.go**
  - Rewrite tests for new struct definitions
  - Test JSON marshaling/unmarshaling with real API response format
  - Test both TIME_LIMIT and TOKENS_LIMIT field sets

### Phase 2: Processing

- [x] **2.1. Update processor/processor.go**
  - Update `ProcessedLimit` struct to include new fields (Label, Total, Used, UsageDetails)
  - Update `ProcessLimits()` to handle new field semantics:
    - TIME_LIMIT: Use CurrentValue for used, Usage for total
    - TOKENS_LIMIT: Only has Percentage
  - Generate correct labels per limit type

- [x] **2.2. Update processor/processor_test.go**
  - Update test fixtures with new Limit structure
  - Test TIME_LIMIT processing (currentValue/usage)
  - Test TOKENS_LIMIT processing (percentage only)

### Phase 3: Formatters

- [ ] **3.1. Update formatter/human.go**
  - Fix labels: TIME_LIMIT → "[5-Hour Prompt Limit]", TOKENS_LIMIT → "[Weekly Quota Limit]"
  - Update TIME_LIMIT display: show "Usage: X / Y" with currentValue/usage
  - Update TOKENS_LIMIT display: show "Usage: X%" only
  - Include usageDetails for TIME_LIMIT

- [ ] **3.2. Update formatter/colored.go**
  - Fix labels (same as human.go)
  - Update display logic for both limit types
  - Include usageDetails display

- [ ] **3.3. Update formatter/json.go**
  - Update `JSONLimit` struct with new fields
  - Include usage_details for TIME_LIMIT
  - Add level to output structure

- [ ] **3.4. Update formatter/yaml.go**
  - Update YAML structure to match JSON
  - Include usage_details for TIME_LIMIT

- [ ] **3.5. Update formatter tests**
  - Update human_test.go fixtures
  - Update colored_test.go fixtures
  - Update json_test.go fixtures
  - Update yaml_test.go fixtures

### Phase 4: TUI

- [ ] **4.1. Update tui/tui.go - Model struct**
  - Replace `quota *models.QuotaResponse` with `limits []models.Limit`
  - Add `level string` field
  - Replace `processed *ProcessedData` with `processed map[string]ProcessedLimitData`
  - Add `expanded map[string]bool` for expand/collapse state

- [ ] **4.2. Update tui/tui.go - Processing**
  - Rewrite `processQuota()` → `processLimits()` to handle []models.Limit
  - Generate ProcessedLimitData for each limit type
  - Calculate warning level per limit

- [ ] **4.3. Update tui/tui.go - View rendering**
  - Update `renderTitle()` to show level: "Z.AI Quota Monitor [Pro]"
  - Rewrite `renderQuotaDisplay()` for multi-limit display
  - Add expand/collapse rendering for TIME_LIMIT usageDetails
  - Update help text to include 'e' key

- [ ] **4.4. Update tui/tui.go - Key handling**
  - Add 'e'/'E' handler to toggle expanded state for TIME_LIMIT

- [ ] **4.5. Update tui/styles.go**
  - Add `levelStyle` for level display in title bar
  - Use gold color (#FFD700) for level

- [ ] **4.6. Update tui/tui_test.go**
  - Update test fixtures with new limit structure
  - Test expand/collapse key handling
  - Test multi-limit display

### Phase 5: Root Command

- [ ] **5.1. Update cmd/zai-quota/root.go**
  - Simplify runCLI() to use quota.Limits directly
  - Remove fake limit construction logic
  - Pass quota.Level to TUI

### Phase 6: API Client Tests

- [ ] **6.1. Update api/client_test.go**
  - Update mock responses to use new API structure
  - Test parsing of limits array
  - Test level field extraction

### Phase 7: Verification

- [ ] **7.1. Run full test suite**
  - `go test ./...` - all tests pass
  - `go vet ./...` - no issues

- [ ] **7.2. Manual verification**
  - Test TUI with real API
  - Verify both limit types display correctly
  - Verify expand/collapse works
  - Verify level shows in title bar
  - Verify --json and --yaml output

- [ ] **7.3. Update README**
  - Document multi-limit display
  - Update TUI screenshot/description
  - Document 'e' key for expand/collapse

## Task Dependencies

```
Phase 1 (Models)
    │
    ├──────────────────┬──────────────────┐
    ▼                  ▼                  ▼
Phase 2           Phase 3            Phase 6
(Processor)       (Formatters)       (API Tests)
    │                  │                  │
    └──────────────────┼──────────────────┘
                       ▼
                  Phase 4 (TUI)
                       │
                       ▼
                  Phase 5 (Root)
                       │
                       ▼
                  Phase 7 (Verify)
```

## Estimated Effort

| Phase | Tasks | Effort |
|-------|-------|--------|
| Phase 1 | Data Models | ~30 min |
| Phase 2 | Processing | ~20 min |
| Phase 3 | Formatters | ~45 min |
| Phase 4 | TUI | ~1.5 hours |
| Phase 5 | Root Command | ~15 min |
| Phase 6 | API Tests | ~20 min |
| Phase 7 | Verification | ~30 min |

**Total: ~4 hours**

## Critical Changes Summary

1. **models/quota.go** - Complete rewrite of struct definitions
2. **processor/processor.go** - Handle new field semantics (currentValue vs usage)
3. **formatter/*.go** - Fix labels (swap TIME_LIMIT/TOKENS_LIMIT descriptions)
4. **tui/tui.go** - Multi-limit display, expand/collapse, level in title
5. **All test files** - Update fixtures to match real API response
