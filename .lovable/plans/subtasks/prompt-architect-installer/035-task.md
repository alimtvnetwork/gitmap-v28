---
plan: .lovable/plans/pending/06-prompt-architect-installer.md
domain: Core
phase: Implement
target_files: ["gitmap/cmd/prompt_version_cmd.go"]
depends_on: ["Task 034" if int(num) > 1 else "None"]
---
# Task 035 — Version Command Entrypoint

## 1. Goal
CLI entrypoint for gitmap ct prompts-version in `cmd/prompt_version_cmd.go`.
