---
plan: .lovable/plans/pending/05-workdir-pull-table-dirty-remedy.md
domain: Core
phase: Implement
target_files: ["gitmap/gitutil/remediation_commit.go"]
depends_on: ["Task 035" if int(num) > 1 else "None"]
---
# Task 036 — Commit Remediation Recipe

## 1. Goal
Generate git commit command sequence in `gitutil/remediation_commit.go`.
