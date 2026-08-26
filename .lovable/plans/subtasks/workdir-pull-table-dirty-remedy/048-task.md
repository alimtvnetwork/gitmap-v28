---
plan: .lovable/plans/pending/05-workdir-pull-table-dirty-remedy.md
domain: Core
phase: Implement
target_files: ["gitmap/tests/cmd_test/cicd_local_verify_test.go"]
depends_on: ["Task 047" if int(num) > 1 else "None"]
---
# Task 048 — Local CI/CD Quality Checks

## 1. Goal
Verify linters, type checks, and static analysis in `tests/cmd_test/cicd_local_verify_test.go`.
