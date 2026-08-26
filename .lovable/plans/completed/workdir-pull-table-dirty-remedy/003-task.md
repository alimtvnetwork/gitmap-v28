---
plan: .lovable/plans/pending/05-workdir-pull-table-dirty-remedy.md
domain: Core
phase: Implement
target_files: ["gitmap/store/workdir_list.go"]
depends_on: ["Task 002" if int(num) > 1 else "None"]
---
# Task 003 — Store WorkDir List

## 1. Goal
List and ensure work directory records in `store/workdir_list.go`.
