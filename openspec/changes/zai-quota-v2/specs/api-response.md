# Spec: API Response Structure

## Capability

Correct parsing of Z.ai quota API response with multiple limit types and subscription level.

## API Endpoint

```
GET /v1/quota
Authorization: Bearer <api_key>
```

## Response Structure

```json
{
  "code": 200,
  "msg": "Operation successful",
  "data": {
    "limits": [
      {
        "type": "TIME_LIMIT",
        "unit": 5,
        "number": 1,
        "usage": 1000,
        "currentValue": 8,
        "remaining": 992,
        "percentage": 1,
        "nextResetTime": 1775186469998,
        "usageDetails": [
          {"modelCode": "search-prime", "usage": 5},
          {"modelCode": "model-alpha", "usage": 3}
        ]
      },
      {
        "type": "TOKENS_LIMIT",
        "unit": 3,
        "number": 5,
        "percentage": 33,
        "nextResetTime": 1773052446291
      }
    ],
    "level": "pro"
  },
  "success": true
}
```

## Requirements

### REQ-API-001: Response Validation

The client SHALL validate:
- `success` field is `true`
- `code` field is present
- `data` field is not null
- `data.limits` array is present and not empty

### REQ-API-002: Limit Types

The API returns two limit types:

| Type | Description | Reset Period |
|------|-------------|--------------|
| TIME_LIMIT | Prompt/request count | 5 hours |
| TOKENS_LIMIT | Token consumption percentage | 7 days |

### REQ-API-003: TIME_LIMIT Fields

TIME_LIMIT response includes:

| Field | Type | Description |
|-------|------|-------------|
| `type` | string | "TIME_LIMIT" |
| `unit` | int | Time unit multiplier (5) |
| `number` | int | Number of units (1) |
| `usage` | int | Total allowed prompts |
| `currentValue` | int | Actually used prompts |
| `remaining` | int | usage - currentValue |
| `percentage` | int | (currentValue / usage) × 100 |
| `nextResetTime` | int64 | Reset time in milliseconds |
| `usageDetails` | array | Per-model usage breakdown |

**Important:** `usage` is the TOTAL ALLOWED, not the amount used. `currentValue` is the amount actually used.

### REQ-API-004: TOKENS_LIMIT Fields

TOKENS_LIMIT response includes:

| Field | Type | Description |
|-------|------|-------------|
| `type` | string | "TOKENS_LIMIT" |
| `unit` | int | Time unit multiplier (3) |
| `number` | int | Number of units (5) |
| `percentage` | int | Pre-calculated usage percentage |
| `nextResetTime` | int64 | Reset time in milliseconds |

TOKENS_LIMIT does NOT include: `usage`, `currentValue`, `remaining`, `usageDetails`

### REQ-API-005: Level Field

The `data.level` field indicates subscription tier:

| Level | 5-Hour Prompts | Weekly Prompts |
|-------|----------------|----------------|
| `lite` | ~80 | ~400 |
| `pro` | ~400 | ~2,000 |
| `max` | ~1,600 | ~8,000 |

### REQ-API-006: Timestamp Format

`nextResetTime` is a Unix timestamp in **milliseconds**, not seconds.

Conversion: `time.Unix(nextResetTime/1000, 0)`

### REQ-API-007: Error Handling

The client SHALL handle:
- HTTP 401: Authentication error
- HTTP 403: Authorization error
- HTTP 5xx: Server error
- Network timeout
- Malformed JSON response

## Data Models

```go
type APIResponse struct {
    Success bool            `json:"success"`
    Code    int             `json:"code"`
    Msg     string          `json:"msg"`
    Data    *QuotaResponse  `json:"data"`
}

type QuotaResponse struct {
    Limits []Limit `json:"limits"`
    Level  string  `json:"level"`
}

type Limit struct {
    Type          string        `json:"type"`
    Unit          int           `json:"unit"`
    Number        int           `json:"number"`
    Usage         int           `json:"usage,omitempty"`
    CurrentValue  int           `json:"currentValue,omitempty"`
    Remaining     int           `json:"remaining,omitempty"`
    Percentage    int           `json:"percentage"`
    NextResetTime int64         `json:"nextResetTime"`
    UsageDetails  []UsageDetail `json:"usageDetails,omitempty"`
}

type UsageDetail struct {
    ModelCode string `json:"modelCode"`
    Usage     int    `json:"usage"`
}
```

## Acceptance Criteria

- [ ] Client parses both TIME_LIMIT and TOKENS_LIMIT
- [ ] TIME_LIMIT correctly maps usage→total, currentValue→used
- [ ] TOKENS_LIMIT correctly extracts percentage only
- [ ] Level field is extracted and available
- [ ] Timestamps converted from milliseconds correctly
- [ ] Error responses handled with appropriate error types
