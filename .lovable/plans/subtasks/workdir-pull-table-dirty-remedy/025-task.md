---
plan: .lovable/plans/pending/05-workdir-pull-table-dirty-remedy.md
domain: Core
phase: Implement
target_files: ["gitmap/model/pull_table_row.go"]
depends_on: ["Task 024" if int(num) > 1 else "None"]
---
# Task 025 — Pull Table Row Model Definition

## 1. Goal
Define struct for rich pull table rows in `model/pull_table_row.go`.
