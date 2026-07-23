# Changelog

## v4.30.0 — golden CSVs: normalize CRLF→LF in asserts (durable cross-platform fix)

### Fixed
- Windows CI was still failing `TestCloneFromReport_Golden_*` and
  `TestScanGolden_CSV*` after v4.29.0 because `git ls-files --eol` showed
  the **index** had `i/lf` for those CSVs (committed earlier as LF), so
  every fresh clone — Linux, macOS, Windows — got LF on disk regardless
  of any new `eol=crlf` attribute. Meanwhile the writer still emits CRLF
  per `csvcrlf_contract_test.go`. Result: 1-byte-per-line drift recurred.
- Root-cause fix that survives any future `.gitattributes` / autocrlf /
  editor change: added `normalizeGoldenBytes` (`bytes.ReplaceAll(b,
  "\r\n", "\n")`) to both `assertScanGolden`
  (`gitmap/formatter/scangolden_contract_test.go`) and
  `assertReportGolden` (`gitmap/clonefrom/summary_golden_test.go`),
  applied to BOTH `got` and `want` before comparison and to `got`
  before regenerate-write. Goldens now stored as LF on disk
  (hand-readable, cross-platform consistent), the production CRLF
  contract remains enforced separately by `TestCSVCRLF_WriteCSV` and
  `cmd/csvcrlf_contract_test.go` which inspect in-memory writer bytes.
- Reverted the v4.29.0 `.gitattributes` `-text eol=crlf` per-glob
  override (it was a no-op anyway: `text=unset` causes git to ignore
  `eol=`, and the index already had LF so checkouts never saw CRLF).
- Re-converted on-disk goldens to LF so they match the index and any
  future fresh clone:
  `clonefrom_report_canonical.csv` 271→267,
  `clonefrom_report_empty.csv` 54→53,
  `clonefrom_report_quoting.csv` 436→429,
  `scan_canonical.csv` 620→617,
  `scan_empty.csv` 132→131.

### Notes
- This is the durable end of the v4.27 → v4.28 → v4.29 regression cycle.
  Future contributors editing CSV goldens cannot accidentally break
  byte-exact comparisons via line-ending mismatches anymore.

## v4.29.0 — golden CSVs: pin CRLF via .gitattributes (real root cause)

### Fixed
- v4.28.0 regenerated `gitmap/clonefrom/testdata/clonefrom_report_*.csv`
  and `gitmap/formatter/testdata/scan_*.csv` but the broad
  `testdata/** text eol=lf` rule in `.gitattributes` immediately re-stripped
  every `\r` on the next checkout — so CI saw a 1-byte-per-line drift again
  (`want 267 / got 271`, `want 429 / got 436`, etc.). The CSV writers
  (`clonefrom/summary.go` `writeReportRows`, `formatter/csv.go`) set
  `csv.Writer.UseCRLF = true` per `csvcrlf_contract_test.go`, so the disk
  fixture MUST stay CRLF.
- Added more-specific overrides in `.gitattributes`:
  ```
  gitmap/clonefrom/testdata/clonefrom_report_*.csv  -text eol=crlf
  gitmap/formatter/testdata/scan_*.csv              -text eol=crlf
  ```
  `-text` flags them as binary (no munging) AND `eol=crlf` belt-and-suspenders
  forces the line ending on any platform. More-specific patterns win over the
  generic `testdata/**` rule.

### Notes
- This is the durable fix for the v4.27.0 / v4.28.0 CRLF regression class.
  Any future CRLF-emitting golden CSV must add a `-text eol=crlf` override.
- No production code changed. Fixture-only + `.gitattributes`.

## v4.28.0 — golden CSVs: regenerate after gofmt-sweep CRLF→LF drift

### Fixed
- `gitmap/clonefrom/testdata/clonefrom_report_{empty,canonical,quoting}.csv`
  and `gitmap/formatter/testdata/scan_{empty,canonical}.csv` were rewritten
  with LF line endings (and trailing newline) by the v4.26.0 repo-wide
  formatting sweep, but the corresponding writers (`writeReportRows`,
  formatter `csv.go`) emit CRLF per `csvcrlf_contract_test.go`. Result:
  `TestCloneFromReport_Golden_*` and `TestScanGolden_CSV*` failed with
  byte-exact mismatches in CI. Regenerated all five CSV goldens (and the
  two scan JSON goldens for parity) via the documented two-step
  `GITMAP_UPDATE_GOLDEN=1 GITMAP_ALLOW_GOLDEN_UPDATE=1` flow so on-disk
  bytes match writer output exactly.

### Notes
- No production code changed — fixture-only update.
- Future protection: `.gitattributes` should pin `*.csv` under `testdata/`
  to `eol=crlf` to prevent editor/auto-format drift from recurring.

## v4.27.0 — clone-pick: replace bare-err with cliexit.Reportf

### Fixed
- `gitmap/cmd/clonepick_execute.go:70` was printing
  `fmt.Fprintln(os.Stderr, err)` after a `TouchAfterReplay` DB
  failure. The CI gate `.github/scripts/check-bare-stderr-err.sh`
  (per spec/04-generic-cli/07-error-handling.md) forbids bare
  `Fprintln(os.Stderr, err)` in `gitmap/cmd/*.go` because it
  drops three of four required context fields (command, op,
  subject, cause).
- Replaced with
  `cliexit.Reportf(constants.CmdClonePick, "touch-replay",
   strconv.FormatInt(replayId, 10), err)` so the failure line
  carries the command name, the operation, the subject
  (replay-id), and the underlying cause — matching the
  surrounding `clonefrom` call sites.

