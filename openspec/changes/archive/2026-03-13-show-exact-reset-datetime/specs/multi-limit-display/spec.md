## MODIFIED Requirements

### Requirement: Reset time display format

The system SHALL display reset times with duration, exact datetime, and timezone.

#### Scenario: Reset display with detected timezone
- **WHEN** the TUI displays a quota reset time
- **AND** timezone is detected
- **THEN** the display SHALL show: `Reset: {duration} ({date}, {time} {offset} {city})`
- **AND** date format SHALL be `Mon D` (e.g., "Mar 13")
- **AND** time format SHALL be 24-hour `HH:MM` (e.g., "20:00")
- **AND** offset format SHALL be `±HH:MM` (e.g., "+04:00")

Example: `Reset: 4h 6m (Mar 13, 20:00 +04:00 Dubai)`

#### Scenario: Reset display with UTC fallback
- **WHEN** the TUI displays a quota reset time
- **AND** timezone detection failed
- **THEN** the display SHALL show: `Reset: {duration} ({date}, {time} UTC)`

Example: `Reset: 4h 6m (Mar 13, 16:00 UTC)`

#### Scenario: Multi-day reset display
- **WHEN** the reset is more than 24 hours away
- **THEN** the date portion SHALL reflect the actual reset date

Example: `Reset: 20 days (Apr 2, 00:00 +04:00 Dubai)`

#### Scenario: Both limit types show enhanced reset
- **WHEN** both `TOKENS_LIMIT` and `TIME_LIMIT` are displayed
- **THEN** both SHALL use the enhanced reset format with timezone
