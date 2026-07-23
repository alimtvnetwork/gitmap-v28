# Pending Issues

## 01 — Unit Test Coverage Gaps
- **Status**: Open since v2.49.0
- **Description**: Missing unit tests for `task`, `env`, and `install` command families
- **Impact**: Low — commands work but lack automated regression coverage
- **Blocked By**: Nothing — can be done anytime
- **Files Affected**: `cmd/task*.go`, `cmd/env*.go`, `cmd/install*.go`

## 02 — Install --check Missing "Not Found" Message
- **Status**: Open since v2.49.0
- **Description**: `gitmap install --check <tool>` doesn't print a distinct message when a tool is not installed; constant was added but wiring is incomplete
- **Impact**: Low — tool still works, just poor UX for missing tools
- **Files Affected**: `cmd/installtools.go`

## 03 — Docs Site Navigation Missing Pages
- **Status**: Open since v2.76.0
- **Description**: `version-history` and `clone` pages exist but are not linked from the sidebar or commands page navigation
- **Impact**: Low — pages exist at `/version-history` and users won't discover them organically
- **Files Affected**: Sidebar component, `src/data/commands.ts`

## 04 — Helptext/env.md Missing --shell Examples
- **Status**: Open since v2.49.0
- **Description**: The `--shell` flag was wired into env commands but the help text file doesn't demonstrate usage
- **Impact**: Low — flag works but users won't know about it from `gitmap help env`
- **Files Affected**: `helptext/env.md`

## 05 — Clone-Next Missing --dry-run Support
- **Status**: Open (feature gap)
- **Description**: The flatten spec (87-clone-next-flatten.md) mentions `--dry-run` for previewing clone-next actions but it's not implemented
- **Impact**: Medium — users can't preview destructive folder removal before it happens
- **Files Affected**: `cmd/clonenext.go`, `cmd/clonenextflags.go`, `constants/constants_clonenext.go`

## 06 — Multi-URL Clone: PowerShell Comma-Splitting Crash (FIXED v3.80.0)
- **Status**: Fixed in v3.80.0
- **Reported**: User ran `gitmap clone url1,url2,url3` in PowerShell on Windows; got `fatal: could not create leading directories of 'D:\...\https:\github.com\alimtvnetwork\email-reader-v3.gitmap-tmp-...': Invalid argument`
- **Root Cause**:
  1. PowerShell on Windows silently splits unquoted comma-separated arguments into multiple `argv` entries when invoking external executables. So `url1,url2,url3` arrived as three separate `os.Args` entries, not one string.
  2. `parseCloneFlags` only inspected the first two positional args: `Arg(0)` became the source URL, `Arg(1)` was treated as the **folder name**.
  3. `executeDirectClone` then called `filepath.Abs("https://github.com/.../email-reader-v3")`, producing the nonsense Windows path `D:\...\https:\github.com\alimtvnetwork\email-reader-v3` (illegal because `:` is reserved after the drive letter).
  4. The replace-strategy code then tried to `os.RemoveAll` and `git clone` into that path, both of which fail with "filename, directory name, or volume label syntax is incorrect" / "could not create leading directories".
  5. Spec `01-app/104-clone-multi.md` and `mem://features/clone-multi` had been **planned for v3.38.0 but never implemented** — the parser still assumed exactly one source.
- **Solution**:
  1. New `flattenURLArgs([]string) []string` (`gitmap/cmd/clonemulti.go`) — splits each positional arg on `,`, trims whitespace, drops empties, dedupes case-insensitively (normalising trailing `.git`), preserving first-seen order. Accepts both `a b c` and `a,b,c` and mixed `a,b c d,e`.
  2. `parseCloneFlags` now returns a `CloneFlags` struct exposing the **full positional slice** (not just `Arg(0)`/`Arg(1)`).
  3. `resolveCloneFolderName` defensively returns `""` when the second positional arg looks like a URL — so even single-URL invocations can't be misinterpreted as `<url> <folder=other-url>`.
  4. `runClone` detects multi-URL form (any positional contains `,`, or 2+ positionals where both Arg(0) and Arg(1) parse as URLs) and dispatches to the new `runCloneMulti` worker which calls a non-fatal `executeDirectCloneOne` per URL, continuing on failure.
  5. Exit codes per spec: `0` all OK, `1` partial failure, `3` all URLs invalid.
- **Files Affected**:
  - `gitmap/cmd/clone.go` — new `runClone` dispatcher + `shouldUseMultiClone` + `runCloneMulti`
  - `gitmap/cmd/clonemulti.go` (new) — `flattenURLArgs`, `classifyURLs`, `executeDirectCloneOne`, `resolveCloneFolder`, `normaliseURLKey`
  - `gitmap/cmd/rootflags.go` — `CloneFlags` struct, `isLikelyURL` guard
  - `gitmap/constants/constants_clone.go` — `MsgCloneInvalidURLFmt`, `MsgCloneSummaryMultiFmt`, `MsgCloneRegisteredInline`, `MsgCloneMultiBegin`, `MsgCloneMultiItem`, `ErrCloneAllInvalid`, `ErrCloneMultiFailedFmt`, `ExitCloneMultiPartialFail`, `ExitCloneMultiAllInvalid`
  - `gitmap/constants/constants.go` — version bumped to `3.80.0`
- **PowerShell Note**: Even after this fix, users should prefer space-separated URLs in PowerShell to avoid relying on PS's implicit comma-splitting (which differs across PS 5.1 / 7.x). Both forms now work either way.

