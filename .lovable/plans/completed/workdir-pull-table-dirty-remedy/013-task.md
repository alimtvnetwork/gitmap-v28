---
plan: .lovable/plans/pending/05-workdir-pull-table-dirty-remedy.md
domain: Core
phase: Implement
target_files: ["gitmap/fsutil/nested_filter.go"]
depends_on: ["Task 012" if int(num) > 1 else "None"]
---

# Task 013 — Nested Submodule Skip Predicates

## 1. Goal

Filter out submodules and nested child repos in `fsutil/nested_filter.go`.
