---
plan: .lovable/plans/pending/06-prompt-architect-installer.md
domain: Core
phase: Implement
target_files: ["gitmap/cmd/prompt_child_repos.go"]
depends_on: ["Task 022" if int(num) > 1 else "None"]
---
# Task 023 — Child Git Repos Auto-Discovery

## 1. Goal
Discover child git repositories for batch prompt installation in `cmd/prompt_child_repos.go`.
