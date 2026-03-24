## Why

The CLI tool has two fully-implemented, tested formatters (`FormatHuman`, `FormatColored`) that are never called. When run without `--json` or `--yaml`, the tool always launches the TUI—even in non-TTY environments (scripts, cron jobs, pipes), which causes failures. Users need a plain-text output mode for quick checks and automation.

## What Changes

- Add `--text` flag for explicit plain-text output
- Add TTY auto-detection: non-TTY environments default to plain-text instead of TUI
- Wire up the existing `FormatHuman` function
- Delete unused `FormatColored` function and its tests

## Capabilities

### New Capabilities

- `plain-text-output`: Support for human-readable plain-text output via `--text` flag or automatic non-TTY detection

### Modified Capabilities

- None

## Impact

- `cmd/zai-quota/root.go` - add flag, TTY detection, wire FormatHuman
- `internal/formatter/colored.go` - DELETE
- `internal/formatter/colored_test.go` - DELETE
- Behavior change: `zai-quota` in non-TTY now outputs plain text instead of attempting TUI
