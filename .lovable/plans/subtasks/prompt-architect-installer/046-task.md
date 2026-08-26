---
plan: .lovable/plans/pending/06-prompt-architect-installer.md
domain: Core
phase: Implement
target_files: ["gitmap/tests/cmd_test/prompt_platform_test.go"]
depends_on: ["Task 045" if int(num) > 1 else "None"]
---
# Task 046 — Cross-Platform Execution Contract Test

## 1. Goal
Verify platform command generator contract in `tests/cmd_test/prompt_platform_test.go`.
