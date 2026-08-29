---
plan: .lovable/plans/pending/05-workdir-pull-table-dirty-remedy.md
domain: Core
phase: Implement
target_files: ["gitmap/cmd/pull_target_resolver.go"]
depends_on: ["Task 015" if int(num) > 1 else "None"]
---

# Task 016 — Pull Target Resolver Upgrade

## 1. Goal

Upgrade pull target resolver with top-level discovery in `cmd/pull_target_resolver.go`.
