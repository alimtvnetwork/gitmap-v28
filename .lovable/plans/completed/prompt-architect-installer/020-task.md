---
plan: .lovable/plans/pending/06-prompt-architect-installer.md
domain: Core
phase: Implement
target_files: ["gitmap/installer/prompt_runner_test.go"]
depends_on: ["Task 019" if int(num) > 1 else "None"]
---

# Task 020 — OS Runner Unit Tests

## 1. Goal

Unit tests for script command runners in `installer/prompt_runner_test.go`.
