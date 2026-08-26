---
plan: .lovable/plans/pending/06-prompt-architect-installer.md
domain: Core
phase: Implement
target_files: ["gitmap/model/prompt_architect.go"]
depends_on: ["Task 001" if int(num) > 1 else "None"]
---
# Task 002 — Prompt Architect Metadata Model

## 1. Goal
Define struct for prompt architect version metadata in `model/prompt_architect.go`.