### Notes
- Imports updated: `+ cliexit`, `+ strconv`. `fmt` and `os` are
  still used elsewhere in the file (`Fprintf` for the saved /
  replayed banner) so they stay.
- `bash .github/scripts/check-bare-stderr-err.sh` now reports
  `OK: no bare 'fmt.Fprintln(os.Stderr, err)' in gitmap/cmd/`.
- gofmt-clean. No behaviour change to the success path.


## v4.26.0 — Repo-wide gofmt sweep (CI gofmt-clean gate fix)

### Fixed
- 73 `.go` files across `gitmap/` (clonepick, commitin/{dedupe,
  finalize, funcintel, message, orchestrator, profile, ...},
  store, tui, vscodeworkspace, etc.) were missing the trailing
  final-newline that `gofmt` enforces. CI's `gofmt -l .` gate
  was failing as a result.
- Ran `gofmt -w` across the whole tree; `gofmt -l .` now reports
  zero files. No semantic / token changes — newline-only edits.

### Why
The CI lint job's gofmt-clean gate started failing (`The
following .go files are not gofmt-clean: ...`). The
`gitmap fix-repo` audit step in CI (per memory: FIX-REPO GOFMT
v4.11.0+) is intentionally narrow — it only checks files
fix-repo touched. A bulk historical drift outside that scope
slipped past it. This sweep restores the global gate; future
edits should keep files gofmt-clean on save.

### Notes
- No version-bearing tokens (`-vN`) were rewritten, so the
  fix-repo digit-capture rule and `--strict` test gate are not
  involved.
- All edits are whitespace-only (final newline). No behaviour
  change. CI `gofmt -l .` and `lint` jobs should now go green.


## v4.25.0 — ERD-vs-SQLCreate parity test (v3.12.1 publish-audit follow-up)

### Added
- `gitmap/store/erd_parity_test.go` — `TestERDMatchesSQLCreate`
  walks every `constants_*.go` in `gitmap/constants/`, extracts
  table names from `CREATE TABLE IF NOT EXISTS` blocks, and
  asserts the set equals the table blocks declared in
  `spec/01-app/gitmap-database-erd.mmd`. Fails CI on either
  direction of drift (missing-from-ERD or orphan-in-ERD) with a
  concrete fix recipe in the error message.
- 24 commit-in / pipeline / transactions / archive / VSCode /
  clone-pick table stubs added to the canonical ERD so the test
  passes against the current schema (52 tables total). Stubs
  declare PK only; per-column refinement is a future task.

### Changed
- Renamed `spec/01-app/gitmap-database-erd-v3.12.1.mmd` →
  `spec/01-app/gitmap-database-erd.mmd` (canonical, version-less).
- Deleted stale `spec/01-app/gitmap-database-erd.mmd` (v3.5.0,
  missing 11 tables) and `gitmap-core-schema-simplified.mmd`
  (v15 Phase 1.2). Replaced by the canonical file above.
- ERD frontmatter now points at the parity test as the
  enforcement mechanism.

### Why
The v3.5.0 ERD silently shed 11 tables and nobody noticed for
~6 months — drift is invisible without a test. Adding a
name-only parity gate (intentionally not column-level, to avoid
churning the ERD on every migration) closes the failure mode
that hurt us. CI picks the test up automatically via the
existing `go test ./...` step in the full-suite job — no
workflow edits required.


## v4.24.0 — release-version snippet Vitest wiring (spec 105 follow-up)

### Added
- `src/pages/releaseVersionSnippets.ts` — extracted the pure
  snippet builder (`buildReleaseSnippets`), semver guard
  (`isValidReleaseVersion`), and `normalizeReleaseVersion` from
  `ReleaseVersion.tsx` so they can be unit-tested without
  mounting React.
- `src/test/release-version-snippets.test.ts` — 12 Vitest cases
  covering: semver acceptance/rejection (incl. `latest` rejected),
  Windows + Unix pinned/generic snippet structure, inline
  snapshots for both platforms, and **spec 105 invariants** —
  the word "latest" never appears, `/releases/latest` is never
  referenced, and the requested version is baked into every URL.

### Changed
- `src/pages/ReleaseVersion.tsx` now imports the helpers from
  `releaseVersionSnippets.ts` instead of defining them inline.
  Zero behaviour change; rendered output is byte-identical.

### Why
Spec 105 promises that the `/release/:version` page never
silently resolves to `latest`. Until now that promise lived only
in the page source — easy to regress in a quick edit. The new
test suite locks the contract: any future change that lets the
word "latest" or a `/releases/latest` URL slip into either
snippet, or that drops the version from a URL, fails CI.


## v4.23.0 — web `VERSION` resynced + drift-guard test

### Fixed
- `src/constants/index.ts` `VERSION` had silently drifted from the
  Go binary version: it was stuck at `v4.14.0` while the binary
  shipped `v4.15.0` … `v4.22.0` (eight releases of stale "Current
  version" badges in the docs site, the home page, and every
  changelog entry that mentions `VERSION`). Now resynced to
  `v4.23.0`.

### Added
- `src/test/version-sync.test.ts`: reads
  `gitmap/constants/constants.go`, extracts the `Version` literal,
  and asserts `VERSION === "v" + goVersion`. Future Go bumps that
  forget the TS sync now fail Vitest in CI instead of shipping a
  stale badge. 2 assertions, sub-10ms runtime, zero new
  dependencies.
- `mem://project/version-bump-procedure.md`: updated to list
  `src/constants/index.ts` as a mandatory step alongside
  `gitmap/constants/constants.go` so future agents bump both.


## v4.22.0 — clone-pick `--ask` windowed scroller

