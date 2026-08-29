---
plan: .lovable/plans/pending/06-prompt-architect-installer.md
domain: Core
phase: Implement
target_files: ["gitmap/cmd/prompt_batch_progress.go"]
depends_on: ["Task 026" if int(num) > 1 else "None"]
---

# Task 027 — Batch Progress Tracker Integration

## 1. Goal

Track batch progress across multi-repo installation in `cmd/prompt_batch_progress.go`.
