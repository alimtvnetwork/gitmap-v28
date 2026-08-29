---
plan: .lovable/plans/pending/06-prompt-architect-installer.md
domain: Core
phase: Implement
target_files: ["gitmap/cmd/prompt_target_resolver.go"]
depends_on: ["Task 020" if int(num) > 1 else "None"]
---

# Task 021 — Single Target Path Resolver

## 1. Goal

Resolve single directory or repository target in `cmd/prompt_target_resolver.go`.
