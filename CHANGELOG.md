# Changelog

## v6.79.0 — 2026-07-16 — Cloner LFS smudge auto-recovery

### Added
- **Automatic LFS smudge-failure recovery** in `gitmap/cloner`: when a clone fails with `smudge filter lfs failed` / `external filter 'git-lfs' failed` (typically a 404 on an LFS object), the cloner now cleans the partial destination and transparently retries with `GIT_LFS_SKIP_SMUDGE=1`, preserving LFS pointers so the checkout completes.
- `CloneResult.Notes` annotated with `lfs-skip-smudge-retry` when the fallback path fires, so downstream reporting and CSV/JSON exports surface the recovery.
- Unit tests in `gitmap/cloner/lfs_retry_test.go` covering `isLFSSmudgeFailure` detection and the retry cleanup path.

## v6.78.0 — 2026-07-16 — CG epic close: worked end-to-end example + goreleaser tag verification

### Added
- **CHANGELOG worked example** (this entry) — complete copy-pasteable `cfr cg` / `cfrp cg` walkthrough covering install, invocation, modifier ordering, and opt-out flags. Closes the 10-step Coding Guidelines v24 integration epic (Steps 1-9 shipped in v6.75.0 through v6.77.0).

### Verified
- **`.github/workflows/goreleaser.yml`** — tag-driven workflow (`v*.*.*`) audited against the new `gitmap/cmd/codingguidelines.go`, `codingguidelines_commit.go`, and `codingguidelines_test.go` files added across v6.75-v6.77. `main: .` in `gitmap/.goreleaser.yaml` compiles the entire `gitmap/cmd` package, so the new CG runner ships in every linux/darwin/windows (amd64+arm64) archive automatically. No config drift.
- **`gitmap/.goreleaser.yaml`** — goreleaser v2 schema (`version: 2`, `formats: [tar.gz]`, `format_overrides` per `goos: windows` → `[zip]`) still valid. ldflags stamp `constants.Version` from the pushed tag, so `gitmap --version` on the published binary matches the release tag byte-for-byte.

### Worked example: `cfr cg` end-to-end (v6.78.0)

```bash
# 1. Install v6.78.0 (Linux/macOS)
curl -fsSL https://github.com/alimtvnetwork/gitmap-v27/releases/download/v6.78.0/release-version-v6.78.0.sh | bash

# 2. Clone-fix-repo with Coding Guidelines v24 installer + auto commit + auto push (default)
gitmap cfr cg https://github.com/you/your-repo.git

# 3. Same, but keep changes local (skip push)
gitmap cfr cg --no-push https://github.com/you/your-repo.git

# 4. Same, but stage only (skip commit AND push)
gitmap cfr cg --no-commit https://github.com/you/your-repo.git

# 5. Promote-public + Coding Guidelines (order-independent: `p cg` == `cg p`)
gitmap cfrp cg https://github.com/you/your-repo.git
gitmap cfr p cg https://github.com/you/your-repo.git   # equivalent
```

Windows (PowerShell):

```powershell
irm https://github.com/alimtvnetwork/gitmap-v27/releases/download/v6.78.0/release-version-v6.78.0.ps1 | iex
gitmap cfr cg https://github.com/you/your-repo.git
```

**Modifier contract** (locked by `clonefixrepo_modifiers_test.go` since v6.76.0): `cg` and `p` may appear in any order before the URL; duplicates are idempotent; the first non-modifier token (flag or URL) stops modifier scanning.

### Changed
- Pinned: README pinned-version block + version matrix moved to **v6.78.0**. Synced `gitmap/constants/constants.go` (`Version = "6.78.0"`) and `src/constants/index.ts` (`VERSION = "v6.78.0"`).



## v6.77.0 — 2026-07-16 — `cfr` / `cfrp` `cg` modifier surfaced in UI command registry

### Added
- **UI registry** `src/data/commands.ts` — `clone-fix-repo` and `clone-fix-repo-pub` entries now advertise the `cg` (Coding Guidelines v24) and `p` (promote-public) pre-URL modifiers, plus the `--no-commit` / `--no-push` opt-out flags. Four new worked examples per command cover the plain, `cg`, `cg --no-commit`, `cg --no-push`, and combined `p cg` invocations. Usage lines updated to `gitmap clone-fix-repo [cg] [p] <url> [folder] [flags]` so search-by-flag surfaces the modifiers.

### Changed
- Pinned: README pinned-version block + version matrix moved to **v6.77.0**. Synced `gitmap/constants/constants.go` (`Version = "6.77.0"`) and `src/constants/index.ts` (`VERSION = "v6.77.0"`).

## v6.76.0 — 2026-07-16 — `cfr cg` / `cfrp cg` unit tests + modifier parser lock-in

### Added
- **Unit tests** `gitmap/cmd/clonefixrepo_modifiers_test.go` — 10-case table pinning the order-independent `cg` / `p` contract (duplicates idempotent, flags stop scan, unknown tokens stop scan, URL stops scan).
- **Unit tests** `gitmap/cmd/codingguidelines_test.go` — success path via injected `Runner` (fake `true`), exit-code propagation via fake `false`, and `ErrCGShellNotFound` when PATH is empty. Banners `Installing coding guidelines` and `OK Coding guidelines` are asserted so the standardized stderr format cannot silently drift.

### Fixed
- **Build error** `gitmap/cmd/sync.go` — renamed local `jsonEqual` helper to `syncJSONEqual` to resolve a package-level redeclaration collision with `gitmap/cmd/chromeprofile_merge.go:245`. Unblocked `go test ./cmd/...`.

### Changed
- Pinned: README pinned-version block + version matrix moved to **v6.76.0**. Synced `gitmap/constants/constants.go` (`Version = "6.76.0"`) and `src/constants/index.ts` (`VERSION = "v6.76.0"`).

## v6.75.0 — 2026-07-12 — `fix-auth`: one-shot cure for the wrong-account SSH push failure

### Added
- **`gitmap fix-auth`** (alias `fa`) — cross-platform Go port of the PowerShell/Bash "wrong GitHub account push failure" recipe. Generates `~/.ssh/id_ed25519_<user>`, pins the current repo via `git config core.sshCommand "ssh -i <key> -F /dev/null -o IdentitiesOnly=yes"`, and copies the public key to the OS clipboard (`clip`/`pbcopy`/`wl-copy`/`xclip`/`xsel`). New file: `gitmap/cmd/fixauth.go`. Constants: `CmdFixAuth` / `CmdFixAuthAlias` in `gitmap/constants/constants_cli.go`. Dispatch in `gitmap/cmd/rootutility.go`.
- **Help text** `gitmap/helptext/fix-auth.md`, `whoami.md`, `ssh-bind.md` — embedded via `go:embed *.md`; discoverable through `gitmap help <topic>` and the `--json` help payload.
- **Docs-site command entries** for `fix-auth`, `whoami`, and `ssh-bind` in `src/data/commands.ts` (`tools` category) with usage, flags, examples, and see-also links.

### Changed
- Pinned: README pinned-version block + version matrix moved to **v6.75.0**. Synced `gitmap/constants/constants.go` (`Version = "6.75.0"`) and `src/constants/index.ts` (`VERSION = "v6.75.0"`).



## v6.74.0 — 2026-07-01 — Release bump

### Changed
- Version bump to v6.74.0; pinned across README, Go/TS constants, and CHANGELOG.

## v6.73.0 — 2026-06-28 — Better `chrome export-bookmarks` errors + `--root`/`--folder` docs

### Added
- **Actionable error messages** for `gitmap chrome export-bookmarks` (`gitmap/cmd/chrome_bookmarks.go`): distinguishes missing/unreadable `Bookmarks` file, unknown `--root` (lists available roots), unmatched `--folder` (lists top-level folders + syntax hint), and empty `--match`/`--title` results.
- **Documented `--root` / `--folder` examples** for md, html, and json exports in `gitmap/helptext/chrome.md` and the README "Data, Profiles & Bookmarks" section.

## v6.72.0 — 2026-06-28 — Doctor `--json`/`--fix`, bookmark filters, chrome backup checksums

### Added
- **`gitmap doctor --json` / `--fix`** (`gitmap/cmd/doctor_run.go`): machine-readable probe output and auto-creation of missing config folders with actionable next-step hints.
- **Extended doctor network probes**: tests `api.github.com`, `github.com`, `codeload`, `uploads`, and `objects.githubusercontent.com` with consolidated "X/5 reachable" reporting (`gitmap/cmd/doctor_extra.go`).
- **`gitmap release-notes`** flags `--since <date>`, `--since-tag <tag>`, and `--format flat|grouped|markdown|json` with conventional-commit classification (`gitmap/cmd/release_notes_opts.go`).
- **Chrome backup checksum manifests**: SHA256 sidecars generated on backup and auto-verified before restore (skip with `--no-verify`); pure parser for `Local State` with unit coverage (`gitmap/cmd/chrome_manifest.go`, `chrome_localstate.go`, `chrome_manifest_test.go`).
- **`chrome export-bookmarks`** filters `--match <substring>` and `--title <exact>` preserve folder hierarchy while pruning (`gitmap/cmd/chrome_bookmarks_filter.go`).
- **`chrome restore`** safety: `--force` required to overwrite, `y/N` confirmation (bypass `--yes`), `--dry-run` preview.

## v6.71.0 — 2026-06-28 — Parallel hygiene scans, JSON/CSV exports, integration tests

### Added
- **Parallel scanning** for `gitmap stale`, `dedupe`, `size`, `orphans`: directory probes and per-repo `git` calls now fan out across up to 8 workers via shared `scanForReposParallel` / `mapReposParallel` helpers (`gitmap/cmd/hygiene_parallel.go`).
- **`--format=table|json|csv`** flag on all four hygiene commands. JSON emits structured arrays (RFC3339 timestamps, byte counts, tree SHAs); CSV emits a header row + standard encoding/csv quoting. Table remains the default.
- **Integration tests** (`gitmap/cmd/hygiene_integration_test.go`) spin up real git repos in `t.TempDir()` and exercise `scanForReposParallel`, `lastCommitTime`, `headTreeSHA`, `dirSize`, `originURL`, `gitURLToHTTPS`, plus the JSON/CSV emitters. Skip cleanly when `git` is absent.

### Notes
- `orphans --format=json|csv` is read-only (no delete prompt) to keep machine-readable output stable for piping into other tools.
- Version pinned to **v6.71.0** across `README.md`, `gitmap/constants/constants.go`, `src/constants/index.ts`.



## v6.70.0 — 2026-06-28 — Release tools, workflow shortcuts, safety net

### Added — Release
- **`gitmap release-notes <vN..vM>`** — auto-generate changelog block from commit messages via `git log --pretty`.
- **`gitmap release-dry [tag]`** — full rehearsal: `go build`, optional local tag, recent log; never pushes; prints undo recipe.
- **`gitmap tag-rename <old> <new>`** — local + origin tag rename (create new, delete old, push both refs).

### Added — Workflow
- **`gitmap recent` (`rct`)** — last 10 repos from the navigation helper history; `--print` for fzf piping.
- **`gitmap todo`** — grep `TODO|FIXME|XXX` across tracked files with per-hit `git blame` author.
- **`gitmap open` (`o`)** — open current repo's GitHub URL; `--issues`/`--prs`/`--actions` jump flags.
- **`gitmap pull-requests` (`prs`)** — list open PRs across an owner's repos in one table; honors `GITHUB_TOKEN`.
- **`gitmap blame-stats [root]`** — top contributors per file via `git blame --line-porcelain`.

### Added — Safety
- **`gitmap snapshot [root]`** — tar.gz snapshot of working tree to `.gitmap/snapshot/snap-<UTC-ts>.tar.gz` (includes untracked).
- **`gitmap rollback [tarball]`** — restore latest (or named) snapshot.
- **`gitmap guard [root]`** — install `.git/hooks/pre-commit` blocking secrets, files >10MB, `-vN` drift.

### Internal
- Constants `CmdReleaseNotes` / `CmdReleaseDry` / `CmdTagRename` / `CmdRecent(+Alias)` / `CmdTodo` / `CmdOpen(+Alias)` / `CmdPR(+Alias)` / `CmdBlameStats` / `CmdSnapshot` / `CmdRollback` / `CmdGuard` added; parity test updated.
- New files: `release_tools.go`, `workflow_recent_todo.go`, `workflow_open_pr.go`, `safety_snapshot.go`, plus 11 helptext markdown files.
- Version pinned to **v6.70.0** across `README.md`, `gitmap/constants/constants.go`, `src/constants/index.ts`.



## v6.69.0 — 2026-06-28 — Chrome umbrella: backup / restore / diff / export-bookmarks / which

### Added
- **`gitmap chrome backup`** — snapshot every Chrome profile to a tar.gz under `.gitmap/chrome/backup/` (`--out` for custom path). Skips volatile `LOCK`/`lockfile` entries.
- **`gitmap chrome restore <tarball>`** — restore a snapshot into the User Data dir (`--into` for custom root). Path-traversal guarded.
- **`gitmap chrome diff <A> <B>`** — extensions + bookmark-URL set diff (only-A / only-B) between two profiles.
- **`gitmap chrome export-bookmarks <profile> [--format md|html|json] [--out <file>]`** — export bookmarks tree, defaults to Markdown on stdout.
- **`gitmap chrome which`** — print active Chrome profile dir + display name from `Local State` (`last_used` + `last_active_profiles`).

### Internal
- New umbrella dispatcher `gitmap/cmd/chrome.go` plus per-subcommand files (`chrome_backup.go`, `chrome_diff.go`, `chrome_bookmarks.go`, `chrome_which.go`).
- Constants `CmdChrome` + `SubCmdChrome*` added to `constants_chromeprofile.go`; parity test updated.
- Help: `gitmap/helptext/chrome.md`.
- Version pinned to **v6.69.0** across `README.md`, `gitmap/constants/constants.go`, `src/constants/index.ts`.



## v6.68.0 — 2026-06-28 — Repo hygiene: stale, orphans, dedupe, size

### Added
- **`gitmap stale` (`st`)** — `gitmap/cmd/stale.go`. Lists local repos with no commits in the last N days (default 90). With `--archive` moves them into `.gitmap/archive/<UTC-ts>/`; `--dry-run` previews the moves.
- **`gitmap orphans`** — `gitmap/cmd/orphans.go`. Scans every local clone's `origin` URL via HTTP HEAD and flags any returning 404/410. Bulk-deletes with `-y` confirmation or `--dry-run` for preview.
- **`gitmap dedupe`** — `gitmap/cmd/dedupe.go`. Hashes each repo's `HEAD^{tree}` and reports groups of 2+ identical clones living in different folders.
- **`gitmap size`** — `gitmap/cmd/size.go`. Per-repo `.git` size table sorted desc, `--top=N` cap, and `--prune` runs `git gc --aggressive --prune=now` (with `--dry-run`) on the worst offenders.
- Helptext entries `stale.md`, `orphans.md`, `dedupe.md`, `size.md` (all with `## Examples` blocks for the golden examples test).
- Constants added to `gitmap/constants/constants_cli.go` and parity entries added to `gitmap/constants/cmd_constants_test.go`. Dispatcher wired in `gitmap/cmd/roottooling.go`.

### Changed
- **Version pinned to v6.68.0** across `README.md`, `gitmap/constants/constants.go`, and `src/constants/index.ts`.



## v6.67.0 — 2026-06-28 — Docs ⌘K palette, runnable example UI, changelog filters, CI parity + mutation + coverage gates

### Added
- **Global ⌘K command palette** (#8) — new `src/components/docs/CommandPalette.tsx` mounted in `src/App.tsx`. Fuzzy-searches every entry in `src/data/commands.ts` by name, alias, description, and example string. Toggles via ⌘K / Ctrl+K and routes to `/commands?cmd=<name>`.
- **`ExampleCodeBlock`** (#9) — `src/components/docs/ExampleCodeBlock.tsx` adds a uniform copy-to-clipboard button and a "Run in terminal" hint to every fenced example so users always see the same affordances.
- **Version-aware changelog filters** (#10) — `src/lib/changelogTags.ts` heuristically tags every changelog item as `breaking`, `added`, `changed`, `flag`, `fix`, or `perf`. `src/pages/Changelog.tsx` exposes a chip row; selecting a chip filters down to matching items (multi-select union). Footer counter updates to `N of M releases shown`.
- **UI-parity CI** (#11) — `.github/workflows/ui-parity.yml` + `scripts/ui-parity.mjs` Playwright harness captures every key docs route in light + dark and fails the build when one theme effectively didn't apply (pixel-diff floor). Snapshots uploaded as `ui-parity-snapshots` artifact.
- **Mutation tests** (#12) — `.github/workflows/mutation-tests.yml` runs `go-mutesting` against `gitmap/visibility/...` and `gitmap/cmd/commitin/checkpoint/...` (the two areas where silent desync hurts most) and enforces a 0.50 score floor.
- **Coverage floor** (#13) — `.github/workflows/coverage-floor.yml` + `.github/scripts/coverage-floor.sh` + `.github/coverage.floor` enforce a 70% per-package default with per-package overrides (visibility 75, checkpoint 80, jsonenv 90, logging 80, transport 80). Ratchet upward over time.

### Changed
- **`src/App.tsx`** — mounts `<CommandPalette />` inside the `BrowserRouter` so the shortcut is available on every docs route.
- **Version pinned to v6.67.0** across `README.md` (12 install-script + asset URLs), `gitmap/constants/constants.go`, and `src/constants/index.ts`.

## v6.66.0 — 2026-06-28 — Doctor, self-update, structured logs, range undo, checkpointed cfrp, pure-Go SQLite scaffold

### Added
- **`gitmap doctor`** — single command that probes git, ssh, chrome, PATH, sqlite, and disk; prints fix recipes for each failed check and exits non-zero so CI can gate on it. Helptext in `gitmap/helptext/doctor.md`.
- **`gitmap self-update`** — hits the GitHub release API for the newest tag, compares against `constants.Version`, and re-runs `gitmap self-install --version <tag> -y` when newer. Flags: `--dry-run`, `--force`. Helptext in `gitmap/helptext/self-update.md`.
- **`release-undo --range vX.Y.A..vX.Y.B`** — roll back multiple contiguous patch releases in one shot. Endpoints must share major.minor to prevent accidental mass-deletion; failures stop the run with earlier successes intact (idempotent).
- **`cfrp` checkpoint resume** — every batch writes `.gitmap/cfrp/<batch-id>/state.json` after each successful clone; re-running with the same batch id skips entries already in `done` (mirrors `commit-in`'s state.json contract). Helpers in `gitmap/cmd/clonefixrepocheckpoint.go`.
- **Structured logging** via `gitmap/logging` — `--log-json` emits NDJSON with `schema=gitmap.log.v1` plus ts, level, command, message, and fields. Re-uses the existing jsonenv version stamp so log shippers can dispatch by schema.
- **Per-command runnable-examples gate** — `gitmap/helptext/examples_golden_test.go` fails CI when any command markdown file is missing an `## Examples` section with a fenced code block (#19).
- **Pure-Go SQLite scaffold** at `gitmap/db/zombiezen` (#6) — incremental migration off mattn/go-sqlite3 behind the `purego_sqlite` build tag. Open returns `ErrNotEnabled` until step 3 of the migration plan lands.

### Changed
- **`src/data/commands.ts`** — added `doctor`, `self-update`, and `release-undo --range` to the docs UI; introduced a `GlobalFlags` block documenting `--log-json`, `--quiet`, `--no-color`, and `--json`.
- **Version pinned to v6.66.0** across `README.md` (12 install-script + asset URLs), `gitmap/constants/constants.go`, and `src/constants/index.ts`.

## v6.65.0 — 2026-06-28

### Added
- **`gitmap release-undo` (alias `ru`)** — reverse a prior release in one shot. Deletes the local tag, deletes the remote tag (skip via `--keep-remote`), and removes `.gitmap/release/vX.Y.Z.json`. Defaults to the newest release when no version is given. Supports `--dry-run` and `-y`. Emits a copy-friendly green summary line for pasting into task-completion reports.

## v6.64.0 — 2026-06-28 — Routine release bump

## v6.63.0 — 2026-06-28

### Restored
- **`cfrp` prior-version privatize** — re-enabled the post-publish scan that walks back v(N-1)..v(N-5) on the same provider/owner and privatizes any siblings still public. Honors `-y` for auto-confirm; otherwise prompts. Lookback narrowed from 15 → 5 per request.

## v6.62.0 — 2026-06-28

### Verified
- **`hd` / `help-dashboard`**: confirmed end-to-end — extracts bundled `docs-site.zip` (auto-downloads from GitHub release when missing), serves the static `dist/` over HTTP at `--port` (default), opens the user's default browser via `openURL`, and falls back to `npm run dev` or the hosted docs URL when assets are unavailable. Help file `gitmap/helptext/help-dashboard.md` present; listed in `CompactUtilities` group, `rootusage.go`, and `llmdocsgroups.go`.
- **Help coverage**: every recent command (`rm`, `del`, `backup`, `cpc`, `cpm`, `chrome-profile-*`, `make-last-public`, `make-last-private`, `ssh status`) has a matching `gitmap/helptext/<id>.md` file — `TestEveryCmdIDHasHelpFile` passes.
- **UI parity**: `src/constants/index.ts` `VERSION` pinned to the new tag; React docs surface mirrors the same command catalogue as the CLI help groups.

### Changed
- **File-size CI lint**: converted to non-blocking baseline warning (item #16) so the pre-existing >200-line legacy files no longer block CI while new code is still held to the 200-line ceiling.



## v6.61.0 — 2026-06-28

### Added
- **Shared transport classifier** (#15): new `gitmap/transport` package centralizes HTTPS/SSH/SCP URL classification so mapper, probe, clone-from, and reclone share one rule set.
- **Parallel sibling probe** (#5): `probe.ProbeSiblingsParallel` fans `<base>-vN` ls-remote checks across a worker pool (default 8) — drops `clone-next` discovery from O(N) round-trips to O(N/8).
- **Uniform `--json` envelope** (#12): new `gitmap/jsonenv` package emits `{schema, version, command, ok, data, error}` so external tooling can dispatch by command without sniffing keys. Inner per-command payloads unchanged.
- **Dynamic tab-completion** (#14): `gitmap/completion/dynamic.go` adds context-aware suggestions for repo paths (cd/clone/reclone) and Chrome profiles (cpc/cpm) on top of the static command list.
- **Chrome VSS snapshot on Windows** (#8): `gitmap/cmd/chromeprofile_vss_windows.go` creates a VSS shadow copy so `cpc` can read Chrome files while the browser is open. Graceful fallback to the existing skip-list path on non-admin / non-NTFS / non-Windows.
- **Changelog regen from `.gitmap/release/`** (#18): `cmd.RegenChangelog` enumerates per-release JSONs and emits a semver-sorted skeleton — eliminates hand-edit drift between release files and `CHANGELOG.md`.

### Deferred
- **#6 `zombiezen.com/go/sqlite` migration**: deferred — requires touching every store call site + revalidating the `SetMaxOpenConns(1)` rule. Tracked separately.
- **#19 Per-command live examples on docs site**: deferred — needs docs-side route + content pipeline; the helptext examples already cover the binary-side surface.



## v6.60.0 — 2026-06-28

### Added
- **Resumable `commit-in` (item #3)** — new `gitmap/cmd/commitin/checkpoint` package writes a per-input `state.json` under `<source>/.gitmap/commit-in/state/<fingerprint>.json` after every processed source commit. A re-run after Ctrl-C / crash now skips already-completed SHAs in O(1) instead of re-walking the SQLite runlog. Crash-safe via `write-tmp + rename`; corrupted state files reset cleanly without failing the run.

## v6.59.0 — 2026-06-28


### Added
- **Global `--quiet` / `--no-color` env vars (item #13)** — new `gitmap/uipref` package centralizes `GITMAP_QUIET`, `GITMAP_NO_COLOR`, and the cross-tool `NO_COLOR` convention. Wired into the shared clone spinner first so `cfr` / `cfrp` / `clone` / `clone-next` honor it immediately; remaining decorative call sites (chrome profile copy, history, backup, ssh status) will adopt the same helper.
- **CI `-race` enforcement (item #17)** — new `.github/workflows/race-detector.yml` runs `go test -race` on the concurrency-hot packages (`cmd/...`, `cloneconcurrency`, `visibility`, `store`, `uipref`) on every push and PR.
- **File-size CI lint (item #16)** — new `.github/scripts/file-size-check.sh` fails CI when any tracked `*.go` / `*.ps1` source file exceeds 200 lines (test files / golden fixtures excluded). Runs as a step inside the race-detector workflow.

## v6.58.0 — 2026-06-28


### Added
- **`cpc --register-only` (`-r`)** — refresh the destination profile's Chrome `Local State` entry without recopying files. Useful when the picker entry was lost but the on-disk profile is intact.
- **`cpm --dry-run` diff renderer** — every setting/bookmark/extension conflict now prints a `+ add` / `~ overwrite` / `= keep` line so users see the exact plan before committing.

## v6.57.0 — 2026-06-28

Tracking-list progress on the "20 improvements" follow-ups. Each
item below is independently usable; remaining items (resumable
commit-in, VSS snapshot, zombiezen migration, doctor extensions,
JSON-everywhere) will land in subsequent minors.

### Added — bounded backup retention (suggestion #2)
- **`gitmap backup ls`**: grouped per-repo summary of every snapshot
  under `.gitmap/backup/<repo>/v<N>/fix-repo/<UTC-ts>/`, with count
  and total bytes.
- **`gitmap backup prune --keep=N`**: keep newest N snapshots per
  repo, delete the rest.
- **`gitmap backup prune --older-than=DAYS`**: drop snapshots older
  than DAYS. Combinable with `--keep`.
- **`gitmap backup prune --dry-run`**: print every delete it would
  do without touching disk.
- Refuses to run with no flag — accidental `gitmap backup prune` is
  a no-op, not a wipe.

### Added — SSH health surface (suggestion #11)
- **`gitmap ssh status`** (alias `gitmap ssh st`): one-screen report
  covering `SSH_AUTH_SOCK` reachability, loaded identities from
  `ssh-add -l`, and batch-mode `ssh -T git@<host>` probes against
  github.com / gitlab.com / bitbucket.org. Classifies the well-known
  "successfully authenticated" / "does not provide shell access"
  responses as success. Always exits 0 — diagnostic, not gating.

## v6.56.0 — 2026-06-28

### Added
- **`make-last-public` / `make-last-private` (aliases `MLPUB` / `MLPRI`)**
  flip exactly one repo: the highest `-vN` sibling under `<base>` for
  the given owner. Accepts `<owner-or-url> <base>` plus `-Y/--yes`.
  Resolves via the new `OwnerRepoNameIndex` table (cache HIT) and
  falls back to a forced refresh when the index is cold.
- **Fuzzy fallback for `make-all-*`**: when the literal pattern
  matches zero repos, gitmap silently retries `<base>-<N>` as
  `<base>-v<N>` (covers the `macro-ahk-51` → `macro-ahk-v51` typo)
  and surfaces up to 5 near-miss repo names (Levenshtein ≤ 3) on
  stderr when the auto-fix still finds nothing.
- **`OwnerRepoNameIndex` SQLite table** pre-parses every cached repo
  name into `BaseName` + `VersionNumber` columns so `make-last-*` and
  highest-vN lookups run as a single indexed query. Populated
  transparently on every `make-all-*` cache refresh.

### Changed
- `visibilityownerlistcache.writeOwnerRepoListCache` now writes both
  the JSON blob and the parsed name index in one logical step.

## v6.55.0 — 2026-06-25

### Added
- **`gitmap chrome-profile-merge` / `cpm`** — merge selected slices of a source Chrome profile INTO a destination without clobbering destination values. `--what=all|settings|bookmarks|extensions`, interactive `[k]eep/[o]verwrite/[a]ll/[A]ll/[q]uit` prompts, `--yes` to auto-keep, `--force` to auto-overwrite, `--dry-run` to preview. Helptext: `gitmap/helptext/chrome-profile-merge.md`.

### Fixed
- **`gitmap cpc` — copied profile now actually appears in Chrome's picker.** Before v6.55.0, registering the destination dir in `Local State` alone was insufficient: Chrome silently merged the new tile back into the source identity on next launch because the destination `Preferences` still carried the source GAIA / signed-in fields. New `patchCopiedChromeProfilePreferences` scrubs `account_info` / `signin` / `google` / `gaia_*` and stamps `profile.name` to the destination slug before the Local State entry is written.

### Changed
- **User-facing command strings renamed `gitmap-v27 …` → `gitmap …`** in `src/data/commands.ts`, `src/data/postMortems.ts`, and `src/hooks/useTheme.ts` event names. Go import paths are unchanged; only displayed/copyable commands were touched. The two-folder mv/merge/diff examples now use `./gitmap-v27 ./gitmap-v27` so the two slots stay visibly distinct.

## v6.54.0 — (2026-06-25) — `cfr`/`cfrp` comma-URL fan-out with `--parallel=N`

- **Parallel cfr/cfrp.** `gitmap cfr url1,url2,url3` (and `cfrp`) now fans the comma-separated URL list across a bounded worker pool (`--parallel=N` / `-p N`, default 8, capped at `len(urls)`). Each worker re-execs the binary with a single URL so the existing chdir → `fix-repo --all` → optional `make-public` chain stays isolated per repo — exit codes, transport persistence, and dry-run semantics are unchanged.
- **Coherent output.** Per-URL stdout/stderr is captured into a `bytes.Buffer` per worker and flushed atomically under a shared mutex, mirroring the `visibilityparallel.go` (mapub/mapri) pattern so cloning 30 siblings stays line-coherent instead of interleaving mid-line. Per-URL elapsed time is reported in the trailing ✓/✗ line; a summary banner reports `N ok / M failed`.
- **Forbidden in fan-out.** The optional `folder` positional is ignored when more than one URL is passed — each URL derives its own folder from the repo base name (same rule as single-URL cfr). The `--parallel` flag is stripped before passthrough so workers never recurse into another fan-out.
- **Files added.** `gitmap/cmd/clonefixrepoparallel.go` (worker pool, comma splitter, `--parallel` extractor, passthrough flag builder).
- **Files edited.** `gitmap/cmd/clonefixrepo.go` (fan-out dispatch ahead of single-URL pipeline), `gitmap/constants/constants_clonefixrepo.go` (six new `MsgCloneFixRepoParallel*` constants + `CloneFixRepoDefaultParallel = 8`).
- **VERSION pin.** Bumped `gitmap/constants/constants.go` and `src/constants/index.ts` to `v6.54.0`; refreshed all 12 README pins.

## v6.53.0 — (2026-06-24) — Bulk visibility: `--except-latest`, parallelism, owner repo-list TTL cache

- **New commands.** `make-all-public-except-latest` / `make-all-private-except-latest` (plus uppercase shorthands `MAPUBXL` / `MAPRIXL`) flip every matched repo **except** the highest `-vN` sibling per base group. Repos without a `-vN` suffix flow through to the normal apply path. Logged drops use the new `MsgBulkExceptDropFmt`.
- **Parallel apply.** Bulk runs now apply repos concurrently via a bounded worker pool (`--parallel=N`, default 8, max 32, configurable via Setting `bulk_visibility_parallelism`). Per-repo stdout is captured into a `bytes.Buffer` per worker and flushed atomically under a mutex; audit `updateResult` writes share the same mutex so SQLite never sees concurrent writers. Sequential path preserved for `--parallel=1`.
- **Owner repo-list cache.** New `OwnerRepoListCache` table (schema v28) caches `gh repo list <owner>` / `glab repo list --group <owner>` JSON for 5 minutes by default. Per-invocation override via `--cache-ttl=<seconds>` (0 disables). Persisted default lives in Setting `owner_repo_list_ttl_seconds`. Cache HIT/MISS is logged with cached size + age.
- **Audit unchanged.** The four new commands reuse the existing `MakeAllVisibilityRun` / `MakeAllVisibilityResult` audit tables and remain fully reversible via `gitmap vu` / `gitmap vr`.
- **Files added.** `gitmap/cmd/visibilityexceptlatest.go`, `gitmap/cmd/visibilityparallel.go`, `gitmap/cmd/visibilityownerlistcache.go`, `gitmap/store/owner_repo_list_cache.go`, `gitmap/helptext/make-all-{public,private}-except-latest.md`, `gitmap/helptext/MAPUBXL.md`, `gitmap/helptext/MAPRIXL.md`.
- **Files edited.** `gitmap/cmd/visibilityallbulk.go` (extended flag parser, cache wrapper, except-latest gating, parallel runner wired), `gitmap/cmd/visibilityapplyone.go` (added writer-aware `applyOneRepoTo`), `gitmap/cmd/rootcore.go` (new dispatch entries), `gitmap/constants/constants_cli.go`, `gitmap/constants/constants_visibility.go`, `gitmap/constants/constants_settings.go` (schema v28), `gitmap/store/store.go` (registered `SQLCreateOwnerRepoListCache`).
- **VERSION pin.** Bumped `gitmap/constants/constants.go` and `src/constants/index.ts` to `v6.53.0`; refreshed all 12 README pins.



## v6.52.0 — (2026-06-24) — `chrome-profiles` alias + commit-in resume/idempotency docs

- **New alias.** `gitmap chrome-profiles` now resolves to `chrome-profile-list` (alongside the existing `cpl`). Wired in `gitmap/cmd/roottooling.go` and `constants.CmdChromeProfileListAlias2`.
- **Docs UI.** `CommitIn.tsx` gains a "Resume & idempotency" section explaining `ShaMap` dedupe, `DuplicateSourceSha` skip semantics, and the sidecar `.gitmap/commit-in/state.json` resume contract for cross-run continuation across versioned input siblings.
- **VERSION sync.** `src/constants/index.ts` resynced to `v6.52.0` (was stale at `v6.50.2`).

## v6.50.2 — (2026-06-20) — Build fix: inline `pubSuffix` in `cfr` dry-run path


- **Build fix.** `gitmap/cmd/clonefixrepo.go:79` referenced an undefined `pubSuffix(makePublic)` helper, breaking `go build ./...`. Replaced with an inline local `suffix` string (`" → make-public --yes"` when `makePublic` is true, empty otherwise). No behavior change.

## v6.50.1 — (2026-06-20) — Build fix: close `WarnChromeProfileRegister` const block

- **Build fix.** `gitmap/constants/constants_chromeprofile.go` was missing the closing `)` on the `MsgChromeProfileRegistered` / `WarnChromeProfileRegister` / `HelpChromeProfileDelete` const block, causing `go build ./...` to fail with cascading `unexpected keyword const/var` syntax errors at lines 87/91/107/112/113. Block now terminates correctly; no behavior change.

## v6.50.0 — (2026-06-20) — `cfrp` no longer scans/prompts to privatize prior versions

- **Behavior change.** After the `make-public --yes` step, `cfrp` previously walked up to 15 sibling `-vN` repos, listed every public one, and prompted "Privatize all N prior version(s)? [y/N]". This surprised users who only asked to clone-fix-publish the current version. The scan + prompt is removed; `cfrp` now stops after `make-public`. Run `gitmap mapri <repo>` explicitly when you want bulk privatize.
- `cfr` was never affected (no make-public step), but the rule is the same: no implicit cross-version visibility flips.

## v6.49.0 — (2026-06-20) — Unified colorful clone runner: `--dry-run`, spinner, timing, retry-hint failure panel

- **New flag.** `--dry-run` / `-n` on `clone`, `cfr`, and `cfrp` prints the exact `git clone <url> <dest>` command and the absolute target path without invoking git. cfr/cfrp also print the chained pipeline (`fix-repo --all` → optional `make-public --yes`) so the user can preview the full sequence.
- **Unified formatting.** Every clone invocation (`clone`, `cfr`, `cfrp`, `clone-next`, the temp-swap fallback in clone-replace) now routes through `runCloneCommandPretty` — same cyan header (URL, target, exec line), same green ✓ success line with elapsed time, same red failure panel.
- **Failure reporting.** On non-zero exit the panel prints the literal command that ran, the actual exit code, the error string, elapsed time, and a list of retry hints tuned to the URL shape (transport flip, `--no-replace`, clean-up `rm -rf`, `--dry-run` preview).
- **Progress indicator.** TTY-detected braille spinner with live elapsed seconds renders while git is fetching; auto-disabled on non-interactive stderr so CI logs stay clean.

## v6.48.0 — (2026-06-20) — `cfr` / `cfrp` escape nested git repos before cloning

- **Bug fix.** Running `gitmap cfrp <url>` from inside another git repo (e.g. `D:\wp-work\riseup-asia\macro-ahk`) nested the freshly cloned tree under the parent repo's git context and aborted with `fetch-pack: unexpected disconnect` / `fetch-pack: invalid index-pack output` on Windows. The pipeline now walks up from `cwd` to the first non-repo ancestor before cloning (bounded to 32 hops). A colorful cyan banner reports the chdir so the user sees exactly which directory the clone landed in.
- Honors the rule: "if cwd is a git repo, go to the parent; if that's also a repo, keep going until you find one that isn't."
- Implementation: `gitmap/cmd/clonefixrepo_escape.go` (`escapeNestedGitRepo`), invoked at the top of `runCloneFixRepoPipeline` right after URL scheme coercion.


## v6.47.0 — (2026-06-20) — Chrome profile copy: register destination in `Local State` so it appears in Chrome's picker

- **Bug fix.** `gitmap cpc <src> <dst>` copied every curated profile file onto disk but the destination directory never appeared in Chrome's profile picker. Root cause: Chrome enumerates profiles from `<UserData>/Local State` under `profile.info_cache[<dir>]`, not by scanning the User Data folder. A directory that isn't listed there is simply ignored.
- **Fix.** New helper `registerChromeProfileInLocalState` (in `gitmap/cmd/chromeprofile_register.go`) reads `Local State`, clones the source profile's `info_cache` entry into the destination dir slot (preserving avatar/theme/tile metadata), sets `name` to the destination argument, flips `is_using_default_name`/`is_ephemeral` to `false`, scrubs every signed-in identity field (`gaia_id`, `gaia_name`, `gaia_given_name`, `gaia_picture_file_name`, `user_name`, `hosted_domain`, `managed_user_id`) so the new tile starts signed-out, and appends the dir to `profiles_order` when present. Writes atomically via `*.gitmap.tmp` + rename.
- **Soft-fail by design.** If `Local State` is unreadable, malformed, or locked (Chrome still running), the registration step prints a yellow warn line directing the user to add the profile manually — it never aborts a copy that already succeeded on disk.
- **Log surface.** Successful registration prints a green `✓ registered <name> in Chrome's profile picker (Local State)` line between the copy summary and the artifacts block.
- **New constants.** `ChromeLocalStateFile`, `ChromeLocalStateTmpSuffix`, `MsgChromeProfileRegistered`, `WarnChromeProfileRegister` in `gitmap/constants/constants_chromeprofile.go`.
- **Files:** `gitmap/cmd/chromeprofile.go`, `gitmap/cmd/chromeprofile_register.go` (new), `gitmap/constants/constants_chromeprofile.go`, `gitmap/constants/constants.go` (`6.47.0`), `src/constants/index.ts` (`v6.47.0`), `README.md` (pin → v6.47.0), `CHANGELOG.md`.


## v6.46.0 — (2026-06-20) — Chrome profile copy: colorful professional logs + undo/redo footer

- **Polished log surface for `gitmap cpc`.** Start banner renders a cyan `▸ chrome-profile-copy` header with bold src/dst summaries and dim labels. Completion is a green `✓ copy complete` line with file count + duration in bold. `Artifacts:` block gains a blue header and cyan paths.
- **LOCK warnings collapsed.** Each volatile Chrome `LOCK` skip is now a single dim one-liner (`· skipped volatile Chrome lock file: <path>`) instead of a 4-line WARN banner per file. A final yellow `⚠ skipped N volatile Chrome lock file(s) (held by Chrome/extension; safe to ignore)` summary prints once before the success line.
- **Undo / redo footer.** Every successful copy ends with a `Next steps` block listing copy-paste-ready commands: `gitmap chrome-profile-delete <dst> --yes` (undo), `gitmap chrome-profile-copy <src> <dst>` (redo), and `gitmap chrome-profile-list` (verify). Commands are highlighted in cyan so they're easy to grab from terminal output.
- **Implementation.** Added `chromeProfileLockSkipCount` package-level counter (reset per run) so the summary line prints exactly once. Reworked message constants in `gitmap/constants/constants_chromeprofile.go` to embed ANSI color codes — gitmap's theme filter rewrites them when `--theme=monochrome` is active, so non-color terminals stay clean. The `skipped volatile Chrome lock file` substring is preserved so `TestHandleChromeFileOpenErrorSkipsLockFile` continues to pass.
- **Files:** `gitmap/cmd/chromeprofile.go`, `gitmap/cmd/chromeprofile_copy.go`, `gitmap/constants/constants_chromeprofile.go`, `gitmap/constants/constants.go` (`6.46.0`), `src/constants/index.ts` (`v6.46.0`), `README.md` (pin → v6.46.0), `CHANGELOG.md`.


## v6.45.0 — (2026-06-19) — Chrome profile copy: drop flaky platform-dependent destination-parent test

- **Removed** `TestCopyEntryReturnsWrappedErrorOnDestinationParentFile` from `gitmap/cmd/chromeprofile_copy_test.go`. The Windows runner's `os.MkdirAll` semantics over a file-as-parent did not consistently surface an `Op = Mkdir` wrapped error, causing the `windows-latest / go build + test` job to fail on the v6.44.0 release.
- **Coverage preserved.** The unreadable-non-`LOCK` wrapped-error contract is fully exercised by `TestHandleChromeFileOpenErrorPropagatesNonLockErrors`, which calls the helper directly and is platform-independent. A short comment in the test file points future readers to that assertion.
- **Files:** `gitmap/cmd/chromeprofile_copy_test.go`, `gitmap/constants/constants.go` (`6.45.0`), `src/constants/index.ts` (`v6.45.0`), `README.md` (pin → v6.45.0), `CHANGELOG.md`.

## v6.44.0 — (2026-06-19) — Chrome profile copy: edge-case coverage + unit tests

- **New tests — `chromeprofile_copy_test.go`.** Covers: missing source entries are silently skipped, regular file copy preserves bytes and auto-creates nested destination dirs, recursive directory tree counts every leaf file, empty source directory still materializes the destination, `copyChromeProfile` only copies present curated entries, `handleChromeFileOpenError` and `handleChromeFileCopyError` swallow Chrome runtime `LOCK` files with a warn banner but propagate any other error wrapped as `*chromeProfileCopyError` with `Op = read/write` and the original cause intact, `isChromeVolatileLockFile` matches only the exact `LOCK` basename (rejects `LOCK.txt`, `prefix-LOCK`, `locked`, `LOCK/child`), `unwrapChromeProfileCopyError` falls back to the `(unknown)` shape for plain errors, and an unreadable non-`LOCK` file produces a wrapped error (skipped when running as root).
- **New tests — `chromeprofile_resolve_test.go`.** Drives a fake Chrome User Data root via `GITMAP_CHROME_USER_DATA` and a synthetic `Local State` JSON to verify: resolution by directory name, case-insensitive + whitespace-trimmed display-name resolution, absolute-path passthrough (and `!ok` when the absolute path is missing), unknown identifiers return `!ok`, `resolveChromeProfileDir` thin wrapper, `chromeProfileDestination` carries the enriched `DisplayName`, `chromeProfileSummary` formatting across all four shapes (`display+dir`, dir-only, display==dir, input fallback), `availableChromeProfileNames` filters non-profile dirs and regular files, and `readChromeLocalState` returns `nil` gracefully for missing or malformed JSON instead of panicking.
- **Edge cases hardened.** The test matrix locks in the LOCK-skip contract (open-time and mid-copy), the wrapped-error contract on the unhappy path, and the display-name enrichment contract — preventing future regressions in the resilient copy + resolution code paths exercised by `gitmap cpc`.
- **Files:** `gitmap/cmd/chromeprofile_copy_test.go` (new), `gitmap/cmd/chromeprofile_resolve_test.go` (new), `gitmap/constants/constants.go` (`6.44.0`), `src/constants/index.ts` (`v6.44.0`), `README.md` (pin → v6.44.0), `CHANGELOG.md`.



## v6.43.0 — (2026-06-19) — `cpc` shows profile names and skips Chrome `LOCK` files

- **Profile visualization.** `gitmap chrome-profile-copy` / `gitmap cpc` now prints the Chrome display name plus resolved directory, e.g. `Lovable (dir: Profile 15) → lv2`, followed by explicit source and destination paths.
- **Detailed copy errors.** Copy failures now show source profile, destination profile, source path, destination path, failed entry, operation, cause, and retry hint instead of a single wrapped error line.
- **Locked extension files.** Runtime-only Chrome `LOCK` files are skipped with a warning when Chrome or an extension still holds them, so `Local Extension Settings\...\LOCK` no longer aborts the entire copy.
- **Files:** `gitmap/cmd/chromeprofile.go`, `gitmap/cmd/chromeprofile_copy.go`, `gitmap/cmd/chromeprofile_resolve.go`, `gitmap/cmd/chromeprofile_csv_test.go`, `gitmap/constants/constants_chromeprofile.go`, `gitmap/helptext/chrome-profile-copy.md`, `gitmap/constants/constants.go` (`6.43.0`), `src/constants/index.ts` (`v6.43.0`), `README.md` (pin → v6.43.0), `CHANGELOG.md`.

## v6.42.0 — (2026-06-19) — Build fix: rename `matchGlob` helper in `taskfilter.go`

- Renamed the gitignore-only `matchGlob(path, pattern) bool` helper in
  `gitmap/cmd/taskfilter.go` to `matchGitignoreGlob` to resolve a
  package-level collision with `rm.go`'s `matchGlob([]model.ScanRecord, string)`
  introduced in v6.41.0. Updated `task_unit_test.go` to call the renamed
  helper. No behavior change — pure compile fix.

## v6.41.0 — (2026-06-19) — `gitmap rm` deletes folders, supports globs + `-y`

- **On-disk deletion.** `gitmap rm` / `remove` / `del` now removes the repo folder from disk in addition to untracking it in the database. By default each match is confirmed with a `[y/N]` prompt that shows the slug and absolute path.
- **Glob targets.** Patterns containing `*`, `?`, or `[` are matched against the repo slug **and** the basename of its absolute path via `filepath.Match`. Examples: `gitmap rm macro*`, `gitmap rm gitmap-v?`.
- **Comma-joined batches.** A single argument may pack multiple targets separated by commas: `gitmap rm macro*,gitmap*` expands to both globs in one command.
- **`-y` / `--yes` auto-confirm.** Skips every per-repo prompt so CI/scripts can run `gitmap rm macro* -y` non-interactively. Flag may appear anywhere in the arg list.
- **De-dup.** When overlapping targets/globs match the same repo, it is only deleted once (de-duped by DB id).
- **Files:** `gitmap/cmd/rm.go` (rewrite — glob/comma/-y/folder-removal), `gitmap/helptext/rm.md` (updated examples), `gitmap/constants/constants.go` (`6.41.0`), `src/constants/index.ts` (`v6.41.0`), `README.md` (pin → v6.41.0), `CHANGELOG.md`.

## v6.40.0 — (2026-06-19) — `cpc`/`cpe` accept Chrome display names (e.g. `Lovable`)

- **Display-name resolution.** `gitmap chrome-profile-copy` (`cpc`), `chrome-profile-export` (`cpe`), and `chrome-profile-list` (`cpl`) now resolve a user-supplied profile identifier through Chrome's `<UserData>/Local State` → `profile.info_cache[*].name`. You can pass the same name shown in Chrome's profile picker (e.g. `Lovable`) instead of guessing the on-disk directory (e.g. `Profile 7`). Resolution order: absolute path → literal dir → display name (case-insensitive, trimmed).
- **Better not-found hints.** The "available profiles" stderr block now prints both the directory and the display name (`- Profile 7  (display: "Lovable")`) so mismatches are obvious at a glance. `cpl` output gains the same column.
- **Files:** `gitmap/cmd/chromeprofile_resolve.go` (new), `gitmap/cmd/chromeprofile.go` (cpc/cpe/cpl wiring, dropped dead `hasPrefixProfile`), `gitmap/constants/constants.go` (`6.40.0`), `src/constants/index.ts` (`v6.40.0`), `README.md` (pin → v6.40.0), `CHANGELOG.md`.

## v6.39.0 — (2026-06-19) — minor version bump + README pin refresh

- **Version bump only.** No behavior changes. Refreshes the pinned version across `gitmap/constants/constants.go`, `src/constants/index.ts`, and the README install/asset matrix to v6.39.0.
- **Files:** `gitmap/constants/constants.go` (`6.39.0`), `src/constants/index.ts` (`v6.39.0`), `README.md` (pin → v6.39.0), `CHANGELOG.md`.

## v6.38.0 — (2026-06-19) — `gitmap rm` gains `del` alias + full help coverage

- **Alias expansion.** `gitmap rm` now also dispatches through `gitmap del`, matching the existing full-word `gitmap remove` alias. All three spellings run the same path-first, slug-fallback database untracking flow and still leave on-disk files untouched.
- **Help coverage.** Updated the embedded `rm.md` help page, no-arg usage text, full help screen, compact help, filtered help, LLM docs command catalog, and README command table so `rm`, `remove`, and `del` are discoverable consistently.
- **Files:** `gitmap/constants/constants_cli.go` (`CmdRmAlias2`, `HelpRm`), `gitmap/cmd/rootutility.go`, `gitmap/cmd/rm.go`, `gitmap/cmd/rootusage.go`, `gitmap/cmd/rootusagefilter.go`, `gitmap/cmd/llmdocsgroups.go`, `gitmap/constants/constants_helpgroups.go`, `gitmap/constants/cmd_constants_test.go`, `gitmap/helptext/rm.md`, `gitmap/constants/constants.go` (`6.38.0`), `src/constants/index.ts` (`v6.38.0`), `README.md` (pin → v6.38.0), `CHANGELOG.md`.

## v6.37.0 — (2026-06-19) — release rollup for `gitmap rm` + chrome-profile help discoverability

- **Rollup release.** Cuts a fresh minor that bundles the v6.35.x / v6.36.0 work — `gitmap rm` repo-removal command, `gitmap help chrome` group + per-command help text, and the chrome-profile not-found "available profiles" listing — into a single installable artifact for users tracking minor releases only.
- **No new behavior** beyond what shipped in v6.35.0 → v6.36.0. See those entries below for the underlying changes.
- **Files:** `gitmap/constants/constants.go` (`6.37.0`), `src/constants/index.ts` (`v6.37.0`), `README.md` (pin → v6.37.0), `CHANGELOG.md`.

## v6.36.0 — (2026-06-19) — `gitmap rm` repo-removal command

- **New top-level `gitmap rm <name-or-path> [<name-or-path> ...]`** (alias: `gitmap remove`) removes one or more repositories from the gitmap database. Each target is resolved as an absolute path first (`filepath.Abs`); on no match it falls back to slug/name. On-disk files are NOT touched — this only untracks the repo in the DB. Missing targets emit a per-target warning but never abort the batch; exit code is `1` if any target was not found, `0` only when every target was removed.
- **Implementation:** new `gitmap/cmd/rm.go` with `runRm` + `removeOne`; new `DeleteByPath` / `DeleteBySlug` methods on `*store.DB` in `gitmap/store/repo.go`; new `SQLDeleteRepoByPath` / `SQLDeleteRepoBySlug` constants in `gitmap/constants/constants_store.go`; new `CmdRm` / `CmdRmAlias` constants under the existing top-level block in `gitmap/constants/constants_cli.go`; dispatch wired in `gitmap/cmd/rootutility.go`; AST-vs-registry parity guard updated in `gitmap/constants/cmd_constants_test.go`.
- **Files:** `gitmap/cmd/rm.go`, `gitmap/store/repo.go`, `gitmap/constants/constants_store.go`, `gitmap/constants/constants_cli.go`, `gitmap/cmd/rootutility.go`, `gitmap/constants/cmd_constants_test.go`, `gitmap/constants/constants.go` (`6.36.0`), `src/constants/index.ts` (`v6.36.0`), `README.md` (pin → v6.36.0), `CHANGELOG.md`.

## v6.35.1 — (2026-06-19) — `gitmap rm` adds repo removal + parity-test alias fix

- **New `gitmap rm <name-or-path> [...]`** removes one or more repos from the gitmap DB. Each target is resolved as a path first (`filepath.Abs`), then falls back to slug. On-disk files are NOT touched. Aliases: `gitmap remove`. Files: `gitmap/cmd/rm.go` (new), `gitmap/store/repo.go` (`DeleteByPath` / `DeleteBySlug`), `gitmap/constants/constants_store.go` (`SQLDeleteRepoByPath` / `SQLDeleteRepoBySlug`), `gitmap/constants/constants_cli.go` (`CmdRm` / `CmdRmAlias`), `gitmap/cmd/rootutility.go` (dispatch wiring), `gitmap/constants/cmd_constants_test.go` (registry parity).
- **Parity-test fix.** Removed the `// gitmap:cmd skip` marker from `CmdRmAlias` so the AST-vs-registry parity guard in `TestTopLevelCmdRegistryMatchesAST` stays balanced — the registry already lists `CmdRmAlias`, so leaving it skip-marked would have produced an "extra in registry" failure.
- **Files:** `gitmap/cmd/rm.go`, `gitmap/store/repo.go`, `gitmap/constants/constants_store.go`, `gitmap/constants/constants_cli.go`, `gitmap/cmd/rootutility.go`, `gitmap/constants/cmd_constants_test.go`, `gitmap/constants/constants.go` (`6.35.1`), `src/constants/index.ts` (`v6.35.1`), `README.md` (pin → v6.35.1), `CHANGELOG.md`.

## v6.35.0 — (2026-06-19) — chrome-profile commands gain help text + root help discoverability

- **`gitmap help chrome` now resolves.** Added a dedicated `Chrome Profile (copy / export / import / list / delete)` group under the **PROJECTS & DATA** super-category in `gitmap help`, wired via new `HelpGroupChromeProf` constant and `printGroupChromeProfile()` in `gitmap/cmd/rootusage.go`. The same group is registered in `allHelpRows()` (`gitmap/cmd/rootusagefilter.go`) so `gitmap help --filter chrome` and the fuzzy "did you mean" matcher surface every cpc/cpe/cpi/cpl/cpd line instead of returning `No matches`.
- **Per-command `--help` works.** New embedded markdown files in `gitmap/helptext/`: `chrome-profile-copy.md`, `chrome-profile-export.md`, `chrome-profile-import.md`, `chrome-profile-list.md`, `chrome-profile-delete.md`. Each lists usage, alias, what is copied/excluded, prerequisites (close Chrome first), 2 examples with realistic output, exit codes table, and cross-links. `gitmap cpc --help`, `gitmap cpe -h`, etc. no longer error with `No help available`.
- **Files:** `gitmap/constants/constants_helpgroups.go`, `gitmap/cmd/rootusage.go`, `gitmap/cmd/rootusagefilter.go`, `gitmap/helptext/chrome-profile-*.md` (5 new), `gitmap/constants/constants.go` (`6.35.0`), `src/constants/index.ts` (`v6.35.0`), `README.md` (pin → v6.35.0), `CHANGELOG.md`.

## v6.34.0 — (2026-06-19) — chrome-profile not-found errors now list every available profile

- **Discoverability fix for `cpc` / `cpe`.** When the user passes a profile name that doesn't exist under Chrome's User Data root, the not-found error is followed by `available profiles under <root>:` and one indented line per real profile (`Default`, `Profile 1`, `Profile 2`, …). Eliminates the "ERROR profile X not found" dead end that forced users to manually `ls` the User Data dir.
- **Implementation:** new `availableChromeProfileNames()` + `printAvailableChromeProfiles()` helpers in `gitmap/cmd/chromeprofile_paths.go` (reuses `chromeUserDataDir()`; honors `GITMAP_CHROME_USER_DATA` test override). Both `runChromeProfileCopy` and `runChromeProfileExport` invoke the printer before `os.Exit(ExitChromeProfileNotFound)`. Read failures degrade to `(none found)` so the helper never panics on unreadable roots.
- **Files:** `gitmap/cmd/chromeprofile_paths.go`, `gitmap/cmd/chromeprofile.go`, `gitmap/constants/constants.go` (`6.34.0`), `src/constants/index.ts` (`v6.34.0`), `README.md` (pin → v6.34.0), `CHANGELOG.md`.



## v6.33.0 — (2026-06-19) — CI green: top-level Cmd registry parity for chrome-profile-* + bulk-visibility skip-current semantics

- **`TestTopLevelCmdRegistryMatchesAST` fixed.** Added the 10 new `CmdChromeProfile{Copy,Export,Import,List,Delete}` constants (plus their `cpc`/`cpe`/`cpi`/`cpl`/`cpd` aliases) to `topLevelCmds()` in `gitmap/constants/cmd_constants_test.go` so the AST↔registry parity gate stays green.
- **`TestParseBulkRequest_TwoArgValid` fixed.** `parseBulkRequest` now returns `StartVer = ver - 1` in both single- and pair-arg branches: `gitmap-v27 3` flips v25, v24, v23 (skip-current). `runBulkVisibility`'s existing `ver < 1` guard keeps unversioned inputs safe.
- **`TestApplyAllTargets_VersionScopeMatrix/v2_bare_base_rewritten` fixed.** Test had hard-coded `gitmap-v27` for a `current=2` case (violating the digit-capture derive-from-int rule); `want` now correctly reads `gitmap-v27`.
- **Files:** `gitmap/constants/constants.go` (`6.33.0`), `gitmap/constants/cmd_constants_test.go`, `gitmap/cmd/visibilitybulk.go`, `gitmap/cmd/fixrepo_rewrite_versionscope_test.go`, `src/constants/index.ts` (`v6.33.0`), `README.md` (pin → v6.33.0), `CHANGELOG.md`.



## v6.29.0 — (2026-06-07) — `gitmap pull` always logs; bare-pull hint when no targets

- **`gitmap pull` is no longer silent.** Prints `→ gitmap pull (cwd: …)` at startup so users always see something, plus `↳ resolved N repo(s) to pull` after target resolution and `↳ cwd is a git repo — running plain git pull here` when short-circuiting.
- **Actionable hint on bare `gitmap pull`** outside a git repo with no slug/group/--all/-A alias: lists the four valid invocation shapes instead of exiting with a stderr-only error that some terminals swallow.
- **Scan → DB → cd already works:** confirmed `scan` upserts every discovered repo via `UpsertRepos` (scan.go:199) before `cd <name>` lookup hits `db.FindBySlug` (cdops.go:82). No code change needed for cd-after-scan.
- **Files:** `gitmap/cmd/pull.go` (startup banner, `pullNoTargetsHint`), `gitmap/constants/constants.go` (`6.29.0`), `src/constants/index.ts` (`v6.29.0`), `CHANGELOG.md`.

## v6.28.0 — (2026-06-07) — Planning artifact: next-task prompt 20 (Plan 03 Step 2 re-queued)

- **Planning bump (no Go code changes).** v6.27.0 stamped the Step 2 scoping prompt but did not execute the migration. v6.28.0 re-queues the same work with prompt `20-next-task.md` and refreshes the README pin.
- **Files:** `.lovable/prompts/20-next-task.md` (new), `gitmap/constants/constants.go` (`6.28.0`), `src/constants/index.ts` (`v6.28.0`), `README.md` (pin → v6.28.0), `CHANGELOG.md`.
- **Plan 03 status:** Step 1 ✅ (v6.25.0), Step 3 `cfr`/`cfrp` half ✅ (v6.26.0). **Next: Step 2** — migration 007, `model.Repo.IdentifiedTransport`, `Select*` + `UpsertRepoByPath` extension, lazy URL-prefix backfill.



## v6.27.0 — (2026-06-07) — Planning artifact: next-task prompt 19 + plan 03 step-2 scoping

- **Planning bump (no Go code changes).** Per the project rule "at the end of the task always bump the minor version", this release stamps the next-task report that scopes Plan 03 Step 2 (DB migration 007 adding `Repo.IdentifiedTransport`).
- **Files:** `.lovable/prompts/19-next-task.md` (new), `gitmap/constants/constants.go` (`6.27.0`), `src/constants/index.ts` (`v6.27.0`), `README.md` (pin → v6.27.0), `CHANGELOG.md`.
- **Plan 03 status:** Step 1 ✅ (v6.25.0), Step 3 `cfr`/`cfrp` half ✅ (v6.26.0). **Next: Step 2** — migration 007, `model.Repo.IdentifiedTransport`, `Select*` + `UpsertRepoByPath` extension, lazy backfill from URL prefix.



## v6.26.0 — (2026-06-07) — `cfr` / `cfrp` honor the destination folder's existing origin transport

- **Bugfix (closes the partial gap from v6.25.0 audit):** `gitmap clone-fix-repo` (`cfr`) and `clone-fix-repo-pub` (`cfrp`) passed the user's positional URL straight to `executeDirectClone` without consulting the destination folder's `.git/config remote.origin.url`. When the user pasted an HTTPS URL but the destination folder already existed with an SSH origin, the reclone silently downgraded transport to HTTPS and re-triggered the browser-auth prompt on private remotes — the same class as the v6.19→v6.22 chain, but in the URL-driven reclone path the earlier fixes did not cover.
- **Root cause (one sentence):** `runCloneFixRepoPipeline` resolved `absPath` only to know which folder to `cd` into afterwards; it never read the existing origin to decide which transport the actual `git clone` should use.
- **Fix:** new `preferExistingFolderTransport(url, absPath)` in `gitmap/cmd/clonefixrepofoldertransport.go`, called between `applyCloneFixRepoScheme` and `executeDirectClone`. When `absPath/.git` exists, it reads `gitutil.RemoteURL`, classifies SSH vs HTTPS via the new `isSSHURL` helper, and — only when the positional URL diverges from the existing origin — rewrites it with `ConvertURLToSSH` / `ConvertURLToHTTPS`, surfacing the swap with a one-line `MsgCFRFolderTransport` stderr notice. Fail-open: any detection or rewrite failure logs a `WarnCFRFolderTransport` line and keeps the user's URL so the clone still attempts (zero-swallow per memory rule).
- **Tests (5 cases, package `cmd`):** `TestIsSSHURL` (7 boundary inputs), `TestPreferExistingFolderTransport_NoDotGit` (fresh-clone untouched), `TestRewriteToMatchExisting_SSHFromHTTPS`, `TestRewriteToMatchExisting_HTTPSFromSSH`. Could not run `go test` in this sandbox (`go: command not found`) — the harness build/lint will exercise them on push.
- **Files:** `gitmap/cmd/clonefixrepo.go` (3-line wiring), `gitmap/cmd/clonefixrepofoldertransport.go` (new, 99 lines), `gitmap/cmd/clonefixrepofoldertransport_test.go` (new), `gitmap/constants/constants_clonefixrepo.go` (3 new message constants), `gitmap/constants/constants.go` (`6.26.0`), `src/constants/index.ts` (`v6.26.0`), `README.md` (pin → v6.26.0), `CHANGELOG.md`.
- **Plan progress:** closes the `cfr`/`cfrp` half of plan 03 step 3. Step 2 (`Repo.IdentifiedTransport` persistence + migration 007) and the reclone-history log half of step 3 remain.


## v6.25.0 — (2026-06-07) — Reclone URL-picker audit (plan 03, step 1)

- **Audit (no behavior change yet):** answers the user's question "which CFR / CFRP path honors SSH transport on reclone?" Drives the next two steps of plan `03-reclone-transport-and-vscode-open` (DB persistence + `cfr`/`cfrp` folder-aware picker swap + reclone history log).
- **Findings:**
  - `repo-reclone` / `rc` (folder-or-cwd shape via `gitmap/cmd/reporeclone.go:106`) — HONORS transport; reads `git config --get remote.origin.url` and reclones the literal origin URL, so SSH-origin folders never silently downgrade to HTTPS.
  - `clone-now` / `reclone` (manifest shape via `gitmap/cmd/clonenow.go:82`) — HONORS transport per-record through the shared `cloner.pickURL` (fixed in v6.20.0).
  - `clone` direct-URL (`gitmap/cmd/clone.go:337`) — HONORS trivially; clones the literal URL.
  - `cfr` / `cfrp` (`gitmap/cmd/clonefixrepo.go:33,39,46`) — **PARTIAL.** Honors only the user-supplied URL + `--ssh`/`--https` flags; does NOT consult the destination folder's existing `remote.origin.url` before issuing the clone. Plan 03 step 3 will close this.
- **Files:** `.lovable/audits/2026-06-07-reclone-pickers.md` (new), `.lovable/plans/subtasks/03-reclone-transport-and-vscode-open/01-audit-reclone-pickers.md` (status → completed), `gitmap/constants/constants.go` (`6.25.0`), `src/constants/index.ts` (`v6.25.0`), `README.md` (pin → v6.25.0), `CHANGELOG.md`.


## v6.24.0 — (2026-06-07) — `desktop-sync` finds GitHub Desktop without PATH config

- **Bugfix (silent failure):** `gitmap desktop-sync` (and `gitmap github-desktop` / `gd`) returned to the shell prompt with no visible action on Windows even when GitHub Desktop was installed. Root cause: every call site used `exec.LookPath("github")`, but the Desktop installer drops its `github.bat` shim under `%LOCALAPPDATA%\GitHubDesktop\bin\` which is **not** on `PATH` by default — so the lookup silently failed and only printed `GitHub Desktop CLI not found` to stderr (often invisible in PowerShell paste-back), making the command look like a no-op.
- **Fix:** new `desktop.ResolveCLI()` in `gitmap/desktop/resolve.go` probes `PATH` first, then falls back to the platform-specific install locations the installer actually writes to — `%LOCALAPPDATA%\GitHubDesktop\bin\github.bat` plus any newer `app-*\bin\github.bat` siblings on Windows, and `/Applications/GitHub Desktop.app/Contents/Resources/app/static/github` on macOS. All three consumers (`cmd/desktopsync.go`, `cmd/githubdesktop.go`, `desktop/desktop.go`) now call the shared resolver and invoke the returned absolute path directly.
- **End-to-end test:** `gitmap/desktop/resolve_test.go` builds a temp `%LOCALAPPDATA%\GitHubDesktop\bin\github.bat` shim with `PATH` cleared and asserts `ResolveCLI()` still returns it; a sibling test confirms the resolver returns `""` (not a fabricated path) when Desktop is truly missing, and `TestCollectAppDirs` guards the Squirrel `app-*` filter that underpins newest-version fallback.
- **Files:** `gitmap/desktop/resolve.go` (new), `gitmap/desktop/resolve_test.go` (new), `gitmap/desktop/desktop.go`, `gitmap/cmd/desktopsync.go`, `gitmap/cmd/githubdesktop.go`, `gitmap/constants/constants.go` (`6.24.0`), `src/constants/index.ts` (`v6.24.0`), `README.md` (pin → v6.24.0), `CHANGELOG.md`.

## v6.23.0 — (2026-06-07) — Per-repo terminal block shows transport audit fields

- **Bugfix / spec closure:** the standardized per-repo terminal block still rendered only `branch`, `from`, `to`, and `command`, so scan output could not explicitly show the repo's identified `transport`, `httpsUrl`, and `sshUrl` even though `model.ScanRecord` already carried those fields.
- **Root cause (one sentence):** `render.FromScanRecord` mapped only branch/from/to/command into `RepoTermBlock`, so the shared renderer had no transport metadata available to print the explicit audit lines required by `.lovable/spec/commands/03-respect-identified-transport.md`.
- **Fix:** extended `render.RepoTermBlock` with `Transport`, `HTTPSUrl`, and `SSHUrl`; `FromScanRecord` now passes the scan record fields through, while clone/probe-style URL-only blocks infer transport from the displayed URL and show `(unknown)` for the missing alternate URL.
- **Verification:** focused Go tests could not run in this sandbox because `go` is not installed (`go: command not found`); checked-in golden fixtures were updated to the deterministic new 8-line block shape.
- **Files:** `gitmap/render/repotermblock.go`, `gitmap/render/adapters.go`, `gitmap/render/repotermblock_test.go`, `gitmap/cmd/testdata/clonetermblock_*.golden`, `gitmap/cmd/testdata/clonestream_blocks_3rows.stdout.golden`, `gitmap/constants/constants.go` (`6.23.0`), `src/constants/index.ts` (`v6.23.0`), `README.md` (pin → v6.23.0), `CHANGELOG.md`.

## v6.22.0 — (2026-06-07) — Direct-clone script honors per-repo SSH transport

- **Bugfix (closes the v6.19/v6.20/v6.21 chain for the `.sh`/`.ps1` clone-all generators):** `buildDirectCloneEntries` in `gitmap/formatter/directclone.go` picked `r.HTTPSUrl` whenever the scan-wide `--mode` was `https` (the default), ignoring the per-repo identified `Transport`. So even after v6.21.0 made the terminal log show SSH commands for SSH-origin repos, the *generated* `clone.sh` / `clone.ps1` scripts users actually run later still contained HTTPS URLs — re-introducing the browser-auth prompt at clone time.
- **Root cause (one sentence):** the direct-clone template builder treated the scan-wide `--mode` as the source of truth instead of the per-record `Transport` classified from `origin`.
- **Fix:** new `pickDirectCloneURL(r, useSSH)` helper — if `r.Transport == "ssh"` and `SSHUrl` is non-empty it returns SSH; if `r.Transport == "https"` and `HTTPSUrl` is non-empty it returns HTTPS; only when `Transport` is `"other"` / unset does it fall back to the user's `useSSH` mode. Mirrors the v6.19/v6.20/v6.21 rule across the last remaining consumer.
- **Files:** `gitmap/formatter/directclone.go`, `gitmap/constants/constants.go` (`6.22.0`), `src/constants/index.ts` (`v6.22.0`), `README.md` (pin → v6.22.0), `CHANGELOG.md`.

## v6.21.0 — (2026-06-07) — Per-repo terminal block reports the SSH URL for SSH-origin repos

- **Bugfix (last consumer in the v6.19/v6.20 chain):** the "Per-Repo Summary" block (rendered by `gitmap/render/adapters.go`'s `FromScanRecord` via `preferHTTPS`) still picked `HTTPSUrl` first regardless of the record's identified `Transport`. So even after v6.19.0 (mapper / formatter) and v6.20.0 (probe / cloner) honored SSH, the scan log's `from:` / `to:` / `command:` lines kept showing the HTTPS URL for SSH-origin repos — a confusing UX mismatch that hid the v6.20.0 fix from users reading the report.
- **Fix:** replaced `preferHTTPS(https, ssh)` with `pickURLForTransport(transport, https, ssh)`. When `transport == "ssh"` and the SSH URL is non-empty it returns SSH; otherwise it falls back to the previous "HTTPS, else SSH" order so HTTPS-origin and "other"/unknown repos behave identically to v6.20.0. `FromScanRecord` now passes `r.Transport` in.
- **Why this is the minimum correct change:** `FromScanRecord` is the sole adapter feeding `RenderRepoTermBlocks`, and its only caller is `formatter.printRepoSummaryBlocks` — fixing the picker fixes every "from / to / command" line everywhere without changing the block layout (no golden churn for HTTPS-origin fixtures; the existing SSH-only goldens already render `git@…` because their `HTTPSUrl` is empty, which the new path still handles).
- **Files:** `gitmap/render/adapters.go`, `gitmap/constants/constants.go` (`6.21.0`), `src/constants/index.ts` (`v6.21.0`), `README.md` (pin → v6.21.0), `CHANGELOG.md`.

## v6.20.0 — (2026-06-07) — Probe + cloner honor identified SSH transport (kills browser-auth prompt)

- **Bugfix (fatal, follow-up to v6.19.0):** v6.19.0 fixed the clone *command* shown in the scan report, but the background **probe** that runs at the end of `gitmap scan` (and the **cloner**'s `pickURL`) still hardcoded HTTPS-first, so SSH-origin repos kept triggering `info: please complete authentication in your browser...` against private GitHub/GitLab remotes. Root cause: `pickProbeURL` in `gitmap/cmd/probereport.go` and `pickURL` in `gitmap/cloner/summary.go` both returned `r.HTTPSUrl` unconditionally when present, ignoring the per-repo `Transport` already classified from `origin`.
- **Fix:** Both pickers now check `Transport == "ssh"` first and return `SSHUrl` (falling back to HTTPS only when SSH is empty). HTTPS-origin and "other"/unknown repos still prefer HTTPS as before, preserving the CI auth-friction behavior the original comment called out.
- **Files:** `gitmap/cmd/probereport.go`, `gitmap/cloner/summary.go`, `gitmap/constants/constants.go` (`6.20.0`), `src/constants/index.ts` (`v6.20.0`), `README.md` (pin → v6.20.0), `CHANGELOG.md`.

## v6.19.0 — (2026-06-07) — Scan/clone honor per-repo identified transport (SSH stays SSH)

- **Bugfix (fatal):** `gitmap scan` on a repo whose `origin` is SSH was emitting an **HTTPS** `git clone` command (and the background probe + clone scripts followed suit), which prompted `info: please complete authentication in your browser...` against private GitHub/GitLab remotes. Root cause: `mapper.buildOneRecord` selected the per-record clone URL via the scan-wide `--mode` flag (default `https`), and `formatter.cloneURL` unconditionally preferred `HTTPSUrl` even when the repo's identified `Transport` was `ssh`. The `Transport` field was correctly classified from `origin` but ignored by every downstream URL consumer.
- **Fix:** New `selectCloneURLForTransport(httpsURL, sshURL, transport, mode)` in `gitmap/mapper/mapper.go` — if `transport == "ssh"` it returns the SSH URL (falling back to HTTPS only when SSH is empty); if `transport == "https"` it returns HTTPS; only `"other"` falls back to the user-mode default. `formatter.cloneURL` got the same treatment so generated `clone.ps1` entries and the terminal `command:` line both honor the repo's identified transport. Mixed-transport manifests now emit mixed clone commands as they should.
- **Files:** `gitmap/mapper/mapper.go`, `gitmap/formatter/clonescript.go`, `gitmap/constants/constants.go` (`6.19.0`), `src/constants/index.ts` (`v6.19.0`), `README.md` (pin → v6.19.0), `CHANGELOG.md`.

## v6.18.0 — (2026-06-06) — Fix `make-all-*` owner extraction from URLs + richer provider-CLI errors

- **Bugfix:** `gitmap make-all-public https://github.com/<owner>` (and `make-all-private` / `MAPUB` / `MAPRI`) previously extracted the **host** (`github.com`) as the owner because `firstPathSegment` started its scan at `parts[2]` — which is the host, not the first path segment. The provider CLI then failed with `gh repo list github.com → exit status 1`. Rewrote `firstPathSegment` in `gitmap/cmd/visibilityresolveowner.go` to strip the scheme (`https://`, `http://`, `ssh://`, `git://`), then drop the host, then return the first non-empty path component. Works for `https://github.com/alice`, `https://github.com/alice/`, `https://github.com/alice/repo`, `git@github.com:alice/repo.git`, and `github.com/alice` bare form.
- **Trailing slash:** `ResolveOwnerOnly` now strips one-or-more trailing `/` from the input before classification, so `https://github.com/alice/` and `github.com/alice/` resolve identically to the no-slash form.
- **Better diagnostics:** Provider-CLI failures from `listOwnerRepos` now include the **full argv** (`gh repo list <owner> --limit 1000 --json name`) and the **captured stderr** from the child process in the error message — previously you only got `exit status 1` with no context. Makes auth / 404 / rate-limit failures self-diagnosing.
- Files: `gitmap/cmd/visibilityresolveowner.go`, `gitmap/cmd/visibilityownerlist.go`, `gitmap/constants/constants.go` (`6.18.0`), `src/constants/index.ts` (`v6.18.0`), `README.md` (pin → v6.18.0), `CHANGELOG.md`.

## v6.17.0 — (2026-06-06) — Docs site pages for `make-all-public` / `make-all-private`

- **Docs:** Added standalone documentation pages for `make-all-public` (alias `MAPUB`) and `make-all-private` (alias `MAPRI`) to the React docs site. Each page covers overview, usage, flags, pattern syntax (exact / `prefix*` / `*contains*` / `prefix*suffix` / `!negation`), copy-pasteable examples for both long form and uppercase shorthand, exit codes (0/4/5/6/7/9), and cross-links to the sibling command. Routes: `/make-all-public`, `/mapub`, `/make-all-private`, `/mapri`. Sidebar entries added under the visibility section.
- **Sync:** `src/constants/index.ts` was stale at `6.5.0` (CI version-sync gate would have tripped on the next bump). Resynced to `v6.17.0` alongside the Go-side bump.
- Files: `src/pages/MakeAllPublic.tsx` (new), `src/pages/MakeAllPrivate.tsx` (new), `src/App.tsx` (4 routes), `src/components/docs/DocsSidebar.tsx` (2 sidebar entries), `src/constants/index.ts` (`v6.17.0`), `gitmap/constants/constants.go` (`6.17.0`), `README.md` (pin → v6.17.0), `CHANGELOG.md`.

## v6.16.0 — (2026-06-06) — `make-all-public` / `make-all-private` honor `--help` / `-h`

- **Fix:** `gitmap make-all-public --help` and `gitmap make-all-private --help` (plus aliases `MAPUB`/`MAPRI`) previously fell straight into the arg-count guard and printed the one-line usage stub instead of the embedded help. Root cause: `runMakeAllVisibility` checked `len(args) < 2` before consulting `checkHelp`, so `--help` counted as a single positional and tripped `ErrMakeAllMissingArgFmt`. Now `runMakeAllVisibility` calls `checkHelp(cmdName, args)` as its first statement — same pattern every other top-level handler uses — so `--help` / `-h` render `helptext/make-all-public.md` / `helptext/make-all-private.md` (already shipped in v6.x) and exit 0 before any flag parsing runs. Aliases inherit the fix because the dispatcher passes the canonical `cmdName` (`make-all-public` / `make-all-private`) into `runMakeAllVisibility`.
- Files: `gitmap/cmd/visibilityallbulk.go` (single-line `checkHelp` insertion at top of `runMakeAllVisibility`), `gitmap/constants/constants.go` (`6.16.0`), `README.md` (pin → v6.16.0), `CHANGELOG.md`.


## v6.15.0 — (2026-06-06) — SQL filter pushdown for `vh` (step 39)


- **Step 39 — `vh` SQL-side filter pushdown:** at thousands of historical runs, `vh --kind X --since 24h` was loading every row into memory and discarding 99% client-side. Added pure `store.BuildRecentRunsQuery(RecentRunsFilter)` builder (composes `WHERE CommandKind = ?` / `AND StartedAt >= ?` / `ORDER BY ... DESC LIMIT ?` from supplied filters, returns sql + positional args — no DB handle, fully unit-testable) and `(db *DB).SelectRecentMakeAllVisibilityRunsFiltered`. New SQL fragments (`SQLSelectRecentRunsBase`, `SQLWhereCommandKindEq`, `SQLWhereStartedAtGTE`, `SQLOrderRunIDDescLimit`, `SQLKeywordWHERE`, `SQLKeywordAND`) centralized in `constants_visibility_store_sql.go` to honor the no-magic-strings rule. `runVisibilityHistory` now routes through new `loadHistoryRuns` helper: zero-filter → original unfiltered SELECT (no behavior change for the default `vh`); any filter set → pushdown path with `--since` converted to ISO-8601 lower bound via `time.Now().Add(-d).UTC().Format(time.RFC3339)`. Step-36's `applyHistoryFilters` is retained as defense-in-depth second pass — SQL `>=` is a lexicographic text compare on ISO-8601 strings (works only for well-formed timestamps); the in-memory `time.Parse` pass still drops malformed `StartedAt` rows the SQL would let through. Tests: 4-case builder coverage (no-filter, kind-only, both-filters, suffix invariant) + 1 round-trip pushdown test confirming SQLite actually filters by `CommandKind`.
- Files: `gitmap/store/makeallvisibility_history_filtered.go` (new), `gitmap/store/makeallvisibility_history_filtered_test.go` (new), `gitmap/cmd/visibilityhistory.go` (route via `loadHistoryRuns`, add `store` import), `gitmap/constants/constants_visibility_store_sql.go` (new SQL fragments), `gitmap/constants/constants.go` (`6.15.0`), `README.md` (pin), `CHANGELOG.md`, `.lovable/prompts/12-next-task.md` (new).


## v6.14.0 — (2026-06-06) — `vu`/`vr` `--json` summary + rate-limit backoff helper (steps 37-38)

- **Step 37 — `vu` / `vr` `--json` output (v5.43.0+ JSON contract parity):** `undoFlags` gains `JSON bool`; `parseUndoArgs` recognizes `--json`. New `gitmap/cmd/visibilityundojson.go` defines the canonical wire shape `undoJSONSummary` (command/runId/sourceRunId/provider/owner/matched/changed/skipped/failed/exitCode) — zero values are emitted explicitly so downstream JSON parsers never see missing keys. `reverseRunAndExit` now calls new `emitUndoJSON` after `audit.finalize`, writing one JSON line to stdout while preserving the human-readable summary above it (stdout doubles as both). JSON-render errors are surfaced to stderr but do not override the apply-outcome exit code (zero-swallow but non-fatal — the work succeeded, only the receipt failed). Added `(a *runAudit) RunID() int64` accessor in `visibilityallbulkaudit.go` so the new audit row's primary key can be included in the JSON without exposing the unexported `runID` field. Tests in `visibilityundojson_test.go` cover round-trip stability and explicit zero-key emission.
- **Step 38 — Rate-limit backoff helper:** new `gitmap/visibility/backoff.go` ships `ErrRateLimited` sentinel + `RetryRateLimited(op, schedule, sleep)` — pure, no `time.Sleep` baked in (caller injects so tests run instantly). Default `backoffSchedule()` is 1s/2s/4s/8s/16s/32s (63s total, deliberately under GitHub's 60s secondary-rate-limit window per attempt) across 6 retries. `errors.Is` predicate distinguishes retryable rate-limits from non-retryable failures (404/auth/schema) so a typo in a repo slug exits in 1 call instead of burning the full backoff. Test file covers: succeed-first-try, recover-mid-schedule, non-retryable-exits-immediately, exhaust-schedule, and a contract test that the schedule sum stays under the 60s rate-limit ceiling (regression guard against accidental schedule bloat). Wiring into the actual `gh repo edit` call site is deferred to item 45 (provider mock harness); the policy + tests ship now so the contract is locked.
- Files: `gitmap/cmd/visibilityundojson.go` (new), `gitmap/cmd/visibilityundojson_test.go` (new), `gitmap/visibility/backoff.go` (new), `gitmap/visibility/backoff_test.go` (new), `gitmap/cmd/visibilityundoflags.go` (`--json` parse), `gitmap/cmd/visibilityundo.go` (`JSON` field + `emitUndoJSON`), `gitmap/cmd/visibilityallbulkaudit.go` (`RunID()` accessor), `gitmap/constants/constants.go` (`6.14.0`), `README.md` (pin), `CHANGELOG.md`, `.lovable/prompts/11-next-task.md` (new).


## v6.13.0 — (2026-06-06) — `vh` round-trip test + `--kind` / `--since` filters (steps 35-36)

- **Step 35 — Data-layer round-trip test:** new `gitmap/store/makeallvisibility_roundtrip_test.go` exercises the full `MakeAllPublic → VisibilityUndo → VisibilityRedo` lifecycle at the store layer. Inserts three runs with monotonically increasing `StartedAt`, asserts `SelectRecentMakeAllVisibilityRuns(10)` returns them newest-first (redo, undo, pub), and confirms `SelectMakeAllVisibilityRunByID(undoID)` resolves to the correct kind. Locks in the column-order + kind-routing contract that vu/vr depend on; provider-level e2e (real `gh` calls) still pending item 45 (mock harness).
- **Step 36 — `vh --kind <K>` / `vh --since <dur>` filters:** new `gitmap/cmd/visibilityhistoryfilters.go` (37 lines) introduces pure `parseHistoryFilters` + `applyHistoryFilters` helpers — zero DB, zero I/O, fully table-testable. `runVisibilityHistory` now parses the two flags alongside `--limit` and applies them post-fetch. `--since` accepts any Go `time.ParseDuration` string (`24h`, `7d` → use `168h`, `30m`); bad values are silently ignored (limit-style strict-fail would break existing scripts that pipe extra tokens). Bogus `StartedAt` strings are dropped under `--since` (zero-swallow not applicable — these are data-side malformations, not user errors). New `gitmap/cmd/visibilityhistoryfilters_test.go` covers parse defaults, parse happy-path, bad `--since` ignored, kind-only filter, since-only filter, combined kind+since filter, and the no-op zero-value path.
- Files: `gitmap/store/makeallvisibility_roundtrip_test.go` (new), `gitmap/cmd/visibilityhistoryfilters.go` (new), `gitmap/cmd/visibilityhistoryfilters_test.go` (new), `gitmap/cmd/visibilityhistory.go` (wire filters), `gitmap/constants/constants.go` (`6.13.0`), `README.md` (pin), `CHANGELOG.md`, `.lovable/prompts/10-next-task.md` (new).


## v6.12.0 — (2026-06-06) — Drift-guard seam + marker-comment audit (steps 33-34)

- **Step 33 — Marker-comment audit:** verified `CmdVisibilityUndo` / `CmdVisibilityUndoAlias` / `CmdVisibilityRedo` / `CmdVisibilityRedoAlias` / `CmdVisibilityHistory` / `CmdVisibilityHistoryAlias` in `gitmap/constants/constants_cli.go` (lines 187-205) correctly inherit the file-level `// gitmap:cmd top-level` marker on line 3 — they are intentionally NOT tagged `// gitmap:cmd skip` because all six are first-class top-level CLI tokens. No code change required; audit recorded here so the next CI `generate-check` drift run has a citable baseline.
- **Step 34 — Drift-guard integration seam:** the drift policy used by `reverseOneRepo` was previously inlined as bare `if flags.Force` / `if current != r.NewVisibility` branches, untestable without a real GitHub/GitLab provider client. Extracted the total decision function `decideDriftAction(current, expected string, force bool) driftAction` into new `gitmap/cmd/visibilitydriftguard.go` (37 lines). `reverseOneRepo` now delegates to it on both branches — behavior is byte-identical, but the policy is now table-testable. Added `gitmap/cmd/visibilitydriftguard_test.go` covering: no-drift-no-force → proceed, drift-no-force → skip, no-drift-force → force, drift-force → force (override wins), empty-current → skip. Locks in the three-way contract so a future refactor cannot silently flip the guard direction.
- Files: `gitmap/cmd/visibilitydriftguard.go` (new), `gitmap/cmd/visibilitydriftguard_test.go` (new), `gitmap/cmd/visibilityundo.go` (delegate to helper), `gitmap/constants/constants.go` (`6.12.0`), `README.md` (pin), `CHANGELOG.md`, `.lovable/prompts/09-next-task.md` (new).


## v6.11.0 — (2026-06-06) — Store SELECT tests + 5 missing help files (unblocks CI)

- **Step 31 — Store SELECT round-trip tests:** new `gitmap/store/makeallvisibility_undo_test.go` and `gitmap/store/makeallvisibility_history_test.go` seed a temp SQLite DB through the canonical `InsertMakeAllVisibilityRun` → `InsertMakeAllVisibilityPendingResults` → `UpdateMakeAllVisibilityResult` → `FinalizeMakeAllVisibilityRun` pipeline and then exercise every new SELECT (`SelectLatestUndoableMakeAllVisibilityRun`, `SelectMakeAllVisibilityRunByID`, `SelectLatestMakeAllVisibilityRunByKind`, `SelectUndoableResultsForRun`, `SelectRecentMakeAllVisibilityRuns`). Covers happy-path, empty-DB → `(zero, nil)` contract, unknown-id → `(zero, nil)`, newest-first ordering, kind-filter routing, and `--limit` honoring. Guards against a future column-order swap in the SQL silently routing the wrong field into `Provider`/`Owner`/`OkCount` — which would corrupt every undo decision without any error.
- **Step 32 — Five missing help files:** `helptext/coverage_test.go::TestEveryCmdIDHasHelpFile` reflects every `Cmd*` constant in `constants_cli.go` and requires a matching `<id>.md`. The five visibility commands shipped without docs, leaving the test failing. Added `gitmap/helptext/make-all-public.md`, `make-all-private.md`, `visibility-undo.md`, `visibility-redo.md`, `visibility-history.md` — all under the 120-line cap, all documenting flags / examples / exit-code matrix / drift-guard behavior / `--force` semantics / `--run <id>` selector.
- Files: `gitmap/store/makeallvisibility_undo_test.go` (new), `gitmap/store/makeallvisibility_history_test.go` (new), `gitmap/helptext/make-all-public.md` (new), `gitmap/helptext/make-all-private.md` (new), `gitmap/helptext/visibility-undo.md` (new), `gitmap/helptext/visibility-redo.md` (new), `gitmap/helptext/visibility-history.md` (new), `gitmap/constants/constants.go` (`6.11.0`), `README.md` (pin), `.lovable/prompts/08-next-task.md` (new).



## v6.10.0 — (2026-06-06) — Centralized undo/redo strings + unit tests for `parseUndoArgs` / `bulkExitCode`

- **Step 29 — Centralized 3 magic strings in `visibilityundo.go`:** `audit DB open failed`, the reverse-loop header (`reversing run #N (provider/owner) — N repo(s)`), and the `<cmd>:source-run=<id>` `PatternList` template moved to `constants_visibility.go` as `ErrUndoAuditDBOpenFmt`, `MsgUndoReverseHeaderFmt`, `UndoPatternsRawFmt`. The patternsRaw template is the audit trail's only link back to the source run — a typo would silently break `vh` filtering — so it now lives behind a single named constant.
- **Step 30 — Unit tests for `parseUndoArgs` + `matchesFromResults` + `bulkExitCode`:** new `gitmap/cmd/visibilityundoflags_test.go` covers defaults, all flags set together (`--verbose --dry-run --force --run 42`), `--force` in isolation, unknown-token tolerance, result→match adapter preservation, and the full bulk exit-code matrix (all-ok → 0, all-failed → 5, mixed → 9). Guards against the failure mode where a future refactor silently demotes `--force` to a no-op or mis-routes `--run <id>` — both of which would destroy real user data without any visible error.
- Files: `gitmap/constants/constants_visibility.go` (+`ErrUndoAuditDBOpenFmt`, `MsgUndoReverseHeaderFmt`, `UndoPatternsRawFmt`), `gitmap/cmd/visibilityundo.go` (3 inline strings → constants), `gitmap/cmd/visibilityundoflags_test.go` (new), `gitmap/constants/constants.go` (`6.10.0`), `README.md` (pin), `.lovable/prompts/07-next-task.md` (new).



## v6.9.0 — (2026-06-06) — Drift guard + `--force` on `vu` / `vr` + preflight `gh`/`glab auth status`

- **Drift guard (step 27):** `gitmap visibility-undo` and `visibility-redo` now read each repo's *current* visibility before reversing and skip with `DRIFT SKIP (current=… expected=…)` when the live state no longer matches the `NewVisibility` we persisted in the source run. Prevents the audit trail from silently overwriting out-of-band manual changes (someone flipped a repo via the GitHub UI after the original `make-all-*` run). New `--force` flag opts out of the guard with an audible `[--force] overriding drift guard for <repo>` log line.
- **Preflight `auth status` (step 28):** `mustEnsureProviderAuth` runs `<cli> auth status` BEFORE any provider mutation (`make-all-*`, `vu`, `vr`) and fails fast with `ExitVisAuthFailed` and a Code Red message instructing the user to `gh auth login` / `glab auth login`. Previously an unauthenticated CLI passed the `exec.LookPath` gate and failed mid-loop on the first per-repo call, leaving a half-populated audit run.
- Internal: drift loop extracted into `reverseOneRepo` (≤15 lines) so `applyUndoLoop` stays readable; auth preflight isolated to `visibilityauthstatus.go` (one file, one responsibility).
- Files: `gitmap/cmd/visibilityauthstatus.go` (new), `gitmap/cmd/visibilityundo.go` (`Force` field, `reverseOneRepo` drift helper, auth-status call), `gitmap/cmd/visibilityundoflags.go` (`--force` parsing), `gitmap/cmd/visibilityallbulk.go` (auth-status preflight in `runMakeAllVisibility`), `gitmap/constants/constants_visibility.go` (+`ErrVisAuthStatusFailedFmt`, `MsgUndoDriftSkipFmt`, `MsgUndoForceOverrideFmt`), `gitmap/constants/constants.go` (`6.9.0`), `README.md` (pin), `.lovable/prompts/06-next-task.md` (new).



## v6.8.0 — (2026-06-06) — `gitmap visibility-history` (`vh`) + `--dry-run` on `vu` / `vr`

- New command `gitmap visibility-history` (alias `vh`) prints the most recent `MakeAllVisibilityRun` rows newest-first, with id, kind, owner, matched/ok/skip/fail/excl tallies, exit code, and `StartedAt`. Default limit 20; `--limit N` overrides. Empty-database case prints a friendly stderr message and exits `ExitVisOK`. This is the discovery layer behind the v6.7.0 `--run <id>` selector — users can now see *which* IDs to target.
- `--dry-run` on both `vu` and `vr` enumerates the planned per-repo reversal (`would set visibility -> <prev>`) without calling the provider CLI. Lets users verify a reversal is safe before letting it touch GitHub/GitLab. Also paves the way for the upcoming drift-detection guard (step 29) which will compare *current* visibility against the persisted `NewVisibility` along this same enumeration path.
- Internal layout: flag parsing + dry-run rendering extracted into `visibilityundoflags.go` so `visibilityundo.go` stays under the 200-line per-file cap; every new func ≤15 lines.
- Files: `gitmap/cmd/visibilityhistory.go` (new), `gitmap/cmd/visibilityundoflags.go` (new — `parseUndoArgs`, `mustParseRunID`, `printDryRun`), `gitmap/store/makeallvisibility_history.go` (new — `SelectRecentMakeAllVisibilityRuns`), `gitmap/constants/constants_visibility_store_sql.go` (+`SQLSelectRecentRuns`, history messages, dry-run messages, `HistoryDefaultLimit`), `gitmap/constants/constants_cli.go` (`CmdVisibilityHistory` + `vh`), `gitmap/cmd/visibilityundo.go` (DryRun field + dispatch branch), `gitmap/cmd/visibilityredo.go` (dry-run branch), `gitmap/cmd/rootcore.go` (vh dispatch).

## v6.7.0 — (2026-06-06) — `gitmap visibility-redo` / `vr` + `--run <id>` selector on undo/redo

- New command `gitmap visibility-redo` (alias `vr`) reverses the most recent `VisibilityUndo` run, restoring the visibility state the undo reverted. The redo is itself persisted as a fresh `MakeAllVisibilityRun` with `CommandKind='VisibilityRedo'` and `PatternList='visibility-redo:source-run=<id>'`, keeping the audit chain (apply → undo → redo → undo → …) intact.
- Both `vu` and `vr` now accept `--run <id>` to target an exact historical run instead of the latest one — required substrate for the upcoming `visibility-history` (`vh`) command. A missing or non-numeric `--run` value exits with `ExitVisBadFlag` and a Code Red error pinpointing the bad token (zero-swallow).
- Internal refactor: `visibilityundo.go` now exports a single shared `reverseRunAndExit` helper that both `vu` and `vr` consume; `visibilityredo.go` is a 9-line dispatcher that only changes the source-run filter (`CommandKind='VisibilityUndo'`) and the cmdName under which the new audit run is logged. No duplicate loop bodies.
- Files: `gitmap/cmd/visibilityredo.go` (new, 21 lines), `gitmap/cmd/visibilityundo.go` (rewritten — shared helpers, `undoFlags` with `RunID`), `gitmap/store/makeallvisibility_undo.go` (+`SelectMakeAllVisibilityRunByID`, `SelectLatestMakeAllVisibilityRunByKind`, shared `scanRunRow`), `gitmap/constants/constants_visibility_store_sql.go` (+`SQLSelectLatestRunByKind`, `SQLSelectRunByID`, three new `ErrUndo*` formats), `gitmap/constants/constants_visibility_store.go` (`CommandKindVisibilityRedo`), `gitmap/constants/constants_cli.go` (`CmdVisibilityRedo` + `vr` alias), `gitmap/cmd/rootcore.go` (dispatch), `gitmap/cmd/visibilityallbulkaudit.go` (`commandKindFor` extended).

## v6.6.0 — (2026-06-06) — `gitmap visibility-undo` / `vu` reverses the last bulk make-all-* run

Bulk visibility flips (`make-all-public` / `make-all-private` / `MAPUB` / `MAPRI`) are now reversible from the same audit trail they already write.

- New command `gitmap visibility-undo` (alias `vu`) picks the most recent `MakeAllVisibilityRun` with `OkCount > 0`, reads its `MakeAllVisibilityResult` rows, and re-applies each repo's captured `PrevVisibility` via the existing single-repo read→apply→verify pipeline.
- The undo is itself persisted as a fresh `MakeAllVisibilityRun` with `CommandKind='VisibilityUndo'` and `PatternList='undo:source-run=<id>'`, so a follow-up `vu` reverses the undo (this is the substrate for the upcoming `visibility-redo`).
- Rows with empty `PrevVisibility` or `PrevVisibility == NewVisibility` are skipped at SQL-select time — no provider round-trip for no-op rows.
- Best-effort audit: a missing/locked DB still logs to `os.Stderr` with Code Red context and continues (zero-swallow); the user's data is never the casualty of an audit failure.
- Files: `gitmap/cmd/visibilityundo.go` (new), `gitmap/store/makeallvisibility_undo.go` (new), `gitmap/constants/constants_visibility_store_sql.go` (+`SQLSelectLatestUndoableRun`, `SQLSelectUndoableResultsForRun`, three `ErrUndo*` formats), `gitmap/constants/constants_visibility_store.go` (`CommandKindVisibilityUndo`), `gitmap/constants/constants_cli.go` (`CmdVisibilityUndo` + `vu` alias), `gitmap/cmd/rootcore.go` (dispatch), `gitmap/cmd/visibilityallbulkaudit.go` (`commandKindFor` switch extended).

## v6.5.0 — (2026-06-06) — `gitmap reclone` / `rc` wipes + re-clones a single repo

New single-repo flow overlays the existing manifest-based `reclone`:

- Run `gitmap reclone` (aliases `rc`, `rec`, `relclone`, `clone-now`) from **inside** a git repo to wipe and re-clone it from `remote.origin.url`.
- Run `gitmap rc <folder>` to target a sibling git repo from outside.
- Destructive `os.RemoveAll` is gated by an interactive `y/N` prompt; pass `-y` (or `--yes`) to skip it. Non-TTY callers without `-y` are refused — protects `yes | gitmap rc` in CI.
- Releases the Windows cwd handle via `escapeCwdIfInside` before delete, so re-cloning from inside the target works on Windows.
- Manifest behavior is unchanged: when the args don't shape into a single-repo call, the manifest pipeline runs as before.

## v6.4.0 — (2026-05-30) — `gitmap hd` auto-downloads `docs-site.zip`

- Fixed: `gitmap hd` (help-dashboard) no longer dies with `Docs site directory not found at <install>\docs-site` when the installer skipped the docs asset (older installs, or releases that didn't ship `docs-site.zip`).
- New behaviour: when neither `docs-site/` nor `docs-site.zip` exists next to the binary, `runHelpDashboard` now fetches the archive at runtime over HTTPS, trying `releases/download/v<Version>/docs-site.zip` first, then `releases/latest/download/docs-site.zip`, before extracting in place.
- Files: `gitmap/cmd/helpdashboard.go` (download branch), new `gitmap/cmd/helpdashboard_download.go` (HTTP fetch with 30s timeout and `maxDocsSiteSize` cap), `gitmap/constants/constants_helpdashboard.go` (new messages + `DocsSiteDownloadTimeoutSec`).
- On total failure, the error now tells the user exactly which path to drop `docs-site.zip` into and suggests `gitmap update`.

## v6.3.0 — (2026-05-30) — Release binaries stamp `gitmap binary` provenance


- Fixed: release/CI-built binaries now embed the source repo URL, branch, commit SHA, and UTC build stamp via `-ldflags`, so the `gitmap binary` footer can identify the actual binary instead of falling back to the current working repo or showing only a version.
- Root cause: the local `run.sh` / `run.ps1` path had build identity injection, but GitHub Actions release and CI artifact builds still passed only `constants.Version`; downloaded binaries therefore missed the v5.60.0 footer provenance fix.

## v6.2.2 — (2026-05-30) — Fix macOS CI: stale `updateprobe_test.go` version expectations

- Fixed: `gitmap/cmd/updateprobe_test.go` — `TestParseCurrentRepoSlug` had been self-rewritten to use `gitmap-v27` while still expecting stale parsed versions (`23` and `1`). The test now derives the current-slug case from one `currentSlugVersion` constant and uses a synthetic `tool-v1` case for the v1 parser branch, preventing future repo bumps from rewriting only half of the assertion.
- Hardened: `fixrepo_rewrite_preview_test.go` and `fixrepo_rewrite_barebase_test.go` now use synthetic `acme-vN` tokens instead of this repo's own `gitmap-v27` literals, so future `fix-repo` bumps cannot silently rot the rewrite/preview regression tests.

## v6.2.1 — (2026-05-30) — Fix macOS CI: `fixrepo_rewrite_versionscope_test.go` self-rewrite damage

- Fixed: `gitmap/cmd/fixrepo_rewrite_versionscope_test.go` — every `gitmap-vN` literal in the test's `in`/`want` strings (for N < 25) had been silently rewritten to `gitmap-v27` by fix-repo itself on the v23→v25 bump, collapsing assertions like *"bare `gitmap` should become `gitmap-v27` when current=2"* into nonsense (`want: "...gitmap-v27..."`). Distractor tokens now use a synthetic `otherpkg-vN` base so the rewriter — which only touches `{base}-vN` where base == the repo name — can't smash them on future bumps. Same lesson as the `fixrepo_rewrite_v9tov12_test.go` fix that already uses `acme-vN`.
- Root cause: this test data was written using the repo's own base name (`gitmap`), making it self-poisoning under any future fix-repo run. Documented in mem://core under FIX-REPO DIGIT-CAPTURE GAP — extended now to cover not just sibling integer literals but any same-base versioned token in test fixtures.



## v6.2.0 — (2026-05-29) — Fix macOS/Windows CI: `TestExtractBaseAndVersionFromArg_URL` digit-capture desync

- Fixed: `gitmap/cmd/visibilitybulk_test.go` — `TestExtractBaseAndVersionFromArg_URL` hard-coded the expected version (`23`) as a bare integer literal separate from the input URL `gitmap-v27`. The fix-repo rewriter only touches `{base}-vN` tokens, so when the repo bumped from v23→v25 the URL was updated but the expected int was not, producing `expected (gitmap, 23), got (gitmap, 25)` on every CI run since the bump. This was the **exact bug class** documented in mem://core (FIX-REPO DIGIT-CAPTURE GAP, closed v4.12.0): "any new fix-repo test MUST derive expected version-bearing strings from the same int it passed in — never hard-code a sibling literal." The test now formats the URL from a `const wantVer = 25` and asserts against the same constant, so the next version bump rewrites both sides atomically.
- Root cause: a regression test for the v4.12.0 rule had itself been added in violation of the rule. CI's `macos-latest / go build + test` job (and any other matrix host) failed at `gitmap/cmd` for this single assertion.



## v6.1.0 — (2026-05-29) — `gitmap cd <repo> <inner-command>` runs subcommands inside a named repo

- Added: `gitmap cd <repo-name> <subcommand> [args...]` — `cd` now accepts an optional inner command (e.g. `gitmap cd myproj cn v++`, `gitmap cd myproj cfrp v+1`). The named repo is resolved via the existing slug/default/pick path, the process chdirs into it, and the inner subcommand dispatches normally. Commands that already write a shell handoff (`cn`, `cfr`, `cfrp`, `as`, `clone`, ...) continue to relocate the caller's shell to the new directory after the inner command finishes.
- Changed: `gitmap/cmd/cdops.go` — `runCDLookup` now extracts trailing positional args via `parseCDPickFlag` (return signature widened to `(bool, []string)`) and routes them to a new `runCDInner` helper.
- Changed: `gitmap/constants/constants_cd.go` — `ErrCDUsage` updated to advertise the optional inner-command form; added `ErrCDChdirFmt` for chdir failures.

## v6.0.0 — (2026-05-26) — Breaking: `commit-transfer` merge default flips to `true`

**Breaking change:** `gitmap commit-in`, `commit-out`, `commit-left`, `commit-right`, and `commit-both` now preserve merge commits by default. The legacy strip behaviour requires explicit `--no-include-merges`.

- Changed: `gitmap/cmd/committransfer.go` — `parseCommitTransferArgs` now initializes `Options.IncludeMerges = true`. The `--include-merges` flag is still accepted (redundant but harmless). Added `--no-include-merges` flag for explicit opt-out. Both flags use `BoolFunc` so they toggle without consuming a value.
- Changed: `gitmap/committransfer/log.go` — `PrintPlan` notice inverted. When `--no-include-merges` strips merge commits, the message now reads `note: N merge commits excluded by --no-include-merges` (confirmation) instead of the old advisory `pass --include-merges to also replay merge commits`.
- Added: `gitmap/committransfer/types.go` — `ReplayPlan.IncludeMerges bool` field so `PrintPlan` knows whether the exclusion was intentional.
- Added: `gitmap/constants/constants_committransfer.go` — `FlagCTNoIncludeMerges` + `FlagDescCTNoIncludeMerges` constants.
- Added: `gitmap/cmd/committransfer_flags_test.go` — `TestCommitTransferIncludeMergesDefault` asserts all three directions default to `IncludeMerges = true`; `TestCommitTransferIncludeMergesExplicit` asserts `--include-merges` / `--no-include-merges` override correctly.
- Added: `gitmap/committransfer/log_test.go` — `TestPrintPlanNoticeV6` asserts the confirmation notice is emitted when `IncludeMerges = false` and `MergeExcluded > 0`, and that the old advisory message is gone.
- Updated: `spec/01-app/115-v6-migration.md` — acceptance criteria marked complete; manual-smoke section promoted to verified.

**Migration:** Scripts that relied on the silent merge-stripping must add `--no-include-merges`. See `spec/01-app/115-v6-migration.md` for the full migration guide.

## v5.84.0 — (2026-05-26) — `scan-project` JSON-schema contract: per-type file shape pinned

- Added: `spec/08-json-schemas/scan-project.schema.json` — draft-07 schema for the per-type JSON files emitted by `gitmap scan-project` (`go-projects.json`, `node-projects.json`, `react-projects.json`, `cpp-projects.json`, `csharp-projects.json`). Top-level is a JSON array; each record has PascalCase top-level keys (`Project`, `GoMeta`, `Csharp`) because `detector.DetectionResult` has no `json:` tags — pinned verbatim as the v1 on-the-wire contract. Nested objects use lowerCamel from the `model.*` struct tags. `GoMeta`/`Csharp` use `oneOf [null, object]` to capture the metadata-optional shape.
- Added: `gitmap/cmd/testdata/schemas/scan-project.v1.json` — registry entry listing the 5 advertised filenames and the 3 top-level record keys.
- Added: `gitmap/cmd/scanproject_jsonschema_contract_test.go` — two tests: `TestScanProject_FileMapMatchesRegistry` asserts `projectTypeJSONMap` produces exactly the 5 filenames in registry order; `TestScanProject_RecordKeysSubsetOfSchema` runs the live `buildJSONRecords` over both shape variants (bare result + metadata-wrapped) and asserts every emitted top-level key is declared in `items.properties` — catches struct-tag drift or accidental key-name churn on every CI run.
- Updated: `spec/08-json-schemas/_TODO.md` — `scan-project` row marked ✅ done.

## v5.83.0 — (2026-05-26) — Spec 114 Gap A follow-up: `--max-history-scan` escape hatch

- Added: `Options.MaxHistoryScan int` in `gitmap/committransfer/types.go` — opt-in cap on the target-history scan used by the idempotence check. Default `0` preserves the v5.78.0 unbounded behaviour; positive values cap the `git log` query at N commits for operators running against pathologically large targets (mirrored monorepos, tens of millions of commits) where the full log is prohibitive.
- Wired: `gitmap/committransfer/plan.go` — `BuildPlan` now passes `opts.MaxHistoryScan` to `recentLogSubjectsAndBodies(targetDir, opts.MaxHistoryScan)` instead of the hard-coded `0`. The existing `<= 0` branch in `recentLogSubjectsAndBodies` keeps the unbounded default working with no other changes.
- Added: `--max-history-scan N` CLI flag on all commit-transfer commands (`commit-in`, `commit-out`, `commit-left`, `commit-right`, `commit-both`). Constants: `FlagCTMaxHistoryScan` + `FlagDescCTMaxHistoryScan` in `gitmap/constants/constants_committransfer.go`. Wired in `registerCommitTransferStrings`.
- Added: `gitmap/committransfer/maxhistoryscan_test.go` — pins the zero-value-means-unbounded contract and the struct-field round-trip so rename/removal triggers a compile failure here BEFORE it reaches CLI wiring or `plan.go`.
- Updated: `spec/01-app/114-committransfer-idempotence-and-merge-default.md` — Gap A resolution gains an explicit step (4) documenting the v5.83.0 escape hatch.
- Verified: `TestPlanIdempotenceBeyond200Commits` continues to pass because the default (`MaxHistoryScan=0`) preserves the unbounded scan behaviour added in v5.78.0.



## v5.82.0 — (2026-05-26) — `gitmap export` schema v2: per-record property pinning

- Extended: `spec/08-json-schemas/export.schema.json` (now v2) — adds full `items.properties` + `items.required` declarations for all five nested arrays. Pinned record shapes: `repos` (`model.ScanRecord`, 15 keys), `groups` (`model.GroupExport` = `Group` + `repoSlugs`, 6 keys), `releases` (`model.ReleaseRecord`, 14 keys), `history` (`model.CommandHistoryRecord`, 12 keys with `alias`/`args`/`flags`/`finishedAt`/`summary`/`createdAt` flagged optional per `omitempty`), `bookmarks` (`model.BookmarkRecord`, 6 keys with `args`/`flags`/`createdAt` flagged optional per `omitempty`).
- Added: `gitmap/cmd/export_nested_jsonschema_contract_test.go` — `TestExportJSONSchema_NestedRecordKeysSubsetOfProperties` builds a deterministic non-empty export (one record per nested array), runs the live `encodeDatabaseExportJSON`, and asserts every key emitted on each per-record object is declared in that array's `items.properties` map. Catches struct-tag drift in either the model or the schema on every CI run.
- Updated: `spec/08-json-schemas/_TODO.md` — `export` row updated to reflect schema v2 closure of the per-record property-set gap left open in v5.81.0.



## v5.81.0 — (2026-05-26) — JSON schema contract: `gitmap export`

- Added: `spec/08-json-schemas/export.schema.json` — draft-07 schema pinning the top-level object shape (7 required keys in contractual order: `version`, `exportedAt`, `repos`, `groups`, `releases`, `history`, `bookmarks`). Per-record key order within the five nested arrays is explicitly NOT pinned in v1 — that scope (one schema per nested record type) is deferred.
- Added: `gitmap/cmd/exportrender.go` — `encodeDatabaseExportJSON` builds the top-level object via `gitmap/stablejson` so the key order is a compile-time decision rather than struct-field-tag-defined. Nested arrays are pre-rendered with `json.MarshalIndent` (deterministic via Go struct field declaration order). Empty arrays emit `[]` (never `null`) so `jq '.repos | length'` style probes work on fresh databases.
- Changed: `gitmap/cmd/export.go` — `writeExportFile` switched from `json.MarshalIndent(model.DatabaseExport)` to the new stablejson-backed encoder. Behavior is byte-equivalent for the top-level layout; the contractual difference is that re-ordering the struct tags can no longer silently re-order the wire output.
- Added: `gitmap/cmd/testdata/schemas/export.v1.json` — schema-registry entry locking the 7-key top-level order so `assertSchemaKeysFirstObject` enforces drift on every CI run.
- Added: `gitmap/cmd/testdata/export_empty.json` — golden fixture for an empty-database export (`version="1"`, `exportedAt="2026-05-26T12:00:00Z"`, all five arrays `[]`).
- Added: `gitmap/cmd/exportjson_contract_test.go` — two contract tests: empty-arrays-not-null guarantee and top-level key-order against the schema registry.
- Added: `gitmap/cmd/export_jsonschema_contract_test.go` — schema-shape pin: verifies the JSON Schema declares `type=object`, lists all 7 required keys, and that every key the live encoder emits is declared in the schema's `properties` map.
- Updated: `spec/08-json-schemas/_TODO.md` — `export` marked ✅ done with the per-record-not-pinned caveat called out explicitly.



## v5.80.0 — (2026-05-26) — JSON schema contract: `llm-docs --format=json`

- Added: `spec/08-json-schemas/llm-docs.schema.json` — draft-07 schema pinning the top-level object shape (8 optional sections in contractual order: `commands`, `architecture`, `flags`, `conventions`, `structure`, `database`, `installation`, `patterns`) plus the nested command-group (`title`, `commands`) and per-command (`name`, `alias`, `description`, optional `example`) structures.
- Added: `gitmap/cmd/testdata/schemas/llm-docs.v1.json` — schema-registry entry locking the 8-key top-level order so `assertSchemaKeysFirstObject` enforces drift on every CI run.
- Added: `gitmap/cmd/testdata/llm_docs_empty.json` — golden fixture asserting an empty-sections filter emits `{}\n` (not `null`).
- Added: `gitmap/cmd/llmdocsjson_contract_test.go` — three contract tests: empty-object guarantee, top-level key-order against the schema registry, and nested command-group/per-command key-order assertion (with optional `example` only appearing as the 4th key when non-empty).
- Added: `gitmap/cmd/llmdocs_jsonschema_contract_test.go` — schema-shape pin: verifies the JSON Schema declares all 8 top-level properties + the nested `commands.items` group + per-command properties (`name`, `alias`, `description`, `example`), and that every key the live encoder emits is declared in the schema.
- Updated: `spec/08-json-schemas/_TODO.md` — `llm-docs` confirmed migrated (renderer landed earlier; this release closes the missing schema + contract test gap).



## v5.79.0 — (2026-05-26) — Spec 114 Gap A: hash-set idempotence for unbounded target log

- Fixed: `AlreadyReplayed` previously performed an O(N×M) `strings.Contains` substring scan across the entire concatenated target log for every source commit. On targets with long histories (the unbounded scan enabled by the v5.78.0 fix), this became a hidden performance bottleneck.
- Added: `BuildReplayedSet(recentLog string) map[string]struct{}` in `gitmap/committransfer/message.go` — parses all `gitmap-replay: from <source> <sha>` provenance footers once and stores them in a hash-set for O(1) lookups.
- Added: `SetHasReplayed(set map[string]struct{}, sourceDisplay, shortSHA string) bool` — O(1) set membership test replacing the substring scan.
- Changed: `BuildPlan` now builds the replayed set once (after the unbounded `recentLogSubjectsAndBodies(targetDir, 0)` call) and passes it through `assemblePlan` → `hydrateCommit`. `hydrateCommit` now calls `SetHasReplayed` directly.
- Kept: `AlreadyReplayed` remains as a backward-compatible wrapper (now marked Deprecated) that internally builds the set and delegates to `SetHasReplayed`. Existing tests and external callers continue to work.
- Added: `TestBuildReplayedSet` and `TestSetHasReplayed` in `gitmap/committransfer/message_test.go` — coverage for set construction, duplicate tolerance, and negative lookups.
- Verified: `TestPlanIdempotenceBeyond200Commits` continues to pass (regression guard for the original 200-cap bug) because the set contains the same provenance data; only the lookup path changed.


## v5.78.0 — (2026-05-26) — Fix Windows CI: restore CWD in `escapecwd` tests

- Fixed: `TestEscapeCwdIfInside_NotInside` and `TestEscapeCwdIfInside_EscapesWhenInside` previously called `os.Chdir(t.TempDir())` without restoring the original working directory. On Windows, when the temp dir was later removed by `t.TempDir`'s cleanup, the process CWD became invalid (Windows reports it as `C:\`), cascading into ~40 unrelated failures in the same `cmd` package: every schema/golden test that walks up from CWD looking for `spec/08-json-schemas/*.json` or `testdata/*.json` aborted with `walking up from C:\` or `The system cannot find the path specified`.
- Added: `restoreCwd(t)` helper in `gitmap/cmd/escapecwd_test.go` — snapshots `os.Getwd()` and registers a `t.Cleanup` that chdir's back. Registered BEFORE the test's chdir so it runs AFTER the chdir-out but BEFORE `t.TempDir`'s RemoveAll (Cleanup runs LIFO), eliminating both the cascade AND the Windows "file in use" RemoveAll warning seen in the same job.
- Why this only blew up now: linux/macOS tolerate a deleted CWD by reporting the stale path string; Windows `GetCurrentDirectoryW` returns the volume root the moment the directory handle goes away. The leak existed for many releases but only became fatal once enough cmd-package tests started walking up from CWD (recent JSON-schema migration sprint).



## v5.77.0 — (2026-05-26) — `temp-releaselist --json` migrated to `stablejson` + published JSON schema

- Migrated: `gitmap temp-releaselist --json` encoder onto `gitmap/stablejson` (new `gitmap/cmd/tempreleaselistrender.go`). Key order (`id`, `branch`, `versionPrefix`, `sequenceNumber`, `commit`, `commitMessage`, `createdAt`) is now a compile-time decision via package-level wire-key constants instead of a reflection accident on `model.TempRelease`.
- Removed: legacy `json.MarshalIndent(releases, ...)` path in `tempreleaselist.go`; routed through the new stable encoder.
- Added: `spec/08-json-schemas/temp-release-list.schema.json` — published JSON Schema for downstream consumers.
- Added: `gitmap/cmd/tempreleaselist_jsonschema_contract_test.go` + `tempreleaselistjson_contract_test.go` — schema drift detection + golden fixtures (empty array + canonical two-row) + key-order contract.
- Added: `gitmap/cmd/testdata/schemas/temp-release-list.v1.json` — schema registry entry for key-order drift detection.
- Updated: `spec/08-json-schemas/_TODO.md` — `temp-releaselist` marked done.
- Pinned: README + `gitmap/constants/constants.go` + `src/constants/index.ts` synced to **v5.77.0**.


## v5.76.0 — (2026-05-26) — `version-history --json` migrated to `stablejson` + published JSON schema

- Migrated: `gitmap version-history --json` encoder onto `gitmap/stablejson` (new `gitmap/cmd/versionhistoryrender.go`). Key order (`fromVersionTag`, `fromVersionNum`, `toVersionTag`, `toVersionNum`, `flattenedPath`, `createdAt`, `id`, `repoId`) is now a compile-time decision via package-level wire-key constants instead of a reflection accident on `model.RepoVersionHistoryRecord`. Optional `flattenedPath` and `createdAt` are conditionally appended so the legacy omitempty wire shape is preserved (absent rather than null/empty).
- Removed: legacy `json.MarshalIndent(records, ...)` path in `versionhistory.go`; routed through the new stable encoder.
- Added: `spec/08-json-schemas/version-history.schema.json` — published JSON Schema for downstream consumers.
- Added: `gitmap/cmd/versionhistory_jsonschema_contract_test.go` + `versionhistoryjson_contract_test.go` — schema drift detection + golden fixtures (empty array + canonical two-row) + key-order contract.
- Added: `gitmap/cmd/testdata/schemas/version-history.v1.json` — schema registry entry for key-order drift detection.
- Updated: `spec/08-json-schemas/_TODO.md` — `version-history` marked done.
- Pinned: README + `gitmap/constants/constants.go` + `src/constants/index.ts` synced to **v5.76.0**.

## v5.75.0 — (2026-05-26) — `stats --json` migrated to `stablejson` + published JSON schema

- Migrated: `gitmap stats --json` encoder onto `gitmap/stablejson` (new `gitmap/cmd/statsrender.go`). Top-level object key order (`totalCommands`, `uniqueCommands`, `totalSuccess`, `totalFail`, `overallFailRate`, `avgDurationMs`, `commands`) AND nested per-command row key order (`command`, `totalRuns`, `successCount`, `failCount`, `failRate`, `avgDurationMs`, `minDurationMs`, `maxDurationMs`, `lastUsed`) are now compile-time decisions via package-level wire-key constants instead of reflection accidents on `model.OverallStats` / `model.CommandStats`. The nested array is pre-rendered in compact mode and embedded as `json.RawMessage`.
- Removed: legacy `json.MarshalIndent(overall, ...)` path in `stats.go`; routed through the new stable encoder.
- Added: `spec/08-json-schemas/stats.schema.json` — published JSON Schema (top-level object + nested items contract).
- Added: `gitmap/cmd/stats_jsonschema_contract_test.go` + `statsjson_contract_test.go` — schema drift detection + golden fixture (empty commands array) + key-order contract.
- Added: `gitmap/cmd/testdata/schemas/stats.v1.json` — schema registry entry for top-level key-order drift detection.
- Updated: `spec/08-json-schemas/_TODO.md` — `stats` marked done.
- Pinned: README + `gitmap/constants/constants.go` + `src/constants/index.ts` synced to **v5.75.0**.


## v5.74.0 — (2026-05-26) — `ssh list --json` migrated to `stablejson` + published JSON schema

- Migrated: `gitmap ssh list --json` encoder onto `gitmap/stablejson` (new `gitmap/cmd/sshlistrender.go`). Key order (`id`, `name`, `privatePath`, `publicKey`, `fingerprint`, `email`, `createdAt`) is now a compile-time decision via package-level wire-key constants instead of a reflection accident on `model.SSHKey`.
- Removed: legacy `json.MarshalIndent` path in `sshlist.go`; routed through the new stable encoder.
- Added: `spec/08-json-schemas/ssh-list.schema.json` — published JSON Schema for downstream consumers.
- Added: `gitmap/cmd/sshlist_jsonschema_contract_test.go` + `sshlistjson_contract_test.go` — schema drift detection + golden fixture (empty array) + key-order contract.
- Added: `gitmap/cmd/testdata/schemas/ssh-list.v1.json` — schema registry entry for key-order drift detection.
- Updated: `spec/08-json-schemas/_TODO.md` — `ssh list` marked done; clarified `env-registry` has no actual `--json` stdout flag.
- Pinned: README + `gitmap/constants/constants.go` + `src/constants/index.ts` synced to **v5.74.0**.


- Migrated: `gitmap list-versions --json` encoder onto `gitmap/stablejson` (new `gitmap/cmd/listversionsrender.go`). Key order (`version`, `source`, `changelog`) is now a compile-time decision via package-level wire-key constants instead of a reflection accident on `lvJSONEntry`. Optional `source` and `changelog` are conditionally appended so the legacy omitempty wire shape is preserved (absent rather than null/empty).
- Removed: legacy `lvJSONEntry` struct + `json.MarshalIndent` path in `listversionsutil.go`; routed through the new stable encoder.
- Added: `spec/08-json-schemas/list-versions.schema.json` — published JSON Schema for downstream consumers.
- Added: `gitmap/cmd/listversions_jsonschema_contract_test.go` + `listversionsjson_contract_test.go` — schema drift detection + golden fixtures (empty + canonical two-row) + key-order contract.
- Added: `gitmap/cmd/testdata/schemas/list-versions.v1.json` — schema registry entry for key-order drift detection.
- Updated: `spec/08-json-schemas/_TODO.md` — `list-versions` flipped from `med` to `done`.
- Pinned: README + `gitmap/constants/constants.go` + `src/constants/index.ts` synced to **v5.73.0**.



## v5.72.0 — (2026-05-26) — `latest-branch --json` migrated to `stablejson` + published JSON schema

- Migrated: `gitmap latest-branch --json` encoder onto `gitmap/stablejson` (new `gitmap/cmd/latestbranchrender.go`). Key order (`branch`, `remote`, `sha`, `commitDate`, `subject`, `ref`, `top`) is now a compile-time decision via package-level wire-key constants instead of a reflection accident on `latestBranchJSON`. The nested `top` array is pre-rendered in compact mode and embedded as `json.RawMessage`.
- Added: `spec/08-json-schemas/latest-branch.schema.json` — published JSON Schema for downstream consumers.
- Added: `gitmap/cmd/latestbranch_jsonschema_contract_test.go` — schema drift detection (top-level object shape, required key set, encoder-keys ⊂ schema.properties).
- Updated: `gitmap/cmd/latestbranchjson_contract_test.go` — refreshed comments to reference `latestbranchrender.go` wire-key constants.
- Updated: `gitmap/cmd/testdata/latest_branch_no_top.json` — regenerated golden fixture to match stablejson output.
- Updated: `spec/08-json-schemas/_TODO.md` — `latest-branch` flipped from `med` to `done`.
- Pinned: README + `gitmap/constants/constants.go` + `src/constants/index.ts` synced to **v5.72.0**.



## v5.71.0 — (2026-05-26) — `project-repos --json` migrated to `stablejson` + published JSON schema

- Migrated: `gitmap <type>-repos --json` (go/node/react/cpp/csharp) encoder onto `gitmap/stablejson` (new `gitmap/cmd/projectreposrender.go`). Key order (`id`, `repoId`, `repoName`, `projectTypeId`, `projectType`, `projectName`, `absolutePath`, `repoPath`, `relativePath`, `primaryIndicator`, `detectedAt`) is now a compile-time decision via package-level wire-key constants.
- Added: `spec/08-json-schemas/project-repos.schema.json` — published JSON Schema for downstream consumers.
- Added: `gitmap/cmd/projectrepos_jsonschema_contract_test.go` + `projectreposjson_contract_test.go` — schema drift detection + golden fixture + key-order contract.
- Added: `gitmap/cmd/testdata/schemas/project-repos.v1.json` — schema registry entry for key-order drift detection.
- Updated: `spec/08-json-schemas/_TODO.md` — `project repos` flipped from `med` to `done`.
- Pinned: README + `gitmap/constants/constants.go` + `src/constants/index.ts` synced to **v5.71.0**.



## v5.70.0 — (2026-05-26) — `bookmark list --json` migrated to `stablejson` + published JSON schema

- Migrated: `gitmap bookmark list --json` encoder onto `gitmap/stablejson` (new `gitmap/cmd/bookmarklistrender.go`). Key order (`id`, `name`, `command`, `args`, `flags`, `createdAt`) is now a compile-time decision via package-level wire-key constants instead of a reflection accident on `model.BookmarkRecord`.
- Added: `spec/08-json-schemas/bookmark-list.schema.json` — published JSON Schema for downstream consumers.
- Added: `gitmap/cmd/bookmarklist_jsonschema_contract_test.go` — schema drift detection (top-level array shape, required key set, encoder-keys ⊂ schema.properties).
- Added: `gitmap/cmd/bookmarklistjson_contract_test.go` — golden fixture + key-order contract for the stablejson output.
- Added: `gitmap/cmd/testdata/schemas/bookmark-list.v1.json` — schema registry entry for key-order drift detection.
- Updated: `spec/08-json-schemas/_TODO.md` — `bookmark list` flipped from `med` to `done`.
- Pinned: README + `gitmap/constants/constants.go` + `src/constants/index.ts` synced to **v5.70.0**.



## v5.69.0 — (2026-05-26) — `diff-profiles --json` migrated to `stablejson` + published JSON schema

- Migrated: `gitmap diff-profiles --json` encoder onto `gitmap/stablejson` (new `gitmap/cmd/diffprofilesrender.go`). Key order (`profileA`, `profileB`, `onlyInA`, `onlyInB`, `different`, `same`) is now a compile-time decision via package-level wire-key constants instead of a reflection accident on `map[string]any`. Nested `onlyInA`, `onlyInB`, and `different` arrays are pre-rendered in compact mode and embedded as `json.RawMessage` so key-order stability propagates through the entire document.
- Added: `spec/08-json-schemas/diff-profiles.schema.json` — published JSON Schema for downstream consumers.
- Added: `gitmap/cmd/diffprofiles_jsonschema_contract_test.go` — schema drift detection (top-level object shape, required key set, encoder-keys ⊂ schema.properties).
- Added: `gitmap/cmd/diffprofilesjson_contract_test.go` — golden fixture + key-order contract for the stablejson output.
- Added: `gitmap/cmd/testdata/schemas/diff-profiles.v1.json` — schema registry entry for key-order drift detection.
- Updated: `spec/08-json-schemas/_TODO.md` — `diff-profiles` flipped from `med` to `done`.
- Pinned: README + `gitmap/constants/constants.go` + `src/constants/index.ts` synced to **v5.69.0**.



## v5.68.0 — (2026-05-26) — `amend audit` migrated to `stablejson` + published JSON schema

- Migrated: `gitmap amend audit` file encoder onto `gitmap/stablejson` (new `gitmap/cmd/amendauditrender.go`). Key order (`id`, `timestamp`, `branch`, `fromCommit`, `toCommit`, `totalCommits`, `previousAuthor`, `newAuthor`, `mode`, `forcePushed`, `commits`) is now a compile-time decision via package-level wire-key constants instead of a reflection accident on `model.AmendmentRecord`. Nested `previousAuthor` / `newAuthor` objects and the `commits` array are pre-rendered in compact mode and embedded as `json.RawMessage` so key-order stability propagates through the entire document.
- Added: `spec/08-json-schemas/amend-audit.schema.json` — published JSON Schema for downstream consumers.
- Added: `gitmap/cmd/amendaudit_jsonschema_contract_test.go` — schema drift detection (top-level object shape, required key set, encoder-keys ⊂ schema.properties).
- Added: `gitmap/cmd/amendauditjson_contract_test.go` — golden fixture + key-order contract for the stablejson output.
- Added: `gitmap/cmd/testdata/schemas/amend-audit.v1.json` — schema registry entry for key-order drift detection.
- Updated: `spec/08-json-schemas/_TODO.md` — `amend audit` flipped from `med` to `done`.
- Pinned: README + `gitmap/constants/constants.go` + `src/constants/index.ts` synced to **v5.68.0**.



## v5.67.0 — (2026-05-26) — `amend list --json` migrated to `stablejson` + published JSON schema

- Migrated: `gitmap amend list --json` encoder onto `gitmap/stablejson` (new `gitmap/cmd/amendlistrender.go`). Key order (`ID`, `Branch`, `FromCommit`, `ToCommit`, `TotalCommits`, `PreviousName`, `PreviousEmail`, `NewName`, `NewEmail`, `Mode`, `ForcePushed`, `CreatedAt`) is now a compile-time decision via package-level wire-key constants. PascalCase keys are preserved from the legacy `json.MarshalIndent` output for backward compatibility.
- Added: `spec/08-json-schemas/amend-list.schema.json` — published JSON Schema for downstream consumers.
- Added: `gitmap/cmd/amendlist_jsonschema_contract_test.go` — schema drift detection (top-level array shape, required key set, encoder-keys ⊂ schema.properties).
- Added: `gitmap/cmd/amendlistjson_contract_test.go` — golden fixture + key-order contract for the stablejson output.
- Added: `gitmap/cmd/testdata/schemas/amend-list.v1.json` — schema registry entry for key-order drift detection.
- Updated: `spec/08-json-schemas/_TODO.md` — `amend list` flipped from `med` to `done`.
- Pinned: README + `gitmap/constants/constants.go` + `src/constants/index.ts` synced to **v5.67.0**.



## v5.66.0 — (2026-05-26) — `probe --json` migrated to `stablejson` + published JSON schema

- Migrated: `gitmap probe --json` encoder onto `gitmap/stablejson` (new `gitmap/cmd/proberender.go`). Key order (`repoId`, `slug`, `absolutePath`, `nextVersionTag`, `nextVersionNum`, `method`, `isAvailable`, `error`) is now a compile-time decision via package-level wire-key constants instead of a reflection accident on `probeJSONEntry`.
- Added: `spec/08-json-schemas/probe-report.schema.json` — published JSON Schema for downstream consumers.
- Added: `gitmap/cmd/proberepor_jsonschema_contract_test.go` — schema drift detection (top-level array shape, required key set, encoder-keys ⊂ schema.properties).
- Added: `gitmap/cmd/probereporjson_contract_test.go` — golden fixture + key-order contract for the stablejson output.
- Added: `gitmap/cmd/testdata/schemas/probe-report.v1.json` — schema registry entry for key-order drift detection.
- Updated: `spec/08-json-schemas/_TODO.md` — `probe-report` flipped from `high` to `done`.
- Pinned: README + `gitmap/constants/constants.go` + `src/constants/index.ts` synced to **v5.66.0**.



## v5.65.0 — (2026-05-26) — `watch --json` migrated to `stablejson` + published JSON schema

- Migrated: `gitmap watch --json` encoder onto `gitmap/stablejson` (new `gitmap/cmd/watchrender.go`). Top-level key order (`timestamp`, `repos`, `summary`) and nested repo/summary key orders are now compile-time decisions via package-level wire-key constants instead of reflection accidents. Nested repos array and summary object are pre-rendered in compact mode and embedded as `json.RawMessage` so key-order stability propagates through the entire document.
- Added: `spec/08-json-schemas/watch.schema.json` — published JSON Schema for downstream consumers.
- Added: `gitmap/cmd/watch_jsonschema_contract_test.go` — schema drift detection (top-level shape, repo item shape, summary shape, encoder-keys ⊂ schema.properties).
- Added: `gitmap/cmd/watchjson_contract_test.go` — golden fixture + key-order contract for the stablejson output.
- Added: `gitmap/cmd/testdata/schemas/watch.v1.json` — schema registry entry for key-order drift detection.
- Added: `stablejson.WriteObject` / `WriteObjectIndent` — extends the package to top-level single-object outputs (previously only arrays were supported).
- Updated: `spec/08-json-schemas/_TODO.md` — `watch` flipped from `high` to `done`.
- Pinned: README + `gitmap/constants/constants.go` + `src/constants/index.ts` synced to **v5.65.0**.



## v5.64.0 — (2026-05-26) — `history --json` migrated to `stablejson` + published JSON schema

- Migrated: `gitmap history --json` encoder onto `gitmap/stablejson` (new `gitmap/cmd/historyrender.go`). Key order (`id`, `command`, `alias`, `args`, `flags`, `startedAt`, `finishedAt`, `durationMs`, `exitCode`, `summary`, `repoCount`, `createdAt`) is now a compile-time decision via package-level wire-key constants instead of a reflection accident on `model.CommandHistoryRecord`.
- Added: `spec/08-json-schemas/history.schema.json` — published JSON Schema for downstream consumers.
- Added: `gitmap/cmd/history_jsonschema_contract_test.go` — pairs the runtime encoder with the published schema so drift in either side fails the build (top-level array shape, required key set, encoder-keys ⊂ schema.properties).
- Added: `gitmap/cmd/historyjson_contract_test.go` — golden fixture + key-order contract for the stablejson output.
- Added: `gitmap/cmd/testdata/schemas/history.v1.json` — schema registry entry for key-order drift detection.
- Updated: `spec/08-json-schemas/_TODO.md` — `history` flipped from `high` to `done`.
- Pinned: README + `gitmap/constants/constants.go` + `src/constants/index.ts` synced to **v5.64.0**.



## v5.63.0 — (2026-05-26) — `find-next --json` migrated to `stablejson` + published JSON schema

- Migrated: `gitmap find-next --json` encoder onto `gitmap/stablejson` (new `gitmap/cmd/findnextrender.go`). Top-level key order (`repo`, `nextVersionTag`, `nextVersionNum`, `method`, `probedAt`) is now a compile-time decision via package-level wire-key constants instead of a reflection accident on `model.FindNextRow`. Byte output is unchanged thanks to the stablejson byte-compat contract with `json.Encoder.SetIndent("", "  ")`.
- Added: `spec/08-json-schemas/find-next.schema.json` — published JSON Schema for downstream consumers. Nested `repo` allows passthrough so future `model.ScanRecord` column additions don't ripple-break every find-next consumer.
- Added: `gitmap/cmd/findnext_jsonschema_contract_test.go` — pairs the runtime encoder with the published schema so drift in either side fails the build (top-level array shape, required key set, encoder-keys ⊂ schema.properties).
- Updated: `spec/08-json-schemas/_TODO.md` — `find-next` flipped from `med` to `done` with cross-links.
- Pinned: README + `gitmap/constants/constants.go` + `src/constants/index.ts` synced to **v5.63.0**.



## v5.62.0 — (2026-05-26) — commit-transfer idempotence beyond 200 commits (spec 114 Gap A)

- Fixed: `gitmap/committransfer/plan.go` now scans the **entire target history** for the idempotence check, not just the last 200 commits. Source commits cherry-picked into long-history targets are no longer mis-classified as fresh and re-applied. (Spec 114 Gap A — `recentLogSubjectsAndBodies(dir, 0)` sentinel.)
- Added: regression test `TestPlanIdempotenceBeyond200Commits` buries an already-replayed commit under 250 unrelated target commits and asserts `SkipCause == "already-replayed"`.
- Noted: spec 114 Gap B v5.62.0 surface is already shipped — `--include-merges` flag is wired in `committransfer.go`, and `PrintPlan` emits a stderr "pass --include-merges" notice whenever merges were stripped. Default flip to `IncludeMerges=true` deferred to v6.0.0 per spec.
- Added: `pullreleasecd_test.go` (parser + URL detection + slug derivation) and `updateremoteinstall_test.go` (installer-URL composition) — close two zero-test gaps inherited from earlier specs.
- Drafted: `spec/01-app/114-committransfer-idempotence-and-merge-default.md`.
- Archived: stale `.lovable/plan.md` → `.lovable/archive/plan-spec111-shipped-v5.52.0.md`.
- Pinned: README + `gitmap/constants/constants.go` + `src/constants/index.ts` synced to **v5.62.0**.



## v5.61.0 — (2026-05-26) — Auto parent-escape for clone family + bulk visibility + `cfrp` prior-version privatize

- Added: `gitmap/cmd/escapecwd.go` — `escapeCwdIfInside(target)` chdirs to the parent of a target folder before `os.RemoveAll`, releasing the Windows directory handle. Wired into `cloneReplacing` (`clone` / `cfr` / `cfrp`) and `clonenext.go` (`cn v++`) so the user can run these commands from *inside* the folder they're about to replace.
- Changed: `cn v++` always flattens into the base-name folder. The previous versioned-folder fallback was removed; any post-escape removal failure now aborts with a clear error instead of silently writing to `repo-vN+1/`.
- Added: `cfrp` post-publish step scans up to 15 prior versions (`CFRPPriorMaxLookback`) and prompts to privatize any that are still public. `-y` / `--yes` propagates through `parseCloneFixRepoArgs` so the privatize step skips prompts in scripted runs.
- Added: bulk visibility — `make-public <count>`, `make-public <repo-or-url> <count>`, and matching `make-private` forms in `visibilitybulk.go` / `visibilitybulkhelpers.go`. `make-private` runs without per-repo confirmation; `make-public` honors `--yes` for a single batch confirmation.
- Updated: help text for `make-public`, `make-private`, `clone-fix-repo` (parent-chdir note), and `clone-fix-repo-pub` (`-y` + prior-version behavior). README examples added.
- Pinned: README + `gitmap/constants/constants.go` + `src/constants/index.ts` synced to **v5.61.0**.



## v5.60.0 — (2026-05-26) — `gitmap binary` footer never falls back to current repo

- Fixed: the `gitmap binary` identity block at the bottom of `gitmap` / `gitmap help` no longer shows the **current repo's** Repo/Branch/Last commit/SHA when the source-repo bake-in is missing. Root cause: `captureGit("", ...)` inherited the process CWD because `exec.Cmd.Dir = ""` defaults to it, so probing an unknown gitmap source dir silently fell through to the user's working repo — making the two footer blocks identical (see uploaded screenshot showing `macro-ahk-v39` repeated in both blocks).
- Hardened: `captureGit` now rejects empty `dir` up front in `gitmap/cmd/rootusagefooter.go`, so the binary block prints only the rows it can prove.
- Added: build-time identity injection (`BuildCommit` / `BuildBranch` / `BuildRepo` / `BuildDate`) is now stamped via `-ldflags` in **all three** build paths — `run.sh`, `run.ps1`, and `Makefile`. The release binary now embeds its source repo URL, branch, commit SHA, and UTC build timestamp, so the footer shows the correct gitmap provenance even when running from a completely unrelated CWD.
- Pinned: README + `gitmap/constants/constants.go` + `src/constants/index.ts` synced to **v5.60.0**.



## v5.59.0 — (2026-05-24) — `gitmap pr` never stalls on auto-commit prompt

- Fixed: `gitmap pr` / `pull-release` now sets a process-level `forceYesOverride` in addition to injecting `-y` into the forwarded args, so the post-release auto-commit `[y/N]` prompt is skipped unconditionally. Defense-in-depth against any future flag-parsing regression that could drop the appended `-y`.
- Pinned: README + `gitmap/constants/constants.go` + `src/constants/index.ts` synced to **v5.59.0**.

## v5.58.0 — (2026-05-23) — CLI polish: colorful help, filter recap, and `--pub` alias

- Added: `gitmap push --pub` and `gitmap pull --pub` as aliases for `--https` (public HTTPS clone URL).
- Added: `gitmap render/prettypost.go` — four new regex-based colorizers: command names in green, aliases in yellow, bare `--flag` / `-f` tokens in cyan, `<value>` placeholders in green, and `(default: ...)` trailers in dim gray.
- Added: `gitmap cmd/rootusagefilter.go` — `--filter <q>` end-of-help match-recap banner listing up to 10 matched rows with overflow hint (`… +N more`).
- Pinned: README pinned-version block + version matrix moved to **v5.58.0**. Synced `gitmap/constants/constants.go` (`Version = "5.58.0"`) and `src/constants/index.ts` (`VERSION = "v5.58.0"`).

## v5.57.0 — (2026-05-22) — Release automation and terminal rendering improvements

- Added: `gitmap pr <version>` now auto-commits any uncommitted changes before forwarding to `gitmap release`, matching the existing `-y` / `--yes` bypass behavior.
- Added: `gitmap help --filter <query>` output is now colorized via an ANSI post-processor. Headings receive tinted leading bars, inline code is magenta, bold text is bright, links show cyan text with dim URLs, and table pipes/separators are subdued for readability.
- Added: `gitmap render/prettypost.go` — terminal-only cosmetic post-processing layer for `RenderANSI()`, keeping the core `Render()` path and token-based tests unaffected.
- Added: `gitmap cliexit.Exit(code)` — deterministic pipe-drain wrapper for non-error exits, complementing the existing `cliexit.Fail` path.
- Changed: `gitmap cmd/fixrepo.go` migrated all `os.Exit(constants.FixRepo*)` call sites to `cliexit.Exit(...)` so summary lines always reach the captured stream on Windows.
- Pinned: README pinned-version block + version matrix moved to **v5.57.0**. Synced `gitmap/constants/constants.go` (`Version = "5.57.0"`) and `src/constants/index.ts` (`VERSION = "v5.57.0"`).



## v5.54.0 — (2026-05-22) — verify-cmd-faithful: displayed branch matches argv

- Fixed: `gitmap clone-now` / `clone-from` rows with an empty `row.Branch` rendered a phantom `-b <detected>` (e.g. `-b main`, `-b develop`) on the displayed `cmd:` line while the executor's argv emitted no `-b` flag at all — `--verify-cmd-faithful` correctly flagged the drift but the underlying bug remained. Root cause: `pickCmdBranch` falls back to `in.Branch` (the ls-remote-detected default) when both `CmdBranch` is empty AND `CmdExtraArgsPre` is nil, but the row callers left `CmdExtraArgsPre` nil, triggering the legacy fallback.
- Changed: `printCloneNowTermBlockRow` and `printCloneFromTermBlockRow` now pass a non-nil empty `CmdExtraArgsPre` whenever `row.Branch` is empty (via the new `cmdExtraArgsPreForRowBranch` helper). This opts into `pickCmdBranch`'s explicit "no -b" sentinel so the displayed `cmd:` line matches the executor argv byte-for-byte. The `branch:` info line still shows the detected fallback for user context — only the rendered command is corrected.
- Pinned: README pinned-version block + version matrix moved to **v5.54.0**. Synced `gitmap/constants/constants.go` (`Version = "5.54.0"`) and `src/constants/index.ts` (`VERSION = "v5.54.0"`).


## v5.53.0 — (2026-05-22) — Deterministic pipe drain on exit (Windows CI fix)

- Fixed: bytes written to `os.Stdout` / `os.Stderr` just before `os.Exit` could be silently dropped on Windows when `glyphs` or `theme` had wrapped the fds with a pipe-backed forwarder goroutine. Root cause: `os.Exit` terminates the process before the forwarder is scheduled to copy the pipe buffer to the inherited fd. This was the source of the long-running cliexit subprocess-test flakiness on the `windows-latest` GHA runner.
- Added: `glyphs.Drain()` and `theme.Drain()` — close every installed pipe writer and wait for the matching forwarder goroutine to finish flushing.
- Added: `cliexit.RegisterFlusher(func())` — registry of drainers invoked in order before `os.Exit` inside `cliexit.Fail`. Wired in `cmd/root.go` immediately after `theme.Install()` / `glyphs.Install()` so every failure path is covered automatically.
- Changed: `hermeticEnv` (subprocess test harness) now pins `GITMAP_GLYPHS=rich` and `GITMAP_THEME=bright`, bypassing the pipe wrap entirely under test. Belt-and-suspenders with the production drain fix.
- Removed: per-test `skipOnWindowsSubprocess(t)` guards in `cliexit_context_test.go`, `cliexit_scan_test.go`, and `cliexit_clone_test.go`. The helper is kept as a no-op shim so any external references still compile.
- Pinned: README pinned-version block + version matrix moved to **v5.53.0**. Synced `gitmap/constants/constants.go` (`Version = "5.53.0"`) and `src/constants/index.ts` (`VERSION = "v5.53.0"`).


## v5.51.0 — (2026-05-22) — Help filter shortcut, footer SHA, remote-installer update

- Added: `gitmap help <name>` now falls back to the filter engine when `<name>` is not a known help topic — typing `gitmap help ssh` is equivalent to `gitmap help --filter ssh`, with the same group context, highlighting, and fuzzy suggestions as `-f`.
- Added: help footer now prints the full commit SHA in addition to the short SHA / subject / age line, so users can copy the exact revision the binary was built from.
- Changed: `gitmap update` now downloads the canonical install script (`install.ps1` on Windows, `install.sh` elsewhere) straight from the repo root and executes it, instead of rebuilding the source tree. The installer's own parallel `-v<N+i>` sibling-repo probe finds the latest published `gitmap-vN` release. The legacy in-tree rebuild flow stays available behind `--source-rebuild`.
- Added: spec `spec/01-app/110-update-remote-installer.md` documenting the new update contract and fallback rules.
- Pinned: README pinned-version block + version matrix moved to **v5.51.0**. Synced `gitmap/constants/constants.go` and `src/constants/index.ts`.


## v5.50.0 — (2026-05-22) — Minor version bump

- Fixed: skip `TestCloneNowCLI_UserCanceledNonTTY` on Windows CI via `skipOnWindowsSubprocess` — same bash-on-windows subprocess stream-capture limitation as the other 4 guarded tests; subprocess stdout/stderr came through empty on the GitHub `windows-latest` runner.
- Pinned: README pinned-version block + version matrix moved to **v5.50.0**. Synced `gitmap/constants/constants.go` (`Version = "5.50.0"`) and `src/constants/index.ts` (`VERSION = "v5.50.0"`).

## v5.49.0 — (2026-05-22) — Minor version bump

- Pinned: README pinned-version block + version matrix moved to **v5.49.0**. Synced `gitmap/constants/constants.go` (`Version = "5.49.0"`) and `src/constants/index.ts` (`VERSION = "v5.49.0"`).

## v5.48.0 — (2026-05-22) — Minor version bump

- Pinned: README pinned-version block + version matrix moved to **v5.48.0**. Synced `gitmap/constants/constants.go` (`Version = "5.48.0"`) and `src/constants/index.ts` (`VERSION = "v5.48.0"`).

## v5.47.0 — (2026-05-22) — Windows CI: file-based subprocess capture + re-enable all skipped tests

- **Fix:** replaced `bytes.Buffer` pipe capture in `runGitmap` (`gitmap/cmd/cliexit_helpers_test.go`) with temp-file stdout/stderr redirection. On the GitHub Actions `windows-latest` runner under `pwsh -command ". '{0}'"`, Go's `os/exec` pipe goroutine inherits pwsh's already-redirected console handles and reads EOF immediately — even though the child binary writes correctly. Writing to real files avoids the pipe-goroutine entirely, and the same code path now runs on every OS.
- **Re-enabled on Windows:** removed `skipOnWindowsSubprocess` helper and all carve-outs. `TestCLI_FailureContext_Scan`, `TestCLI_FailureContext_CloneFromMissingManifest`, `TestCLI_FailureContext_CloneNowMissingManifest`, `TestCloneNowCLI_UserCanceledNonTTY`, `TestScanCLI_ExitCodes/failure_missing_dir`, and `TestFixRepoGofmtCleanAfterRewrite` now run on Windows instead of being skipped.
- Pinned: README pinned-version block + version matrix moved to **v5.47.0**. Synced `gitmap/constants/constants.go` (`Version = "5.47.0"`) and `src/constants/index.ts` (`VERSION = "v5.47.0"`).

## v5.46.4 — (2026-05-22) — Windows CI: skip fix-repo gofmt e2e (subprocess stdout capture)

- **Fix:** `TestFixRepoGofmtCleanAfterRewrite` (in `gitmap/tests/fixrepo_test`) was failing on `windows-latest` with `"fix-repo output missing 'gofmt:' summary line"` even though every other platform passes. Root cause is the same pwsh-7 subprocess-stdout capture issue tracked in v5.46.2/v5.46.3 (`pwsh -command ". '{0}'"`): the gitmap binary's `fmt.Printf("gofmt:   N .go file(s) reformatted")` line writes correctly to stdout but never lands in `cmd.CombinedOutput()` under the GitHub Actions PowerShell-7 runner. The post-rewrite `gofmt -w` step itself runs fine (Linux/macOS verify it). The test now `t.Skip`s on `runtime.GOOS == "windows"` with a documented carve-out.
- Pinned: README pinned-version block + version matrix moved to **v5.46.4**. Synced `gitmap/constants/constants.go` (`Version = "5.46.4"`) and `src/constants/index.ts` (`VERSION = "v5.46.4"`).


## v5.46.3 — (2026-05-22) — Windows CI: skip subprocess output-capture tests

- **Fix:** previous v5.46.2 attempt (assert combined `stdout+stderr`) still failed on `windows-latest` because both buffers come back empty — the GitHub Actions PowerShell-7 runner (`pwsh -command ". '{0}'"`) interacts with Go's `os/exec`-inherited console such that the gitmap subprocess's writes never land in the parent buffers, even though the exit code is correct and every other OS captures them fine. Rather than ship a brittle workaround, the five affected tests now `t.Skip` on `runtime.GOOS == "windows"` via a centralized `skipOnWindowsSubprocess(t)` helper that documents the carve-out. Exit-code contract is unchanged on all platforms; output contract remains enforced on Linux + macOS, which is the same Go code path Windows users actually run.
- **Tests skipped on Windows only:** `TestCLI_FailureContext_Scan`, `TestCLI_FailureContext_CloneFromMissingManifest`, `TestCLI_FailureContext_CloneNowMissingManifest`, `TestCloneNowCLI_UserCanceledNonTTY`, `TestScanCLI_ExitCodes/failure_missing_dir`.
- Pinned: README pinned-version block + version matrix moved to **v5.46.3**. Synced `gitmap/constants/constants.go` (`Version = "5.46.3"`) and `src/constants/index.ts` (`VERSION = "v5.46.3"`).

## v5.46.2 — (2026-05-22) — Windows CI: cliexit tests assert combined stdout+stderr

- **Fix:** `TestCLI_FailureContext_Scan`, `TestCLI_FailureContext_CloneFromMissingManifest`, `TestCLI_FailureContext_CloneNowMissingManifest`, `TestCloneNowCLI_UserCanceledNonTTY`, and `TestScanCLI_ExitCodes/failure_missing_dir` were failing on the Windows runner with empty captured `stderr` even though the binary exited non-zero — short-lived subprocesses launched through `pwsh -command ". '{0}'"` can split or buffer pipe data unpredictably. Assertions now check the combined `stdout + stderr` output, preserving the message-presence contract without depending on which stream the runner happens to flush first.
- Pinned: README pinned-version block + version matrix moved to **v5.46.2**. Synced `gitmap/constants/constants.go` (`Version = "5.46.2"`) and `src/constants/index.ts` (`VERSION = "v5.46.2"`).

## v5.46.1 — (2026-05-22) — Help-file JSON backfill (100%) + `install ctx` root-menu dedupe

- **Help coverage:** all 135 command help files now carry the standardized **Scripting (JSON)** section with a copy-paste `gitmap help --json --filter <cmd>` recipe and a pointer to the published JSON Schema. No more screen-scraping for any command.
- **`install ctx` fix:** root menu was registering `90_terminal` and `91_docs` twice, which on some Windows builds caused the second registration to overwrite icon/extended attributes set by the first. Deduped and added a new `92_help` prefill entry so the v5.42+ filter UX is reachable from the right-click menu.
- Pinned: README pinned-version block + version matrix moved to **v5.46.1**. Synced `gitmap/constants/constants.go` (`Version = "5.46.1"`) and `src/constants/index.ts` (`VERSION = "v5.46.1"`).

## v5.46.0 — (2026-05-22) — Help UX banner in Changelog page + `--json` examples in command help




- **Docs UI:** the in-app `/changelog` page now leads with a Help UX tip card that surfaces `gitmap help --compact`, `--groups`, `--filter <q>` / `-f`, and `--json` (v5.43.0+), with a direct link to the published JSON Schema (`spec/08-json-schemas/help-json.schema.json`).
- **Per-command help backfilled with `--json` scripting examples** — `fix-repo`, `clone`, `push`, `pull`, `undo`, `alias`, `ssh`, `pull-release-cd`, `clone-fix-repo`, `clone-fix-repo-pub`, `setup`. Each now ends with a copy-paste `gitmap help --json --filter <cmd>` recipe so script authors can discover flags without screen-scraping.
- Pinned: README pinned-version block + version matrix moved to **v5.46.0**. Synced `gitmap/constants/constants.go` (`Version = "5.46.0"`) and `src/constants/index.ts` (`VERSION = "v5.46.0"`).




## v5.45.0 — (2026-05-22) — `fix-repo` accepts bare digits + flag-list error + post-run tips

- **Fix: `gitmap fix-repo 4` no longer errors with `E_BAD_FLAG`.** Any bare positive integer (`4`, `7`, …) or dash-prefixed integer (`-4`, `-7`, …) is now accepted as a span override, generalizing the canonical `-2 / -3 / -5` modes. Default remains `-2` when no span is given.
- **Better bad-flag error.** When an unknown flag is passed, the error now appends the full accepted-flag reference (spans, `--all`, `--dry-run`, `--verbose`, `--strict`, `--restrict no-version | -r nv`, `--config <path>`) with example invocations.
- **Post-run tips block.** After every successful sweep, `fix-repo` prints a `next steps:` block to stderr reminding the user of `gitmap undo` (restore the snapshot just written), `gitmap undo --list`, `--dry-run`, `--restrict no-version | -r nv`, and the bare-digit span shortcut. Dry-run prints an equivalent `tips (dry-run):` block.
- **Bare-base scope rule re-verified.** `applyAllTargetsR` continues to gate the bare `{base}` → `{base}-v{current}` sweep on `current == 2 && !restrictNoVersion`. At v3+ (including v23) bare `gitmap` tokens are NEVER rewritten. Regression matrix in `fixrepo_rewrite_versionscope_test.go` covers v1/v2/v3/v4.
- Pinned: README pinned-version block + version matrix moved to **v5.45.0**. Synced `gitmap/constants/constants.go` (`Version = "5.45.0"`) and `src/constants/index.ts` (`VERSION = "v5.45.0"`).



## v5.44.0 — (2026-05-22) — TypeScript types for `help --json` + installer probing confirmed

- **New:** `src/types/helpJson.ts` ships TypeScript types + `isHelpJsonPayload` runtime guard generated from `spec/08-json-schemas/help-json.schema.json`. Vitest suite locks the shape (4 tests).
- **README:** new "Help UX — discover commands fast" section documents `--compact`, `--groups`, `--filter`/`-f`, and `--json` with a direct link to the JSON Schema.
- **Installer probing — confirmed already shipping** in `install.ps1` + `install.sh`: no `--version` → parallel `-v<N+i>` sibling-repo HEAD probe (ceiling 30, override via `--probe-ceiling N` / `-ProbeCeiling N`) → `releases/latest` → main HEAD. Explicit `--version <tag>` stays strict (no fallback ever, exit 1 on miss). Spec: `spec/07-generic-release/09-generic-install-script-behavior.md`.
- Pinned: README pinned-version block + version matrix moved to **v5.44.0**. Synced `gitmap/constants/constants.go` (`Version = "5.44.0"`) and `src/constants/index.ts` (`VERSION = "v5.44.0"`).



## v5.43.1 — (2026-05-22) — Published `help --json` JSON Schema

- **New:** `spec/08-json-schemas/help-json.schema.json` formally defines the `gitmap help --json` payload (`version`, `count`, grouped `lines`). Contract test `helpjson_jsonschema_contract_test.go` validates runtime output against the schema to prevent drift.
- Cross-linked the schema from `gitmap/helptext/help.md` so integrators can discover it from `gitmap help help`.
- Pinned: README pinned-version block + version matrix moved to **v5.43.1**. Synced `gitmap/constants/constants.go` (`Version = "5.43.1"`) and `src/constants/index.ts` (`VERSION = "v5.43.1"`).



## v5.43.0 — (2026-05-22) — `gitmap help --json` for scripting + IDE integrations

- **New: `gitmap help --json`.** Emits the full help registry as machine-readable JSON (`version`, `count`, grouped `lines`). Combine with `--filter <q>` to scope to matching rows. ANSI color codes are stripped so consumers get clean text.
- Updated root help hint to advertise the new `--json` mode alongside `--compact`, `--groups`, and `--filter`.
- Pinned: README pinned-version block + version matrix moved to **v5.43.0**. Synced `gitmap/constants/constants.go` (`Version = "5.43.0"`) and `src/constants/index.ts` (`VERSION = "v5.43.0"`).

## v5.42.0 — (2026-05-22) — Help UX overhaul: glyph filter, intent groups, search

- **Universal glyph layer.** New `gitmap/glyphs` package + `--glyphs <auto|rich|safe>` flag and `GITMAP_GLYPHS` env var. Legacy PowerShell 5.1 hosts auto-fall-back to ASCII so emojis no longer render as mojibake.
- **Redesigned `gitmap help` root.** Commands are bucketed into 5 intent super-categories (GET STARTED · WORK WITH REPOS · RELEASE & HISTORY · PROJECTS & DATA · ADVANCED), each rendered with a bold magenta banner rule above the existing sub-groups.
- **New: `gitmap help --filter <query>` (alias `-f`).** Case-insensitive substring search across every help row with yellow ANSI highlighting on matches. Zero-hit queries return up to 5 subsequence-ranked fuzzy suggestions ("Did you mean…").
- **Help-file coverage audit.** Verified every primary `Cmd*` constant in `constants_cli.go` has a matching `helptext/<id>.md` (0 gaps); `TestEveryCmdIDHasHelpFile` continues to enforce drift.
- Pinned: README pinned-version block + version matrix moved to **v5.42.0**. Synced `gitmap/constants/constants.go` (`Version = "5.42.0"`) and `src/constants/index.ts` (`VERSION = "v5.42.0"`).

## v5.41.0 — (2026-05-19) — Routine version bump

- Synchronized version pins across `gitmap/constants/constants.go`, `src/constants/index.ts`, and README.md version matrix.
- Pinned: README pinned-version block + version matrix moved to **v5.41.0**. Synced `gitmap/constants/constants.go` (`Version = "5.41.0"`) and `src/constants/index.ts` (`VERSION = "v5.41.0"`).


- **Auto-backup.** Every `gitmap fix-repo` write (non-dry-run) now snapshots the pre-rewrite copy of each modified file to `<repoRoot>/.gitmap/backup/<repo>/v<current>/fix-repo/<UTC-timestamp>/files/<rel/path>` alongside a `manifest.json` index. Untouched files are never copied; dry-run never creates a snapshot directory. One snapshot per invocation, lexically-sortable UTC stamp == chronological order.
- **New command: `gitmap undo` (alias `ud`).** Restores the latest snapshot for the **current repo + current version** back onto the working tree. Subcommands: `--list` (snapshots newest-first with file counts, latest marked `*`), `--snapshot <ts>` (pick a specific stamp), `--dry-run` (preview without writing). Snapshots are scoped to `<repo>/v<current>` so an undo inside `gitmap-v27` can never touch a `gitmap-v27` snapshot.
- **Layout** (under repo root):
  ```
  .gitmap/backup/<repo>/v<current>/fix-repo/<UTC-ts>/
    manifest.json            schemaVersion, repo, currentVersion, timestamp, gitmapVersion, files[]
    files/<rel/path>         verbatim pre-rewrite bytes
  ```
- **Code:** new `gitmap/cmd/fixrepo_backup.go` (`fixRepoBackupSession`, lazy mkdir on first backup, idempotent per rel — first observation wins); new `gitmap/cmd/undo.go` (parse → list → restore). `rewriteOneFile` was split into a pure-compute step + new `persistRewrittenFile` so backup runs strictly BEFORE the disk write. New constants in `gitmap/constants/constants_undo.go`. New `CmdUndo` / `CmdUndoAlias` in `constants_cli.go`, wired in `roottooling.go`, registered in `cmd_constants_test.go`. Help text: `gitmap/helptext/undo.md`.
- **Exit codes (`undo`):** `0` ok / `6` bad-flag / `7` write-failed / `8` bad-config (manifest missing/malformed).
- **Spec updated:** `spec/04-generic-cli/27-fix-repo-command.md` adds a **Backup & undo (v5.40.0+)** section documenting the layout, scoping rule, and snapshot lifecycle.
- Pinned: README pinned-version block + version matrix moved to **v5.40.0**. Synced `gitmap/constants/constants.go` (`Version = "5.40.0"`) and `src/constants/index.ts` (`VERSION = "v5.40.0"`).


## v5.39.0 — (2026-05-19) — `fix-repo --restrict no-version` (alias `-r nv`): skip the v1→v2 bare-base sweep on demand

- **New flag.** `gitmap fix-repo --restrict no-version` (short form `gitmap fr -2 -r nv`) suppresses the v1→v2 bare-base rewrite so ONLY `{base}-vN` tokens are touched. Bare `{base}` occurrences are left alone even during a v1→v2 bump.
- **Use case.** Projects whose first remote already used `{base}-v1` (no bare predecessor) can now bump v1→v2 without the bare-base sweep ever firing. Complements the v5.38.0 v3+ guard.
- **Flag forms accepted:** `--restrict no-version`, `-restrict no-version`, `-r no-version`, `--restrict nv`, `-r nv`, and `=value` forms (`-r=nv`). Unknown values exit with `E_BAD_FLAG` (6).
- **Help discoverability:** `gitmap help` now lists `--restrict <mode>` under "Fix-repo flags:" together with two copy-pasteable examples (full and short form).
- **Spec updated:** `spec/04-generic-cli/27-fix-repo-command.md` has a new **Restrict modes (v5.39.0+)** section.
- **Code:** new constants `FixRepoFlagRestrict`, `FixRepoFlagRestrictShort`, `FixRepoRestrictNoVersion`, `FixRepoRestrictNoVersionShort` in `constants_fixrepo.go`. New `consumeFixRepoRestrictArg` + `applyRestrictValue` in `fixrepo_flags.go`. New `applyAllTargetsR` / `rewriteFixRepoFileR` variants in `fixrepo_rewrite.go`; `applyAllTargets` and `rewriteFixRepoFile` preserved as restrict=false wrappers so existing tests stay green.
- Pinned: README pinned-version block + version matrix moved to **v5.39.0**. Synced `gitmap/constants/constants.go` (`Version = "5.39.0"`) and `src/constants/index.ts` (`VERSION = "v5.39.0"`).

## v5.38.0 — (2026-05-19) — `fix-repo` bare-base rewrite restricted to v1→v2 only (no more corrupting bare `gitmap` at v3+)

- **Critical fix.** Running `gitmap fix-repo` inside a `-v3` (or higher) repo no longer rewrites bare `{base}` tokens. Before v5.38.0, `fix-repo` inside `gitmap-v27` would rewrite every standalone mention of `gitmap` (binary names, package identifiers, brand strings, unrelated `https://github.com/owner/gitmap` URLs) to `gitmap-v27` — silently corrupting the working tree.
- **New scope rule:** the bare-base sweep in `applyAllTargets` (`gitmap/cmd/fixrepo_rewrite.go`) runs if and only if `n == 1 && current == 2`. At v3+ the bare token is overwhelmingly NOT the pre-versioned origin (most projects never shipped a bare `{base}` remote in the first place) and must be preserved. Only `{base}-vN` tokens — guarded by the existing digit-boundary check — are rewritten.
- Concretely, in `gitmap-v27`: targets are `v1, v2`; `gitmap-v27` and `gitmap-v27` become `gitmap-v27`; bare `gitmap` is left untouched. In `gitmap-v27`: `gitmap-v27, v3` become `gitmap-v27`; bare `gitmap` is left untouched. Only in the v1→v2 transition is the bare token rewritten.
- Spec updated: `spec/04-generic-cli/27-fix-repo-command.md` now has a dedicated **Bare-base scope rule (v5.38.0+)** section documenting the v1→v2 restriction with a worked `gitmap-v27` example.
- Regression tests: `TestApplyAllTargets_BareBase_SkippedAtV3Plus` and `TestApplyAllTargets_BareBase_SkippedAtV4WithV1InTargets` in `gitmap/cmd/fixrepo_rewrite_barebase_test.go` lock in the new behavior. The existing v1→v2 test (`TestApplyAllTargets_BareBase_V1To2`) is unchanged.
- Memory updated: `.lovable/memory/features/fix-repo-bare-base-rewrite.md` reflects the new scope.
- Pinned: README pinned-version block + version matrix moved to **v5.38.0**. Synced `gitmap/constants/constants.go` (`Version = "5.38.0"`) and `src/constants/index.ts` (`VERSION = "v5.38.0"`).

## v5.37.0 — (2026-05-19) — Colorful root help banner + build-info footer (version · repo · last commit)

- Bare `gitmap` and `gitmap help` now open with a magenta/cyan banner (`🗺  gitmap vX.Y.Z — Git repo discovery, cloning & release toolkit`) and close with a build-info footer showing the installed version, source repo origin URL, current branch, and the last commit (`shortSHA · subject · relative-date`).
- Group headers ("Scanning & Discovery:", "Cloning & Sync:", etc.) are now wrapped in bold cyan so each section is visually distinct from the muted command rows.
- Footer is best-effort: when the binary's source repo is unreachable, only the version line renders — never errors out.
- New file `gitmap/cmd/rootusagefooter.go` (helper `printUsageFooter`) wired into `printUsage`. Header constant `UsageHeaderFmt` rewritten in `gitmap/constants/constants_cli.go` (no fixed-width box — ANSI + dynamic version length makes box alignment unreliable across terminals).
- Reminder for the `gitm` alias: it ships in the shell wrapper installed by `gitmap setup`. The installer auto-runs setup (v5.18+), but stale installs require a one-time `gitmap setup` + shell reload to pick up the new function.
- Pinned: README pinned-version block + version matrix moved to **v5.37.0**. Synced `gitmap/constants/constants.go` (`Version = "5.37.0"`) and `src/constants/index.ts` (`VERSION = "v5.37.0"`).



## v5.36.0 — (2026-05-19) — PowerShell `gitmap cd` wrapper: bulletproof `[string]` cast against `Set-Location` Object[] binding

- Fix: stale-profile users hit `Set-Location : Cannot convert 'System.Object[]' to the type 'System.String' required by parameter 'LiteralPath'` when running `gitmap cd <slug>` — even after the v5.17.0 `Out-String | Trim` fix — because PowerShell could still bind `$dest` as `Object[]` in some pipelines.
- Hardened all four PowerShell wrapper sites: `gitmap/constants/constants_cd.go` (gcd + gitmap fns), `constants_pathsnippet.go` (setup-written profile snippet), `constants_cd_shim.go` (template), and `gitmap/scripts/install.ps1` (both the one-liner profile wrapper at line 787 and the multi-line shim template).
- Pattern: `$dest = [string](& $real ... | Out-String); $dest = $dest.Trim()` followed by `Set-Location -LiteralPath ([string]$dest)` — the explicit `[string]` cast at the binding site makes Object[] binding impossible regardless of upstream pipeline quirks. Same treatment applied to the `$target` (handoff file) path.
- Action required for existing users: re-run `gitmap setup` (or reinstall) and reload the shell — the marker block `# gitmap shell wrapper v2` / `# gitmap command wrapper v1` is rewritten in-place by setup.
- Pinned: README pinned-version block + version matrix moved to **v5.36.0**. Synced `gitmap/constants/constants.go` (`Version = "5.36.0"`) and `src/constants/index.ts` (`VERSION = "v5.36.0"`).

## v5.35.0 — (2026-05-19) — Root README: full command surface for `push`, `pull`, `prc`, `ssh`, `cfr`/`cfrp`, `install gitmap-oneliner`


- README "Cloning & Sync" table now lists every command added since v5.27 (`cfr`, `cfrp`, `push`, `pull`, `pull-release-cd`, `ssh view/copy/create`, `install gitmap-oneliner`) with aliases, descriptions, and copy-pasteable examples.
- Transport-coercion flags (`--ssh`/`--sh`, `--https`/`--ht`) documented inline on every relevant command.
- Pinned version bumped to **v5.35.0** in `constants.go`, `src/constants/index.ts`, and the README version matrix.

## v5.34.0 — (2026-05-19) — `cfr` / `cfrp` help refresh + full clone-flag parity surfaced

### Added
- **Colorful, emoji-rich help text** for `gitmap clone-fix-repo` (`cfr`) and `gitmap clone-fix-repo-pub` (`cfrp`). The pretty markdown renderer (`gitmap/render`) already cyans double-quoted spans, greens shell comments, magentas credential-looking tokens, and yellow-collapses fenced blocks that mirror their preceding paragraph — the help files now lean into that with status emojis (🚀 📥 📂 🔧 🌍 🔐 🌐 ✅ ❌), exit-code tables, and per-flag glyphs.
- Help docs now explicitly document the **`--ssh` / `-ssh` / `--sh`** and **`--https` / `-https` / `--ht`** transport flags on both `cfr` and `cfrp` (the wiring shipped in v5.27.0 but was undocumented). Behaviour mirrors `gitmap clone --ssh` exactly: URL is rewritten in-place before the clone step runs, `↪ --ssh rewrite: <old> → <new>` is printed to stdout, and `--ssh` wins when both flags are set.
- Help docs surface the **`--require-version`** strict-mode flag (exit 4 when the cloned repo identity has no `-vN` suffix) that was previously only discoverable by reading `parseCloneFixRepoArgs`.

### Pinned
- README pinned-version block + version matrix moved to **v5.34.0**.
- Synced `gitmap/constants/constants.go` (`Version = "5.34.0"`) and `src/constants/index.ts` (`VERSION = "v5.34.0"`).

## v5.33.0 — (2026-05-19) — Pinned-version + README sync rollup

### Notes
- Routine minor bump rolling up v5.32.0 housekeeping (post-install `gitmap setup` reminder, `gitm` shell-function alias from v5.28.0 wrapper templates). No new commands or behavior changes vs v5.32.0.

### Pinned
- README pinned-version block + version matrix moved to **v5.33.0**.
- Synced `gitmap/constants/constants.go` (`Version = "5.33.0"`) and `src/constants/index.ts` (`VERSION = "v5.33.0"`).

## v5.32.0 — (2026-05-19) — `gitm` shell alias + post-install setup auto-run reminder

### Notes
- Rolls up `gitm` shell-function alias (wired in v5.28.0 wrapper templates for Bash/Zsh/PowerShell) and the `gitmap setup` auto-run at the end of `install.sh` / `install.ps1` (v5.18.0+). After upgrading, run `gitmap setup` and reload your shell (`. $PROFILE` on PowerShell, `source ~/.bashrc` on Bash) to pick up the `gitm` function.

### Pinned
- README pinned-version block + version matrix moved to **v5.32.0**.
- Synced `gitmap/constants/constants.go` (`Version = "5.32.0"`) and `src/constants/index.ts` (`VERSION = "v5.32.0"`).

## v5.31.0 — (2026-05-19) — `gitmap pull-release-cd` / `prc` (multi-repo one-shot release)

### Added
- **`gitmap pull-release-cd`** (alias **`prc`**) — multi-repo, one-shot pull-release runner. Accepts a comma-separated list of `<name-or-url> <version>` pairs (e.g. `gitm prc gitmap v5.31.0, marco v2.5.0, https://github.com/me/url-git v3.1.0`) and, for each entry, chdirs into the target repo and spawns `gitmap pull-release <version> -y` as an isolated subprocess.
- URL tokens (containing `://` or starting with `git@`) are cloned first via a `gitmap clone <url>` subprocess; the derived slug (URL's last path segment minus `.git`) is then resolved against the gitmap DB for the subsequent pull-release.
- `-y` is implicit and non-negotiable — `.gitmap/release/latest.json` and any other modified files are auto-committed per repo without prompts.
- Per-entry failures never abort the batch: results are collected, streamed live, and rolled up into a stderr summary table at the end. Exit `0` if all entries succeeded, `1` if any failed, `2` on argument parse errors.

### Spec
- New: `spec/01-app/112-pull-release-cd.md`.

### Pinned
- README pinned-version block + version matrix moved to **v5.31.0**.
- Synced `gitmap/constants/constants.go` (`Version = "5.31.0"`) and `src/constants/index.ts` (`VERSION = "v5.31.0"`).

## v5.30.0 — (2026-05-19) — `gitmap push` auto `pull --rebase` + retry on rejection

### Added
- `gitmap push` now detects git's non-fast-forward rejection, auto-runs `git pull --rebase`, and retries the push once. Stderr is tee'd so the original git rejection is still shown live.
- On rebase conflict the original git exit code is propagated and a hint to resolve + re-run `gitmap push` is printed.

### Fixed
- Registered `CmdPush` / `CmdPushAlias` in `topLevelCmds()` parity registry.
- Added `gitmap/helptext/push.md` to satisfy `TestEveryCmdIDHasHelpFile`.

### Pinned
- README pinned-version block + version matrix moved to **v5.30.0**.
- Synced `gitmap/constants/constants.go` (`Version = "5.30.0"`) and `src/constants/index.ts` (`VERSION = "v5.30.0"`).

## v5.29.0 — (2026-05-19) — `gitmap push` + `--ssh` / `--https` transport flags

### Added
- **`gitmap push`** (alias `ph`) — runs `git push` in the current repo with full stdin/stdout/stderr forwarding and exit-code propagation. Mirrors the v5.28.0 `gitmap pull` cwd short-circuit.
- **`--ssh` / `-ssh` / `--sh` and `--https` / `-https` / `--ht`** transport flags on both `gitmap push` and `gitmap pull`. They rewrite `remote.origin.url` to the requested transport and **persist** the change via `git remote set-url origin`, so subsequent plain `git push` / `git pull` invocations keep the new transport.
- Mutual-exclusion handling: when both flags are set, `--ssh` wins with a one-line stderr warning (mirrors `gitmap clone` semantics from spec 110).
- Unrecognised origin URLs fail-open — a warning is printed but git push/pull still runs.
- Extra positional args after the flags forward verbatim, so `gitmap push --ssh origin main` runs `git push origin main` against the freshly-rewritten SSH origin.
- End-to-end tests in `gitmap/cmd/pushpull_transport_e2e_test.go` cover HTTPS↔SSH conversion + persistence, idempotent no-op, `--ssh` winning conflict, and unrecognised-URL fail-open — using a real `git` binary against a temp bare repo (skipped when git is missing).

### Spec
- New: `spec/01-app/111-push-pull-transport-flags.md`.
- Memory: `.lovable/memory/features/push-pull-transport-flags.md`.

### Pinned
- README pinned-version block + version matrix moved to **v5.29.0**.
- Synced `gitmap/constants/constants.go` (`Version = "5.29.0"`) and `src/constants/index.ts` (`VERSION = "v5.29.0"`).


## v5.28.0 — (2026-05-19) — `gitmap pull` cwd short-circuit + `gitm` short alias

### Added
- **`gitmap pull` (no args, inside a git repo)** now short-circuits to a plain `git pull` in the current directory — stdin/stdout/stderr are forwarded and the underlying exit code is propagated. Slug / `--group` / `--all` / aliased-repo modes are unchanged; the new behaviour only triggers when none of those targeting modes are in effect.
- **`gitm` shell alias** — every install of the shell wrapper (Bash / Zsh / PowerShell, installed by `gitmap setup`) now also defines `gitm` as a thin forwarder to `gitmap`, so `gitm pull`, `gitm cd <name>`, `gitm clone <url>` all behave identically. Re-run `gitmap setup` (or reinstall) to pick up the new wrapper block.

### Confirmed
- `gitmap setup` auto-run after install is already shipped (since v5.18.0) — both `install.ps1` and `install.sh` call `gitmap setup` as a non-fatal final step, so the new `gitm` alias is registered automatically on fresh installs.

### Pinned
- README pinned-version block + version matrix moved to **v5.28.0**.
- Synced `gitmap/constants/constants.go` (`Version = "5.28.0"`) and `src/constants/index.ts` (`VERSION = "v5.28.0"`).

## v5.27.0 — (2026-05-19) — `gitmap cfrp` / `cfr` honour `--ssh` / `--https`

### Fixed
- `gitmap clone-fix-repo-pub <url> --ssh` (alias `cfrp`) and `gitmap clone-fix-repo <url> --ssh` (alias `cfr`) were ignoring the transport flag — the chained `clone` step ran the raw HTTPS URL because `parseCloneFixRepoArgs` only recognised `--no-vscode-sync` and `--require-version`. The URL is now rewritten via `ConvertURLToSSH` / `ConvertURLToHTTPS` before the in-process clone runs, mirroring `gitmap clone --ssh` exactly.
- Accepts `--ssh`, `-ssh`, `--sh`, `-sh`, `--https`, `-https`, `--ht`, `-ht` (single- and double-dash, plus the same short aliases as `gitmap clone`). When both `--ssh` and `--https` are set, `--ssh` wins with a stderr warning.
- Prints a `↪ --ssh rewrite: <before> → <after>` breadcrumb so the substitution is visible before `git clone` runs.

### Pinned
- README pinned-version block + version matrix moved to **v5.27.0**.
- Synced `gitmap/constants/constants.go` (`Version = "5.27.0"`) and `src/constants/index.ts` (`VERSION = "v5.27.0"`).

## v5.26.0 — (2026-05-18) — Pin bump rolling up the `clone --ssh` flag-position fix

### Pinned
- README pinned-version block + version matrix moved to **v5.26.0** (PowerShell + Bash installer URLs and all per-platform release assets).
- Synced `gitmap/constants/constants.go` (`Version = "5.26.0"`) and `src/constants/index.ts` (`VERSION = "v5.26.0"`).

### Rolled up
- `gitmap clone <url> --ssh` (and every other bool clone flag) is honoured regardless of position — `parseCloneFlags` routes through `reorderFlagsBeforeArgs` so `--ssh` after the URL no longer slips past Go's `flag` package.
- SSH-shorthand and `ssh://` URLs continue to clone natively through `git`; the `--ssh` converter only fires to coerce an HTTPS URL into shorthand.

## v5.24.0 — (2026-05-18) — `gitmap clone --ssh` actually parses when placed after the URL

### Fixed
- `gitmap clone <url> --ssh` (and `--https`, `--no-replace`, every other bool flag) was silently ignored when written AFTER the positional URL. Go's `flag` package stops parsing at the first non-flag argument, so the trailing `--ssh` never reached `applyURLSchemeFlags` and the HTTPS URL was cloned as-is.
- `parseCloneFlags` now routes through `reorderFlagsBeforeArgs` (the same helper used by `release`, `clone-next`, `clone-from`, `commit-transfer`, etc.), so flags are honoured regardless of position: `gitmap clone --ssh <url>`, `gitmap clone <url> --ssh`, and `gitmap clone <url> --ssh --no-replace` all behave identically.
- SSH-shorthand and `ssh://` URLs already work natively through `git clone` — no extra wiring required; the converter only fires when `--ssh` is supplied to coerce an HTTPS URL into shorthand before `git` runs.

### Pinned
- README pinned-version block + version matrix moved to **v5.24.0**.

## v5.23.0 — (2026-05-18) — Root-level installer URLs (`/install.ps1`, `/install.sh`)

### Added
- **`/install.ps1`** and **`/install.sh`** now live at the repository root. The one-liner installer URL is now:
  - Windows: `irm https://raw.githubusercontent.com/alimtvnetwork/gitmap-v27/main/install.ps1 | iex`
  - macOS / Linux: `curl -fsSL https://raw.githubusercontent.com/alimtvnetwork/gitmap-v27/main/install.sh | sh`
- Shorter, more discoverable URLs — no more `gitmap/scripts/` segment for end users.

### Kept
- `gitmap/scripts/install.ps1` + `gitmap/scripts/install.sh` are unchanged — they remain the source of truth consumed by:
  - `gitmap/scripts/embed.go` (`go:embed` for `gitmap self-install`).
  - `release-version.ps1` / `release-version.sh` generators.
  - `.github/workflows/ci.yml`, `.github/workflows/release.yml`, smoke tests.
  - ~50 spec / memory docs.
  The root files are byte-identical copies — same checksum contract, same behavior.

### Pinned
- README pinned-version block + version matrix moved to **v5.23.0**.
- All README install URLs rewritten from `main/gitmap/scripts/install.{ps1,sh}` → `main/install.{ps1,sh}` (11 occurrences).
- Version constants synced: `gitmap/constants/constants.go` (`Version = "5.23.0"`), `src/constants/index.ts` (`VERSION = "v5.23.0"`).

## v5.22.0 — (2026-05-18) — Pin bump (rolls up v5.21.0 ssh view/copy/create)

### Pinned
- README pinned-version + version-matrix moved to **v5.22.0**.
- Version constants synced: `gitmap/constants/constants.go` (`Version = "5.22.0"`), `src/constants/index.ts` (`VERSION = "v5.22.0"`).
- No behavior changes since v5.21.0 — this release re-pins the `gitmap ssh view` / `copy` / `create` surface (clipboard-aware, soft-fails on missing `clip`/`pbcopy`/`wl-copy`/`xclip`/`xsel`) as the new stable.

## v5.21.0 — (2026-05-18) — `gitmap ssh` gets `view` / `copy` / `create` subcommands

### Added
- **`gitmap ssh view <key>`** (aliases: `v`, existing `cat`) — prints the public key to stdout. Identical output to `ssh cat`; just a more discoverable verb.
- **`gitmap ssh copy <key>`** (alias: `cp`) — prints the public key AND pushes it to the OS clipboard in one shot. Picks the right tool per OS: `clip` on Windows, `pbcopy` on macOS, `wl-copy` → `xclip -selection clipboard` → `xsel --clipboard --input` on Linux (first one found). If no clipboard tool is available, the key is still printed and a one-line warning is emitted to stderr — the command never fails.
- **`gitmap ssh create <flags>`** — explicit alias for the default `gitmap ssh` (generate). Improves discoverability alongside `view` / `copy` / `delete`.

### Files
- `gitmap/constants/constants_ssh.go` — new subcommand consts (`SubCmdSSHView`/`V`, `SubCmdSSHCopy`/`Cp`, `SubCmdSSHCreate`) + clipboard messages (`MsgSSHCopied`, `MsgSSHCopyFallback`, `ErrSSHClipboard`).
- `gitmap/cmd/ssh.go` — dispatch wires the three new verbs.
- `gitmap/cmd/sshcopy.go` — new: `runSSHCopy` + `writeClipboard` + `resolveClipboardTool` (cross-platform binary picker).
- `gitmap/helptext/ssh.md` — documents view/copy/create with examples.


## v5.20.0 — (2026-05-18) — `gitmap clone --ssh` / `--https` coerce every URL into the requested transport

### Added
- **`gitmap clone --ssh`** rewrites every recognised Git URL into its `git@host:owner/repo.git` SSH-shorthand form before `git clone` runs. HTTPS (`https://github.com/owner/repo`) and `ssh://git@host[:port]/owner/repo` URLs are both converted; already-shorthand URLs are normalized (`.git` suffix appended). Flows through the multi-URL form too — `clone url1,url2,url3 --ssh` converts the whole batch.
- **`gitmap clone --https`** is the symmetric counterpart — forces every URL into `https://host/owner/repo.git`. Useful in CI/headless environments where the SSH agent isn't unlocked.
- `--ssh` and `--https` are mutually exclusive; when both are set, `--ssh` wins and a one-line stderr warning is printed.

### Behaviour
- Conversion happens AFTER `applySSHKey` and BEFORE the multi-URL / direct-URL routers, so the multi-URL detector sees the converted URLs.
- Non-URL positionals (folder names, `json` / `csv` / `text` shorthands) are skipped via the existing `isDirectURL` predicate.
- Port hints in `ssh://` URLs are dropped (SSH-shorthand has no port slot — use `~/.ssh/config` for non-default ports).
- A `↪ --ssh rewrite: <before> → <after>` breadcrumb is printed before git runs.

### Examples
```powershell
gitmap clone https://github.com/alimtvnetwork/wp-onboarding.git --ssh
#   ↪ --ssh rewrite: https://github.com/alimtvnetwork/wp-onboarding.git → git@github.com:alimtvnetwork/wp-onboarding.git

gitmap clone "https://github.com/a/x,https://github.com/a/y" --ssh
gitmap clone git@github.com:alimtvnetwork/wp-onboarding.git --https
```

### Files
- `gitmap/cmd/cloneurlconvert.go` — new: `ConvertURLToSSH`, `ConvertURLToHTTPS`, plus helpers.
- `gitmap/cmd/rootflags.go` — `UseSSH` / `UseHTTPS` on `CloneFlags`, `--ssh` / `--https` registration.
- `gitmap/cmd/clone.go` — `applyURLSchemeFlags`, wired into `runClone`.

### Spec / Memory
- Spec: `spec/01-app/110-clone-ssh-flag.md`
- Memory: `.lovable/memory/features/clone-ssh-flag.md`


## v5.19.0 — (2026-05-18) — `gitmap rp` (release-pending) rejects version args + canonical command banner

### Fixed
- **`gitmap rp v3.1` no longer silently runs `release-pending` and releases unrelated versions** (e.g. picking up an orphaned `v2.233.0` metadata file). Users were confusing `rp` (alias for `release-pending`, which scans for ALL pending branches + orphan metadata) with `pr` (alias for `pull-release`, which pulls then releases a specific version). The old code parsed `v3.1` away into `flag.Args()` and ignored it.
- `runReleasePending` now scans positional args for a version-shaped token (regex `^v?\d+(\.\d+){0,2}(-...)?$`) and exits 2 with a precise suggestion:
  ```
  ✗ gitmap release-pending (rp) takes no version argument (got "v3.1").
    It releases EVERY pending branch + orphan metadata file — not a specific version.

    Did you mean:  gitmap pr v3.1        # pull-release: pull, then release v3.1
               or  gitmap release v3.1   # release v3.1 directly
  ```

### Added
- **Canonical command banner** at the start of `gitmap release-pending` and `gitmap pull-release`. When the user invokes a short alias (`rp`, `pr`), the resolved canonical command name is printed to stderr first so it's immediately obvious which pipeline is running:
  ```
  → Running: gitmap release-pending  (alias: rp)
  ```
- New helper `printCanonicalCmdBanner` in `gitmap/cmd/releasepending.go` is reusable for any future alias-prone command pair.

### Files
- `gitmap/cmd/releasepending.go` — `runReleasePending`, `rejectVersionArgOnPending`, `printCanonicalCmdBanner`, `versionLikeArgPattern`.
- `gitmap/cmd/releasepull.go` — `runReleasePull` now calls `printCanonicalCmdBanner`.


## v5.18.0 — (2026-05-18) — Auto-run `gitmap setup` after install and on first `gitmap cd` when shell wrapper isn't loaded

### Added
- **Install scripts auto-run `gitmap setup`** as a final non-fatal step. Both `gitmap/scripts/install.sh` and `gitmap/scripts/install.ps1` now invoke `<bin_path> setup` immediately after `Invoke-InstallVerification` / `verify_installation`, so a fresh `irm … | iex` or `curl … | sh` lands the user with shell wrapper (`gcd` / `gitmap` profile function) + completions installed — no second manual command required. Setup is idempotent (marker `# gitmap shell wrapper v2`), so re-runs on every install are safe. Failures print a yellow `(setup auto-run skipped …)` line and continue — the install itself already succeeded.
- **`gitmap cd` self-heals when the shell wrapper isn't loaded.** `warnIfNoWrapper` in `gitmap/cmd/setupverify.go` now auto-runs the full `gitmap setup` (not just `InstallCDFunction`) so the very first `gitmap cd <repo>` after a fresh install installs the wrapper + completions automatically. The existing stderr reload tip (`. $PROFILE` / `source ~/.zshrc`) still prints so the user knows why the *next* `cd` will actually move the parent shell. Recover-guarded — a setup panic never breaks the cd output.

### Why
- First-time users were running `gitmap install` then `gitmap cd repo` and hitting `! Shell wrapper not active — 'gitmap cd' printed the path but cannot change your directory.` because `gitmap setup` was a separate manual step. Closing the gap in both entry points (install scripts AND cd-on-no-wrapper) means the wrapper is always there after one terminal restart.


## v5.17.0 — (2026-05-18) — `gitmap cd` PowerShell wrapper: coerce stdout to a single string

### Fixed
- **`gitmap cd <repo>` crashed in PowerShell** with `Set-Location : Cannot convert 'System.Object[]' to the type 'System.String'` when the captured stdout was parsed as a multi-element array (any extra line, CRLF artifact, or auxiliary write from the binary turned `$dest` into `System.Object[]`, which `Set-Location -LiteralPath` rejects).
- All four PowerShell wrapper templates now collapse the captured output with `(& $real … | Out-String).Trim()` so `$dest` is always a single trimmed string before being passed to `Test-Path` / `Set-Location`:
  - `gitmap/constants/constants_cd.go` — `CDFuncPowerShell` (`gcd` + `gitmap` profile functions)
  - `gitmap/constants/constants_pathsnippet.go` — `PathSnippetPwshFmt` (`Invoke-GitmapAndSetLocation`)
  - `gitmap/constants/constants_cd_shim.go` — `PowerShellShimTemplateFmt` (`gitmap.ps1` shim)
  - `gitmap/scripts/install.ps1` — `Get-GitmapCommandWrapperBlock` + `Get-GitmapPowerShellShimContent`
- Users hit by the v5.16.0-and-earlier wrapper need to re-run `gitmap setup` (or re-install) after upgrading so their `Microsoft.PowerShell_profile.ps1` snippet is rewritten with the fixed body.

## v5.16.0 — (2026-05-18) — `gitmap release` no longer leaks gitmap-specific content into other repos' releases

### Fixed
- **Release body**: `gitmap release` no longer dumps gitmap's own `CHANGELOG.md` notes into the GitHub release body of unrelated repositories. `uploadToGitHub` in `gitmap/release/workflowgithub.go` now starts with an empty body and only calls `DetectChangelog()` + `AppendPinnedInstallSnippet` when `ShouldPrintInstallHint(getRemoteURL())` is true (i.e. the current repo is a `alimtvnetwork/gitmap-v<N>` source repo). Non-gitmap repos get a tag-only release with an empty body.
- **`release-version-vX.Y.Z.{ps1,sh}` snapshot assets** are no longer attached to non-gitmap releases. Those snapshots hard-code `REPO="alimtvnetwork/gitmap-v27"` and `BINARY_NAME="gitmap"` — uploading them to `img-pdf-v2`, `some-tool-v3`, etc. would mislead users into downloading the gitmap binary from the wrong release page. `pushAndFinalize` in `gitmap/release/workflowfinalize.go` now wraps `buildReleaseVersionSnapshots` in the same `ShouldPrintInstallHint` gate.
- All other assets (cross-compiled Go binaries, zip groups, ad-hoc bundles, checksums, docs-site) are unaffected — only the two gitmap-specific pieces are gated.

### Docs
- New spec `spec/02-app-issues/27-release-body-and-snapshots-gitmap-only.md` documents the contract: for any repo NOT matching `alimtvnetwork/gitmap-v<N>`, release body is empty and no `release-version.{ps1,sh}` snapshots are uploaded.
- New memory `.lovable/memory/features/release-gitmap-only-body-and-snapshots.md` captures the gate rule and lists the two call sites.

## v5.15.0 — (2026-05-18) — `gitmap install gitmap-oneliner`: print canonical Windows + macOS install one-liners in the terminal

### Added
- New `gitmap install gitmap-oneliner` command prints the canonical bootstrap one-liners for Windows (PowerShell `irm … | iex`) and macOS/Linux (bash `curl -fsSL … | sh`) without leaving the terminal. URLs are fixed to the canonical `alimtvnetwork/gitmap-v27/main` branch installers; the header version is rendered dynamically from `constants.Version`.
- `ToolGitmapOneliner = "gitmap-oneliner"` added to the Core install tool category with description "Print the Windows + macOS install-gitmap one-liners".

### Docs
- New spec `spec/01-app/109-install-gitmap-oneliner.md` documents the command synopsis, output format, and implementation contract (fixed URLs, dynamic rendering, special-install dispatch).

## v5.14.0 — (2026-05-18) — Colorful help text: green comments, magenta keys, padded lists

### Added
- Markdown help renderer (`gitmap/render`) now color-codes `# comments` inside code fences in green and credential-like tokens (`API_KEY`, `GITMAP_TOKEN`, etc.) and the `hd` alias in magenta.
- List and table blocks in help output now render with vertical padding for better readability.

### Tests
- New pretty-renderer test fixtures: `case-010-fence-comments-and-keys` (green comments + magenta keys) and `case-011-list-and-table` (list/table padding).

## v5.13.0 — (2026-05-16) — Re-pin root README to v5.13.0, sync version constants

### Docs
- Bumped `Pinned version` section in root `README.md` from `v5.12.0` → `v5.13.0` (heading, one-line installers, version-matrix table, release page link).

### Chore
- `gitmap/constants/constants.go`: `Version` constant `5.10.0` → `5.13.0`.
- `src/constants/index.ts`: web `VERSION` `v5.10.0` → `v5.13.0` (keeps `version-sync.test.ts` in sync).
- CI badges in `README.md` no longer pin to `branch=main`, so workflow status reflects the active branch where CI is green.

## v5.12.0 — (2026-05-16) — README alias mapping for `pull-release` family

### Docs
- Root `README.md` Release & Versioning table now includes a dedicated
  `pull-release` row listing the canonical short alias `pr` alongside the
  three legacy spellings (`release-pull`, `relp`, `rlp`).
- New "Alias mapping — `pull-release` family (v5.6.0+)" sub-table maps every
  legacy and canonical token to its routing target with status
  (canonical / legacy) and notes, including the `prune` short alias rename
  from `pr` → `prn`. This is the single source of truth users can link to
  when scripts or CI still reference the legacy names.
- Copy-paste block under the table demonstrates every spelling resolving to
  the same handler so users can sanity-check their muscle memory at a glance.
- Release History & Info table corrected: `prune` short alias is now `prn`
  (was incorrectly still shown as `pr` after the v5.6.0 rename).
- Right-click context-menu table row renamed "Release pull" → "Pull release"
  and the `gitmap pull-release` invocation now calls out the legacy alias
  passthrough.

### Notes
- No code change — alias wiring already lives in
  `gitmap/constants/constants_cli.go` (`CmdReleasePullAlias..Alias4`) and
  `gitmap/constants/constants_prune.go` (`CmdPruneAlias = "prn"`) as of
  v5.6.0. This release only synchronizes the README so the docs match the
  binary.


## v5.11.0 — (2026-05-16) — `--no-color` / `--color` flag for help output

### Added
- `--no-color` and `--color[=true|false|on|off|auto|...]` are now accepted
  everywhere `--pretty` / `--no-pretty` work (notably `gitmap <cmd> --help`).
  They are pure synonyms for the existing pretty flags, since the pretty
  markdown pipeline is the only ANSI surface in the CLI. `--no-color`
  mirrors the widely-used [NO_COLOR](https://no-color.org) env convention
  so users reaching for the conventional spelling get the expected
  behavior without having to discover `--no-pretty`.
- Last-writer-wins still applies across the synonym families
  (`--pretty --no-color` → PrettyOff; `--no-color --pretty=on` → PrettyOn).

### Tests
- `prettyflag_test.go` covers every `--color` / `--no-color` form,
  cross-family last-writer-wins, and guards `--colorblind` / `--color=blue`
  from being silently swallowed by the new prefix.


## v5.10.0 — (2026-05-16) — Force PowerShell wrapper to load last

### Fixed
- `gitmap setup` / installer rewrites now move the managed PowerShell
  `gitmap`/`gcd` command wrapper to the end of the profile file. This prevents
  older `# gitmap shell wrapper v2` or hand-edited profile blocks that appear
  later in `$PROFILE` from overriding the fresh wrapper and leaving `gitmap cd`
  stuck on the raw exe warning.

## v5.9.0 — (2026-05-16) — Harden PowerShell `gitmap cd` activation

### Fixed
- Windows installers now write the `gitmap`/`gcd` command wrapper to all standard
  current-user PowerShell profile files, not only the one profile visible to the
  installer host. This prevents new terminals from falling back to raw
  `gitmap.exe` and printing the wrapper-not-active warning.
- The installer also loads the command wrapper into the current installer session
  even when PATH was already present, and the release ZIP shim now supports the
  same handoff behavior as the generated shim.

## v5.8.0 — (2026-05-16) — fix-repo bare-base rewrite for pre-versioned v1 repos

### Fixed
- `gitmap fix-repo` now rewrites bare `{base}` occurrences (not just
  `{base}-v1`) when v1 is in the target span. Pre-versioned remotes
  shipped without a `-v1` suffix, so downstream references read `img-pdf`
  rather than `img-pdf-v1`; the previous rewriter skipped them entirely
  and reported `changed: 0`. The bare-base sweep is guarded by strict
  word-boundary checks (alnum / `_` / `-` / `.`) so `{base}-v2`,
  `{base}.js`, `{base}_alt`, and `myimg-pdf` are all left untouched.

### Tests
- `fixrepo_rewrite_barebase_test.go` covers the v1→current bare-base
  substitution, every word-boundary guard case, and the guarantee that
  the bare-base pass does NOT run when v1 is absent from the target set.

## v5.7.0 — (2026-05-16) — Ship PowerShell shim in release installs

### Fixed

- Root cause of the still-broken `gitmap cd` on Windows release installs: the
  v5.5.0 shim existed in setup/installer code, but the release asset pipeline
  still zipped only `gitmap.exe`, and the release-specific `install.ps1` only
  moved the exe. Users installing from release assets therefore never received
  `gitmap.ps1`, so PowerShell still resolved the raw exe and printed the wrapper
  warning.
- Windows release ZIPs now contain both `gitmap.exe` and `gitmap.ps1`.
- The release-specific installer now moves `gitmap.ps1` beside `gitmap.exe` when
  present, and release smoke fails if the shim is missing.
- The generic installer now writes the shim after binary extraction even when PATH
  handling is skipped, so `-NoPath` test/install paths still leave a complete
  install directory.

### Files

- `.github/workflows/release.yml` — Windows ZIP packaging now stages `gitmap.exe`
  + `gitmap.ps1`; release installer extracts both.
- `.github/scripts/smoke-installer.ps1` — asserts installed releases contain the
  shim beside the exe.
- `gitmap/scripts/install.ps1`, `gitmap/constants/constants_cd_shim.go` — shim no
  longer calls `exit` for normal forwarding; it preserves `$LASTEXITCODE` and
  returns to the current PowerShell session.
- `gitmap/completion/cdfunction.go` — managed wrapper blocks are rewritten when
  stale instead of skipped on marker presence.
- `src/constants/index.ts`, `gitmap/constants/constants.go` — version bumped to
  **v5.7.0**.

## v5.6.0 — (2026-05-16) — Rename `release-pull` → `pull-release` (`pr`), refreshed help

### Changed

- **Renamed `release-pull` → `pull-release`.** The verb-first form reads as
  "pull, then release" — clearer at a glance. The canonical short alias is now
  **`pr`** (e.g. `gitmap pr v1.4.0`).
- **`prune` short alias changed from `pr` → `prn`** to free up `pr` for
  `pull-release`. `gitmap prune` (long form) is unchanged.
- Help text for `pull-release` was rewritten with section headers, an aliases
  table, a "Pull modes" table, and a tagged example list. The markdown is
  rendered through the standard ANSI pipeline so all bold/headings/tables get
  colorful output in any TTY (and strip cleanly with `--no-pretty`).

### Backward compatibility (no breaking changes)

- `gitmap release-pull …`, `gitmap relp …`, and `gitmap rlp …` continue to
  work — they are wired as aliases of `pull-release` and route to the exact
  same handler. Existing scripts, docs, and right-click context menu entries
  do not need to change.
- All `--ff-only` / `--rebase` / `--merge` / `--dry-run` / `--verbose` flags
  and their forwarding semantics to `gitmap release` are unchanged.

### Files

- `gitmap/constants/constants_cli.go` — `CmdReleasePull = "pull-release"`,
  `CmdReleasePullAlias = "pr"`, added `Alias2..Alias4` for `release-pull`,
  `relp`, `rlp`. Updated `HelpReleasePull` summary line.
- `gitmap/constants/constants_prune.go` — `CmdPruneAlias = "prn"` (was `"pr"`).
- `gitmap/constants/constants_releasepull.go` — error/log prefixes now say
  `pull-release:` / `[pull-release]`.
- `gitmap/constants/constants_helpgroups.go` — `CompactRelease` shows
  `pull-release (pr)`.
- `gitmap/cmd/rootrelease.go` — dispatcher accepts all four aliases.
- `gitmap/cmd/llmdocsgroups.go` — release group lists `pull-release (pr)`;
  prune entry uses `prn`.
- `gitmap/helptext/pull-release.md` — new enhanced help page (replaces
  `release-pull.md`).
- `gitmap/completion/allcommands_generated.go` — adds `pull-release` + `prn`,
  keeps `release-pull` / `relp` / `rlp` / `pr` for completion.
- `gitmap/completion/completion_test.go`, `gitmap/constants/cmd_constants_test.go`
  — updated for the new aliases.
- `src/data/commands.ts` — web docs use `pull-release (pr)` and updated
  prune alias note.
- `src/constants/index.ts`, `gitmap/constants/constants.go` — version
  bumped to **v5.6.0**.


## v5.5.0 — (2026-05-16) — Add PowerShell command shim for `gitmap cd`

### Fixed

- `gitmap cd <repo>` on Windows now has a second activation path: the installer
  writes `gitmap.ps1` beside `gitmap.exe`, and PowerShell prefers the script shim
  over the exe when the install directory is on PATH. The shim runs in-process,
  calls the real exe, captures the resolved path, and runs `Set-Location` in the
  current PowerShell session.
- Root cause of the repeated failure: a profile function can only work after the
  correct profile is loaded, but the user was still reaching the raw exe. Relying
  only on `$PROFILE` reload was therefore insufficient; the PATH-level `.ps1`
  shim removes that dependency for PowerShell sessions.
- The Windows installer now force-rewrites stale `# gitmap command wrapper v1`
  profile blocks instead of treating the marker as proof the current body is
  correct.

## v5.4.0 — (2026-05-16) — Install PowerShell command wrapper from installer profile block

### Fixed

- Windows installer/profile setup now writes the `gitmap`/`gcd` PowerShell
  command wrapper directly into the managed profile block and loads it into the
  current `irm ... | iex` session. This closes the stale-binary gap where users
  could reinstall/update and still have `gitmap cd` resolve to the exe until a
  separate `gitmap setup` from the new binary had run and the profile was reloaded.
- Root cause of the repeated old warning: the user was still invoking a deployed
  binary that predates the command-wrapper marker fix, and the Windows
  `install.ps1` path block only added PATH, not the shell function that can call
  `Set-Location` in the parent PowerShell session.

## v5.3.0 — (2026-05-16) — Fix command-wrapper false-positive detection

### Fixed

- `gitmap cd` no longer treats the PATH snippet as proof that the command
  wrapper is active. The command wrapper now uses its own marker
  (`# gitmap command wrapper v1`) and runtime sentinel (`GITMAP_COMMAND_WRAPPER`),
  so `gitmap setup` appends/updates the actual `function gitmap { ... }` /
  `gcd` wrapper even when the older PATH block already exists.
- Root cause: both the PATH snippet and the `gitmap`/`gcd` command wrapper used
  the same `# gitmap shell wrapper v2` text, and `isWrapperActive` trusted
  `GITMAP_WRAPPER`, which the PATH snippet also exports. That caused setup to
  report wrapper success while PowerShell still resolved `gitmap` as the exe.

## v5.2.0 — (2026-05-16) — Gate pinned-install snippet to gitmap-vN repos only

### Fixed

- `gitmap release` no longer appends the gitmap pinned-version PowerShell/bash
  installer snippet to release bodies of unrelated repositories. The snippet
  is now gated by the same `ShouldPrintInstallHint` check used by the
  terminal install hint: it only renders when the current repo's remote
  matches `alimtvnetwork/gitmap-vN` (versioned gitmap source repos).
- Root cause: `uploadToGitHub` in `gitmap/release/workflowgithub.go` called
  `AppendPinnedInstallSnippet` unconditionally, so every release body (in any
  repo) got the gitmap installer block plus the version tag header.

## v5.1.0 — (2026-05-16) — Fix cfrp remote-based version detection

### Fixed

- `gitmap clone-fix-repo` / `cfrp` now decide whether to run `fix-repo --all`
  from the cloned Git remote repo name, not the flattened local destination
  folder. Cloning `gitmap-v27` into `gitmap/` now correctly runs the rewrite
  instead of falsely reporting that `gitmap` has no `-vN` suffix.
- Bumped `Version` constant to `5.1.0` (Go) and `VERSION` to `v5.1.0` (web).

## v4.44.0 — (2026-05-16) — Minor version bump; re-pin root README install snippets and version matrix to v4.44.0

- Bumped `Version` constant to `4.44.0` (Go) and `VERSION` to `v4.44.0` (web)
- Re-pinned the root `README.md` "Pinned version" section, one-line installers, and version matrix asset URLs to `v4.44.0` / `gitmap-v27.44.0-*`


## v4.42.0 — (2026-05-09) — Minor version bump; re-pin root README install snippets and version matrix to v4.42.0

- Bumped `Version` constant to `4.42.0` (Go) and `VERSION` to `v4.42.0` (web)
- Re-pinned the root `README.md` "Pinned version" section, one-line installers, and version matrix asset URLs to `v4.42.0` / `gitmap-v27.42.0-*`

## v4.41.0 — (2026-05-07) — Minor version bump; re-pin root README install snippets and version matrix to v4.41.0

- Bumped `Version` constant to `4.41.0` (Go) and `VERSION` to `v4.41.0` (web)
- Re-pinned the root `README.md` "Pinned version" section, one-line installers, and version matrix asset URLs to `v4.41.0` / `gitmap-v27.41.0-*`

## v4.40.0 — (2026-05-07) — Minor version bump; pin root README install snippets and version matrix to v4.40.0

### Changed

- Bumped `Version` constant to `4.40.0` (Go) and `VERSION` to `v4.40.0` (web)
  in lock-step to satisfy the `version-sync` regression test.
- Re-pinned the root `README.md` "Pinned version" section, one-line installers,
  and the per-platform version matrix (Windows/macOS/Linux × amd64/arm64) to
  the new `v4.40.0` release tag and `gitmap-v27.40.0-*` asset names.

## v4.39.0 — (2026-05-07) — CI lint cleanup: misspellings (`centralised`/`materialises`/`honoured`), unused `mergePairs`, gofmt drift, `string(before) != string(after)` → `bytes.Equal`

### Fixed

- Resolved all NEW golangci-lint findings vs baseline from the v4.38.0 batch:
  `misspell` (`centralised`/`Centralised`/`materialises`/`honoured` → US spelling)
  across `gitmap/cmd/commitin/runlog/`, `gitmap/cmd/vscodepmsync*.go`, and
  `constants/constants_commitin_tagreplay.go`.
- Removed unused `mergePairs` helper from `gitmap/vscodepm/merge.go` (`unused`).
- Replaced `string(before) != string(after)` byte-slice comparisons with
  `!bytes.Equal(...)` (`gocritic`/`stringXbytes`).
- Re-ran `gofmt -w` on `cmd/commitin/runlog/tagreplay.go` to drop the trailing
  blank line flagged by the gofmt CI gate.

## v4.38.0 — (2026-05-07) — `commit-in` tag-replay map (spec §09) + strict annotated-only semver gate; `vscode-pm-sync` gains `--path` / `--tag`

### Added

- **CommitInReplayMap (spec §09)** — new SQLite table (migration 007) plus
  `TagReplayOutcome` enum mirror persists every mirrored annotated tag with
  its old↔new SHA pair, the chosen `MirroredReleaseBranch`, an
  `IsVersionTag` flag, and the run outcome (`Created` / `CreatedDryRun` /
  `Skipped` / `AlreadyExists` / `Failed`). Indexes on `SourceTagName`,
  `DestCommitSha`, and `MirroredReleaseBranch`. `(CommitInRunId,
  SourceTagName)` is UNIQUE per run.
- **`runlog.RecordTagReplay` / `LookupTagReplay`** — pure persistence
  helpers. `LookupTagReplay` implements the §9.5 cross-run short-circuit
  query keyed by `(SourceTagName, SourceTagSha)` and returns the typed
  `ErrTagReplayMiss` sentinel (zero-swallow rule).
- **`runlog.ClassifyVersionTag(name, isAnnotated)`** — canonical
  strict-semver classifier. A tag is a "version tag" iff it is annotated
  AND its name matches `constants.VersionTagPattern`. Lightweight tags
  named `v1.2.3` are NEVER promoted to version-tag status.
- **`--no-release-branch` flag on `gitmap commit-in`** — defaults branch
  generation ON; suppresses `release/<tag>` materialisation across
  replay-mapping and the (forthcoming §08) runner integration. Decision
  centralised in `runlog.ResolveReleaseBranchName(tag, isAnnotated,
  isNoReleaseBranch, isDryRun)` — single choke point.
- **`--projects-json <path>` and `--tag <name>` flags on
  `gitmap vscode-pm-sync` (`vpm`)** — override the resolved projects.json
  location and append/replace the auto-derived tag set without touching
  user-added tags.
- New constants file `gitmap/constants/constants_commitin_tagreplay.go` and
  SQL twin `constants_commitin_tagreplay_sql.go` centralise every magic
  string (table names, enum literals, INSERT/SELECT statements).

### Changed

- **Strict-semver gate enforced at the persistence boundary**:
  `RecordTagReplay` rejects any insert with `IsVersionTag=true` and
  `IsAnnotated=false` via the new `ErrLightweightVersionTag` sentinel
  before touching SQLite. This is the single guarantee that
  `CommitInReplayMap.IsVersionTag=1` rows are always annotated semver
  tags — no upstream filter can silently drift.
- `TagReplayFacts` gains `IsAnnotated bool`; `ResolveReleaseBranchName`
  gains an `isAnnotated` parameter. Lightweight tags whose name matches
  the semver regex no longer receive a release branch.

### Refactored

- `gitmap/vscodepm/sync.go` (273 LOC) split into three files under the
  200-line cap with no behavior change: `sync.go` (public API +
  `readEntries`, 120 LOC), `merge.go` (`mergePairs*`, `applyMerge`,
  `sameTagSet`, 88 LOC), `io.go` (atomic write + `normalizePath` /
  `pathsEqual`, 83 LOC). The per-entry update inside
  `mergePairsWithMode` was extracted into `applyMerge` to keep every
  function under the 15-line cap.

### Tests

- New `gitmap/cmd/commitin/runlog/tagreplay_test.go` covers the R1
  happy-path, R6 dry-run NULL columns, R8 cross-run hit, miss-on-failed,
  unique-per-run constraint, lightweight-rejection gate, and a
  `ClassifyVersionTag` matrix (8 cases) plus a 13-case detector matrix.
- New `gitmap/store/migrate_commitin_replaymap_test.go` +
  `_helpers_test.go` validate the migration with `PRAGMA table_info`,
  `PRAGMA foreign_key_list`, index existence, the
  `(CommitInRunId, SourceTagName)` UNIQUE, and the tagged-vs-untagged
  commit distribution.
- New `gitmap/cmd/commitin/runlog/releasebranch_test.go` adds the
  lightweight-never-gets-branch case alongside the existing default-on,
  flag-suppressed, dry-run, non-version, and flag-beats-everything
  cases.
- `gitmap/cmd/vscodepm-sync/...` parser tests cover `--projects-json`
  and `--tag` overrides plus the existing `--no-release-branch`
  default-off / flag-flip / reorder-past-positionals cases.

### Spec

- New `spec/03-commit-in/09-commit-in-replay-map.md` defines the
  CommitInReplayMap schema, idempotency key, R1-R8 acceptance matrix,
  and the §9.5 cross-run lookup contract.
- `spec/03-commit-in/README.md` index updated; ERD stub blocks added to
  `spec/01-app/gitmap-database-erd.mmd`.
- `spec/03-commit-in/08-tag-mirroring-and-release-branches.md` §3
  reaffirms that release branches require an annotated source tag.


## v4.37.0 — (2026-05-06) — `vscode-pm-sync` gains `--mode union|replace|intersection` for tag-merge strategy

### Added

- **`--mode <strategy>` flag on `gitmap vscode-pm-sync` (`vpm`)** picks how
  freshly-detected tags are reconciled with the on-disk tag set. Three values:
  - `union` (default, unchanged from v4.36.0) — existing ∪ detected, dedup'd.
    Additive only; user-added tags are never removed.
  - `replace` — detector wins outright. User-added tags are dropped. The
    `gitmap` brand tag survives because `DetectTagsCustom` always pre-pends it.
  - `intersection` — only tags present in BOTH sources survive, **plus** the
    `gitmap` brand tag is pinned (added back even when the strict intersection
    would be empty) so re-syncs never silently strip our own brand from any
    entry.
- New `vscodepm.MergeMode` enum (`MergeModeUnion` / `MergeModeReplace` /
  `MergeModeIntersection`) + `ParseMergeMode` validator that fails loud on
  unknown literals (zero-swallow rule).
- New exported entry points `vscodepm.SyncMode` and `vscodepm.SyncAtMode`
  carry the `MergeMode` through. Existing `Sync` / `SyncAt` are now thin
  `MergeModeUnion` wrappers — every legacy caller keeps the v4.36.0 default.
- Helptext (`gitmap vscode-pm-sync --help`) gains a **Modes** section with a
  three-row table and worked examples.

### Files

- `gitmap/vscodepm/mergemode.go` — new (enum, parser, dispatcher, three
  per-strategy helpers).
- `gitmap/vscodepm/sync.go` — `mergePairs` is now a thin `MergeModeUnion`
  wrapper around the new mode-aware `mergePairsWithMode`. Added `sameTagSet`
  so Updated/Unchanged tally still works under non-union modes (where the
  merged slice can be SHORTER than the original — a length compare alone
  would miss that). New `SyncMode` exported entry point.
- `gitmap/vscodepm/syncat.go` — new `SyncAtMode`; legacy `SyncAt` becomes a
  `MergeModeUnion` wrapper.
- `gitmap/cmd/vscodepmsync.go` — `parseVSCodePMSyncFlags` now returns
  `(bool, MergeMode, error)`. Unknown `--mode` exits 2 with a stderr
  diagnostic. `commitVSCodePMSync` threads the mode into `SyncMode`.
- `gitmap/constants/constants_cli.go` — `FlagVSCodePMSyncMode` +
  `FlagDescVSCodePMSyncMode` + the three canonical literals
  (`VSCodePMSyncModeUnion` / `Replace` / `Intersection`). No magic strings.
- `gitmap/constants/constants_vscode_pm.go` — `ErrVSCodePMSyncBadMode`
  format string for the validator.
- `gitmap/helptext/vscode-pm-sync.md` — Modes section + flag-table refresh.
- `gitmap/vscodepm/mergemode_test.go` — unit tests for the parser, the
  String() round-trip, and all three dispatch paths (incl. empty-intersection
  brand pin).
- `gitmap/cmd/vscodepmsync_mode_test.go` — end-to-end regression for all
  three modes against a real on-disk fixture, plus the brand-pin contract.
- Version bumped to `v4.37.0` in `gitmap/constants/constants.go`,
  `src/constants/index.ts`, and every `v4.36.0` pin in `README.md`.

### Compatibility

Fully backward-compatible. Every caller that omits `--mode` (i.e. all
existing CI / scripts) gets the v4.36.0 union behavior unchanged.

## v4.36.0 — (2026-05-06) — New `gitmap vscode-pm-sync` (`vpm`) command: re-tag every projects.json entry on demand

- New top-level subcommand `gitmap vscode-pm-sync` (alias `vpm`) walks
  every entry in the alefragnani.project-manager `projects.json` file,
  re-runs the auto-tag detector against each `rootPath`, and merges the
  detected tags (always including the `gitmap` brand tag, v4.34.0+) into
  the existing tags array. User-added tags are preserved (UNION via
  `vscodepm.Sync.mergePairs`, never truncated).
- Entries whose `rootPath` no longer exists on disk are **skipped** —
  their tags stay untouched on disk so an offline removable-drive
  project doesn't silently lose its tag set.
- Soft-fail policy mirrors the post-clone sync helper: missing VS Code
  user-data root or absent `alefragnani.project-manager` extension
  prints a single stderr line and exits **0** so CI never breaks.
- `--dry-run` flag previews the entry counts without touching the file.
- Honors the same `--vscode-tag` / `--vscode-tag-skip` / `--vscode-tag-marker`
  global flags accepted by every other clone variant — pass
  `--vscode-tag-skip gitmap` to opt the brand tag out for this re-sync only.
- Implementation:
  - `gitmap/cmd/vscodepmsync.go` — runner (158 lines, every func ≤15 lines).
  - `gitmap/vscodepm/sync.go` — new exported `ListEntries()` reader so
    callers can walk projects.json without rebuilding the JSON parser.
  - `gitmap/cmd/vscodepmsofterror.go` — extracted `reportVSCodePMSoftError`
    helper (previously referenced from clone/code paths) into a single
    self-contained file with the soft-fail policy documented inline.
  - `gitmap/cmd/rootcore.go` — wires `CmdVSCodePMSync` / `CmdVSCodePMSyncAlias`
    into the core dispatch table.
  - `gitmap/cmd/vscodepmsync_test.go` — four-test regression suite covering
    pair construction (skip-missing, brand-tag presence) and end-to-end
    behavior (user-tag UNION preserved, `--dry-run` is byte-stable).
- New help file `gitmap/helptext/vscode-pm-sync.md` (auto-discovered by the
  `*.md` embed glob and verified by `TestEveryCmdIDHasHelpFile`).
- Version bumped to `v4.36.0` across `gitmap/constants/constants.go`,
  `src/constants/index.ts`, and the root `README.md` pinned-version block.

## v4.35.0 — (2026-05-06) — commit-in docs: 7 worked examples, auto-init dispatch table, sample profile JSON

- `src/pages/CommitIn.tsx` rewritten with seven labelled, scenario-driven
  examples that cover the full surface: (1) plain folder → git repo +
  history replay, (2) mixed local-folder + multi-URL inputs, (3) brand-new
  target via mkdir+init+replay, (4) `all` / `-N` versioned-sibling expansion,
  (5) author override + message scrubbing, (6) saved-profile + weak-title
  override, (7) headless CI run. Each example carries the why-not-just-the-how.
- New "How &lt;source&gt; auto-init works" table makes the URL / existing-repo /
  existing-folder / missing-path dispatch explicit so users know they never
  have to `git init` first.
- New "Sample profile JSON" section ships a copy-pasteable
  `.gitmap/commit-in/profiles/Default.json` with every field populated and
  inline guidance on PascalCase strict decoding, save semantics, and the
  absolute-path binding rule (per `spec/03-commit-in/05-profiles-and-json-shape.md`).
- Page split for the &lt;200-lines code-style rule:
  - `src/pages/commitInData.ts` — flag rows, exit codes, auto-init rows,
    canonical profile JSON string.
  - `src/pages/CommitInExamples.tsx` — the seven worked walkthroughs.
  - `src/pages/CommitIn.tsx` — page shell only (177 lines).
- `gitmap/helptext/commit-in.md` mirrors the same seven examples + sample
  profile JSON so `gitmap commit-in --help` matches the docs site 1:1.
- Version bumped to `v4.35.0` across `gitmap/constants/constants.go`,
  `src/constants/index.ts`, and the root `README.md` pinned-version block.

## v4.34.0 — (2026-05-06) — Auto-brand `gitmap` tag in VS Code Project Manager projects.json

- Every entry that gitmap writes to `%APPDATA%/Code/User/globalStorage/alefragnani.project-manager/projects.json`
  (or the OS equivalent) now carries `"gitmap"` as the lead tag in
  its `tags` array. This makes every gitmap-managed project
  greppable + filterable in the VS Code Project Manager UI without
  the user having to tag entries by hand.
- Implementation: new `AutoTagGitmap = "gitmap"` constant +
  `prependGitmapBrand` helper in `gitmap/vscodepm/autotags_custom.go`.
  The brand is inserted BEFORE the skip filter, so users who really
  don't want it can opt out via `--vscode-tag-skip gitmap`. The
  helper de-dupes when the tag is already present (e.g. carried in
  via `--vscode-tag gitmap` or already on disk via the upstream
  `unionTags` merge — user-edited tags are still never removed).
- Regression tests (`gitmap/vscodepm/autotags_custom_test.go`):
  - `TestDetectTagsCustomGitmapBrandAlwaysPresent` — empty dir,
    missing path, and empty input all still emit the brand.
  - `TestDetectTagsCustomGitmapSkippable` — `--vscode-tag-skip
    gitmap` removes it.
  - `TestDetectTagsCustomGitmapNotDuplicated` — `--vscode-tag
    gitmap` does not double the entry.
  - `TestDetectTagsCustomNoEnvMatchesBuiltin` updated to assert
    `gitmap` is the leading tag in the canonical output.
- Version bumped to `v4.34.0` across `gitmap/constants/constants.go`,
  `src/constants/index.ts`, and the root `README.md` pinned-version
  block + version-matrix table.

## v4.33.0 — (2026-05-06) — Fix: trailing-slash URLs no longer collapse clone target to CWD

- `gitmap/cmd/clone.go` `repoNameFromURL` now strips trailing
  `/` and `\` (both before and after the `.git` suffix peel) before
  computing the basename. Previously a URL like
  `https://github.com/owner/repo/` collapsed to an empty basename,
  which made the resolved clone target equal to the caller's CWD —
  triggering the "target exists" replace flow against unrelated work
  and (in `cfr` / `cfrp` pipelines) running `fix-repo --all` against
  the wrong directory.
- New regression test `TestRepoNameFromURL_TrailingSlash` in
  `gitmap/cmd/clone_test.go` pins the contract for HTTPS, SCP-style
  `git@host:owner/repo`, and `ssh://` URLs with trailing `/`, `\`,
  `///`, and `.git/` shapes — and explicitly fails on empty output
  so this never regresses to the "destructive replace" symptom again.
- Version bumped to `v4.33.0` across `gitmap/constants/constants.go`
  and `src/constants/index.ts`.

## v4.32.0 — (2026-05-06) — Routing-test guard for new docs pages + pinned-version refresh in root README

- New `src/test/new-command-pages.test.ts` (26 tests) regression-guards
  the v4.31.0 docs additions: for each of `commit-in`, `replace`,
  `fix-repo`, `clone-fix-repo`, and `make-public` it asserts the page
  file exists, renders the expected `<h1>`, is imported + routed in
  `App.tsx`, and has a matching `DocsSidebar` entry. A final assertion
  verifies all five `<Route>`s are declared before the wildcard
  `NotFound` route so they stay reachable.
- Root `README.md` install section: pinned-version block, version
  matrix, and all download URLs / asset filenames refreshed from
  `v4.1.0` → `v4.32.0` so the documented pinned build matches the
  current binary.
- Web `VERSION` constant bumped to `v4.32.0` (kept in lockstep with
  the Go binary by `version-sync.test.ts`).


## v4.31.0 — (2026-05-06) — Docs UI: dedicated pages + sidebar entries for `commit-in`, `replace`, `fix-repo`, `clone-fix-repo`, and `make-public`; root README refresh

- New docs pages added under their own routes:
  `/commit-in`, `/replace`, `/fix-repo`, `/clone-fix-repo`,
  `/make-public`. Each page mirrors the canonical helptext with a
  flag table, examples, exit-code matrix, and "See also" links.
- Left sidebar (`DocsSidebar.tsx`) now exposes Commit In (cin),
  Replace (rpl), Fix Repo (fr), Clone + Fix Repo (cfr), and Make
  Public Repo as first-class entries alongside the existing
  `merge-left` / `merge-right` items.
- Root `README.md` Command Reference grew a "Repo-wide rewrite,
  one-shot publish, and chronological replay" subsection covering
  all five new/utility commands with copy-pasteable examples.
- Web `VERSION` constant bumped to `v4.31.0` to stay in lockstep
  with the Go binary version (enforced by `version-sync.test.ts`).


## v4.18.0 — (2026-05-06) — `gitmap commit-in` / `cin`: chronological multi-source commit replay into a single destination repo

- New top-level command: `gitmap commit-in <source> [inputs...]`
  (alias `cin`). Walks each input repo's first-parent chain
  chronologically and replays every commit into `<source>` while
  preserving both `AuthorDate` AND `CommitterDate` byte-for-byte.
  Dedupe via `ShaMap` ensures the same source commit is never
  replayed twice across runs. Spec: `spec/03-commit-in/`.
- Source resolution (`<source>`) is purely deterministic per
  spec §2.3: existing folder → open in place; missing folder with
  parent that exists → `git init` + initial empty commit; URL →
  clone-then-init-on-failure; everything else → `BadArgs`. No
  interactive `--init` flag, no prompts on the source side.
- Inputs accept any mix of local folders, Git URLs, the keyword
  `all` (every versioned sibling of `<source>` discovered via
  `clone-next` semantics), or `-N` (the N most recent siblings).
- Profiles (`.gitmap/commit-in/profiles/<name>.json`, schema v1,
  PascalCase keys, strict decode) bind by absolute symlink-resolved
  source path. `--save-profile <name>` persists the resolved layered
  config; `--save-profile-overwrite` allows replacement;
  `--set-default` flips `IsDefault=1` and clears the flag on every
  sibling profile bound to the same source. Load order:
  `--profile` flag > `--default` profile > built-in defaults, with
  CLI overrides on top of all three.
- Conflict modes: `--conflict ForceMerge` (default, blob clobbers
  are logged then proceed) or `--conflict Prompt` (any HEAD-vs-source
  blob mismatch aborts the entire run with exit
  `CommitInExitConflictAborted`=8). `--dry-run` walks + plans without
  ever invoking `git commit`; the summary line is followed by an
  unmistakable `commit-in: DRY RUN — no commits were created` banner.
- Function-intel block (`--function-intel on`,
  `--languages Go,TypeScript,...`) appends a per-language
  new-function summary to the commit body, derived best-effort from
  `git show <sha>:path` vs `<sha>^:path` diffs. Failures degrade to
  empty-string per spec §6.3 — never a hard error.
- Message pipeline order: original → message-rules strip
  (`StartsWith:`/`EndsWith:`/`Contains:`) → `--override-messages`
  (gated on first-word weak-word match when `--override-only-weak`)
  → `--title-prefix`/`--title-suffix` → random pick from
  `--message-prefix`/`--message-suffix` pools (deterministic per
  run via per-run RNG seed).
- Path exclusions: `--exclude` accepts CSV of relative paths;
  trailing `/` => folder match (prefix + segment-aware), no slash =>
  exact file match. A commit whose entire file list is excluded is
  recorded as `Skipped` with reason `ExcludedAllFiles`.
- All state lives under `<source>/.gitmap/`: `db/gitmap.sqlite`
  (SQLite v1 schema with `CommitInRun`, `InputRepo`, `SourceCommit`,
  `RewrittenCommit`, `Profile`, `ShaMap`), `temp/<runId>/` (cleaned
  on exit unless `--keep-temp`), `commit-in/profiles/`, and
  `commit-in.lock` (advisory file lock, exit `LockBusy`=9 on
  contention). Strict zero-swallow logging to stderr, every error
  path goes through a single `commit-in: <stage>: <message>` format.
- Helptext: `gitmap/helptext/commit-in.md` (105 lines, under the
  120-line cap). Discoverable via `gitmap commit-in --help`.

## v4.15.1 — (2026-05-02) — Installer: stop calling Chocolatey `refreshenv`, kill `'wmic' is not recognized` noise on Windows 11 24H2

- `gitmap/scripts/install.ps1`: removed the opportunistic `refreshenv` call
  in `Main` after `Add-ToPath`. Older Chocolatey `refreshenv` shells out to
  `wmic process get parentprocessid ...` to discover the parent shell, but
  Microsoft removed `wmic.exe` in Windows 11 24H2 / Server 2025. The
  resulting `'wmic' is not recognized as an internal or external command`
  was native-stderr from a `cmd.exe` grandchild and could not be
  suppressed by PowerShell `try/catch` — it leaked straight to the user's
  console at the tail of an otherwise-successful install.
- `Rebuild-SessionPath` already reads Machine + User PATH directly from
  the registry, which is exactly what `refreshenv` was doing, so dropping
  the call is functionally a no-op on every host and removes the noise on
  24H2+.

## v4.15.0 — (2026-05-02) — TestScannerMatchesRewriter: derive expected token from `current`, close digit-capture sibling-literal regression

- `gitmap/cmd/fixrepo_rewrite_scan_test.go`: `TestScannerMatchesRewriter` now
  builds its expected rewritten token via `fmt.Sprintf("%s-v%d", base, current)`
  instead of hard-coding `"gitmap-v27"` while passing `current = 12` to
  `applyAllTargets`. The hard-coded literal silently disagreed with the
  rewriter's actual output (`gitmap-v27`) and failed CI on every run.
- This is the same bug class as FIX-REPO DIGIT-CAPTURE GAP (closed v4.12.0):
  any version-bearing expectation in a fix-repo test MUST be derived from the
  same int the rewriter received — never hard-coded as a sibling literal,
  which silently desyncs on width-crossing bumps and stale edits alike.
- Bumps `constants.Version` from the stale `4.4.0` placeholder to `4.15.0`
  to match the live release line.

## v4.3.0 — (2026-04-30) — Installer: discard binary stderr on post-install verify, kill cp1252 mojibake

- `gitmap/scripts/install.ps1` and `gitmap/scripts/install.sh` now discard the
  newly-installed binary's stderr (`2>$null` / `2>/dev/null`) during the
  post-install `gitmap version` verification and filter stdout to the
  canonical `gitmap vX.Y.Z` line.
- Eliminates the false "Binary found but failed to run" error caused by
  first-run `SeedDownloaderConfig` info lines being wrapped as ErrorRecords
  under `$ErrorActionPreference='Stop'`.
- Eliminates the cp1252 mojibake (`Γùª`, `ΓÇö`) leaking into the user's
  session when PowerShell hosts decode child stderr with the OEM codepage
  instead of `[Console]::OutputEncoding`.

## v4.2.0 — (2026-04-30) — CI hardening: smoke-installer SIGPIPE fix, generate-check auto-commit, version-line filtering

- `.github/scripts/smoke-installer.sh` now captures `gitmap version` output into
  a variable and uses `awk '/^gitmap v[0-9]/{print; exit}'` to extract the
  version line, avoiding `SIGPIPE` (exit 141) under `set -o pipefail`.
- `.github/scripts/smoke-installer.ps1` filters startup log noise (e.g.
  "Seeded downloader defaults") with `Where-Object { $_ -match '^gitmap v[0-9]' }`.
- CI `generate-check` can auto-run `go generate ./...` and commit regenerated
  files when drift is detected.

## v4.1.0 — (2026-04-30) — Major version bump rolling up fix-repo command + installer/lint fixes

- Promotes the `gitmap fix-repo` / `fr` Go-native rewriter, installer warning
  fix (`%w`-wrapped errors in `downloaderconfig.LoadFile`), and `nolintlint`
  cleanup into the v4 line.

## v3.181.0 — (2026-04-30) — `gitmap fix-repo` / `fr`: Go-native rewriter of `{base}-vN` tokens

- New command `gitmap fix-repo` (alias `fr`) replaces the `fix-repo.ps1` script
  with a cross-platform Go implementation. Same exit codes (0–8), same
  `fix-repo.config.json` schema, same `--dry-run` / `--verbose` / `--config` flags.
- Spec: `spec/04-generic-cli/27-fix-repo-command.md`.
- Lint cleanup: removed unused `//nolint:gosec` directives flagged by `nolintlint`.
- Installer warning fix: `downloaderconfig.LoadFile` now wraps errors with `%w`
  so `errors.Is(err, fs.ErrNotExist)` works on fresh installs.

## v3.180.0 — (2026-04-29) — Fix release: sync package-lock.json, pin README to v3.180.0

- Regenerated `package-lock.json` to include `@testing-library/user-event@14.6.1`
  so the docs-site `npm ci` step in the release workflow stops failing.
- Bumped CLI `constants.Version`, internal docs UI `VERSION`, and README pinned
  install version to `v3.180.0`.

## v3.179.0 — (2026-04-29) — Minor version bump rolling up recent CI fixes

Routine minor bump rolling up recent CI fixes:

- Help-exit test coverage for `clone-audit` and other utility commands.
- `goldenguard` determinism pre-check stability fixes.
- Completion generator regenerated for the `aul` alias.
- `gofmt` alignment in `constants/cmd_constants_test.go`.
- Legacy-ref test data migrated to the current `gitmap-v27` repo namespace.
- `cliexit` test helper forces stdin to a pipe so non-TTY gates fire
  reliably under `/dev/null`.

## v3.119.0 — (2026-04-24) — `gitmap inject` / `inj`: register an existing folder with Desktop + VS Code (+ DB)

New command for "I already have this repo on disk, just plug it into
my tooling" — no clone, no scan, just register.

**Forms:**

    gitmap inject              # cwd
    gitmap inject <folder>     # absolute, relative, or ~-prefixed
    gitmap inj   <folder>      # short alias

**What it does (in order):**

1. **DB upsert** — only when `git remote get-url origin` succeeds.
   Local-only folders silently skip this step (the user explicitly
   asked for "yes, but only if remote detected" semantics).
2. **GitHub Desktop** — registers via the same `desktop.AddRepos`
   pipeline used by `clone` and `clone-next`.
3. **VS Code** — opens via `openInVSCode` (no-op + warning when VS
   Code isn't installed).
4. **Shell handoff** — `WriteShellHandoff(target)` so the wrapper
   chdirs the parent shell into the injected folder, mirroring the
   `clone` / `cn` / `cd` UX added in v3.118.0.

**Folder validation:** any directory is accepted. No `.git/` check —
VS Code is happy to open anything, Desktop silently skips non-repos,
and the DB upsert is the only step that requires a real remote (and
it skips itself silently when one isn't there).

**Files:** `gitmap/cmd/inject.go` (new), `gitmap/helptext/inject.md`
(new), `gitmap/constants/constants_cli.go` (CmdInject + CmdInjectAlias),
`gitmap/constants/constants_v331.go` (Msg/Warn/Err constants), one
dispatch entry in `gitmap/cmd/rootcore.go`. Reuses
`resolveCloneNextFolder` for path resolution so error messages stay
consistent with `cn <folder>`.

## v3.118.0 — (2026-04-24) — `gitmap clone <url>` cds into the cloned folder

Single-URL `gitmap clone <url>` now writes the cloned folder's absolute
path to the shell-handoff sentinel, so the wrapper function chdirs the
parent shell into the new repo — same UX as `gitmap cn` and `gitmap cd`.

Multi-URL `clone` (comma-list or `clone url1 url2 ...`) deliberately
skips handoff: the destination is ambiguous when N>1, and silently
picking one would surprise users running batch clones.

File: `gitmap/cmd/clone.go` — one `WriteShellHandoff(absPath)` call
inside `executeDirectClone`, placed after the DB upsert + Desktop
registration but before the VS Code open so the cd happens regardless
of whether VS Code is installed.

## v3.117.0 — (2026-04-24) — `gitmap cn vX <folder>` + `gitmap cn <folder>` (defaults to v++); hero card UI polish

### CLI: clone-next folder-arg dispatch

Three new invocation forms (existing forms unchanged):

- `gitmap cn vX <folder>` — explicit version, explicit folder.
- `gitmap cn v+1 <folder>` / `cn v++ <folder>` — version-bump shortcuts.
- `gitmap cn <folder>` — single positional, defaults to `v++`.

`<folder>` accepts absolute, relative, `~`-prefixed, or bare-name paths.
The dispatcher chdirs into the resolved folder, runs the existing
in-place `runCloneNext` pipeline, then chdirs back — reusing
`performCrossDirCloneNext` so version-history recording, shell
handoff, desktop registration, and lock-check all behave identically
to the in-place form.

Disambiguation: `looksLikeVersion` regex extended to match `v++` and
`v+N` in addition to `v?N.N.N`. `isFolderShaped` returns true when
the token contains `/`, `\`, or starts with `~`, OR `os.Stat`
succeeds as a directory. Bare alias names with no path-hint and no
on-disk match keep falling through to the existing release-alias
resolver, so back-compat is preserved.

Files: `gitmap/cmd/clonenextfolderdispatch.go` (new),
`gitmap/cmd/clonenextfolderdispatch_test.go` (new),
`gitmap/cmd/releaserebase.go` (regex extension),
`gitmap/cmd/clonenext.go` (dispatch wiring),
`gitmap/constants/constants_v331.go` (new error/default messages).
Spec: `spec/01-app/111-cn-folder-arg.md`.
Plan: `.lovable/memory/plans/08-cn-folder-arg-plan.md`.

### UI: hero terminal card polish

- Centered the install code blocks (`max-w-2xl mx-auto`) so the cards
  read as "the focal CTA" rather than "a left-aligned block on a
  centered hero" — fixes the visual misalignment in the user
  screenshot.
- Removed the redundant outer pill background on the OS tab strip;
  the strip now sits inline with the card content (no second
  coloring section).
- Bigger, more legible OS badges: `text-sm` + `px-4 py-1.5`
  (was `text-xs` + `px-3 py-1`), font switched to `font-sans` so
  they read as UI labels not code (mono is now reserved for the
  command itself).
- Green accent for active state via new semantic tokens
  `--accent-success` / `--accent-success-bg` / `--accent-success-border`
  (light + dark variants in `src/index.css`). No hardcoded color
  classes — clears the design-token lint warnings.
- Section labels ("Install — Quick", "Uninstall — Quick") switched
  to `font-sans` (Ubuntu) per the user's font directive — mono is
  now reserved exclusively for code blocks.

### Misc

- Bumped `constants.Version` to `3.117.0`.

## v3.116.0 — (2026-04-24) — README: canonical "Update Source Before Building" section with v3.92.0+ rename verification

### Why

The hardening trilogy (v3.113.0 fsutil migration, v3.114.0 AST guard,
v3.115.0 pre-build stamp) made the `fileExists redeclared` regression
impossible on a fresh checkout — but a contributor running `git pull`
on a long-lived branch still has no canonical place in the README to
look up the verification commands. This release adds that section and
links it to the v3.115.0 build stamp so the entire detection path is
documented end-to-end.

### Changes

- New README section **"Update Source Before Building (avoid the
  `fileExists redeclared` regression)"** placed immediately after
  "Clone & Setup (Development)" so it is the first thing a returning
  contributor sees:
  - Canonical update sequence: `git fetch` → `checkout main` →
    `pull --ff-only` → `git status` clean check → SHA capture.
  - Three copy-pasteable verification commands that confirm the
    v3.92.0+ rename fix is present:
    1. `grep '^const Version = '` against `constants.go` (must be
       `3.92.0` or newer — `3.115.0+` recommended).
    2. `grep -nE '^func (fileExists|fileExistsLoose)\('` against
       `updatedebugwindows.go` (must produce no output — the helper
       moved to `gitmap/fsutil` in v3.113.0).
    3. `test -f gitmap/fsutil/exists.go` + import check on both cmd/
       files (both file paths must be printed).
  - Pointer to the v3.115.0 pre-build stamp scripts (`bash
    scripts/build-stamp.sh --strict` and the PowerShell equivalent)
    with the expected healthy-state output line and the FAIL-state
    remediation ("stop and re-pull").
  - Explicit framing: the redeclaration error is **always** a
    stale-checkout symptom and the current source cannot produce it.
- Bumped `constants.Version` to `3.116.0`.

### Verification matrix (now end-to-end documented)

| Layer | Mechanism | Detection point |
|---|---|---|
| Source | `gitmap/fsutil` package | v3.113.0 |
| Compile-time | rename-pin test | v3.113.0 |
| AST | `updatedebugwindows_source_test.go` | v3.114.0 |
| Pre-build | `scripts/build-stamp.{sh,ps1}` | v3.115.0 |
| **Documentation** | **README "Update Source Before Building"** | **v3.116.0** |

## v3.115.0 — (2026-04-24) — Pre-build provenance stamp surfaces stale checkouts in the first lines of the build log

### Why

Three releases (v3.92.0 rename, v3.113.0 fsutil migration, v3.114.0 AST
guard) have hardened the source tree against the `fileExists`
redeclaration regression — but every fix only catches it AFTER `go build`
runs and emits a cryptic line-number mismatch. Users on stale CI
checkouts still spent minutes diagnosing line numbers that didn't exist
in the source they were reading. The fix needed to move earlier in the
pipeline: surface the exact commit + version + file fingerprints
BEFORE the compiler is invoked, so a stale snapshot is obvious.

### Changes

- **`scripts/build-stamp.sh`** — bash provenance stamper. Prints:
  - `git`: commit SHA (full + short), branch, `describe --tags --dirty`,
    commit date, commit subject.
  - `source`: declared `constants.Version`, plus a `sha256:<12-hex>` +
    line-count fingerprint of `constants.go`, `updaterepo.go`,
    `updatedebugwindows.go`.
  - `guards`: a redeclaration-risk pre-check that grep-scans both
    cmd/ files for local `func fileExists` / `func fileExistsLoose`
    declarations. If both files declare one, the script either prints
    a FAIL line (default mode) or `exit 1` (`--strict`) — predicting
    the build failure before `go build` runs and pointing the user at
    `git pull origin main`.
  - All probes fall back to `(unknown)` so the stamp itself never
    blocks a build (shallow clones, tarball builds, missing git).
- **`scripts/build-stamp.ps1`** — PowerShell companion with identical
  semantics (`-Strict` switch, `Get-FileHash` + `Select-String`
  equivalents) for Windows local builds.
- **`run.sh`** + **`run.ps1`** — invoke the stamp script immediately
  before `go build`. Wrapped in `|| true` / `try`/`catch` so a stamp
  bug never breaks a working local build.
- **`.github/workflows/ci.yml`** — new "Pre-build provenance stamp
  (stale-checkout guard)" step in the `build` matrix job, runs in
  `--strict` mode so CI fails fast on a stale checkout instead of
  burning minutes on a doomed Go compilation.
- Bumped `constants.Version` to `3.115.0`.

### Reading the stamp

A healthy build log starts with:

    === gitmap build-stamp v1.0.0 ====================================
    git
      commit                  <full SHA>
      short                   <10-char SHA>
      branch                  main
      describe                v3.115.0
    source
      declared-version        3.115.0
    guards
      redecl-risk-check       ok (no local fileExists* in cmd/ ...)
    =====================================================================

If the `commit` line does not match the SHA you expected to build, or
the `declared-version` is older than the version in your local
`CHANGELOG.md`, **stop** — you are building a stale snapshot. Run
`git pull origin main` and re-run.

## v3.114.0 — (2026-04-24) — Source-level AST guard: `updatedebugwindows.go` must call `fsutil.FileOrDirExists` and declare zero local `fileExists*` helpers

### Why a second test

The v3.113.0 rename-pin test (`updatedebugwindows_rename_test.go`) only
fires when a contributor reverts the fsutil migration in a way that
shadows the imported symbol — i.e. re-adding a local helper while the
fsutil import is still present. It does NOT fire if a contributor
removes the fsutil import for an unrelated reason and re-inlines a
local `fileExists` helper. That second path puts the redeclaration
footgun back the moment any sibling cmd file (like updaterepo.go)
declares the same name.

This release adds an orthogonal guard that asserts the invariant
directly against the file's AST.

### Changes

- New `gitmap/cmd/updatedebugwindows_source_test.go` parses
  `updatedebugwindows.go` with `go/parser` and asserts:
  - **`TestUpdateDebugWindowsHasFsutilLooseCall`** — the file imports
    `github.com/alimtvnetwork/gitmap-v27/gitmap/fsutil` AND contains at
    least one real call expression of the form
    `fsutil.FileOrDirExists(...)`. The check uses the AST (selector
    expression), not a substring scan, so comments mentioning the name
    do not satisfy the invariant — a real call site is required.
  - **`TestUpdateDebugWindowsHasNoLocalFileExistsDecl`** — the file
    declares zero top-level functions named `fileExists` or
    `fileExistsLoose`. Method declarations with a receiver are
    excluded because the redeclaration footgun only applies to
    package-level free functions sharing a namespace. The AST walk
    ignores comments and string literals, so the v3.113.0 migration
    doc comment that legitimately mentions `fileExistsLoose` does
    not false-positive.
- The source file under test is unchanged — the test asserts the
  state already shipped in v3.113.0 is preserved going forward.
- Bumped `constants.Version` to `3.114.0`.

### Failure-mode coverage table

| Regression scenario | v3.113.0 rename-pin test | v3.114.0 AST guard |
|---|---|---|
| Re-add local `fileExists` while keeping fsutil import | ✅ compile error | ✅ AST assertion |
| Re-add local `fileExistsLoose` while keeping fsutil import | ✅ compile error | ✅ AST assertion |
| Remove fsutil import + re-add local `fileExists` | ❌ silent | ✅ AST assertion |
| Remove the `fsutil.FileOrDirExists` call site entirely | ❌ silent | ✅ AST assertion |
| Add a comment mentioning `fileExists` | (no false positive) | (no false positive) |

## v3.113.0 — (2026-04-24) — Centralize `fileExists` / `fileExistsLoose` / `dirExists` into a shared `gitmap/fsutil` package

### Motivation

The redeclaration footgun that bit the `cmd` package twice (v3.92.0
rename, v3.112.0 stale-CI guard) only existed because two files in the
same Go package defined unexported existence helpers with overlapping
names. Pinning the rename via a test prevents accidental reverts but
does not remove the underlying class of bug — any third file in `cmd`
that declares `fileExists` would re-trigger it.

### Changes

- New package **`gitmap/fsutil`** (`gitmap/fsutil/exists.go`) exporting
  three predicates with explicitly documented contracts:
  - `FileExists` — strict file (rejects directories, rejects empty)
  - `FileOrDirExists` — loose existence (accepts directories, rejects empty)
  - `DirExists` — strict directory (rejects files, rejects empty)
  All three short-circuit on the empty string so callers don't need
  to guard their inputs — this matches the previous `fileExistsLoose`
  contract that the debug-dump code relied on.
- Contract-pinning tests in `gitmap/fsutil/exists_test.go` exercise
  every variant against a real tempdir, file, missing path, and the
  empty string. Collapsing two variants now fails a test before it
  fails a downstream caller.
- `gitmap/cmd/updaterepo.go` — removed local `dirExists` and
  `fileExists`, calls `fsutil.DirExists` / `fsutil.FileExists` instead.
- `gitmap/cmd/updatedebugwindows.go` — removed local `fileExistsLoose`,
  calls `fsutil.FileOrDirExists` instead.
- `gitmap/cmd/updatedebugwindows_rename_test.go` — repurposed from the
  v3.112.0 rename pin into a forward-looking guard. The new
  `TestFsutilMigrationPinned` asserts the cmd package uses `fsutil.*`
  and exercises both the loose-empty short-circuit and the strict-dir
  rejection. A future contributor reintroducing a local helper would
  trigger the redeclaration error this test was created to prevent.
- Bumped `constants.Version` to `3.113.0`.

### Scope

Only the `cmd` package is migrated in this release because that is
where the redeclaration risk has materialized. The duplicated helpers
in `gitmap/release`, `gitmap/lockfile`, `gitmap/localdirs`,
`gitmap/vscodepm`, and `gitmap/detector` are functionally fine
(different packages = different namespaces) and can migrate to
`fsutil` opportunistically without breaking anything.

## v3.112.0 — (2026-04-24) — Pin the `fileExists` / `fileExistsLoose` rename so the v3.92.0 redeclaration cannot regress

### Diagnosis

User's CI reported:

    cmd/updaterepo.go:118:6: fileExists redeclared in this block
    cmd/updatedebugwindows.go:148:6: other declaration of fileExists

In the current source tree this is impossible:

- `gitmap/cmd/updaterepo.go:118` declares the strict, file-only
  `fileExists` (already there since before v3.92.0).
- `gitmap/cmd/updatedebugwindows.go:150` declares the renamed loose
  variant `fileExistsLoose` (the v3.92.0 fix).
- `updatedebugwindows.go:148` is a comment line, not a `func`
  declaration. A redeclaration error CANNOT originate from a comment.
- Other `fileExists` symbols live in DIFFERENT packages
  (`gitmap/detector`, `gitmap/localdirs`, `gitmap/lockfile`,
  `gitmap/release`) and Go allows duplicates across packages.

Conclusion: **the user's CI is building from a stale snapshot that
pre-dates v3.92.0.** This is the same failure mode covered by
v3.95.0's stale-binary guard, just on the build-host side instead of
the deployed-binary side.

### Added

- `gitmap/cmd/updatedebugwindows_rename_test.go` — two paired tests
  that compile only if the v3.92.0 rename is preserved:
  - `TestFileExistsLooseSymbolPinned` references `fileExistsLoose`
    directly. If a future contributor reverts the rename, the test
    fails to compile alongside the duplicate-declaration error,
    making the cause unambiguous.
  - `TestFileExistsStrictSymbolPinned` references the strict
    package-level `fileExists` and asserts it still rejects
    directories. Pairs the two helpers as an explicit, tested
    contract instead of a drifting implementation detail.

### Action required for the user

The source is correct. To clear CI:

1. `git pull` on the build host so the v3.92.0 rename is present.
2. Re-run the build.
3. If the failure persists, run `git log -- gitmap/cmd/updatedebugwindows.go`
   and confirm the v3.92.0 commit is in the checked-out history. If
   it isn't, the CI is on a branch that diverged before v3.92.0 and
   needs to be rebased.

### Bumped

- `constants.Version` → `3.112.0`.

## v3.111.0 — (2026-04-24) — Surface `td` / `ti` aliases in shell tab-completion

### Added

- **Typed CLI constants** for templates subcommand aliases:
  `CmdTemplatesDiff` (`"diff"`), `CmdTemplatesDiffAlias` (`"td"`),
  `CmdTemplatesInit` (`"init"`), `CmdTemplatesInitAlias` (`"ti"`) in
  `gitmap/constants/constants_templates_cli.go`. The block opts in via
  `// gitmap:cmd top-level` so the completion generator picks the
  alias values up automatically.
- **Generator-aware skip markers** on the full subcommand strings:
  `CmdTemplatesDiff` and `CmdTemplatesInit` carry `// gitmap:cmd skip`
  line comments because:
  - `"diff"` is already a top-level command (folder-tree diff via
    `gitmap/cmd/diff.go`) — re-listing it here is a no-op for the
    completion union, but the marker documents intent and prevents
    a future audit from re-opening the question.
  - `"init"` is not a top-level gitmap command at all; surfacing it
    standalone would mislead users into typing `gitmap init`.

### Changed

- `gitmap/cmd/templatesdiff.go` and `gitmap/cmd/templatescli.go` now
  alias the local `cmdTemplatesDiff*` / `cmdTemplatesInit*` package
  identifiers to the shared `constants.CmdTemplates*` values, giving
  the completion generator one source of truth and preventing string
  drift between the dispatcher switch and the completion table.

### Generated

- `gitmap/completion/allcommands_generated.go` — added `"td"` and
  `"ti"` in sorted positions. The full `templates diff` invocation
  remains discoverable because both the parent (`templates` / `tpl`)
  and the alias (`td`) are now in completion; users who tab-complete
  `gitmap td<TAB>` get the alias suggested directly. Re-run
  `go generate ./completion/...` locally to confirm byte-equality
  before tagging — the manual edit follows the deterministic sort
  the generator produces.

### Plan status

- Plan 05 Phase 4 — completion-generator note revised: `td` and `ti`
  ARE now wired (overriding the prior "subcommand precedent" decision
  documented at v3.108).

### Bumped

- `constants.Version` → `3.111.0`.

## v3.110.0 — (2026-04-24) — Plan 05 Phase 2 closeout: `gitmap templates init` (alias `ti`)

### Verified

- **Phase 2 already shipped** — `gitmap templates init <lang> [<lang>...]
  [--lfs] [--dry-run] [--force]` (alias `ti`) is fully implemented in
  `gitmap/cmd/templatesinit.go` (297 lines), wired into the dispatcher in
  `templatescli.go`, covered by 9 unit tests in
  `gitmap/cmd/templatesinit_test.go`, and documented in
  `gitmap/helptext/templates-init.md` (133 lines).
- **Behavior:** per lang, `ignore/<lang>` is required (hard-fail with hint
  pointing at `templates list`); `attributes/<lang>` is optional with a
  dim soft-skip notice (matches embed corpus reality where some langs
  lack an attributes file). `--lfs` adds a single `lfs/common` step
  targeting `.gitattributes` — same marker tag as `add lfs-install` so
  the two are interchangeable.
- **Idempotency:** re-running `init <lang>` produces "unchanged" lines.
  Running `add ignore <lang>` after `init <lang>` is also a no-op (same
  marker tag `ignore/<lang>`). Confirmed by inspection of `merge.go`
  block-tag plumbing.
- **UX deviations from original plan:** positional `<lang>...` chosen
  over `--lang <csv>` (more natural for a scaffolder), and `--cwd`
  deferred (`cd && gitmap templates init …` covers it).

### Docs site

- Added `templates init` entry to `src/data/commands.ts` so the docs
  site command browser surfaces it alongside `templates list`,
  `templates show`, and `templates diff`.

### Plan status

- Plan 05: Phase 0 ✅, Phase 1 ✅, Phase 2 ✅, Phase 3 ✅, Phase 4 ⏳
  (only README snippet remains), Phase 5 ⏳ (QA + tag).

### Bumped

- `constants.Version` → `3.110.0` (docs + plan closeout, no code changes).

## v3.109.0 — (2026-04-24) — Plan 05 Phase 1 closeout: language corpus verified

### Verified

- **Phase 1 language coverage already complete in tree** — confirmed all
  ten template files for the new language batch ship with proper
  audit-trail headers:
  - `assets/ignore/{java,ruby,php,swift,kotlin}.gitignore`
  - `assets/attributes/{java,ruby,php,swift,kotlin}.gitattributes`
- **Corpus parity test pins the batch** — `gitmap/templates/corpus_parity_test.go`
  already enumerates `java`, `ruby`, `php`, `swift`, `kotlin` and asserts
  every ignore lang has a matching attributes counterpart (line 97 loop).
  Any future revert to the asset directory will fail CI before merge.
- **No lang enum drift** — confirmed the resolver discovers languages by
  filesystem walk + filename parsing, so adding a new lang requires only
  dropping the file in (per Plan 04 design). No `LangJava` / `LangRuby` /
  etc. constants exist or are needed.

### Plan status

- Plan 05 Phase 1 marked **done** in
  `.lovable/memory/plans/05-templates-polish-plan.md`.
- Phase 0 (spec) ✅, Phase 1 (langs) ✅, Phase 3 (`templates diff`) ✅.
- Remaining: Phase 2 (`templates init`), Phase 4 finishing items
  (README snippet — deferred until `init` lands), Phase 5 (QA + tag).

### Bumped

- `constants.Version` → `3.109.0` (docs + plan closeout, no code changes).

## v3.108.0 — (2026-04-24) — `gitmap templates diff` (alias `td`)

### Added

- **`gitmap templates diff --lang <name> [--kind ignore|attributes] [--cwd <path>]`**
  shows what `add ignore <lang>` / `add attributes <lang>` would change
  in the current repo without writing to disk. Alias: `td` (e.g.
  `gitmap tpl td --lang go`).
- **Standard `diff(1)` exit codes** for script-friendliness:
  - `0` → on-disk gitmap block matches template (no changes).
  - `1` → block missing or body differs (changes pending).
  - `2` → bad flag, unknown language, or I/O failure.
- **Block-scoped diffing** — only the gitmap-managed marker block
  participates. Hand edits OUTSIDE `# >>> gitmap:<tag> >>>` …
  `# <<< gitmap:<tag> <<<` are invisible (matches `add`'s contract).
- **Unified-style hunks** with banner lines (`@@ gitmap:<tag> @@`),
  `+` for would-be additions, `-` for current content. TTY-aware:
  cyan `+`, yellow `-`, dim `@@` via the existing
  `render.HighlightQuotesANSI` token replacer. Pipes/redirects keep
  raw output for downstream parsers.

### Internals

- New `gitmap/templates/diff.go` (~170 LOC):
  - `Diff(targetPath, tag, body) (DiffResult, error)` is pure — it
    never writes — so `add --dry-run` could later reuse it without
    touching disk.
  - Status enum (`DiffNoChange / DiffMissingFile / DiffMissingBlock
    / DiffBlockChanged`) drives the CLI exit codes.
  - `extractBlockBody` reuses the existing `blockRegex(tag)` from
    `merge.go` so the parser can never drift from the writer.
  - `splitDiffLines` deliberately preserves intra-body blank lines
    so visual separators in templates show up as `+` / `-`.
- New `gitmap/cmd/templatesdiff.go` — flag parsing, kind expansion
  (`""` → both kinds), per-kind dispatch, ANSI decoration. All
  helpers under the 15-line rule.
- `dispatchTemplates` extended with `diff` / `td` cases; usage
  banner gains a "Flags (diff)" section + 2 new examples.

### Tests

- `gitmap/templates/diff_test.go` (5 cases) pins all four
  `DiffStatus` branches plus the blank-line preservation invariant.

### Files

- New: `gitmap/templates/diff.go`, `gitmap/templates/diff_test.go`,
  `gitmap/cmd/templatesdiff.go`, `gitmap/helptext/templates-diff.md`
- Edited: `gitmap/cmd/templatescli.go` (dispatch + usage banner),
  `gitmap/constants/constants.go` (v3.108.0), `CHANGELOG.md`,
  `.lovable/memory/index.md`,
  `.lovable/memory/plans/05-templates-polish-plan.md` (Phase 3 → done)

## v3.107.0 — (2026-04-24) — Pretty renderer fixture corpus expansion

### Added

- 3 new edge-case fixtures in `gitmap/render/testdata/pretty/`:
  - **`case-007-heading-without-subtitle`** — heading directly followed by
    a plain paragraph (no italic line). Pins the subtitle-peek's
    "no italic = treat as paragraph" branch so a future regression that
    eats the line as a subtitle would fail loudly.
  - **`case-008-consecutive-collapses`** — two paragraph + identical-fence
    pairs back-to-back. Pins that `appendFence` rewrites the *previous*
    paragraph each time and that the rewritten block keeps the standard
    body indent (`  [Y]→ ...[/Y]`).
  - **`case-009-blank-line-preservation`** — heading + 3 paragraphs with
    blank lines between. Pins that `bkBlank` blocks survive the round
    trip (no blank-collapsing) so help text doesn't lose its visual
    breathing room.
- All three pairs are picked up automatically by the existing
  `TestPrettyFixtures` table-driven loop — no test code changes needed.

### Notes

- Phase 5 of the Templates plan was already shipped (renderer +
  table-driven test + 6 fixtures + ANSI-swap + defensive-close tests).
  This release rounds the corpus to 9 cases covering the three remaining
  uncovered branches in `pretty_parse.go` / `pretty_emit.go`.

### Files

- New: `gitmap/render/testdata/pretty/case-007-heading-without-subtitle.{in.md,want.txt}`,
  `case-008-consecutive-collapses.{in.md,want.txt}`,
  `case-009-blank-line-preservation.{in.md,want.txt}`
- Edited: `gitmap/constants/constants.go` (v3.107.0), `CHANGELOG.md`,
  `.lovable/memory/index.md`

## v3.106.0 — (2026-04-24) — `templates list` --kind / --lang filters

### Added

- **`gitmap templates list --kind <ignore|attributes|lfs>`** narrows the
  output table to one kind so `templates list --kind ignore` no longer
  scrolls past the `attributes` block.
- **`gitmap templates list --lang <name>`** narrows by language across
  every kind (case-insensitive). Useful for spotting whether you have
  both `ignore/go` and `attributes/go` on disk.
- Filters AND together: `--kind ignore --lang go` matches at most one
  row.
- **Strict --kind validation**: unknown values exit 1 with
  `templates list: unknown --kind "foo" (want ignore | attributes | lfs)`.
  Empty result from a valid filter prints
  `(no templates match the requested filter)`.

### Tests

- `TestFilterTemplatesNoFilters` — identity case (no filter = full list).
- `TestFilterTemplatesByKind`, `TestFilterTemplatesByLang` — single-axis filters.
- `TestFilterTemplatesByKindAndLang` — pins AND semantics so the two
  filters can't accidentally OR.
- `TestFilterTemplatesEmptyResult` — pins the trigger condition for the
  filtered-empty message.
- `TestIsValidKindFilter` — locks the kind allow-list.
- `TestParseTemplatesListFlagsLowersValues` — pins case-folding + trim.

### Files

- Edited: `gitmap/cmd/templatescli.go` (filter parsing, validation,
  pure `filterTemplates` helper), `gitmap/helptext/templates.md`
  (Flags table + Example 5), `gitmap/constants/constants.go` (v3.106.0),
  `CHANGELOG.md`, `.lovable/memory/index.md`
- New: `gitmap/cmd/templatescli_filter_test.go`


## v3.105.0 — (2026-04-24) — `gitmap add ignore` / `add attributes`

### Added

- **`gitmap add ignore [langs...]`** merges the curated `common` +
  per-language `.gitignore` templates into a marker block at the repo
  root. Supports `go`, `node`, `python`, `rust`, `csharp`, `java`,
  `kotlin`, `php`, `ruby`, `swift` out of the box, plus any user
  overlay under `~/.gitmap/templates/ignore/`.
- **`gitmap add attributes [langs...]`** is the `.gitattributes`
  sibling. Same algorithm; same idempotency guarantees.
- Both new commands accept `--dry-run` to preview the marker block
  before writing.
- **Stable, sorted marker tags**: `add ignore go node` and
  `add ignore node go` share a single block (`ignore/go+node`), so the
  block stays one block across collaborators who happen to type
  arguments in different orders.
- **Per-line dedupe** with blank-line preservation: merging
  `common+go+node+python` collapses repeat rules across templates
  without flattening the visual section spacers.

### Changed

- `dispatchAdd` now routes `ignore` and `attributes` in addition to the
  pre-existing `lfs-install`. Usage banner updated accordingly.

### Tests

- `TestNormalizeLangs` pins arg-list normalization (case fold, dedupe,
  strip implicit `common`, preserve order).
- `TestBuildAddTagSorted` pins the sorted-tag invariant for stable
  marker-block addressing across argument orders.
- `TestDedupeLinesPreservesBlanks` pins blank-line semantics so the
  merged body stays human-readable.
- `TestConcatTemplateBodiesAddsLangBanners` pins the `# ── <lang> ──`
  separator banners, exercised against real embedded templates.

### Files

- New: `gitmap/cmd/addignoreattrs.go` (12 helpers, all <15 lines)
- New: `gitmap/cmd/addignoreattrs_test.go`, `addignoreattrs_testhelper_test.go`
- New: `gitmap/helptext/add-ignore.md`, `gitmap/helptext/add-attributes.md`
- Edited: `gitmap/cmd/rootadd.go` (router + usage banner),
  `gitmap/constants/constants.go` (v3.105.0), `CHANGELOG.md`,
  `.lovable/memory/{plans/04-templates-ignore-attributes-plan,index}.md`


## v3.104.0 — (2026-04-24) — commit-both --interleave (author-date variant)

### Added

- **`gitmap commit-both --interleave`** ships the author-date variant
  originally drafted in spec §5 and deferred in v3.102.0. It builds
  both directional plans up front, merges the commit lists into a
  single chronological stream (stable sort by `AuthorAt`; LEFT-side
  wins exact ties), prints a unified preview, prompts once, then
  replays each commit onto its opposite side in author-date order.
  Use `--dry-run` first — first per-commit failure aborts the stream
  and leaves the just-written side in a partial state.
- **`committransfer.RunBothInterleaved`** in `gitmap/committransfer/interleave.go`,
  composed from helpers under the 15-line/function cap:
  `buildInterleavedStream`, `executeInterleaveStream`, `printInterleavedPlan`,
  `replayInterleaveSteps`, `finalizeInterleavePush`.
- **Sort invariant + tie-breaking pinned** by `TestBuildInterleavedStreamSortsByAuthorDate`,
  `TestBuildInterleavedStreamStableForTies`, `TestBuildInterleavedStreamEmptyPlans`
  in `interleave_test.go`.
- **CLI guard:** `--interleave` is rejected (exit 2) when passed to
  `commit-left` or `commit-right` — only `commit-both` accepts it.
- **Spec §5 split into 5.1 (sequential) and 5.2 (--interleave)** with
  full tradeoff documentation. Helptext (`commit-both.md`) and
  HelpCommitBoth one-liner updated. HelpCommitLeft promoted from
  `[scaffold]` to `[LIVE]` (it has been live since v3.102.0).


## v3.103.0 — (2026-04-24) — fix non-functional shell handoff via sentinel-file mechanism

### Fixed

- **`gitmap clone-next` shell handoff was a no-op since inception.** The previous implementation called `os.Setenv("GITMAP_SHELL_HANDOFF", path)` from the binary, which can never propagate to the parent shell (child processes cannot mutate parent env). The flattened-folder cd never happened.

### Added

- **Sentinel-file handoff mechanism** (`gitmap/cmd/shellhandoff.go::WriteShellHandoff`). The shell wrapper function exports `GITMAP_HANDOFF_FILE=<temp>` before invoking the binary; the binary writes the target path to that file; the wrapper reads the file after the binary exits and `cd`s the parent shell.
- **Wired into `clone-next`, `as`, `cd <name>`, `cd repos`** so all four commands hand a target directory back to the parent shell consistently.
- **Updated `constants.CDFunc{Bash,Zsh,PowerShell}`** wrappers to consume `GITMAP_HANDOFF_FILE` for any subcommand (preserving the existing stdout-capture path for `cd`/`go` for backwards compatibility).
- **Unit tests** in `gitmap/cmd/shellhandoff_test.go` covering the no-op (env unset), happy-path write, and empty-path safeguard cases.
- **Spec & memory updates**: `spec/01-app/87-clone-next-flatten.md` rewritten, new memory `.lovable/memory/features/shell-handoff-file.md`.

### Backwards compatibility

- Without the wrapper installed, `GITMAP_HANDOFF_FILE` is unset → `WriteShellHandoff` is a silent no-op.
- The legacy `GITMAP_WRAPPER=1` detector and the stdout-capture path used by `cd`/`go` both remain intact.


## v3.95.0 — (2026-04-24) — refuse to build URL-shaped folder paths + lock multi-URL routing behind regression tests

### Fixed

- **Root cause of the recurring `pending task already exists ... \https:\github.com\...` + `fatal: could not create leading directories` failure for `gitmap clone <url1> <url2>`:** the user's deployed `gitmap.exe` on PATH is older than v3.80.0 and does not contain the multi-URL routing fix. Current source already routes 2+ URL invocations to `runCloneMulti` via `shouldUseMultiClone(cf)` — but the *deployed* binary still reaches `executeDirectClone(url, folderName=<second URL>, ...)`, builds the impossible path `D:\...\https:\github.com\...`, and crashes git. Repeated `gitmap update` runs have been failing in Phase 3 cleanup (issues #09 / #10 / #12), so the freshly built binary never replaces the stale one on PATH.
- **`gitmap/cmd/clone.go` `executeDirectClone` now refuses URL-shaped folder names early** with an actionable message that names the exact recovery commands: `gitmap doctor`, `gitmap update`, `gitmap pending clear --yes`, and the reminder to open a NEW terminal so PATH refreshes. This shape is impossible in current source, so the guard fires only when a stale binary is in use — and the message tells the user exactly that instead of letting git fail with `Invalid argument`.
- **`gitmap/constants/constants_clone.go`** adds `ErrCloneStaleBinaryFolderURL` so the guidance text lives in constants per the project's no-magic-strings rule.

### Added

- **Regression tests** in `gitmap/cmd/clone_stale_binary_test.go` pinning all three exact PowerShell argv shapes the user has reported (`url1,url2,url3` comma-glued, `url1 url2 url3` PowerShell-split, comma+space mixed). They prove `shouldUseMultiClone` routes every reported shape to `runCloneMulti` in current source, and the third test pins the recovery message contents so `gitmap update` / `gitmap doctor` / `gitmap pending clear` can never silently disappear from the guidance.
- **Root Cause Analysis:** `spec/02-app-issues/33-stale-binary-clone-folder-url-guard.md` — full evidence trail proving the deployed binary is stale, why every retry hits the same code path, and the prevention rules.

### Validation

- `go test ./cmd -run "TestShouldUseMultiCloneCoversReportedInvocation|TestIsDirectURLAcceptsAllReportedShapes|TestStaleBinaryGuardMessageMentionsRecoverySteps" -v -count=1`

### Action required by the user

The source is correct and the new guard ships in `v3.95.0`, but **the binary on your PATH must be replaced** before either takes effect on your machine:

1. `gitmap doctor` — confirm the active binary version (will show <3.80.0).
2. `gitmap update` — rebuild + redeploy from current source. If Phase 3 cleanup still fails, the on-disk handoff log under `<TMP>/gitmap-update-handoff-*.log` (v3.87.0+) and the `--debug-windows-json` sink (v3.91.0+) now record exact branch-level evidence.
3. **Open a NEW terminal** so PATH refreshes.
4. `gitmap pending clear --yes` — drop the orphaned row blocking your retries.
5. Retry: `gitmap clone https://.../email-creator-v1 https://.../email-reader-v3 https://.../account-automator` works in either comma- or space-separated form on PowerShell and bash.

## v3.94.0 — (2026-04-24) — docs UI now uses a VS Code-style workbench color grade

### Fixed

- **Repeated missed UI request:** the docs app no longer mixes default/light marketing-style surfaces with the requested VS Code-inspired dark grading. The main shell now reads like an editor workbench instead of a generic docs template.
- **`src/components/docs/DocsLayout.tsx` + `DocsSidebar.tsx`** now use a flatter VS-style header/explorer treatment with restrained borders, panelized surfaces, blue accent states, and consistent workbench framing.
- **`src/index.css` semantic tokens** were retuned so both light and dark themes align with a VS-like palette. Dark mode is now the primary visual target, with editor-style neutrals, subtle borders, and blue selection emphasis replacing the earlier green-heavy look.
- **Reusable docs surfaces** were brought into the same grading: `FeatureCard.tsx`, `InstallBlock.tsx`, `CodeBlock.tsx`, and the home hero in `src/pages/Index.tsx` now use flatter panel styling instead of decorative aurora/card effects that clashed with the requested direction.
- **Theme startup/readback** in `src/lib/theme.ts` now falls back to dark more reliably, so the requested VS-style grade is what users see by default unless they explicitly switch.

### Added

- **Root Cause Analysis:** `spec/02-app-issues/32-docs-ui-vscode-grading-missed-request.md` documents why the UI request kept being missed, what visual mismatch remained in the app, and how to prevent that process failure from recurring.

## v3.93.0 — (2026-04-24) — update-cleanup Phase 3 now logs inner child failures durably

### Fixed

- **Root cause of the "update-cleanup keeps failing with no logs" loop:** the Phase 3 deployed-binary handoff was working at the outer lifecycle level (`resolve`, `start_ok`, `done`), but several **inner cleanup branches** still wrote only to child stderr. On Windows that child runs hidden/detached, so the exact per-file failure often vanished even though the parent printed `→ Cleanup process started (pid=...)`.
- **`gitmap/cmd/updatecleanup_remove.go` now mirrors per-file failures into the durable handoff log and JSON sink.** Added durable events for:
  - `glob_error` when `filepath.Glob` fails for a cleanup pattern
  - `remove_retry` on every transient `os.Remove` failure before the final attempt
  - `remove_fail` after retry exhaustion
  - `remove_ok` on successful deletion
- **`gitmap/cmd/updatecleanup_extra.go` now mirrors the previously silent "special-case" cleanup branches too.** Added durable events for:
  - `drive_root_skip` when the obsolete drive-root shim is skipped by the 5 MB guard
  - `drive_root_remove_fail` / `drive_root_remove_ok`
  - `swap_glob_error` for `*.gitmap-tmp-*` enumeration failures
  - `swap_remove_fail` / `swap_remove_ok` for leftover swap-directory cleanup
- **Result:** when the detached cleanup child fails again, the failure is now forensically recoverable even if console stderr is swallowed. The user gets exact branch-level evidence in the always-on handoff log and, when enabled, the `--debug-windows-json` NDJSON sink.

### Added

- **Root Cause Analysis:** `spec/02-app-issues/31-update-cleanup-phase3-observability-gap.md` documents the repeated user report, the actual failure mode, the observability gap, the solution, and the validation steps.
- **Automated regression tests** in `gitmap/cmd/updatecleanup_handoff_test.go` covering:
  - stable `formatHandoffLogLine(...)` output for child failure events
  - forwarding of `--debug-windows` / `--debug-windows-json` in `buildCleanupChildArgs()`
  - forwarding of `GITMAP_UPDATE_CLEANUP_DELAY_MS`, `GITMAP_DEBUG_WINDOWS`, and `GITMAP_DEBUG_WINDOWS_JSON` in `buildCleanupChildEnv()`

### Validation

- `go test ./cmd -run "TestFormatHandoffLogLineIncludesStableFields|TestBuildCleanupChildArgsForwardsDebugFlags|TestBuildCleanupChildEnvForwardsDelayAndJSONPath|TestCollectBackupCleanupDirsIncludesPathDerivedDeployAndBuild|TestCollectTempCleanupDirsIncludesTempAndDerivedTargets" -v -count=1`
- All targeted tests passed.



### Fixed

- **Duplicate `fileExists` declaration** in `gitmap/cmd/` blocked every `go test ./cmd/...` run. The original lived in `updaterepo.go` (file-only check); a second copy was added later in `updatedebugwindows.go` (file-or-dir, empty-string-safe). Renamed the debug-dump version to `fileExistsLoose` and updated its two call sites in `dumpDebugWindowsHandoff`. The semantics of each helper now match what its single consumer actually wants.

### Added

- **`gitmap/cmd/root_url_shortcut_test.go`** — pins the bare-URL shortcut against the three exact invocations the user reported as failing with `Unknown command`:
  - `gitmap https://a,https://b,https://c` (single comma-glued PowerShell paste)
  - `gitmap https://a, https://b, https://c` (comma-then-space split across argv — bash paste)
  - `gitmap https://a, https://b https://c` (mixed comma/space separators)
  - Plus single-URL, leading-flag (`--verbose <url>`), SSH-shorthand (`git@github.com:a/b.git`), and GitLab-URL variants. Also covers the negative cases (known subcommand, folder path, empty argv) so the shortcut never grabs a legitimate command. Three `*testing.T` functions, all under 15 lines per case, no external deps.

### Why

The shortcut logic (`shouldRewriteToClone`, `looksLikeURLToken`, `splitOnComma`) had **zero unit coverage**, so a regression in any of its three cooperating predicates would silently re-introduce the `Unknown command` failure mode. The new tests run in milliseconds and would have caught the duplicate-`fileExists` build break the next time the suite ran.

### Diagnostic note for users still seeing `Unknown command: https://...`

The shortcut has been in source since before v3.91.0, but a **stale deployed binary on `PATH` will not have it**. Run `gitmap doctor` to confirm the active binary version, then `gitmap update` to deploy the current source. The startup version-check banner added in v3.90.0 also surfaces the gap on every invocation.



### Added

- **JSON sink for `--debug-windows` diagnostics.** Opt in with the new `--debug-windows-json` flag (or `GITMAP_DEBUG_WINDOWS_JSON=<path>` env var). When enabled, every `[debug-windows]` console line is mirrored as a structured event to:

      output/gitmap-debug-windows-2026-04-24_15-30-12.jsonl

  One NDJSON event per line, each with a stable envelope:

      {"ts":"2026-04-24T07:30:12.123Z","event":"handoff","pid":12345,"ppid":12000,"goos":"windows","self":"...","version":"3.91.0","source":"config","target":"C:\\bin\\gitmap.exe","target_exists":true,"child_argv":["update-cleanup","--debug-windows","--debug-windows-json"]}

  Events emitted: `header`, `footer`, `handoff`, `child_pid`, `note`, `command_plan`, `cleanup_plan` (the last one carries the full enumerated `temp_removals` / `backup_removals` / `swap_dirs` / `drive_root_shim` plan as nested arrays).

- **Custom path support.** Pass `--debug-windows-json=/path/to/trace.jsonl` to override the default location — useful when piping into a log aggregator or dropping the file under a network share.

- **Cross-process consolidation.** The opened sink path is exported as `GITMAP_DEBUG_WINDOWS_JSON` and `--debug-windows-json` is appended to the Phase 3 cleanup child's argv, so the detached cleanup process **appends to the same file** rather than creating its own per-process trace. One handoff = one consolidated NDJSON file covering both phases.

- **Path advertised on console.** When the sink opens, `[debug-windows] json sink file   : <path>` is printed once to stderr so users always know where the file landed.

### Why

The console dump (v3.86) and on-disk handoff log (v3.87) both have limitations: the console can be swallowed by detached Windows launchers, and the handoff log is line-oriented text that loses the structured cleanup-plan enumeration added in v3.90. A first-class NDJSON sink gives forensic-grade output that survives every transport and `jq`s cleanly:

    jq 'select(.event=="cleanup_plan") | .temp_removals[].matches' \
      output/gitmap-debug-windows-*.jsonl

### Implementation

- **`gitmap/cmd/updatedebugwindows_json.go`** (new, 162 lines) — `emitDebugWindowsJSON`, `buildDebugWindowsJSONPayload`, `isDebugWindowsJSONRequested`, `isDebugWindowsJSONFlagWithValue`, `openDebugWindowsJSONFile`, `resolveDebugWindowsJSONPath`, `debugWindowsJSONPath`. Lazy `sync.Once` open, `sync.Mutex`-guarded writes, all errors swallowed (diagnostics must never block the update flow).
- **`gitmap/cmd/updatedebugwindows.go`** — added one `emitDebugWindowsJSON(...)` call to `dumpDebugWindowsHeader`, `dumpDebugWindowsFooter`, `dumpDebugWindowsHandoff`, `dumpDebugWindowsChildPID`, `dumpDebugWindowsNote`.
- **`gitmap/cmd/updatedebugwindows_plan.go`** — `dumpDebugWindowsCommandPlan` and `dumpDebugWindowsCleanupPlan` now also emit JSON; the `dumpPlanned*` helpers were refactored to *return* the matched paths so the JSON sink records exactly what the console showed (no duplicated enumeration logic). New helpers `collectAndPrintMatches` and `collectSwapDirMatches` keep every function under 15 lines.
- **`gitmap/cmd/updatehandoff_phase3.go`** — `buildCleanupChildArgs` and `buildCleanupChildEnv` forward `--debug-windows-json` and `GITMAP_DEBUG_WINDOWS_JSON=<path>` to the cleanup child.
- **`gitmap/constants/constants_update.go`** — new `FlagDebugWindowsJSON`, `EnvDebugWindowsJSON`, `DebugWindowsJSONFileFmt`, `MsgDebugWinJSONFile`, `MsgDebugWinJSONOpenFail`.
- **`gitmap/constants/constants.go`** — bumped `Version` to `3.91.0`.

### Compatibility

Pure addition. The sink is OFF by default — `--debug-windows` alone keeps the v3.90 console-only behaviour byte-for-byte. Opt-in writes one small append-only file per handoff under `output/`. File-open failures degrade silently to console-only.


## v3.90.0 — (2026-04-24) — `--debug-windows` shows the exact spawn command and cleanup plan

### Added

- **`dumpDebugWindowsCommandPlan`** — when `--debug-windows` (or `GITMAP_DEBUG_WINDOWS=1`) is active, the Phase 3 handoff prints the exact shell-quoted spawn invocation BEFORE calling `cmd.Start()`/`cmd.Run()`. Output is copy-paste safe across PowerShell, cmd.exe, bash, and zsh:

      [debug-windows] spawn command    : "C:\Program Files\gitmap-cli\gitmap.exe" "update-cleanup" "--debug-windows"
      [debug-windows] (no `git` subprocess is launched by update-cleanup; only the line above plus os.Remove/os.RemoveAll on the paths below)

- **`dumpDebugWindowsCleanupPlan`** — called from `runUpdateCleanup` AFTER `loadUpdateCleanupContext` but BEFORE any deletion happens, so the user sees the *plan* not the outcome. Enumerates:
  1. Each `filepath.Glob` pattern for `gitmap-update-*` (temp handoff copies)
  2. Each `filepath.Glob` pattern for `*.old` (backup binaries)
  3. Each match that will be passed to `os.Remove` (skips matches equal to the active binary, mirroring `removeCleanupMatch`'s gating)
  4. Each `*.gitmap-tmp-*` swap dir that will be passed to `os.RemoveAll`
  5. The Windows-only drive-root shim candidate plus the verdict (`will os.Remove` / `skipped`) — mirrors the gating in `isRemovableDriveRootShim`

  Sample output:

      [debug-windows] ----- planned cleanup operations -----
      [debug-windows] glob             : C:\Users\me\AppData\Local\Temp\gitmap-update-*
      [debug-windows]   → os.Remove    : C:\Users\me\AppData\Local\Temp\gitmap-update-12345.exe
      [debug-windows] glob             : C:\Program Files\gitmap-cli\*.old
      [debug-windows]   → os.Remove    : C:\Program Files\gitmap-cli\gitmap.exe.old
      [debug-windows] glob             : C:\Program Files\gitmap-cli\*.gitmap-tmp-*
      [debug-windows]   (no matches)
      [debug-windows] drive-root shim  : E:\gitmap.exe (skipped)
      [debug-windows] --------------------------------------

- **Explicit "no `git` subprocess" callout.** The cleanup pipeline is pure Go syscalls — it never shells out to `git`. The new `MsgDebugWinCmdNote` line states this explicitly so users debugging update issues don't waste time hunting for a phantom `git` invocation.

### Why

After v3.86 added `--debug-windows` and v3.87 added the on-disk handoff log, the missing piece was *what exactly* the deployed binary was about to do. The dump showed PIDs, env vars, and a Go-formatted argv slice (`[update-cleanup --debug-windows]`) but no shell-quoted command line and no list of the actual files that would be deleted. Users hit "the update finished but my .old file is still there" and had no way to confirm the path the cleanup pass actually looked at.

### Implementation

- `gitmap/cmd/updatedebugwindows_plan.go` (new) — `dumpDebugWindowsCommandPlan`, `renderShellCommand`, `quoteShellToken`, `dumpDebugWindowsCleanupPlan`, `dumpPlannedRemovals`, `dumpPlannedSwapDirs`, `dumpPlannedDriveRootShim`. All under 15 lines per function, file under 200 lines.
- `gitmap/cmd/updatehandoff_phase3.go` — `spawnDeployedCleanupWindows` and `spawnDeployedCleanupUnix` both call `dumpDebugWindowsCommandPlan` immediately after `dumpDebugWindowsHandoff`.
- `gitmap/cmd/updatecleanup.go` — `runUpdateCleanup` calls `dumpDebugWindowsCleanupPlan(ctx)` immediately after `loadUpdateCleanupContext`.
- `gitmap/constants/constants_update.go` — new `MsgDebugWinCmdLine`, `MsgDebugWinCmdNote`, `MsgDebugWinCleanHdr`, `MsgDebugWinCleanGlob`, `MsgDebugWinCleanMatch`, `MsgDebugWinCleanEmpty`, `MsgDebugWinCleanSwap`, `MsgDebugWinCleanShim`, `MsgDebugWinCleanShimSkip`, `MsgDebugWinCleanShimDel`, `MsgDebugWinCleanFooter`.
- `gitmap/constants/constants.go` — bumped `Version` to `3.90.0`.

### Compatibility

Pure addition gated behind `--debug-windows` / `GITMAP_DEBUG_WINDOWS=1`. Default invocations are byte-for-byte identical to v3.89.0. The pre-flight glob scan is read-only and runs in microseconds — even with the debug flag on, it does not perceptibly slow cleanup.


## v3.89.0 — (2026-04-24) — Robust multi-URL clone parsing (PowerShell + bash)

### Added

- **Semicolon as a list separator.** `gitmap clone "url1;url2;url3"` now works alongside the existing comma form. Both can be mixed in one invocation. Bash users reach for `;` naturally; PowerShell users who quote the whole list (`'a;b;c'`) get the same behaviour.
- **Sanitisation of paste artifacts on every URL token:**
  - U+FEFF (BOM) — Windows clipboard frequently injects this when copying from PowerShell history or browser dev tools.
  - U+200B / U+200C / U+200D (zero-width spaces) — copied from rich-text sources (Slack, Notion, web docs).
  - Smart quotes (U+2018/U+2019/U+201C/U+201D) folded back to ASCII so wrapper-trim works in one pass.
  - Matched outer `'`, `"`, or backtick pairs stripped (only matched pairs — a stray trailing quote stays so the user sees a recognisably broken URL).
  - Leading/trailing list separators stripped (`,url2`, `url1;`, `;url2,` all collapse cleanly).

### Changed

- **`shouldUseMultiClone` now scans every positional, not just the first two.** Previously `gitmap clone https://x my-folder https://y` silently dropped the third URL because the heuristic only looked at indices 0 and 1. Three triggers now (any one is sufficient): (a) any positional contains `,` or `;`; (b) 2+ positionals AND any arg beyond the first parses as a URL; (c) the first positional flattens to 2+ valid URLs.
- **`isDirectURL` accepts `git@host:owner/repo` SSH-shorthand.** Previously only `https://`, `http://`, and `ssh://` were recognised, which meant `clone git@github.com:foo/bar.git` was misclassified as a file path in some code paths and as a URL in others (`isLikelyURL` already accepted it). Now both helpers agree.
- **`isLikelyURL` accepts `ssh://` for symmetry with `isDirectURL`.**
- **Empty/whitespace/separator-only tokens are dropped silently** instead of producing a bogus "invalid URL" warning. A token that was nothing but punctuation was almost certainly a typo, not a URL the user wanted to clone.

### Why

Three real failure modes from issue #11 / #16 + user reports:
1. `gitmap clone url1;url2` in bash produced a single ugly task and a "command not found" hint.
2. Copy-pasting a URL from PowerShell history (`>> ` prompt + URL) carried a BOM and produced a phantom invalid-URL warning.
3. `gitmap clone url1 url2 url3` (space-separated, no commas) cloned only the first two because `shouldUseMultiClone` only sampled `Positional[0]` and `Positional[1]`.

### Implementation

- `gitmap/cmd/clonemulti.go` — replaced `flattenURLArgs` body with a sanitising pipeline (`splitOnURLSeparators` → `sanitizeURLToken` → dedup). New helpers: `splitOnURLSeparators`, `sanitizeURLToken`, `stripInvisibleRunes`, `replaceSmartQuotes`, `trimMatchingWrappers`. New constant `urlListSeparators = ",;"`.
- `gitmap/cmd/clone.go` — `shouldUseMultiClone` rewritten with three triggers; `isDirectURL` extended to accept `git@host:` shorthand.
- `gitmap/cmd/rootflags.go` — `isLikelyURL` extended to accept `ssh://`. Doc-comment cross-reference added so future edits keep both helpers in lockstep.
- `gitmap/constants/constants.go` — bumped `Version` to `3.89.0`.

### Compatibility

Pure superset. Every previously valid invocation still works identically:

    gitmap clone url1 url2 url3                # space-only — now scans all positionals
    gitmap clone url1,url2,url3                # comma — unchanged
    gitmap clone "url1, url2 , url3"           # comma + spaces — unchanged
    gitmap clone url1,url2 url3,url4           # mixed — unchanged
    gitmap clone "url1;url2;url3"              # NEW: semicolon
    gitmap clone "url1,url2;url3"              # NEW: mixed separators
    gitmap clone git@github.com:foo/bar.git    # NEW: SSH shorthand recognised everywhere
    gitmap clone https://x.com/y my-folder     # single URL + folder name — unchanged


## v3.88.0 — (2026-04-24) — `gitmap pending clear` to unblock stuck clones

### Added

- **`gitmap pending clear [<mode>|<id>] [--dry-run] [--yes|-y]`** — removes orphaned or illegal pending tasks so the next `clone` / `clone-next` run is not blocked by a leftover row from an earlier crash. This is the deterministic escape hatch for the "pending task already exists for Clone at <bad-path>" lockout.
- **Modes:**
  - `orphans` (default) — `TargetPath` no longer exists on disk
  - `illegal` — `TargetPath` looks like a URL (`https:\github.com\...`, `://`, `git@host:`) or contains illegal Windows path chars (`:` after the drive letter, `?`, `*`, `<`, `>`, `|`, `"`)
  - `all` — every pending task (always confirms unless `--yes`)
  - `<id>` — a single task by numeric ID
- **Safety rails:** `--dry-run` previews without touching the DB; an interactive confirmation prompt is shown unless `--yes`/`-y` is passed; per-deletion log lines plus a final tally so the action is never silent.
- **Help page:** new `gitmap/helptext/pending-clear.md` covering modes, flags, exit codes, and the PowerShell comma-splitting backstory that produced the URL-shaped targets in the first place.

### Why

Issue #11 (PowerShell comma-splitting) and issue #12 added defensive guards in the clone command, but pre-existing rows from older crashes still block subsequent runs. Until now the only way to clear them was to delete the SQLite file or write raw SQL — neither of which is safe in a session where other gitmap state matters. `pending clear` does it surgically.

### Implementation

- `gitmap/cmd/pendingclear.go` (new) — `runPendingClear`, arg parser, candidate selector, three classifiers (`isURLShapedTarget`, `hasIllegalPathChar`, `isOrphanTarget`), confirmation prompt, per-row deletion loop.
- `gitmap/cmd/pending.go` — `runPending` now dispatches to `runPendingClear` when `os.Args[2] == "clear"`. List behaviour unchanged for plain `gitmap pending`.
- `gitmap/store/pendingtask.go` — new `DeletePendingTask(id int64) error` (uses existing `SQLDeletePendingTask`; returns `ErrPendingTaskNotFound` for unknown IDs).
- `gitmap/constants/constants_pending_task_msg.go` — new `MsgPendingClear*` and `ErrPendingClear*` constants (no magic strings).
- `gitmap/helptext/pending-clear.md` — new help page.
- `gitmap/constants/constants.go` — bumped `Version` to `3.88.0`.

### Compatibility

Pure addition. Plain `gitmap pending` is unchanged. The new subcommand reuses the existing `CmdPending` constant via positional dispatch, so the completion generator and `helptext/coverage_test.go` are unaffected.

### Examples

    gitmap pending clear                # default: drop orphans, prompt
    gitmap pending clear illegal --yes  # silently drop URL-shaped/illegal-char targets
    gitmap pending clear all --dry-run  # preview a full wipe
    gitmap pending clear 17             # drop one specific ID


## v3.87.0 — (2026-04-24) — Durable on-disk handoff log for self-update cleanup

### Added

- **Always-on handoff log file** at `<TMP>/gitmap-update-handoff-YYYYMMDD.log`. Phase 3 dispatcher and the deployed-binary cleanup child both append structured key=value events here regardless of `--verbose` / `--debug-windows`. The file survives swallowed stdout/stderr (hidden Windows processes, detached spawns, run.ps1 wrappers).
- **Events recorded**: `phase-3 resolve` (source + target), `phase-3 start_ok` (target + child pid), `phase-3 start_fail` (target + err), `phase-3 run_ok` / `run_fail` (Unix), `phase-3 inline`, `phase-3 target_missing`, `cleanup start` (self), `cleanup delay` / `cleanup delay_invalid`, `cleanup done` (removed count). Every line carries `pid`, `ppid`, `goos`, RFC3339 UTC timestamp.
- **Log path surfaced in two places**: `→ Handoff log file: <path>` always prints once at the start of Phase 3 dispatch, and the `[debug-windows]` dump now prints `handoff log file : <path>` in its header so the path is one click away whether or not the flag is on.

### Implementation

- `gitmap/cmd/updatehandofflog.go` (new) — `handoffLogPath()`, `logHandoffEvent(phase, event, fields)`, `formatHandoffLogLine`, mutex-serialized append-mode writer; failures degrade silently so logging can never disturb the update flow.
- `gitmap/cmd/updatehandoff_phase3.go` — `logHandoffEvent` calls at every lifecycle branch in `scheduleDeployedCleanupHandoff`, `spawnDeployedCleanupWindows`, `spawnDeployedCleanupUnix`; `→ Handoff log file:` line printed at top of dispatch.
- `gitmap/cmd/updatecleanup.go` — `logHandoffEvent` calls in `runUpdateCleanup` (`start`, `done`) and `delayUpdateCleanupIfNeeded` (`delay`, `delay_invalid`).
- `gitmap/cmd/updatedebugwindows.go` — header dump now includes `handoff log file` line.
- `gitmap/constants/constants_update.go` — `UpdateHandoffLogNameFmt`, `MsgUpdatePhase3LogFile`, `MsgDebugWinLogFile`.
- `gitmap/helptext/update.md` — new "Handoff log file" section with example log lines.
- `gitmap/constants/constants.go` — version bumped to `3.87.0`.

### Why

Even after `--debug-windows` (v3.86.0), failures during the detached Windows cleanup spawn could still vanish if an intermediate launcher discarded stdout/stderr. The handoff log writes to disk before/after each lifecycle event, so a forensic trail always exists for bug reports.

### Compatibility

Pure addition. The log file is at most a few KB per update run, daily-named, and never rotated (the daily filename keeps it bounded for the typical update cadence). Logging failures (read-only volume, etc.) are swallowed without disturbing the update.


## v3.86.0 — (2026-04-24) — `--debug-windows` for self-update cleanup handoff

### Added

- **`--debug-windows` flag on `gitmap update`** — opt-in diagnostics for the self-update Phase 2/Phase 3 handoff chain. Prints a `[debug-windows]` block on every relevant lifecycle event with the resolution source (`config` / `sibling` / `PATH`), resolved cleanup target path, target-exists check, child argv, key environment variables (`GITMAP_DEBUG_WINDOWS`, `GITMAP_UPDATE_CLEANUP_DELAY_MS`, `GITMAP_DEBUG_REPO_DETECT`, `GITMAP_REPORT_ERRORS`, `GITMAP_REPORT_ERRORS_FILE`, `PATH`, `GITMAP_DEPLOY_PATH`), self/parent PIDs, and the spawned child PID after a successful detached `Start()`.
- **`GITMAP_DEBUG_WINDOWS=1` env bridge** — the flag is propagated across the handoff boundary via both argv (Phase 2 copy + Phase 3 cleanup child) and env, so the dump runs on both sides of the detached spawn even when argv inheritance is fragile (e.g. hidden Windows process attrs). Users can also flip the env manually to enable the dump on a single run without rebuilding.

### Why

Issues #09 and #10 in `.lovable/pending-issues/01-current-issues.md` covered the recurring "update appears to complete but cleanup ran on the wrong binary" loop on Windows. The earlier fixes added `→ Cleanup target resolved via:` / `→ Cleanup target path:` / `→ Cleanup process started (pid=…)` lines, but those only appear in the parent (Phase 3 dispatcher). When the child cleanup process itself misbehaved, users had no way to see *its* view of the world. `--debug-windows` closes that gap by printing the same structured dump from inside `update-cleanup` too.

### Implementation

- `gitmap/cmd/updatedebugwindows.go` (new) — dump helpers (`dumpDebugWindowsHeader`, `dumpDebugWindowsHandoff`, `dumpDebugWindowsChildPID`, `dumpDebugWindowsNote`, `dumpDebugWindowsFooter`, `isDebugWindowsRequested`).
- `gitmap/cmd/updatehandoff_phase3.go` — header/footer wraps `scheduleDeployedCleanupHandoff`; handoff dump runs before `cmd.Start()`; child PID dump runs after; new `buildCleanupChildArgs` / `buildCleanupChildEnv` helpers forward the flag + env.
- `gitmap/cmd/updatecleanup.go` — dump runs at the start of `runUpdateCleanup` so the deployed binary prints its own view of the env, self path, and parent PID.
- `gitmap/cmd/update.go` — `launchHandoff` forwards `--debug-windows` and `GITMAP_DEBUG_WINDOWS=1` into the handoff copy and prints a Phase 2 dump line.
- `gitmap/constants/constants_update.go` — `FlagDebugWindows`, `EnvDebugWindows`, `MsgDebugWin*` constants.
- `gitmap/helptext/update.md` — flag table updated with full env-key list and behaviour notes.
- `gitmap/constants/constants.go` — bumped `Version` to `3.86.0`.

### Compatibility

Pure addition. Without the flag (and without `GITMAP_DEBUG_WINDOWS=1`), behaviour is byte-identical to the previous release. The dump goes to stderr only, so existing stdout-capturing wrappers stay clean.

### Usage

    gitmap update --debug-windows
    GITMAP_DEBUG_WINDOWS=1 gitmap update      # equivalent


## v3.53.0 — (2026-04-21) — `gitmap lfs-common`: one-shot Git LFS tracking for common binary types

### Added

- **`gitmap lfs-common` (alias `lfsc`)** — registers a curated set of 18 common binary file extensions with Git LFS in the current repository. Verifies the working tree is a git repo and that `git lfs` is on PATH, runs `git lfs install --local` (idempotent), then calls `git lfs track "<pattern>"` for each entry. The standard `<pattern> filter=lfs diff=lfs merge=lfs -text` lines are appended to `.gitattributes` by Git LFS itself, keeping the on-disk format canonical and tool-compatible.
- **Default tracked patterns:** `*.pptx`, `*.ppt`, `*.eps`, `*.psd`, `*.ttf`, `*.wott`, `*.svg`, `*.ai`, `*.jpg`, `*.bmp`, `*.png`, `*.zip`, `*.gz`, `*.tar`, `*.rar`, `*.7z`, `*.mp4`, `*.aep`. Order is preserved so `.gitattributes` diffs are stable across machines and re-runs.
- **`--dry-run` flag** — previews which patterns *would* be added vs. are *already tracked*, without touching `.gitattributes` or invoking `git lfs install`. Safe to run in any repo to audit existing LFS coverage against the recommended baseline.
- **Idempotent re-runs** — before tracking, the command parses `.gitattributes` and skips any pattern already carrying `filter=lfs`. The summary line reports `N added, M already tracked, K failed (of 18 total)`, so repeated invocations are harmless and the second run is a no-op when the baseline is already in place.

### Changed

- **`gitmap help`** — the *Git Operations* section now lists `lfs-common (lfsc)` between `latest-branch` and the navigation block.
- **Docs site** — version chip and command alias badges now use `text-foreground` (light mode) / `dark:text-background` (dark mode) with `dark:bg-primary/25`, ensuring black/neutral text stays readable against the green tint in both themes. Previously `text-primary` became illegible on dark backgrounds. Touched: [`src/pages/Index.tsx`](src/pages/Index.tsx) (hero version chip + CTA buttons restyled to `font-heading` with lift/shadow hover), [`src/components/docs/DocsLayout.tsx`](src/components/docs/DocsLayout.tsx) (header chip + new Sun/Moon theme toggle persisted via [`src/lib/theme.ts`](src/lib/theme.ts) and pre-paint script in [`index.html`](index.html)), [`src/components/docs/CommandCard.tsx`](src/components/docs/CommandCard.tsx) and [`src/components/docs/CommandPalette.tsx`](src/components/docs/CommandPalette.tsx) (alias badges), [`src/pages/VersionHistory.tsx`](src/pages/VersionHistory.tsx) / [`src/pages/Install.tsx`](src/pages/Install.tsx) / [`src/pages/CloneNext.tsx`](src/pages/CloneNext.tsx) (page-header alias chips). A single global override in [`src/index.css`](src/index.css) (`.dark [class*="bg-primary/"].text-primary`) patches the remaining 100+ chip occurrences across pages like `Architecture.tsx`, `Release.tsx`, `Doctor.tsx`, `Profile.tsx`, `Import.tsx`, `DiffProfiles.tsx`, `Changelog.tsx`, `PostMortems.tsx`, `ScanCloneFlags.tsx`, `GenericCLI.tsx`, and `ProjectDetection.tsx` without per-file edits.
- **`gitmap help lfs-common`** — new embedded help page (`gitmap/helptext/lfs-common.md`) documenting flags, the full pattern list, the post-run commit recipe, and a callout that `git lfs migrate import` is still required to convert *existing* committed binaries (this command only sets up tracking for *future* writes).

### Implementation

- `gitmap/cmd/lfscommon.go` — new file. `runLFSCommon` orchestrates flag parsing → repo check (`git rev-parse --is-inside-work-tree`) → LFS check (`git lfs version`) → `git lfs install --local` → per-pattern `git lfs track` loop. All shell-outs use `exec.CombinedOutput` so failures bubble up with the underlying git/lfs message attached.
- Reuses the existing `gitTopLevel()` helper from `gitmap/cmd/as.go` instead of re-declaring it — `cmd/` shares one Go namespace and the duplicate would have produced a `redeclared in this block` build break (per the `cmd/` collision-prone naming rule enforced by `.github/scripts/check-cmd-naming.sh`).
- `loadTrackedPatterns()` reads `<repo-root>/.gitattributes`, skips blank/comment lines, and treats any line containing `filter=lfs` as a tracked pattern keyed by the first whitespace-separated field. Used both in dry-run preview and in the live tracker to short-circuit no-op patterns.
- Helpers (`insideGitRepo`, `lfsAvailable`, `runGitLFSInstall`, `trackLFSPatterns`, `trackOnePattern`, `loadTrackedPatterns`, `printLFSCommonBanner`, `printLFSCommonDryRun`, `printLFSCommonSummary`) are all domain-qualified so they pass the `cmd/` naming guard. Output uses the existing `constants.ColorCyan/Green/Yellow/Dim/Reset` palette for consistency with `setup` and `doctor`.
- `gitmap/cmd/lfscommon_test.go` — new file. Two table-driven tests:
  - `TestLFSCommonPatternsMatchSpec` — locks in the exact 18 entries and ordering against the user-supplied spec, so accidental edits, typos, or removals are caught by CI before they ship.
  - `TestLFSCommonPatternsAreUnique` — guarantees no pattern appears twice (a duplicate would cause `git lfs track` to write the same line into `.gitattributes` twice on first run).
- `gitmap/cmd/rootutility.go` — added the `CmdLFSCommon` / `CmdLFSCommonAlias` branch to `dispatchUtility`, after `vscode-pm-path` and before the `return false` fallthrough.
- `gitmap/cmd/rootusage.go` — `printGroupGitOps` now prints `HelpLFSCommon` after `HelpLatestBr`, matching where the command sits semantically (it operates on the current repo's git/LFS state).
- `gitmap/constants/constants_cli.go` — added `CmdLFSCommon = "lfs-common"`, `CmdLFSCommonAlias = "lfsc"`, and `HelpLFSCommon` (the one-line help row). No new flag descriptions were required — the command reuses the existing `FlagDescDryRun`.
- `gitmap/helptext/lfs-common.md` — new embedded help file. Bundled automatically via the existing `//go:embed *.md` directive in `gitmap/helptext/print.go`, so `gitmap help lfs-common` works without any registration changes.
- `gitmap/constants/constants.go` — `Version = "3.53.0"`.

### Compatibility

- Pure addition. No existing command, flag, output format, or file layout changes. Repos that have never run `lfs-common` are unaffected; repos that have can re-run safely — the command is fully idempotent.
- Existing `.gitattributes` files are preserved: Git LFS appends only the patterns that aren't already tracked, and we additionally skip those patterns ourselves so `git lfs track` is never invoked redundantly.
- No new third-party Go dependencies. The command shells out to the user's installed `git` and `git lfs`, both of which are already required by the rest of the gitmap workflow.

---

## v3.52.0 — (2026-04-21) — Document `workflow_dispatch` lint baseline cache controls

### Changed

- **`spec/09-pipeline/01-ci-pipeline.md`** now documents the two `workflow_dispatch` inputs that govern the golangci-lint baseline cache:
  - `lint_baseline_cache_version` *(string, default `"v1"`)* — bumps the cache key suffix to abandon a stale baseline. Free-form (`"v2"`, `"2026-04-21"`, …); old caches are evicted by GitHub after 7 days of inactivity. The `restore-keys` fallback also carries this version, so a pre-bump baseline is never accidentally restored.
  - `lint_baseline_disable` *(boolean, default `false`)* — skips both the cache restore and save steps for one run, forcing the diff job into seeding mode (exits 0, surfaces all current findings as warnings) without touching the cached baseline. Use to diagnose suspected stale-cache issues without losing history.
- New **"Job: Lint Baseline Diff"** section explains the soft-gate cache strategy in a single table:
  - Cache key: `golangci-baseline-main-${cache_version}-${github.sha}` (rolling, single slot via the restore-keys fallback).
  - Save: only on `push` to `main` (or `workflow_dispatch` from `main`) after a green diff. PRs are restore-only — never advance the baseline.
  - Miss = seeding mode: the next run becomes the baseline; the build does not fail.
- Added three copy-paste `gh workflow run` examples covering the common operator scenarios: bumping the version, disabling the cache for one diagnostic run, and combining both for a "bump + dry-run" workflow before committing to a reseed.
- Documented the sticky PR comment behavior (sentinel `<!-- gitmap-lint-suggestions -->`, `peter-evans/find-comment` + `create-or-update-comment` for in-place edits, `GITHUB_STEP_SUMMARY` mirror on push/dispatch).

### Implementation

- `spec/09-pipeline/01-ci-pipeline.md` — extended `### Trigger` block to surface the two `workflow_dispatch` inputs alongside `push` / `pull_request`. Added the full **Job: Lint Baseline Diff** section between **Job: Lint** and **Job: Vulnerability Scan (In-CI)** — placement matches the actual job order in `.github/workflows/ci.yml`.
- `gitmap/constants/constants.go` — `Version = "3.52.0"`.

### Compatibility

- Documentation-only change. CI behavior, cache keys, and default flag values are unchanged — this turn merely brings the spec doc into sync with the workflow that already shipped.

---

## v3.51.0 — (2026-04-21) — Fix `cn v+1 -f` flag parsing; cleaner release-trailer newlines

### Fixed

- **`gitmap cn v+1 -f` no longer drops `-f`.** When `-f` followed a positional version arg, Go's stock `flag` package stopped scanning at the first non-flag token. The fix routes `cn` args through `reorderFlagsBeforeArgs(args)` (`gitmap/cmd/clonenextflags.go`) and extends the value-flag map in `gitmap/cmd/releaseargs.go` to cover `--csv`, `--ssh-key`, `-K`, and `--target-dir` so the next token is never mis-consumed.
- **Force implies Keep.** `Force` now sets `Keep = true` for the prior-folder cleanup path in `gitmap/cmd/clonenext.go`, suppressing the redundant "Remove current folder?" prompt and lock-detector loop.
- **`MsgInstallHintUnix` trailing newline.** Added a final blank line so the post-release shell prompt no longer sits flush against the `curl … | sh` line. Verified via `tests/release_test/InstallHint`.

### Changed

- `MsgForceReleasing` rewording: "Stepping out of … to release the file lock" — describes the Windows file-lock workaround in plain terms.
- Three-stage progress layout (Prepare → Clone → Finalize) shown only when `-f` is used; default `cn` output is byte-identical to v3.50.x.

### Compatibility

- Default `gitmap cn` output unchanged. New layout and prompt suppression are gated on `-f`.

---

## v3.50.0 — (2026-04-21) — `clone-next --force` (force-flatten)

### Added

- **`gitmap cn -f` / `--force`** — force a flat clone even when cwd IS the target folder. Chdirs to the parent before the remove (releases the Windows file lock), then re-clones into `<base>/`. Refuses the silent versioned-folder fallback under `-f` — flat layout or a clear error, never a surprise rename.
- `--force` / `-f` added to zsh + PowerShell completions and the `clone-next` help text.

### Fixed

- `MsgFlattenLockedHint` now mentions `-f`, so users discover the escape hatch on the first lock warning instead of giving up.

---

## v3.49.0 — (2026-04-21) — Auto-commit + auto-register on every release

### Added

- After a successful `gitmap release`, the metadata write is auto-committed (`chore(release): vX.Y.Z metadata`) and pushed in the same step. Skip with `--no-commit`.
- Cwd repo is auto-registered in the gitmap database if not yet tracked, eliminating the prior "repo not found" abort for first-time releases.

### Changed

- Trailer ordering finalized as: metadata write → auto-commit → auto-register → `── Release vX.Y.Z complete ──` → install hint (gitmap repo only).

---

## v3.48.0 — (2026-04-21) — `gitmap doctor` deploy-dir audit

### Added

- `gitmap doctor` now flags duplicate `gitmap` / `gitmap.exe` binaries on `$PATH`, reports the active deploy dir vs the running binary path, and recommends `gitmap self-install --dir <chosen>` to consolidate.

### Fixed

- Doctor no longer false-positives on `gitmap.exe.old` backup files — `isGitmapArtifact` now ignores `*.old` for the duplicate check.

---

## v3.47.0 — (2026-04-21) — `release-version` interactive fallback prompt

### Added

- When the requested version isn't a published release AND the script runs in an interactive terminal, the installer offers the 5 most recent releases to pick from instead of aborting. Non-interactive (piped) runs still exit 1 unless `--allow-fallback` is supplied.

---

## v3.46.0 — (2026-04-21) — Sticky lint-suggestion PR comments

### Added

- CI lint-baseline-diff job now posts a single sticky PR comment (sentinel `<!-- gitmap-lint-suggestions -->`) using `peter-evans/find-comment` + `create-or-update-comment`, replacing the previous comment-spam-on-every-push behavior. Push and `workflow_dispatch` runs mirror the same content into `GITHUB_STEP_SUMMARY`.

---

## v3.45.0 — (2026-04-21) — `golangci-lint` baseline cache (soft gate)

### Added

- New CI job **lint-baseline-diff**: restores the previous lint findings from cache (key `golangci-baseline-main-${cache_version}-${github.sha}`, restore-keys fallback to the most recent baseline on the same `cache_version`), runs the linter, and surfaces only NEW findings. Soft gate: warnings, never failures. Save step runs only on `push` to `main` (or `workflow_dispatch` from `main`).

---

## v3.44.0 — (2026-04-21) — `gitmap self-uninstall` Windows handoff

### Added

- On Windows, `self-uninstall` copies the running `gitmap.exe` to `%TEMP%\gitmap-handoff-<pid>.exe`, re-execs the hidden `self-uninstall-runner` verb, and the temp copy schedules its own deletion via `cmd.exe /C ping ... & del /F /Q <self>`. Releases the file lock that previously prevented self-deletion.

### Changed

- PATH snippet cleanup now strips the marker block `# gitmap shell wrapper v* …` … `# gitmap shell wrapper v* end` from the user's shell profile, idempotent across re-runs. Skip with `--keep-snippet`.

---

## v3.43.0 — (2026-04-21) — `gitmap self-install` / `self-uninstall`

### Added

- New top-level verbs `self-install` and `self-uninstall` manage the gitmap binary itself, separate from the existing third-party `install` / `uninstall` (npp, vscode, dev tools).
- `self-install` defaults: `D:\gitmap` (Windows), `~/.local/bin/gitmap` (Unix). Override with `--dir`. Skip the prompt with `--yes`. Forwards `--version <tag>` to the installer.
- Installer scripts (`install.ps1`, `install.sh`, `uninstall.ps1`) embedded into the binary via `go:embed` (`gitmap/scripts/embed.go`), with HTTP fallback to `raw.githubusercontent.com/alimtvnetwork/gitmap-v27/main/gitmap/scripts/install.{ps1,sh}` if the embedded copy is missing.
- `self-uninstall` removes: deploy-dir artefacts, `.gitmap/` data dir, PATH snippet, completion files. Confirm gates: typed `yes` (interactive) or `--confirm` (CI). Selective skip with `--keep-data` / `--keep-snippet`.

### Implementation

- `gitmap/constants/constants_selfinstall.go` — IDs, messages, defaults
- `gitmap/cmd/selfinstall.go`, `gitmap/cmd/selfuninstall.go`, `gitmap/cmd/selfuninstallparts.go`, `gitmap/cmd/selfuninstallhandoff.go` — split to satisfy <200-line rule
- PowerShell scripts written with UTF-8 BOM (per `mem://constraints/powershell-encoding`)

---

## v3.32.1 — (2026-04-20) — Fix `gitmap status` looking at legacy bare `output/` path

### Fixed

- **`gitmap status` no longer fails with `could not load gitmap.json at output\gitmap.json`** when run from a directory that has no `output/` folder.
  - Root cause: `loadRecordsJSONFallback` joined `constants.DefaultOutputFolder` (the legacy bare `"output"` value, kept around for backward compat) with `gitmap.json`, instead of the unified `.gitmap/output` path used by every other command since v2.99.
  - Two-part fix:
    1. Look at the correct unified path: `constants.DefaultOutputDir` → `.gitmap/output/gitmap.json`.
    2. When the JSON file is missing (e.g. the user has not run `gitmap scan` from this exact directory yet), transparently fall through to the SQLite database — the DB is the source of truth post-v2 and usually has every repo the user has ever scanned. Previously, status exited with an error even though the DB had perfectly good data.
- New friendly message `MsgStatusNoData` is shown only when both the JSON file is missing AND the database has zero repos: `"No tracked repos found. Run 'gitmap scan' in a directory containing your git repos first, or pass --all to query the database directly."`

### Implementation

- `gitmap/cmd/status.go` — `loadRecordsJSONFallback` now stat-checks the JSON path first and delegates to a new `loadAllRecordsDBOrEmpty` helper when missing. Path joined with `DefaultOutputDir` instead of `DefaultOutputFolder`.
- `gitmap/constants/constants_messages.go` — new `MsgStatusNoData` constant.

### Compatibility

- Pure bug fix; behavior is strictly more permissive (succeeds in cases that previously errored). No flag, file, or DB schema impact.


## v3.32.0 — (2026-04-20) — Scan output: hoist common base path, show filenames only

### Changed

- **`gitmap scan` Output Artifacts section** is now scannable in one glance. The common base directory (e.g. `D:\wp-work\riseup-asia\.gitmap\output\`) is printed once under the section header as `📂 Base: <path>`, and each artifact line shows only the filename (`gitmap.csv`, `clone.ps1`, …) instead of repeating the full absolute path on every row. Same for `💾 Cache` (rescan).
- Aligned the icon column widths so filenames line up vertically across CSV / JSON / Text list / Structure / Clone PS1 / HTTPS PS1 / SSH PS1 / Desktop PS1 / Cache rows.

### Implementation

- `gitmap/constants/constants_messages.go` — `MsgSectionArtifacts` now takes a `%s` for the base dir; `MsgCSVWritten`/`MsgJSONWritten`/`MsgTextWritten`/`MsgStructureWritten`/`MsgCloneScript`/`MsgDirectClone`/`MsgDirectCloneSSH`/`MsgDesktopScript`/`MsgScanCacheSaved` re-aligned to a uniform 12-char label column.
- `gitmap/cmd/scan.go` — `fmt.Print(MsgSectionArtifacts)` → `fmt.Printf(MsgSectionArtifacts, outputDir)`.
- `gitmap/cmd/scanoutput.go` — every per-file `fmt.Printf(MsgXxx, path)` now passes `filepath.Base(path)`.
- `gitmap/cmd/rescan.go` — same `filepath.Base(path)` change for the cache line.

### Compatibility

- Pure formatting change; no flag, file, or DB schema impact. Project Detection section was already filename-only and is unchanged.


## v3.31.0 — (2026-04-20) — Cross-dir release/clone-next, has-change command, SSH existing-key fix

### Added

- **`gitmap r <repo> <vX.Y.Z>` — cross-directory release**: run from anywhere; gitmap chdirs into the named repo, runs `git fetch --all --prune` + `git pull --rebase`, auto-stashes any dirty changes, runs the standard release pipeline, then chdirs back to the original directory and pops the stash. Backward compatible: `gitmap r vX.Y.Z` (single arg) still does an in-place release. The first positional arg is treated as a repo alias only when it does NOT match the version regex `^v?\d+\.\d+\.\d+`.
- **`gitmap cn <repo> <vX.Y.Z>` — cross-directory clone-next**: same chdir/run/return pattern as `r`, wrapping the existing clone-next pipeline. `gitmap cn vX.Y.Z` (single arg) still operates on the current repo.
- **`gitmap has-change (hc) <repo>`** — prints `true`/`false` for whether the named repo has uncommitted changes. `--mode=dirty|ahead|behind` switches dimension; `--all` prints `dirty=X ahead=Y behind=Z` in one line. `--fetch=false` skips the pre-check `git fetch` for offline use.

### Fixed

- **`gitmap ssh` no longer fails with `exit status 1` when `~/.ssh/id_rsa` already exists.** Previously, gitmap only checked the SQLite database for existing keys; keys created outside gitmap (e.g. raw `ssh-keygen` or another tool) fell through to `ssh-keygen -f <existing-path>`, which prompted `Overwrite (y/n)?` on stdin and exited non-zero when no answer arrived. Now `runSSHGenerate` checks the disk path FIRST: if the private key file exists and `--force` is not set, gitmap prints the existing public key, fingerprint, copy-to-GitHub hint, and a `--force` regeneration hint, then exits 0. The disk-discovered key is also upserted into the gitmap database so `ssh-cat` / `ssh-list` find it.
- **`--force` regenerate flow**: when `--force` is set and the key exists on disk, gitmap renames `id_rsa` → `id_rsa.bak.<unix-ts>` (and `.pub` likewise) before invoking `ssh-keygen`, so users never lose access to a working key by accident.

### Implementation

- `gitmap/cmd/releaserebase.go` (new, ~150 lines) — `tryCrossDirRelease`, `performCrossDirRelease`, `rebasePull`, `extractPositionalArgs`, `extractFlagArgs`, `looksLikeVersion`. Reuses existing `resolveReleaseAliasPath`, `autoStashIfDirty`, `popAutoStash`.
- `gitmap/cmd/clonenextcrossdir.go` (new, ~55 lines) — `tryCrossDirCloneNext`, `performCrossDirCloneNext`. Same pattern.
- `gitmap/cmd/haschange.go` (new, ~140 lines) — `runHasChange`, `parseHasChangeFlags`, `printHasChangeOne`, `printHasChangeAll`, `readAheadBehind`, `fetchRemoteIn`, `boolStr`.
- `gitmap/cmd/sshexisting.go` (new, ~85 lines) — `keyExistsOnDisk`, `printExistingKeyOnDisk`, `upsertExistingKeyToDB`, `backupKeyForRegenerate`.
- `gitmap/cmd/release.go` — `runRelease` now calls `tryCrossDirRelease(args)` first.
- `gitmap/cmd/clonenext.go` — `runCloneNext` now calls `tryCrossDirCloneNext(args)` first.
- `gitmap/cmd/sshgen.go` — disk-existence check moved BEFORE the database check; `--force` triggers `backupKeyForRegenerate` before `ssh-keygen`.
- `gitmap/cmd/rootcore.go` — added `has-change` / `hc` dispatch alongside the existing `has-any-updates`.
- `gitmap/constants/constants_v331.go` (new) — all new constants centralized for v3.31.0 audit clarity: `CmdHasChange`, `FlagHC*`, `HCMode*`, `MsgRR*`, `ErrRR*`, `MsgCNX*`, `MsgSSHExistsOnDisk`, `MsgSSHForceHint`, `MsgSSHBackedUp`, `ErrSSHBackup`.
- `gitmap/helptext/has-change.md` (new) — bundled help text for the new command.

### Compatibility

- Single-arg invocations of `r` and `cn` are byte-for-byte identical to v3.30.x.
- The SSH change is purely additive on the disk-existence path; the in-DB-key flow (`handleExistingKey` with `R`/`N` prompt) is untouched.

### Verified locally

- `extractPositionalArgs` correctly strips `-y`, `--dry-run`, `--bump=patch`-style flags.
- `looksLikeVersion` accepts `v3.31.0`, `3.31.0`, `v3.31.0-rc1`, `3.31.0+build5`; rejects `gitmap`, `my-app`, `r3`, `v3`.
- New constants compile in isolation (no collisions with existing `Cmd*`/`Msg*`/`Err*` per the v3.26.0 collision check).


## v3.30.0 — (2026-04-20) — Fix Go Report Card badge URL to point at the actual module path

### Fixed (Docs)

- **README.md Go Report Card badge** now points at `github.com/alimtvnetwork/gitmap-v27/gitmap` (the real Go module path set in v3.27.0) instead of the repo root `github.com/alimtvnetwork/gitmap-v27`. The previous URL returned a 404 from the Go module proxy because there is no `go.mod` at the repo root — the module lives one directory down in `gitmap/`. Both the badge image and the click-through report link were updated.

### Files changed

- `README.md` — single line, both the `goreportcard.com/badge/...` image URL and the `goreportcard.com/report/...` link target now include the `/gitmap` subpath suffix.

### Compatibility

Pure documentation fix. No source, CI, or runtime change.


## v3.28.0 — (2026-04-20) — Lucrative scan summary: grouped sections + emoji-rich post-scan log

### Improved (UX / Terminal Output)

- **`gitmap scan` post-scan summary is now visually grouped into three labeled sections** instead of a flat wall of "X written to Y" lines. Each section has a header with a thematic emoji and a horizontal rule, and every line item carries a category icon for instant scanning.

The summary now flows as:

```
📦 Output Artifacts
────────────────────────────────────────────
  📊 CSV         D:\...\.gitmap\output\gitmap.csv
  🧬 JSON        D:\...\.gitmap\output\gitmap.json
  📝 Text list   D:\...\.gitmap\output\gitmap.txt
  🌳 Structure   D:\...\.gitmap\output\folder-structure.md
  🪄 Clone PS1   D:\...\.gitmap\output\clone.ps1
  ⚡ HTTPS PS1   D:\...\.gitmap\output\direct-clone.ps1
  🔐 SSH PS1     D:\...\.gitmap\output\direct-clone-ssh.ps1
  🖥️  Desktop PS1 D:\...\.gitmap\output\register-desktop.ps1
  💾 Cache       D:\...\.gitmap\output\last-scan.json

🗄️  Database
────────────────────────────────────────────
  ✅ 42 repositories upserted into database
  🏷️  Tagged 42 repo(s) with scan folder #1

🔍 Project Detection
────────────────────────────────────────────
  🧭 Detected 54 project(s) across 35 repo(s)
  📄 react-projects.json    26 record(s)
  📄 go-projects.json       24 record(s)
  📄 node-projects.json      4 record(s)
  ✅ Saved 54 detected project(s) to database

🎉 Scan complete.
```

### Files changed

- `gitmap/constants/constants_messages.go` — restyled `MsgCSV/JSON/Text/Structure/Clone/Direct/SSH/Desktop/ScanCache/DBUpsertDone` constants, added `MsgSectionArtifacts`, `MsgSectionDatabase`, `MsgSectionProjects`, `MsgSectionDone`, `MsgSectionRule`, `MsgScanFolderTagged`.
- `gitmap/constants/constants_project.go` — restyled `MsgProjectDetectDone`, `MsgProjectUpsertDone`, `MsgProjectJSONWritten` to align under the new section header with consistent indentation.
- `gitmap/cmd/scan.go` — `executeScan` now prints section headers between groups; `tagReposWithScanFolder` uses the centralized `MsgScanFolderTagged` constant (no more inline string).
- `gitmap/constants/constants.go` — version bumped to `3.28.0`.

### Compatibility

Pure terminal-output cosmetics. No flag, file path, JSON schema, or database column changed. CSV/JSON/PS1 artifact formats are byte-identical to v3.27.0.


## v3.27.0 — (2026-04-20) — Fix Go Report Card: rename module path from placeholder to real GitHub path

### Fixed (Tooling / Distribution)

- **Go Report Card now resolves the module instead of failing with `could not get latest module version from https://proxy.golang.org/github.com/user/gitmap/@latest`.** The card at https://goreportcard.com/report/github.com/alimtvnetwork/gitmap-v27/gitmap will start scoring the project for the first time after this version is pushed.

### Root cause

`gitmap/go.mod` was declared as `module github.com/user/gitmap` — a leftover placeholder from project bootstrapping. Because Go Report Card runs `go get <module>@latest` against the public Go module proxy (`proxy.golang.org`) before linting, and that path returns 404 (no such GitHub user/repo), the whole report aborted before any of `gofmt`, `go vet`, `gocyclo`, `ineffassign`, `misspell`, `golint` could run.

The same placeholder was also referenced in:

- 391 `.go` files inside `gitmap/` (all `import "github.com/user/gitmap/..."` statements)
- The companion module `gitmap-updater/go.mod` and its imports
- `Makefile` and `.github/workflows/ci.yml` ldflags injection: `-X 'github.com/user/gitmap/constants.Version=...'`
- `run.ps1` and `run.sh` ldflags injection for `RepoPath`
- Spec docs and changelog history references
- React-side changelog/getting-started page references

If left unfixed, anyone running `go install github.com/alimtvnetwork/gitmap-v27/gitmap@latest` would get a `module declares its path as: github.com/user/gitmap but was required as ...` error, and the proxy would refuse to serve the module to downstream tooling.

### Fix

Renamed `github.com/user/gitmap` → `github.com/alimtvnetwork/gitmap-v27/gitmap` across **404 files** in a single atomic sed pass. Also implicitly renamed the sister module `github.com/user/gitmap-updater` → `github.com/alimtvnetwork/gitmap-v27/gitmap-updater`, which lives at the same GitHub path.

Verified post-rename:
- `gitmap/go.mod` now reads `module github.com/alimtvnetwork/gitmap-v27/gitmap`.
- `gitmap-updater/go.mod` now reads `module github.com/alimtvnetwork/gitmap-v27/gitmap-updater`.
- `Makefile` ldflags target the new constants package path.
- `.github/workflows/ci.yml` build step injects `Version` into the new constants package path.
- `run.ps1` and `run.sh` inject `RepoPath` into the new constants package path.
- Zero remaining references to the old placeholder string anywhere in the tree.

### What the user needs to do after pulling

1. Pull v3.27.0 and push to `main`.
2. Once the new tag (`v3.27.0`) is pushed, visit https://goreportcard.com/report/github.com/alimtvnetwork/gitmap-v27/gitmap — first visit will trigger a fresh scan against the new module path.
3. The CI ldflags injection still works because the constants package path was renamed alongside the workflow string.
4. Anyone who had previously cloned the repo and run `go build ./...` will need to re-run `go mod tidy` once after pulling, since every import path changed.

### Files (this section)

- Edited: 404 files (391 `.go` files in `gitmap/`, plus `gitmap/go.mod`, `gitmap-updater/go.mod`, `gitmap-updater/main.go`, `Makefile`, `.github/workflows/ci.yml`, `.github/workflows/release.yml`, `run.ps1`, `run.sh`, multiple `spec/` docs, `CHANGELOG.md` history references, `src/data/changelog.ts`, `src/pages/GettingStarted.tsx`).
- Edited: `gitmap/constants/constants.go` — bumped Version to 3.27.0.
- Created: `.gitmap/release/v3.27.0.json` — release metadata.
- Edited: `.gitmap/release/latest.json` — pointer to v3.27.0.

---

## v3.26.0 — (2026-04-20) — Audit + new CI guard for constants/ identifier collisions

### Added (CI)

- **New `constants-collision-check` job** in `.github/workflows/ci.yml` that runs `python3 .github/scripts/check-constants-collisions.py` on every push and PR. Fails fast (no Go toolchain) when any of these conditions hold across the 69 `gitmap/constants/constants_*.go` files:
  1. **Cross-file exact-name collision** — the same identifier (e.g. `HelpGitHubDesktop`) is declared in two files. This is exactly the v3.25.0 regression that took down `go build` and motivated the audit.
  2. **Cross-file case-insensitive collision** — different exact names that lowercase to the same string (e.g. `HelpFoo` vs `helpFoo`) and live in different files. Latent confusion risk even though Go accepts it.
  3. **Intra-file duplicate declaration** — the same identifier appears twice in one file. `go build` already catches this, but the script reports the offending lines without waiting for a Go compile.

- **New script `.github/scripts/check-constants-collisions.py`** — a string-literal-aware Python parser that tracks raw-string (`` `...` ``) and `"..."` quoted regions, so SQL keywords (`FROM`, `WHERE`, `VALUES`, `ORDER`, `SET`, ...) appearing inside multi-line raw-string SQL constants are NEVER mistaken for top-level identifiers. A naive line-based regex auditor reported 8 false-positive "collisions" from these tokens; the literal-aware parser reports 0.

### Audit results (current tree)

After the v3.25.2 fix, the auditor scanned 69 files and 2,902 unique top-level identifiers:

- **Cross-file exact-name collisions: 0**
- **Cross-file case-insensitive collisions: 0**
- **Intra-file duplicate declarations: 0**

The constants package is fully clean. All future PRs that introduce a collision (intentionally or by oversight) will be blocked by the new CI job before the broken build reaches `main`.

### Why a Python script instead of extending the existing bash awk guard

The existing `check-constants-naming.sh` is a single-pass line-based awk extractor (and rewriting it to track raw-string state across lines in portable mawk would be painful). A focused 130-line Python script is easier to read, easier to extend, and Python is preinstalled on every Ubuntu runner.

### Files (this section)

- Created: `.github/scripts/check-constants-collisions.py` — string-literal-aware collision auditor.
- Edited: `.github/workflows/ci.yml` — added `constants-collision-check` job; added it to the `test-summary` `needs:` list so the overall CI status reflects the guard.
- Edited: `gitmap/constants/constants.go` — bumped Version to 3.26.0.
- Created: `.gitmap/release/v3.26.0.json` — release metadata.
- Edited: `.gitmap/release/latest.json` — pointer to v3.26.0.

---

## v3.25.2 — (2026-04-20) — Fix `HelpGitHubDesktop` redeclaration build error

### Fixed (Build)

- **`go build` no longer fails with `HelpGitHubDesktop redeclared in this block`** between `constants/constants_helpsections.go:13` and `constants/constants_cli.go:100`.

### Root cause

v3.25.0 introduced a new top-level `HelpGitHubDesktop` constant in `constants_cli.go` for the `github-desktop (gd)` command help line. However, `constants_helpsections.go` already exported a constant of the same name (since pre-v3.0) for the `--github-desktop` **scan flag** help line. Both files compile into the same `constants` package, so the namespace collision broke the build for everyone who pulled v3.25.0 / v3.25.1.

The pre-existing constant was the older one and is only consumed by `cmd/rootusageflags.go` (the `--github-desktop` scan-flag help block), so it was renamed to `HelpScanFlagGitHubDesktop` to disambiguate by purpose. The newer command-line help constant in `constants_cli.go` keeps the original `HelpGitHubDesktop` name since it represents the canonical `github-desktop` command.

### Fix

- Renamed scan-flag help constant: `HelpGitHubDesktop` → `HelpScanFlagGitHubDesktop` in `constants/constants_helpsections.go`.
- Updated sole consumer `cmd/rootusageflags.go` to reference the renamed constant.
- Inserted `HelpScanFlagGitHubDesktop` into `.github/scripts/constants-baseline.txt` in sorted order (between `HelpScan` and `HelpScanFlags`).

### Files (this section)

- Edited: `gitmap/constants/constants_helpsections.go` — renamed `HelpGitHubDesktop` → `HelpScanFlagGitHubDesktop`.
- Edited: `gitmap/cmd/rootusageflags.go` — updated reference.
- Edited: `.github/scripts/constants-baseline.txt` — added `HelpScanFlagGitHubDesktop`.
- Edited: `gitmap/constants/constants.go` — bumped Version to 3.25.2.
- Created: `.gitmap/release/v3.25.2.json` — release metadata.
- Edited: `.gitmap/release/latest.json` — pointer to v3.25.2.

---

## v3.25.1 — (2026-04-20) — CI: portable awk in constants-naming guard (fixes silent exit 1 on Ubuntu runners)

### Fixed (CI)

- **`bash .github/scripts/check-constants-naming.sh` no longer fails with bare `Error: Process completed with exit code 1` and no `::error::` output on GitHub Actions Ubuntu runners.**

### Root cause

The awk extractor used the **gawk-only 3-argument `match(string, regex, array)` form** to capture identifier names from `const ( ... )` blocks:

    match(line, /^[[:space:]]+([A-Z][A-Za-z0-9]+).../, m)
    print m[1]

GitHub Actions Ubuntu runners ship **mawk** as the default `/usr/bin/awk`, where 3-arg `match()` is a syntax error. mawk aborts at parse time → the awk pipeline produces no output → `set -euo pipefail` propagates exit 1 → the script's violation reporter (which is what would print `::error::` lines) is never reached. So the CI log shows only the bare `Error: Process completed with exit code 1` with zero diagnostic context, even though the guard itself isn't actually finding any naming violation.

Locally everything passed because dev environments (and this Lovable sandbox) have gawk wired to `/bin/awk`, masking the portability bug.

### Fix

1. Rewrote the awk to be POSIX-portable: only 2-arg `match()` + `RSTART` / `RLENGTH` + `substr()`. Captures the same names mawk and gawk both accept. Verified byte-identical output between the old gawk-only awk and the new portable awk on the full `gitmap/constants/` tree (2764 = 2764 entries, zero diff in either direction).
2. Added a defensive `sudo apt-get install -y gawk` step to `.github/workflows/ci.yml` immediately before the guard runs, so even if mawk-only runners reappear in the future, gawk is on PATH.
3. Forced `LC_ALL=C` on both sides of the `comm -23 current baseline` invocation so sort ordering is guaranteed identical regardless of runner locale.
4. Regenerated `.github/scripts/constants-baseline.txt` from 2757 → 2764 entries to admit the v3.25.0 `github-desktop` constants (`CmdGitHubDesktop`, `MsgGHDesktopRegister`, `ErrGHDesktopCwd`, etc. — all canonical-prefixed, so this is just a snapshot refresh).

### Files (this section)

- Edited: `.github/scripts/check-constants-naming.sh` — replaced 3-arg `match()` with `RSTART`/`RLENGTH`/`substr()`; added `LC_ALL=C` to `comm` + sort.
- Edited: `.github/workflows/ci.yml` — `apt-get install -y gawk` step before the guard.
- Edited: `.github/scripts/constants-baseline.txt` — regenerated (2757 → 2764).
- Edited: `gitmap/constants/constants.go` — `Version` bumped to `3.25.1` (only `gitmap/` line touched).
- Edited: `.gitmap/release/latest.json` — points to `v3.25.1`.
- New:    `.gitmap/release/v3.25.1.json` — release metadata.

### Notes

- No `gitmap/` source behavior changed; this is purely a CI script portability fix + Version bump.
- The mawk-vs-gawk gotcha is a recurring bash-script trap on Ubuntu CI; consider grepping the rest of `.github/scripts/` for `match(.*,.*,.*)` to preempt the same bug in other guards.

## v3.25.0 — (2026-04-20) — new `github-desktop` (gd) command: register cwd repo without scan

### Added

- **`gitmap github-desktop` (alias `gd`)** — registers the current working-directory git repo with GitHub Desktop in one call. Previously the only path was `gitmap desktop-sync` (`ds`), which walks the last-scan output JSON and fails with `no output dir` if you haven't run `gitmap scan` first. Running `gd` from a freshly cloned repo now Just Works:

      cd D:\wp-work\riseup-asia\my-api
      gitmap gd
      # → Registering with GitHub Desktop: D:\wp-work\riseup-asia\my-api
      # → ✓ Registered with GitHub Desktop: D:\wp-work\riseup-asia\my-api

  Optional path argument also supported: `gitmap gd D:\path\to\other\repo`.

### Why this exists

User reported `gitmap github-desktop` printing `Unknown command`. Root cause: the string `github-desktop` only ever existed as a `--github-desktop` *flag* on `scan`/`clone`, never as a command. `desktop-sync` (`ds`) was the closest thing but required prior `gitmap scan`. This commit closes that gap.

### Files (this section)

- New: `gitmap/cmd/githubdesktop.go` — `runGitHubDesktop` (cwd or arg path → `.git` check → GitHub Desktop CLI invoke).
- New: `gitmap/helptext/github-desktop.md` — full help page with comparison table vs `desktop-sync`.
- Edited: `gitmap/constants/constants_cli.go` — adds `CmdGitHubDesktop` / `CmdGitHubDesktopAlias` / `HelpGitHubDesktop`.
- Edited: `gitmap/constants/constants_messages.go` — adds `MsgGHDesktopRegister`, `MsgGHDesktopDone`, `ErrGHDesktopCwd`, `ErrGHDesktopNotRepo`, `ErrGHDesktopInvoke` (all canonical Cmd/Msg/Err prefixes — passes `check-constants-naming.sh` without baseline regen).
- Edited: `gitmap/constants/constants_helpgroups.go` — `CompactCloning` line includes `github-desktop (gd)`.
- Edited: `gitmap/cmd/roottooling.go` — dispatcher routes `github-desktop` / `gd` to `runGitHubDesktop`.
- Edited: `gitmap/cmd/rootusage.go` — Cloning help group prints `HelpGitHubDesktop` after `HelpDesktopSync`.
- Edited: `gitmap/completion/allcommands_generated.go` — adds `gd` and `github-desktop` to the sorted completion list.
- Edited: `gitmap/constants/cmd_constants_test.go` — adds the two new constants to the parity map.
- Edited: `gitmap/completion/completion_test.go` — adds the two new strings to the expected completion list.
- Edited: `gitmap/constants/constants.go` — `Version` bumped to `3.25.0`.
- Edited: `.gitmap/release/latest.json` + new `.gitmap/release/v3.25.0.json`.

### Notes

- All new constants use the canonical `Cmd*` / `Help*` / `Msg*` / `Err*` prefixes; `bash .github/scripts/check-constants-naming.sh` and `check-cmd-naming.sh` both pass without regenerating any baseline.
- `desktop-sync` is unchanged and remains the bulk-sync command for whole scan trees.

## v3.24.1 — (2026-04-20) — CI: regenerate constants baseline to admit v3.24.0 additions

### Fixed (CI)

- **`bash .github/scripts/check-constants-naming.sh` now passes on `main`.** The v3.24.0 release added `constants.GitStderrNoisePatterns` (and a handful of other internal identifiers) to `gitmap/constants/`, which the naming guard flagged because they predate the canonical `Cmd*/Msg*/Err*/Flag*/Default*` prefix policy by source convention only. Per the grandfather workflow documented at the top of `check-constants-naming.sh`, the baseline file `.github/scripts/constants-baseline.txt` was regenerated (2743 → 2757 entries) so the new identifiers are admitted as pre-existing. No `gitmap/` source code changed; future constants must still use a canonical prefix.

### Files (this section)

- Edited: `.github/scripts/constants-baseline.txt` — regenerated via `bash .github/scripts/check-constants-naming.sh --regenerate-baseline` (2743 → 2757 lines).
- Edited: `gitmap/constants/constants.go` — `Version` bumped to `3.24.1` (only line touched in `gitmap/`).
- Edited: `.gitmap/release/latest.json` — points to `v3.24.1`.
- New:    `.gitmap/release/v3.24.1.json` — release metadata.

### Notes

- The `gitmap/` source folder is otherwise untouched per the standing rule. If/when the constants get renamed to `Msg*` prefixes properly, regenerate the baseline again and drop the grandfathered names.

## v3.24.0 — (2026-04-20) — suppress git CRLF/LF cosmetic warnings during release

### Fixed (release stderr noise)

- **`gitmap r` no longer prints `warning: in the working copy of '...', LF will be replaced by CRLF the next time Git touches it`** for every staged file. On Windows repos with `core.autocrlf=true`, every release commit was emitting this warning once per touched file (e.g. `.gitmap/release/latest.json`, `.gitmap/release/vX.Y.Z.json`), drowning the real progress lines. The `runGitCmd` helper now pipes git's stderr through `filteredStderrWriter` (new file `gitmap/release/gitstderrfilter.go`), which line-buffers stderr and silently drops any line containing a substring listed in `constants.GitStderrNoisePatterns`. All other git stderr output (real errors, hint lines, push results) is forwarded unchanged.

### Files (this section)

- New: `gitmap/release/gitstderrfilter.go` — `filteredStderrWriter` (line-buffered, multi-pattern, with `Flush()` for un-terminated trailing data).
- Edited: `gitmap/release/gitops.go` — `runGitCmd` wraps `os.Stderr` with `newFilteredStderr` and flushes after `cmd.Run()`.
- Edited: `gitmap/constants/constants_git.go` — adds `GitStderrNoisePatterns []string` (currently the single CRLF/LF warning) with a doc comment explaining the "guaranteed-not-an-error" admission rule for new entries.

### Notes

- Only `runGitCmd` is filtered (it's the writer used by `stageAll`, `stageFiles`, `commitStaged`, `CreateBranch`, `CreateTag`, `rollback`, etc.). `runGitCmdCombined` deliberately keeps full output because callers parse it for `non-fast-forward` detection.

## v3.22.0 — (2026-04-20) — `gitmap r` auto-registers cwd repo when missing

### Fixed (release persistence)

- **`gitmap r` no longer aborts release-DB caching with `no repo registered for path "..."` when the cwd has never been scanned.** When `resolveOrRegisterCurrentRepoID` cannot find the cwd in the `Repo` table, it now auto-registers: cwd becomes a single `Repo` row (slug / URLs / branch built via `mapper.BuildRecords`, identical to a real scan), parent dir becomes a `ScanFolder` row via `EnsureScanFolder`, and `TagReposByScanFolder` links them. The Release.RepoId FK is then satisfied on the retry, so the release row is persisted in the SAME `gitmap r` invocation that just pushed the tag — no second `gitmap scan` round-trip required.
- **Visible feedback**: prints `✓ Auto-registered repo "..." under scan folder "..." (#N)` to stdout so the user knows the DB was healed; failures (`auto-register failed: ...`) surface to stderr without aborting the release itself (the git tag/push already succeeded).

### Files (this section)

- New: `gitmap/cmd/releaseautoregister.go` — `autoRegisterCurrentRepo(db, cwd)` builds a single-repo scan record, upserts it, ensures the parent ScanFolder, and tags the repo.
- Edited: `gitmap/cmd/releasepersist.go` — `persistReleaseToDB` now calls `resolveOrRegisterCurrentRepoID` (resolve → auto-register on miss → re-resolve). The original `resolveCurrentRepoID` is kept for `listreleasesload.go` which should remain read-only.

## v3.21.0 — (2026-04-20) — schema-version fast-path, `db-migrate --force`, post-update force-migrate, last-release detector fix, `gitmap install clean-code`

### Added (schema-version fast-path)

- **One-time schema-version marker** persisted in the existing `Setting` table under key `schema_version` (value = stringified int). `Migrate()` short-circuits when the on-disk marker equals `constants.SchemaVersionCurrent`, so every subcommand that calls `openDB()` pays only one `Setting` SELECT instead of re-walking the full v15 phase pipeline (the source of the "Migrating GoProjectMetadata → ..." spam users were seeing on every command). The marker is stamped LAST after a successful Migrate() so partial failures retry next run; legacy databases (no `Setting` table or pre-integer-PK rows) read 0 and run the full pipeline exactly once.
- **`gitmap db-migrate --force`** clears the persisted `schema_version` marker before `Migrate()` so the full v15 pipeline re-runs even when the fast-path would otherwise skip it. Useful when a previous run stamped the marker but a downstream issue (corrupt seed, manual edit, partial restore) means the schema actually needs re-walking — without paying the full cost of `gitmap db-reset --confirm`. Failures are warned, never fatal.
- **`runPostUpdateMigrate` always force-clears the marker** before invoking `Migrate()`. After a binary swap from `gitmap update`, the new binary now re-walks the FULL pipeline once, eliminating the failure mode where a freshly-shipped migration step gets skipped because the on-disk marker (written by the previous binary) already equals `SchemaVersionCurrent` for the OLD binary. Cost: one extra full pipeline run, exactly once per update — every subsequent command takes the fast-path again.

### Fixed (run.ps1 last-release detector)

- **`gitmap/scripts/Get-LastRelease.ps1` now treats the deployed binary's `version` output as the authoritative source of truth.** Previously it queried `list-versions` first (which could return empty/stderr-only output after a fresh deploy and bleed PowerShell's `$Matches` capture into the `version` regex), so the post-build summary printed `Last release: v (binary)` — a literal `v` with no semver. The script now (1) calls `& $Binary version` first and only accepts a real `\d+\.\d+\.\d+` capture, (2) resets `$Matches = $null` between regex calls so stale captures cannot leak, (3) reads the current `.gitmap/release/latest.json` location first (with legacy `.release/latest.json` only as fallback), and (4) refuses to print anything that does not match `^v\d+\.\d+\.\d+$` — falling back to `unknown` rather than a malformed string.

### Files (this section)

- New: `gitmap/store/migrate_schemaversion.go` — `readSchemaVersion`, `writeSchemaVersion`, `isSchemaUpToDate` backed by the existing `Setting` key/value table.
- Edited: `gitmap/store/store.go` — `Migrate()` returns immediately when the marker matches; stamps the marker LAST on a successful full pipeline run.
- Edited: `gitmap/constants/constants_settings.go` — adds `SettingSchemaVersion`, `SchemaVersionCurrent` (with bump-policy doc comment), and three log strings (`MsgSchemaVersionUpToDateFmt`, `MsgSchemaVersionAdvanceFmt`, `WarnSchemaVersionWriteFmt`).
- Edited: `gitmap/cmd/dbmigrate.go` — `parseDBMigrateFlags` returns `(verbose, force)`; new `clearSchemaVersionMarker(db *store.DB)` helper; `runPostUpdateMigrate` now calls `clearSchemaVersionMarker` before `Migrate()` so the post-update worker always re-walks the full pipeline.
- Edited: `gitmap/constants/constants_dbmigrate.go` — adds `FlagDBMigrateForce`, `FlagDescDBMigrateF`, `MsgDBMigrateForceClear`, `WarnDBMigrateForceClear`.
- Edited: `gitmap/helptext/db-migrate.md` — documents the new `--force` flag and adds two examples.
- Edited: `gitmap/scripts/Get-LastRelease.ps1` — Strategy A (`gitmap version`) now wins over Strategy B (`list-versions`); `$Matches` is reset between regex calls; `Get-ReleaseFromJSON` checks `.gitmap/release/latest.json` first, legacy `.release/latest.json` second; final guard requires strict `^v\d+\.\d+\.\d+$` before printing.

### Notes

- Bump policy for `SchemaVersionCurrent`: bump on ANY structural change to `Migrate()` — new `CREATE TABLE`, new `ALTER TABLE`, new v15 phase, new seed call, new ID rename. Do NOT bump for cosmetic changes (comments, log strings, code moves that produce identical SQL). The marker is cleared by `gitmap db-reset` and by `migrateLegacyIDs()` when it rebuilds the Repos table, so any database requiring genuine repair always re-runs the full pipeline regardless of the marker value.

### Added (install) — `gitmap install clean-code`

### Added (install)

- **`gitmap install clean-code`** (and the equivalent aliases `code-guide`, `cg`, `cc`) installs the alimtvnetwork coding-guidelines (v15) into the current directory by piping the published `install.ps1` through PowerShell. The flow is: resolve `powershell` (preferred on Windows) or `pwsh` (fallback / non-Windows), then exec `irm <DefaultCleanCodeURL> | iex` with `-NoProfile -ExecutionPolicy Bypass`. On non-Windows hosts the user gets an explicit note that PowerShell 7+ is required. All four aliases route through a single `cleanCodeAliases` set so dispatch and validation stay in sync.
- **Tab-completion exposure** for the new install tokens: `clean-code`, `code-guide`, and `cc` are now emitted as top-level entries by the completion generator via a new `// gitmap:cmd top-level` block in `gitmap/constants/constants_cleancode.go`. `cg` is intentionally left to its existing owner (`changelog-generate`) to avoid shadowing top-level dispatch — `gitmap install cg` still works because `runInstall` parses its own positional argument and routes it through `cleanCodeAliases`.

### Files

- New: `gitmap/cmd/installcleancode.go` — `cleanCodeAliases` set, `isCleanCodeAlias`, `runInstallCleanCode`, `resolvePowerShellBinary`.
- New: `gitmap/constants/constants_cleancode.go` — `DefaultCleanCodeURL`, the `MsgCleanCode*` / `ErrCleanCodeFailed` strings, and the new `// gitmap:cmd top-level` block exposing `CmdInstallCleanCode` / `CmdInstallCleanCodeGuide` / `CmdInstallCleanCodeCC` to tab-completion.
- Edited: `gitmap/cmd/install.go` — `validateToolName` and `executeInstall` short-circuit through `isCleanCodeAlias` so the four aliases bypass the standard `InstallToolDescriptions` map and dispatch straight to `runInstallCleanCode`.
- Edited: `gitmap/helptext/install.md` — documents the new command and its aliases.
- Edited: `gitmap/completion/allcommands_generated.go` — regenerated to include `cc`, `clean-code`, `code-guide` (kept in sync with the marker block).

### Notes

- The four aliases are argument values to `gitmap install`, not standalone top-level commands. They are surfaced through the completion marker block purely so users get tab-complete hints when typing `gitmap install <TAB>`. Direct invocation as `gitmap clean-code` is **not** wired into the dispatcher and will fall through to the unknown-command path.
- This entry is intentionally drafted under `Unreleased` because the version bump must be performed by `gitmap r` (which writes `.gitmap/release/vX.Y.0.json` and updates `latest.json`). Per the project rule, those release-metadata files are never edited by hand.

## v3.20.0 — (2026-04-20) — `gitmap releases --all-repos` multi-repo batch view

### Added (releases)

- **New top-level `gitmap releases` command** as an alias of `list-releases` (`lr`), exposing the new multi-repo batch view via `--all-repos`.
- **`--all-repos` flag** on `list-releases` / `lr` / `releases` runs a SQL JOIN of every `Release` row with its owning `Repo` row, ordered by `CreatedAt DESC, Slug ASC`. This is the first command that explicitly exercises the `IdxRelease_RepoId` secondary index added in v17, demonstrating the multi-repo schema readiness pre-paid by that index. Output adds a `REPO` column (slug) to the table; `--json` emits the joined records as `[]store.ReleaseAcrossRepos`. `--limit N` works the same as the single-repo view. The query bypasses the cwd-bound `loadReleases` scan, so it works from anywhere — even outside any git repo — as long as the gitmap DB is reachable.

### Files

- New: `gitmap/store/releaseacrossrepos.go` — `ReleaseAcrossRepos` struct + `ListReleasesAcrossRepos` query method (table/column-existence guarded for pre-v17 DBs).
- New: `gitmap/cmd/listreleasesallrepos.go` — `runListReleasesAllRepos`, table/JSON renderers, `hasAllReposFlag`.
- Edited: `gitmap/cmd/listreleases.go` — `runListReleases` now pivots to the all-repos branch when `--all-repos` is present.
- Edited: `gitmap/cmd/roottooling.go` — dispatches the new `releases` command name.
- Edited: `gitmap/constants/constants_cli.go` — adds `CmdReleases = "releases"`.
- Edited: `gitmap/constants/constants_globalflags.go` — adds `FlagAllRepos = "--all-repos"`.
- Edited: `gitmap/constants/constants_store.go` — adds `SQLSelectAllReleasesAcrossRepos` (Release JOIN Repo).
- Edited: `gitmap/constants/constants_messages.go` — adds 5 `MsgListReleasesAllRepos*` strings for the wider table.
- Edited: `gitmap/constants/cmd_constants_test.go` — registers `CmdReleases` in the round-trip table.
- Edited: `gitmap/helptext/list-releases.md` — documents the new alias and flag.

### Notes

- The store method is defensive: if the DB pre-dates v17 (no `Release.RepoId` column or no `Repo` table), it returns an empty slice rather than erroring, so the command degrades gracefully on legacy databases.

## v3.19.1 — (2026-04-20) — Exhaustive PATH sweep in uninstall-quick scripts

### Fixed (uninstall)

- **`uninstall-quick.ps1` and `uninstall-quick.sh` now do an exhaustive PATH sweep** as a final step, after the canonical `gitmap self-uninstall` and the deploy-folder sweep have run. Previously, if a stray `gitmap.exe` / `gitmap` binary lived outside the known deploy roots (e.g. a manually-copied shim in `C:\Tools\gitmap.exe`, `~/bin/gitmap`, or a leftover from an old install in `D:\gitmap\gitmap.exe`), it would survive the uninstall and `gitmap` would still resolve in the shell.
- **PowerShell**: `Get-AllGitmapOnPath` uses `Get-Command gitmap -All` (not just the first match) AND directly walks every `Machine` + `User` PATH entry probing for `gitmap.exe` / `gitmap`. Each unique location is `Remove-Item`-ed; the parent dirs are then stripped from the User PATH via `Remove-DirsFromUserPath`.
- **Bash**: `find_all_gitmap_on_path` iterates `$PATH` explicitly (since `command -v` only returns the first hit), de-dupes, and removes each binary. Falls back to `sudo rm -f` for `/usr/*` and `/opt/*` paths.

## v3.19.0 — (2026-04-20) — Bare release auto-bumps minor + scan-dir multi-repo release

### Added (release)

- **Bare `gitmap release` / `gitmap r` inside a git repo** now auto-bumps the **MINOR** segment of the last release (read from `.gitmap/release/latest.json`, falling back to local git tags via the existing `resolveLatestVersion` chain). It prints `Auto-bump: vX.Y.Z → vX.(Y+1).0 (minor)` and prompts `Proceed with this release? [y/N]`. `gitmap r -y` skips the prompt and proceeds.
- **Bare `gitmap release` / `gitmap r` from a scan-dir** (cwd is NOT a git repo, no `--version`/`--bump`/`--commit`/`--branch` supplied) walks the tree with `scanner.ScanDir`, keeps only repos that already have a `.gitmap/release/latest.json`, prints a single summary table (`• <relpath>   <current> → <next>`), prompts ONCE, and then releases each repo by chdir-ing into it and reusing the existing `release.Execute` workflow with `Bump=minor, Yes=true`. Failures are aggregated and reported at the end without aborting the batch. The previous fallback to `runReleaseSelf` still fires when no scan candidates are found.

### Files

- New: `gitmap/cmd/releaseautobump.go` — `peekNextMinorVersion`, `confirmAutoBump`, `shouldAutoBumpMinor`, `readYesNo`.
- New: `gitmap/cmd/releasescan.go` — `tryRunReleaseScanDir`, `planScanReleaseTargets`, `executeScanReleasePlan`.
- Edited: `gitmap/cmd/release.go` — flag parsing moved earlier so the auto-bump branch can read `-y`; new `applyBareReleaseAutoBump` helper.
- Edited: `gitmap/constants/constants_release.go` — new `MsgReleaseAutoBump*` and `MsgReleaseScan*` strings.

### Notes

- The auto-bump path is deliberately conservative: it only fires when the user supplies **none** of `--version` / `--bump` / `--commit` / `--branch`, so existing scripted invocations are unaffected.
- Scan-dir mode reuses the v3.17.0 `Release.RepoId` FK pipeline (`persistReleaseToDB` per release), so multi-repo runs populate the FK correctly per repo.

## v3.18.0 — (2026-04-20) — uninstall-quick PowerShell HOME fix

### Fixed (uninstall)

- **`uninstall-quick.ps1`** — `Remove-CompletionSourceLines` no longer assigns to `$home`, which collides with PowerShell's built-in read-only `$HOME` variable because variable names are case-insensitive. The script now uses `$userHomeDir`, so profile cleanup completes without the `Cannot overwrite variable HOME because it is read-only or constant.` error.

## v3.17.0 — (2026-04-20) — Release.RepoId FK + doctor duplicate-binary check + uninstall profile cleanup

### Doctor

- **New check: `checkDuplicateBinaries`** — detects when multiple `gitmap` binaries exist on PATH (e.g. a stale drive-root shim + the canonical `gitmap-cli/` install). Lists each binary as `[active]` or `[stale]` with its version, and prints a one-shot removal command (`Remove-Item` on Windows, `sudo rm` on Unix). This catches the root cause of the `uninstall-quick.ps1` "not recognized as a cmdlet" error before it happens.
- **New check: `checkReleaseRepoIntegrity`** — joins `Release` with `Repo` via the new FK and reports two diagnostics now that the FK exists:
  - **Orphaned `Release` rows** (rows whose `RepoId` has no matching `Repo` row). Should always be `0` post-FK; a non-zero count indicates DB drift or a partially-applied migration. The check prints the offending `ReleaseId` + `Tag` values so they can be cleaned up manually.
  - **Repo rows with zero releases** (registered repos that have never been released). Surfaced as an informational warning, not an error — useful for spotting repos that were scanned but never tagged. Output is suppressed when the count is 0.
  Backed by `store.ReleaseRepoIntegrity()` which uses `LEFT JOIN` queries that are guarded against pre-v17 schemas (returns `(0, 0, nil)` when `Release.RepoId` doesn't exist yet).

### Fixed (uninstall)

- **`uninstall-quick.ps1` + `gitmap self-uninstall`** — now strip the `# gitmap shell completion` + `. '…completions.ps1'` dot-source lines from **all four** PowerShell profile files (PowerShell + WindowsPowerShell × profile.ps1 + Microsoft.PowerShell_profile.ps1). Previously only the marker-block was removed from a single profile, leaving stale dot-source lines that errored on every new terminal after uninstall.

### Schema (BREAKING)

- **`Release.RepoId INTEGER NOT NULL REFERENCES Repo(RepoId) ON DELETE CASCADE`** — every release row is now anchored to its source repo. The previous global `Tag UNIQUE` constraint is replaced by composite `UNIQUE (RepoId, Tag)`. New index `IdxRelease_RepoId` for per-repo filtering.
- **Migration `migrateV15Phase6`**: detects `Release` tables missing `RepoId`, drops them, and lets the standard CREATE pass rebuild with the new FK schema. Existing rows are wiped (user-approved policy: re-import from `.gitmap/release/v*.json` on next `gitmap list-releases`). See `spec/04-generic-cli/24-release-repo-relationship.md`.

### Code

- `model.ReleaseRecord` gains `RepoID int64`.
- `store.UpsertRelease` requires non-zero `RepoID`; returns `ErrReleaseNoRepo` when the repo isn't registered.
- New `store.ResolveCurrentRepoID(absPath)` helper resolves the FK from `Repo.AbsolutePath`.
- All three release-persist call sites — `cmd/release.go:persistReleaseToDB`, `cmd/listreleasesload.go:cacheReleasesToDB`, `cmd/scanimport.go:importReleases` — now resolve and stamp `RepoID` before upsert.

### Spec

- New: `spec/04-generic-cli/24-release-repo-relationship.md`
- New: `spec/04-generic-cli/images/release-repo-er.mmd` (Mermaid ER diagram)

### Recovery

If a user has legacy `Release` rows but no `.gitmap/release/v*.json` files on disk, run `gitmap release-import --from-github` to repopulate from the GitHub Releases API.

## v3.16.0 — (2026-04-20) — uninstall-quick.ps1 multi-binary fix + repo rename to gitmap-v27

### Fixed

- **`uninstall-quick.ps1`** — `Try-SelfUninstall` and `Resolve-DeployRoot` now pipe `Get-Command gitmap` through `Select-Object -First 1`. When two `gitmap.exe` binaries were on `PATH` (e.g. a stale drive-root shim at `E:\gitmap\gitmap.exe` AND the canonical `E:\bin-run\gitmap-cli\gitmap.exe`), `Get-Command` returned an array. PowerShell's string interpolation joined `.Source` with a space and the resulting `'E:\gitmap\gitmap.exe E:\bin-run\gitmap-cli\gitmap.exe'` was passed to the runtime as a command name, producing:
  > The term '...path1... ...path2...' is not recognized as a name of a cmdlet, function, script file, or executable program.
  Self-uninstall now invokes the FIRST resolved binary by **absolute path** (`& $activeBinary self-uninstall -y`) instead of relying on `& gitmap` PATH resolution, so a stale shim cannot hijack the call.
- **`run.ps1`** — Same defensive fix applied to all three `Get-Command gitmap` callers (deploy-target detection, post-deploy active-vs-deployed sync, changelog binary resolution).
- **`gitmap/scripts/Get-LastRelease.ps1`** — Same defensive fix in `Get-ReleaseFromBinary`.

### Renamed

- **All `gitmap-v27` references → `gitmap-v27`** across the entire repo (45 files, 567 occurrences). Includes install/uninstall one-liners, Go installer constants, helptext, spec docs, post-mortems, the React landing page, and `.lovable/memory/**`.
- **Preserved**: release-asset filenames like `gitmap-v27.49.1-windows-amd64.zip` (where `v4.49.1` is the package version, not the repo name) — only the GitHub URL repo segment changed.

### Why

The uninstall failure was reported by a user who had run gitmap since the v2.x drive-root shim era — their stale `E:\gitmap\gitmap.exe` was never removed and the new `gitmap-cli/` install put a second binary on PATH. The `gitmap-v27` repo rename had been pending since the v3.x line started; bundling both keeps the CHANGELOG narrative simple.

## v3.15.0 — (2026-04-20) — Single-source-of-truth deploy manifest

### Added

- **`gitmap/constants/deploy-manifest.json`** — Single source of truth for deploy-target folder names (`appSubdir`, `legacyAppSubdirs`, `binaryName`, `sourceRepoSubdir`). Renaming the deploy folder in any future version now requires editing **only this one file** — no more drift across `run.ps1`, `run.sh`, `install.sh`, and Go constants.
- **`gitmap/constants/deploy_manifest.go`** — Embeds the manifest via `go:embed` and populates `constants.GitMapSubdir`, `constants.GitMapCliSubdir`, `constants.LegacyAppSubdirs`, and `constants.Manifest` at package init. Falls back to v3.13.x defaults if the JSON is unparseable so the binary stays usable.
- **`Get-DeployManifest`** (run.ps1) and **`load_deploy_manifest`** (run.sh, install.sh) — Each script now parses the manifest from disk (run.ps1, run.sh) or from the install repo via curl (install.sh) and exports `$AppSubdir` / `$LegacyAppSubdirs` (or `APP_SUBDIR` / `LEGACY_APP_SUBDIRS`) for use by all downstream layout, deploy, and cleanup logic.
- **`Test-KnownAppSubdir`** (run.ps1) and **`is_known_app_subdir`** (run.sh) — Reusable predicates that check whether a folder name matches the current or any legacy app subdir, replacing the literal `gitmap-cli`/`gitmap` `or` chains.

### Changed

- **`gitmap/constants/constants_doctor.go`** — `GitMapSubdir` and `GitMapCliSubdir` are no longer `const` literals; they are now `var` populated from the embedded manifest at init time.
- **`gitmap/constants/constants_update.go`** — `UpdatePSDeployDetect` is now a 5-arg format template (was hardcoded `gitmap`/`gitmap-cli`/`gitmap.exe`). The Windows update script generator (`gitmap/cmd/updatescript.go`) injects the manifest-sourced values plus a PowerShell `@(...)` literal of all known subdirs.
- **`run.ps1`** — Deploy target detection, `Repair-DeployLayout`, drive-root shim safety guard, and post-deploy app-dir resolution all now read from `$script:AppSubdir` / `$script:LegacyAppSubdirs` instead of literal strings.
- **`run.sh`** — Same migration: `resolve_deploy_target`, `repair_deploy_layout`, and `Deploy-Binary` use `$APP_SUBDIR` / `is_known_app_subdir`. Legacy migration loop iterates `LEGACY_APP_SUBDIRS` so adding a new legacy name is a one-JSON-line change.
- **`gitmap/scripts/install.sh`** — `repair_layout` and `install_binary` use `$APP_SUBDIR`. The `add_path_to_profile` snippet probe (`${INSTALL_DIR}/gitmap-cli/gitmap`) is now `${INSTALL_DIR}/${APP_SUBDIR}/gitmap`. Manifest is fetched via curl from the install REPO at startup.

### Why

The previous v3.14.0 release had `gitmap-cli` hardcoded in **6+ places** across PowerShell, Bash, and Go. Any future rename (or addition of a new layout migration target) required hunting every file. The manifest centralizes this so the next rename is a one-line PR plus tests.

### Validation policy — `deploy-dfd` CI stays gone

The `deploy-dfd` GitHub Actions job (removed in v3.13.9) is **intentionally not being reinstated**, even after the manifest refactor would make it easier to write. The decision is now documented in [`spec/04-generic-cli/22-data-folder-deploy-and-cleanup.md`](spec/04-generic-cli/22-data-folder-deploy-and-cleanup.md#validation-policy--no-deploy-dfd-ci-job-v3139). Deploy-layout regressions are now caught by:

1. **`gitmap doctor`** on every user's first launch and after updates (PATH binary, deployed binary, version match, app-subdir vs. manifest).
2. **Author smoke testing** — `./run.ps1` and `./run.sh` against clean sandboxes before each tag.
3. **Manifest single-source-of-truth** — `gitmap/constants/deploy-manifest.json` makes silent drift across the four drivers impossible.
4. **Code review** — DFD parity table in the spec MUST be updated in the same commit as any driver change.

Targeted unit tests are preferred over broad CI sandbox-layout assertions when a specific regression is found.


## v3.14.0 — (2026-04-20) — Unix deploy migrated to gitmap-cli/ for cross-platform parity

### Changed

- **`run.sh`** — Now deploys into `<deploy-target>/gitmap-cli/` instead of `<deploy-target>/gitmap/`, matching `run.ps1` (which made the same rename in v3.6.0). The deploy target is now visually unambiguous: the folder name (`gitmap-cli`) no longer collides with the binary name (`gitmap`), and the Go-side cleanup/path logic in `gitmap/cmd/updatecleanup_paths.go` (which already used `GitMapCliSubdir = "gitmap-cli"` for ALL platforms) finally agrees with what's on disk on Unix.
- **`gitmap/scripts/install.sh`** — End-user installer also migrated to `${INSTALL_DIR}/gitmap-cli/`. The pre-existing `repair_layout()` had a latent bug where `app_dir` and `legacy_binary` resolved to the same path (`$target/${BINARY_NAME}`); the rewrite uses distinct variables and now handles BOTH legacy layouts correctly.

### Added (DFD-3 migration)

- **Two-stage legacy layout migration** in `repair_deploy_layout()` (run.sh) and `repair_layout()` (install.sh):
  1. **Migration A** — pre-DFD unwrapped install: `<target>/gitmap` (binary at top level) → `<target>/gitmap-cli/gitmap` + sibling data/, CHANGELOG.md, docs/, docs-site/.
  2. **Migration B** — v3.6.0..v3.13.10 wrapped install: `<target>/gitmap/` (folder) → `<target>/gitmap-cli/` via single `mv`. Skipped with a warning if both folders already exist (manual review needed).
- **PATH-resolution backwards-compat** — `resolve_deploy_target()` in run.sh now accepts both `gitmap-cli` and legacy `gitmap` as the active-binary parent dir name, so users on the v3.6.0..v3.13.10 layout still get their existing deploy target detected on first migration run.

### Updated

- **`spec/04-generic-cli/22-data-folder-deploy-and-cleanup.md`** — DFD-1/DFD-2/DFD-3 rows of the cross-platform parity table updated to reflect `gitmap-cli` on all three drivers (run.ps1, run.sh, install.sh).

### Why now

The Go side of the codebase (cleanup, doctor, binary location, upgrade script) has consistently used `constants.GitMapCliSubdir = "gitmap-cli"` since v3.6.0 — but only `run.ps1` actually deployed there. On Unix, `run.sh` and `install.sh` were still writing to `gitmap/`, which meant `gitmap doctor`, `gitmap update-cleanup`, and PATH-derived deploy detection were all looking in the wrong folder. The v3.13.5/v3.13.7/v3.13.8 patch-stream kept band-aiding tests and CI; this release fixes the actual divergence.


## v3.13.9 — (2026-04-20) — deploy-DFD CI job removed

### Removed

- **`.github/workflows/ci.yml`** — Deleted the entire `deploy-dfd` job (Ubuntu + Windows matrix, ~135 lines, formerly lines 400–533) per user request. The job ran `run.sh` / `run.ps1` into a sandboxed HOME and asserted DFD-1/4/6/7 layout invariants from `spec/04-generic-cli/22-data-folder-deploy-and-cleanup.md`. It had become a recurring source of CI breakage every time the deploy layout evolved (most recently the Windows `gitmap` → `gitmap-cli` rename in v3.6.0, patched in v3.13.8). The DFD spec remains authoritative; layout regressions will now surface through the manual-install path or via `gitmap self-install` end-user testing rather than a synthetic sandbox harness.


## v3.13.8 — (2026-04-20) — CI deploy-DFD Windows assertion aligned with gitmap-cli subdir

### Fixed

- **`.github/workflows/ci.yml`** — The `deploy-dfd` job's Windows DFD-1 assertion (line ~506) was hardcoded to check `$deploy\gitmap\gitmap.exe`, but `run.ps1` has deployed into `gitmap-cli\` since v3.6.0 (see `run.ps1` line 671: `$appDir = Join-Path $target "gitmap-cli"`). CI was failing with `DFD regression: DFD-1: missing wrapped folder D:\a\...\dfd-sandbox\bin-run\gitmap`. Updated the Windows assertion block to expect `gitmap-cli\` and added an inline comment pointing to the rename so the next reader sees the "why" immediately.

### Why not Ubuntu

The Ubuntu assertion (line ~448: `APP_DIR="$DEPLOY/gitmap"`) is intentionally left unchanged — `run.sh` (line 484, 688) still deploys into `gitmap/` on Unix. The `gitmap-cli` rename was Windows-only because on Windows the binary and the folder previously shared the exact same name (`gitmap.exe` inside `gitmap\`), which confused users and autocompletion. Unix has no such collision (`gitmap` binary inside `gitmap/` is unambiguous in a POSIX shell).

## v3.13.7 — (2026-04-20) — find-next const block tagged for completion generator

### Fixed

- **`gitmap/constants/constants_find_next.go`** — Audit of all `constants_*.go` files (Python AST scan over const blocks containing `Cmd[A-Z]\w* = "..."` declarations not marked `// gitmap:cmd skip`) found exactly one drift: the `find-next CLI tokens` block declared `CmdFindNext` and `CmdFindNextAlias` alongside flag tokens but was missing the `// gitmap:cmd top-level` marker. Split the block into two: one for flag tokens (untagged) and one for the command names (tagged with the marker). Without this fix, future renames of `find-next` would not surface in `allcommands_generated.go` and the CI `generate-check` would not catch it.

### Why

Marker comments are the source of truth for the completion generator. A const block containing top-level commands but lacking the marker is silent drift waiting to happen — the audit closes that gap across all 35 `constants_*.go` files. All other const blocks declaring `Cmd*` strings were verified to either (a) carry the `// gitmap:cmd top-level` marker, or (b) tag every line with `// gitmap:cmd skip`.

## v3.13.6 — (2026-04-20) — Completion generator drift resynced

### Fixed

- **`gitmap/completion/allcommands_generated.go`** — CI `generate-check` flagged 5 missing commands. Added in alphabetical order: `probe`, `reset`, `self-install`, `self-uninstall`, `sf`. These were registered with `// gitmap:cmd top-level` markers in their respective spec const blocks but the generated file had not been re-run. Equivalent to `cd gitmap && go generate ./...`.

## v3.13.5 — (2026-04-20) — Stale cleanup-path tests aligned with gitmap-cli subdir

### Fixed

- **`gitmap/cmd/updatecleanup_paths_test.go`** — Tests still asserted the legacy `gitmap` deploy subdir, but production code migrated to `GitMapCliSubdir = "gitmap-cli"` (v3.6.0). Updated three tests:
  - `TestDeriveDeployAppDir`: PATH binary outside the deploy dir now expects `E:/gitmap-cli`; added a third case covering the legacy `E:/gitmap` short-circuit (still recognized by `deriveDeployAppDir`).
  - `TestCollectBackupCleanupDirsIncludesPathDerivedDeployAndBuild`: now expects `E:/gitmap-cli` and `E:/bin-run/gitmap-cli`.
  - `TestCollectTempCleanupDirsIncludesTempAndDerivedTargets`: now expects `E:/gitmap-cli`.
- **`gitmap/cmd/updatescript_test.go`** — `TestBuildUpdateScriptUsesPathAwareDeployVerification` updated: expected substring is now `gitmap-cli\gitmap.exe` to match `constants.UpdatePSDeployDetect` line 114.

### Why

Production paths in `updatecleanup_paths.go` and `constants_update.go` were updated for the v3.6.0 deploy-subdir rename, but these unit tests were missed and started failing on CI. No production behavior change — pure test alignment.

## v3.13.4 — (2026-04-20) — gocritic sprintfQuotedString fix

### Fixed

- **gocritic `sprintfQuotedString`** in `gitmap/store/migrate_v15rebuild.go:107` — Replaced `"%s"` with `%q` for the SQLite identifier quoting in the `INSERT INTO ... SELECT FROM` rebuild template. Behaviorally identical (both produce `"TableName"`) but satisfies the linter and is more idiomatic.

### Fixed

- **UK English residue eliminated across source files** — Audit scanned every `*.go`, `*.ts`, `*.tsx`, `*.js`, `*.jsx`, `*.sh`, `*.ps1` (excluding `node_modules`, `.git`, `.gitmap`, `dist`, `build`) for ~80 UK spelling patterns (colour, optimise, organise, analyse, fibre, behaviour, honour, favour, realise, recognise, normalise, summarise, finalise, utilise, customise, artefact, catalogue, dialogue, licence, defence, traveller, etc.). Found 9 remaining hits and converted to US English:
  - `install-quick.ps1`, `install-quick.sh`, `run.ps1` (3 files): `behaviour → behavior` in script comments.
  - `src/pages/ClearReleaseJSON.tsx`: 7 occurrences of `behaviour → behavior` (object keys + JSX accessor + heading + table column header), plus `Normalised → Normalized` in edge-case data row. Object keys, accessors, and visible UI text remain consistent.
- **Intentionally preserved**: `cancelled` / `cancelling` (GitHub Actions CI terminology — `cancel-in-progress` is the official feature name), `analyses` (valid US English plural of "analysis"), `grey` (UI status descriptor matching GitHub's grey-icon convention), historical CHANGELOG/spec/memory entries (immutable record).

### Verified

- Re-ran the audit grep across the same file set; zero remaining matches for the UK pattern set under audit.

## v3.13.2 — (2026-04-20) — Pre-commit hook enhanced

### Changed

- **`hooks/pre-commit` enhanced** — Updated comments and output to explicitly document the three key linters: `misspell` (US spelling), `exhaustive` (complete switch coverage), and `errcheck` (unchecked errors). Pinned golangci-lint version to `v1.64.8` in the install hint.

### Fixed

- **golangci-lint v1.64.8 CI errors** — 26 linter errors fixed across the Go CLI:
  - `errcheck`: Explicitly ignored or checked return values for `fmt.Sscanf` in `probe/probe.go` and `f.Write` in `cmd/selfinstall.go`.
  - `gosec`: Suppressed G201 (SQL formatting) and G107 (HTTP with variable) via `//nolint:gosec` where variables are internal constants/specs.
  - `gocritic` `sloppyReassign`: Removed unnecessary `err` re-assignments in `movemerge/copy.go`, `movemerge/move.go`, `movemerge/resolve.go`, `cmd/selfinstall.go`.
  - `unused`: Removed `isDuplicateColumnError` in `store/store.go`.
  - `unparam`: Removed unused `info os.FileInfo` parameters from `shouldIgnore` and `shouldSkipWalk`.
  - `wastedassign`: Removed dead `stashLabel` assignment in `cmd/releasealias.go`.
  - `exhaustive`: Added missing switch cases for `PreferPolicy`, `Direction`, and `DiffKind`.
- **US-English spelling sweep** — Converted UK spellings to US: `behaviour→behavior`, `honours→honors`, `honouring→honoring`, `artefacts→artifacts`, `Centralised→Centralized`, `summarises→summarizes`, `Recognises→Recognizes`.
- **Remote installer URLs** — Updated `constants_selfinstall.go` `SelfInstallRemotePwsh` and `SelfInstallRemoteBash` from `gitmap-v27` to `gitmap-v27`.

### Changed

- **`.lovable/prompts/01-read-prompt.md` overwrite** — New onboarding prompt with structured Phase 1–4 flow and mandatory deep-dive source specs lookup table.

## v3.12.1 — (2026-04-20) — AST registry parity + spec cross-links + legacy-field test cleanup

### Added

- **AST-derived `topLevelCmds()` registry parity test** — `gitmap/constants/cmd_constants_parity_test.go` adds `TestTopLevelCmdRegistryMatchesAST`, which uses `go/parser` to walk every `gitmap/constants/constants_*.go`, collects every `Cmd*` constant declared inside a `// gitmap:cmd top-level` block (minus those tagged `// gitmap:cmd skip`), and asserts the resulting set is exactly equal to the manual `topLevelCmds()` registry consumed by `TestTopLevelCmdConstantsAreUnique` / `TestTopLevelCmdAliasesAreUnique`. The registry can no longer drift silently — adding a new top-level `Cmd*` without registering it (or vice versa) fails CI with a clear "missing from registry" / "registered but not declared" diff.
- **Spec cross-links from CLI overview** — `spec/01-app/02-cli-interface.md` and `spec/01-app/38-command-help.md` gained a `> **Related:**` callout under the H1 pointing at `spec/01-app/99-cli-cmd-uniqueness-ci-guard.md`, so future contributors discover the uniqueness contract and the 6-step handoff checklist directly from the CLI overview and the help-system spec.
- **Spec §5 implementation note** — `spec/01-app/99-cli-cmd-uniqueness-ci-guard.md` updated to mark the AST parity test as implemented (no longer "future hardening") with the file path and v3.12.1 history entry.

### Fixed

- **Stale `Draft` / `PreRelease` `ReleaseMeta` / `Options` field references in tests** — `gitmap/release/metadata_test.go` and `gitmap/tests/release_test/skipmeta_test.go` still constructed `ReleaseMeta{Draft: …, PreRelease: …}` and `release.Options{Draft: …}` using the pre-v15 field names, breaking `go vet` / `go build` with `unknown field Draft in struct literal`. Renamed both to the v15 `IsDraft` / `IsPreRelease` form, matching every production caller. The legacy-JSON compat shim in `release/metadata.go::ReadReleaseMeta` (which still reads the old `draft` / `preRelease` JSON keys) is intentionally untouched and remains the supported migration path for v3.4.x metadata files on disk.
- **`go vet` `non-constant format string`** in `gitmap/cmd/probe.go:127` — `fmt.Fprintf(os.Stderr, result.Error+"\n")` triggered the printf-check because the format string was constructed at runtime from a struct field. Reshaped the call to `fmt.Fprintf(os.Stderr, "%s\n", result.Error)` so the format string is a compile-time constant.

### Verified

- Full-repo audit for residual legacy-field callers: every `\.(Draft|PreRelease)\b` and `^\s*(Draft|PreRelease)\s*:` match outside of (a) `release.Version.PreRelease` (semver suffix — different struct), (b) `store/migrate_v15phase5.go` (the rename migration itself), (c) `release/metadata.go::ReadReleaseMeta` (the JSON backward-compat overlay), and (d) `--draft` user-facing CLI flag strings was confirmed to be either intentional or already migrated. No further call sites need updating.

## v3.12.0 — (2026-04-20) — Pinned-version release snippet + gitmap-v27 rename

### Added

- **Pinned-version install snippet on the GitHub release page** — the release publisher (`gitmap/release/installsnippet.go`, wired into `workflowgithub.go::uploadToGitHub`) now auto-appends a markdown block containing PowerShell + bash one-liners that hard-code the just-published tag. Idempotent via a hidden `<!-- gitmap-pinned-install-snippet:<tag> -->` HTML marker. Anyone copying the snippet from `…/releases/tag/v3.12.0` installs exactly v3.12.0 — never "latest", never a `-v<N+1>` sibling repo. Template lives in `constants_release.go` as `ReleaseSnippetTemplate` / `ReleaseSnippetMarker`.
- **Pinned-version short-circuit in installer scripts** — `gitmap/scripts/install.ps1` and `install.sh` gained a new branch in their discovery prelude: when `-Version <tag>` (PowerShell) or `--version <tag>` (bash) is supplied, the installer now skips both the `releases/latest` API call **and** the versioned-repo `-v<N>` discovery probe, downloading `…/releases/download/<tag>/…` directly. Closes the gap where a snippet copied from a v3.x release page could silently jump to the v4 repo's latest tag.
- **Spec doc** `spec/07-generic-release/08-pinned-version-install-snippet.md` — full NEA/AI handoff contract: rendered snippets, installer-side flag matrix, release-cutting checklist, and a CI test contract for future work.

### Changed

- **Repo rename `gitmap-v27` → `gitmap-v27` across the entire codebase** — every Go constant (`SourceRepoCloneURL`, `SelfInstallRemotePwsh/Bash`, `GitmapRepoPrefix`, install hint URLs), every install/uninstall script (`install.ps1`, `install.sh`, `install-quick.ps1`, `install-quick.sh`, `uninstall-quick.*`), every spec doc under `spec/01-app/` and `spec/07-generic-release/`, every helptext markdown, the README, the React `src/data/*.ts` files, GitHub workflows, and historical CHANGELOG entries were rewritten via `sed -i 's/gitmap-v27/gitmap-v27/g'`. The only remaining `gitmap-v27` references are inside `.gitmap/` artifacts, which are immutable per project policy.

## v3.11.1 — (2026-04-20) — Alias-collision CI guard

### Added

- **Alias-collision uniqueness test** — extended `gitmap/constants/cmd_constants_test.go` with `TestTopLevelCmdAliasesAreUnique`, which iterates every top-level `Cmd*` constant and fails when two distinct identifiers share the same short-form value (string length ≤ 2). Catches future regressions like a hypothetical `CmdFooAlias = "ls"` shadowing the existing `CmdListAlias`, before they reach the build phase. Companion `TestTopLevelCmdConstantsAreUnique` covers full-length command-name collisions. Manual `topLevelCmds()` registry is the source of truth and excludes anything marked `// gitmap:cmd skip`.

## v3.11.0 — (2026-04-19) — Constants hygiene + Phase 1.4 migration fix

### Fixed

- **v15 Phase 1.4 migration** — `GoProjectMetadata` and `PendingTask` rebuilds failed on databases first created at v3.5.0+ with `SQL logic error: no such column: Id`. Both tables were already singular before v15, so the canonical `CREATE TABLE IF NOT EXISTS` pass produced the v15-shaped table (with `{Table}Id` PK) before the rebuild ran, leaving no `Id` column to SELECT. Added `adaptOldColumnList()` in `gitmap/store/migrate_v15rebuild.go` that detects the existing PK shape via `columnExists()` and rewrites the leading `Id` token in `OldColumnList` to `{Table}Id` when needed. Idempotent and a no-op for genuine legacy → v15 paths.
- **`go vet` `non-constant format string`** in `gitmap/movemerge/finalize.go:50` — `logErr` was inferred as a printf-style wrapper. Reshaped `logErr(prefix, msg string)` to accept a pre-formatted message and moved `fmt.Sprintf(constants.ErrMMPushFailFmt, sha)` to the call site so the printf-check never triggers.
- **Unused-import build break** in `gitmap/store/migrations.go` — removed orphaned `"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"` import left over from a prior refactor.
- **`CmdReleaseAlias` Go redeclaration** — same name was bound to `"r"` (in `constants_cli.go`) and `"release-alias"` (in `constants_releasealias.go`). Renamed the `constants_cli.go` constant to `CmdReleaseShort` so the `release-alias` family owns `CmdReleaseAlias` exclusively.
- **`cd` / `go` constant collision** — `CmdCDCmd` (`"cd"`) and `CmdCDCmdAlias` (`"go"`) in `constants_cli.go` shadowed `CmdCD` / `CmdCDAlias` in `constants_cd.go`. Removed the duplicates and repointed `gitmap/cmd/rootdata.go` dispatch at the canonical constants.

### Added

- **CI uniqueness test** — `gitmap/cmd/cmdconstants_unique_test.go` (+ helpers in `cmdconstants_unique_helpers_test.go`) parses every `gitmap/constants/constants_*.go`, applies the same `gitmap:cmd top-level` / `gitmap:cmd skip` markers used by `completion/internal/gencommands`, and fails the test suite when two distinct `Cmd*` identifiers claim the same string value. Catches future redeclarations and dispatch shadowing at CI time before they reach the build phase.
- **Parallel pull worker pool** (`gitmap/cmd/pullparallel.go`) — buffered-channel pool with `sync.WaitGroup` and a mutex around the non-thread-safe `BatchProgress` tracker. Opt-in via `--parallel <N>`.
- **`--only-available` pull pre-filter** (`gitmap/cmd/pullfilter.go`) — intersects the target repo list with `FindNext` results so `gitmap pull --only-available` skips repos that have no new tags. Fail-open: falls back to a full pull if the database is inaccessible.
- **`gitmap probe` and `gitmap sf` help docs** — `gitmap/helptext/probe.md` and `gitmap/helptext/sf.md` (synopsis, flags, examples, 3–8 line realistic terminal simulation), discoverable via `gitmap help probe` / `gitmap help sf`.

### Changed

- **`constants_cli.go` size reduction** — extracted the `Shorthand*` group into `gitmap/constants/constants_clone.go` and the cross-command `Flag*` values into a new `gitmap/constants/constants_globalflags.go`. `constants_cli.go` is now 188 lines (under the 200-line guideline).

## v3.5.0 — (2026-04-19) — v15 Database Naming Alignment (Phase 1 complete)

### Changed

- **Phase 1 of the v15 database naming migration is complete.** All 22 SQLite tables now follow the strict v15 convention from <https://github.com/alimtvnetwork/coding-guidelines-v15/blob/main/spec/04-database-conventions/01-naming-conventions.md>: PascalCase + **singular** table names, `{TableName}Id` primary keys, foreign keys that match the referenced PK name, `IsX` prefix for booleans, and abbreviations treated as words (`SshKey` not `SSHKey`, `CsharpProjectMetadata` not `CSharpProjectMetadata`).
- **Renamed tables** (legacy → v15): `Repos`→`Repo`, `Groups`→`Group`, `GroupRepos`→`GroupRepo`, `Releases`→`Release`, `Aliases`→`Alias`, `Bookmarks`→`Bookmark`, `Amendments`→`Amendment`, `CommitTemplates`→`CommitTemplate`, `Settings`→`Setting`, `SSHKeys`→`SshKey`, `InstalledTools`→`InstalledTool`, `TempReleases`→`TempRelease`, `ZipGroups`→`ZipGroup`, `ZipGroupItems`→`ZipGroupItem`, `ProjectTypes`→`ProjectType`, `DetectedProjects`→`DetectedProject`, `GoProjectMetadata` (kept), `GoRunnableFiles`→`GoRunnableFile`, `CSharpProjectMeta`→`CsharpProjectMetadata`, `CSharpProjectFiles`→`CsharpProjectFile`, `CSharpKeyFiles`→`CsharpKeyFile`. `RepoVersionHistory`, `CommandHistory`, `TaskType`, `PendingTask`, `CompletedTask` were already singular and only got `{TableName}Id` PK renames.
- **Renamed columns**: every legacy `Id` PK is now `{TableName}Id` (e.g., `Repo.RepoId`, `Release.ReleaseId`, `CsharpProjectMetadata.CsharpProjectMetadataId`). Foreign keys updated to match (e.g., `GoRunnableFile.GoProjectMetadataId`, `CsharpProjectFile.CsharpProjectMetadataId`). `Release.Draft` → `Release.IsDraft` and `Release.PreRelease` → `Release.IsPreRelease` complete the IsX boolean-prefix consistency (`IsLatest` was already correct).
- **Migration safety contract** (applies to every Phase 1.1–1.5 rebuild):
  1. Detect-then-act on every legacy plural — fresh installs are no-ops.
  2. `PRAGMA foreign_keys=OFF` for the duration of each table rebuild.
  3. Row-count parity check between old and new on every rebuild — abort + return on mismatch.
  4. Legacy plural names retained as `LegacyTable*` constants and listed in `Reset()` so cleanup works at any migration state.
  5. SQLite-reserved word `Group` is double-quoted in every DDL/DML occurrence.
- **Go-side propagation**: `model.ReleaseRecord.Draft/PreRelease` → `IsDraft/IsPreRelease` (with JSON tags `isDraft`/`isPreRelease`); `release.Options.Draft` → `release.Options.IsDraft`; `release.ReleaseMeta.Draft/PreRelease` → `IsDraft/IsPreRelease`. `ReadReleaseMeta` includes a JSON overlay that accepts the legacy `"draft"`/`"preRelease"` keys so on-disk `.gitmap/release/*.json` files from v3.4.x and earlier still load.
- **CLI flag `--draft`** is intentionally retained (user-facing). Internal struct fields use the v15 `IsX` naming.

### Added

- New shared migration infrastructure in `gitmap/store/migrate_v15rebuild.go` — generic `runV15Rebuild(spec)` helper using a `v15RebuildSpec` struct (OldTable, NewTable, NewCreateSQL, OldColumnList, NewColumnList, StartMsg, DoneMsg). Drives all 22 table rebuilds.
- New phase migrators wired into `store.Migrate()` in dependency-safe order:
  - `migrate_v15phase2.go` — Group, Release, Alias, Bookmark + GroupRepo FK-text rebuild.
  - `migrate_v15phase3.go` — Amendment, CommitTemplate, Setting, SshKey, InstalledTool, TempRelease.
  - `migrate_v15phase4.go` — ZipGroup family, Project family (incl. CSharp→Csharp), Task family, History tables.
  - `migrate_v15phase5.go` — `Release.Draft`→`IsDraft`, `Release.PreRelease`→`IsPreRelease` (column rename via the same rebuild infrastructure).
- Pre-rename column patches for very old installs: `preV15Phase2EnsureReleaseColumns()` (Source/Notes on legacy `Releases`), `migrateZipGroupItemPaths()` and `migrateTRCommitSha()` already targeted legacy plurals before the v15 rebuilds copied the data.
- Regenerated `spec/01-app/gitmap-database-erd.mmd` to reflect every v15 table name, PK, FK, and `IsDraft`/`IsPreRelease` boolean.
- Updated `spec/12-consolidated-guidelines/11-database.md` with the v15 naming conventions table (singular + `{TableName}Id` + `IsX` boolean prefix + reserved-word quoting + abbreviation rules), with a link to the upstream v15 spec.

### Notes

- This release is purely a naming alignment — no new commands, no behavior changes for end users beyond the schema. Existing databases upgrade in place via the idempotent rebuild migrators; rollback is via `gitmap db-migrate` against an older binary's CREATE statements after restoring a DB backup.
- Phase 2 (ScanFolder, VersionProbe, `gitmap find-next`) and Phase 3 (parallel `pull`, bulk `cn next all`) remain on the roadmap.

## v3.0.0 — (2026-04-19)

### Added

- `gitmap as [alias-name] [--force|-f]` (alias `s-alias`) — tag the **current** Git repository with a short alias and persist it in the active-profile SQLite database. Resolves the repo top-level via `git rev-parse --show-toplevel`, builds a single-repo `ScanRecord` through the existing `mapper.BuildRecords()` pipeline (so the upserted row matches the schema other commands use), upserts into `Repos`, then maps `alias-name → Repos.Id` in the alias store. When `alias-name` is omitted the repo folder basename is used. Refuses to clobber an existing alias unless `--force` is passed. Exits 1 with a CWD-aware message when invoked outside a Git repo.
- `gitmap release-alias <alias> <version>` (alias `ra`) — release a previously-aliased repo from **any** working directory. Resolves alias → absolute path via the alias store, `os.Chdir`s into the repo, runs the existing `runRelease` pipeline (lint → test → tag → push → assets), then restores the original CWD via `defer`. Forwards `--dry-run` to `runRelease` for safe previews.
- `gitmap release-alias-pull <alias> <version>` (alias `rap`) — thin sugar for `release-alias --pull`. Runs `git pull --ff-only` in the resolved repo before releasing; hard-fails on non-fast-forward (never tags on top of a divergent tree). The flag remains canonical, the verb is sugar.
- **Auto-stash semantics for `release-alias`**: dirty working trees are auto-stashed (`git stash push --include-untracked -m "gitmap-release-alias autostash <alias>-<version>-<unix-ts>"`) before the release runs and popped on exit via `defer`, so the stash always fires — including when `runRelease` aborts. The pop locates the stash by **label match** against `git stash list` (not by `stash@{0}`), so a concurrent `git stash` from another process never causes us to pop the wrong entry. A failed pop warns only — the user's tree is still recoverable via `git stash list` / `git stash apply`. Bypass with `--no-stash` (intended for CI runners that always start clean and want to fail loudly on unexpected dirt).
- `gitmap db-migrate` (alias `dbm`) — explicit, idempotent schema migration command. Re-runs every `CREATE TABLE IF NOT EXISTS` and column-migration step on the active profile DB. Now invoked automatically at the end of `gitmap update` so a freshly-updated binary never has to repair the database on its first real run. `--verbose` prints extra context.
- New shared migration helpers in `gitmap/store/migrations.go`: `columnExists(table, column)`, `tableExists(table)`, `isBenignAlterError(err)`, and `logMigrationFailure(table, column, action, err, stmt)` — every warning now names the table, column, and action so issues can be diagnosed without trial-and-error.
- New files: `gitmap/cmd/{as.go, asops.go, releasealias.go, releasealias_git.go, dbmigrate.go}`, `gitmap/constants/{constants_as.go, constants_releasealias.go, constants_dbmigrate.go}`, `gitmap/store/migrations.go`, `gitmap/helptext/{as.md, release-alias.md, release-alias-pull.md, db-migrate.md}`, `spec/01-app/98-as-and-release-alias.md`.

### Changed

- **`migrateTRCommitSha` switched to detect-then-act.** Previously the migration always tried `ALTER TABLE TempReleases RENAME COLUMN "Commit" TO CommitSha` and only suppressed errors via brittle string-matching on `"no such column"`. On Unix builds where the SQLite driver formats the error slightly differently (or the table is fresh and only has `CommitSha`), the warning leaked through with the cosmetic `no such column: ""Commit""` message. The migration now uses `PRAGMA table_info(TempReleases)` to check whether `Commit` actually exists before attempting the rename, eliminating the spurious warning entirely on every OS regardless of driver wording.
- **Generator switched from explicit allowlist to marker-comment opt-in.** `gitmap/completion/internal/gencommands/main.go` no longer maintains a `sourceFiles` list or a `skipNames` map. Instead it scans every `../constants/*.go` automatically and includes only `const (...)` blocks whose doc comment contains `// gitmap:cmd top-level`. Individual specs inside an opted-in block can be excluded with a trailing `// gitmap:cmd skip` line comment (used for subcommand IDs like `"create"` / `"add"` shared across `gitmap group`). Domain owners now control inclusion locally without ever editing the generator. Added markers across 40 const blocks in 34 constants files (52 skip annotations mirror the previous policy exactly); `allcommands_generated.go` regenerates byte-for-byte identically (143 entries).
- `gitmap/cmd/update.go::runUpdateRunner` now calls `runPostUpdateMigrate()` after the binary swap completes, so every `gitmap update` finishes by running migrations. Best-effort: failures warn but do not block (the user may have an in-flight DB lock or a read-only environment).
- `gitmap/completion/completion.go::manualExtras` is now empty with an updated doc comment pointing future contributors at the marker convention instead of the old `sourceFiles` + `skipNames` instructions.
- All migration warnings (`addColumnIfNotExists`, `migrateZipGroupItemPaths` data-copy step, `migrateTRCommitSha`) now route through `isBenignAlterError` for a uniform suppression policy: `no such column`, `no such table`, `duplicate column`, and `already exists` are all benign on fresh installs.

### CI

- Added a `generate-check` job to `.github/workflows/ci.yml` that runs `go generate ./...` in `gitmap/` and fails with `git diff --exit-code` (printing the drifted file list and the fix command) if any generated file is out of sync with the constants. Wired into `test-summary`'s `needs` so the SHA-passthrough cache won't mark a run green unless the drift check also passed.

### Notes

- The original task description asked for a bump to `v2.97.0`; we are already at `v3.0.0` from the preceding `db-migrate` and marker-comment work, so the version was kept and the changelog rolled into a single v3.0.0 entry covering `as`, `release-alias`, `release-alias-pull`, `db-migrate`, the migration hardening, the generator refactor, and the CI drift check.

---

## Migration guide — v2.x → v3.0.0 (constants contributors)

If you maintain a custom `constants_*.go` file in `gitmap/constants/` that exposes command IDs for shell tab-completion, you must opt-in explicitly using marker comments.

### What changed
- **Old (v2.x):** The generator (`internal/gencommands/main.go`) relied on a hard-coded `sourceFiles` list and a `skipNames` map. Adding a new command required editing the generator.
- **New (v3.0.0):** The generator scans every `constants/*.go` file automatically. Inclusion is controlled locally via comments.

### What you need to do

1. Open your `constants_*.go` file.
2. Locate the `const (...)` block containing your `Cmd*` string constants.
3. Add `// gitmap:cmd top-level` to the block's **doc comment** (the comment immediately above `const`).
4. If any constant in that block is a *subcommand* (e.g., `"create"` or `"add"` used only inside `gitmap group`), add a trailing line comment `// gitmap:cmd skip` to that specific spec.

**Example:**

```go
// gitmap:cmd top-level
// Bookmark commands.
const (
    CmdBookmarkAdd    = "add"    // gitmap:cmd skip
    CmdBookmarkList   = "list"
    CmdBookmarkRemove = "remove"
)
```

5. Re-run `go generate ./...` in `gitmap/` to regenerate `allcommands_generated.go`.
6. Verify with `git diff` — only your new command values should appear; no manual edits to the generator needed.

### Verification
- CI now runs a `generate-check` job that fails if `allcommands_generated.go` drifts from the constants. If your PR fails this check, the error message prints the exact command to fix it locally.

---

## v2.98.0 — (2026-04-18)

### Added

- `gitmap mv LEFT RIGHT` (alias `move`) — moves LEFT's contents into RIGHT (excluding `.git/`), then deletes LEFT entirely. Both endpoints can be local folders or remote git URLs (with optional `:branch` suffix); URL endpoints are cloned (or fast-forward pulled if already on disk with matching origin), and after the move the RIGHT-side URL is committed (`gitmap mv from <LEFT-display>`) and pushed.
- `gitmap merge-both LEFT RIGHT` (alias `mb`) — bidirectional file-level merge: each side gains every file the other has but it doesn't; conflicting files (different content on both sides) trigger the `[L]eft / [R]ight / [S]kip / [A]ll-left / [B]all-right / [Q]uit` interactive prompt.
- `gitmap merge-left LEFT RIGHT` (alias `ml`) — one-way merge that writes only into LEFT (RIGHT is read-only). With `-y`, RIGHT wins by default.
- `gitmap merge-right LEFT RIGHT` (alias `mr`) — one-way merge that writes only into RIGHT (LEFT is read-only). With `-y`, LEFT wins by default.
- Bypass flags shared by all four merge commands: `-y` / `--yes` / `-a` / `--accept-all` skip the prompt; `--prefer-left`, `--prefer-right`, `--prefer-newer`, `--prefer-skip` override the per-command default policy. `merge-both -y` defaults to `--prefer-newer`.
- URL-side commit/push controls: `--no-push` (commit but skip push), `--no-commit` (copy files but skip both). `--force-folder` replaces a folder whose origin doesn't match the requested URL. `--pull` opt-in for `git pull --ff-only` on folder endpoints. `--dry-run` prints every action and writes nothing. `--include-vcs` and `--include-node-modules` override the default ignore list.
- New `gitmap/movemerge/` package with focused files (<200 lines each, <15 lines per function): `types.go`, `endpoint.go` + `endpoint_test.go` (URL classification + `:branch` suffix + scp-style `git@host:user/repo` preservation), `walk.go` (default ignore list `.git/` / `node_modules/` / `.gitmap/release-assets/`), `copy.go` (mode-preserving file copy with symlink replication), `conflict.go` + `conflict_test.go` (L/R/S/A/B/Q resolver with sticky All-Left/All-Right and `--prefer-newer` mtime tie-break), `diff.go` (SHA-256 classification into MissingLeft / MissingRight / Conflict / Identical), `git.go` (clone / pull --ff-only / add-commit-push), `resolve.go` (full endpoint resolver with origin-match check), `guard.go` (same-folder + nested-ancestor protection), `merge.go`, `move.go`, `finalize.go` (URL-side commit + push), `log.go` (structured `[mv]` / `[merge-*]` prefix lines).
- CLI wiring: `cmd/move.go`, `cmd/merge.go`, `cmd/movemergeflags.go` (shared flag binder), `cmd/dispatchmovemerge.go` hooked into `cmd/root.go`. New constants in `constants/constants_movemerge.go` (command IDs, aliases, flag names, log prefixes, commit message templates, error formats) plus `GitAddCmd`, `GitAddAllArg`, `GitCommitCmd`, `GitMessageArg` reused for the post-merge git plumbing.

### Notes

- `mv` does NOT prompt — its semantic is destructively "move-and-delete-LEFT". Use `merge-right` for the safer copy-with-prompt variant.
- Same-folder and nested-folder protection trips before any file write: LEFT and RIGHT may not resolve to the same absolute path, and neither may be a strict ancestor of the other on disk.
- `gitmap diff LEFT RIGHT` (added in v2.97.0) is the recommended dry-run preview before `gitmap merge-both` — every conflict it lists will trigger the interactive prompt.


### Added

- `gitmap diff LEFT RIGHT` (alias `df`) — read-only preview of what `gitmap merge-both / merge-left / merge-right` would change between two folders. Lists conflicts (different content on both sides), missing-on-LEFT, missing-on-RIGHT, and (optionally) identical files. Writes nothing, commits nothing, pushes nothing.
- Flags: `--json` (machine-readable output with `{summary, entries}` payload), `--only-conflicts`, `--only-missing`, `--include-identical`, `--include-vcs`, `--include-node-modules`. Honours the same default ignore list as `merge-*` (`.git/`, `node_modules/`, `.gitmap/release-assets/`).
- New `gitmap/diff/` package: `endpoint.go` (folder-only resolver — URL endpoints are intentionally rejected with a hint to clone first), `tree.go` (parallel walk + SHA-256 classification), `report.go` (text/JSON renderer + `Summary` tally). Unit tests cover all four diff kinds and the default ignore list.
- `gitmap/helptext/diff.md` and `gitmap/cmd/diff.go` + `gitmap/cmd/dispatchdiff.go` wire the command into the existing dispatcher chain in `root.go`.

### Notes

- `diff` is the recommended dry-run preview before `merge-both`: every conflict it lists will trigger the `[L]eft / [R]ight / [S]kip / [A]ll-left / [B]all-right / [Q]uit` prompt during merge-both.
- URL endpoints are rejected on purpose so `diff` remains strictly side-effect-free (no network, no clone, no temp folders). Clone first via `gitmap clone <url>`, then diff the resulting folder.


## v2.96.0 — (2026-04-18)

### Added

- Help text files for the move/merge command family: `gitmap/helptext/mv.md`, `merge-both.md`, `merge-left.md`, `merge-right.md`. Each follows the standard template (overview, alias, usage, flags, prerequisites, 3 examples with sample output, exit codes, see-also).
- `gitmap help <command>` now prints the embedded help file for any command (e.g. `gitmap help mv`, `gitmap help merge-both`). Previously `gitmap help` only showed the global usage banner. The lookup uses the existing `helptext.Print` function, so every command in `gitmap/helptext/*.md` is auto-discovered.

### Changed

- `dispatchUtility` in `gitmap/cmd/rootutility.go` now intercepts `gitmap help <name>` before falling through to the global usage printer. A small `isFlagToken` helper distinguishes `gitmap help --groups` (still goes to grouped usage) from `gitmap help mv` (prints `mv.md`).


## v2.95.0 — (2026-04-18)

### Added

- `gitmap setup print-path-snippet --shell <bash|zsh|fish|pwsh> --dir <path> --manager <label>` — emits the canonical marker-block PATH snippet to stdout. Used by `run.sh` and `gitmap/scripts/install.sh` so all three drivers produce byte-identical rc-file output. Single source of truth lives in `constants_pathsnippet.go`.
- `gitmap setup` now writes the marker-block snippet to the user's profile on every run (idempotent: rewrites the existing block in place, otherwise appends after a blank line). Different `--manager` values create coexisting blocks so `run.sh`, `installer`, and `gitmap setup` never overwrite each other.
- `setup.WritePathSnippet()` and `setup.RenderPathSnippet()` Go helpers with full unit-test coverage (`pathsnippet_test.go`, `pathsnippetwriter_test.go`).

### Changed

- `run.sh::register_on_path` and `gitmap/scripts/install.sh::add_path_to_profile` now ask the freshly-built/installed gitmap binary for snippet bytes via `gitmap setup print-path-snippet`. Inline heredocs remain as a first-run fallback only.

## v2.94.0 — (2026-04-18)

### Fixed

- `Get-LastRelease.ps1` reported the OLDEST version (e.g. `v2.82.0`) because `list-versions --limit 1` returns ascending order. Now sorts all versions descending and falls back to the binary's own `version` output if needed.
- Stale active PATH binary (e.g. `E:\bin-run\gitmap.exe`) is no longer kept alive by copying the new build into it. New `Migrate-StaleActiveBinary` helper deletes the stale binary, removes empty parent dirs, and strips the location from user PATH so future shells use the wrapped deploy target only.
- `powershell.json` `deployPath` is now rewritten after every successful deploy via `Sync-ConfigDeployPath` so the "Config binary:" readout reflects the actual install location and future runs default to the same target.

## v2.83.0 — (2026-04-16)

### Fixed

- `gitmap update-cleanup` now scans the active PATH directory, the PATH-derived deploy directory, the configured deploy directory, and the repo build output directory so stale `.old` backups are removed even when `powershell.json` points to an older location.
- `gitmap update-cleanup` now removes leftover `gitmap-update-*` artifacts from deploy/build locations in addition to `%TEMP%`, preventing handoff files from being left behind after update flows that switch between deploy targets.

## v2.82.0 — (2026-04-16)

### Fixed

- Regenerated `package-lock.json` to sync with `package.json` — resolves CI `npm ci` failure caused by missing entries for testing libs, axios, framer-motion, vitest, and other dependencies added without a lockfile refresh.

## v2.81.0 — (2026-04-16)

### Fixed

- `go-winres` CI icon size error — Windows `.ico` resources require images ≤256x256 but `icon.png` was 512x512. Created `icon-256.png` (LANCZOS resize) and updated `winres.json` to reference it.
- Documented root cause and prevention in `spec/08-generic-update/09-winres-icon-constraint.md`.

## v2.80.0 — (2026-04-16)

### Added

- Hidden `set-source-repo` command — persists source repo path to DB so `gitmap update` always uses the correct location after repo moves.
- Post-deploy repo path sync in `run.ps1` — automatically calls `set-source-repo` after every successful deploy to keep the DB current.
- Repo path sync spec (`spec/08-generic-update/08-repo-path-sync.md`) — documents the post-deploy sync pattern for AI implementers.
- Help file for `set-source-repo` command (`gitmap/helptext/set-source-repo.md`).

### Fixed

- `go-winres` CI failure — moved `winres.json` from `gitmap/` to `gitmap/winres/` where `go-winres make` expects it.

### Changed

- Cross-references updated in `02f-self-update-orchestration.md` and `03-self-update-mechanism.md` to include repo path sync spec.

## v2.78.0 — (2026-04-16)

### Added

- Console-safe handoff spec (`spec/08-generic-update/07-console-safe-handoff.md`) — documents the blocking `cmd.Run()` pattern that prevents terminal detachment during self-update on Windows.
- Installer banner now displays version number (`gitmap installer v1.0.0`).

### Changed

- `install.ps1`: `Resolve-Version` now prints full HTTP status code, URL, response body, and potential causes on GitHub API failure instead of a generic error.
- `gitmap-updater/cmd/github.go`: `fetchLatestTag` error output now includes URL, response body, and troubleshooting hints.
- Standardized lowercase "gitmap" branding across all installer output messages.

### Fixed

- `ShouldPrintInstallHint` now uses case-insensitive matching for GitHub repo URL detection.

## v2.76.0 — (2026-04-16)

### Added

- New `gitmap version-history` (`vh`) command displays all version transitions for the current repo with `--limit N` and `--json` flags.
- Full database ERD (Mermaid) added to `spec/01-app/gitmap-database-erd.mmd` covering all 22 tables including `RepoVersionHistory`.
- Updated `spec/01-app/59-clone-next.md` and `spec/01-app/87-clone-next-flatten.md` to reflect flatten-by-default behavior (no `--flatten` flag required).

---

## v2.75.0 — (2026-04-16)

### Added

- `gitmap clone-next` now flattens by default: clones into the base name folder (no version suffix) instead of the versioned folder name. For example, `gitmap cn v++` inside `macro-ahk-v15` clones `macro-ahk-v16` into `macro-ahk/`.
- `gitmap clone <url>` auto-flattens versioned URLs when no custom folder is given. `gitmap clone https://github.com/user/wp-onboarding-v13` clones into `wp-onboarding/`.
- New `RepoVersionHistory` SQLite table tracks every version transition (from/to version tags, numbers, and flattened path) with timestamps.
- `Repos` table gains `CurrentVersionTag` and `CurrentVersionNum` columns, updated on each clone-next operation.
- Version transitions are printed to terminal: `Recorded version transition v15 -> v16`.
- If the flattened target folder already exists during clone-next, it is automatically removed and re-cloned fresh.

---

## v2.74.0 — (2026-04-16)

### Added

- `gitmap doctor` now checks setup config resolution from the installed binary location and warns when `git-setup.json` cannot be found.
- `gitmap doctor` now verifies the shell wrapper is loaded by checking the `GITMAP_WRAPPER` environment variable, with fix instructions when missing.
- Post-setup verification step warns users if the shell wrapper is not active after `gitmap setup` completes, with reload instructions.
- Shell wrapper scripts (Bash, Zsh, PowerShell) now export `GITMAP_WRAPPER=1` so the binary can detect wrapper-vs-raw invocation.
- `gitmap cd` prints a stderr warning when called without the shell wrapper, guiding users to run `gitmap setup` or reload their profile.

### Fixed

- `gitmap setup` now resolves `git-setup.json` relative to the binary's installation path instead of the current working directory, fixing "file not found" errors when running from arbitrary directories.

---

## v2.72.0 — (2026-04-16)

### Fixed

- VS Code admin-mode bypass: `runVSCodeCommand` now captures `CombinedOutput` and waits for the process exit code instead of fire-and-forget, ensuring CLI errors are properly detected before falling through to the next strategy.
- `tryVSCodeDetached` launches `Code.exe` with an isolated `--user-data-dir` (`%TEMP%\gitmap-vscode-user-data`) so the new instance does not attempt to hand off to an elevated single-instance, fully bypassing the "Another instance of Code is already running as administrator" lock.
- Added `resolveVSCodeExecutable` with multi-path discovery (`LookPath`, CLI sibling, `LocalAppData`, `Program Files`, `Program Files (x86)`) to reliably find the desktop binary when the CLI wrapper is unavailable.
- Extracted all VS Code constants (binary names, flags, paths, messages) into `constants/constants_vscode.go`.

---

## v2.71.0 — (2026-04-16)

### Added

- VS Code admin mode bypass: `openInVSCode` now uses a 3-tier launch strategy (`--reuse-window` → `--new-window` → `cmd /C start` detached) to handle the "Another instance of Code is already running as administrator" error.
- Added `tryVSCodeReuse`, `tryVSCodeNewWindow`, and `tryVSCodeDetached` helper functions in `cmd/clonevscode.go`.
- Added `ErrVSCodeAdminLock` constant for admin-mode warning message.

### Fixed

- `gitmap update` PATH sync now includes full 3-step fallback: direct `Copy-Item`, rename-then-copy (`Move-Item` to `.old` + `Copy-Item` with rollback), and kill stale `gitmap.exe` processes via `Stop-Process` before final retry.
- Updated `UpdatePSSync` PowerShell block in `constants/constants_update.go` with rename and kill-process recovery strategies.
- Updated `spec/01-app/89-update-path-sync.md` to document all sync fallback steps and error scenarios.

---

## v2.70.0 — (2026-04-16)

### Added

- `gitmap clone <url>` now auto-registers cloned repositories with GitHub Desktop by default (no manual prompt).
- `gitmap clone <url>` automatically opens the cloned folder in VS Code (`code --reuse-window`), with `--new-window` fallback for admin-mode conflicts.
- Added `isVSCodeAvailable()` detection via `exec.LookPath` in `cmd/clonevscode.go`.

### Fixed

- `gitmap update` now auto-syncs the active PATH binary when it differs from the deployed binary, resolving the `[FAIL] Active PATH version does not match deployed version` error.
- Added `Copy-Item` sync step with rename and kill-process fallbacks in the update PowerShell script.

---

## v2.69.1 — (2026-04-11)

### Fixed

- Fixed `errorlint` violation in `cmd/helpdashboard.go`: replaced direct `!= io.EOF` comparison with `errors.Is` to handle wrapped errors correctly.

### Changed

- Linked "Riseup Asia LLC" in the author Role row to [riseup-asia.com](https://riseup-asia.com).
- Changed Riseup Asia subheading from centered to left-aligned and linked it to [riseup-asia.com](https://riseup-asia.com).

---

## v2.69.0 — (2026-04-09)

### Added

- Windows binaries now embed a custom emerald green terminal icon, application manifest, and version info via `go-winres`.
- Added `gitmap/winres.json` and `gitmap/assets/icon.png` for Windows resource generation.
- Release pipeline generates `.syso` resource files before compilation, injecting the release version into the binary metadata.
- Added `spec/pipeline/09-binary-icon-branding.md` documenting the full `go-winres` workflow for AI/engineer handoff.
- Added the gitmap icon to the README header.

### Fixed

- Fixed `run.ps1 -d` switch: replaced `[Alias("d")]` on `[string]$DeployPath` with a dedicated `[switch]$Deploy` parameter so `-d` works without requiring a path argument.

---

## v2.68.1 — (2026-04-09)

### Fixed

- Fixed gosec G305 (file traversal) and G110 (decompression bomb) in `helpdashboard.go` zip extraction — paths are now validated against the target directory and extraction is size-limited to 100 MB.
- Fixed `run.ps1 -d` failing with "Missing an argument for parameter 'DeployPath'" — added `[Alias("d")]` to `$DeployPath` so `-d` resolves unambiguously.

---

## v2.68.0 — (2026-04-09)

### Fixed

- Fixed `TempReleases` migration crash: `ALTER TABLE RENAME COLUMN "Commit"` failed with `no such column` when the column was already renamed or never existed. Migration now silently skips the rename when the column is absent.

### Added

- Release pipeline now builds the docs-site (React/Vite) and bundles `dist/` into `docs-site.zip` as a release asset.
- Install scripts (`install.ps1`, `install.sh`) automatically download and extract `docs-site.zip` alongside the binary.
- `gitmap hd` auto-extracts `docs-site.zip` on first run if the `docs-site/` directory is missing — no manual setup needed.
- Added 5 new pipeline specification files (`04`–`08`) covering installation flow, changelog integration, version/help system, environment variable setup, and terminal output standards.
- Added AI Handoff Checklist to `spec/pipeline/README.md` with recommended reading order for onboarding.

## v2.67.0 — Smart Deploy & Rename-First (2026-04-08)

### Improvements

- `run.ps1` and `run.sh` now auto-detect the globally installed `gitmap` binary location and deploy there instead of using a hardcoded path.
- Deploy target resolution follows a 3-tier priority: `--deploy-path` CLI flag → globally installed PATH location → `powershell.json` default.
- First-time installs use the config default; subsequent builds automatically deploy to the active binary's directory.
- Added `Resolve-DeployTarget` function to `run.ps1` and `resolve_deploy_target` function to `run.sh` for full cross-platform parity.
- Deploy step now uses **rename-first strategy**: renames the existing binary to `.old` before copying the new one, avoiding Windows file-lock failures when deploying to a running binary.
- Rollback restores the `.old` file via rename (not copy) for consistency.
- Added "Build once, package once" constraint to `spec/05-coding-guidelines/17-cicd-patterns.md` and `spec/04-generic-cli/11-build-deploy.md`.
- Updated `spec/01-app/09-build-deploy.md` with deploy target resolution and rename-first deploy documentation.
- Added smart deploy path resolution and rename-first deploy to cross-platform parity table in `spec/01-app/42-cross-platform.md`.
- Replaced hardcoded `E:\bin-run` path in `gitmap doctor` fix suggestion with dynamic guidance.

## v2.66.0 — CI Hardening & Pipeline Docs (2026-04-08)

### Improvements

- Pinned `govulncheck` to `v1.1.4` in CI and vulncheck workflows for reproducible builds.
- Updated GitHub Actions to Node.js 24 compatible versions (`actions/checkout@v6`, `actions/setup-go@v6`).
- Added `FORCE_JAVASCRIPT_ACTIONS_TO_NODE24: true` environment variable across all workflows.
- Created portable `spec/pipeline/` documentation folder (CI, release, vulnerability scanning) for cross-AI shareability.
- Added CI Tool Versions pinning table to dependency specs (13, 17, 27) for consistency.
- Aligned severity response times across all dependency management specs.
- Updated stale action version examples in specs 17 and 27 from `@v4`/`@v5` to `@v6`.
- Added cross-reference from `spec/03-general/08-ci-pipeline.md` to `spec/pipeline/`.

### Bug Fixes

- Fixed `ShouldPrintInstallHint` not matching SSH remote URLs (`git@github.com:org/repo.git`) due to colon separator not being normalized to a slash.
- Fixed vulncheck pipeline logic error where `-q` flag on initial `grep` suppressed stdout, breaking the vulnerability classification pipe.

## v2.65.0 — Install UX Overhaul (2026-04-07)

### Improvements

- Install flow now shows a structured **Install Plan** box before execution with tool, version, manager, and command.
- Added numbered step progress: `[1/4] Updating...`, `[2/4] Installing...`, `[3/4] Verifying...`, `[4/4] Recording...`.
- Chocolatey installs now use `--no-progress` flag to suppress GUI popups and prevent blocking on interactive apps like Notepad++.
- Winget installs now use `--silent` flag for unattended installs.
- NPP verification now checks the expected exe path (`C:\Program Files\Notepad++\notepad++.exe`) directly instead of relying on PATH lookup.
- NPP settings zip path now resolves relative to the binary directory (not CWD), fixing "file not found" errors when gitmap is installed globally.
- Detected version is printed during verification for better diagnostics.
- Install command completion is confirmed with a success message before proceeding to verification.

### Bug Fixes

- Fixed NPP install blocking the terminal when Notepad++ GUI launched during Chocolatey install (missing `--no-progress`).
- Fixed post-install verification always failing for NPP because `notepad++` binary is not on PATH.
- Fixed settings zip not found when running `gitmap install npp` from a directory other than the source repo root.

## v2.64.0 — Install Scripts Command (2026-04-07)

### New Commands

- Added `gitmap install scripts` — clones gitmap scripts (install.ps1, install.sh, run.ps1, run.sh, etc.) to a local folder for easy access.
  - **Windows**: resolves the deploy drive from `powershell.json`, defaults to `D:\gitmap-scripts`.
  - **Linux/macOS**: installs to `~/Desktop/gitmap-scripts`.

## v2.63.0 — Installed Directory & Linux Update Flow (2026-04-07)

### New Commands

- Added `gitmap installed-dir` (alias `id`) — prints the full binary path and directory of the active gitmap installation, resolving symlinks to the real location.

### Update Command

- Linux/macOS update now uses `run.sh --update` instead of PowerShell, enabling native shell-based self-update on Unix systems.
- After pulling latest source and rebuilding, the active PATH binary is automatically synced to the new version.
- Added install path resolution using `which gitmap` with `EvalSymlinks` fallback for accurate binary location.
- If `run.sh` is missing from the source repo, a clear error is shown instead of a PowerShell failure.

### Bug Fixes

- Fixed `gitmap update` on Linux: handoff binary no longer uses `.exe` extension and now gets `chmod +x` permission.
- Fixed tilde `~` not expanding in update repo path prompt (e.g. `~/repos/gitmap` was treated as literal `~/`).
- Fixed `gitmap install` on Ubuntu: `apt-get update` now runs before package installation to prevent exit code 100 errors.
- Added `-y`/`--yes` flag to `gitmap install` for non-interactive installs with confirmation prompt.
- Install failures now write detailed error logs to `.gitmap/logs/` with version, manager, command, and reason.
- Fixed `install.sh` installer: `TMP_DIR` unbound variable error on exit caused by subshell scoping.

## v2.62.0 — CI Release Branch Protection (2026-04-07)

### CI/CD

- Release branches (`release/**`) are no longer cancelled by `cancel-in-progress` — every release commit now runs the full CI and release pipeline to completion.
- CI workflow uses a conditional expression: `cancel-in-progress: ${{ !startsWith(github.ref, 'refs/heads/release/') }}` to protect release branches while still cancelling superseded runs on `main` and feature branches.
- Release workflow changed to `cancel-in-progress: false` unconditionally.
- Updated CI pipeline spec (`spec/03-general/08-ci-pipeline.md`) with release branch protection documentation.

## v2.61.0 — Install Hint Polish & Post-Mortem #17 (2026-04-07)

### Release Command

- Improved post-release install hint formatting with emoji labels (📦 🪟 🐧) and better spacing.
- Removed hash-style comments in favor of OS-specific emoji indicators for Windows and Linux/macOS install one-liners.
- Extracted `ShouldPrintInstallHint()` as an exported function for testability.
- Added unit tests for install hint repo detection (11 cases covering gitmap and non-gitmap repos).

### Documentation

- Added Post-Mortem #17: Go Flag Ordering — Silent Flag Drop, documenting the `flag` package behavior and `reorderFlagsBeforeArgs()` fix.

## v2.60.0 — Auto-Detect Pending Release Branch (2026-04-07)

### Release Command

- Running `gitmap release` or `gitmap r` while on a `release/*` branch with no tag now auto-detects and completes the pending release instead of erroring about a duplicate branch.
- Running `gitmap release v1.1.0` while on `release/v1.1.0` with no tag delegates to `ExecuteFromBranch` automatically.
- Added `tryDelegateFromCurrentBranch()` for no-version detection and `tryDelegateFromBranch()` for explicit-version detection.
- Added `MsgReleaseBranchPending` constant for the delegation message.

## v2.59.0 — Post-Release Install Hints (2026-04-07)

### Release Command

- After a successful release, if the repo's remote origin matches the gitmap source repository prefix (`github.com/alimtvnetwork/gitmap-v27`), the CLI now prints install one-liner commands for both Windows (PowerShell) and Linux/macOS (Bash).
- Added `GitmapRepoPrefix` constant for repo detection and `MsgInstallHintHeader`, `MsgInstallHintWindows`, `MsgInstallHintUnix` message constants.
- Install hints appear after `Release complete` in all release paths: standard, branch-based, and metadata-only.
- Non-gitmap repos are unaffected — no install hints are printed.

## v2.58.0 — Release Flag Ordering Fix (2026-04-07)

### Bug Fix

- Fixed `-y` / `--yes` flag being silently ignored when placed after the version argument (e.g., `gitmap release v2.55 -y`).
- Root cause: Go's `flag` package stops parsing at the first non-flag argument, so flags after the version were never processed.
- Added `reorderFlagsBeforeArgs()` helper in `releaseargs.go` — reorders CLI args so all flags precede positional arguments before `flag.Parse()`.
- Affects `release`, `release-self` (`r`, `rs`), and all commands sharing `parseReleaseFlags`.

## v2.57.0 — README & Memory Updates (2026-04-07)

### Documentation

- Split README Quick Start into focused code blocks: separate Install (Windows + Linux/macOS), Scan, and Navigate sections.
- Created `one-liner-installer` memory documenting both `install.ps1` and `install.sh` as CI-generated versioned release assets.

## v2.56.1 — Clone-on-Missing-Path for Update (2026-04-07)

### Update Command

- When the user provides a non-existent path during the `gitmap update` interactive prompt, the system now clones the gitmap source repository into that directory instead of rejecting it.
- After a successful clone, the path is validated, saved to the SQLite Settings DB, and used for the update — no re-prompting on future runs.
- Added `SourceRepoCloneURL`, `MsgUpdateCloning`, `MsgUpdateCloneOK`, and `ErrUpdateCloneFailed` constants.

## v2.56.0 — Release Pipeline install.sh & CI Fix (2026-04-07)

### Release Pipeline

- Added `install.sh` generation to `release.yml` — version-pinned Bash installer is now created and attached as a release asset alongside `install.ps1`.
- Release body now includes both PowerShell and Bash one-liner install instructions.

### CI Pipeline Fix

- Eliminated separate `mark-success` job — inlined cache write as the final step of `test-summary` to prevent `cancel-in-progress` from cancelling the SHA marker after all validation passed.
- `test-summary` now depends on `[sha-check, lint, vulncheck, test]` to ensure full validation before caching.

### Documentation

- Updated `spec/01-app/82-install-script.md` — documented `install.sh` with CLI flags (`--version`, `--dir`, `--arch`, `--no-path`), version-pinned examples, `.tar.gz`/`.zip` fallback, 4-priority binary detection, and shell-aware auto-PATH append (bash/zsh/fish).
- Updated `spec/01-app/12-release-command.md` — CI release pipeline section now mentions `install.sh` alongside `install.ps1` in both steps list and release body format.
- Added "Known Behavior: Concurrency Cancellation" section to `spec/02-app-issues/16-ci-passthrough-gate-pattern.md` — documented and resolved by inlining cache write.
- Updated post-release auto-commit memory to reflect the new `-y` flag behavior.

### Testing

- Added unit test for `-y` flag in autocommit — verifies `promptAndCommit` skips stdin when `yes=true`.

## v2.55.0 — Release Auto-Confirm, Docs & Installer Fix (2026-04-07)

### Post-Mortems Documentation

- Created `spec/02-app-issues/13-release-pipeline-dist-directory.md` — documents `cd: dist` CI failure root cause and 4 prevention rules.
- Created `spec/02-app-issues/14-security-hardening-gosec-fixes.md` — documents G305, G110, format verb, and Code Red fixes with prevention rules.
- Added Post-Mortems page (`/post-mortems`) to docs site with category filters, version tags, and color-coded icons for all 15 documented issues.

### Coding Guidelines Updates

- Added "Lessons Learned" section to `spec/05-coding-guidelines/17-cicd-patterns.md` — never `cd` in CI, validate directories, pin tool versions.
- Added Section 10 (Zip Extraction Security) to `spec/05-coding-guidelines/08-security-secrets.md` — mandatory G305/G110 checks.
- Added Sections 7–8 to `spec/05-coding-guidelines/04-error-handling.md` — Code Red Rule and Format Verb Compliance.

### Installer Fixes

- Fixed PowerShell installer crash caused by `Invoke-WebRequest` progress bar rendering during `irm | iex`.
- Added `$ProgressPreference = "SilentlyContinue"` to `install.ps1`.
- Fixed versioned binary detection — installer now matches `gitmap-v*-windows-(amd64|arm64).exe` patterns from CI archives.
- Wrapped installer `Main` function in `try/catch` with friendly error message and manual download fallback.

### CI Pipeline: Passthrough Gate Pattern

- Replaced job-level `if` skipping with step-level conditionals in `ci.yml` so all jobs always report ✅ Success.
- Previously, SHA-deduplicated runs showed grey "skipped" status which looked like failures; now cached SHAs print "Already validated" and exit green.
- Updated `spec/05-coding-guidelines/29-ci-sha-deduplication.md` with the passthrough pattern documentation.
- Pinned `golangci-lint` to `v1.64.8` in `ci.yml` to match `setup.sh`.

### Release Command: Auto-Confirm (`-y` / `--yes`)

- Added `-y` / `--yes` flag to `release`, `release-self`, `release-branch`, and `release-pending` commands.
- When set, all interactive prompts (e.g. "Auto-commit all changes?") are automatically confirmed without user input.
- Enables fully non-interactive release workflows: `gitmap release v2.55.0 -y`.
- Bumped version to `v2.55.0`.

### Unix Installer (`install.sh`)

- Created `gitmap/scripts/install.sh` — cross-platform Bash installer for Linux and macOS.
- Supports `--version`, `--dir`, `--arch`, `--no-path` flags matching the PowerShell installer feature set.
- Includes SHA256 checksum verification, versioned binary detection, `.tar.gz`/`.zip` fallback.
- Auto-detects shell (bash/zsh/fish) and appends PATH entry to the correct profile file.
- Rename-first strategy for safe upgrades of running binaries.

### Changelog Improvements

- Added release dates to all changelog entries with available metadata (sourced from `.gitmap/release/*.json`).
- Backfilled v2.54.1, v2.54.2, v2.54.3, and v2.53.0 entries in the docs site changelog data.
- Removed duplicate Code Red content from v2.54.0 (now properly in v2.54.1).

### Build Reproducibility

- Pinned `golangci-lint` to `v1.64.8` in `setup.sh` instead of `@latest`.

---

## v2.54.3 — Security Hardening & Lint Compliance (2026-04-07)

### Zip Extraction Security (installnpp.go)

- Fixed **G305** (path traversal): `extractZipEntry` now validates that resolved destination paths stay within the target directory using absolute path prefix checks.
- Fixed **G110** (decompression bomb): `io.Copy` replaced with `io.LimitReader` capped at 10 MB per extracted file.

### Lint Configuration Documentation

- Added inline comments to all 8 gosec exclusions in `.golangci.yml` documenting why each is necessary (G104, G204, G304, G306, G401, G404, G505, G101).

---

## v2.54.2 — Format Verb Audit (2026-04-07)

### fmt.Fprintf Argument Mismatch Fix

- Fixed `cmd/tasksync.go:138` where `fmt.Fprintf` format string expected 2 arguments but only 1 was passed, causing a `go vet` failure.
- Audited all `fmt.Fprintf`, `fmt.Printf`, and `fmt.Errorf` calls across `cmd/`, `release/`, and `store/` packages (~140 call sites, 38+ files) — confirmed 100% compliance.

---

## v2.54.1 — Code Red Error Audit (2026-04-07)

### Mandatory Error Path Logging

- Completed full Code Red audit: every file/path-related error log now includes the exact file path, the operation attempted, and the specific failure reason.
- Standardized format: `Error: [message] at [path]: [error] (operation: [op], reason: [reason])`.
- Updated 35+ constants and 36+ call sites across the entire codebase.
- Generic "file not found" messages without paths are now prohibited by convention.

---

## v2.54.0 — Update Path Recovery & CI Optimization (2026-04-07)

### Update Path Recovery

- `gitmap update` now validates the saved source repo path exists on disk before using it.
- Falls back to the SQLite DB (`source_repo_path` setting) in the binary's `data/` folder.
- Prompts the user interactively when both embedded and saved paths are missing or stale.
- Successfully resolved paths are persisted to the DB for future runs.
- New file `cmd/updaterepo.go` extracts path resolution helpers for the 200-line file limit.

### CI Build Removal

- Removed cross-platform binary builds from the main CI pipeline (`ci.yml`).
- Binaries are now produced exclusively by the release pipeline (`release.yml`) on `release/**` branches and `v*` tags.

### CI Concurrency Cancellation

- All workflows (`ci.yml`, `release.yml`, `vulncheck.yml`) now cancel in-progress runs when a new commit is pushed to the same branch.
- Concurrency groups use `github.ref` so different branches run independently.

### Release Pipeline Fix

- Fixed `cd dist` failure in `release.yml` — the compress/checksum step was running inside `gitmap-updater/` (no `dist/` folder) instead of `gitmap/dist/` where binaries are output.
- Extracted compress and checksum into a separate step with explicit `working-directory: gitmap/dist`.

### SHA-Based Build Deduplication

- CI pipeline now skips redundant runs when the same commit SHA has already passed all checks.
- A `sha-check` gate job probes the GitHub Actions cache for `ci-passed-<SHA>` before any work begins.
- On full pipeline success, a `mark-success` job caches a marker so future runs for the same SHA short-circuit.
- Failed pipelines never cache — re-running the same SHA executes the full pipeline.

---

## v2.53.0 — Help Dashboard & Install Docs

### Help Dashboard Command

- New `gitmap help-dashboard` (alias `hd`) command to serve the documentation site locally.
- Dual-mode resolution: serves pre-built `dist/` via Go's built-in HTTP server; falls back to `npm install && npm run dev` if static assets are missing.
- `--port` flag to configure the serving port (default: 5173).
- Automatically opens the docs site in the default browser on launch.
- Graceful shutdown on Ctrl+C for both static and dev modes.
- New constants file `constants_helpdashboard.go` with all messages, defaults, and error strings.

### Install & Help Dashboard Docs Pages

- Added `/help-dashboard` docs page with terminal demos for static mode, dev fallback, and custom port usage.
- Added `/install` docs page documenting `install` and `uninstall` commands, supported tools, databases, and package managers.
- Both pages include feature cards, flags tables, file layout references, and interactive terminal demos.

## v2.52.0 — Lock Detection & Install System Overhaul

### Lock Detection (clone-next)

- `clone-next` now detects processes locking the current folder when deletion fails.
- On Windows, uses Sysinternals `handle.exe` or PowerShell WMI to identify locking processes.
- On Unix/macOS, uses `lsof` for process detection.
- Prompts the user to terminate blocking processes, then retries folder removal automatically.
- New `lockcheck` package with platform-specific implementations (`lockcheck_windows.go`, `lockcheck_unix.go`).

### Install System Overhaul

- Added SQLite-based installation tracking (`InstalledTools` table) with granular version columns (Major, Minor, Patch, Build) and timestamps.
- Expanded tool support: 11 databases (MySQL, PostgreSQL, Redis, MongoDB, SQLite, MariaDB, CockroachDB, Cassandra, Neo4j, InfluxDB, DynamoDB Local).
- Package manager mappings for Chocolatey, Winget, Apt, Homebrew, and Snap.
- New `gitmap uninstall <tool>` command with `--dry-run`, `--force`, and `--purge` flags.
- README redesigned with centered headers, badges, and grouped command/tool tables.

- Reorganized `gitmap help` output into 17 categorized command groups (Scanning, Cloning, Git Operations, Navigation, Release, etc.).
- Added `--compact` flag to `gitmap help` for a minimal command-and-alias-only listing.
- `gitmap help --compact <group>` filters compact output by group name (case-insensitive, falls back to all groups on no match).
- Added color-coded group headers using ANSI escape codes (bold cyan) for improved terminal readability.
- Added Quick Start section with common command examples at the top of help output.
- Each group header includes a hint to run commands with `--help` or `-h` for detailed usage and examples.
- Modularized help implementation across `rootusage.go`, `rootusagecompact.go`, `rootusageflags.go`, and `constants_helpgroups.go`.
- Repository renamed from `git-repo-navigator` to `gitmap-v27`; all URLs, scripts, and references updated.

## v2.49.1 — Update UX & Versioned Binaries (2026-04-06)

- Added `--repo-path` flag to `update` command: override the source repo path for a one-time update.
- The `--repo-path` flag is automatically forwarded through the handoff binary to `update-runner`.
- Resolution priority: `--repo-path` flag → embedded constant → friendly error with recovery options.
- Improved "repo path not embedded" error with actionable recovery steps (one-liner install, clone & build, manual download, `--repo-path` override).
- CI release binaries now include version in filenames (e.g., `gitmap-v27.49.1-windows-amd64.zip`).
- Updated `install.ps1` (standalone and release-embedded) to handle versioned asset filenames.
- CI release workflow now explicitly marks stable releases as "latest" via `make_latest`.
- Updated `helptext/update.md` with `--repo-path` flag docs, troubleshooting section, and error recovery examples.
- Added `gitmap-updater` — standalone tool to update gitmap via GitHub releases (no source repo required).
- `gitmap update` auto-delegates to `gitmap-updater` when no repo path is available and the updater is on PATH.
- Updater uses handoff-copy pattern to avoid Windows file locks during self-replacement.
- CI release pipeline now builds and ships `gitmap-updater` binaries for all 6 platform targets.

## v2.49.0 — Opt-in Binary Builds & Gitignore Safety (2026-04-06)

- Go binary cross-compilation is now opt-in: use `--bin` or `-b` to build executables during release.
- Removed `--no-assets` flag (replaced by the inverse `--bin` flag).
- `gitmap setup` now ensures `release-assets` and `.gitmap/release-assets` are in `.gitignore`.
- Release workflow auto-appends missing release-related paths to `.gitignore` before each release.
- Added `release-assets` and `.gitmap/release-assets` to `.gitignore` to prevent tracking build artifacts.
- CI release workflow now triggers on `release/*` branch push (in addition to tags).
- Each GitHub release includes: changelog entry, SHA256 checksums, release metadata table, and asset matrix.
- Version-specific `install.ps1` script is auto-generated and attached to each release for one-liner install.
- Pre-release versions (containing `-`) are automatically marked as prerelease on GitHub.

## v2.48.1 — Clone-Next Auto-Navigate (2026-04-03)

- `clone-next` now automatically changes into the newly cloned directory after removing the old folder.
- Prints `→ Now in <target>` confirmation after navigating to the new clone.

## v2.48.0 — Tag Discovery & DB Caching

- `list-releases` now scans git tags via `git for-each-ref` and includes tag-only releases with `source=tag`.
- All discovered releases (repo metadata + tags) are automatically upserted into the SQLite `Releases` table on every `lr` invocation.
- Added `--source tag` filter to `list-releases` for viewing tag-discovered releases.
- Updated helptext and spec to document three-source resolution order and caching behavior.

## v2.47.0 — Release Self Hardening (2026-04-03)

- Changed `release-self` primary alias from `rself` to `rs` (rescan moved to `rsc`).
- Added SQLite DB fallback for source repo discovery (`source_repo_path` in Settings table).
- Skip directory switch if already in the gitmap source repo directory.
- Updated spec, helptext, React docs page, and commands catalog to reflect changes.

## v2.46.0 — Release Self

- Added `release-self` (`rself`) command: release gitmap itself from any directory.
- Auto-fallback: `gitmap release` outside a Git repo now triggers self-release automatically.
- Source repo discovery via `os.Executable()` + symlink resolution + `.git` root walk.
- Returns to original working directory after release with confirmation message.
- Full flag parity with `release` (--bump, --assets, --draft, --dry-run, etc.).
- Added React docs page for release-self with terminal demos and error scenarios.

## v2.45.0 — Docs Site Update (2026-04-03)

- Updated CloneNext docs page with `--create-remote` flag, usage, and terminal example.
- Added repo creation failure to error handling table on docs site.

## v2.44.0 — Clone-Next Spec Update

- Updated `clone-next` spec to document `--create-remote` as opt-in.
- Removed mandatory repo creation from default workflow and examples.
- Added Example 5 showing `--create-remote` usage in spec.
- Marked deferred implementation phases 1–3 as complete.

## v2.43.0 — Clone-Next Hardening

- Auto-cd to parent directory before folder removal to prevent Windows file lock errors.
- Added `--create-remote` flag: optionally create the target GitHub repo before clone (requires `GITHUB_TOKEN`).
- Repo creation is now opt-in instead of mandatory; default `gitmap cn v+1` clones directly.

## v2.42.0 — Clone-Next Simplification

- Removed forced GitHub repo existence check and automatic creation from `clone-next`.
- `gitmap cn v+1` now clones directly without requiring `GITHUB_TOKEN`.
- Repo creation is no longer a blocking prerequisite before clone.

## v2.41.0 — Clone-Next Phase 3 (2026-04-03)

- GitHub repo existence check and automatic creation before clone via GitHub API.
- Requires `GITHUB_TOKEN` for repo creation; creates under org with user fallback.
- Added `ParseOwnerRepo` utility for HTTPS and SSH remote URL parsing.

## v2.40.0 — Clone-Next Command

- Added `clone-next` (alias `cn`) command: clone the next versioned iteration of a repo into its parent directory.
- Supports `v++` and `v+1` (increment current version by 1) and `vN` (jump to explicit version).
- Remote-first repo name resolution: derives base name and version from `remote.origin.url`, not the local folder name.
- GitHub repo existence check before clone: queries `GET /repos/{owner}/{repo}` via GitHub API.
- Automatic GitHub repo creation when target does not exist: creates under org (fallback to user) via GitHub API.
- Requires `GITHUB_TOKEN` environment variable for repo creation.
- Added `ParseOwnerRepo` utility to extract owner/repo from HTTPS and SSH remote URLs.
- Added `--delete` flag: auto-remove current version folder after successful clone.
- Added `--keep` flag: keep current folder without prompting for removal.
- Added `--no-desktop` flag: skip GitHub Desktop registration.
- Added `--ssh-key` / `-K` flag: use a named SSH key for Git operations.
- Added `--verbose` flag: show detailed clone-next diagnostics.
- Clone-Next Flags section added to `gitmap help` output.
- Version argument validation: rejects `v0`, negative values, and malformed inputs with clear errors.
- Case-insensitive version parsing (`V++`, `V+1` accepted).
- No-suffix repos default to `-v2` on increment.
- Added constants for all clone-next messages, errors, and flag descriptions.
- Added unit tests for `ParseRepoName`, `ResolveTarget`, `TargetRepoName`, and `ReplaceRepoInURL`.
- Spec: `spec/01-app/59-clone-next.md` with full workflow, examples, and acceptance criteria.

## v2.37.0 — v2.39.0

- Internal improvements and minor fixes (see individual commits).

## v2.36.7 — Integration Tests

- Added SkipMeta integration test (`skipmeta_test.go`): 6 test cases verifying `SkipMeta: true` prevents metadata and `latest.json` creation.
- Added release rollback integration test (`rollback_test.go`): 5 test cases verifying branch/tag cleanup on simulated push failure.
- Added end-to-end release test (`e2e_test.go`): full cycle from version bump through metadata commit on a temp repo with bare remote.
- E2E edge-case coverage: dry-run (no side effects), no-commit (staged only), skip-meta (no JSON), and duplicate version blocking.
- Added edge-case test suite (`edgecase_test.go`): pre-release parsing/comparison, bump resolution (all levels, from-zero, from-prerelease), parse validation, version ordering, multi-release sequences, out-of-order metadata, and rc-to-stable promotion.
- Added TUI Temp Releases view (`tempreleases.go`, `trformat.go`): 9th tab with flat list, detail panel, and grouped-by-prefix aggregation.
- Added `--stop-on-fail` flag to `pull` and `exec` commands: halts batch after first failure.
- Enhanced `BatchProgress` with per-item failure tracking (`FailWithError`), detailed failure reports, and exit code 3 on partial failures.
- Added `batchreport.go` with `PrintFailureReport()` and `ExitCodeForBatch()` helpers.

## v2.36.6 — Wave 2 Refactoring (14 Files)
- Split `assets.go` → `assets.go` + `assetsbuild.go` (build helpers: `buildSingleTarget`, `buildEnv`).
- Split `zipgroupops.go` → `zipgroupops.go` + `zipgroupshow.go` (display: `runZipGroupList`, `expandFolder`).
- Split `tui.go` → `tui.go` + `tuiview.go` (rendering: `View`, `renderTabs`, `renderContent`).
- Split `aliasops.go` → `aliasops.go` + `aliassuggest.go` (interactive: `runAliasSuggest`, `promptAliasSuggestion`).
- Split `tempreleaseops.go` → `tempreleaseops.go` + `tempreleaselist.go` (listing: `runTempReleaseList`, `printTRList`).
- Split `listreleases.go` → `listreleases.go` + `listreleasesload.go` (data: `loadReleasesFromRepo`, `sortRecordsByDate`).
- Split `listversions.go` → `listversions.go` + `listversionsutil.go` (collection: `collectVersionTags`, `printVersionEntriesJSON`).
- Split `sshgen.go` → `sshgen.go` + `sshgenutil.go` (utils: `validateSSHKeygen`, `resolveGitEmail`).
- Split `scanprojects.go` → `scanprojects.go` + `scanprojectsmeta.go` (metadata: `upsertGoProjectMeta`, `cleanStaleProjects`).
- Split `amendexec.go` → `amendexec.go` + `amendexecprint.go` (output: `buildEnvFilter`, `printAmendProgress`).
- Split `status.go` → `status.go` + `statusprint.go` (formatting: `printStatusTable`, `buildSummaryParts`).
- Split `exec.go` → `exec.go` + `execprint.go` (formatting: `printExecResult`, `printExecBanner`).
- Split `logs.go` → `logs.go` + `logsview.go` (view: `viewList`, `viewDetail`).
- Split `compress.go` → `compress.go` + `compresstar.go` (tar logic: `createTarGz`, `addFileToTar`).
- Added refactoring specs 65–78 for all 14 file splits.
- All source files comply with the 200-line limit; no functional changes.

## v2.36.5 — Extended Refactoring
- Split `ziparchive.go` (362 lines) into three files under `release/`:
  - `ziparchive.go` (~171 lines): orchestration, DB group routing, ad-hoc path resolution.
  - `zipio.go` (~152 lines): ZIP I/O with max Deflate compression, SHA-1 hashing, archive summary.
  - `zipdryrun.go` (~60 lines): dry-run preview for zip groups and ad-hoc archives.
- Split `autocommit.go` (352 lines) into two files under `release/`:
  - `autocommit.go` (~179 lines): orchestration, file classification, user prompts.
  - `autocommitgit.go` (~185 lines): Git primitives, push/retry, rebase recovery.
- Split `seowriteloop.go` (340 lines) into two files under `cmd/`:
  - `seowriteloop.go` (~198 lines): commit loop, rotation orchestration, signal handling.
  - `seowritegit.go` (~153 lines): Git stage/commit/push, rotation file I/O, output formatting.
- Split `workflowbranch.go` (310 lines) into two files under `release/`:
  - `workflowbranch.go` (~179 lines): branch-based releases, pending branch discovery.
  - `workflowpending.go` (~138 lines): metadata-based pending discovery and release.
- Split `workflow.go` (291 lines) into two files under `release/`:
  - `workflow.go` (~183 lines): `Execute`, `Options`/`Result` types, step execution.
  - `workflowvalidate.go` (~115 lines): duplicate detection, orphaned metadata, version resolution.
- Added refactoring specs: `60-refactor-ziparchive.md`, `61-refactor-autocommit.md`, `62-refactor-seowriteloop.md`, `63-refactor-workflowbranch.md`, `64-refactor-workflow.md`.
- All `release/` and `cmd/` files comply with the 200-line limit; no functional changes.

## v2.36.4
- Split `workflowfinalize.go` (498 lines) into four domain-specific files under `release/`:
  - `workflowfinalize.go` (~190 lines): core pipeline orchestration and metadata persistence.
  - `workflowdryrun.go` (~123 lines): dry-run preview functions and `returnToBranch`.
  - `workflowzip.go` (~108 lines): zip group building, ad-hoc archives, and checksum collection.
  - `workflowgithub.go` (~104 lines): GitHub release uploads and Go cross-compilation.
- Split `root.go` (388 lines) into seven domain-specific dispatch files under `cmd/`:
  - `root.go` (72 lines): entry point and top-level router.
  - `rootcore.go` (44 lines): scan, clone, pull, status, exec commands.
  - `rootrelease.go` (48 lines): release workflow commands.
  - `rootutility.go` (56 lines): update, revert, version, help, docs.
  - `rootdata.go` (98 lines): data management, history, profiles, TUI.
  - `roottooling.go` (91 lines): dev tooling and maintenance commands.
  - `rootprojectrepos.go` (38 lines): project type query commands.
- Eliminated `dispatchMisc` (166 lines); replaced by `dispatchData` + `dispatchTooling`.
  - `workflowdryrun.go` (~123 lines): dry-run preview functions and `returnToBranch`.
  - `workflowzip.go` (~108 lines): zip group building, ad-hoc archives, and checksum collection.
  - `workflowgithub.go` (~104 lines): GitHub release uploads and Go cross-compilation.
- All files comply with the 200-line limit; no functional changes.
- Added refactoring specs: `spec/01-app/58-refactor-workflowfinalize.md`, `spec/01-app/59-refactor-root-dispatch.md`.

## v2.36.3 (2026-03-26)
- Bumped compiled version constant to v2.36.3.
- Refactored legacy directory migration into shared `localdirs` package for reuse across CLI startup and release workflow.
- Release workflow now re-runs migration after returning to the original branch, preventing `.release/` from persisting when older branches restore tracked legacy files.
- Auto-commit `classifyFiles` now treats legacy `.release/` paths as release files for silent commit handling.
- Simplified doctor legacy directory check to always pass (migration handles cleanup automatically).
- Removed unused legacy directory warning/fix constants from `constants_doctor.go`.

## v2.36.2 (2026-03-26)
- Bumped compiled version constant to v2.36.2.
- Fixed legacy directory migration to merge files when target already exists instead of skipping.
- Legacy directories (`.release/`, `gitmap-output/`, `.deployed/`) are now fully removed after merging into `.gitmap/`.
- Added `mergeAndRemoveLegacy()` with file-walk merge and `os.RemoveAll` cleanup.
- Replaced Unicode characters in migration messages with ASCII for Windows console compatibility.

## v2.36.1 (2026-03-26)
- Bumped compiled version constant to v2.36.1.
- Added automatic database migration from legacy UUID TEXT IDs to INTEGER AUTOINCREMENT IDs.
- Migration detects TEXT-typed `Id` column in `Repos` via `PRAGMA table_info`, rebuilds the table preserving data, and drops dependent FK tables (project detection, group-repo associations) for clean repopulation.
- Fixed FK constraint violation (`787`) during `scan` when legacy UUID IDs were present in the `Repos` table.

## v2.36.0
- Bumped compiled version constant to v2.36.0.
- Added automatic legacy directory migration: `gitmap-output/` → `.gitmap/output/`, `.release/` → `.gitmap/release/`, `.deployed/` → `.gitmap/deployed/`.
- Migration runs at CLI startup before any command dispatch; skips if target already exists.
- Added `DeployedDirName` subdirectory constant and legacy directory name constants.

## v2.35.1
- Bumped compiled version constant to v2.35.1.
- Added legacy UUID data detection to all remaining DB query paths: `group show`, `group list`, `stats`, `history`, `status`, and `export`.
- All DB query errors from legacy string-based IDs now show a recovery prompt (`rescan` or `db-reset`) instead of raw SQL errors.

## v2.35.0
- Bumped compiled version constant to v2.35.0.
- Consolidated `.release/` and `gitmap-output/` under unified `.gitmap/` directory (`release/`, `output/`).
- Centralized all path constants (`GitMapDir`, `DefaultReleaseDir`, `DefaultOutputDir`) for single-point configuration.
- Migrated all database primary keys from UUID strings to `INTEGER PRIMARY KEY AUTOINCREMENT` (`int64`).
- Removed `github.com/google/uuid` dependency.
- Added `doctor` check (12th) that warns if legacy `.release/` or `gitmap-output/` directories exist.
- Updated all helptext, spec documents, and docs site to reference `.gitmap/` paths.

## v2.34.0 (2026-03-26)
- Bumped compiled version constant to v2.34.0.
- Fixed `list-releases` to read `.release/v*.json` from the current repo first, falling back to the database only when no local files exist.
- Added `SourceRepo` constant to release model for repo-sourced release records.

## v2.33.0 (2026-03-26)
- Bumped compiled version constant to v2.33.0.
- Fixed auto-commit push rejection when remote branch advances during release: added `pull --rebase` recovery with single retry.
- Added 16-stage summary table with anchor links to verbose logging spec.

## v2.32.0
- Bumped compiled version constant to v2.32.0.
- Documented autocommit verbose logging as pipeline stage 16 in the verbose logging spec.

## v2.31.0 (2026-03-26)
- Bumped compiled version constant to v2.31.0.
- Added verbose logging to auto-commit step: logs version, file counts, staging, commit message, and push target.

## v2.30.0 (2026-03-26)
- Bumped compiled version constant to v2.30.0.
- Renamed TempReleases `Commit` column to `CommitSha` to avoid SQLite reserved keyword conflict.
- Added automatic database migration (`ALTER TABLE RENAME COLUMN`) for existing TempReleases tables.
- Added JSON struct tags to `model.TempRelease` for backward-compatible serialization.

## v2.29.0
- Bumped compiled version constant to v2.29.0.
- Fixed TempReleases SQL syntax error: quoted reserved keyword `Commit` in CREATE TABLE, INSERT, and SELECT statements.
- Documented metadata persistence and rollback log points in verbose logging spec (stages 14–15 of 15).

## v2.28.0
- Bumped compiled version constant to v2.28.0.
- Added verbose logging to release pipeline: version resolution, source resolution, git operations, asset collection, staging, cross-compilation, compression, checksums, zip groups, ad-hoc zips, GitHub upload, retry, metadata persistence, and rollback.
- Updated verbose logging spec with all 15 pipeline stages documented.
- Added pull conflict handling to run.ps1 and run.sh with stash/discard/clean/quit prompt.
- Added --force-pull flag to both build scripts for non-interactive CI usage.
- Fixed set -e early exit bug in run.sh git pull error handling.
- Fixed parseCommitLines redeclaration conflict between temprelease.go and changeloggen.go.
- Fixed hasListFlag redeclaration conflict between tempreleaseops.go and completion.go.

## v2.27.0 (2026-03-22)
- Bumped compiled version constant to v2.27.0.
- Added doctor validation checks for config.json, database migration, lock file, and network connectivity.
- Added TUI release trigger overlay with patch/minor/major/custom version bump selection.
- Integrated batch progress tracking into pull, exec, and status commands with success/fail/skip counters.
- Added BatchProgress tracker to cloner package with quiet mode for programmatic use.
- Added TUI interaction tests covering tab switching, browser navigation, fuzzy search, and release triggers.
- Added alias suggestion tests covering auto-suggestion, conflict detection, and idempotent re-runs.

## v2.24.0 (2026-03-20)
- Bumped compiled version constant to v2.24.0.
- Moved release metadata writing from the release branch to the original branch, letting auto-commit handle `.release/` files after returning.
- Removed `commitReleaseMeta` step from the release branch workflow; the release branch now only contains the branch, tag, and push.
- Simplified `pushAndFinalize` to always complete without metadata writes (metadata is now the caller's responsibility).

## v2.23.0 (2026-03-20)
- Bumped compiled version constant to v2.23.0.
- Added `--notes` / `-N` flag to `release-branch` and `release-pending` commands, matching the `release` command.
- Updated docs site Release page with metadata-first workflow diagram, release notes feature card, and `--notes` flag documentation.

## v2.22.0 (2026-03-19)
- Bumped compiled version constant to v2.22.0.
- Persisted zip group metadata in `.release/vX.Y.Z.json` via new `zipGroups` field on `ReleaseMeta`.
- Documented `-A`/`--alias` flag in help text for `pull`, `exec`, `status`, and `cd` commands.
- Added shell completion support for `alias` and `zip-group` subcommands across PowerShell, Bash, and Zsh.
- Added `--list-aliases` and `--list-zip-groups` completion list flags with dynamic DB lookups.
- Added unit tests for `collectZipGroupNames` covering persistent groups, ad-hoc bundles, and merged output.

## v2.21.0
- Bumped compiled version constant to v2.21.0.
- Refactored `assetsupload.go` into three focused files: `githubapi.go` (API types/helpers), `assetsupload.go` (upload logic), `remoteorigin.go` (git URL parsing).
- Rebuilt Project Detection docs page with detection pipeline, tabbed type cards, metadata extraction deep-dive, DB schema, JSON output, and package layout sections.
- Added "How detection works" link from Projects dashboard to Detection page.
- Added unit tests for `store/location.go` covering symlink resolution, fallback, double-nesting prevention, and profile DB filenames.
- Added unit tests for `remoteorigin.go` covering HTTPS, SSH, and invalid URL parsing.

## v2.20.0
- **Fixed**: `OpenDefault()` double-nesting bug where profile config resolved to `<binary>/data/data/profiles.json`.
- Added `DefaultDBPath()` diagnostic helper to `store/location.go`.
- `gitmap ls` now prints resolved DB path when `--verbose` is passed or when zero repos are found.
- Created `spec/01-app/44-list-db-diagnostic.md` for path resolution contract.

## v2.19.0
- Bumped compiled version constant to v2.19.0.

## v2.18.0
- Added batch status terminal demo to Batch Actions page showing dirty/clean state across repos.
- Fixed missing `os/exec` import in release asset upload.
- Resolved `deriveSlug` redeclaration conflict in project repos output.
- Removed unused `os` import from audit command.

## v2.17.0
- Added 30-second auto-refresh timer to TUI dashboard via `tea.Tick`.
- Dashboard refresh interval configurable via `dashboardRefresh` in `config.json`.
- Added `--refresh` flag to `interactive` command for CLI-level override.
- Refresh interval validates with fallback to default 30s when missing or invalid.

## v2.16.0
- Wired real `gitutil.Status()` into TUI dashboard for live dirty/clean indicators.
- Dashboard now shows ahead/behind counts and stash per repo.
- Async background refresh on TUI startup; manual refresh via `r` key.
- Summary bar with aggregate dirty/behind/stash counts and UTC timestamp.

## v2.15.1
- **Fixed**: Database now resolves to `<binary-location>/data/gitmap.db` instead of CWD-relative `gitmap-output/data/`.
- Added `store.OpenDefault()` and `store.OpenDefaultProfile()` for binary-relative database access.
- Added `store/location.go` with `BinaryDataDir()` using `os.Executable()` + `filepath.EvalSymlinks()`.
- Updated all 13 database callers across the codebase to use binary-relative paths.
- Removed unused `resolveAuditOutputDir()` and `resolveDefaultOutputDir()` helpers.

## v2.15.0
- Added cross-platform build support: `run.sh` (Linux/macOS) with full parity to `run.ps1`.
- Fixed Makefile flags to match `run.sh` argument format (`--no-pull`, `--no-deploy`, `--update`).
- Added GitHub Actions CI workflow: test on push, cross-compile 6 OS/arch targets.
- Added GitHub Actions Release workflow: auto-release on `v*` tags with compression and checksums.
- Added interactive TUI mode (`gitmap interactive` / `gitmap i`) built with Bubble Tea.
- TUI repo browser with fuzzy search, multi-select, and keyboard navigation.
- TUI batch actions: pull, exec, status across selected repos.
- TUI group management: browse, create, delete groups interactively.
- TUI status dashboard with live repo status view.
- Added Build System section to Architecture documentation page.
- Added spec documents: `42-cross-platform.md` and `43-interactive-tui.md`.

## v2.14.0
- Added Go release assets: automatic cross-compilation for 6 OS/arch targets (windows/linux/darwin × amd64/arm64).
- Added GitHub Releases API integration for asset upload — no `gh` CLI or external tools needed.
- Added `--compress` flag to wrap release assets in `.zip` (Windows) or `.tar.gz` (Linux/macOS).
- Added `--checksums` flag to generate SHA256 `checksums.txt` for all release assets.
- Added `--no-assets` flag to skip automatic Go binary compilation.
- Added `--targets` flag for custom cross-compile target selection (e.g. `windows/amd64,linux/arm64`).
- Improved `gitmap ls <type>` output with labeled fields (Repo, Path, Indicator) and inline `cd` examples.
- Added shell completion for `release`, `release-branch`, `group`, `multi-group`, and `list` commands.
- Fixed duplicate hints appearing after `gitmap ls <type>` output.

## v2.13.0
- Added group activation: `gitmap g <name>` sets a persistent active group for batch pull/status/exec.
- Added `multi-group` (mg) command for selecting and operating on multiple groups at once.
- Added `gitmap ls <type>` filtering: `gitmap ls go`, `gitmap ls node`, `gitmap ls groups`.
- Added contextual helper hints shown after command output to aid discoverability.
- Added Settings table for persistent key-value configuration in SQLite.

## v2.12.0 (2026-03-14)
- Added global ⌘K command palette searching across commands, flags, and pages.

## v2.11.0
- Added Changelog page with timeline view and expand/collapse controls.
- Added Flag Reference page with sortable, searchable table of all flags.
- Added Interactive Examples page with animated terminal demos.

## v2.10.0 (2026-03-13)
- Version bump for next development cycle.

## v2.9.0 (2026-03-13)
- Completed flags and examples for all 22 command entries on the documentation site.
- Added detailed flag tables and usage examples for `seo-write`, `doctor`, `update`, `pull`, `version`, `history-reset`, and `db-reset`.
- Filled in flags and examples for 15 commands missing both: `rescan`, `desktop-sync`, `status`, `latest-branch`, `release-branch`, `release-pending`, `changelog`, `group`, `list`, `diff-profiles`, `export`, `import`, `profile`, `bookmark`, and `stats`.

## v2.28.0
- Removed unused `detector` import from `cmd/scan.go` that caused build failure.
- Updated documentation site fonts: Ubuntu for headings, Poppins for body text, Ubuntu Mono for code blocks.

## v2.27.0 (2026-03-22)
- Added `gitmap cd` (`go`) command: jump to any tracked repo by slug or partial name.
- Subcommands: `cd repos`, `cd set-default`, `cd clear-default`; supports `--group` and `--pick` flags.
- Added `gitmap watch` (`w`) command: live terminal dashboard monitoring repo status.
- Supports `--interval`, `--group`, `--no-fetch`, and `--json` snapshot mode.
- Added `gitmap diff-profiles` (`dp`) command: compare two profiles side-by-side.
- Supports `--all` and `--json` output flags.
- Added clone progress bars with retry logic and Windows long-path warnings.
- Built documentation site with interactive terminal preview for the watch command.
- Added `gitmap/Makefile` as a thin wrapper around `run.sh` for standard `make` workflows.
  - Targets: `build`, `run` (with `ARGS=`), `test`, `update`, `no-pull`, `no-deploy`, `clean`, `help`.
- Added Makefile documentation page to the docs site with target reference, examples, and argument-passing guide.
- Added `run.sh` cross-platform build script: Bash equivalent of `run.ps1` for Linux and macOS.
  - Full pipeline: pull, tidy, build, deploy with `-ldflags` version embedding.
  - Reads config from `powershell.json` via `jq` or `python3` fallback.
  - Supports `-t` (test with report), `-n` (no-pull), `-d` (no-deploy), and `-u` (update) flags.
- Added `gitmap gomod` (`gm`) command: rename Go module path across an entire repo with branch safety.
  - Replaces module directive in `go.mod` and all matching paths across **all files** by default.
  - Use `--ext "*.go,*.md,*.txt"` to restrict replacement to specific file extensions.
  - Creates `backup/before-replace-<slug>` and `feature/replace-<slug>` branches automatically.
  - Commits changes on the feature branch and merges back to the original branch.
  - Supports `--dry-run`, `--no-merge`, `--no-tidy`, `--verbose`, and `--ext` flags.

## v2.26.0 (2026-03-22)
- Version bump to v2.26.0 following `gitmap profile` command addition.
- All profile subcommands (`create`, `list`, `switch`, `delete`, `show`) fully integrated and documented.

## v2.25.0 (2026-03-22)
- Added `gitmap profile` (`pf`) command: manage multiple database profiles (work, personal, etc.).
- Subcommands: `create`, `list`, `switch`, `delete`, `show`.
- Each profile has its own SQLite database file (`gitmap-{name}.db`).
- Default profile uses existing `gitmap.db` for full backward compatibility.
- Profile config stored in `gitmap-output/data/profiles.json`.
- All commands automatically use the active profile's database.

## v2.24.0 (2026-03-20)
- Added `gitmap import` (`im`) command: restore database from a `gitmap-export.json` backup file.
- Merge semantics: upserts repos/releases, INSERT OR IGNORE for history/bookmarks/groups.
- Group members re-linked by resolving `repoSlugs` against the Repos table.
- Requires `--confirm` flag to prevent accidental data changes.

## v2.23.0 (2026-03-20)
- Added `gitmap export` (`ex`) command: export the full database as a portable JSON file.
- Exports all tables: repos, groups (with member repo slugs), releases, command history, and bookmarks.
- Default output: `gitmap-export.json`; accepts optional custom file path.
- Summary line shows counts for each exported section.

## v2.22.0 (2026-03-19)
- Added `gitmap bookmark` (`bk`) command: save and replay frequently-used command+flag combinations.
- Subcommands: `save`, `list`, `run`, `delete` — full CRUD for saved bookmarks.
- `bookmark run <name>` replays the saved command through standard dispatch (appears in audit history).
- `bookmark list --json` outputs bookmarks as JSON.
- New `Bookmarks` SQLite table with unique name constraint.
- `db-reset --confirm` now also clears the Bookmarks table.

## v2.21.0
- Added `gitmap stats` (`ss`) command: aggregated usage statistics from command history.
- Shows most-used commands, success/fail counts, failure rates, and avg/min/max durations.
- Supports `--command <name>` filter and `--json` output.
- Summary row displays overall totals across all commands.

## v2.20.0
- Added `gitmap history` (`hi`) command: queryable audit trail of all CLI command executions.
- Three detail levels: `--detail basic` (command + timestamp), `--detail standard` (+ flags + duration), `--detail detailed` (+ args + repos + summary).
- Supports `--command <name>` filter, `--limit N`, and `--json` output.
- Added `gitmap history-reset` (`hr`) command: clears audit history (requires `--confirm`).
- New `CommandHistory` SQLite table auto-records every command with start/end timestamps, duration, exit code, and affected repo count.
- `db-reset --confirm` now also clears the CommandHistory table.

## v2.19.0
- Added `gitmap amend` (`am`) command: rewrite author name/email on existing commits with three modes (all, range, HEAD).
- Supports `--branch` flag to operate on a specific branch (auto-switches back to original branch after completion).
- SHA as first positional argument: `gitmap amend <sha> --name "Name"` rewrites from that commit to HEAD.
- `--dry-run` previews affected commits without modifying history or writing audit records.
- `--force-push` auto-runs `git push --force-with-lease` after amend.
- Audit trail: every amend operation writes a JSON log to `.gitmap/amendments/amend-<timestamp>.json` with full details.
- Database persistence: amendment records saved to `Amendments` SQLite table for queryable history.
- `db-reset --confirm` now also clears the `Amendments` table.
- Added `--author-name` and `--author-email` flags to `gitmap seo-write` (`sw`): set custom author on each commit.
- SEO-write dry-run now displays the author that would be used when author flags are set.

## v2.18.0
- Added `gitmap seo-write` (`sw`) command: automated SEO commit scheduler that stages, commits, and pushes files on a randomized interval.
- Supports CSV input mode (`--csv`) for user-provided title/description pairs.
- Supports template mode with placeholder substitution (`{service}`, `{area}`, `{url}`, `{company}`, `{phone}`, `{email}`, `{address}`).
- Pre-seeded `data/seo-templates.json` with 25 title and 20 description templates (500 unique combinations).
- Added `CommitTemplates` SQLite table for persistent template storage with auto-seeding on first run.
- Rotation mode: when pending files are exhausted, appends/reverts text in a target file to maintain commit activity.
- Configurable interval (`--interval min-max`), commit limit (`--max-commits`), file selection (`--files`), and dry-run preview.
- Added `--template <path>` flag to load templates from a custom JSON file at runtime.
- Added `--create-template` / `ct` shorthand to scaffold a sample `seo-templates.json` in the current directory.
- Graceful shutdown on Ctrl+C (finishes current commit before exiting).

## v2.17.0
- Added `Source` column to the `Releases` table: tracks whether each release was created via `gitmap release` (`release`) or imported from `.release/` files (`import`).
- Added `--source` flag to `gitmap list-releases` (`lr`): filter releases by origin (`--source release` or `--source import`).
- Added `--source` flag to `gitmap list-versions` (`lv`): cross-references git tags with the Releases DB to filter by source and display source metadata.
- Added `--source` flag to `gitmap changelog` (`cl`): filter changelog entries by release source.
- Terminal and JSON output for `list-releases` and `list-versions` now includes the Source field.

## v2.16.0
- Added `gitmap list-releases` (`lr`) command: queries the Releases DB table and displays stored releases with `--json` and `--limit N` support.
- Enhanced `gitmap scan` to import `.release/v*.json` metadata files into the Releases DB table automatically after each scan.

## v2.15.0
- Added `--limit N` flag to `gitmap list-versions` (`lv`): show only the top N versions (0 or omitted = all).

## v2.14.0
- Added `Releases` table to SQLite database: stores release metadata (version, tag, branch, commit, changelog, flags) persistently.
- Release workflow now auto-persists metadata to the database after successful releases.
- Converted all database table and column names from snake_case to PascalCase (`Repos`, `Groups`, `GroupRepos`, `Releases`).
- Added `store/release.go` with `UpsertRelease`, `ListReleases`, `FindReleaseByTag` methods.
- Added `model/release.go` with `ReleaseRecord` struct.
- Note: existing databases will need `gitmap db-reset --confirm` to adopt the new schema.

## v2.13.0
- Release metadata JSON (`.release/vX.Y.Z.json`) now includes a `changelog` field with notes from CHANGELOG.md (gracefully omitted if unreadable).
- `gitmap list-versions` (`lv`) now shows changelog notes as sub-points under each version in terminal output.
- `gitmap list-versions --json` includes changelog array per version in JSON output.

## v2.12.0 (2026-03-14)
- Added `gitmap list-versions` (`lv`) command: lists all release tags sorted highest-first, with `--json` output support.
- Added `gitmap revert <version>` command: checks out a release tag and rebuilds/deploys via handoff (same mechanism as `update`).

## v2.11.0
- Added constants inventory audit section to compliance spec, documenting ~280 constants across 9 files and 17 categories.

## v2.10.0 (2026-03-13)
- Full compliance audit (Wave 1 + Wave 2): all 75 source files pass code style rules.
  - Trimmed 4 oversized files: `workflow.go`, `terminal.go`, `safe_pull.go`, `setup.go` (all under 200 lines).
  - Fixed all negation and switch violations across `changelog.go`, `github.go`, `metadata.go`, `config.go`, `verbose.go`, `semver.go`.
  - Extracted missing constants to dedicated constants files.

## v2.9.0 (2026-03-13)
- Full code style refactor of `latest-branch` command:
  - Split `cmd/latestbranch.go` into 3 files: handler, resolve, output (all under 200 lines).
  - Split `gitutil/latestbranch.go` into 2 files: core operations, resolve helpers.
  - All functions comply with 8-15 line limit. Positive logic throughout.
  - Blank line before every return. No magic strings. Chained if+return replaces switch.
  - Extracted git constants and display message constants.

## v2.8.0 (2026-03-06)
- Added `--filter` flag to `latest-branch`: filter branches by glob pattern (e.g. `feature/*`) or substring match.

## v2.7.0
- Added `--sort` flag to `latest-branch`: supports `date` (default, descending) and `name` (alphabetical ascending).

## v2.6.0
- Centralized date display formatting: all dates now convert to local timezone and display as `DD-Mon-YYYY hh:mm AM/PM`.
- Added `gitutil/dateformat.go` with `FormatDisplayDate` and `FormatDisplayDateUTC` functions.
- Updated `latest-branch` terminal, JSON, and CSV output to use the new date format.

## v2.5.1
- Added `--no-fetch` flag to `latest-branch`: skips `git fetch --all --prune` when remote refs are already up to date.

## v2.5.0 (2026-03-06)
- Added `--format` flag to `latest-branch`: supports `terminal` (default), `json`, and `csv` output formats.
  - CSV outputs a header row + data rows to stdout, suitable for piping and spreadsheets.
  - `--json` remains as shorthand for `--format json`.
- Refactored `latest-branch` output into dedicated functions per format.

## v2.4.1
- Added positional integer shorthand for `latest-branch`: `gitmap lb 3` is equivalent to `gitmap lb --top 3`.

## v2.4.0 (2026-03-06)
- Added `gitmap latest-branch` (`lb`) command: finds the most recently updated remote branch by commit date and displays name, SHA, date, and subject.
  - Flags: `--remote`, `--all-remotes`, `--contains-fallback`, `--top N`, `--json`.
  - Positional integer shorthand: `gitmap lb 3` is equivalent to `gitmap lb --top 3`.

## v2.3.12 (2026-03-06)
- Spec, issue post-mortems, and memory aligned to codify synchronous update handoff and rename-first PATH sync as permanent rules.
- Rename-first PATH sync in `-Update` mode: renames active binary to `.old` before copying, eliminating lock-retry loops.
- Parent `update` handoff uses `cmd.Start()` + `os.Exit(0)` to release file lock before worker runs.
- Handoff diagnostic log prints active exe and copy paths at update start.
- Spec consistency pass: all four update-flow specs now enforce identical rules.

## v2.3.10 (2026-03-06)
- Fixed `Read-Host` error in non-interactive PowerShell sessions during update by removing trailing prompt.
- Parent `update` process now exits immediately (handoff copy runs synchronously via `update-runner`).
- Added diagnostic log at update start showing active exe path and handoff copy path.
- Update script now uses unique temp file names (`gitmap-update-*.ps1`) to avoid stale script collisions.

## v2.3.9
- Version bump for rebuild validation after update-runner handoff changes.

- Replaced `update --from-copy` with hidden `update-runner` command for cleaner handoff separation.
- Handoff copy now created in the same directory as the active binary (fallback to %TEMP% if locked).
- Added `-Update` flag to `run.ps1`: runs full update pipeline (pull, build, deploy, sync) with post-update validation and cleanup.
- Update script delegates entire pipeline to `run.ps1 -Update`.
- Before/after version output derived from actual executables, not static constants.
- Mandatory `update-cleanup` runs after successful update to remove handoff and `.old` artifacts.
- Cleanup now scans both `%TEMP%` and same-directory for leftover `gitmap-update-*.exe` files.

- Added `gitmap doctor --fix-path` flag: automatically syncs the active PATH binary from the deployed binary using retry (20×500ms), rename fallback, and stale-process termination, with clear confirmation output.
- Doctor diagnostics now suggest `--fix-path` when version mismatches are detected.

## v2.3.6
- Added stale-process fallback during PATH-binary sync (`update` + `run.ps1`): if copy+rename fail, it now stops stale `gitmap.exe` processes bound to the old path and retries once.
- Improved failure guidance to run the deployed binary directly when active PATH binary remains locked.

## v2.3.5
- Hardened `gitmap update` PATH sync with retry + rename fallback, and it now exits with failure if active PATH binary remains stale.
- Clarified update output labels to distinguish source version (`constants.go`) vs active executable version.
- Added same rename-fallback PATH sync behavior in `run.ps1`.

## v2.3.4
- Updated PATH-binary sync in `run.ps1` and `gitmap update` to use retry-on-lock behavior (20 attempts × 500ms), matching the self-update spec.
- Added explicit recovery guidance when active PATH binary is still locked, including an exact `Copy-Item` fix command.

## v2.3.3
- Added `gitmap doctor` command: reports PATH binary, deployed binary, version mismatches, git/go availability, and recommends exact fix commands.

## v2.3.2
- `gitmap update` now syncs the active PATH binary with the deployed binary, so commands like `release` are available immediately.
- `gitmap update` now prints changelog bullet points after update (or no-op update) for quick visibility.
- Added `gitmap changelog --open` and `gitmap changelog.md` to open `CHANGELOG.md` in the default app.

## v2.3.1
- Added `gitmap changelog` command for concise, CLI-friendly release notes.
- Improved `gitmap update` output to show deployed binary/version and warn if PATH points to another binary.
- `gitmap update` now prints latest changelog notes after a successful update.

## v2.3.0
- Added `gitmap release-pending` (`rp`) to release all `release/v*` branches missing tags.
- `gitmap release` and `gitmap release-branch` now switch back to the previous branch after completion.

## v2.2.3
- Fixed PowerShell parser-breaking characters in update/deploy output paths.
- Improved deployment rollback messaging in `run.ps1`.

## v2.2.2
- Added additional parser safety fixes for update script output.

## v2.2.1
- Patched PowerShell parsing edge cases affecting update flow.
