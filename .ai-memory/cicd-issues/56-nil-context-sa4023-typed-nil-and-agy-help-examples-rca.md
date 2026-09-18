# RCA: Nil Context (SA1012), Typed-Nil Comparison (SA4023), Gofmt Drift, and AGY Help Examples

## 1. Symptom

GitHub Actions CI run `#35343946879` for commit `d616f68` failed across multiple jobs and steps:

1. **Lint / golangci-lint**:
   - `cli/cmdssh/ssh_install_actions.go:38:40`: SA1012: do not pass a nil Context, even if a function permits it; pass context.TODO if a Context is not available (staticcheck)
   - `cli/cmdssh/ssh_target_nodes.go:25:40`: SA1012: do not pass a nil Context, even if a function permits it (staticcheck)
   - `cli/cmdssh/sshexec.go:127:40`: SA1012: do not pass a nil Context, even if a function permits it (staticcheck)
   - `cli/cmdagy/agy_conv_ws.go:16:5`: SA4023: this comparison is never true; the value is not nil (staticcheck)

2. **Lint Script Unit Tests**:
   - `test_gofmt_check_clean_repo` failed with `AssertionError: 1 != 0` due to unformatted Go source files:
     - `cli/cmdagy/agy_list_prompts.go`
     - `cli/cmdagy/agy_rerun.go`
     - `cli/cmdprompttemplate/template_crud.go`
     - `cli/cmdssh/sshexec.go`

3. **Golden Test / Build Guards Across OS Platforms (Windows, macOS, Ubuntu)**:
   - `TestEveryHelpFileHasExamples` in `cli/helptext/examples_golden_test.go` failed:
     ```text
     help files missing `## Examples` section (1):
       - agy.md
     ```

## 2. Root Cause

1. **SA1012 (Nil Context)**: Calls to `CheckConnLiveness(nil, ...)` passed a literal `nil` for the `context.Context` parameter instead of a non-nil context such as `context.Background()`. Staticcheck strictly forbids passing `nil` contexts.
2. **SA4023 (Typed-Nil Interface Trap)**: In `resolveConvWorkspace` (`cli/cmdagy/agy_conv_ws.go`), `err` was originally declared as type `error` (interface). Later, `conn, err := store.OpenSQLiteDB(dbPath)` assigned a concrete pointer `*apperror.AppError` to `err`. In Go, storing a typed nil pointer into an interface produces a non-nil interface (`err != nil` is always true).
3. **Gofmt Drift**: Files modified in recent commits were not formatted using `gofmt -s -w` before committing, leaving spacing and struct alignment inconsistencies.
4. **Golden Help Examples Section**: `cli/helptext/examples_golden_test.go` verifies that all command help files have a top-level `\n## Examples` markdown heading containing fenced code blocks. `cli/helptext/agy.md` contained sub-level `### Examples` under each subcommand section, but lacked the required top-level `## Examples` section.

## 3. Resolution

1. **Fixed SA1012**: Replaced `nil` with `context.Background()` in calls to `CheckConnLiveness` in:
   - `cli/cmdssh/ssh_install_actions.go`
   - `cli/cmdssh/ssh_target_nodes.go`
   - `cli/cmdssh/sshexec.go`
2. **Fixed SA4023**: Decoupled `error` variable assignments in `cli/cmdagy/agy_conv_ws.go` by introducing a distinct `dbErr` variable for `store.OpenSQLiteDB`.
3. **Formatted Go Code**: Executed `python 03-ai-scripts/26-go-code-formatter.py` and `gofmt -s -w`, ensuring 0 unformatted files remain (`gofmt -l cli/` outputs nothing).
4. **Added `## Examples` in `agy.md`**: Added a level-2 `## Examples` section with fenced code blocks to `cli/helptext/agy.md`, satisfying `TestEveryHelpFileHasExamples`.
5. **Verified Linters & Gates**:
   - `golangci-lint run ./...` in `cli/`: Passed with 0 errors.
   - `go vet ./...` in `cli/`: Passed with 0 errors.
   - `go test ./helptext -run TestEveryHelpFileHasExamples`: Passed (`ok`).
   - `python linter-scripts/check-nested-ifs.py`: Passed (0 violations across 3,142 files).
   - `python linter-scripts/check-enum-and-boolean.py`: Passed (0 violations across 2,368 files).
   - `python linter-scripts/check-relative-paths.py`: Passed (0 violations across 7,093 files).

## 4. Prevention & Learnings

- **Never pass `nil` to `context.Context`**: Always use `context.Background()` or `context.TODO()` to prevent SA1012 staticcheck failures.
- **Beware the Go typed-nil interface trap**: When a function returns a concrete error type (such as `*apperror.AppError`), never assign it to a pre-existing `error` interface variable using `:=` or `=` if `err != nil` will subsequently be checked. Use a dedicated variable with the concrete error type.
- **Always format Go files before commit**: Run `python 03-ai-scripts/26-go-code-formatter.py` or `gofmt -s -w` across all touched Go files.
- **Enforce `## Examples` on all command markdown files**: Any new or updated command documentation under `cli/helptext/` must have a top-level `## Examples` header followed by runnable code blocks.
