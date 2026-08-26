---
plan: .lovable/plans/pending/05-workdir-pull-table-dirty-remedy.md
domain: Core
phase: Implement
target_files: ["gitmap/cmd/pull_table_style.go"]
depends_on: ["Task 027" if int(num) > 1 else "None"]
---
# Task 028 — Pull Table Styling & Glyphs

## 1. Goal
Apply color tokens and status icons in `cmd/pull_table_style.go`.
