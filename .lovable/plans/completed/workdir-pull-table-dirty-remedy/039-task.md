---
plan: .lovable/plans/pending/05-workdir-pull-table-dirty-remedy.md
domain: Core
phase: Implement
target_files: ["gitmap/cmd/remediation_box.go"]
depends_on: ["Task 038" if int(num) > 1 else "None"]
---

# Task 039 — Interactive Remediation Suggestion Box

## 1. Goal

Render bordered remediation hint box in `cmd/remediation_box.go`.
