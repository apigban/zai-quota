# Design: zai-quota-v2

## Architecture Overview

```
┌─────────────────────────────────────────────────────────────────┐
│  cmd/zai-quota/root.go                                          │
│  ├── Parse flags (--json, --yaml, --debug, --help, --version)   │
│  ├── If machine-readable flag → existing formatter path         │
│  └── Else → launch TUI                                          │
└─────────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────────┐
│  internal/api/client.go                                         │
│  ├── FetchQuota() returns *models.QuotaResponse                 │
│  └── Parses actual API structure (limits array, level)          │
└─────────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────────┐
│  internal/models/quota.go                                       │
│  ├── QuotaResponse { Limits []Limit, Level string }             │
│  ├── Limit { Type, Usage, CurrentValue, Remaining, ... }        │
│  └── UsageDetail { ModelCode, Usage }                           │
└─────────────────────────────────────────────────────────────────┘
                              │
          ┌───────────────────┼───────────────────┐
          ▼                   ▼                   ▼
┌─────────────────┐ ┌─────────────────┐ ┌─────────────────┐
│ internal/tui    │ │ internal/tui/   │ │ internal/       │
│ (tui.go)        │ │ styles.go       │ │ processor       │
│                 │ │                 │ │ (processor.go)  │
│ Multi-limit     │ │ Lipgloss defs   │ │ ProcessLimits() │
│ display         │ │ Level style     │ │                 │
└─────────────────┘ └─────────────────┘ └─────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────────┐
│  internal/formatter/                                            │
│  ├── human.go    - Plain text output                            │
│  ├── colored.go  - Colored terminal output                      │
│  ├── json.go     - JSON output                                  │
│  └── yaml.go     - YAML output                                  │
└─────────────────────────────────────────────────────────────────┘
```

## Data Models

### QuotaResponse

```go
type QuotaResponse struct {
    Limits []Limit `json:"limits"`
    Level  string  `json:"level"`
}
```

### Limit

```go
type Limit struct {
    Type          string        `json:"type"`
    Unit          int           `json:"unit"`
    Number        int           `json:"number"`
    Usage         int           `json:"usage,omitempty"`        // TIME_LIMIT: total allowed
    CurrentValue  int           `json:"currentValue,omitempty"` // TIME_LIMIT: actually used
    Remaining     int           `json:"remaining,omitempty"`
    Percentage    int           `json:"percentage"`
    NextResetTime int64         `json:"nextResetTime"`          // Milliseconds
    UsageDetails  []UsageDetail `json:"usageDetails,omitempty"` // TIME_LIMIT only
}
```

### UsageDetail

```go
type UsageDetail struct {
    ModelCode string `json:"modelCode"`
    Usage     int    `json:"usage"`
}
```

### Field Semantics by Limit Type

| Field | TIME_LIMIT | TOKENS_LIMIT |
|-------|------------|--------------|
| `type` | "TIME_LIMIT" | "TOKENS_LIMIT" |
| `unit` + `number` | 5×1 = 5 hours | 3×5 = 15 days (weekly) |
| `usage` | Total allowed prompts | Not present |
| `currentValue` | Actually used prompts | Not present |
| `remaining` | usage - currentValue | Not present |
| `percentage` | (currentValue / usage) × 100 | Pre-calculated by API |
| `nextResetTime` | Millisecond timestamp | Millisecond timestamp |
| `usageDetails` | Per-model breakdown | Not present |

## TUI Model

```go
type state int

const (
    stateLoading state = iota
    stateLoaded
    stateRefreshing
    stateError
)

type Model struct {
    state     state
    cfg       *config.Config
    client    *api.Client
    
    // API response data
    limits    []models.Limit
    level     string
    
    // Processed data for display
    processed map[string]ProcessedLimitData  // keyed by limit type
    
    // UI state
    expanded  map[string]bool  // keyed by limit type (for usageDetails)
    err       error
    width     int
    height    int
    debug     bool
    debugLog  []string
}

type ProcessedLimitData struct {
    Type          string
    Label         string  // "[5-Hour Prompt Limit]" or "[Weekly Quota Limit]"
    Percentage    int
    Total         int                     // TIME_LIMIT only
    Used          int                     // TIME_LIMIT only
    Remaining     int                     // TIME_LIMIT only
    ResetDisplay  string
    WarningLevel  string
    UsageDetails  []models.UsageDetail    // TIME_LIMIT only
}
```

