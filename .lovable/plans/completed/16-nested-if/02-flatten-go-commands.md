# Subtask 02: Flatten Go Command Implementations

## Description

Eliminate nested `if` statements in `gitmap/cmd/fix_cmd.go` and `gitmap/cmd/rootusage.go` using guard clauses and helper function extraction.

## Citations

- [Braces and Nesting](spec/02-coding-guidelines/01-cross-language/04-code-style/01-braces-and-nesting.md)
- [Canonical Size Tier](spec/02-coding-guidelines/00-canonical-size-tier.md)
- [Master Coding Guidelines](.lovable/coding-guidelines/coding-guidelines.md)

## Acceptance Criteria

- [ ] No nested `if` statements in command files.
- [ ] All functions <= 15 lines (target <= 8 lines).
