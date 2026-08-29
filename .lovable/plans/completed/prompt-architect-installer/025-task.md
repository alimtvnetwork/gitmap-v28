---
plan: .lovable/plans/pending/06-prompt-architect-installer.md
domain: Core
phase: Implement
target_files: ["gitmap/cmd/prompt_exclusion_filter.go"]
depends_on: ["Task 024" if int(num) > 1 else "None"]
---

# Task 025 — Target Exclusion Filter

## 1. Goal

Filter out excluded repositories from prompt installation in `cmd/prompt_exclusion_filter.go`.
