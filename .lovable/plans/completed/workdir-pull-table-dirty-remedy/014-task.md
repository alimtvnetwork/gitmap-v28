---
plan: .lovable/plans/pending/05-workdir-pull-table-dirty-remedy.md
domain: Core
phase: Implement
target_files: ["gitmap/fsutil/discovery_cache.go"]
depends_on: ["Task 013" if int(num) > 1 else "None"]
---

# Task 014 — Discovery Cache & Deduplicator

## 1. Goal

Cache and deduplicate discovered repository paths in `fsutil/discovery_cache.go`.
