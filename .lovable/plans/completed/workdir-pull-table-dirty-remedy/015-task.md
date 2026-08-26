---
plan: .lovable/plans/pending/05-workdir-pull-table-dirty-remedy.md
domain: Core
phase: Implement
target_files: ["gitmap/fsutil/git_dir_validator.go"]
depends_on: ["Task 014" if int(num) > 1 else "None"]
---
# Task 015 — Git Directory Safety Validator

## 1. Goal
Validate directory accessibility and .git integrity in `fsutil/git_dir_validator.go`.
