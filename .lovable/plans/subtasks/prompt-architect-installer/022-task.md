---
plan: .lovable/plans/pending/06-prompt-architect-installer.md
domain: Core
phase: Implement
target_files: ["gitmap/cmd/prompt_args_parser.go"]
depends_on: ["Task 021" if int(num) > 1 else "None"]
---
# Task 022 — Multi-Target Argument Parser

## 1. Goal
Parse arguments and flags for ct install-prompts in `cmd/prompt_args_parser.go`.
