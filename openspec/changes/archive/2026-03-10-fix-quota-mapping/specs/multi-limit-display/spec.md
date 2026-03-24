## MODIFIED Requirements

### Requirement: Quota type labels match Z.AI website

The system SHALL display quota labels that match the Z.AI website:
- `TOKENS_LIMIT` → "[5-Hour Prompt Limit]"
- `TIME_LIMIT` → "[Tool Quota]"

#### Scenario: Display correct labels in TUI
- **WHEN** the TUI renders quota information
- **THEN** `TOKENS_LIMIT` displays as "[5-Hour Prompt Limit]"
- **AND** `TIME_LIMIT` displays as "[Tool Quota]"

#### Scenario: Display correct labels in CLI output
- **WHEN** the processor formats quota for CLI output (JSON, YAML, human)
- **THEN** `TOKENS_LIMIT` has label "[5-Hour Prompt Limit]"
- **AND** `TIME_LIMIT` has label "[Tool Quota]"

### Requirement: Quota display order prioritizes prompts

The system SHALL display the prompt quota (`TOKENS_LIMIT`) before the tool quota (`TIME_LIMIT`) in the TUI.

#### Scenario: TUI display order
- **WHEN** both `TOKENS_LIMIT` and `TIME_LIMIT` are present
- **THEN** `TOKENS_LIMIT` section appears first
- **AND** `TIME_LIMIT` section appears second
