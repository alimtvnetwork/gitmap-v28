---
plan: .lovable/plans/pending/06-prompt-architect-installer.md
domain: Core
phase: Implement
target_files: ["gitmap/cmd/prompt_status_row.go"]
depends_on: ["Task 031" if int(num) > 1 else "None"]
---
# Task 032 — Status Row Formatter

## 1. Goal
Format individual prompt status rows in `cmd/prompt_status_row.go`.
