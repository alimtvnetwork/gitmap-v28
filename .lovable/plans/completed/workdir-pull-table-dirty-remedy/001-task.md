---
plan: .lovable/plans/pending/05-workdir-pull-table-dirty-remedy.md
domain: Core
phase: Implement
target_files: ["gitmap/model/workdir.go"]
depends_on: ["Task 000" if int(num) > 1 else "None"]
---

# Task 001 — WorkDir Model Definition

## 1. Goal

Define WorkDir record struct in `model/workdir.go`.
