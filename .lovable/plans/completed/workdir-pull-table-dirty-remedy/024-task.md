---
plan: .lovable/plans/pending/05-workdir-pull-table-dirty-remedy.md
domain: Core
phase: Implement
target_files: ["gitmap/gitutil/upstream_delta.go"]
depends_on: ["Task 023" if int(num) > 1 else "None"]
---

# Task 024 — Upstream Tracking Delta Calculator

## 1. Goal

Calculate ahead and behind counts relative to upstream in `gitutil/upstream_delta.go`.
