---
plan: .lovable/plans/pending/06-prompt-architect-installer.md
domain: Core
phase: Implement
target_files: ["gitmap/cmd/prompt_interactive_confirm.go"]
depends_on: ["Task 038" if int(num) > 1 else "None"]
---

# Task 039 — Interactive Confirmation Prompt

## 1. Goal

Prompt user before batch installation if required in `cmd/prompt_interactive_confirm.go`.
