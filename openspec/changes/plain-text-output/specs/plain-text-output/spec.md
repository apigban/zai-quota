## ADDED Requirements

### Requirement: Plain text output via --text flag
The system SHALL output human-readable plain text when the `--text` flag is provided.

#### Scenario: User requests plain text output
- **WHEN** user runs `zai-quota --text`
- **THEN** system outputs quota information in plain text format without ANSI codes
- **AND** system exits immediately after output

#### Scenario: Plain text output contains no ANSI codes
- **WHEN** user runs `zai-quota --text`
- **THEN** output contains no ANSI escape sequences
- **AND** output is suitable for piping to files or other commands

### Requirement: Automatic TTY detection
The system SHALL detect whether stdout is connected to a terminal and adjust output mode accordingly.

#### Scenario: Interactive terminal without format flags
- **WHEN** user runs `zai-quota` with stdout connected to a TTY
- **AND** no `--json`, `--yaml`, or `--text` flag is provided
- **THEN** system launches the interactive TUI

#### Scenario: Non-interactive environment without format flags
- **WHEN** user runs `zai-quota` with stdout NOT connected to a TTY (pipe, redirect, cron)
- **AND** no `--json`, `--yaml`, or `--text` flag is provided
- **THEN** system outputs plain text format
- **AND** system exits immediately after output

### Requirement: Explicit flags override TTY detection
The system SHALL honor explicit format flags regardless of TTY status.

#### Scenario: --text in TTY environment
- **WHEN** user runs `zai-quota --text` in an interactive terminal
- **THEN** system outputs plain text (not TUI)
- **AND** system exits immediately

#### Scenario: --json in non-TTY environment
- **WHEN** user runs `zai-quota --json` in a script
- **THEN** system outputs JSON format
- **AND** system exits immediately

### Requirement: Mutually exclusive format flags
The system SHALL reject combinations of format flags.

#### Scenario: User provides both --text and --json
- **WHEN** user runs `zai-quota --text --json`
- **THEN** system returns an error
- **AND** error message indicates flags are mutually exclusive

#### Scenario: User provides multiple format flags
- **WHEN** user runs `zai-quota --json --yaml --text`
- **THEN** system returns an error
- **AND** error message indicates only one format flag is allowed
