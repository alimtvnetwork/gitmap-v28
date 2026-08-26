---
plan: .lovable/plans/pending/06-prompt-architect-installer.md
domain: Core
phase: Implement
target_files: ["gitmap/tests/cmd_test/prompt_cicd_test.go"]
depends_on: ["Task 046" if int(num) > 1 else "None"]
---
# Task 047 — Local CI/CD Type Checking & Linter Pass

## 1. Goal
Verify clean lint and type checks in `tests/cmd_test/prompt_cicd_test.go`.
