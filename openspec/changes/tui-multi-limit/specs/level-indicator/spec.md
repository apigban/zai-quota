# Spec: level-indicator

## ADDED Requirements

### Requirement: Display subscription level in title bar

The TUI SHALL display the user's subscription level in the title bar, capitalized and enclosed in brackets.

#### Scenario: Pro subscription level
- **WHEN** API returns level="pro"
- **THEN** title bar displays "Z.AI Quota Monitor [Pro]"

#### Scenario: Lite subscription level
- **WHEN** API returns level="lite"
- **THEN** title bar displays "Z.AI Quota Monitor [Lite]"

#### Scenario: Max subscription level
- **WHEN** API returns level="max"
- **THEN** title bar displays "Z.AI Quota Monitor [Max]"

#### Scenario: Unknown subscription level
- **WHEN** API returns level="" (empty string) or unrecognized value
- **THEN** title bar displays "Z.AI Quota Monitor" (no level indicator)

### Requirement: Level capitalization

The TUI SHALL capitalize the first letter of the level value for display.

#### Scenario: Level value is lowercase
- **WHEN** API returns level="pro"
- **THEN** TUI capitalizes to display "[Pro]"

#### Scenario: Level value is already capitalized
- **WHEN** API returns level="Pro"
- **THEN** TUI displays "[Pro]" (no double capitalization)

#### Scenario: Level value is mixed case
- **WHEN** API returns level="pRo"
- **THEN** TUI normalizes to display "[Pro]"

### Requirement: Level styling

The TUI SHALL style the level indicator distinctly from the title text using gold color.

#### Scenario: Level visual distinction
- **WHEN** level is displayed in title bar
- **THEN** level text is styled with gold color (#FFD700)
- **AND** level text is bold
- **AND** level is separated from title by a space

### Requirement: Level persists across refreshes

The TUI SHALL maintain level display across data refreshes.

#### Scenario: Level remains after refresh
- **WHEN** user has viewed quota with level="pro"
- **AND** user presses 'r' to refresh
- **THEN** after refresh completes, level indicator still shows "[Pro]"

#### Scenario: Level updates if changed
- **WHEN** user's subscription changes from "pro" to "max"
- **AND** user refreshes the TUI
- **THEN** level indicator updates to show "[Max]"
