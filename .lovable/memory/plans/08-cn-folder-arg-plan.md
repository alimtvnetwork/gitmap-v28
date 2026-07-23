# Plan 08 — `gitmap cn` folder-arg dispatch

Spec: [spec/01-app/111-cn-folder-arg.md](../../../spec/01-app/111-cn-folder-arg.md)

## Goal

Add three new invocation forms to `gitmap cn` without disturbing the
existing in-place / cross-dir-alias forms:

- `gitmap cn vX <folder>` (or `v+1`, `v++`)
- `gitmap cn <folder>` (defaults to `v++`)
- Both honor relative, absolute, `~`-prefixed paths.

## Phases

### Phase 1 — Spec & memory (DONE in this session)
- [x] Write `spec/01-app/111-cn-folder-arg.md`
- [x] Write this plan file (`08-cn-folder-arg-plan.md`)
- [x] Append memory entry to `.lovable/memory/index.md`

### Phase 2 — Dispatcher
- [ ] Extend `versionPattern` in `releaserebase.go` to accept
      `v++` and `v+N` (currently only matches `v?N.N.N`).
      Verify `looksLikeVersion("v+1") == true` via existing
      `releaserebase_test.go` style test.
- [ ] New file `gitmap/cmd/clonenextfolderdispatch.go`:
  - `tryFolderArgCloneNext(args []string) bool` — runs BEFORE
    `tryCrossDirCloneNext` in `runCloneNext`. Handles the three
    new forms, returns true when it dispatched, false to fall
    through.
  - `resolveCloneNextFolder(token string) (string, error)` —
    expands `~`, joins to cwd, stats, returns absolute path.
  - `isFolderShaped(token string) bool` — true if path-separator
    present OR `os.Stat` succeeds as a directory.
- [ ] Wire `tryFolderArgCloneNext` into `runCloneNext` ABOVE the
      existing `tryCrossDirCloneNext` call. Order matters: the new
      form is more specific (path-separator presence), so it must
      win over the alias fallback.

### Phase 3 — Tests
- [ ] `gitmap/cmd/clonenextfolderdispatch_test.go` covering the
      eight cases enumerated in the spec test matrix.
- [ ] `gitmap/clonenext/version_test.go` extended to assert
      `looksLikeVersion("v++")` and `looksLikeVersion("v+1")`.

### Phase 4 — Help + docs
- [ ] Update `gitmap/helptext/clone-next.md` to document the two
      new forms with realistic example simulations.
- [ ] Add a "Folder-arg forms" subsection to `README.md` Command
      Reference → Cloning & Sync block.
- [ ] Add a `commands.ts` entry on the docs site for the two new
      forms (label them as v3.117.0+).

### Phase 5 — QA + tag
- [ ] `go test ./gitmap/...`
- [ ] `golangci-lint run`
- [ ] Bump `constants.Version` to `3.117.0` (done in implementation
      session) and add CHANGELOG entry.

## Risks & mitigations

| Risk | Mitigation |
|---|---|
| Release-alias resolver matches a bare folder name before our path-stat does | `isFolderShaped` requires path-sep OR successful `os.Stat`-as-dir; bare alias names with no slash and no on-disk dir keep the existing semantics |
| `v+1` regex change accidentally matches an alias like `v+abc` | Pattern anchors to `v\+\+` OR `v\+\d+` only — letters after `+` reject |
| Two-positional folder-then-version arg reverses today's contract | Kept as-is — that's the cross-dir release-alias form; no breakage |

## Done when

A user can run `gitmap cn ~/dev/macro-ahk-v11` from any cwd and the
v12 clone appears alongside it, with the same UX as the in-place form.
