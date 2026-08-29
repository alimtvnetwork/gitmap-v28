---
plan: .lovable/plans/pending/06-prompt-architect-installer.md
domain: Core
phase: Implement
target_files: ["gitmap/cmd/prompt_remediation.go"]
depends_on: ["Task 037" if int(num) > 1 else "None"]
---

# Task 038 — Remediation Advice Formatter

## 1. Goal

Format manual installation one-liner advice in `cmd/prompt_remediation.go`.
