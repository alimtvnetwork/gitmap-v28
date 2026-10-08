# Subtask 02 — Consolidate cmdpipeline shards (~300-line concern groups)

## Objective
Wave D over-split `cli/cmdpipeline/pipeline_logs.go` (1713 lines) into 17 shards.
Regroup them into ~300-line concern-grouped files. Behavior byte-identical: no signature or logic changes, purely file-level regrouping.

## Read first
`.ai-memory/what-to-read.md` — "The Mindset" section (point 5: ~300 lines max per file, split by concern, never fragment into tiny shards).

## Owned files (ONLY these — do not touch any other file in this package)
Source shards to read, regroup, then remove (removal explicitly authorized for this consolidation — nothing beyond it):
- cli/cmdpipeline/pipeline_extract_clean.go, cli/cmdpipeline/pipeline_extract_compact.go, cli/cmdpipeline/pipeline_extract_correlate.go, cli/cmdpipeline/pipeline_extract_format.go, cli/cmdpipeline/pipeline_extract_scan.go, cli/cmdpipeline/pipeline_extract_summarize.go
- cli/cmdpipeline/pipeline_logs_cards.go, cli/cmdpipeline/pipeline_logs_clipboard.go, cli/cmdpipeline/pipeline_logs_entry.go, cli/cmdpipeline/pipeline_logs_fetch.go, cli/cmdpipeline/pipeline_logs_help.go, cli/cmdpipeline/pipeline_logs_output.go, cli/cmdpipeline/pipeline_logs_payload.go, cli/cmdpipeline/pipeline_logs_render.go, cli/cmdpipeline/pipeline_logs_repo.go, cli/cmdpipeline/pipeline_logs_sections.go, cli/cmdpipeline/pipeline_logs_target.go
Target files to create (all paths relative to repo root):
- cli/cmdpipeline/pipeline_extract_clean.go, cli/cmdpipeline/pipeline_extract_compact.go, cli/cmdpipeline/pipeline_extract_correlate.go, cli/cmdpipeline/pipeline_extract_format.go, cli/cmdpipeline/pipeline_extract_scan.go, cli/cmdpipeline/pipeline_extract_summarize.go
- cli/cmdpipeline/pipeline_logs_entry.go, cli/cmdpipeline/pipeline_logs_target.go, cli/cmdpipeline/pipeline_logs_cards.go, cli/cmdpipeline/pipeline_logs_fetch.go, cli/cmdpipeline/pipeline_logs_repo.go, cli/cmdpipeline/pipeline_logs_output.go, cli/cmdpipeline/pipeline_logs_payload.go, cli/cmdpipeline/pipeline_logs_sections.go

## Target groupings (shard → target; adjust only if a target would exceed 300 lines)
extract group (1246 lines → 6 files, one concern each — unchanged):
- pipeline_extract_clean.go (256), pipeline_extract_compact.go (194), pipeline_extract_correlate.go (254), pipeline_extract_format.go (179), pipeline_extract_scan.go (186), pipeline_extract_summarize.go (177)
logs group (1854 lines → 8 files):
- pipeline_logs_entry.go (282): unchanged
- pipeline_logs_target.go (284): unchanged
- pipeline_logs_cards.go (259): pipeline_logs_cards.go + pipeline_logs_clipboard.go
- pipeline_logs_fetch.go (186): pipeline_logs_fetch.go + pipeline_logs_help.go
- pipeline_logs_repo.go (195): unchanged
- pipeline_logs_output.go (255): pipeline_logs_output.go + pipeline_logs_render.go
- pipeline_logs_payload.go (173): unchanged
- pipeline_logs_sections.go (151): unchanged

## Method
1. For each target: read the listed shards fully (via `gitmap cat`), write ONE file with a single `package` clause at top followed by the concatenated bodies (strip duplicate package clauses and duplicate file-header comments; keep all code, comments, and blank-line structure otherwise byte-identical).
2. Verify no duplicate top-level symbol definitions within each merged file (shards came from one original file, so there should be none — if you find any, STOP and report BLOCKED).
3. Verify each target file is ≤300 lines (`wc -l`).
4. Verify target filenames do not collide with pre-existing files in the package (checked 2026-10-09: no collisions — re-verify before writing).
5. Only after ALL targets are written and verified: remove the 17 shard files.

## Done criteria
- [ ] 14 cmdpipeline files created, each ≤300 lines, concern-grouped.
- [ ] All 17 source shards removed; no other files touched.
- [ ] Report lists every created file with its line count.

## Hard rules
- Search ONLY via `gitmap aum search` / `gitmap find` / `gitmap cat`. TOTAL BAN on grep, rg, git grep, Select-String.
- NEVER run git commands. NEVER run `go build` / `go test` (lead builds centrally).
- Relative paths only. Python only via `gitmap py` (not needed here).
- Do NOT prune aliases. Do NOT touch versioning.
