---
plan: .lovable/plans/pending/06-prompt-architect-installer.md
domain: Core
phase: Implement
target_files: ["gitmap/model/prompt_status.go"]
depends_on: ["Task 007" if int(num) > 1 else "None"]
---
# Task 008 — Prompt Installation Status Struct

## 1. Goal
Define installation status data structure in `model/prompt_status.go`.
