---
plan: .lovable/plans/pending/05-workdir-pull-table-dirty-remedy.md
domain: Core
phase: Implement
target_files: ["gitmap/version.json"]
depends_on: ["Task 048" if int(num) > 1 else "None"]
---

# Task 049 — Version Minor Bump

## 1. Goal

Bump minor version in version.json in `version.json`.
