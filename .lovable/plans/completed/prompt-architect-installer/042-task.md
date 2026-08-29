---
plan: .lovable/plans/pending/06-prompt-architect-installer.md
domain: Core
phase: Implement
target_files: ["gitmap/cmd/roottooling.go"]
depends_on: ["Task 041" if int(num) > 1 else "None"]
---

# Task 042 — Top-Level install-prompts Alias

## 1. Goal

Wire install-prompts and ct into root router in `cmd/roottooling.go`.
