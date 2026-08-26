---
plan: .lovable/plans/pending/05-workdir-pull-table-dirty-remedy.md
domain: Core
phase: Implement
target_files: ["gitmap/cmd/push_target_resolver.go"]
depends_on: ["Task 016" if int(num) > 1 else "None"]
---
# Task 017 — Push Target Resolver Upgrade

## 1. Goal
Upgrade push target resolver with top-level discovery in `cmd/push_target_resolver.go`.
