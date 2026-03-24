## 1. Cleanup

- [x] 1.1 Delete `internal/formatter/colored.go`
- [x] 1.2 Delete `internal/formatter/colored_test.go`

## 2. Root Command Changes

- [x] 2.1 Add `textOutput` flag variable and `--text` flag to root command
- [x] 2.2 Add TTY detection using `golang.org/x/term.IsTerminal`
- [x] 2.3 Update `run()` function to handle `--text` flag and TTY detection logic
- [x] 2.4 Update mutual exclusivity check to include `--text` flag
- [x] 2.5 Wire `FormatHuman` into `runCLI()` for text output

## 3. Testing

- [x] 3.1 Update `cmd/zai-quota/root_test.go` to test `--text` flag
- [x] 3.2 Add tests for TTY detection logic
- [x] 3.3 Add tests for mutual exclusivity with `--text`
- [x] 3.4 Run full test suite and verify all tests pass

## 4. Documentation

- [x] 4.1 Update help text in `root.go` to document `--text` flag
- [x] 4.2 Verify behavior matches proposal's decision matrix
