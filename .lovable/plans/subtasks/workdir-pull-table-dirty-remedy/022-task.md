---
plan: .lovable/plans/pending/05-workdir-pull-table-dirty-remedy.md
domain: Core
phase: Implement
target_files: ["gitmap/gitutil/active_branch.go"]
depends_on: ["Task 021" if int(num) > 1 else "None"]
---
# Task 022 — Git Active Branch Resolver

## 1. Goal
Extract current active branch name or detached HEAD in `gitutil/active_branch.go`.
