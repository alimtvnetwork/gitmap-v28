# Subtask 03 — Consolidate cmdupdate + pipelinedb shards (~300-line concern groups)

## Objective
Wave D over-split `cli/cmdupdate/update_fleet.go` (1288 lines) into 9 shards and `cli/pipelinedb/pipeline_split_ops.go` (1041 lines) into 8 shards.
Regroup them into ~300-line concern-grouped files. Behavior byte-identical: no signature or logic changes, purely file-level regrouping.

## Read first
`.ai-memory/what-to-read.md` — "The Mindset" section (point 5: ~300 lines max per file, split by concern, never fragment into tiny shards).

## Owned files (ONLY these — do not touch any other file in these packages)
Source shards to read, regroup, then remove (removal explicitly authorized for this consolidation — nothing beyond it):
- cli/cmdupdate/fleet_dispatch.go, cli/cmdupdate/fleet_execute.go, cli/cmdupdate/fleet_remote.go, cli/cmdupdate/fleet_remote_zip.go, cli/cmdupdate/fleet_ssh.go, cli/cmdupdate/fleet_summary.go, cli/cmdupdate/fleet_targets.go, cli/cmdupdate/fleet_telemetry.go, cli/cmdupdate/fleet_zip.go
- cli/pipelinedb/pipeline_split_clear.go, cli/pipelinedb/pipeline_split_errors.go, cli/pipelinedb/pipeline_split_info.go, cli/pipelinedb/pipeline_split_maintenance.go, cli/pipelinedb/pipeline_split_purge.go, cli/pipelinedb/pipeline_split_queries.go, cli/pipelinedb/pipeline_split_runs.go, cli/pipelinedb/pipeline_split_stats.go
Target files to create (all paths relative to repo root):
- cli/cmdupdate/fleet_dispatch.go, cli/cmdupdate/fleet_execute.go, cli/cmdupdate/fleet_remote.go, cli/cmdupdate/fleet_targets.go, cli/cmdupdate/fleet_zip.go, cli/cmdupdate/fleet_telemetry.go
- cli/pipelinedb/pipeline_split_cleanup.go, cli/pipelinedb/pipeline_split_maintenance.go, cli/pipelinedb/pipeline_split_info.go, cli/pipelinedb/pipeline_split_queries.go, cli/pipelinedb/pipeline_split_runs.go

## Target groupings (shard → target; adjust only if a target would exceed 300 lines)
cmdupdate (1383 lines → 6 files):
- fleet_dispatch.go (248): unchanged
- fleet_execute.go (213): fleet_execute.go + fleet_summary.go
- fleet_remote.go (292): fleet_remote.go + fleet_remote_zip.go + fleet_ssh.go
- fleet_targets.go (183): unchanged
- fleet_zip.go (207): unchanged
- fleet_telemetry.go (140): unchanged
pipelinedb (1031 lines → 5 files):
- pipeline_split_cleanup.go (199): pipeline_split_clear.go + pipeline_split_errors.go — NOTE: do NOT name this pipeline_split_ops.go; a pre-existing cli/pipelinedb/pipeline_split_ops.go (102 lines) already exists and must not be touched.
- pipeline_split_maintenance.go (184): pipeline_split_maintenance.go + pipeline_split_purge.go
- pipeline_split_info.go (229): pipeline_split_info.go + pipeline_split_stats.go
- pipeline_split_queries.go (185): unchanged
- pipeline_split_runs.go (194): unchanged

## Method
1. For each target: read the listed shards fully (via `gitmap cat`), write ONE file with a single `package` clause at top followed by the concatenated bodies (strip duplicate package clauses and duplicate file-header comments; keep all code, comments, and blank-line structure otherwise byte-identical).
2. Verify no duplicate top-level symbol definitions within each merged file (shards came from one original file, so there should be none — if you find any, STOP and report BLOCKED).
3. Verify each target file is ≤300 lines (`wc -l`).
4. Verify target filenames do not collide with pre-existing files in the package (checked 2026-10-09: no collisions except pipeline_split_ops.go, which is why the cleanup target uses a different name — re-verify before writing).
5. Only after ALL targets in a package are written and verified: remove that package's shard files.

## Done criteria
- [ ] 6 cmdupdate files + 5 pipelinedb files created, each ≤300 lines, concern-grouped.
- [ ] All 17 source shards removed; no other files touched.
- [ ] Report lists every created file with its line count.

## Hard rules
- Search ONLY via `gitmap aum search` / `gitmap find` / `gitmap cat`. TOTAL BAN on grep, rg, git grep, Select-String.
- NEVER run git commands. NEVER run `go build` / `go test` (lead builds centrally).
- Relative paths only. Python only via `gitmap py` (not needed here).
- Do NOT prune aliases. Do NOT touch versioning.
