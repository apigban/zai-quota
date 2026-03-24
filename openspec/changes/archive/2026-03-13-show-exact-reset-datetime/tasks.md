## 1. Timezone Detection

- [x] 1.1 Create `internal/processor/timezone.go` with `GetTimezone()` function
- [x] 1.2 Implement TZ environment variable check
- [x] 1.3 Implement `/etc/timezone` file reading
- [x] 1.4 Implement `/etc/localtime` symlink resolution
- [x] 1.5 Add `extractCity(iana string) string` helper function
- [x] 1.6 Add unit tests for timezone detection

## 2. Time Formatting

- [x] 2.1 Add `FormatResetDateTime(t time.Time) string` to `internal/processor/time.go`
- [x] 2.2 Integrate timezone detection into formatting function
- [x] 2.3 Add unit tests for datetime formatting with various timezones

## 3. TUI Integration

- [x] 3.1 Add `ResetDateTime string` field to `ProcessedLimitData` struct
- [x] 3.2 Update `processLimits()` to populate `ResetDateTime`
- [x] 3.3 Update `renderLimit()` to display new format: `Reset: {duration} ({datetime})`
- [x] 3.4 Verify both TOKENS_LIMIT and TIME_LIMIT show enhanced format

## 4. Testing & Verification

- [x] 4.1 Run existing tests to ensure no regressions
- [x] 4.2 Test with timezone detected (TZ env var or /etc/timezone)
- [x] 4.3 Test UTC fallback (unset TZ, mock failed detection)
- [x] 4.4 Manual TUI verification with `./zai-quota`
