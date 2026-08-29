---
plan: .lovable/plans/pending/06-prompt-architect-installer.md
domain: Core
phase: Implement
target_files: ["gitmap/cmd/prompt_status_layout.go"]
depends_on: ["Task 030" if int(num) > 1 else "None"]
---

# Task 031 — Status Table Layout

## 1. Goal

Calculate table column dimensions for prompt status in `cmd/prompt_status_layout.go`.
