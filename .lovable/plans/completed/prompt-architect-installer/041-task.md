---
plan: .lovable/plans/pending/06-prompt-architect-installer.md
domain: Core
phase: Implement
target_files: ["gitmap/cmd/ct.go"]
depends_on: ["Task 040" if int(num) > 1 else "None"]
---

# Task 041 — Top-Level ct Command Dispatcher

## 1. Goal

Implement gitmap ct [install-prompts|status|version] in `cmd/ct.go`.