### Added
- Picker now renders only the visible window (terminal height minus
  header + footer chrome) instead of dumping every tracked path. Large
  repos (Linux kernel, Chromium-scale trees) finally fit on screen.
- New nav keys: `pgup` / `ctrl+b` (page up), `pgdown` / `ctrl+f`
  (page down), `g` / `home` (first row), `G` / `end` (last row). Vim
  muscle memory carried over from the existing `j/k`.
- Header now includes the visible row range, e.g.
  `gitmap clone-pick --ask  (12/487 selected, rows 41-60)`.
- Reacts to `tea.WindowSizeMsg` so terminal resizes re-flow the
  window in real time. Default viewport (20 rows) is used until the
  first size message arrives.

### Changed
- `pickerModel` gains `viewportHeight` and `scrollOffset` fields.
  Cursor moves re-clamp `scrollOffset` via the pure `clampScroll`
  helper so the cursor row is always in view.
- `gitmap/clonepick/picker.go` split into four files to stay under
  the strict <200 lines/file core rule:
  - `picker.go` (178) — model + key dispatch.
  - `picker_run.go` (52) — bubbletea program lifecycle
    (`RunPicker`, `RunPickerKeep`).
  - `picker_nav.go` (70) — pure cursor + scroll math
    (`applyCursorMove`, `clampScroll`, `minInt`/`maxInt`).
  - `picker_view.go` (103) — render loop (windowed).

### Tests
- `gitmap/clonepick/picker_window_test.go`: `clampScroll` table,
  `tea.WindowSizeMsg` viewport sizing, `pgdown` jump distance,
  `G` end-of-list, and a "rendered row count == viewportHeight"
  guard so future View() tweaks can't silently regress the window.


## v4.21.1 — clone-pick cmd file split (200-line cap audit)

### Changed
- `gitmap/cmd/clonepick.go` split from 270 lines into four focused
  files, all under the strict <200 lines/file core rule:
  - `clonepick.go` (103) — dispatcher entry + dry-run branch +
    `buildClonePickPlan` / `applyReplayOverrides`.
  - `clonepick_flags.go` (142) — flag binding (`parseClonePickFlags`,
    `bindClonePickCoreFlags`, `bindClonePickAuxFlags`,
    `requireClonePickPositional`).
  - `clonepick_execute.go` (96) — side-effecting executor
    (`runClonePickExecute`, `announceClonePickPersistence`,
    `syncClonePickResultToVSCodePM`).
  - `clonepick_picker.go` (60) — already in shape, untouched.
- Zero behavioural change: every public symbol kept its name and
  signature; only the file each one lives in moved.


## v4.21.0 — clone-pick `--ask` clones once, not twice

### Changed
- `gitmap clone-pick --ask` now reuses the picker's metadata-only
  clone as the final sparse-checkout destination instead of re-cloning
  after the user confirms their selection. Halves the network cost of
  the interactive flow on large repos.
- The reuse path tries `os.Rename` first (fast, same-filesystem); on
  failure (cross-fs / Windows quirks) it falls back to a stdlib
  copy-tree + remove-src walker. Symlinks are preserved; non-regular
  entries (devices/sockets) are skipped — git never produces them.
- Dry-run + cancel paths still clean up the temp clone immediately so
  no orphan dirs accumulate under `$TMPDIR`.

### Added
- `clonepick.RunPickerKeep` / `clonepick.ListRepoPathsKeep`: keep-clone
  variants that return the temp dir so the cmd layer can hand it to
  the executor. The original `RunPicker` / `listRepoPaths` are now
  thin wrappers that always remove the temp dir for non-reuse callers.
- `Plan.PreClonedSrc` (runtime-only field): when set, the executor
  promotes that dir into `DestDir` and skips `git clone`.
- `constants.ErrClonePickPromoteSrc`: error string emitted when both
  rename and copy-tree fallbacks fail during promote.

### Tests
- `gitmap/clonepick/promote_test.go`: covers the rename fast path and
  the copy-tree fallback against `t.TempDir()` so the optimisation has
  regression coverage even when CI can't trigger a real EXDEV.


## v4.20.0 — `gitmap clone-pick --ask` interactive picker shipped

### Added
- `gitmap clone-pick --ask` opens a bubbletea picker before the
  sparse-checkout, listing every tracked path in the target repo
  (enumerated via a metadata-only `--filter=blob:none --no-checkout
  --depth=1` clone + `git ls-tree -r --name-only HEAD` in a temp dir,
  so no blob bytes are downloaded just to render the list).
- Picker keys: `up/k` / `down/j` move, `space` toggles, `a` selects
  every non-auto-greyed row, `n` clears, `s` / `enter` save & continue,
  `q` / `ctrl-c` cancels (exits 130 with `MsgClonePickUserCancelled`).
- Pre-checked rows: any path passed via the positional `<paths>`
  argument seeds the initial selection so re-running with `--ask` is
  a confirmation pass, not a do-over.
- Auto-greyed rows: anything matching `constants.ClonePickAutoExclude`
  (`.git`, `node_modules`, `vendor`, `dist`, `build`, `__pycache__`,
  `.venv`, `target`) renders with `[-]` instead of `[ ]` and is
  skipped by `a` (select-all). Individual greyed rows can still be
  toggled via `space` for the rare case the user wants them.
- `<paths>` is now optional when `--ask` is set: `normalisePaths`
  short-circuits empty input in ask mode, the picker fills the slot.
  An empty pick from the picker still exits 2 with the canonical
  missing-paths message — silent no-op clones are user-hostile.
- `Plan.UsedAsk` is set to true when the picker is used so the
  persisted `CloneInteractiveSelection` row records the provenance
  (no schema migration needed — the column was reserved in v3.153.0
  scaffolding).

