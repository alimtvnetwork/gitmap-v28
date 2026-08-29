---
plan: .lovable/plans/pending/06-prompt-architect-installer.md
domain: Core
phase: Implement
target_files: ["gitmap/cmd/prompt_status_test.go"]
depends_on: ["Task 039" if int(num) > 1 else "None"]
---

# Task 040 — Status & Dashboard Unit Tests

## 1. Goal

Unit tests for prompt status display in `cmd/prompt_status_test.go`.
