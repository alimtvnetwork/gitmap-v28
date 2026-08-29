---
plan: .lovable/plans/pending/05-workdir-pull-table-dirty-remedy.md
domain: Core
phase: Implement
target_files: ["gitmap/cmd/pull_table_summary.go"]
depends_on: ["Task 028" if int(num) > 1 else "None"]
---

# Task 029 — Pull Table Summary Aggregator

## 1. Goal

Aggregate totals, times, and pull results in `cmd/pull_table_summary.go`.
