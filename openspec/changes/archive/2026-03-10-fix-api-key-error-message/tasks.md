## 1. Implementation

- [x] 1.1 Update API key validation error message in `cmd/zai-quota/root.go:52` to mention both configuration options (environment variable and config file)

## 2. Verification

- [x] 2.1 Build the binary and verify no compilation errors
- [x] 2.2 Test error message appears correctly when running without API key configured
- [x] 2.3 Verify error message mentions both `ZAI_API_KEY` and `~/.zai-quota.yaml`
