# Post-Mortem: Unused `cdProfilePath` Function Lint Failure

## Summary

CI lint stage failed because `completion/cdfunction.go` exports an unused function `cdProfilePath`. `golangci-lint` (v1.64.8) flagged it via the `unused` linter, breaking the build.

## Error

```
Error: gitmap-v27/completion/cdfunction.go:37:6: func `cdProfilePath` is unused (unused)

func cdProfilePath(shell string) string {
     ^

Error: issues found
```

## Root Cause

The `cdProfilePath` function was defined in `completion/cdfunction.go:37` but never called from any other code path. This is likely a remnant from a refactor or a function prepared for future use that was never wired in.

Go's `unused` linter (enabled via `golangci-lint`) treats unreferenced unexported functions as errors, not warnings.

## Resolution Options

### Option A — Remove the Function (Preferred)

If `cdProfilePath` is no longer needed, delete the function entirely from `completion/cdfunction.go`.

### Option B — Wire It In

If the function is needed for shell profile path resolution (e.g., determining where to write `cd` wrapper functions), connect it to the appropriate call site in the `cd` or `clone-next` shell integration flow.

### Option C — Suppress (Not Recommended)

Add a `//nolint:unused` directive. This violates the project's zero-suppression policy unless justified with an inline comment explaining the deferral reason.

## Prevention Rules

1. **Run `golangci-lint` locally before pushing** — all new functions must have at least one call site or test reference.
2. **Refactors that remove call sites must also remove orphaned functions** — search for all references before deleting a caller.
3. **Do not add speculative functions** — implement only when a call site exists.

## Related

- `spec/02-app-issues/23-go-build-copyfile-redeclared.md` — Similar Go build failure from namespace conflicts
- `spec/05-coding-guidelines/01-code-quality-improvement.md` — Lint enforcement rules
- `spec/03-general/06-code-style-rules.md` — Code style rules