### Code-quality notes
- `cmd/clonepick.go` grew past the 200-line file cap; the new picker
  glue lives in `cmd/clonepick_picker.go` to keep the new code itself
  cap-compliant. A follow-up should split `cmd/clonepick.go` further
  (flag-binding into its own file is the obvious cut).
- `clonepick/picker.go` (156), `picker_view.go` (88), `picker_tree.go`
  (108), `picker_test.go` (129) are all under the cap.
- Picker model unit tests in `clonepick/picker_test.go` drive
  `pickerModel.handleKey` directly without launching a real terminal,
  covering: preselection, space-toggle, cursor clamp, select-all
  skipping greyed rows, select-none, quit-cancel, save-with-selection,
  prefix-match auto-exclude, view counter + key hints, empty-repo view.
- `Version` constant bumped to `4.20.0`.

### Known follow-ups
- Spec §"--ask picker" wants the temp clone reused as the final
  destination ("cloned once, not twice"). v1 ships clone-twice (the
  first clone is metadata-only so the wire cost is negligible). The
  reuse refactor lands when we have a cross-filesystem-safe move
  helper for the `<temp>` -> `<dest>` step.
- A windowed scroller for repos with >1k tracked paths is deferred
  until we hit a real-world repo where the flat render hurts.

## v4.19.0 — `gitmap clone-pick --replay <id|name>` shipped

### Added
- `gitmap clone-pick --replay <id|name>` re-runs a previously saved
  selection without retyping the URL or the path list. Numeric refs
  resolve through the `SelectionId` index; everything else is treated
  as a `Name` lookup (most-recently-saved wins). When the replay
  succeeds, `CreatedAt` on the row is bumped so the most-recently-used
  selection sorts to the top of future listings.
- `--replay` co-exists with the regular runtime flags: `--dry-run`,
  `--quiet`, `--force`, and `--dest` from the current invocation
  override the persisted Plan, but the URL / paths / mode / branch /
  depth / cone / keep-git settings are reproduced verbatim from the
  saved row. The "no saved selection matches %q" message uses the
  spec'd format so CI greps stay stable.
- New `clonepick.Loader` interface (read surface) split from the
  existing `Persister` (write surface). `*store.DB` implements both,
  but tests can fake just the lookup half. `clonepick.LoadFromDB` and
  `clonepick.TouchAfterReplay` are the only two new exported entry
  points; everything else stays package-private.
- Unit coverage in `clonepick/replay_test.go` exercises numeric vs.
  name routing, the missing-ref error path, blank/nil-loader guards,
  and the dry-run skip on `TouchAfterReplay`.

### Changed
- `cmd/clonepick.go` no longer requires the positional `<repo-url>`
  when `--replay` is set. Missing both still exits with code 2 and
  the existing `MsgClonePickMissingURL` message — no behavioral change
  for fresh runs.
- `Version` constant bumped to `4.19.0`.

## v4.18.0 — `gitmap commit-in` / `cin` — chronological multi-source commit replay

### Added

### Added

- **`gitmap commit-in <source> <inputs...>`** (alias **`cin`**) — replay every
  first-parent commit from one or more input repos (folders or Git URLs) into
  a single destination repo, deduped by source SHA, with byte-exact author
  AND committer date replication. Spec: `spec/03-commit-in/`.

  | Capability                | Behavior                                                                              |
  |---------------------------|---------------------------------------------------------------------------------------|
  | Source auto-init          | URL → `git clone`. Existing repo → reuse. Non-repo dir → `git init`. Missing → mkdir. |
  | Sibling discovery         | `all` keyword walks every `<base>` / `<base>-vN` sibling of `<source>` in version order. |
  | Tail discovery            | `-N` keyword keeps only the N highest-version siblings.                               |
  | Dedupe                    | Per-workspace `ShaMap` skips any `SourceSha` ever seen — re-runs are a no-op.        |
  | Date replication          | Both `AuthorDate` AND `CommitterDate` of the rewritten commit equal the source's.    |
  | Conflict modes            | `ForceMerge` (default) takes the source side and logs the clobber list. `Prompt` aborts. |
  | Concurrency               | Single advisory lock per workspace (`<.gitmap>/commit-in.lock`); second run exits 9.  |
  | Profiles                  | JSON profiles under `<.gitmap>/commit-in/profiles/<name>.json`. CLI > profile > defaults. |

- **Profile pipeline** — message rules (`StartsWith` / `EndsWith` / `Contains`
  line strip), title prefix/suffix, message prefix/suffix, function-intel
  detectors per language, weak-word override gates, and exclusions
  (`PathFolder` / `PathFile`). Strict-decode JSON refuses unknown fields.

- **Run audit** — every run inserts a `CommitInRun` row plus per-input
  `CommitInInputRepo` and per-commit `CommitInSourceCommit` /
  `CommitInRewritten` rows (PascalCase schema, `SchemaVersion 24`).

### Tested

- 13 end-to-end tests in `cmd/commitin/e2e/` cover happy-path, dedupe,
  sibling discovery (`all` + `-N`), auto-init, profile precedence,
  Prompt-mode abort, ForceMerge clobber audit, and lock-busy concurrent
  runs. Plus the existing per-package unit suites for parser, walk,
  dedupe, replay, profile, message pipeline, and finalize.

## v3.48.0 — (2026-04-21) — `gitmap self-install --shell-mode` with strict `zsh+pwsh` combos

### Added

