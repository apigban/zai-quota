# Spec: usage-details-expansion

## ADDED Requirements

### Requirement: Expand usage details for TIME_LIMIT

The TUI SHALL allow users to expand TIME_LIMIT to view per-model usage breakdown when usageDetails are available.

#### Scenario: TIME_LIMIT has usage details available
- **WHEN** TIME_LIMIT includes usageDetails array with model breakdowns
- **THEN** TUI displays "[e] Expand details" hint below the limit section
- **AND** pressing 'e' key expands the view to show usage breakdown

#### Scenario: TIME_LIMIT has no usage details
- **WHEN** TIME_LIMIT has empty or nil usageDetails
- **THEN** TUI does NOT display the expand hint
- **AND** pressing 'e' key has no effect on this limit

### Requirement: Collapse expanded usage details

The TUI SHALL allow users to collapse expanded usage details by pressing 'e' again.

#### Scenario: Collapse after expansion
- **WHEN** usage details are currently expanded
- **AND** user presses 'e' key
- **THEN** usage details collapse and display changes to "[e] Expand details"

### Requirement: Display usage breakdown

The TUI SHALL display per-model usage breakdown in a hierarchical format when expanded.

#### Scenario: Multiple models in breakdown
- **WHEN** TIME_LIMIT has usageDetails with multiple models
- **THEN** TUI displays each model on a separate line with indentation
- **AND** format is "└─ <modelCode>: <usage>"

#### Scenario: Usage breakdown format
- **WHEN** TIME_LIMIT has usageDetails: [{"modelCode": "search-prime", "usage": 5}, {"modelCode": "model-alpha", "usage": 3}]
- **THEN** expanded view shows:
  ```
  Usage breakdown:
    └─ search-prime: 5
    └─ model-alpha: 3
  ```

### Requirement: Update help text for expand key

The TUI SHALL include 'e' key in help text when usage details are available.

#### Scenario: Help text with expandable details
- **WHEN** TIME_LIMIT has usageDetails available
- **THEN** help text displays "[r] Refresh  [e] Expand  [q] Quit"

#### Scenario: Help text without expandable details
- **WHEN** no limits have usageDetails
- **THEN** help text displays "[r] Refresh  [q] Quit"

### Requirement: Expand state persists across refreshes

The TUI SHALL preserve expand/collapse state when data is refreshed.

#### Scenario: Expanded state maintained after refresh
- **WHEN** user has expanded TIME_LIMIT usage details
- **AND** user presses 'r' to refresh
- **THEN** after refresh completes, usage details remain expanded

#### Scenario: Collapsed state maintained after refresh
- **WHEN** user has collapsed TIME_LIMIT usage details
- **AND** user presses 'r' to refresh
- **THEN** after refresh completes, usage details remain collapsed
