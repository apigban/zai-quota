## Why

When users run `zai-quota` without an API key configured, the error message states: "ZAI_API_KEY environment variable is required". This is misleading because the tool actually supports two configuration methods: environment variables AND a config file (`~/.zai-quota.yaml`). Users who read the README and prefer the config file approach are confused when the error only mentions environment variables.

## What Changes

- Update the API key validation error message in `root.go` to mention both configuration options
- The new message will guide users to either set `ZAI_API_KEY` environment variable OR add `api_key` to `~/.zai-quota.yaml`

## Capabilities

### New Capabilities

- `error-message`: Guidance for error messages that mention all valid configuration options

### Modified Capabilities

- None (no existing specs to modify)

## Impact

- **Affected files**: `cmd/zai-quota/root.go` (line 52)
- **User experience**: Users will see helpful, accurate error messages
- **Documentation alignment**: Error message will be consistent with README documentation
