---
plan: .lovable/plans/pending/05-workdir-pull-table-dirty-remedy.md
domain: Core
phase: Implement
target_files: ["gitmap/cmd/workdir_flags.go"]
depends_on: ["Task 005" if int(num) > 1 else "None"]
---
# Task 006 — CLI WorkDir Flag Parser

## 1. Goal
Parse flags for gitmap workdir / wd in `cmd/workdir_flags.go`.
