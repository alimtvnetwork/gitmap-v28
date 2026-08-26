---
plan: .lovable/plans/pending/05-workdir-pull-table-dirty-remedy.md
domain: Core
phase: Implement
target_files: ["gitmap/cmd/pullparallel_integrated.go"]
depends_on: ["Task 041" if int(num) > 1 else "None"]
---
# Task 042 — Pull Parallel Integration

## 1. Goal
Integrate rich table into parallel pull runner in `cmd/pullparallel_integrated.go`.
