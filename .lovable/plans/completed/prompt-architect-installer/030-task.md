---
plan: .lovable/plans/pending/06-prompt-architect-installer.md
domain: Core
phase: Implement
target_files: ["gitmap/cmd/prompt_target_test.go"]
depends_on: ["Task 029" if int(num) > 1 else "None"]
---
# Task 030 — Target Resolution Test Suite

## 1. Goal
Unit tests for target resolution and args parsing in `cmd/prompt_target_test.go`.