## 07 — URL Shortcut: `gitmap <url>` Should Auto-Clone (FIXED v3.81.0)
- **Status**: Fixed in v3.81.0
- **Reported**: User ran `gitmap https://github.com/...,https://...,https://...` (omitting the `clone` subcommand) and got `Unknown command: https://github.com/...`. Same with single-URL `gitmap <url>` and any GitHub/GitLab/SSH URL.
- **Root Cause**: `Run()` treated `os.Args[1]` strictly as a subcommand name and dispatched it through `dispatchCore`/`dispatchRelease`/etc. A bare URL has no matching subcommand, so it fell through to `ErrUnknownCommand`. There was no shortcut layer between argv and dispatch.
- **Solution**: In `gitmap/cmd/root.go` `Run()`, immediately after migration runs, check if `os.Args[1]` looks like a git URL via the existing `isLikelyURL` helper (matches `https://`, `http://`, `git@`). If yes, rewrite argv to `[binary, "clone", <original args...>]` so the existing multi-URL clone pipeline (issue 06) handles it. Single URL, comma-list, or space-separated URLs all work — `runCloneMulti`'s `flattenURLArgs` covers all forms.
- **Files Affected**:
  - `gitmap/cmd/root.go` — argv-rewrite shortcut before alias extraction and dispatch
  - `gitmap/constants/constants.go` — version bumped to `3.81.0`
- **UX Note**: The shortcut only fires for URLs (HTTPS/SSH git). Local file paths, shorthands (`json`/`csv`/`text`), and all existing subcommands keep their current behaviour.

## 08 — CI Lint Failures: errorlint / gocritic / unparam (FIXED v3.81.1)
- **Status**: Fixed in v3.81.1
- **Reported**: `golangci-lint run` failed in CI with 3 NEW findings vs main baseline:
  1. `cmd/reinstall.go:125` — `errorlint`: `err.(*exec.ExitError)` type assertion fails on wrapped errors
  2. `committransfer/env.go:6` — `gocritic` (unlambda): `func() []string { return os.Environ() }` should be `os.Environ`
  3. `committransfer/replay.go:126` — `unparam`: `shouldSkipPath` parameter `info os.FileInfo` is never read
- **Root Cause**:
  1. **errorlint**: Direct type assertion on `error` only matches the outermost concrete type. If any wrapper (e.g. `fmt.Errorf("%w", err)`) sits between, the assertion silently fails and we'd report exit code `1` instead of the real script exit code. The project's `.golangci.yml` enables `errorlint` precisely to forbid this pattern (memory rule: "Use `errors.Is`" — same family applies for `errors.As`).
  2. **gocritic unlambda**: Wrapping a parameterless, same-signature function in another lambda is dead indirection — `os.Environ` already satisfies `func() []string`. The wrapper was a leftover from an earlier refactor that briefly took arguments.
  3. **unparam**: `shouldSkipPath` historically accepted `info os.FileInfo` to check `IsDir()`, but that check was lifted into both call sites (so the caller can return `filepath.SkipDir`). The parameter became dead weight; `unparam` correctly flagged it.
- **Solution**:
  1. `cmd/reinstall.go`: replaced the type assertion with `var exitErr *exec.ExitError; if errors.As(err, &exitErr) { ... }` and added `"errors"` to imports. Now correctly unwraps any future wrapping.
  2. `committransfer/env.go`: simplified to `var currentEnv = os.Environ` — same behaviour, no allocation, no indirection. Tests can still stub it (`currentEnv = func() []string { return ... }`).
  3. `committransfer/replay.go`: removed the unused `info os.FileInfo` parameter from `shouldSkipPath`; updated both call sites in `snapshotCopy` and `mirrorPrune`. Caller still has its own `info` in scope for the `IsDir()` branch after the skip check.
- **Files Affected**:
  - `gitmap/cmd/reinstall.go` — `errors.As` + import
  - `gitmap/committransfer/env.go` — direct method-value assignment
  - `gitmap/committransfer/replay.go` — signature + 2 call sites
  - `gitmap/constants/constants.go` — version bumped to `3.81.1`
- **Prevention**: All three rules (`errorlint`, `gocritic`, `unparam`) are already enabled in `.golangci.yml` — the issue was that they passed silently before the offending code was introduced. Going forward, run `golangci-lint run --path-prefix=gitmap` locally before pushing (or rely on the CI diff-vs-baseline gate which now catches this).

## 09 — Windows Update Cleanup Popup: `Windows cannot find '\\'` (FIXED v3.82.0)
- **Status**: Fixed in v3.82.0
- **Reported**: After a successful `gitmap update`, the terminal showed `→ Handing off cleanup to deployed binary: gitmap.exe update-cleanup`, then Windows displayed a popup: `Windows cannot find '\\'`. This repeated across multiple update attempts and the terminal showed no useful diagnostics.
- **Reproduction Context**:
  1. Run `gitmap update` on Windows from a deployed binary setup.
  2. Allow the update runner to finish build/deploy + migrations.
  3. At phase 3, the handoff copy tries to invoke the newly deployed binary with `update-cleanup`.
  4. Instead of cleanup running quietly, Windows Shell/`cmd` surfaces `Windows cannot find '\\'`.
