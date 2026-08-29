---
plan: .lovable/plans/pending/06-prompt-architect-installer.md
domain: Core
phase: Implement
target_files: ["gitmap/installer/prompt_update_guard.go"]
depends_on: ["Task 017" if int(num) > 1 else "None"]
---

# Task 018 — In-Place Update Overwrite Guard

## 1. Goal

Handle in-place updates and overwrite existing files in `installer/prompt_update_guard.go`.
