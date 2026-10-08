# Subtask 01 — Consolidate cmdpull + cmdui shards (~300-line concern groups)

## Objective
Wave D over-split `cli/cmdpull/pull.go` (1949 lines) into 14 shards and `cli/cmdui/ui_assets.go` (1179 lines) into 13 shards.
Regroup them into ~300-line concern-grouped files. Behavior byte-identical: no signature or logic changes, purely file-level regrouping.

## Read first
`.ai-memory/what-to-read.md` — "The Mindset" section (point 5: ~300 lines max per file, split by concern, never fragment into tiny shards).

## Owned files (ONLY these — do not touch any other file in these packages)
Source shards to read, regroup, then remove (removal explicitly authorized for this consolidation — nothing beyond it):
- cli/cmdpull/pull_arg_normalize.go, cli/cmdpull/pull_batch.go, cli/cmdpull/pull_batch_output.go, cli/cmdpull/pull_batch_remediation.go, cli/cmdpull/pull_cwd.go, cli/cmdpull/pull_dispatch.go, cli/cmdpull/pull_dispatch_routes.go, cli/cmdpull/pull_execute.go, cli/cmdpull/pull_flags.go, cli/cmdpull/pull_ignore_remediation.go, cli/cmdpull/pull_ignore_scan.go, cli/cmdpull/pull_targets.go, cli/cmdpull/pull_targets_cwd.go, cli/cmdpull/pull_transport.go
- cli/cmdui/ui_assets_head.go, cli/cmdui/ui_assets_js_commit.go, cli/cmdui/ui_assets_js_core.go, cli/cmdui/ui_assets_js_fleet.go, cli/cmdui/ui_assets_js_nodes.go, cli/cmdui/ui_assets_js_settings.go, cli/cmdui/ui_assets_js_terminal.go, cli/cmdui/ui_assets_markup_misc.go, cli/cmdui/ui_assets_markup_ops.go, cli/cmdui/ui_assets_markup_settings_a.go, cli/cmdui/ui_assets_markup_settings_b.go, cli/cmdui/ui_assets_markup_shell.go, cli/cmdui/ui_assets_tail.go
Target files to create (all paths relative to repo root):
- cli/cmdpull/pull_cli_input.go, cli/cmdpull/pull_batch.go, cli/cmdpull/pull_batch_output.go, cli/cmdpull/pull_cwd.go, cli/cmdpull/pull_dispatch.go, cli/cmdpull/pull_routes.go, cli/cmdpull/pull_targets.go, cli/cmdpull/pull_ignore_remediation.go, cli/cmdpull/pull_ignore_scan.go
- cli/cmdui/ui_assets_js_fleet.go, cli/cmdui/ui_assets_js_nodes.go, cli/cmdui/ui_assets_js_terminal.go, cli/cmdui/ui_assets_markup_misc.go, cli/cmdui/ui_assets_markup_ops.go, cli/cmdui/ui_assets_markup_settings.go

## Target groupings (shard → target; adjust only if a target would exceed 300 lines)
cmdpull (2224 lines → 9 files):
- pull_cli_input.go (263): pull_arg_normalize.go + pull_flags.go + pull_transport.go
- pull_batch.go (193): pull_batch.go + pull_batch_remediation.go
- pull_batch_output.go (300): pull_batch_output.go (unchanged)
- pull_cwd.go (297): pull_cwd.go + pull_targets_cwd.go
- pull_dispatch.go (287): pull_dispatch.go + pull_execute.go
- pull_routes.go (129): pull_dispatch_routes.go (unchanged)
- pull_targets.go (223): pull_targets.go (unchanged)
- pull_ignore_remediation.go (211): unchanged
- pull_ignore_scan.go (163): unchanged
cmdui (1289 lines → 6 files):
- ui_assets_js_fleet.go (192): ui_assets_js_core.go + ui_assets_js_commit.go + ui_assets_js_fleet.go
- ui_assets_js_nodes.go (181): ui_assets_js_nodes.go + ui_assets_js_settings.go
- ui_assets_js_terminal.go (124): unchanged
- ui_assets_markup_misc.go (240): ui_assets_head.go + ui_assets_tail.go + ui_assets_markup_shell.go + ui_assets_markup_misc.go
- ui_assets_markup_ops.go (278): ui_assets_markup_ops.go + ui_assets_markup_settings_a.go
- ui_assets_markup_settings.go (203): ui_assets_markup_settings_b.go (renamed)

## Method
1. For each target: read the listed shards fully (via `gitmap cat`), write ONE file with a single `package` clause at top followed by the concatenated bodies (strip duplicate package clauses and duplicate file-header comments; keep all code, comments, and blank-line structure otherwise byte-identical).
2. Verify no duplicate top-level symbol definitions within each merged file (shards came from one original file, so there should be none — if you find any, STOP and report BLOCKED).
3. Verify each target file is ≤300 lines (`wc -l`).
4. Verify target filenames do not collide with pre-existing files in the package (checked 2026-10-09: no collisions — re-verify before writing).
5. Only after ALL targets in a package are written and verified: remove that package's shard files.

## Done criteria
- [ ] 9 cmdpull files + 6 cmdui files created, each ≤300 lines, concern-grouped.
- [ ] All 27 source shards removed; no other files touched.
- [ ] Report lists every created file with its line count.

## Hard rules
- Search ONLY via `gitmap aum search` / `gitmap find` / `gitmap cat`. TOTAL BAN on grep, rg, git grep, Select-String.
- NEVER run git commands. NEVER run `go build` / `go test` (lead builds centrally).
- Relative paths only. Python only via `gitmap py` (not needed here).
- Do NOT prune aliases. Do NOT touch versioning.
