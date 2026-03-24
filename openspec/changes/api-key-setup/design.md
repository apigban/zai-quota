## Context

The TUI currently has no interactive way to configure the API key. Users must either:
1. Set `ZAI_API_KEY` environment variable
2. Manually create `~/.zai-quota.yaml`

This works but isn't discoverable or user-friendly for first-time users.

The existing spec at `openspec/specs/api-key-setup/spec.md` defines the core flow but needs enhancement for:
- Masked input (security)
- Existing config file handling (user choice)
- Clear failure messaging

## Goals / Non-Goals

**Goals:**
- Interactive setup flow when no API key is detected
- Secure masked input for API key entry
- Persist configuration to `~/.zai-quota.yaml` with proper permissions (`0600`)
- Handle existing config files with user choice (append vs overwrite)
- Clear error messaging and recovery path on save failures

**Non-Goals:**
- Key rotation or update flow (MVP only)
- Multiple profile support
- Validation of key format before API call

## Decisions

### 1. State Machine Design

Add three new states to `internal/tui/tui.go`:

```
┌─────────────────────────────────────────────────────────────────────────┐
│                     SETUP MODE STATE MACHINE                            │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                          │
│  ┌─────────────┐                                                        │
│  │ stateSetup  │  Text input active, key being entered                  │
│  └──────┬──────┘                                                        │
│         │                                                                │
│         │ Enter pressed                                                  │
│         ▼                                                                │
│  ┌──────────────────────┐  (if file exists)                             │
│  │ stateSetupMergeChoice │  [A]ppend / [O]verwrite prompt               │
│  └──────────┬───────────┘                                                │
│             │                                                            │
│             │ A or O pressed                                             │
│             ▼                                                            │
│  ┌─────────────────────┐                                                │
│  │ SaveConfig() call   │                                                │
│  └──────────┬──────────┘                                                │
│             │                                                            │
│       ┌─────┴─────┐                                                     │
│       │           │                                                     │
│    Success      Failure                                                 │
│       │           │                                                     │
│       ▼           ▼                                                     │
│  ┌──────────┐  ┌───────────────────┐                                    │
│  │stateLoad │  │ stateSetupClosing │  Show error + "restart to retry"   │
│  │  ing     │  └───────────────────┘                                    │
│  └──────────┘                                                            │
│                                                                          │
│  ┌────────────────────┐                                                 │
│  │ stateSetupClosing  │  Also reached via Esc in stateSetup             │
│  └────────────────────┘                                                 │
│                                                                          │
└─────────────────────────────────────────────────────────────────────────┘
```

### 2. Input Masking

Use `charm.sh/bubbles/textinput` with `EchoMode` set to `EchoPassword` (or equivalent):
- Characters displayed as `•` 
- Actual value stored securely in model

### 3. Config Persistence

Add `SaveConfig(cfg *Config) error` to `internal/config/loader.go`:

```go
func SaveConfig(cfg *Config) error {
    homeDir, err := os.UserHomeDir()
    if err != nil {
        return err
    }
    
    configPath := filepath.Join(homeDir, configFileName+"."+configFileExt)
    
    data, err := yaml.Marshal(cfg)
    if err != nil {
        return err
    }
    
    // Write with restrictive permissions (secret file)
    return os.WriteFile(configPath, data, 0600)
}
```

### 4. Existing Config Handling

When `~/.zai-quota.yaml` exists:
- Show prompt: "Existing config found. [A]ppend key [O]verwrite all?"
- **Append**: Read existing config, merge `api_key` field, write back
- **Overwrite**: Write new config with defaults + entered key

Implementation uses `LoadConfig` logic to read existing, then `SaveConfig` to write.

### 5. Failure Handling

On `SaveConfig` failure:
1. Display error message in setup view
2. Show: "Could not save configuration. Please restart the TUI to retry."
3. Transition to `stateSetupClosing`
4. Exit after 3 seconds or on any key press

## Risks / Trade-offs

| Risk | Mitigation |
|------|------------|
| File permission issues on save | Clear error message with restart instruction |
| User loses other config on overwrite | Explicit prompt before overwrite |
| Race condition if config modified externally | Accept limitation; user can restart |
| API key leaked in process list | Input masking helps; env var still visible in ps |
