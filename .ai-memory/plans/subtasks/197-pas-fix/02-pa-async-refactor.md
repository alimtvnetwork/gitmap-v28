# Subtask 02: PA Async Refactor

**Objective**: Decouple ignore duplicate and tracking checks from synchronous git pull workflow in `gitmap pa`.

## Requirements
- Move `.gitignore` inspection into background async workers in `cli/cmdpull/pull.go`.
- Ensure pull begins immediately without blocking on ignore parsing.
- Aggregate ignore duplicate summaries at the completion of all pulls.

## Constraints
- Bounding Box: `cli/cmdpull/*.go`
- Coding Rules: Positive booleans only, AppError wrappers, functions <= 15 lines.