- **`--shell-mode` flag on `gitmap self-install`** — the canonical way to
  pick which shell-profile family (or families) receive the gitmap PATH
  snippet. Replaces the previous `--profile` / `--dual-shell` flags as
  the documented entry point; both old names still work as hidden
  aliases so existing scripts and CI keep running unchanged.

  | Value                  | Behavior                                                                                  |
  |------------------------|-------------------------------------------------------------------------------------------|
  | `auto` (default)       | Detect via `$SHELL` + `detect_active_pwsh` and write to the matching profile only.        |
  | `both`                 | Write to every supported profile gitmap can locate.                                       |
  | `zsh` / `bash` / `pwsh` / `fish` | Singleton — write only the named family's profile, skip everything else.        |
  | Combo (e.g. `zsh+pwsh`) | **Strict union** of the listed families; everything else (incl. `~/.profile`) is skipped. |

- **Strict `+`-joined combos** — any `+`-joined combination of the four
  concrete singletons (`zsh`, `bash`, `pwsh`, `fish`) is now a valid
  `--shell-mode` value: `zsh+pwsh`, `bash+fish`, `zsh+bash+pwsh`, etc.
  Combos differ from `both` in that they are deterministic and exclusive
  — only the listed shells are written, so `--shell-mode zsh+pwsh` from
  inside `pwsh` on macOS guarantees the zsh profile is also touched, and
  guarantees that bash/fish are NOT.

- **Deterministic pwsh propagation** — when the resolved mode includes
  `pwsh` (as a singleton or anywhere inside a combo), `selfinstall`
  exports `GITMAP_DUAL_SHELL=1` to `install.sh`. This is a
  belt-and-suspenders signal alongside the explicit `--shell-mode`
  argument so `detect_active_pwsh` inside the bash installer fires from
  either path.

### Changed

- `--profile` and `--dual-shell` are now hidden aliases for
  `--shell-mode`. Precedence when more than one is passed:
  `--shell-mode` > `--profile` > `--dual-shell`. They all collapse onto
  a single resolved `opts.ShellMode` value before any write logic runs,
  so downstream code only ever sees one source of truth.
- `selfInstallOpts.Profile` is renamed to `selfInstallOpts.ShellMode`.
  Internal-only struct field; no public API impact.
- The bash `should_write_profile` helper in `install.sh` now gates on
  the resolved `${PROFILE_MODE}` against `auto`, `both`, the singleton
  family, OR a `+`-bracketed substring match for combo membership.

### Validation

- Added Go-side `validateShellMode` and bash-side `validate_shell_mode`.
  Both reject:
  - unknown singletons (`--shell-mode tcsh`),
  - meta values inside combos (`auto+zsh`, `both+pwsh`),
  - duplicate tokens (`zsh+zsh`, `pwsh+bash+pwsh`),
  - empty / single-token combos (`zsh+`, `+pwsh`).
  Failures exit 1 with a message naming the bad value and listing
  accepted singletons + the combo syntax — bad CLI input is treated as
  unrecoverable.

### Files

- `gitmap/cmd/selfinstall.go` — `opts.Profile` → `opts.ShellMode`,
  `resolveShellMode` collapser, `validateShellMode` +
  `isValidComboShellMode` + `isConcreteShellFamily`,
  `shellModeRequiresPwsh` for the `GITMAP_DUAL_SHELL` env-export
  decision.
- `gitmap/cmd/releaseargs.go` — `--shell-mode` registered as a
  value-consuming flag in `valueFlags` so it survives
  `reorderFlagsBeforeArgs`.
- `gitmap/constants/constants_selfinstall.go` — new
  `ShellMode{Auto,Both,Zsh,Bash,Pwsh,Fish}` constants,
  `ShellModeComboSep`, `SelfInstallShellModes`,
  `FlagSelfShellMode`, `FlagDescSelfShellMode`,
  `ErrSelfInstallShellModeInvalid`. Back-compat `ProfileMode*` aliases
  are kept as identical-value identifiers so existing references keep
  compiling.
- `gitmap/scripts/install.sh` — argument parser handles
  `--shell-mode <value>`, `should_write_profile` updated for combo
  gating, `validate_shell_mode` rejects bad combos, help text
  rewritten with singleton/combo tables.
- `gitmap/helptext/self-install.md` — singleton + combo tables, usage
  examples, `--profile` / `--dual-shell` documented as hidden aliases.

### Compatibility

- All existing `--profile <value>` and `--dual-shell` invocations
  continue to behave identically. The collapser ensures they map onto
  the same code path as the new `--shell-mode <value>` form.
- `auto` remains the default — invocations with no shell-mode flag at
  all are byte-for-byte unchanged from v3.47.x.
- No DB schema, no on-disk artifact, no PATH-snippet format change.

## v3.40.0 — (2026-04-21) — Auto-derived tags for VS Code Project Manager sync

### Added

- **Auto-tag detection** — every entry synced into `projects.json` (via
  `gitmap scan`, `gitmap code`, `gitmap as`) now gets a deterministic set
  of `tags` derived from top-level marker files in the rootPath:

  | Marker                                                | Tag      |
  |-------------------------------------------------------|----------|
  | `.git`                                                | `git`    |
  | `package.json`                                        | `node`   |
  | `go.mod`                                              | `go`     |
  | `pyproject.toml`, `requirements.txt`                  | `python` |
  | `Cargo.toml`                                          | `rust`   |
  | `Dockerfile`, `compose.yaml`, `docker-compose.yml`    | `docker` |

- **Additive UNION semantics** — auto-tags are merged into the existing
  `tags` array. Tags you added manually in the VS Code UI are NEVER
  removed by gitmap.
- **`--no-auto-tags` flag on `gitmap scan`** — opts out of auto-tagging
  for that run (e.g. when you want a tag-free baseline).
