## Context

The `zai-quota` CLI currently defaults to launching a Bubble Tea TUI when no output format flag is provided. This breaks in non-interactive environments (scripts, cron jobs, CI/CD pipelines) where no TTY is available. Additionally, two fully-implemented formatters (`FormatHuman`, `FormatColored`) exist but are never invoked.

## Goals / Non-Goals

**Goals:**
- Enable non-interactive usage without TUI failures
- Expose the existing `FormatHuman` output path
- Clean up dead code (`FormatColored`)

**Non-Goals:**
- Changing TUI behavior for interactive users
- Adding new formatting features
- Modifying API or data processing logic

## Decisions

### 1. TTY Detection Strategy

**Decision:** Use `golang.org/x/term.IsTerminal(int(os.Stdout.Fd()))` to detect TTY.

**Rationale:** The dependency is already available (via `github.com/charmbracelet/x/term`). This is the standard Go approach for TTY detection.

**Alternatives considered:**
- Checking `os.Stdin.Stat()` for character device - less reliable across platforms

### 2. Flag Naming

**Decision:** Use `--text` as the explicit plain-text flag.

**Rationale:** Short, intuitive, parallels `--json`/`--yaml` pattern. Common in other CLIs.

### 3. Colored Formatter

**Decision:** Delete `FormatColored` and its tests.

**Rationale:** It provides "TUI-like output but non-interactive" which has no clear use case. If users want colors, they get the TUI. If they want scriptable output, they get plain text or JSON/YAML.

### 4. Priority Logic

**Decision:** Explicit flags override TTY detection.

```
--json → JSON (always)
--yaml → YAML (always)
--text → plain text (always)
(no flag, TTY) → TUI
(no flag, non-TTY) → plain text
```

**Rationale:** Predictable behavior. Users can force plain text in a TTY with `--text`.

## Risks / Trade-offs

| Risk | Mitigation |
|------|------------|
| Breaking users who pipe TUI output | Non-TTY now outputs plain text, which is more useful than broken TUI |
| `golang.org/x/term` platform edge cases | Well-tested library; fallback to plain text on error |
