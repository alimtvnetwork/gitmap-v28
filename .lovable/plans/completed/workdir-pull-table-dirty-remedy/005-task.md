---
plan: .lovable/plans/pending/05-workdir-pull-table-dirty-remedy.md
domain: Core
phase: Implement
target_files: ["gitmap/store/workdir_delete.go"]
depends_on: ["Task 004" if int(num) > 1 else "None"]
---
# Task 005 — Store WorkDir Delete

## 1. Goal
Delete work directory record in `store/workdir_delete.go`.
