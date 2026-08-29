---
plan: .lovable/plans/pending/05-workdir-pull-table-dirty-remedy.md
domain: Core
phase: Implement
target_files: ["gitmap/gitutil/pr_detector.go"]
depends_on: ["Task 022" if int(num) > 1 else "None"]
---

# Task 023 — Git PR Tracking Detector

## 1. Goal

Detect upstream tracking and PR branch status in `gitutil/pr_detector.go`.
