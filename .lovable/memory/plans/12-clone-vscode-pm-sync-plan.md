# Plan — Clone → projects.json sync (v4.16.0)

## Status: implemented v4.16.0

## Goal
Wire every clone variant into the existing `vscodepm.Sync` so that
`projects.json` is updated automatically after any successful clone.

## Pieces
1. **Shared helper** `cmd/clonepmsync.go`:
   - `buildClonePMPair(absPath, repoName) vscodepm.Pair` — adds auto-tags.
   - `syncClonedReposToVSCodePM(pairs []vscodepm.Pair, skip bool)` — single
     `vscodepm.Sync` call + soft-fail via `reportVSCodePMSoftError`.

2. **Flag plumbing.** Add `NoVSCodeSync bool` to:
   - `CloneFlags` (rootflags.go)
   - `CloneNextFlags` (clonenextflags.go)
   - `cloneFromFlags` (clonefrom_flags.go)
   - `cloneNowFlags` (clonenow.go)
   - `clonePickParsed` (clonepick.go)

   Each flag set binds `constants.FlagNoVSCodeSync` and passes the value
   into the helper at the call site.

3. **Call sites:**
   - `clone.go::executeDirectClone` — single pair after upsert. Covers
     direct URL clone, cfr, cfrp.
   - `clone.go::runCloneMulti` — collect successful URLs, single batch
     sync at end.
   - `clone.go::executeClone` — after `cloner.CloneFromFileWithOptions`,
     iterate `summary.Cloned` for the manifest path.
   - `clonenext.go::runCloneNext` — after `recordVersionHistory`, before
     `openInVSCode`.
   - `clonefrom.go::runCloneFromExecute` — after `results` returned,
     filter by `CloneFromStatusOK`, build pairs from `Row.URL` +
     `Result.Dest` (resolve relative to cwd).
   - `clonenow.go::runCloneNowExecute` — after `results`, filter by
     `CloneNowStatusOK`, build pairs from `Row.RepoName` + `Result.Dest`.
   - `clonepick.go::runClonePickExecute` — after `clonepick.Execute`,
     when `Status == StatusOK`, sync 1 pair using `Result.Detail`
     (the dest path) + `plan.Name` (or fallback).

4. **Version bump + changelog.** v4.15.1 → v4.16.0 (minor: new feature,
   no breaking changes). CHANGELOG entry under "Added".

## Out of scope (followups)
- `--no-auto-tags` per-clone-command flag.
- Mirroring DB upsert (`UpsertVSCodeProject`) from clone variants.
