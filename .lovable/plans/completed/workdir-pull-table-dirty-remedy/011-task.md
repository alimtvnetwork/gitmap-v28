---
plan: .lovable/plans/pending/05-workdir-pull-table-dirty-remedy.md
domain: Core
phase: Implement
target_files: ["gitmap/fsutil/recursive_top_level.go"]
depends_on: ["Task 010" if int(num) > 1 else "None"]
---

# Task 011 — Recursive Top-Level Scanner Core

## 1. Goal

Scan directories and prune at top-level .git in `fsutil/recursive_top_level.go`.
