---
plan: .lovable/plans/pending/06-prompt-architect-installer.md
domain: Core
phase: Implement
target_files: ["gitmap/installer/prompt_ps_checker.go"]
depends_on: ["Task 015" if int(num) > 1 else "None"]
---
# Task 016 — PowerShell Availability Checker

## 1. Goal
Check if powershell or pwsh is available on PATH in `installer/prompt_ps_checker.go`.