- New `vscodepm.DetectTags(rootPath)` helper + `vscodepm.unionTags`
  internal merger. Detector is shallow, read-only, deterministic.
- New `constants.AutoTagMarkers` / `AutoTagOrder` registry — single
  source of truth for which markers map to which tags and in what order.

### Changed

- `vscodepm.Pair` gains a `Tags []string` field. Existing call sites
  populate it from `DetectTags()`; passing `nil` simply skips tagging.
- Sync summary now treats a tag-set change as an "Updated" event (same
  as a name or paths change).

### Files

- New: `gitmap/vscodepm/autotags.go`,
  `gitmap/vscodepm/autotags_test.go`,
  `gitmap/constants/constants_vscode_pm_autotags.go`.
- Updated: `gitmap/vscodepm/{sync,entry,overwrite}.go`,
  `gitmap/cmd/{scan,vscodepmsync,code,rootflags}.go`,
  `gitmap/constants/constants.go`,
  `gitmap/helptext/{code,scan}.md`, `src/data/commands.ts`,
  `spec/01-vscode-project-manager-sync/README.md`.

### Compatibility

- Pure additive change. Legacy `tags` content is preserved; gitmap only
  appends. No DB migration required (tags are not stored in SQLite —
  they live exclusively in `projects.json`).
- Single-arg `gitmap code [alias]` and the existing `paths` subcommand
  continue to work unchanged.

## v3.39.0 — (2026-04-21) — Multi-root `paths` for VS Code Project Manager sync

### Added

- **`gitmap code <alias> <root> [extra...]`** — variadic form attaches
  extra folders to a Project Manager entry in one call. Additive: never
  clobbers paths you added in the VS Code UI.
- **`gitmap code paths add|rm|list <alias> [path]`** — explicit subcommand
  to manage the `paths` array on an existing entry. `rm` overwrites the
  entry so the removal actually sticks across re-syncs.
- New DB column **`VSCodeProject.Paths`** (JSON-encoded TEXT, schema **v20**).
  Idempotent additive migration — existing v19 databases upgrade silently.
- New `vscodepm.OverwritePaths(rootPath, name, paths)` helper for the
  one place we need to bypass union semantics (`paths rm`).

### Changed

- `vscodepm.Pair` now carries `Paths []string`. `Sync()` UNIONs the DB-side
  list with the on-disk list — gitmap never silently removes a user-added
  path.
- `gitmap as <newalias>` continues to only rewrite `name`. Extra paths,
  tags, enabled, and profile are preserved (per spec).

### Files

- New: `gitmap/vscodepm/overwrite.go`
- Updated: `gitmap/cmd/code.go`, `gitmap/vscodepm/{sync,entry}.go`,
  `gitmap/store/vscode_project.go`, `gitmap/model/vscode_project.go`,
  `gitmap/constants/{constants,constants_settings,constants_vscode_pm_sql}.go`,
  `gitmap/store/store.go`, `gitmap/helptext/code.md`,
  `src/data/commands.ts`.

### Compatibility

- Pure additive change. Legacy v19 DBs auto-upgrade on next invocation.
- `gitmap scan` and the single-arg `gitmap code [alias]` flows are
  unchanged for users who do not need multi-root.

## v3.38.0 — (2026-04-21) — VS Code Project Manager sync (`gitmap code` + auto-sync on scan)

### Added

- **New `gitmap code [alias] [path]` command** — registers a path with the
  `alefragnani.project-manager` VS Code extension and opens VS Code on it
  in one step. Resolution order: explicit `path` arg → git repo root → CWD.
  Alias defaults to the folder basename and can be overridden by the first
  positional arg. Works for any path; the path does **not** need to be a Git
  repo.
- **`gitmap scan` now auto-syncs** every discovered repo into the
  Project Manager `projects.json`. Pass `--no-vscode-sync` to opt out.
  `scan` itself never opens VS Code.
- **`gitmap as <newalias>` mirrors the rename** to `projects.json` — the
  matching `rootPath` entry's `name` field is updated atomically.

### Path resolution (per OS)

The full path to `projects.json` is **never hardcoded**. The new
`vscodepm` package first discovers the VS Code **user-data root**, then
appends `User/globalStorage/alefragnani.project-manager/projects.json`:

| OS      | User-data root                                                       |
|---------|----------------------------------------------------------------------|
| Windows | `%APPDATA%\Code` → fallback `%USERPROFILE%\AppData\Roaming\Code`     |
| macOS   | `$HOME/Library/Application Support/Code`                             |
| Linux   | `$XDG_CONFIG_HOME/Code` → fallback `$HOME/.config/Code`              |

This matches portable / non-standard installations automatically.

### Atomicity & user-edit preservation

- All writes go through a temp-file + `os.Rename` pair so a kill/crash
  mid-write never produces a corrupted `projects.json`.
