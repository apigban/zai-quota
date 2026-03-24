## Requirements

### Requirement: Timezone detection from system

The system SHALL detect the local timezone from the system configuration.

#### Scenario: TZ environment variable set
- **WHEN** the `TZ` environment variable is set to a valid IANA timezone (e.g., "America/New_York")
- **THEN** the system SHALL use that timezone

#### Scenario: /etc/timezone file exists
- **WHEN** the `TZ` environment variable is not set
- **AND** `/etc/timezone` file exists with a valid IANA timezone
- **THEN** the system SHALL use that timezone

#### Scenario: /etc/localtime symlink resolves
- **WHEN** neither `TZ` nor `/etc/timezone` are available
- **AND** `/etc/localtime` is a symlink to a zoneinfo file
- **THEN** the system SHALL extract the IANA timezone from the symlink path

#### Scenario: Timezone detection fails
- **WHEN** all detection methods fail
- **THEN** the system SHALL fall back to UTC

### Requirement: Timezone display format

The system SHALL format timezone information as numeric offset plus city name.

#### Scenario: Standard timezone display
- **WHEN** timezone is detected
- **THEN** the display SHALL show offset and city (e.g., "+04:00 Dubai")

#### Scenario: UTC fallback display
- **WHEN** timezone detection fails
- **THEN** the display SHALL show "UTC" only

### Requirement: City name extraction from IANA

The system SHALL extract a human-readable city name from IANA timezone identifiers.

#### Scenario: City with underscore
- **WHEN** IANA timezone is "America/New_York"
- **THEN** the city name SHALL be "New York"

#### Scenario: Simple city name
- **WHEN** IANA timezone is "Asia/Dubai"
- **THEN** the city name SHALL be "Dubai"

#### Scenario: UTC special case
- **WHEN** timezone is "UTC"
- **THEN** the display SHALL be "UTC" (no city extraction)
