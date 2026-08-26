---
plan: .lovable/plans/pending/06-prompt-architect-installer.md
domain: Core
phase: Implement
target_files: ["gitmap/cmd/prompt_cmd_builder.go"]
depends_on: ["Task 006" if int(num) > 1 else "None"]
---
# Task 007 — Command Execution Parameter Formatter

## 1. Goal
Build OS-specific execution command strings in `cmd/prompt_cmd_builder.go`.
