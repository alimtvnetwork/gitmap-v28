---
plan: .lovable/plans/pending/05-workdir-pull-table-dirty-remedy.md
domain: Core
phase: Implement
target_files: ["gitmap/cmd/pull_table_layout.go"]
depends_on: ["Task 025" if int(num) > 1 else "None"]
---
# Task 026 — Pull Table Header & Column Layout

## 1. Goal
Calculate column widths and table headers in `cmd/pull_table_layout.go`.
