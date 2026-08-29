---
plan: .lovable/plans/pending/05-workdir-pull-table-dirty-remedy.md
domain: Core
phase: Implement
target_files: ["gitmap/cmd/pull_table_row.go"]
depends_on: ["Task 026" if int(num) > 1 else "None"]
---

# Task 027 — Pull Table Row Formatter

## 1. Goal

Format individual pull table rows with lipgloss in `cmd/pull_table_row.go`.