- **Root Cause**:
  1. The bug was **not** inside `runUpdateCleanup()` itself. The failure happened *before cleanup started*, in `gitmap/cmd/updatehandoff_phase3.go` during the Windows phase-3 launch.
  2. `spawnDeployedCleanupWindows` originally built one flat shell command string and passed it to `cmd.exe /C`: `ping 127.0.0.1 -n 3 >nul & start "" /B "<deployed>" update-cleanup`.
  3. That pattern depended on fragile `cmd.exe` quoting semantics (`start` treats the first quoted token as a window title, and Go's Windows argument escaping adds another parsing layer). External Go/Windows reports match this exact failure mode, including the popup `Windows cannot find '\\'`.
  4. The handoff also discarded stdout/stderr and ignored the returned `Start()` error, so the CLI emitted **no useful diagnostics** even when the detached launch failed.
- **Solution**:
  1. Removed the fragile `cmd.exe /C ... start ...` handoff from `gitmap/cmd/updatehandoff_phase3.go`.
  2. Windows now launches the deployed binary directly with `exec.Command(deployed, constants.CmdUpdateCleanup)` instead of routing through `cmd/start`.
  3. Added a Windows-only hidden-process helper (`gitmap/cmd/processattr_windows.go`) so the cleanup process stays unobtrusive without embedding Windows-only fields in shared code.
  4. Added `GITMAP_UPDATE_CLEANUP_DELAY_MS=1500` plus `delayUpdateCleanupIfNeeded()` so the cleanup process waits briefly before deleting temp `.exe` / `.old` files.
  5. Cleanup handoff now prints the resolved target path and reports launch failures to `os.Stderr` instead of failing silently.
- **Files Affected**:
  - `gitmap/cmd/updatehandoff_phase3.go`
  - `gitmap/cmd/updatecleanup.go`
  - `gitmap/cmd/processattr_windows.go`
  - `gitmap/cmd/processattr_other.go`
  - `gitmap/cmd/selfuninstallhandoff.go`
  - `gitmap/constants/constants_update.go`
  - `gitmap/constants/constants.go`
- **Prevention**:
  1. Avoid string-built `cmd.exe /C start ...` launchers for internal handoffs.
  2. Never silence detached-launch failures in update-critical paths.
  3. Keep Windows-only process attributes in `_windows.go` files so non-Windows builds stay clean.

## 10 — Windows Update Cleanup Repeats After "Fix": PATH-First Handoff Targets Wrong Binary (FIXED v3.83.0)
- **Status**: Fixed in v3.83.0
- **Reported**: User kept seeing the update-cleanup failure repeatedly even after the earlier `cmd.exe` popup fix. The update appeared to complete, but cleanup still did not reliably run, and the console still lacked enough evidence to show *which binary* actually received `update-cleanup`.
- **Root Cause**:
  1. The earlier fix removed the fragile `cmd.exe /C start ...` launcher, but `resolveDeployedBinaryPath()` still resolved the cleanup target via `exec.LookPath("gitmap")` **before** checking the config-declared deployed location.
  2. On Windows machines with duplicate or stale `gitmap.exe` installs on `PATH`, Phase 3 could hand `update-cleanup` to the **wrong binary** — not the freshly deployed one that had just been updated.
  3. When the wrong binary was launched hidden/detached, the user saw the same cleanup problem repeat, but the terminal did not clearly reveal the resolution source (`PATH` vs config vs sibling) or the child PID, so the failure looked mysterious and "random".
  4. The embedded PowerShell comment in `constants_update.go` had also become stale: it still described the old `cmd.exe`-based handoff, which made future debugging and reasoning about the real runtime path harder.
- **Why Logs Still Felt Missing**:
  1. The handoff printed the target path, but **not** how that target was chosen.
  2. It did not print a started PID for the detached cleanup child.
  3. Invalid cleanup delay values were silently ignored.
  4. Without `--verbose`, there was no durable handoff trace explaining whether the cleanup ran inline, from config, from a sibling binary, or from a PATH fallback.
- **Solution**:
  1. Reordered deployed-binary resolution in `gitmap/cmd/updatehandoff_phase3.go` to prefer the config-declared deployed binary (`powershell.json` / deploy path) first, then sibling binary next to the handoff copy, and only use `PATH` as a last resort.
  2. Added explicit console output for the **resolution source** and the resolved cleanup target path.
  3. Added explicit console output for the detached cleanup child **PID** after a successful `Start()`.
  4. Added verbose-log entries for target resolution, inline cleanup, child start success, start failure, and missing-target cases.
  5. Added an explicit stderr error when no deployed cleanup target can be resolved at all.
  6. Added an explicit stderr + verbose warning when `GITMAP_UPDATE_CLEANUP_DELAY_MS` contains an invalid value instead of silently ignoring it.
  7. Corrected the stale embedded cleanup comment in `gitmap/constants/constants_update.go` so docs now match the real implementation.
- **Files Affected**:
  - `gitmap/cmd/updatehandoff_phase3.go` — prefer config/sibling over PATH; print source/path/pid; log handoff lifecycle
  - `gitmap/cmd/updatecleanup.go` — log cleanup start/finish and invalid delay values
  - `gitmap/constants/constants_update.go` — new handoff/log strings + corrected embedded update cleanup note
  - `gitmap/constants/constants.go` — version bumped to `3.83.0`
- **Console Evidence Added**:
  1. `→ Cleanup target resolved via: config|sibling|PATH`
  2. `→ Cleanup target path: ...`
  3. `→ Cleanup process started (pid=...)`
  4. `→ Cleanup binary: ...` inside `update-cleanup`
  5. explicit stderr message when the handoff target cannot be resolved or the delay env is invalid
- **Prevention**:
  1. In self-update flows, never trust `PATH` first when a config-declared deployed binary exists.
  2. Every detached handoff must log **which target** was selected and **why**.
  3. Embedded script comments/docs must be updated together with orchestration changes so future debugging is based on reality, not stale notes.
  4. Best-effort cleanup may stay non-fatal, but target-resolution failures must always be visible in the console and verbose log.

## 11 — `Unknown command: https://...` Recurs Even After v3.81.0 URL-Shortcut Fix (FIXED v3.84.0)
- **Status**: Fixed in v3.84.0
- **Reported (4th time)**: User typed `gitmap https://github.com/.../email-creator-v1,https://github.com/.../email-reader-v3,https://github.com/.../account-automator` (and space-separated and mixed comma+space variants). All three forms produced `Unknown command: https://github.com/alimtvnetwork/email-creator-v1`. Issue #07 logged this as fixed in v3.81.0 — yet it kept happening.
- **Root Cause** (two layers, why the same error keeps showing up):
  1. **Stale binary on PATH.** The URL-shortcut + multi-URL clone code IS present in the source tree (`gitmap/cmd/root.go`'s `isLikelyURL` rewrite + `flattenURLArgs` in `clonemulti.go`). The user's installed `gitmap.exe` on PATH is *older than v3.81.0* and simply does not contain that shortcut. Every recent `gitmap update` has been failing in phase-3 cleanup (issues #09 / #10), so the freshly built binary never actually reaches the deployed location — the user keeps running an old one.
  2. **Original v3.81.0 shortcut was too narrow.** `Run()` only checked `isLikelyURL(os.Args[1])`. If the user prepends a flag (`gitmap --verbose https://...`) the URL is in `os.Args[2]`, the shortcut misses, dispatch fails, and the user sees the same `Unknown command: --verbose` / `Unknown command: https://...` style failure. The shortcut needed to scan the full argv slice.
  3. **Error message gave no actionable hint.** `Unknown command: https://...` looked like a dead-end. Nothing pointed the user at `gitmap clone <url>` or `gitmap update`, so each retry produced the same opaque error and the same frustration.
- **Solution**:
  1. Replaced the single-position `isLikelyURL(os.Args[1])` check with `shouldRewriteToClone(os.Args[1:])`, which scans every positional arg (skipping leading flags) and accepts any token whose comma-split pieces look like a git URL. All four reported forms now redirect to `clone` automatically:
     - `gitmap url1,url2,url3`
     - `gitmap url1, url2, url3`
     - `gitmap url1 url2 url3`
     - `gitmap --verbose url1,url2`
  2. Added `ErrUnknownCommandURLHint` so when the offending token IS URL-shaped, the CLI now prints the explicit `gitmap clone <url>` form AND the `gitmap update` instruction with a note to reopen the terminal so PATH refreshes — instead of a dead-end error.
  3. Bumped version to `v3.84.0` so users can confirm via `gitmap version` whether their binary actually contains this fix.
- **Files Affected**:
  - `gitmap/cmd/root.go` — `shouldRewriteToClone` / `looksLikeURLToken` / `looksLikeFlag` helpers; argv scan instead of `os.Args[1]` only; URL-aware unknown-command branch
  - `gitmap/constants/constants_messages.go` — new `ErrUnknownCommandURLHint` constant
  - `gitmap/constants/constants.go` — version bumped to `3.84.0`

## 12 — Phase 3 `update-cleanup` Child Failures Were Not Logged Durably (FIXED v3.93.0)
- **Status**: Fixed in v3.93.0
- **Reported**: User repeatedly hit the same `update-cleanup` failure after the deployed-binary handoff and reported that there were still "no logs" even after earlier cleanup/handoff fixes.
- **Root Cause**:
  1. The Phase 3 deployed-binary handoff itself was already working at the outer lifecycle level: the code logged `resolve`, `start_ok`, `inline`, and `done` transitions.
  2. The missing evidence was **inside the detached cleanup child**. Several inner cleanup branches still wrote only to child stderr: `filepath.Glob` failures, `os.Remove` retry exhaustion, drive-root shim skip/remove branches, and `*.gitmap-tmp-*` swap-dir cleanup failures.
  3. On Windows that child is launched hidden/detached, so those stderr-only lines were exactly the diagnostics most likely to be lost. The user saw the handoff start but had no durable branch-level record of *which actual cleanup operation* failed.
- **Solution**:
  1. `gitmap/cmd/updatecleanup_remove.go` now mirrors `glob_error`, per-attempt `remove_retry`, final `remove_fail`, and `remove_ok` into the always-on handoff log and the optional `--debug-windows-json` sink.
  2. `gitmap/cmd/updatecleanup_extra.go` now mirrors the special-case cleanup branches too: `drive_root_skip`, `drive_root_remove_fail` / `drive_root_remove_ok`, `swap_glob_error`, and `swap_remove_fail` / `swap_remove_ok`.
  3. Added automated tests for stable handoff-log formatting and cleanup-child env/argv forwarding.
  4. Wrote RCA report `spec/02-app-issues/31-update-cleanup-phase3-observability-gap.md` so future work does not regress the diagnosis.
- **Files Affected**:
  - `gitmap/cmd/updatecleanup_remove.go`
  - `gitmap/cmd/updatecleanup_extra.go`
  - `gitmap/cmd/updatecleanup_handoff_test.go` (new)
  - `spec/02-app-issues/31-update-cleanup-phase3-observability-gap.md` (new)
  - `gitmap/constants/constants.go` — version bumped to `3.93.0`
- **Prevention**:
  1. Detached update-cleanup children must mirror **inner branch failures**, not just outer lifecycle events, into a durable channel.
  2. Any cleanup branch that still uses ad-hoc `fmt.Fprintf(os.Stderr, ...)` without `logHandoffEvent(...)` and structured JSON is a regression risk on Windows hidden-process launches.
  3. The handoff log + `--debug-windows-json` sink must evolve in lockstep with new cleanup branches so forensic coverage stays complete.
  4. Every recurrence report for update-cleanup should start by checking the durable handoff log and JSON sink before changing the handoff architecture itself.

## 13 — Docs UI Did Not Switch to VS Code Grading Despite Repeated Requests (FIXED v3.94.0)
- **Status**: Fixed in v3.94.0
- **Reported**: User repeatedly asked for the docs UI to use a VS Code-style color grading, but the work kept returning to update-cleanup/root-cause debugging instead of changing the visible frontend shell.
- **Root Cause**:
  1. The request was mis-prioritized because the surrounding conversation contained backend/update failures, so implementation focus drifted away from the explicit UI ask.
  2. The codebase already contained partial VS Code Dark+ inspired tokens in `src/index.css`, which created a false assumption that the visual request was already satisfied.
  3. The visible app still mixed non-VS styling cues: light/default shell states, green-primary emphasis, aurora hero treatment, and generic cards/buttons that did not read like an editor/workbench UI.
  4. Theme startup logic did not make the requested dark grading the primary/default docs experience.
- **Solution**:
  1. Reworked the docs shell (`DocsLayout`, `DocsSidebar`) into a VS Code-like explorer/header treatment with restrained borders, panel surfaces, blue-primary accents, and flatter controls.
  2. Retuned semantic tokens in `src/index.css` so both light and dark themes align with a VS/workbench palette, with dark mode now centered on editor-style neutrals and blue highlights.
  3. Updated the home page hero and reusable docs surfaces (`FeatureCard`, `InstallBlock`, `CodeBlock`) to remove the mismatched decorative styling and align cards/code panels with editor-like framing.
  4. Updated theme helpers so dark is the default fallback and the applied theme is read more reliably from DOM/storage state.
  5. Wrote RCA report `spec/02-app-issues/32-docs-ui-vscode-grading-missed-request.md` and synced release metadata.
- **Files Affected**:
  - `src/index.css`
  - `src/lib/theme.ts`
  - `src/components/docs/DocsLayout.tsx`
  - `src/components/docs/DocsSidebar.tsx`
  - `src/components/docs/FeatureCard.tsx`
  - `src/components/docs/InstallBlock.tsx`
  - `src/components/docs/CodeBlock.tsx`
  - `src/pages/Index.tsx`
  - `spec/02-app-issues/32-docs-ui-vscode-grading-missed-request.md` (new)
  - `gitmap/constants/constants.go` — version bumped to `3.94.0`
- **Prevention**:
  1. Treat explicit UI restyle requests as frontend work first, even if earlier messages involved backend bugs.
  2. Never infer “already done” from token comments alone; verify the rendered shell actually matches the requested visual direction.
  3. Sync version + both changelog sources for every user-visible fix so release state stays trustworthy.

## 17 — `clone <url1> <url2>` Treats URL #2 as Folder Name on Stale Binary (FIXED v3.95.0)
- **Status**: Fixed in v3.95.0
- **Reported (5th time)**: User ran `gitmap clone https://.../email-creator-v1,https://.../email-reader-v3,https://.../account-automator` in PowerShell. Output:
  ```
  pending task already exists for Clone at D:\...\https:\github.com\alimtvnetwork\email-reader-v3 (Id 1)
  Cloning email-creator-v1 into https://github.com/alimtvnetwork/email-reader-v3...
  fatal: could not create leading directories of 'D:\...\https:\github.com\alimtvnetwork\email-reader-v3.gitmap-tmp-...': Invalid argument
  ```
- **Root Cause**: The deployed `gitmap.exe` on the user's PATH is older than v3.80.0 and does not contain the multi-URL routing fix. Current source already has `runClone` → `shouldUseMultiClone(cf)` → `runCloneMulti` for any 2+ URL invocation, so reaching `executeDirectClone(url, folderName=<another URL>, ...)` is impossible in current code. Three pieces of evidence prove the binary is stale:
  1. The `Cloning email-creator-v1 into https://github.com/...` log line is the literal `MsgCloneURLCloning` format from `executeDirectClone` — the path current source refuses to take for multi-URL input.
  2. The orphaned `pending task already exists` row was inserted by an even older crashed run with the same broken URL-as-folder shape.
  3. The user's repeated `gitmap update` runs have been failing in Phase 3 cleanup (issues #09 / #10 / #12), so the freshly built binary never replaces the stale one on PATH.
- **Solution**:
  1. **Defensive guard in `executeDirectClone`** — if the resolved folder name itself parses as a URL (impossible in current source), refuse to build the broken `D:\...\https:\github.com\...` path and exit with `ErrCloneStaleBinaryFolderURL` naming the exact recovery commands: `gitmap doctor`, `gitmap update`, `gitmap pending clear --yes`, and the reminder to open a NEW terminal so PATH refreshes. The guard fires only when a stale binary is in use — and the message tells the user exactly that instead of letting git fail with `Invalid argument`.
  2. **Pinned regression tests** in `gitmap/cmd/clone_stale_binary_test.go` covering all three reported PowerShell argv shapes (comma-glued, space-split, comma+space mixed) plus the URL detector and the recovery message contents. CI now fails before the broken binary can ship.
  3. **No changes to multi-URL routing itself** — `shouldUseMultiClone` / `runCloneMulti` / `flattenURLArgs` have been correct since v3.80.0. The fix targets the failure mode (silent fall-through to a broken code path) and the user-facing recovery story.
  4. RCA: `spec/02-app-issues/33-stale-binary-clone-folder-url-guard.md`.
- **Files Affected**:
  - `gitmap/cmd/clone.go` — stale-binary guard in `executeDirectClone`
  - `gitmap/constants/constants_clone.go` — `ErrCloneStaleBinaryFolderURL`
  - `gitmap/cmd/clone_stale_binary_test.go` (new) — three pinned tests
  - `spec/02-app-issues/33-stale-binary-clone-folder-url-guard.md` (new) — RCA
  - `gitmap/constants/constants.go` — version bumped to `3.95.0`
  - `src/constants/index.ts` — UI version sync to `v3.95.0`
- **Action Required by User** (the source is correct, but PATH binary must be replaced):
  1. `gitmap doctor` — confirm the active binary version (will show <3.80.0).
  2. `gitmap update` — rebuild + redeploy. If Phase 3 cleanup still fails, the durable handoff log under `<TMP>/gitmap-update-handoff-*.log` (v3.87.0+) and `--debug-windows-json` sink (v3.91.0+) record exact branch-level evidence.
  3. **Open a NEW terminal** so PATH refreshes.
  4. `gitmap pending clear --yes` — drop the orphaned row.
  5. Retry — comma- or space-separated, PowerShell or bash, all work.
- **Prevention**:
  1. Any new "argv → command" routing must have a regression test pinning the exact PowerShell/bash invocations real users have reported.
  2. Code paths that would build provably-impossible filesystem paths (URL chars on Windows, control chars, etc.) must refuse early with an actionable message — never silently hand them to git.
  3. Recurrence reports for fixes already shipped in source must immediately check the deployed binary version (`gitmap doctor`), not the source.

## 14 — `--debug-windows` flag added for self-update handoff diagnostics (FIXED v3.86.0)
- **Status**: Fixed in v3.86.0
- **Reported**: Follow-up to #09 / #10. Even with the cleanup-target resolution lines (`→ Cleanup target resolved via: …`, `→ Cleanup target path: …`, `→ Cleanup process started (pid=…)`), the *child* `update-cleanup` process printed almost nothing about its own environment, so when cleanup misbehaved on Windows the user could not tell which env vars, deploy path, or PID the child actually saw. There was also no way to enable richer diagnostics ad-hoc without rebuilding with `--verbose` plumbed through.
- **Root Cause**:
  1. Phase-3 dispatcher (`scheduleDeployedCleanupHandoff`) printed resolution + child PID, but the child cleanup process (`runUpdateCleanup`) did not echo back its self path, parent PID, env, or the GOOS it observed.
  2. Phase-2 handoff (`launchHandoff`) had no diagnostic output at all — users could not see what argv/env was about to be passed to the handoff copy.
  3. `--verbose` writes to a log file, which is awkward for one-off Windows debugging where the user wants console output they can paste into a bug report.
- **Solution**:
  1. New `--debug-windows` flag on `gitmap update` (also activated by `GITMAP_DEBUG_WINDOWS=1`).
  2. Structured `[debug-windows]` dump printed to **stderr** at three lifecycle points:
     - Phase 2 (`launchHandoff` in `gitmap/cmd/update.go`) — before spawning the handoff copy.
     - Phase 3 dispatcher (`scheduleDeployedCleanupHandoff` in `gitmap/cmd/updatehandoff_phase3.go`) — wraps the entire dispatch with header/footer; per-spawn details printed by `dumpDebugWindowsHandoff` immediately before `cmd.Start()`; spawned child PID printed by `dumpDebugWindowsChildPID` immediately after.
     - Phase 3 child (`runUpdateCleanup` in `gitmap/cmd/updatecleanup.go`) — prints the same dump from inside the deployed binary so the user sees its own view.
  3. Flag/env propagation: `--debug-windows` is forwarded into the Phase 2 handoff copy and the Phase 3 cleanup child via **both** argv (`buildCleanupChildArgs`) and env (`buildCleanupChildEnv` sets `GITMAP_DEBUG_WINDOWS=1`). Either signal alone activates the dump, which makes the flag survive intermediate launchers that strip argv.
  4. Env keys printed are explicit and small (`GITMAP_DEBUG_WINDOWS`, `GITMAP_UPDATE_CLEANUP_DELAY_MS`, `GITMAP_DEBUG_REPO_DETECT`, `GITMAP_REPORT_ERRORS`, `GITMAP_REPORT_ERRORS_FILE`, `PATH`, `GITMAP_DEPLOY_PATH`) so the dump never leaks unrelated secrets from the process environment.
- **Files Affected**:
  - `gitmap/cmd/updatedebugwindows.go` (new) — dump helpers + flag/env detection.
  - `gitmap/cmd/updatehandoff_phase3.go` — header/footer + handoff dump + child PID dump + `buildCleanupChildArgs`/`buildCleanupChildEnv`.
  - `gitmap/cmd/updatecleanup.go` — dump runs at the start of `runUpdateCleanup`.
  - `gitmap/cmd/update.go` — `launchHandoff` forwards flag + env and prints Phase 2 dump.
  - `gitmap/constants/constants_update.go` — `FlagDebugWindows`, `EnvDebugWindows`, `MsgDebugWin*` constants.
  - `gitmap/helptext/update.md` — flag table updated.
  - `gitmap/constants/constants.go` — version bumped to `3.86.0`.
- **Prevention**:
  1. Every detached spawn in the self-update flow must carry an explicit, opt-in stderr-only diagnostic mode that survives the spawn boundary via both argv and env.
  2. Diagnostic env-key lists must be hand-curated, never `os.Environ()` in full — that would leak credentials into bug-report pastes.
  3. Any future addition to the cleanup handoff (extra phases, extra spawns) must extend the `[debug-windows]` dump in lockstep so the trace stays complete.

## 15 — Durable on-disk handoff log for self-update cleanup (FIXED v3.87.0)
- **Status**: Fixed in v3.87.0
- **Reported**: Follow-up to #14. Even with `--debug-windows`, failures during the detached Windows cleanup spawn could still vanish if an intermediate launcher (run.ps1 wrapper, hidden process attr, third-party AV) discarded stdout/stderr. There was no on-disk forensic trail for the dispatcher's resolution decision or the child's start status.
- **Root Cause**:
  1. Phase 3 dispatcher and the cleanup child only wrote to stdout/stderr and to the optional `--verbose` log. When stdout/stderr was redirected to NUL (or the verbose log wasn't enabled) every diagnostic disappeared.
  2. The verbose log is opt-in and per-process — there was no shared sink that both the dispatcher and the spawned child wrote to, so even with `--verbose` you'd get two separate files and have to correlate them by timestamp.
- **Solution**:
  1. New `gitmap/cmd/updatehandofflog.go` — daily, append-mode, mutex-serialized writer at `<TMP>/gitmap-update-handoff-YYYYMMDD.log`. Always-on; failures swallowed so logging can never disturb the update flow.
  2. Every Phase 3 lifecycle branch (`resolve`, `start_ok`, `start_fail`, `inline`, `target_missing`, `run_ok`/`run_fail` on Unix) and every cleanup-child branch (`start`, `delay`, `delay_invalid`, `done`) calls `logHandoffEvent(phase, event, fields)`.
  3. Each line carries `pid`, `ppid`, `goos`, RFC3339 UTC timestamp + sorted key=value fields, so dispatcher + child entries interleave cleanly in one file and can be diffed across runs.
  4. Path is surfaced in TWO always-visible places: `→ Handoff log file: <path>` printed once at the start of Phase 3, and `[debug-windows] handoff log file : <path>` inside the `--debug-windows` dump header.
- **Files Affected**:
  - `gitmap/cmd/updatehandofflog.go` (new)
  - `gitmap/cmd/updatehandoff_phase3.go` — `logHandoffEvent` calls + log-path print at top of dispatch
  - `gitmap/cmd/updatecleanup.go` — `logHandoffEvent` calls in cleanup + delay branches
  - `gitmap/cmd/updatedebugwindows.go` — log file path added to dump header
  - `gitmap/constants/constants_update.go` — `UpdateHandoffLogNameFmt`, `MsgUpdatePhase3LogFile`, `MsgDebugWinLogFile`
  - `gitmap/helptext/update.md` — new "Handoff log file" section with example log lines
  - `gitmap/constants/constants.go` — version bumped to `3.87.0`
- **Prevention**:
  1. Any handoff that crosses a process boundary needs a durable, always-on, shared on-disk sink — stdout/stderr alone is not forensically sufficient.
  2. Log file paths must be discoverable without reading source: print at runtime, document in helptext.
  3. Logger writes must never block or fail the caller; degrade silently on disk errors.
  4. Daily-named log files keep the file bounded without needing rotation logic.

## 16 — `gitmap pending clear` to remove orphaned/illegal pending tasks (FIXED v3.88.0)
- **Status**: Fixed in v3.88.0
- **Reported**: Follow-up to #11/#12. Even after the defensive guards in `executeDirectClone` / `executeDirectCloneOne`, **pre-existing** rows from older crashes still blocked subsequent runs with `pending task already exists for Clone at <bad-path>`. There was no surgical way to drop one row — only nuking the SQLite file or running raw SQL.
- **Root Cause**: The clone pipeline records a pending task before it begins. A crash mid-clone (file-lock, broken FS path, OS reboot) can leave the DB pointing at a target that no longer makes sense (or never made sense, e.g. a URL accidentally treated as a folder name). `runPending` could only **list** rows, and `do-pending` would just retry them — neither could selectively delete.
- **Solution**:
  1. New subcommand `gitmap pending clear [<mode>|<id>] [--dry-run] [--yes|-y]` dispatched from `runPending` when `os.Args[2] == "clear"`.
  2. Modes: `orphans` (default) drops rows whose `TargetPath` is missing on disk; `illegal` drops URL-shaped or Windows-illegal-char targets; `all` drops everything; `<id>` drops one row.
  3. Three classifiers in `pendingclear.go`:
     - `isURLShapedTarget`: matches `://` anywhere, or any of `http:`, `https:`, `ssh:`, `git:` followed by `\` or `/`. Catches the exact corruption pattern from issue #11 (`https:\github.com\...`).
     - `hasIllegalPathChar`: scans for `:` after the drive letter, plus `?`, `*`, `<`, `>`, `|`, `"`.
     - `isOrphanTarget`: `os.Stat`-based, with `filepath.Abs` resolution; resolution failures are conservatively treated as orphans.
  4. Safety rails: confirmation prompt unless `--yes`/`-y`; `--dry-run` previews; per-deletion log lines + final tally.
  5. New `DeletePendingTask(id) error` on `*store.DB` reuses existing `SQLDeletePendingTask` and returns `ErrPendingTaskNotFound` for unknown IDs.
- **Files Affected**:
  - `gitmap/cmd/pendingclear.go` (new)
  - `gitmap/cmd/pending.go` — dispatcher branch for `clear`
  - `gitmap/store/pendingtask.go` — new `DeletePendingTask`
  - `gitmap/constants/constants_pending_task_msg.go` — `MsgPendingClear*` and `ErrPendingClear*` constants
  - `gitmap/helptext/pending-clear.md` (new) — help page
  - `gitmap/constants/constants.go` — version bumped to `3.88.0`
- **Prevention**:
  1. Every queue-style table needs a deterministic, scoped purge command — not just a list view + a retry-all.
  2. Purge commands must default to the safest possible mode (here: orphans only) and require an explicit opt-in for destructive variants (`all`).
  3. Confirmation prompts are mandatory for any DB-mutating command unless `--yes` is explicitly passed.
  4. Path classifiers (`isURLShapedTarget`, `hasIllegalPathChar`) belong in the cleanup command, not in the clone path — the clone path already rejects bad inputs at the entry point (issues #11/#12); the cleanup command exists specifically to handle rows that predate those guards.

## 17 — Robust multi-URL clone parsing (PowerShell + bash) (FIXED v3.89.0)
- **Status**: Fixed in v3.89.0
- **Reported**: Follow-up to #11/#16. Three real failure modes still bit users after the v3.80 multi-URL feature shipped:
  1. `gitmap clone url1;url2` in bash (bash users naturally reach for `;`) produced a single ugly task and a "command not found" hint because semicolon wasn't a list separator.
  2. Copy-pasting a URL from PowerShell history or the browser carried a U+FEFF BOM or a U+200B zero-width space and produced a phantom invalid-URL warning.
  3. `gitmap clone url1 url2 url3` (space-only, no commas) cloned only the first two because `shouldUseMultiClone` only sampled `Positional[0]` and `Positional[1]`.
  4. `gitmap clone git@github.com:foo/bar.git` was misclassified as a file path by `isDirectURL` (which only knew `https://`/`http://`/`ssh://`) even though `isLikelyURL` already accepted it — silent disagreement between two helpers that should have been in lockstep.
- **Root Cause**:
  1. `flattenURLArgs` only split on `,` — no `;` support; no sanitisation of invisible runes; no smart-quote folding; no leading/trailing-separator stripping.
  2. `shouldUseMultiClone` had a 2-positional ceiling on its URL detection, so the third+ args were silently treated as folder names by the single-clone path.
  3. `isDirectURL` and `isLikelyURL` were defined separately and drifted: `isLikelyURL` accepted `git@`, `isDirectURL` did not.
- **Solution**:
  1. New `urlListSeparators = ",;"` constant; `splitOnURLSeparators` uses `strings.FieldsFunc` so both characters act as boundaries simultaneously.
  2. New `sanitizeURLToken` pipeline: `stripInvisibleRunes` (BOM, U+200B/C/D) → `replaceSmartQuotes` (curly → ASCII so wrapper-trim works) → `TrimSpace` → `trimMatchingWrappers` (only matched `'`/`"`/backtick pairs) → strip leading/trailing separators → final `TrimSpace`.
  3. `shouldUseMultiClone` rewritten with three triggers (any one sufficient): (a) any positional contains `,` or `;`; (b) 2+ positionals AND any arg beyond the first parses as a URL; (c) the first positional flattens to 2+ valid URLs.
  4. `isDirectURL` extended to accept `git@host:owner/repo` shorthand; doc-comment cross-references added between the two helpers so future edits keep them in lockstep.
  5. `isLikelyURL` also gained `ssh://` for symmetry.
- **Files Affected**:
  - `gitmap/cmd/clonemulti.go` — sanitisation pipeline + new helpers
  - `gitmap/cmd/clone.go` — three-trigger detection + SSH-shorthand recognition
  - `gitmap/cmd/rootflags.go` — `isLikelyURL` extension + cross-ref comment
  - `gitmap/constants/constants.go` — version bumped to `3.89.0`
- **Prevention**:
  1. Two helpers that classify the same thing (`isDirectURL` vs `isLikelyURL`) MUST cross-reference each other in their doc comments. Drift between them is a Code Red.
  2. Any "scan the positional args" heuristic that hardcodes `[0]` or `[1]` indices is a bug waiting to happen — always iterate the slice.
  3. Real-world URL input is never clean: BOM, smart quotes, zero-width spaces, and stray wrappers are the norm, not edge cases. Sanitise on every entry point.
  4. Empty/separator-only tokens after sanitisation must be dropped silently — emitting "invalid URL: ``" is worse than emitting nothing.
