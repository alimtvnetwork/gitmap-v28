---
plan: .lovable/plans/pending/06-prompt-architect-installer.md
domain: Core
phase: Implement
target_files: ["gitmap/cmd/prompt_status_cmd.go"]
depends_on: ["Task 033" if int(num) > 1 else "None"]
---

# Task 034 — Status Command Entrypoint

## 1. Goal

CLI entrypoint for gitmap ct prompts-status in `cmd/prompt_status_cmd.go`.