## TUI State Machine

```
           ┌──────────────────────────────────────┐
           │                                      │
           ▼                                      │
┌─────────────┐     success      ┌─────────────┐ │
│   Loading   │ ──────────────▶ │   Loaded    │ │
└─────────────┘                 └─────────────┘ │
      │                               │         │
      │ fail                          │ 'r'     │
      ▼                               ▼         │
┌─────────────┐               ┌─────────────┐   │
│    Error    │◀──────────────│ Refreshing  │───┘
│  (empty)    │     fail      └─────────────┘
└─────────────┘
      │
      │ Retry success
      └──────────────▶ Loaded
```

## TUI Key Handling

```go
func (m Model) handleKeyPress(msg tea.KeyPressMsg) (tea.Model, tea.Cmd) {
    switch msg.String() {
    case "q", "Q", "ctrl+c":
        return m, tea.Quit
    case "r", "R":
        if m.state != stateRefreshing {
            m.state = stateRefreshing
            return m, fetchQuotaCmd(...)
        }
    case "e", "E":
        // Toggle expand for TIME_LIMIT if it has usageDetails
        if details, ok := m.processed["TIME_LIMIT"]; ok && len(details.UsageDetails) > 0 {
            if m.expanded == nil {
                m.expanded = make(map[string]bool)
            }
            m.expanded["TIME_LIMIT"] = !m.expanded["TIME_LIMIT"]
        }
    }
    return m, nil
}
```

## TUI View Rendering

### Normal State

```
┌──────────────────────────────────────────────────────────────────┐
│ Z.AI Quota Monitor [Pro]                                 14:32  │
├──────────────────────────────────────────────────────────────────┤
│                                                                  │
│  TIME LIMIT [5-Hour Prompt Limit]                                │
│  █░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░  0.8% │
│  8 / 1,000                                                       │
│  Remaining: 992                                                  │
│  Reset: 4h 23m                                                   │
│  [e] Expand details                                              │
│                                                                  │
│  TOKENS LIMIT [Weekly Quota Limit]                               │
│  ██████████████████████████████████░░░░░░░░░░░░░░░░░░░░  33%   │
│  Reset: 5 days                                                   │
│                                                                  │
│  [r] Refresh  [e] Expand  [q] Quit                               │
└──────────────────────────────────────────────────────────────────┘
```

### Expanded State (after pressing 'e')

```
│  TIME LIMIT [5-Hour Prompt Limit]                                │
│  █░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░  0.8% │
│  8 / 1,000                                                       │
│  Remaining: 992                                                  │
│  Reset: 4h 23m                                                   │
│  Usage breakdown:                                                │
│    └─ search-prime: 5                                            │
│    └─ model-alpha: 3                                             │
│  [e] Collapse                                                    │
```

### Refreshing State

Same as normal, but header shows "Refreshing..." instead of time.

### Error State (Overlay with Data)

```
┌──────────────────────────────────────────────────────────────────┐
│ Z.AI Quota Monitor [Pro]                                 14:33  │
├──────────────────────────────────────────────────────────────────┤
│ ┌──────────────────────────────────────────────────────────┐     │
│ │ Error: Connection timeout                                 │     │
│ │                                          [Dismiss] [Retry]│     │
│ └──────────────────────────────────────────────────────────┘     │
│  TIME LIMIT [5-Hour Prompt Limit]                                │
│  █░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░  0.8% │
│  8 / 1,000                                                       │
│  ...                                                             │
└──────────────────────────────────────────────────────────────────┘
```

### Empty State (Initial Load Failure)

```
┌──────────────────────────────────────────────────────────────────┐
│ Z.AI Quota Monitor                                               │
├──────────────────────────────────────────────────────────────────┤
│ ┌──────────────────────────────────────────────────────────┐     │
│ │ Error: Connection timeout                                 │     │
│ │                                          [Dismiss] [Retry]│     │
│ └──────────────────────────────────────────────────────────┘     │
│                                                                  │
│  No quota data available                                         │
│                                                                  │
│  [r] Retry  [q] Quit                                             │
└──────────────────────────────────────────────────────────────────┘
```

## Labels

