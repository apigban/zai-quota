# Spec: multi-limit-display

## ADDED Requirements

### Requirement: Display TIME_LIMIT quota

The TUI SHALL display TIME_LIMIT quota with:
- Label: "[5-Hour Prompt Limit]"
- Progress bar colored based on percentage (Green <80%, Amber 80-89%, Orange 90-94%, Red ≥95%)
- Usage display showing "X / Y" format where X is currentValue and Y is usage (total)
- Remaining count
- Time until reset

#### Scenario: TIME_LIMIT at low usage
- **WHEN** TIME_LIMIT has currentValue=8, usage=1000, remaining=992, percentage=1, nextResetTime=4 hours from now
- **THEN** TUI displays:
  - Label: "[5-Hour Prompt Limit]"
  - Progress bar with 1% filled in green
  - "8 / 1,000"
  - "Remaining: 992"
  - "Reset: 4h 0m"

#### Scenario: TIME_LIMIT at warning level
- **WHEN** TIME_LIMIT has percentage=85
- **THEN** progress bar displays in amber color

#### Scenario: TIME_LIMIT at critical level
- **WHEN** TIME_LIMIT has percentage=92
- **THEN** progress bar displays in orange color

#### Scenario: TIME_LIMIT at emergency level
- **WHEN** TIME_LIMIT has percentage=97
- **THEN** progress bar displays in red color

### Requirement: Display TOKENS_LIMIT quota

The TUI SHALL display TOKENS_LIMIT quota with:
- Label: "[Weekly Quota Limit]"
- Progress bar colored based on percentage (same thresholds as TIME_LIMIT)
- Percentage display only (no usage/remaining counts)
- Time until reset

#### Scenario: TOKENS_LIMIT at moderate usage
- **WHEN** TOKENS_LIMIT has percentage=33, nextResetTime=5 days from now
- **THEN** TUI displays:
  - Label: "[Weekly Quota Limit]"
  - Progress bar with 33% filled in green
  - "33%"
  - "Reset: 5 days"

#### Scenario: TOKENS_LIMIT with percentage only
- **WHEN** TOKENS_LIMIT is the only data available
- **THEN** TUI does NOT display usage, total, or remaining fields

### Requirement: Display both limits simultaneously

The TUI SHALL display both TIME_LIMIT and TOKENS_LIMIT when both are present in the API response.

#### Scenario: Both limits present
- **WHEN** API returns both TIME_LIMIT and TOKENS_LIMIT
- **THEN** TUI displays TIME_LIMIT first, then TOKENS_LIMIT below it
- **AND** each limit has its own section with proper spacing

#### Scenario: Only TIME_LIMIT present
- **WHEN** API returns only TIME_LIMIT
- **THEN** TUI displays only TIME_LIMIT section

#### Scenario: Only TOKENS_LIMIT present
- **WHEN** API returns only TOKENS_LIMIT
- **THEN** TUI displays only TOKENS_LIMIT section
