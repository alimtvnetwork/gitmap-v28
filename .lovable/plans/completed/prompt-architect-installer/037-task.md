---
plan: .lovable/plans/pending/06-prompt-architect-installer.md
domain: Core
phase: Implement
target_files: ["gitmap/cmd/prompt_failure_reporter.go"]
depends_on: ["Task 036" if int(num) > 1 else "None"]
---
# Task 037 — Failure Diagnostic Reporter

## 1. Goal
Report failed repositories with error explanations in `cmd/prompt_failure_reporter.go`.
