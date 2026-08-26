---
plan: .lovable/plans/pending/06-prompt-architect-installer.md
domain: Core
phase: Implement
target_files: ["gitmap/installer/prompt_exec_context.go"]
depends_on: ["Task 013" if int(num) > 1 else "None"]
---
# Task 014 — Process Execution Timeout & Context

## 1. Goal
Manage process execution context and timeouts in `installer/prompt_exec_context.go`.
