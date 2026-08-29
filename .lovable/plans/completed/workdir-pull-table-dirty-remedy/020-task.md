---
plan: .lovable/plans/pending/05-workdir-pull-table-dirty-remedy.md
domain: Core
phase: Implement
target_files: ["gitmap/tests/cmd_test/recursive_discovery_test.go"]
depends_on: ["Task 019" if int(num) > 1 else "None"]
---

# Task 020 — Recursive Discovery Test Suite

## 1. Goal

Unit tests for top-level repo discovery in `tests/cmd_test/recursive_discovery_test.go`.
