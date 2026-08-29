---
plan: .lovable/plans/pending/06-prompt-architect-installer.md
domain: Core
phase: Implement
target_files: ["gitmap/installer/prompt_runner_dispatch.go"]
depends_on: ["Task 012" if int(num) > 1 else "None"]
---

# Task 013 — OS Dispatcher for Prompt Installation

## 1. Goal

Route execution based on runtime.GOOS in `installer/prompt_runner_dispatch.go`.
