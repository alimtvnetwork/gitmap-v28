# Plan 271 — Improvement Program: workflow fixes + nc auto-discovery + dead code + deprecation enforcement + release

Spec: `02-spec/21-app/271-improvement-program/` (to be written by WS leads)
Backup: `backup/271-improvement-program` (pushed to origin)

## Wave 1 — WS1: Workflow bug fixes ✅ (Worker A, commits c6c26a0c5, lead-verified build exit 0)
- [x] (a) Dedup cache: gate 'already built' cache write on test success
- [x] (b) E2E short-circuit: fix working-directory assumption on re-runs
- [x] Verify no other workflow has same pattern (aum search)
- [x] YAML syntax validated, committed + pushed

## Wave 4 — WS4: Deprecation enforcement ✅ (Worker A, commit 585471f71, lead-verified build exit 0)
- [x] `Deprecated` marker mechanism (stderr warning + replacement + removal version, still executes)
- [x] Wired into runDispatchTable + demo tests (no renamed top-level command exists; used sanctioned fallback)
- [x] Committed + pushed

## Wave 2 — WS2: nc auto-discovery ✅ (Worker B, commit 0941e3e7c, lead-verified: 916 discovered, spec next correct)
- [x] Discovery from dispatcher registration (not hardcoded catalog)
- [x] Flags/format/`nc` alias unchanged; hardcoded entries as fallback only
- [x] `gitmap nc --limit 5`, `gitmap nc --json` verified; `spec next` + `te all` correct
- [x] Committed + pushed

## Wave 3 — WS3: Dead code removal ✅ (Worker B, commit d1ab49d88, lead-verified build exit 0)
- [x] Recount `unused` findings: 36 → 2 (2 build-tag shims deliberately left + noted)
- [x] Removed 34 dead funcs/vars/types across 13 files; no files deleted
- [x] Lead-verified `go build ./...` exit 0
- [x] Committed + pushed

## Wave 4 — WS4: Deprecation enforcement ✅ (Worker A, commit 585471f71, lead-verified build exit 0)
- [x] `Deprecated` marker mechanism (stderr warning + replacement + removal version, still executes)
- [x] Wired into runDispatchTable + demo tests (no renamed top-level command exists; used sanctioned fallback)
- [x] Committed + pushed

## Wave 5 — WS5: Release ✅
- [x] All work committed + pushed
- [x] Minor bump to 6.551.0; tag v6.551.0 + pushed
- [x] GitHub release published: https://github.com/alimtvnetwork/gitmap-v28/releases/tag/v6.551.0
- [x] ~/.local/bin/gitmap rebuilt → v6.551.0 verified
- [ ] `gitmap pe -t 900` CI wait — in progress
