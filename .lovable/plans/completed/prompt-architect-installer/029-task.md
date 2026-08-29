---
plan: .lovable/plans/pending/06-prompt-architect-installer.md
domain: Core
phase: Implement
target_files: ["gitmap/cmd/prompt_inplace_policy.go"]
depends_on: ["Task 028" if int(num) > 1 else "None"]
---

# Task 029 — In-Place Update Policy Evaluator

## 1. Goal

Verify update policy when prompts already exist in `cmd/prompt_inplace_policy.go`.
