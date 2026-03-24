## ADDED Requirements

### Requirement: API key input must be masked
The system SHALL mask the API key input field to prevent visual exposure of the secret.

#### Scenario: User enters API key
- **WHEN** user types characters into the API key field
- **THEN** each character SHALL be displayed as `•` (bullet character)
- **AND** the actual key value SHALL be stored securely in memory

### Requirement: User must choose append or overwrite for existing config
When the configuration file already exists, the system SHALL prompt the user to choose how to handle the existing configuration.

#### Scenario: Config file exists
- **WHEN** user submits an API key AND `~/.zai-quota.yaml` already exists
- **THEN** system SHALL display: "Existing config found. [A]ppend key [O]verwrite all?"
- **AND** wait for user input

#### Scenario: User chooses append
- **WHEN** user presses `A` at the merge choice prompt
- **THEN** system SHALL read existing configuration
- **AND** merge the new `api_key` value into existing config
- **AND** preserve all other existing values
- **AND** write merged config to file

#### Scenario: User chooses overwrite
- **WHEN** user presses `O` at the merge choice prompt
- **THEN** system SHALL write a new config with defaults plus the entered API key
- **AND** replace the entire existing config file

### Requirement: Save failure must show clear recovery instructions
When configuration cannot be saved, the system SHALL display a clear error message with recovery instructions.

#### Scenario: SaveConfig fails
- **WHEN** `SaveConfig()` returns an error
- **THEN** system SHALL display the error message
- **AND** display: "Could not save configuration. Please restart the TUI to retry."
- **AND** exit after 3 seconds or on any key press

## MODIFIED Requirements

### Requirement: TUI must have setup-related states
The system SHALL support the following states for setup flow:
- `stateSetup`: Active text input for API key entry
- `stateSetupMergeChoice`: Prompting user to choose append or overwrite (NEW)
- `stateSetupClosing`: Displaying closing message before exit

#### Scenario: User submits key and file exists
- **WHEN** user presses Enter in `stateSetup` AND config file exists
- **THEN** system SHALL transition to `stateSetupMergeChoice`

#### Scenario: User makes merge choice
- **WHEN** user presses `A` or `O` in `stateSetupMergeChoice`
- **THEN** system SHALL proceed to save configuration
- **AND** transition to `stateLoading` on success

## REMOVED Requirements

### Requirement: Allow session to proceed on save failure
**Reason**: MVP simplicity - if config can't be saved, user should fix the issue and restart rather than proceeding with an unsaved key.
**Migration**: Always exit on save failure with clear instructions to restart.
