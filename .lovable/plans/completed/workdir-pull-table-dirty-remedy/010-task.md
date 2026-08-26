---
plan: .lovable/plans/pending/05-workdir-pull-table-dirty-remedy.md
domain: Core
phase: Implement
target_files: ["gitmap/cmd/workdir_cmd.go"]
depends_on: ["Task 009" if int(num) > 1 else "None"]
---
# Task 010 — WorkDir Root Dispatcher Wiring

## 1. Goal
Wire workdir and wd into root CLI dispatcher in `cmd/workdir_cmd.go`.
