---
plan: .lovable/plans/pending/06-prompt-architect-installer.md
domain: Core
phase: Implement
target_files: ["gitmap/cmd/prompt_version_validator.go"]
depends_on: ["Task 004" if int(num) > 1 else "None"]
---
# Task 005 — Version JSON Validation

## 1. Goal
Validate metadata schema and fallback keys in `cmd/prompt_version_validator.go`.
