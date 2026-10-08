# Subtask 04 — Consolidate dbengine + cmdchromeprofile shards, file-size gate 500→300

## Objective
Wave D over-split `cli/dbengine/dbengine_test.go` + `cli/dbengine/query.go` (1474 + 1147 lines) into 14 shards and `cli/cmdchromeprofile/chromeprofile_smart_import.go` (1214 lines) into 8 shards.
Regroup them into ~300-line concern-grouped files. Behavior byte-identical: no signature or logic changes, purely file-level regrouping.
Also lower the file-size gate in `03-ai-scripts/52-file-size-check.py` from 500 to 300 lines (owner's rule).

## Read first
`.ai-memory/what-to-read.md` — "The Mindset" section (point 5: ~300 lines max per file, split by concern, never fragment into tiny shards).

## Owned files (ONLY these — do not touch any other file)
Source shards to read, regroup, then remove (removal explicitly authorized for this consolidation — nothing beyond it):
- cli/dbengine/dbengine_compiler_test.go, cli/dbengine/dbengine_query_builder_test.go, cli/dbengine/dbengine_query_joins_test.go, cli/dbengine/dbengine_repository_test.go, cli/dbengine/dbengine_test_helpers_test.go, cli/dbengine/dbengine_transaction_modes_test.go, cli/dbengine/dbengine_transaction_test.go, cli/dbengine/dbengine_types_test.go, cli/dbengine/query_build_clauses.go, cli/dbengine/query_build_select.go, cli/dbengine/query_builder.go, cli/dbengine/query_builder_compile.go, cli/dbengine/query_builder_filters.go, cli/dbengine/query_join.go
- cli/cmdchromeprofile/chromeprofile_bookmarks.go, cli/cmdchromeprofile/chromeprofile_destination.go, cli/cmdchromeprofile/chromeprofile_import_check.go, cli/cmdchromeprofile/chromeprofile_import_flow.go, cli/cmdchromeprofile/chromeprofile_import_run.go, cli/cmdchromeprofile/chromeprofile_snapshot_email.go, cli/cmdchromeprofile/chromeprofile_snapshot_files.go, cli/cmdchromeprofile/chromeprofile_snapshot_meta.go
File to edit (threshold change only):
- 03-ai-scripts/52-file-size-check.py — change the 500-line limit to 300 (limit constant, messages, comments, docstring — every reference).

Target files to create (all paths relative to repo root):
- cli/dbengine/dbengine_compiler_test.go, cli/dbengine/dbengine_query_joins_test.go, cli/dbengine/dbengine_transaction_test.go, cli/dbengine/dbengine_repository_test.go, cli/dbengine/dbengine_types_test.go, cli/dbengine/dbengine_test_helpers_test.go, cli/dbengine/query_build_clauses.go, cli/dbengine/query_builder.go, cli/dbengine/query_builder_compile.go, cli/dbengine/query_builder_filters.go
- cli/cmdchromeprofile/chromeprofile_import_check.go, cli/cmdchromeprofile/chromeprofile_import_run.go, cli/cmdchromeprofile/chromeprofile_snapshot_meta.go, cli/cmdchromeprofile/chromeprofile_snapshot_files.go, cli/cmdchromeprofile/chromeprofile_destination.go, cli/cmdchromeprofile/chromeprofile_bookmarks.go

## Target groupings (shard → target; adjust only if a target would exceed 300 lines)
dbengine (2413 lines → 10 files):
- dbengine_compiler_test.go (261): dbengine_compiler_test.go + dbengine_query_builder_test.go
- dbengine_query_joins_test.go (206): unchanged
- dbengine_transaction_test.go (185): unchanged
- dbengine_repository_test.go (274): unchanged
- dbengine_types_test.go (215): unchanged
- dbengine_test_helpers_test.go (247): dbengine_test_helpers_test.go + dbengine_transaction_modes_test.go
- query_build_clauses.go (294): unchanged
- query_builder.go (287): query_build_select.go + query_builder.go
- query_builder_compile.go (278): query_builder_compile.go + query_join.go
- query_builder_filters.go (178): unchanged
cmdchromeprofile (1346 lines → 6 files):
- chromeprofile_import_check.go (250): chromeprofile_import_check.go + chromeprofile_import_flow.go
- chromeprofile_import_run.go (230): unchanged
- chromeprofile_snapshot_meta.go (161): unchanged
- chromeprofile_snapshot_files.go (298): chromeprofile_snapshot_files.go + chromeprofile_snapshot_email.go
- chromeprofile_destination.go (238): unchanged
- chromeprofile_bookmarks.go (69): unchanged
file-size gate:
- 03-ai-scripts/52-file-size-check.py: 500 → 300 everywhere (read the file first; update limit, messages, comments).

## Method
1. For each target: read the listed shards fully (via `gitmap cat`), write ONE file with a single `package` clause at top followed by the concatenated bodies (strip duplicate package clauses and duplicate file-header comments; keep all code, comments, and blank-line structure otherwise byte-identical).
2. Verify no duplicate top-level symbol definitions within each merged file (shards came from one original file, so there should be none — if you find any, STOP and report BLOCKED).
3. Verify each target file is ≤300 lines (`wc -l`).
4. Verify target filenames do not collide with pre-existing files in the package (checked 2026-10-09: no collisions — re-verify before writing).
5. Only after ALL targets in a package are written and verified: remove that package's shard files.
6. Edit 52-file-size-check.py via `gitmap py` (read it first).

## Done criteria
- [ ] 10 dbengine files + 6 cmdchromeprofile files created, each ≤300 lines, concern-grouped.
- [ ] All 22 source shards removed; no other files touched.
- [ ] 52-file-size-check.py enforces 300 lines (no remaining 500 references).
- [ ] Report lists every created file with its line count.

## Hard rules
- Search ONLY via `gitmap aum search` / `gitmap find` / `gitmap cat`. TOTAL BAN on grep, rg, git grep, Select-String.
- NEVER run git commands. NEVER run `go build` / `go test` (lead builds centrally).
- Relative paths only. Python ONLY via `gitmap py`.
- Do NOT prune aliases. Do NOT touch versioning.