| Limit Type | Label |
|------------|-------|
| TIME_LIMIT | `[5-Hour Prompt Limit]` |
| TOKENS_LIMIT | `[Weekly Quota Limit]` |

## Level Display

Level is capitalized and shown in title bar:
- `"pro"` → `[Pro]`
- `"lite"` → `[Lite]`
- `"max"` → `[Max]`

## Progress Bar Colors

| Percentage | Color | Level |
|------------|-------|-------|
| 0-79% | Green | Safe |
| 80-89% | Amber | Warning |
| 90-94% | Orange | Critical |
| 95-100% | Red | Emergency |

## Styling (styles.go)

```go
var (
    colorSafe     = lipgloss.Color("#4CAF50")  // Green
    colorWarning  = lipgloss.Color("#FFC107")  // Amber
    colorCritical = lipgloss.Color("#FF5722")  // Orange
    colorEmergency = lipgloss.Color("#F44336") // Red
    colorPurple   = lipgloss.Color("#7c3aed")
    colorGray     = lipgloss.Color("#666666")
    colorGold     = lipgloss.Color("#FFD700")
    
    titleStyle = lipgloss.NewStyle().
        Bold(true).
        Foreground(colorPurple).
        Padding(0, 1)
    
    levelStyle = lipgloss.NewStyle().
        Foreground(colorGold).
        Bold(true)
    
    errorBoxStyle = lipgloss.NewStyle().
        Border(lipgloss.RoundedBorder()).
        BorderForeground(lipgloss.Color("#ff0000")).
        Padding(0, 1)
    
    helpStyle = lipgloss.NewStyle().
        Foreground(colorGray)
)
```

## JSON Output Structure

```json
{
  "limits": [
    {
      "type": "TIME_LIMIT",
      "label": "5-Hour Prompt Limit",
      "percentage": 1,
      "current_value": 8,
      "total": 1000,
      "remaining": 992,
      "next_reset": "2025-03-09T18:04:00Z",
      "next_reset_local": "2025-03-09 18:04",
      "usage_details": [
        {"model_code": "search-prime", "usage": 5},
        {"model_code": "model-alpha", "usage": 3}
      ]
    },
    {
      "type": "TOKENS_LIMIT",
      "label": "Weekly Quota Limit",
      "percentage": 33,
      "next_reset": "2025-03-16T12:00:00Z",
      "next_reset_local": "2025-03-16 12:00"
    }
  ],
  "level": "pro"
}
```

## Human Output Format

```
[5-Hour Prompt Limit]
Type: TIME_LIMIT
Usage: 8 / 1,000
Remaining: 992
Next Reset: 2025-03-09 18:04

[Weekly Quota Limit]
Type: TOKENS_LIMIT
Usage: 33%
Next Reset: 2025-03-16 12:00
```

## File Changes Summary

| File | Change Type |
|------|-------------|
| `internal/models/quota.go` | **Rewrite** - New struct definitions |
| `internal/processor/processor.go` | **Update** - Handle new fields |
| `internal/formatter/human.go` | **Update** - New display format |
| `internal/formatter/colored.go` | **Update** - New display format |
| `internal/formatter/json.go` | **Update** - New output structure |
| `internal/formatter/yaml.go` | **Update** - New output structure |
| `internal/tui/tui.go` | **Major** - Multi-limit, expand/collapse |
| `internal/tui/styles.go` | **Minor** - Add level style |
| `cmd/zai-quota/root.go` | **Simplify** - Use quota.Limits directly |
| All test files | **Update** - New fixtures |

## Error Handling

| Error Type | Overlay Message | Source |
|------------|-----------------|--------|
| Network timeout | "Connection timeout" | `apierrors.IsNetworkError()` |
| Auth failure | "Authentication failed" | `apierrors.IsAuthError()` |
| Server error | "Server error" | `apierrors.IsGenericError()` |
| Unknown | "Failed to fetch quota" | Fallback |

## Testing Strategy

1. **Unit tests for models** - JSON marshaling/unmarshaling
2. **Unit tests for processor** - Field mapping validation
3. **Unit tests for formatters** - Output format verification
4. **Unit tests for TUI** - State transitions, key handling
5. **Integration tests** - Full lifecycle with mocked API
6. **Manual verification** - Test with real API
