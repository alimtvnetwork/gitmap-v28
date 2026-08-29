---
plan: .lovable/plans/pending/05-workdir-pull-table-dirty-remedy.md
domain: Core
phase: Implement
target_files: ["gitmap/fsutil/workdir_match.go"]
depends_on: ["Task 011" if int(num) > 1 else "None"]
---

# Task 012 — WorkDir Auto-Match Engine

## 1. Goal

Match active directory against registered work paths in `fsutil/workdir_match.go`.
