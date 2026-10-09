# Plan 250 — Package consolidation + stack-trace setting + stdout writer refactor + guideline reuse

## User request (verbatim)
Owner (2026-10-09): reduce package count by grouping logically-together packages; keep packages small (~300-line files) but consolidated; mine coding-guidelines Go packages for reuse; stack-trace output as a user setting (default ON — it's for developers); kill the global stdout pipe hack.

## Scope
- Repo: `~/workspace/repos/gitmap-v28`, module `github.com/alimtvnetwork/gitmap-v28/cli`. Main at 9f75232 (pulled 2026-10-09).
- Backup branch `backup/250-package-consolidation` created from 9f75232 and pushed to origin BEFORE any changes (git-direct; `space backup-branch` refused on dirty tree from parallel 246 stream's ledger file — left untouched).
- Spec: `02-spec/21-app/250-package-consolidation/` (4 files). Subtasks: `.ai-memory/plans/subtasks/250-package-consolidation/` (4 files).

## Decisions (from spec 01-overview.md)
- D1: merge, never delete logic. D2: orphan triage before merges. D3: foundation-first order. D4: hubs untouched (constants/apperror/store/cliexit/model/cmdssh/cmdagy). D5: behavior identical, lead-run build per wave. D6: aliases preserved where external importers exist.

## Wave plan
- **Wave 1** (2 workers): Worker 1 = subtask 01 orphan triage (20 orphans: merge-or-keep verdicts; NOTHING deleted without owner ask) + subtask 03 stack-trace setting (`model.Config.ShowStackTrace`, default true, gate in `handleGlobalError`). Worker 2 = standby verification prep (reads merge map, pre-stages nothing; assists on triage evidence). → lead: `go build ./...`, commit `gitmap cpf`.
- **Wave 2** (2 workers, disjoint): Worker 1 = groups 1 (diag), 2 (jsonx), 10 (helpdoc). Worker 2 = groups 3 (termout), 4 (secrets), 5 (fspath). → lead: build, cycle check via `go list`, commit.
- **Wave 3** (2 workers, disjoint): Worker 1 = groups 6 (scanpipe), 7 (committransfer), 8 (cmdcommitin). Worker 2 = group 9 (24 dispatcher micros → `cli/cmd` as `cmd_*.go`). → lead: build, commit.
- **Wave 4** (2 workers): Worker 1 = subtask 04 stdout writer refactor (new `cli/output` package, synchronous FilterWriter, Install delegation, context plumbing, `output.Raw()`, cliexit migration; highest risk). Worker 2 = final sweep (empty dirs, leftover references, targeted linters). → lead: build + fresh-binary smoke (colors, glyph filter, cat byte-faithful, error path), commit.

## Checkboxes
- [x] Step 0: git pull (9f75232), backup branch `backup/250-package-consolidation` pushed.
- [x] Research 01 (package graph, 313 pkgs, 0 cycles) + Research 02 (stdout/settings/guidelines) DONE.
- [x] Spec 250 (4 files) + subtasks (4 files) DONE.
- [x] Wave 1: orphan triage + stack-trace setting. Build exit 0 (lead-run). Committed 390cc2f, pushed. NOTE: `gitmap cpf` auto-add swept the parallel 246 stream's ledger file into the commit despite targeted `git add` — content is 246's own (program complete), harmless; lesson re-logged.
- [x] Wave 2: foundation merges (groups 1,2,3,4,5,10) DONE. 6 new packages: diag, jsonx, termout, secrets, fspath, helpdoc. Corrections: appfault stays standalone (cycle); jsonenv.Envelope→EnvEnvelope; helpdoc/docs/cmd kept as subpackage; termout filename clashes resolved via pkg-prefix rename (lead fallback); fileExists collision (fspath) resolved via rename (both were dead); 51-helptext-generator.py out-path → cli/helpdoc. 302 importer files rewritten (scripted, dedup handled). Build exit 0 (lead-run), 0 import cycles. 313→292 package dirs.
- [ ] Wave 2: foundation merges (groups 1,2,3,4,5,10).
- [ ] Wave 3: middle + dispatcher merges (groups 6,7,8,9).
- [ ] Wave 4: stdout writer refactor + final sweep.
- [ ] Final: registers updated, atomic push per wave, completion report.

## Conflicts
- Skill R1 bans `go build`; PARENT explicitly overrides: lead runs `go build ./...` personally after every wave. Workers must NOT build/test; they verify via `go list` (read-only) and `gofmt -l` only. Owner instruction wins.
- Skill R8/R9 wants ONE atomic commit at Phase 3; PARENT overrides: atomic `gitmap cpf` per wave (established repo practice). Owner instruction wins.

## Assumptions
- Merge map file counts measured 2026-10-09; workers re-verify before moving.
- `apperror` NOT folded (191 importers); `cli/diag` takes appfault+suggestion+errreport.

## Follow-ups (out of scope)
- Splitting giant hubs (constants/store/cmdssh/cmdagy).
- Full call-site migration to explicit writers (follow-up program after 250).
- `errcmd` / `applogger` adoption (deferred per spec 04).
