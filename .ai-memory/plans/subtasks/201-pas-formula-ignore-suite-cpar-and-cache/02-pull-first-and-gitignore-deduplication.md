# Subtask 02: Pull-First Workflow & GitIgnore Deduplication Sanitizer

## Scope
- Files:
  - `cli/cmdpull/pull.go`
  - `cli/cmdpull/pull_efficient.go`
  - `cli/gitignoreagm/remediate.go`
  - `cli/release/gitignore.go`
- Objective:
  - Defer `checkAgmResumeTaskBeforePull` and `checkAgmResumeTaskEfficient` from pre-pull execution to post-pull summary.
  - Pull first across all repositories without blocking interactive prompts upfront.
  - Implement `.gitignore` deduplication sanitizer that removes duplicate entries idempotently while preserving section comments and ordering.
  - Ensure `.gitmap/` and `.gitmap/backup/` are present in default ignores.
