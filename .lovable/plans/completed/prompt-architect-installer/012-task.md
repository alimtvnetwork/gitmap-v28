---
plan: .lovable/plans/pending/06-prompt-architect-installer.md
domain: Core
phase: Implement
target_files: ["gitmap/installer/prompt_runner_windows.go"]
depends_on: ["Task 011" if int(num) > 1 else "None"]
---

# Task 012 — Windows PowerShell Remote Script Runner

## 1. Goal

Execute remote PowerShell script on Windows in `installer/prompt_runner_windows.go`.
