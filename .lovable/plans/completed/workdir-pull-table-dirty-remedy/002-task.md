---
plan: .lovable/plans/pending/05-workdir-pull-table-dirty-remedy.md
domain: Core
phase: Implement
target_files: ["gitmap/store/workdir_schema.go"]
depends_on: ["Task 001" if int(num) > 1 else "None"]
---

# Task 002 — WorkDir SQLite Schema

## 1. Goal

WorkDir table migrations and queries in `store/workdir_schema.go`.
