---
plan: .lovable/plans/pending/05-workdir-pull-table-dirty-remedy.md
domain: Core
phase: Implement
target_files: ["gitmap/cmd/dir_scope_matcher.go"]
depends_on: ["Task 018" if int(num) > 1 else "None"]
---

# Task 019 — CLI Directory Scope Matcher

## 1. Goal

Match current directory against workdir scopes in `cmd/dir_scope_matcher.go`.
