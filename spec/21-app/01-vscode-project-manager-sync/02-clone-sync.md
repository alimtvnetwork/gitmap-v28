# 02 — Clone-Time Sync to projects.json

> Status: **Spec locked, implementation in progress (v4.16.0)**
> Extends: [`README.md`](./README.md) (the master spec for the
> `alefragnani.project-manager` `projects.json` schema).
> Origin: user report — "when I do clone or CFRP, the new repo never
> appears in Project Manager. Clone is the root thing; behavior should
> be the same everywhere."

## 1. Goal

Every gitmap-v28 command that **lands a fresh repo on disk** must, as part
of the same invocation, append (or update) the matching entry in the
VS Code Project Manager `projects.json` file documented in
`README.md`. No second command should be required.

Today only `gitmap-v28 code <path>` and `gitmap-v28 as <newname>` write to
`projects.json`. The clone family does not — that is the bug this
spec closes.

## 2. Scope — every clone variant

| Command                          | Entry point                          | Notes                                                          |
|----------------------------------|--------------------------------------|----------------------------------------------------------------|
| `gitmap-v28 clone <url>`             | `cmd/clone.go::executeDirectClone`   | Single direct URL. Already opens VS Code; now also syncs.      |
| `gitmap-v28 clone <url1> <url2> ...` | `cmd/clone.go::runCloneMulti`        | Each successful per-URL sub-clone is synced individually.      |
| `gitmap-v28 clone <manifest>`        | `cmd/clone.go::executeClone`         | After the cloner returns, every `Cloned[i]` is synced.         |
| `gitmap-v28 clone-next vX`           | `cmd/clonenext.go::runCloneNext`     | Single repo, after the flattened clone succeeds.               |
| `gitmap-v28 clone-from <file>`       | `cmd/clonefrom.go::runCloneFromExecute` | After execute, iterate `results` and sync every status=ok row. |
| `gitmap-v28 clone-now <file>` / `reclone` | `cmd/clonenow.go::runCloneNowExecute` | Same pattern as clone-from.                                    |
| `gitmap-v28 clone-pick <url> <paths>`| `cmd/clonepick.go::runClonePickExecute` | Single repo, after Status=ok. `Detail` carries the dest path.  |
| `gitmap-v28 clone-fix-repo` (cfr)    | routes through `executeDirectClone`  | Inherits the sync automatically — no extra wiring.             |
| `gitmap-v28 clone-fix-repo-pub` (cfrp) | routes through `executeDirectClone`| Same as cfr — the user-named root cause case.                  |

Skipped clones (status `skipped` because the destination already
exists) are NOT synced, because the entry should already exist from a
prior run; pushing again would be a no-op but would still tick the
`Updated` counter.

Failed clones (status `failed`) are NEVER synced — there is no real
folder to point Project Manager at.

## 3. Behavior contract

1. **One projects.json transaction per clone command, not per repo.**
   The cmd-side code aggregates every successful clone in a single
   `[]vscodepm.Pair` slice and calls `vscodepm.Sync` exactly once at
   the end of the command. This matches what `gitmap-v28 scan` already
   does and avoids racing the atomic-rename writer in
   `vscodepm/sync.go::writeEntriesAtomic` against itself.

2. **Auto-tags ON by default.** Each pair's `Tags` is populated from
   `vscodepm.DetectTags(absPath)` — same convention as `gitmap-v28 code`.
   Mirrors the existing `--no-auto-tags` semantics; out of scope for
   this spec to add a `--no-auto-tags` flag to clone variants (can be
   added later if requested).

3. **Soft-fail.** When the user-data root or extension dir is missing
   (CI, headless, no VS Code installed), the helper logs a one-line
   note via the existing `reportVSCodePMSoftError` and the clone
   command exits successfully. Sync errors NEVER turn a successful
   clone into a failed exit code.

4. **Opt-out:** every clone command accepts `--no-vscode-sync`
   (the same flag name `gitmap-v28 scan` already uses; constant
   `constants.FlagNoVSCodeSync`). When set, the helper prints the
   existing `MsgVSCodePMSyncSkipped` message and returns immediately.

5. **Path discipline.** The pair's `RootPath` MUST be the absolute
   path written to disk (the same value the cloner used as `git
   clone <url> <here>`). Never pass a relative path — the
   normalization in `vscodepm.normalizePath` is case-folded on
   Windows, but only after `filepath.Clean`, which can collapse the
   wrong segments if a relative path slips through.

6. **Name field.** Use the repo name as derived by the existing
   helpers (`repoNameFromURL`, the manifest row's `RepoName`, or
   `clonenext.ParseRepoName(...).BaseName` for the flattened
   clone-next case). This matches the user's `projects.json` sample
   where the `name` is the folder basename, not the full URL slug.

## 4. Shared helper

New file `cmd/clonepmsync.go`:

```go
// syncClonedReposToVSCodePM runs vscodepm.Sync for every pair, honoring
// the --no-vscode-sync flag and soft-failing when VS Code or the extension
// is missing. Wired into every clone variant after a successful clone.
func syncClonedReposToVSCodePM(pairs []vscodepm.Pair, skip bool)
```

The single-repo clone variants build a 1-element slice; the manifest
variants build the slice from successful results.

## 5. Out of scope (for v4.16.0)

- A dedicated `--no-auto-tags` flag on each clone command.
- DB-side `UpsertVSCodeProject` upsert from clone (kept scoped to
  `gitmap-v28 scan` and `gitmap-v28 code` for now). The `projects.json`
  reconciliation is the user-visible behavior; the DB mirror can
  follow in a later patch.
- Migrating the existing `gitmap-v28 code` `syncCodeEntry` to use the new
  helper — the two have slightly different shapes (paths union vs not).

## 6. Test surface

- Unit: `vscodepm.Sync` already has merge/atomic-write tests.
- Integration: a smoke test per clone variant that asserts a
  newly-cloned repo appears as a new entry under a tmpdir-rooted
  `projects.json` (use `t.Setenv("APPDATA", ...)` on Windows / `HOME`
  on Unix).
- Soft-fail: a test that points the env at an empty dir and asserts
  the clone still exits 0 with the `MsgVSCodePMSectionHeader`
  "VS Code not detected — sync skipped" line on stdout.
