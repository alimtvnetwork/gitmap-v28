# Subtask 01 - Diagnose CI failures

Date: 31-Oct-26
Owner: solo worker (parent task forbids subagents).

- Record live CI before state: CI 41198431481, Cross-Platform
  41198431561, both on main 58292e98 (Release v6.643.0).
- Reproduce locally: gofmt, go test for failing packages, the
  policy-check linters, test_ci_scripts.py.
- Output: exact failing files, tests, and linter rules, filed in
  the spec before fixes start.
