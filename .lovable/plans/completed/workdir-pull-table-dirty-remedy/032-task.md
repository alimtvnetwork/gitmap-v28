---
plan: .lovable/plans/pending/05-workdir-pull-table-dirty-remedy.md
domain: Core
phase: Implement
target_files: ["gitmap/gitutil/uncommitted_classifier.go"]
depends_on: ["Task 031" if int(num) > 1 else "None"]
---
# Task 032 — Uncommitted File Classifier

## 1. Goal
Classify modified, untracked, and deleted files in `gitutil/uncommitted_classifier.go`.
