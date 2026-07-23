# Plan 03 — `cn` no-args → `find-next` bridge

Spec: [`spec/01-app/107-cn-find-next-bridge.md`](../../../spec/01-app/107-cn-find-next-bridge.md)

User-locked decisions:

- **Probe semantics**: spec 103 model — `ls-remote --sort=-v:refname` then
  shallow-clone-verify newest→oldest, depth=5.
- **Empty/stale cache**: auto-probe via the **`cn` bridge only**. Plain
  `gitmap find-next` stays strictly read-only.

## Phasing

### Phase 0 — Unblock spec 103 (prerequisite)

Spec 103 (`probe --depth N`) is currently planned/unshipped. The bridge
cannot work without it. This phase ships spec 103 verbatim:

- `gitmap/probe/walk.go` — new file: `WalkRemote(url, currentTag string,
  depth int) ([]VerifiedTag, error)`. Implements `ls-remote
  --sort=-v:refname` + shallow-clone-verify loop, stop on first failure.
- `gitmap/cmd/probe.go` — add `--depth N` flag, clamp to
  `[1, ProbeMaxDepth=10]`, default `1` (back-compat).
- `gitmap/store/version_probe.go` — keep insert API; bridge will call
  it once per verified tag.
- Migration: `ALTER TABLE VersionProbe ADD COLUMN IsPreRelease INTEGER
  NOT NULL DEFAULT 0` (idempotent).
- Tests: `TestProbeDepth_ClampMaxDepth`, `TestProbeDepth_StrictlyNewer`,
  `TestProbeDepth_StopOnFirstFailure`.

**Acceptance**: `gitmap probe --depth 5 <path>` inserts up to 5 rows.

### Phase 1 — Path-type detection

- `gitmap/cmd/cnscope.go` — new file: `detectCnScope(cwd string,
  db *store.DB) (cnScope, error)`. Order: git-repo first, then
  ScanFolder lookup, then `none`.
- `gitmap/store/scan_folder.go` — add `FindScanFolderByPath(abs
  string) (*model.ScanFolder, error)` if not already present.
- Constants: `ErrCnBridgeNoScope` in `constants_messages.go`.

**Acceptance**: unit tests for all three branches with temp git inits +
in-memory SQLite.

### Phase 2 — Cache hydration + staleness

- `gitmap/store/find_next.go` — extend `FindNext` to also return the
  set of repos with **no** probe row or **stale** probe row, so the
  bridge can compute the probe set in one query.
- New constant: `FindNextStaleAfter = 24 * time.Hour` in
  `constants_cn_bridge.go`.
- `gitmap/cmd/cnbridge.go` — new file: `resolveProbeSet(scope cnScope,
  flags cnBridgeFlags) (toProbe []model.ScanRecord, fromCache
  []model.FindNextRow, error)`.

**Acceptance**: unit tests covering empty cache, fresh cache,
mixed-stale cache, `--no-probe`, `--refresh`.

### Phase 3 — Parallel probe execution

- `gitmap/probe/pool.go` — new file: `RunWalkPool(targets
  []model.ScanRecord, depth, workers int) <-chan WalkResult`.
  Outer parallel, inner sequential (per spec 103).
- Wire from the bridge; persist each `VerifiedTag` via
  `db.RecordVersionProbe`.

**Acceptance**: `TestRunWalkPool_Parallel` (10 repos, 4 workers,
fakes return predictable timings, asserts wall-time < sequential).

### Phase 4 — Interactive TUI summary

- `gitmap/tui/cnbridge_model.go` — new Bubble Tea model based on the
  existing `interactive` infra (spec 43).
- Keys: `↑/↓ j/k space a n enter r q esc`.
- Non-TTY fallback: print plain summary, require `--yes` to act.

**Acceptance**: `TestCnBridgeModel_KeyBindings` table-driven
covering toggle, select-all, one-by-one, refresh, cancel.

### Phase 5 — Parallel update execution

- `gitmap/cmd/cnbridgeupdate.go` — new file: `runUpdates(selected
  []bridgeRow, workers int) updateSummary`. Each worker shells out
  to the existing `runCloneNext` path with `v++` resolved against the
  highest verified tag (no need to re-derive).
- Stream progress through a channel into the TUI region.
- Failures recorded but non-fatal (Code-Red logged).

**Acceptance**: `TestRunUpdates_PartialFailure` (3 repos, middle one
fails, asserts other two succeed and summary is accurate).

### Phase 6 — Persistence trace column

- Migration: `ALTER TABLE RepoVersionHistory ADD COLUMN
  TriggeredByProbeId INTEGER DEFAULT NULL`.
- `store/version_history.go` — extend insert API with optional
  trigger id.
- Bridge passes the verified `VersionProbeId` when invoking
  `runCloneNext`.

**Acceptance**: integration test inserts probe → upgrade → asserts
`RepoVersionHistory.TriggeredByProbeId` is populated.

### Phase 7 — React UI

- `src/pages/FindNext.tsx` — new page with mermaid workflow diagram,
  flag table, examples.
- `src/App.tsx` — register `/find-next` route.
- `src/pages/Commands.tsx` — add `find-next` (`fn`) row + `cn`
  (no-args) callout linking to `/find-next`.
- `src/pages/CloneNext.tsx` — add "No-args bridge" section.

**Acceptance**: `bun run build` clean, manual smoke of the route.

### Phase 8 — Help + completion + memory

- `gitmap/helptext/find-next.md` — extend with bridge mention.
- `gitmap/helptext/clone-next.md` — document `cn` no-args bridge.
- `gitmap/completion/powershell.go` — extend `cn` completion with
  `--yes --select --refresh --no-probe --probe-workers --update-workers`.
- `mem://features/cn-find-next-bridge` — feature memory file (this plan
  links to it on completion).
- `.lovable/memory/index.md` — add bullet under Memories.

**Acceptance**: completion test passes, help renders ≤120 lines.

## Risks

- **Spec 103 size**: not yet sized. If it slips, the bridge slips too.
- **TTY detection on Windows Terminal**: `term.IsTerminal(os.Stdout.Fd())`
  is reliable, but pipelines like `gitmap cn | tee log.txt` need the
  non-TTY fallback to be obvious. Print a hint when forced
  non-interactive: "stdout is not a tty; pass --yes to update or
  --select to force interactive".
- **Worker count tuning**: `runtime.NumCPU()` may be wrong for I/O-bound
  git ops. Expose `--probe-workers` and `--update-workers` early so
  users can tune; default may need revisiting after telemetry.
- **Cross-platform shallow-clone behavior**: spec 103's
  `--filter=blob:none --no-checkout` requires git ≥2.27. Detect and
  fall back to `--depth 1 --no-checkout` on older clients.

## Out of scope (will not do as part of this plan)

- Changing the existing `find-next` reader contract.
- Touching `cn vN` / `cn v++` / `cn --all` / `cn --csv` paths.
- Auto-resolving merge conflicts during update (matches spec 59 — clone
  is always fresh into base-name folder).
- Probing intermediate versions' CHANGELOG content.
