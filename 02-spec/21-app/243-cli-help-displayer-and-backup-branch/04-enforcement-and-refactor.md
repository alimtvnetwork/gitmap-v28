# Spec 243.4 — Enforcement, CMD Refactor, Enums, Stale Docs

## 1. File-size cap enforcement
The small-files rule exists on paper but 87 files exceed 500 lines and 10 exceed
1000 lines (verified 2026-10-08):
`cmdpull/pull.go` (1949), `cmdpipeline/pipeline_logs.go` (1713),
`cmdupdate/update_fleet.go` (1288), `cmdpipeline/pipeline_error_extract.go` (1269),
`cmdchromeprofile/chromeprofile_smart_import.go` (1214), `dbengine/query.go` (1147),
`cmd/clihelpers.go` (1129), plus 3 more over 1000.
- Split target: **~300 lines max per file, grouped by concern** — NOT fragmentation into tiny shards. 9 files at ~1000+ lines (≈9000 lines) → ~30 files, not 80+. Each split file owns one concern (e.g. pull.go → pull-flow, arg-normalization, gitignore-remediation as separate files). Functions are already tiny; this is purely file-level surgery.
- Then the 500–1000 band, largest first, same 300-line target.
- Enforcement mechanism: the check script in `03-ai-scripts/` (run via `gitmap py`) fails past **300** lines — wired into the repo's existing pre-commit path (staged-only, so grandfathered files can't block unrelated commits).
- New files: hard guidance — past ~300 lines, split by concern early.

## 2. Stale docs ("Stelldocs" fix)
- `.ai-memory/overview.md` still claims v3.1.0 — update to current (read `version.json`, don't hand-type).
- Write the missing state-ownership doc: which system owns what —
  `.gitmap/` files vs `gitmap.json` vs SQLite/split-DBs. One page under `docs/architecture/`
  (implementer to verify the docs location first).

## 3. CMD package split
`cli/cmd/` holds 741 files (206 tests) — a flat glue package undermining package-by-feature.
- Keep in `cmd/`: thin dispatch only — `root*.go` tables, argv preprocessing (`root.go`), alias context.
- Move every command implementation file into its `cmdX` package (most already exist:
  `cmdpull/`, `cmdscan`, …).
- Split `clihelpers.go` (1129 lines) by concern into the packages owning those concerns.
- Implementer: document the exact file-move map in the plan BEFORE moving anything;
  move in small waves with `go build` verification after each wave.

## 4. Enums package
Per the Go coding guidelines, enums live in their own package. Create `cli/enums/`
and migrate enum-like constant groups out of `cli/constants/` incrementally,
starting with the most self-contained group. Read the guideline's Go enum pattern
first and cite it in the plan.

## 5. New-command help technique
Every new command from now on uses the spec-03 `HelpDisplay` technique:
struct-first, generated `helptext/*.md`, single-touch registration. Update the
new-command checklist (implementer to locate it) to point at spec 03.
