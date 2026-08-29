---
plan: .lovable/plans/pending/05-workdir-pull-table-dirty-remedy.md
domain: Core
phase: Implement
target_files: ["gitmap/gitutil/remediation_discard.go"]
depends_on: ["Task 036" if int(num) > 1 else "None"]
---

# Task 037 — Discard Clean Remediation Recipe

## 1. Goal

Generate git reset / clean command sequence in `gitutil/remediation_discard.go`.
