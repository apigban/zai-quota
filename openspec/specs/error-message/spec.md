### Requirement: Missing API key error shows all configuration options

When the API key is not configured, the error message SHALL inform users of both valid configuration methods: environment variable and config file.

#### Scenario: User runs without any API key configured
- **WHEN** user runs `zai-quota` without `ZAI_API_KEY` set and without `~/.zai-quota.yaml` containing an api_key
- **THEN** the error message SHALL mention both the `ZAI_API_KEY` environment variable AND the `~/.zai-quota.yaml` config file option
- **AND** the error message SHALL include the config key name (`api_key`)

#### Scenario: Error message is actionable
- **WHEN** user sees the missing API key error
- **THEN** the message SHALL provide enough detail for the user to fix the issue without consulting documentation
- **AND** the message SHALL accurately reflect the configuration precedence documented in README.md
