---
plan: .lovable/plans/pending/05-workdir-pull-table-dirty-remedy.md
domain: Core
phase: Implement
target_files: ["gitmap/gitutil/last_sha.go"]
depends_on: ["Task 020" if int(num) > 1 else "None"]
---

# Task 021 — Git Short SHA Resolver

## 1. Goal

Extract 7-character short commit hash from HEAD in `gitutil/last_sha.go`.
