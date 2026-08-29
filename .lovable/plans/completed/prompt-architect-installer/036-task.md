---
plan: .lovable/plans/pending/06-prompt-architect-installer.md
domain: Core
phase: Implement
target_files: ["gitmap/cmd/prompt_summary_card.go"]
depends_on: ["Task 035" if int(num) > 1 else "None"]
---

# Task 036 — Installation Success Summary Card

## 1. Goal

Render success summary banner after installation in `cmd/prompt_summary_card.go`.
