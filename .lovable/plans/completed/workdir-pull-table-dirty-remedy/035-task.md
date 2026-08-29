---
plan: .lovable/plans/pending/05-workdir-pull-table-dirty-remedy.md
domain: Core
phase: Implement
target_files: ["gitmap/gitutil/remediation_stash.go"]
depends_on: ["Task 034" if int(num) > 1 else "None"]
---

# Task 035 — Stash Remediation Recipe

## 1. Goal

Generate git stash / pop command sequence in `gitutil/remediation_stash.go`.
