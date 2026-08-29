---
plan: .lovable/plans/pending/06-prompt-architect-installer.md
domain: Core
phase: Implement
target_files: ["gitmap/cmd/prompt_target_sanitizer.go"]
depends_on: ["Task 005" if int(num) > 1 else "None"]
---

# Task 006 — Target Directory Sanitizer

## 1. Goal

Sanitize target directory paths for remote execution in `cmd/prompt_target_sanitizer.go`.
