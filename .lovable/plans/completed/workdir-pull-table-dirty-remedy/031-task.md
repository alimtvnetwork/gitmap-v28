---
plan: .lovable/plans/pending/05-workdir-pull-table-dirty-remedy.md
domain: Core
phase: Implement
target_files: ["gitmap/gitutil/dirty_inspect.go"]
depends_on: ["Task 030" if int(num) > 1 else "None"]
---
# Task 031 — Dirty Status Detailed Inspector

## 1. Goal
Inspect uncommitted files and unstaged changes in `gitutil/dirty_inspect.go`.
