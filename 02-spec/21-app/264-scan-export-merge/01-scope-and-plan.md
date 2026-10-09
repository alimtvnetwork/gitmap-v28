# 264 — `scan export` / `scan merge`: portable repo-set JSON (scope and plan)

## 1. Problem statement

`gitmap scan` keeps a repo cache in SQLite, but the cache is not portable: there
is no way to take the scanned repo set to another machine and clone exactly
those repos. The user wants:

1. `scan export`: dump the repo cache to a portable JSON file inside a folder
   named by machine slug (`<machine-slug>/repos.json`), so the folder can be
   copied to any machine and used to clone precisely those repos.
2. `scan merge`: combine several such export folders into one deduped JSON,
   then clone everything in one place.

Research (2026-10-09) established two-thirds already exists: `gitmap scan`
writes `gitmap.json` (scan-record array) and `gitmap clone-from <file>` (alias
`cf`) batch-clones a JSON plan (dry-run default, `--execute` to run). The
genuinely new work is `scan export` and `scan merge`. No `scan export` /
`scan merge` exists today (full-repo search, 0 hits).

## 2. Design

### 2.1 `gitmap scan export [--machine <name>] [--out <dir>]`

- Read path (no rescan): `store.OpenDefault()` → `(*DB).ListRepos()`
  (`cli/store/repo.go:91`) returns `[]model.ScanRecord`.
- Machine name: `--machine` flag, else `os.Hostname()` with `"localhost"`
  fallback (pattern at `cli/cmdssh/ssh_export_import_bundle.go:28`).
- Slug: `store.SanitizeSlug(machine)` (`cli/store/split_db_path.go:25`).
- Output: `mkdir -p <out>/<slug>/`, write `repos.json` (indented
  `[]model.ScanRecord` — the exact shape `clone-from` consumes).
- `--out` defaults to cwd. Print: machine, slug, repo count, file path.

### 2.2 `gitmap scan merge <dir>... [--out <file>]`

- Each positional arg is an export folder (reads `<dir>/repos.json`; a direct
  `.json` file path is also accepted).
- Dedupe key per record: `httpsUrl`, else `sshUrl`, else `discoveredUrl`
  (first occurrence wins, order preserved).
- Write merged array to `--out` (default `./merged-repos.json`).
- Clone side needs **no new code**: `gitmap clone-from <merged.json> --execute`
  (or drop it in cwd as `gitmap.json` and run bare `gitmap clone`, which
  auto-discovers the manifest).

### 2.3 Subcommand interception

`gitmap scan` has no subcommands today — `scan export` would parse as
`dir="export"`. Mirror clone's pattern (`isCloneSubcommand` /
`resolveCloneSubcommand` / `dispatchCloneSubcommand` at
`cli/cmdclone/clone.go:89-139`): intercept `export` / `merge` at the top of
`RunScan` (`cli/cmdscan/exports.go:22`) before `ParseScanFlags`, with
`--help` handling per verb. Known precedent limitation (same as clone's
subcommands): a directory literally named `export`/`merge` can no longer be
scanned by bare position — use `gitmap scan ./export` in that case.

### 2.4 Explicitly out of scope

- Any change to `clone-from`, `scan`, `rescan` behavior.
- CSV/other export formats (JSON only; it is the clone-plan format).

## 3. Implementation plan

1. Spec (this dir).
2. New file `cli/cmdscan/scan_export_merge.go` (<300 lines): subcommand
   detection, `RunScanExport`, `RunScanMerge`, help texts; hook into `RunScan`.
3. Build (`go build ./...`), then E2E in temp dirs:
   export (default + `--machine`) → inspect JSON → merge two exports (with an
   overlapping repo) → `clone-from --dry-run` on the merged file →
   `gitmap scan <dir>` regression (normal scan still works).
4. Skills: document in `.agents/skills/gitmap/SKILL.md` + `.cursor` mirror.
5. Commit + push (atomic, hyphen-format).
## 4. Known issue (out of scope)

`03-ai-scripts/51-helptext-generator.py` is broken at HEAD: its embedded Go
runner imports `cli/tool/helptextemitter`, which does not exist in the repo
(likely lost in the wave-2 cmd split). The generator fails identically on the
pristine tree, so this is pre-existing and unrelated to this spec.
`cli/helpdoc/scan.md` could not be regenerated; `cmdscan.ScanHelpDisplay`
already carries the new `export`/`merge` verbs and usage lines, so the next
successful regen picks them up. Fixing the generator is a separate task.

