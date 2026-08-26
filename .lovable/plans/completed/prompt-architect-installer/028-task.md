---
plan: .lovable/plans/pending/06-prompt-architect-installer.md
domain: Core
phase: Implement
target_files: ["gitmap/cmd/prompt_parallel_runner.go"]
depends_on: ["Task 027" if int(num) > 1 else "None"]
---
# Task 028 — Multi-Repo Parallel Runner

## 1. Goal
Run prompt installer concurrently across repositories in `cmd/prompt_parallel_runner.go`.
