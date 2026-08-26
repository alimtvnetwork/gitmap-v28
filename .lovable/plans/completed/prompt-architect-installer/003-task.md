---
plan: .lovable/plans/pending/06-prompt-architect-installer.md
domain: Core
phase: Implement
target_files: ["gitmap/cmd/prompt_version_reader.go"]
depends_on: ["Task 002" if int(num) > 1 else "None"]
---
# Task 003 — Version JSON Reader

## 1. Goal
Read promptArchitectByRiseupAsia section from version.json in `cmd/prompt_version_reader.go`.
