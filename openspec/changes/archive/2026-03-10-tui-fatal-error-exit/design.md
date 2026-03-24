## Context

The TUI currently handles all errors uniformly: display a red error box and remain in `stateError`, requiring users to manually press 'q' to quit. This is appropriate for transient errors (network timeouts, 5xx server errors) but traps users when the error is unrecoverable (401 Unauthorized, 403 Forbidden - bad API key).

The error classification system already exists in `internal/errors/errors.go` with exit codes:
- `ExitAuth (3)` - Authentication errors
- `ExitNetwork (2)` - Network/timeout errors  
- `ExitError (1)` - Generic errors

The API client (`internal/api/client.go`) already wraps errors with these codes.

## Goals / Non-Goals

**Goals:**
- Exit TUI automatically after fatal errors (401, 403) with proper exit code
- Show 3-second countdown so users can read the error message
- Allow immediate quit with 'q' during countdown
- Maintain existing retry behavior for recoverable errors

**Non-Goals:**
- Changing CLI mode error handling (already works correctly)
- Adding new error types or exit codes
- Retry logic for fatal errors

## Decisions

### D1: Fatal error definition
**Decision:** Only `ExitAuth` errors (401, 403) are fatal. Network and server errors remain recoverable.

**Rationale:** Auth errors require user action outside the app (fix config/env). Network issues might resolve themselves.

**Alternatives considered:**
- All errors fatal → poor UX for transient issues
- Config errors only → 401/403 already covers this

### D2: Countdown mechanism
**Decision:** 3-second countdown using `tea.Tick` with 1-second intervals.

**Rationale:** Gives users time to read error but doesn't trap them. Matches patterns in other CLI tools.

**Alternatives considered:**
- Immediate exit → user may miss the error message
- 5+ seconds → too long for simple "fix your API key" message
- No countdown, just message → user still trapped

### D3: Error return flow
**Decision:** Store fatal error in model, return from `tui.Run()` after `tea.Quit`, let `main.go` handle exit code via existing `ClassifyError()`.

**Rationale:** Maintains clean separation. Bubble Tea's `tea.Quit` doesn't carry error info, so we use model state.

```
handleQuotaFetched (fatal error)
  → m.state = stateFatalError
  → m.countdown = 3
  → countdownTickCmd()
  
countdownTickMsg
  → m.countdown--
  → if countdown == 0: m.fatalErr = m.err; tea.Quit()
  → else: countdownTickCmd()

tui.Run()
  → finalModel, _ := p.Run()
  → if finalModel.fatalErr != nil: return finalModel.fatalErr
```

### D4: Manual quit during countdown
**Decision:** 'q' or 'ctrl+c' during countdown sets `fatalErr` and quits immediately.

**Rationale:** User can skip waiting if they understood the error.

## Risks / Trade-offs

| Risk | Mitigation |
|------|------------|
| Countdown feels slow for users who read fast | Allow 'q' to skip |
| User might miss error if they look away | 3 seconds is reasonable; error also printed to stderr on exit |
| New state adds complexity to state machine | State is isolated; only entered from stateLoading on fatal error |

## Migration Plan

No migration needed - this is additive behavior. Deploy in single release.

## Open Questions

None. Design is complete based on exploration session.
