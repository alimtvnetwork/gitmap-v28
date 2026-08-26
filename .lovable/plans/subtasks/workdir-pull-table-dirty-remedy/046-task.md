---
plan: .lovable/plans/pending/05-workdir-pull-table-dirty-remedy.md
domain: Core
phase: Implement
target_files: ["gitmap/tests/cmd_test/workdir_cli_test.go"]
depends_on: ["Task 045" if int(num) > 1 else "None"]
---
# Task 046 — WorkDir CLI Integration Tests

## 1. Goal
E2E tests for workdir commands in `tests/cmd_test/workdir_cli_test.go`.
