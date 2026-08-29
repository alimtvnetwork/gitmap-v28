---
plan: .lovable/plans/pending/05-workdir-pull-table-dirty-remedy.md
domain: Core
phase: Implement
target_files: ["gitmap/store/workdir_default.go"]
depends_on: ["Task 003" if int(num) > 1 else "None"]
---

# Task 004 — Store WorkDir Set Default

## 1. Goal

Set and query active default work directory in `store/workdir_default.go`.
