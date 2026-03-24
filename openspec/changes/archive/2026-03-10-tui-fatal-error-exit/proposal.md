## Why

The TUI traps users when fatal configuration errors occur (bad API key). Users must manually press 'q' to quit even though the error is unrecoverable without fixing their config. This creates friction for a common setup mistake.

## What Changes

- Add fatal error detection that distinguishes between recoverable (network, 5xx) and fatal (auth) errors
- Add 3-second countdown with visual feedback before auto-exit on fatal errors
- Allow users to press 'q' to skip countdown and exit immediately
- Return proper exit code (ExitAuth=3) for fatal auth errors

## Capabilities

### New Capabilities
- `fatal-error-exit`: Defines behavior for fatal vs recoverable errors in TUI, countdown display, and exit handling

### Modified Capabilities
- None. This is additive behavior for the TUI only.

## Impact

- `internal/tui/tui.go` - Add state machine for fatal errors, countdown tick handling
- `internal/errors/errors.go` - Already has `ExitAuth` code, no changes needed
- `internal/api/client.go` - Already returns `QuotaError` with auth codes, no changes needed
