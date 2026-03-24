## Why

First-time users have no way to configure their API key through the TUI - they must manually create `~/.zai-quota.yaml` or set the `ZAI_API_KEY` environment variable. This creates friction for new users and isn't discoverable.

## What Changes

- Add "Setup Mode" to TUI that activates when no API key is configured
- Add masked text input for entering API key securely
- Add `SaveConfig()` function to persist configuration to `~/.zai-quota.yaml`
- Add prompt for append vs overwrite when config file already exists
- Handle save failures gracefully with clear user messaging

## Capabilities

### New Capabilities

(None - this enhances existing capabilities)

### Modified Capabilities

- `api-key-setup`: Adding interactive setup flow, config persistence, masked input, and existing-file handling. Original spec covered detection and prompting; this adds the implementation requirements for persistence and user choice on file conflicts.

## Impact

- `internal/config/loader.go` - Add `SaveConfig()` function
- `internal/tui/tui.go` - Add `stateSetup`, `stateSetupMergeChoice`, `stateSetupClosing` states; add text input handling; add setup view rendering
- `internal/tui/styles.go` - May need additional styles for setup UI
- `~/.zai-quota.yaml` - Will be created with `0600` permissions if it doesn't exist
