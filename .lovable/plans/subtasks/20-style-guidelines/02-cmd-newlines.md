# Subtask 20.02: Enforce Blank Lines and Guard Clause Spacing in `gitmap/cmd/`

## Target

- Directory: `gitmap/cmd/` (.go)

## Violations to Fix

- [ ] Blank lines before `if` statements following variable declarations.
- [ ] Blank lines after closing `}` when followed by more code.
- [ ] Blank lines before `return` statements in multi-line functions.
- [ ] Guard clause separation.

## Acceptance Criteria

- [ ] Clean vertical spacing across all Go command files.
- [ ] `gofmt -w gitmap/` passes cleanly.
