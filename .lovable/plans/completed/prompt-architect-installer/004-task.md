---
plan: .lovable/plans/pending/06-prompt-architect-installer.md
domain: Core
phase: Implement
target_files: ["gitmap/cmd/prompt_version_writer.go"]
depends_on: ["Task 003" if int(num) > 1 else "None"]
---
# Task 004 — Version JSON Writer

## 1. Goal
Write promptArchitectByRiseupAsia section to version.json in `cmd/prompt_version_writer.go`.
