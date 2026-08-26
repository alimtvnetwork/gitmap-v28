---
plan: .lovable/plans/pending/06-prompt-architect-installer.md
domain: Core
phase: Implement
target_files: ["gitmap/cmd/prompt_workdir_resolver.go"]
depends_on: ["Task 023" if int(num) > 1 else "None"]
---
# Task 024 — WorkDir Context Target Extractor

## 1. Goal
Extract targets from registered work directory context in `cmd/prompt_workdir_resolver.go`.
