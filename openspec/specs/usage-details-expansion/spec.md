## Requirements

### Requirement: Expandable usage details for all quota types

The system SHALL allow users to expand/collapse usage details for both `TOKENS_LIMIT` and `TIME_LIMIT` quota types.

#### Scenario: Expand TOKENS_LIMIT details
- **GIVEN** `TOKENS_LIMIT` has `UsageDetails` with at least one entry
- **WHEN** user presses 'e' key
- **THEN** `TOKENS_LIMIT` usage breakdown is displayed
- **AND** each entry shows model code and usage count

#### Scenario: Expand TIME_LIMIT details
- **GIVEN** `TIME_LIMIT` has `UsageDetails` with at least one entry
- **WHEN** user presses 'e' key
- **THEN** `TIME_LIMIT` usage breakdown is displayed
- **AND** each entry shows tool code and usage count

#### Scenario: Toggle both quota details independently
- **GIVEN** both quota types have usage details
- **WHEN** user presses 'e' key
- **THEN** both `TOKENS_LIMIT` and `TIME_LIMIT` expand/collapse together

### Requirement: Fallback display for missing total/remaining

The system SHALL display percentage-based information when `total` or `remaining` fields are not available.

#### Scenario: Display percentage when total is zero
- **GIVEN** a quota has `total` of 0 or less
- **WHEN** the quota is rendered
- **THEN** system displays "X% used" instead of "X / Y"

#### Scenario: Hide remaining when zero or negative
- **GIVEN** a quota has `remaining` of 0 or less
- **WHEN** the quota is rendered
- **THEN** the "Remaining: X" line is not displayed
