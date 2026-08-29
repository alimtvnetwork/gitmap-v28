---
plan: .lovable/plans/pending/05-workdir-pull-table-dirty-remedy.md
domain: Core
phase: Implement
target_files: ["gitmap/cmd/pull_error_formatter.go"]
depends_on: ["Task 037" if int(num) > 1 else "None"]
---

# Task 038 — Pull Failure Error Formatter

## 1. Goal

Format pull error with explanation and recipe in `cmd/pull_error_formatter.go`.
