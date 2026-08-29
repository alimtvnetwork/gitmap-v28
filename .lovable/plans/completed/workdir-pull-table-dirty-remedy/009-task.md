---
plan: .lovable/plans/pending/05-workdir-pull-table-dirty-remedy.md
domain: Core
phase: Implement
target_files: ["gitmap/cmd/workdir_add_rm.go"]
depends_on: ["Task 008" if int(num) > 1 else "None"]
---

# Task 009 — CLI WorkDir Add and RM Commands

## 1. Goal

CLI for gitmap workdir add and rm in `cmd/workdir_add_rm.go`.
