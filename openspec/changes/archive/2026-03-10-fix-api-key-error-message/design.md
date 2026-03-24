## Context

The `zai-quota` CLI tool validates API key presence in `cmd/zai-quota/root.go:52`. When missing, it returns an error that only mentions the `ZAI_API_KEY` environment variable, despite the configuration loader (`internal/config/loader.go`) supporting both:

1. Environment variable: `ZAI_API_KEY`
2. Config file: `~/.zai-quota.yaml` with `api_key: <value>`

This creates a disconnect between documentation (README.md lists both options) and the error message.

## Goals / Non-Goals

**Goals:**
- Update error message to accurately reflect all valid configuration methods
- Align error output with README documentation
- Improve user experience by reducing confusion

**Non-Goals:**
- Changes to the config loader logic
- Adding new configuration methods
- Changes to exit codes or error handling structure

## Decisions

### Keep validation in root.go (not loader.go)

**Rationale:** The loader's responsibility is loading configuration from available sources. Business logic ("API key is required") belongs at the call site where the decision is made about what to do with the loaded config.

**Alternative considered:** Move validation to loader.go so it could return richer context about what sources were attempted. Rejected because:
- Adds business logic to a general-purpose loader
- The simple string fix provides 80% of the value with 5% of the work

### Error message format

```
API key required. Set ZAI_API_KEY environment variable or add 'api_key: YOUR_KEY' to ~/.zai-quota.yaml
```

**Rationale:**
- Concise but complete
- Mentions both valid sources
- Shows exact config file path and key name
- Follows existing error message style in the codebase

## Risks / Trade-offs

**Risk:** None significant. This is a low-risk text change.

**Trade-off:** Slightly longer error message vs. completeness. Completeness wins for user experience.
