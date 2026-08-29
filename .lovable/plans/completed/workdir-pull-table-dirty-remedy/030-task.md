---
plan: .lovable/plans/pending/05-workdir-pull-table-dirty-remedy.md
domain: Core
phase: Implement
target_files: ["gitmap/cmd/pull_table_test.go"]
depends_on: ["Task 029" if int(num) > 1 else "None"]
---

# Task 030 — Pull Table Test Suite

## 1. Goal

Unit tests for rich pull table rendering in `cmd/pull_table_test.go`.
