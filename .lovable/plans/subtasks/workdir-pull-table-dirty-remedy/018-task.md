---
plan: .lovable/plans/pending/05-workdir-pull-table-dirty-remedy.md
domain: Core
phase: Implement
target_files: ["gitmap/cmd/status_target_resolver.go"]
depends_on: ["Task 017" if int(num) > 1 else "None"]
---
# Task 018 — Status Target Resolver Upgrade

## 1. Goal
Upgrade status target resolver with top-level discovery in `cmd/status_target_resolver.go`.
