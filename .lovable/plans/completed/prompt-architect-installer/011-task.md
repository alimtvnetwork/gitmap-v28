---
plan: .lovable/plans/pending/06-prompt-architect-installer.md
domain: Core
phase: Implement
target_files: ["gitmap/installer/prompt_runner_unix.go"]
depends_on: ["Task 010" if int(num) > 1 else "None"]
---

# Task 011 — Unix Bash Remote Script Runner

## 1. Goal

Execute remote bash script on Unix/macOS/Linux in `installer/prompt_runner_unix.go`.
