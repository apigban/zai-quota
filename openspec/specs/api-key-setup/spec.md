# Specification: API Key Setup in TUI

## 1. Overview
The TUI must detect the absence of a configured API key upon startup. If no API key is found (via flags, environment variables, or configuration file), the TUI will enter a "Setup Mode" to prompt the user for their API key. Once provided, the key will be persisted to the local configuration file (`~/.zai-quota.yaml`) for future use.

## 2. User Flow

### 2.1 First-Time Setup
1. **Launch**: User executes `zai-quota`.
2. **Detection**: The application checks `Config.APIKey`. If empty, it enters `stateSetup`.
3. **Prompt**: The TUI displays a screen:
   ```
   Welcome to Z.AI Quota Monitor!
   
   To get started, please provide your API Key.
   You can find this in your z.ai dashboard.
   
   API Key: [________________________________]
   
   (Enter to confirm, Esc to quit)
   ```
4. **Submission**:
   - User enters the key and presses `Enter`.
   - The application saves the key to `~/.zai-quota.yaml`.
   - The TUI transitions to `stateLoading` and proceeds with the normal flow.
5. **Cancellation**:
   - User presses `Esc` or `Ctrl+C`.
   - The TUI displays: "API Key is required to use this tool. Closing..."
   - The application exits after a brief delay or immediately.

### 2.2 Subsequent Launches
1. **Launch**: User executes `zai-quota`.
2. **Detection**: `Config.APIKey` is found (from file).
3. **Normal Flow**: The TUI enters `stateLoading` immediately and fetches quota data.

## 3. Technical Requirements

### 3.1 Configuration Persistence (`internal/config`)
- **File**: `~/.zai-quota.yaml` (defined as `configFileName` + `configFileExt` in `loader.go`).
- **Functionality**:
  - Implement `SaveConfig(cfg *Config) error` in `internal/config/loader.go`.
  - Use `yaml.Marshal` to serialize the `Config` struct.
  - Write to the home directory path.
  - Create the file with appropriate permissions (e.g., `0600` since it contains a secret).

### 3.2 TUI Enhancements (`internal/tui`)
- **States**:
  - Add `stateSetup` and `stateSetupClosing` to the `state` enum.
- **Model**:
  - Add `textinput.Model` from `charm.sh/bubbles/textinput` to the `Model` struct.
- **Initialization**:
  - In `NewModel`, if `cfg.APIKey == ""`, set `state: stateSetup`.
  - Initialize the text input model.
- **Update Logic**:
  - In `Update`, if `state == stateSetup`:
    - Handle `textinput.Msg`.
    - Handle `Enter`:
      - Update `m.cfg.APIKey`.
      - Re-initialize `m.client` with the new key.
      - Call `config.SaveConfig(m.cfg)`.
      - Transition to `stateLoading` and return `fetchQuotaCmd`.
    - Handle `Esc`:
      - Set `state: stateSetupClosing`.
      - Return a command to quit after a short delay.
- **View Logic**:
  - Render a setup-specific view when in `stateSetup`.
  - Show a clear "Closing..." message when in `stateSetupClosing`.

### 3.3 Error Handling
- If `SaveConfig` fails, the TUI should display an error message but might still allow the current session to proceed if the key is valid in memory. However, for a "setup" flow, persistence is the primary goal.

## 4. Verification Plan

### 4.1 Automated Tests
- **Config**: Test `SaveConfig` to ensure it correctly writes YAML to the expected path.
- **TUI**: Mock the `APIKey` being empty and verify the initial state is `stateSetup`.

### 4.2 Manual Verification
1. Remove `~/.zai-quota.yaml` and unset `ZAI_API_KEY`.
2. Run `zai-quota`.
3. Verify the setup prompt appears.
4. Enter a valid (or dummy) API key and press Enter.
5. Verify `~/.zai-quota.yaml` exists and contains the key.
6. Verify the TUI proceeds to fetch quota (will show auth error if dummy, which is expected).
7. Restart `zai-quota` and verify the setup prompt does *not* appear.
