---
plan: .lovable/plans/pending/05-workdir-pull-table-dirty-remedy.md
domain: Core
phase: Implement
target_files: ["gitmap/cmd/pull_integrated.go"]
depends_on: ["Task 040" if int(num) > 1 else "None"]
---

# Task 041 — Pull Command Full Integration

## 1. Goal

Integrate rich table and remediation into pull command in `cmd/pull_integrated.go`.
