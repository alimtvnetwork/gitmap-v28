---
plan: .lovable/plans/pending/05-workdir-pull-table-dirty-remedy.md
domain: Core
phase: Implement
target_files: ["gitmap/gitutil/remediation_generator.go"]
depends_on: ["Task 033" if int(num) > 1 else "None"]
---
# Task 034 — Remediation Recipe Generator

## 1. Goal
Generate structured remediation recipes in `gitutil/remediation_generator.go`.
