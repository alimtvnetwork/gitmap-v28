---
plan: .lovable/plans/pending/05-workdir-pull-table-dirty-remedy.md
domain: Core
phase: Implement
target_files: ["gitmap/cmd/workdir_ls.go"]
depends_on: ["Task 006" if int(num) > 1 else "None"]
---
# Task 007 — CLI WorkDir LS Table Renderer

## 1. Goal
Render work directories status table in `cmd/workdir_ls.go`.
