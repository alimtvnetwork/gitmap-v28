# Strictly Avoid

Hard prohibitions. Violating ANY of these is a critical failure.

---

## File and Directory Rules

- **NEVER** manually create, modify, or delete files within `.gitmap/release/` or `.gitmap/release-assets/`.
- **NEVER** touch the `.release/` legacy folder (if it exists, migration handles it).
- **NEVER** create files exceeding 200 lines. Split by responsibility.
- **NEVER** create functions exceeding 15 lines.

## Code Style Prohibitions

- **NEVER** use negation in `if` conditions (`!`, `!=`). Use positive logic.
- **NEVER** use `switch` statements. Use `if`/`else if` chains.
- **NEVER** use magic strings. All literals must be in the `constants` package.
- **NEVER** use single-character variable names (`s`, `x`, `d`).
- **NEVER** use inline type definitions. Extract named types/interfaces.
- **NEVER** skip the blank line before `return` (unless return is the sole line in an `if`).

## Error Handling Prohibitions

- **NEVER** swallow errors silently. Every error must be explicitly logged to `os.Stderr`.
- **NEVER** use generic error messages. Include path, operation, and reason context.

## Naming Prohibitions

- **NEVER** use snake_case or camelCase for DB table/column names. PascalCase only.
- **NEVER** use boolean variables/functions without `is`/`has` prefix.
- **NEVER** place CLI command IDs anywhere except `constants_cli.go`.

## Dependency and Tooling Prohibitions

- **NEVER** use `@latest` for tool installations in CI. Pin exact versions.
- **NEVER** use `go install tool@latest`. All tools pinned: `golangci-lint@v1.64.8`, `govulncheck@v1.1.4`.
- **NEVER** use range operators for npm dependencies. Pin exact versions.

## Database Prohibitions

- **NEVER** open more than 1 SQLite connection (`SetMaxOpenConns(1)`).
- **NEVER** use UUID strings for primary keys. Use `INTEGER PRIMARY KEY AUTOINCREMENT`.

## CI/CD Prohibitions

- **NEVER** use `cd` in CI scripts. Use `working-directory` in workflow steps.
- **NEVER** use job-level `if` for SHA deduplication. Use the passthrough gate pattern.
- **NEVER** cancel release branch CI runs (`cancel-in-progress: false` for `release/**`).

## Communication Prohibitions

- **NEVER** append boilerplate like "Let me know if you have questions!" or "Hope this helps!"
- **NEVER** append the "If you have any question and confusion..." block.
- **NEVER** append the "Do you understand? Always add this part..." code block.

## Platform Prohibitions

- **NEVER** use non-ASCII punctuation in PowerShell scripts or CLI terminal messages. ASCII only.
- **NEVER** store roles on the profile or users table. Use a separate roles table.
