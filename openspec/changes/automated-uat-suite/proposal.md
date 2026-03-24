## Why

The Z.ai Quota Monitor needs an automated User Acceptance Testing (UAT) suite that can be executed by AI coding assistants with access to bash, curl, and other testing tools. Currently, testing is limited to unit tests (`make test`) which verify internal logic but don't validate end-to-end behavior of the compiled binary.

An automated UAT suite enables:
- **Release validation**: Verify all capabilities work correctly before releases
- **Regression detection**: Catch breaking changes early
- **AI-executable testing**: Allow AI assistants to validate changes autonomously
- **Documentation by example**: Tests serve as executable documentation

## What Changes

- New `tests/uat/` directory containing the complete UAT suite
- Mock Z.ai API server (`tests/uat/mock/`) as a standalone Go binary
- Shell scripts (`tests/uat/tests/`) for each test scenario
- Test library (`tests/uat/lib/`) with shared utilities
- Main runner script (`tests/uat/run_uat.sh`) as single entry point
- Expected outputs (`tests/uat/expected/`) for comparison assertions

## Capabilities

### New Capabilities

- `automated-uat`: Automated end-to-end testing suite for all Z.ai Quota Monitor capabilities

### Modified Capabilities

(None - this is a new testing infrastructure with no changes to application capabilities)

## Impact

- **New directory**: `tests/uat/` with full test suite
- **New binary**: `tests/uat/mock/mock-server` (compiled mock API)
- **Build changes**: Add `make uat` target to Makefile
- **CI/CD ready**: Scripts return proper exit codes for CI integration
- **No application code changes**: Pure testing infrastructure
