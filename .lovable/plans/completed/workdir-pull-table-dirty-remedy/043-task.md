---
plan: .lovable/plans/pending/05-workdir-pull-table-dirty-remedy.md
domain: Core
phase: Implement
target_files: ["gitmap/cmd/push_integrated.go"]
depends_on: ["Task 042" if int(num) > 1 else "None"]
---

# Task 043 — Push Command Full Integration

## 1. Goal

Integrate top-level discovery into push command in `cmd/push_integrated.go`.
