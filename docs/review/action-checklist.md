# Action Checklist — 30 Codebase Improvements

Work in priority order: P0 → P1 → P2. Check off with PR/commit link when done.
Format: `- [ ] ID — item (area)`

## P0 — Do first (correctness, trust, hygiene)

- [ ] P0-01 — Consolidate byte-identical `readme.md` / `what-to-read.md` (2,522 lines each) into one canonical file (docs)
- [ ] P0-02 — Fix `version.json` identity: Title/RepoSlug say "Coding Guidelines" / `coding-guidelines-v24` (versioning)
- [ ] P0-03 — Remove committed `cli/cli.exe` and `.syso` binaries; build in CI, ship via releases (hygiene)
- [ ] P0-04 — Audit `repo-secrets/` for committed secrets; purge and rotate anything sensitive (security)
- [ ] P0-05 — Verify the "Security Token Purge" plan is fully complete; grep state, logs, goldens for token leaks (security)
- [ ] P0-06 — Delete or archive root clutter: `fix_*.py`, `update_runner*.py`, stale `benchmark.md`, audit CSVs (hygiene)

## P1 — Consolidation and structure

- [ ] P1-01 — Merge `cloner`, `clonefrom`, `clonenext`, `clonenow`, `clonepick` behind one engine with strategy options (architecture)
- [ ] P1-02 — Shrink the 169KB readme into an index; move detail into `docs/commands/` (docs)
- [ ] P1-03 — Update `.ai-memory/overview.md` (says v3.1.0 / 60+ commands; actual v6.496.0) (docs)
- [ ] P1-04 — Fix "root readme must stay in sync with this file" note in `.ai-memory/what-to-read.md` (docs)
- [ ] P1-05 — Deduplicate 100+ `cli/constants_*.go` files; group or generate; check overlapping names (architecture)
- [ ] P1-06 — Document one source of truth per datum across `.gitmap/` state, `gitmap.json`, SQLite/split-DBs (architecture)
- [ ] P1-07 — Reconsider bump-on-every-change; adopt meaningful SemVer or date-based versions (versioning)
- [ ] P1-08 — Deduplicate version stamps (`version.json`, `package.json`, Go constants, badges) behind one verified sync script (versioning)
- [ ] P1-09 — Split 75KB `install.ps1` / 69KB `run.ps1` into sourced modules (detect, download, verify, install) (architecture)
- [ ] P1-10 — Audit install scripts for unnecessary elevation, network fetches, unpinned URLs (security)
- [ ] P1-11 — Stop committing multi-MB `test-inventory.json` / `test-heatmap.json`; move to ignored artifacts (hygiene)
- [ ] P1-12 — Add end-to-end smoke test for scan → manifest → parallel re-clone (testing)

## P2 — Quality, process, and follow-through

- [ ] P2-01 — Revisit no-negation / no-`switch` rules against real code samples; relax where Go idioms suffer (process)
- [ ] P2-02 — Enforce 200-line-file / 15-line-function limits with a linter check, not convention (process)
- [ ] P2-03 — Verify `setup.sh` hooks run `make lint test` pre-commit; wire into CI if not (process)
- [ ] P2-04 — Extend golden-fixture contract tests beyond the current subset (testing)
- [ ] P2-05 — Audit for bare `fmt.Errorf` / `os.Exit` paths outside `apperror` + `cliexit` (robustness)
- [ ] P2-06 — Confirm `govulncheck` (`vulncheck` Make target) gates PRs in CI (security)
- [ ] P2-07 — Rename `package.json` from `vite_react_shadcn_ts` to the real docs-site name (hygiene)
- [ ] P2-08 — Consolidate 100+ `cmd*` micro-packages into domain groups (architecture)
- [ ] P2-09 — Decide if the 70-page docs site belongs in this repo or a separate site repo (architecture)
- [ ] P2-10 — Triage `.ai-memory/pending-issues/` and `ambiguous-questions/`; close or schedule (process)
- [ ] P2-11 — Keep `docs/benchmarks/` current; add scanner perf-regression test on large trees (testing)
- [ ] P2-12 — Check binary startup cost; lazy-load heavy subsystems (TUI, SQLite) (robustness)

## Progress log

| Date | ID | Change | Link |
|---|---|---|---|
| | | | |
