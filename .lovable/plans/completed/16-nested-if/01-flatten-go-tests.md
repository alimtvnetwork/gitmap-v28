# Subtask 01: Flatten Go Test Conditionals

## Description
Eliminate all nested `if` statements across Go test files (`gitmap/cmd/comp_053_test.go`, `gitmap/cmd/comp_102_test.go`, `gitmap/cmd/comp_109_test.go`, `gitmap/cmd/ip_cmd_test.go`, `gitmap/cmd/ssh_client_test.go`, `gitmap/scanner/progress_test.go`, `gitmap/store/ssh_repo_test.go`) using early assertions and guard clauses.

## Citations
- [Braces and Nesting](spec/02-coding-guidelines/01-cross-language/04-code-style/01-braces-and-nesting.md)
- [Canonical Size Tier](spec/02-coding-guidelines/00-canonical-size-tier.md)
- [Master Coding Guidelines](.lovable/coding-guidelines/coding-guidelines.md)

## Acceptance Criteria
- [ ] No nested `if` statements in Go test files.
- [ ] All unit tests pass.
