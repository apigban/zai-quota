## 1. Config Persistence

- [x] 1.1 Add `SaveConfig(cfg *Config) error` function to `internal/config/loader.go`
- [x] 1.2 Use `yaml.Marshal` to serialize Config struct
- [x] 1.3 Write to `~/.zai-quota.yaml` with `0600` permissions
- [x] 1.4 Add unit tests for `SaveConfig` in `internal/config/loader_test.go`

## 2. TUI State Machine

- [x] 2.1 Add `stateSetup`, `stateSetupMergeChoice`, `stateSetupClosing` constants to state enum in `internal/tui/tui.go`
- [x] 2.2 Add `textinput.Model` field to Model struct (import from `charm.sh/bubbles/textinput`)
- [x] 2.3 Add `existingConfigExists` bool field to Model struct for tracking merge choice flow
- [x] 2.4 Add `saveError` error field to Model struct for displaying save failures

- [x] 2.5 Add styles for setup UI in `internal/tui/styles.go` (setupBoxStyle, setupTitleStyle, etc.)

## 3. TUI Initialization

- [x] 3.1 In `NewModel`, check if `cfg.APIKey == ""` and set initial state to `stateSetup`
- [x] 3.2 Initialize textinput.Model with password echo mode (masked input)
- [x] 3.3 Focus the text input on initialization
- [x] 3.4 Check if `~/.zai-quota.yaml` exists to set `existingConfigExists` flag

- [x] 3.5 Add unit test: `NewModel` returns `stateSetup` when API key is empty
- [x] 3.6 Add unit test: `NewModel` returns `stateLoading` when API key is set

- [x] 4.1 Handle `textinput.Msg` when in `stateSetup`
- [x] 4.2 Handle `Enter` key in `stateSetup` - check if file exists, transition to `stateSetupMergeChoice` or proceed to save
- [x] 4.3 Handle `Esc` key in `stateSetup` - transition to `stateSetupClosing`
- [x] 4.4 Handle `A` key in `stateSetupMergeChoice` - load existing config, merge api_key, save, transition to `stateLoading`
- [x] 4.5 Handle `O` key in `stateSetupMergeChoice` - create new config with defaults + key, save, transition to `stateLoading`
- [x] 4.6 On save failure, set `saveError` and transition to `stateSetupClosing`
- [x] 4.7 On save success, re-initialize client with new API key and return `fetchQuotaCmd`
- [x] 4.8 On save failure, set `saveError` and transition to `stateSetupClosing`

- [x] 4.9 On save success, re-initialize client with new API key and return `fetchQuotaCmd`
- [x] 4.10 Handle `textinput.Msg` when in `stateSetup`
- [x] 4.2 Handle `Enter` key in `stateSetup` - check if file exists, transition to `stateSetupMergeChoice` or proceed to save
- [x] 4.3 Handle `Esc` key in `stateSetup` - transition to `stateSetupClosing`
    - [x] 4.4 Handle `A` key in `stateSetupMergeChoice` - load existing config, merge api_key, save, transition to `stateLoading`
    - [x] 4.5 Handle `O` key in `stateSetupMergeChoice` - create new config with defaults + key, save, transition to `stateLoading`
    - [x] 4.6 On save failure, set `saveError` and transition to `stateSetupClosing`
    - [x] 4.7 On save success, re-initialize client with new API key and return `fetchQuotaCmd`
    - [x] 4.8 On save failure, set `saveError` and transition to `stateSetupClosing`
    - [x] 4.10 Handle `textinput.Msg` when in `stateSetup`
    - [x] 4.2 Handle `Enter` key in `stateSetup` - check if file exists, transition to `stateSetupMergeChoice` or proceed to save
    - [x] 4.3 Handle `Esc` key in `stateSetup` - transition to `stateSetupClosing`
    - [x] 4.4 Handle `A` key in `stateSetupMergeChoice` - load existing config, merge api_key, save, transition to `stateLoading`
    - [x] 4.5 Handle `O` key in `stateSetupMergeChoice` - create new config with defaults + key, save, transition to `stateLoading`
    - [x] 4.6 On save failure, set `saveError` and transition to `stateSetupClosing`
    - [x] 4.7 On save success, re-initialize client with new API key and return `fetchQuotaCmd`
    - [x] 4.8 On save failure, set `saveError` and transition to `stateSetupClosing`
    - [x] 4.9 Handle `textinput.Msg` when in `stateSetup`
    - [x] 4.10 Handle `Enter` key in `stateSetup` - check if file exists, transition to `stateSetupMergeChoice` or proceed to save
    - [x] 4.11 Handle `Esc` key in `stateSetup` - transition to `stateSetupClosing`
    - [x] 4.12 Handle `A` key in `stateSetupMergeChoice` - load existing config, merge api_key, save, transition to `stateLoading`
    - [x] 4.13 Handle `O` key in `stateSetupMergeChoice` - create new config with defaults + key, save, transition to `stateLoading`
    - [x] 4.14 On save failure, set `saveError` and transition to `stateSetupClosing`
    - [x] 4.15 On save success, re-initialize client with new API key and return `fetchQuotaCmd`
    - [x] 4.16 Handle `textinput.Msg` when in `stateSetup`
    - [x] 4.17 Handle `Enter` key in `stateSetup` - check if file exists, transition to `stateSetupMergeChoice` or proceed with save
    - [x] 4.18 Handle `Esc` key in `stateSetup` - transition to `stateSetupClosing`
    - [x] 4.19 Handle `A` key in `stateSetupMergeChoice` - load existing config, merge api_key, save, transition to `stateLoading`
    - [x] 4.20 Handle `O` key in `stateSetupMergeChoice` - create new config with defaults + key, save, transition to `stateLoading`
    - [x] 4.21 On save failure, set `saveError` and transition to `stateSetupClosing`
    - [x] 4.22 On save success, re-initialize client with new API key and return `fetchQuotaCmd`
    - [x] 4.23 Add `renderSetupView()` method for `stateSetup` - show welcome message, prompt, and masked input
    - [x] 5.2 Add `renderMergeChoiceView()` method for `stateSetupMergeChoice` - show append/overwrite prompt
    - [x] 5.3 Add `renderSetupClosingView()` method for `stateSetupClosing` - show error or cancellation message
    - [x] 5.4 Update `View()` to route to appropriate render method based on state
    - [x] 5.5 Add styles for setup UI in `internal/tui/styles.go` (setupBoxStyle, setupTitleStyle, etc.)

## 6. Testing

    - [x] 6.1 Add unit test: `SaveConfig` writes correct YAML
    - [x] 6.2 Add unit test: `SaveConfig` creates file with `0600` permissions
    - [x] 6.3 Add unit test: `NewModel` returns `stateSetup` when API key is empty
    - [x] 6.4 Add unit test: `NewModel` returns `stateLoading` when API key is set
    - [ ] 6.5 Manual verification: Remove config, run TUI, verify setup prompt appears
    - [ ] 6.6 Manual verification: Enter key, verify config saved and TUI proceeds
    - [ ] 6.7 Manual verification: Restart TUI, verify setup prompt does not appear
