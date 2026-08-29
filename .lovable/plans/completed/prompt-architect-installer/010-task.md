---
plan: .lovable/plans/pending/06-prompt-architect-installer.md
domain: Core
phase: Implement
target_files: ["gitmap/cmd/prompt_version_test.go"]
depends_on: ["Task 009" if int(num) > 1 else "None"]
---

# Task 010 — Unit Tests for Constants & Metadata

## 1. Goal

Unit tests for version.json prompt metadata in `cmd/prompt_version_test.go`.
