---
plan: .lovable/plans/pending/05-workdir-pull-table-dirty-remedy.md
domain: Core
phase: Implement
target_files: ["gitmap/cmd/workdir_set_default.go"]
depends_on: ["Task 007" if int(num) > 1 else "None"]
---
# Task 008 — CLI WorkDir Set Default Command

## 1. Goal
CLI for gitmap workdir set-default in `cmd/workdir_set_default.go`.
