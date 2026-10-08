# Plan 243b — Wave D fix (consolidation) + pending completion

## User request (verbatim)
Owner feedback 2026-10-09: the Wave D split of 9 files into 83 files is over-fragmented ("that is definitely wrong"). Rule: **~300 lines max per file, grouped by concern** (≈30–50 files for ~13,800 lines, not 83). Improve Wave D, `git pull` before work, build the code, end-to-end test every new/changed feature via the CLI, and complete the pending items (cmd-split waves 2–4, 500–1000-line band, HelpDisplay SampleOutput gap). Owner's overall summary requested at the end.

## Scope
- Repo: `~/workspace/repos/gitmap-v28`, main at f127f7a (pulled 2026-10-09, tree clean).
- Backup branch `backup/pre-consolidation-rework` created from f127f7a via dogfooded `gitmap space backup-branch "pre-consolidation-rework"` and pushed to origin.
- Mindset updated: point 5 now reads "~300 lines max per file … never fragment into tiny shards". Spec 243 §04 updated to the 300-line target (commit f127f7a).

## Conflicts
- V6 skill R1 bans `go build` / `go test`; the OWNER explicitly authorized builds + CLI-level end-to-end testing for this task. Top-Instruction Priority Mandate: owner instruction wins. `go test` suites stay off-limits; `go vet` + CLI E2E is the verification bar. Logged, not hidden.

## Wave plan
- Wave 1 — consolidation (this dispatch): subtasks 01–04, 83 shards → ~56 files, each ≤300 lines.
- Wave 2 — cmd-split waves 2–4 per the move map in `.ai-memory/plans/subtasks/243-cli-help-displayer-and-backup-branch/04-enforcement-refactor-docs.md` (lead validates each file before its wave; workers move + rename package + manifest caller rewires; lead rewires dispatch centrally).
- Wave 3 — 500–1000-line band (78 files → ~300-line concern groups) + HelpDisplay `SampleOutput` field gap.
- Wave 4 (lead) — `go build`, `go vet`, E2E per feature (space backup-branch lifecycle, help rendering diffs, scan, pull-all, space common --dry-run), fix misses, atomic `gitmap cpf` + push.

## Checkboxes
- [x] Step 0: git pull (f127f7a), binary built (/tmp/gitmap-e2e, exit 0), backup branch `backup/pre-consolidation-rework` pushed via dogfooded command.
- [x] `subtasks/243b-consolidation-and-completion/01-consolidate-cmdpull-cmdui.md` — Wave 1 DONE (27 shards → 15 files, all ≤300)
- [x] `subtasks/243b-consolidation-and-completion/02-consolidate-cmdpipeline.md` — Wave 1 DONE (17 shards → 14 files, all ≤300)
- [x] `subtasks/243b-consolidation-and-completion/03-consolidate-cmdupdate-pipelinedb.md` — Wave 1 DONE (17 shards → 11 files, all ≤300)
- [x] `subtasks/243b-consolidation-and-completion/04-consolidate-dbengine-cmdchromeprofile-filesize.md` — Wave 1 DONE (22 shards → 16 files, all ≤300; file-size gate 500→300)
- [x] Wave 1 verified: `go build ./...` exit 0, `go vet` exit 0 on all 7 touched packages (lead fixed 18 files with mid-file import blocks via import-hoist script — LEAD_FALLBACK).
- [x] `subtasks/243b-consolidation-and-completion/05-wave2-moves-g1.md` — Wave 2 DONE (G1: 117 moved; 2 BLOCKED by pre-existing dest files `cmdapps/apps.go`, `cmdsync/sync.go` — lead to move as `apps_dispatch.go`/`sync_dispatch.go` after symbol check = no collisions; 1 SKIPPED `bash_runner.go` absent)
- [x] `subtasks/243b-consolidation-and-completion/06-wave2-moves-g2.md` — Wave 2 DONE (G2: 128 moved; 1 SKIPPED `gitrm.go` absent)
- [x] `subtasks/243b-consolidation-and-completion/07-wave2-moves-g3.md` — Wave 2 DONE (G3: 214 files verified byte-clean by lead; worker 02 errored on runtime drain, completion verified independently)
- [x] `subtasks/243b-consolidation-and-completion/08-wave2-moves-g4.md` — Wave 2 DONE (G4 included in the 214)
- [x] Wave 2 build-fix (lead): 459 files moved; resolved ~150 undefined symbols via export-and-qualify + DI-hook pattern + per-package helper copies. `go build ./...` exit 0 (2026-10-09, commit 3f4f8d7, pushed). `go vet` clean on main code (test-file issues pre-existing, out of scope).
- [x] Wave 2: cmd-split waves 2–4 moved per move map; dispatch rewired; build clean.
- [ ] Wave 3: 500–1000 band split; HelpDisplay SampleOutput gap closed.
- [ ] Wave 4: E2E pass/fail per feature with evidence; atomic commit pushed.
- [ ] Overall summary + coordinator's honest take delivered to owner.

## Assumptions
- Shard files came from single originals → no duplicate symbols within a merge group (workers verify; BLOCKED if found).
- Moved cmd files keep filenames; destination collision re-check per wave (rule from move map).

## Follow-ups (out of scope)
- 300–500-line band (pre-commit gate ratchets it on touch; staged-only so it never blocks unrelated commits).
- cmd-split wave 4 keep-list is permanent (thin dispatch stays in `cli/cmd/`).
