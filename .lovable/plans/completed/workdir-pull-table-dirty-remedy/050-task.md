---
plan: .lovable/plans/pending/05-workdir-pull-table-dirty-remedy.md
domain: Core
phase: Implement
target_files: ["gitmap/tests/cmd_test/final_sanity_test.go"]
depends_on: ["Task 049" if int(num) > 1 else "None"]
---

# Task 050 — Final Sanity & Binary Build

## 1. Goal

Final compilation and validation check in `tests/cmd_test/final_sanity_test.go`.
