---
plan: .lovable/plans/pending/06-prompt-architect-installer.md
domain: Core
phase: Implement
target_files: ["gitmap/installer/prompt_bash_checker.go"]
depends_on: ["Task 016" if int(num) > 1 else "None"]
---

# Task 017 — Curl / Bash Availability Checker

## 1. Goal

Check if curl and bash are available on PATH in `installer/prompt_bash_checker.go`.
