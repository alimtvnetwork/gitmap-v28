---
plan: .lovable/plans/pending/05-workdir-pull-table-dirty-remedy.md
domain: Core
phase: Implement
target_files: ["gitmap/tests/cmd_test/pull_workflow_e2e_test.go"]
depends_on: ["Task 046" if int(num) > 1 else "None"]
---

# Task 047 — Pull Workflow E2E Test Suite

## 1. Goal

E2E tests for auto-pull in parent directory in `tests/cmd_test/pull_workflow_e2e_test.go`.
