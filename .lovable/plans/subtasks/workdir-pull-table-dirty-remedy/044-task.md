---
plan: .lovable/plans/pending/05-workdir-pull-table-dirty-remedy.md
domain: Core
phase: Implement
target_files: ["gitmap/cmd/status_integrated.go"]
depends_on: ["Task 043" if int(num) > 1 else "None"]
---
# Task 044 — Status Command Integration

## 1. Goal
Integrate rich columns and dirty diagnostics in `cmd/status_integrated.go`.