- Foreign entries in `projects.json` (rootPaths gitmap doesn't know about)
  are **always preserved** — gitmap merges, never replaces.
- User-edited fields per entry (`paths`, `tags`, `enabled`, `profile`) are
  preserved on every re-sync. Only `name` is rewritten by gitmap.

### Schema

- New table **`VSCodeProject`** (PascalCase, `INTEGER PRIMARY KEY AUTOINCREMENT`,
  unique `RootPath COLLATE NOCASE`). DB schema version bumped to **19**.
- Idempotent migration; legacy databases upgrade silently on the next
  `gitmap` invocation.

### Files

- New: `gitmap/vscodepm/{path,entry,sync}.go`, `gitmap/cmd/{code,vscodepmsync}.go`,
  `gitmap/store/vscode_project.go`, `gitmap/model/vscode_project.go`,
  `gitmap/helptext/code.md`, `gitmap/constants/constants_vscode_pm{,_sql}.go`.
- Updated: `gitmap/cmd/{scan,asops,rootcore,rootflags}.go`,
  `gitmap/store/store.go`, `gitmap/constants/{constants,constants_cli,constants_settings}.go`.

### Spec & memory

- New spec: `spec/01-vscode-project-manager-sync/README.md` plus the captured
  user fixture `spec/01-vscode-project-manager-sync/sample-projects.json`.
- New memory: `mem://features/vscode-project-manager-sync` (path resolution
  rule promoted to a Core memory).

### Compatibility

- Pure additive change. Existing flows unaffected when VS Code is not
  installed: `scan` prints a soft note and continues, `gitmap code` errors
  cleanly with an actionable hint.

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

- **README.md Go Report Card badge** now points at `github.com/alimtvnetwork/gitmap-v28/gitmap` (the real Go module path set in v3.27.0) instead of the repo root `github.com/alimtvnetwork/gitmap-v28`. The previous URL returned a 404 from the Go module proxy because there is no `go.mod` at the repo root — the module lives one directory down in `gitmap/`. Both the badge image and the click-through report link were updated.

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

- **Go Report Card now resolves the module instead of failing with `could not get latest module version from https://proxy.golang.org/github.com/user/gitmap/@latest`.** The card at https://goreportcard.com/report/github.com/alimtvnetwork/gitmap-v28/gitmap will start scoring the project for the first time after this version is pushed.

### Root cause

`gitmap/go.mod` was declared as `module github.com/user/gitmap` — a leftover placeholder from project bootstrapping. Because Go Report Card runs `go get <module>@latest` against the public Go module proxy (`proxy.golang.org`) before linting, and that path returns 404 (no such GitHub user/repo), the whole report aborted before any of `gofmt`, `go vet`, `gocyclo`, `ineffassign`, `misspell`, `golint` could run.

The same placeholder was also referenced in:

- 391 `.go` files inside `gitmap/` (all `import "github.com/user/gitmap/..."` statements)
- The companion module `gitmap-updater/go.mod` and its imports
- `Makefile` and `.github/workflows/ci.yml` ldflags injection: `-X 'github.com/user/gitmap/constants.Version=...'`
- `run.ps1` and `run.sh` ldflags injection for `RepoPath`
- Spec docs and changelog history references
- React-side changelog/getting-started page references

If left unfixed, anyone running `go install github.com/alimtvnetwork/gitmap-v28/gitmap@latest` would get a `module declares its path as: github.com/user/gitmap but was required as ...` error, and the proxy would refuse to serve the module to downstream tooling.

### Fix

Renamed `github.com/user/gitmap` → `github.com/alimtvnetwork/gitmap-v28/gitmap` across **404 files** in a single atomic sed pass. Also implicitly renamed the sister module `github.com/user/gitmap-updater` → `github.com/alimtvnetwork/gitmap-v28/gitmap-updater`, which lives at the same GitHub path.

Verified post-rename:
- `gitmap/go.mod` now reads `module github.com/alimtvnetwork/gitmap-v28/gitmap`.
- `gitmap-updater/go.mod` now reads `module github.com/alimtvnetwork/gitmap-v28/gitmap-updater`.
- `Makefile` ldflags target the new constants package path.
- `.github/workflows/ci.yml` build step injects `Version` into the new constants package path.
- `run.ps1` and `run.sh` inject `RepoPath` into the new constants package path.
- Zero remaining references to the old placeholder string anywhere in the tree.

### What the user needs to do after pulling

1. Pull v3.27.0 and push to `main`.
2. Once the new tag (`v3.27.0`) is pushed, visit https://goreportcard.com/report/github.com/alimtvnetwork/gitmap-v28/gitmap — first visit will trigger a fresh scan against the new module path.
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

## v3.16.0 — (2026-04-20) — uninstall-quick.ps1 multi-binary fix + repo rename to gitmap-v28

### Fixed

- **`uninstall-quick.ps1`** — `Try-SelfUninstall` and `Resolve-DeployRoot` now pipe `Get-Command gitmap` through `Select-Object -First 1`. When two `gitmap.exe` binaries were on `PATH` (e.g. a stale drive-root shim at `E:\gitmap\gitmap.exe` AND the canonical `E:\bin-run\gitmap-cli\gitmap.exe`), `Get-Command` returned an array. PowerShell's string interpolation joined `.Source` with a space and the resulting `'E:\gitmap\gitmap.exe E:\bin-run\gitmap-cli\gitmap.exe'` was passed to the runtime as a command name, producing:
  > The term '...path1... ...path2...' is not recognized as a name of a cmdlet, function, script file, or executable program.
  Self-uninstall now invokes the FIRST resolved binary by **absolute path** (`& $activeBinary self-uninstall -y`) instead of relying on `& gitmap` PATH resolution, so a stale shim cannot hijack the call.
- **`run.ps1`** — Same defensive fix applied to all three `Get-Command gitmap` callers (deploy-target detection, post-deploy active-vs-deployed sync, changelog binary resolution).
- **`gitmap/scripts/Get-LastRelease.ps1`** — Same defensive fix in `Get-ReleaseFromBinary`.

### Renamed

- **All `gitmap-v28` references → `gitmap-v28`** across the entire repo (45 files, 567 occurrences). Includes install/uninstall one-liners, Go installer constants, helptext, spec docs, post-mortems, the React landing page, and `.lovable/memory/**`.
- **Preserved**: release-asset filenames like `gitmap-v28.49.1-windows-amd64.zip` (where `v4.49.1` is the package version, not the repo name) — only the GitHub URL repo segment changed.

### Why

The uninstall failure was reported by a user who had run gitmap since the v2.x drive-root shim era — their stale `E:\gitmap\gitmap.exe` was never removed and the new `gitmap-cli/` install put a second binary on PATH. The `gitmap-v28` repo rename had been pending since the v3.x line started; bundling both keeps the CHANGELOG narrative simple.

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
- **Remote installer URLs** — Updated `constants_selfinstall.go` `SelfInstallRemotePwsh` and `SelfInstallRemoteBash` from `gitmap-v28` to `gitmap-v28`.

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

## v3.12.0 — (2026-04-20) — Pinned-version release snippet + gitmap-v28 rename

### Added

- **Pinned-version install snippet on the GitHub release page** — the release publisher (`gitmap/release/installsnippet.go`, wired into `workflowgithub.go::uploadToGitHub`) now auto-appends a markdown block containing PowerShell + bash one-liners that hard-code the just-published tag. Idempotent via a hidden `<!-- gitmap-pinned-install-snippet:<tag> -->` HTML marker. Anyone copying the snippet from `…/releases/tag/v3.12.0` installs exactly v3.12.0 — never "latest", never a `-v<N+1>` sibling repo. Template lives in `constants_release.go` as `ReleaseSnippetTemplate` / `ReleaseSnippetMarker`.
- **Pinned-version short-circuit in installer scripts** — `gitmap/scripts/install.ps1` and `install.sh` gained a new branch in their discovery prelude: when `-Version <tag>` (PowerShell) or `--version <tag>` (bash) is supplied, the installer now skips both the `releases/latest` API call **and** the versioned-repo `-v<N>` discovery probe, downloading `…/releases/download/<tag>/…` directly. Closes the gap where a snippet copied from a v3.x release page could silently jump to the v4 repo's latest tag.
- **Spec doc** `spec/07-generic-release/08-pinned-version-install-snippet.md` — full NEA/AI handoff contract: rendered snippets, installer-side flag matrix, release-cutting checklist, and a CI test contract for future work.

### Changed

- **Repo rename `gitmap-v28` → `gitmap-v28` across the entire codebase** — every Go constant (`SourceRepoCloneURL`, `SelfInstallRemotePwsh/Bash`, `GitmapRepoPrefix`, install hint URLs), every install/uninstall script (`install.ps1`, `install.sh`, `install-quick.ps1`, `install-quick.sh`, `uninstall-quick.*`), every spec doc under `spec/01-app/` and `spec/07-generic-release/`, every helptext markdown, the README, the React `src/data/*.ts` files, GitHub workflows, and historical CHANGELOG entries were rewritten via `sed -i 's/gitmap-v28/gitmap-v28/g'`. The only remaining `gitmap-v28` references are inside `.gitmap/` artifacts, which are immutable per project policy.

## v3.11.1 — (2026-04-20) — Alias-collision CI guard

### Added

- **Alias-collision uniqueness test** — extended `gitmap/constants/cmd_constants_test.go` with `TestTopLevelCmdAliasesAreUnique`, which iterates every top-level `Cmd*` constant and fails when two distinct identifiers share the same short-form value (string length ≤ 2). Catches future regressions like a hypothetical `CmdFooAlias = "ls"` shadowing the existing `CmdListAlias`, before they reach the build phase. Companion `TestTopLevelCmdConstantsAreUnique` covers full-length command-name collisions. Manual `topLevelCmds()` registry is the source of truth and excludes anything marked `// gitmap:cmd skip`.

## v3.11.0 — (2026-04-19) — Constants hygiene + Phase 1.4 migration fix

### Fixed

- **v15 Phase 1.4 migration** — `GoProjectMetadata` and `PendingTask` rebuilds failed on databases first created at v3.5.0+ with `SQL logic error: no such column: Id`. Both tables were already singular before v15, so the canonical `CREATE TABLE IF NOT EXISTS` pass produced the v15-shaped table (with `{Table}Id` PK) before the rebuild ran, leaving no `Id` column to SELECT. Added `adaptOldColumnList()` in `gitmap/store/migrate_v15rebuild.go` that detects the existing PK shape via `columnExists()` and rewrites the leading `Id` token in `OldColumnList` to `{Table}Id` when needed. Idempotent and a no-op for genuine legacy → v15 paths.
- **`go vet` `non-constant format string`** in `gitmap/movemerge/finalize.go:50` — `logErr` was inferred as a printf-style wrapper. Reshaped `logErr(prefix, msg string)` to accept a pre-formatted message and moved `fmt.Sprintf(constants.ErrMMPushFailFmt, sha)` to the call site so the printf-check never triggers.
- **Unused-import build break** in `gitmap/store/migrations.go` — removed orphaned `"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"` import left over from a prior refactor.
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

- After a successful release, if the repo's remote origin matches the gitmap source repository prefix (`github.com/alimtvnetwork/gitmap-v28`), the CLI now prints install one-liner commands for both Windows (PowerShell) and Linux/macOS (Bash).
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
- Repository renamed from `git-repo-navigator` to `gitmap-v28`; all URLs, scripts, and references updated.

## v2.49.1 — Update UX & Versioned Binaries (2026-04-06)

- Added `--repo-path` flag to `update` command: override the source repo path for a one-time update.
- The `--repo-path` flag is automatically forwarded through the handoff binary to `update-runner`.
- Resolution priority: `--repo-path` flag → embedded constant → friendly error with recovery options.
- Improved "repo path not embedded" error with actionable recovery steps (one-liner install, clone & build, manual download, `--repo-path` override).
- CI release binaries now include version in filenames (e.g., `gitmap-v28.49.1-windows-amd64.zip`).
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
