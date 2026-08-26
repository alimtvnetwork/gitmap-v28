---
plan: .lovable/plans/pending/05-workdir-pull-table-dirty-remedy.md
domain: Core
phase: Implement
target_files: ["gitmap/gitutil/remediation_test.go"]
depends_on: ["Task 039" if int(num) > 1 else "None"]
---
# Task 040 — Remediation Test Suite

## 1. Goal
Unit tests for remediation recipe generation in `gitutil/remediation_test.go`.
