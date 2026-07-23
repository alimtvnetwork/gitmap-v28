# 107 — `cn` (no args) → `find-next` Bridge with Auto-Probe and Interactive Update

> Status: **planned** — targets v3.55.0
> Depends on: spec 59 (`clone-next`), spec 90 (`ScanFolder` + `VersionProbe`),
> spec 103 (`probe --depth N`, planned v3.36.0), v3.9.0 `find-next`,
> spec 43 (`interactive` TUI).

## Goal

Today, `gitmap cn` with no arguments errors out (`ErrCloneNextUsage`).
Today, `gitmap find-next` is read-only and never triggers a probe — the
user must run `gitmap probe --all` first.

This spec wires both into a single ergonomic flow:

```
gitmap cn         (no args, anywhere)
  └─► detect cwd type (git repo | scan-folder root | neither)
  └─► hydrate VersionProbe state from DB for that scope
  └─► if empty/stale → auto-run a parallel probe walk (spec 103, depth=5)
  └─► render interactive summary
  └─► user picks: yes-all | one-by-one | selected | default-all | cancel
  └─► parallel `cn v++` per selected repo
  └─► persist every step (probe rows + RepoVersionHistory rows)
```

`find-next` itself **stays read-only**. The new behaviour lives only in
the `cn`-no-args bridge.

## CLI Surface

| Form | Behavior |
|------|----------|
| `gitmap cn` (in a git repo) | Single-repo bridge: probe-walk that one repo, prompt to upgrade. |
| `gitmap cn` (in a scan-folder root) | Multi-repo bridge: probe-walk every repo under that ScanFolder, prompt for selection. |
| `gitmap cn` (neither) | Clear error: "Not a git repo and not a registered scan folder. Run `gitmap sf add <path>` first." Exit 1. |
| `gitmap cn --yes` | Non-interactive default-all: auto-update every available repo. |
| `gitmap cn --select` | Force interactive multi-select even when only one update is available. |
| `gitmap cn --no-probe` | Use only cached `VersionProbe` rows; never touch the network. Empty cache → "no updates available". |
| `gitmap cn --refresh` | Force a fresh probe walk even if cache is fresh. |
| `gitmap cn vX` / `cn v++` | **Unchanged** — existing single-repo explicit-version path (spec 59). |
| `gitmap cn --all` / `--csv <path>` | **Unchanged** — existing batch path (`clonenextbatch.go`). |

The bridge is triggered **only** by the empty-args case, so every existing
script keeps working.

## Workflow

### 1. Path-Type Detection

```go
type cnScope struct {
    Kind        cnScopeKind  // singleRepo | scanFolder | none
    RepoPaths   []string     // absolute paths to operate on
    ScanFolderID int64       // 0 when Kind == singleRepo
}
```

Detection order:

1. Is cwd inside a git repo (walk up looking for `.git`)? → `singleRepo`.
2. Is `filepath.Abs(cwd)` registered in `ScanFolder.AbsolutePath`?
   → `scanFolder`, load every `Repo` with that `ScanFolderId`.
3. Otherwise → `none`, exit with the clear error above.

If both are true (rare — a scan-folder root that is itself a git repo),
`singleRepo` wins. Document this in helptext.

### 2. DB Hydration

Query `VersionProbe` for the resolved repo set via the existing
`SQLSelectFindNext` (no scan filter) or `SQLSelectFindNextByScanFolder`
(scoped). This already returns latest-probe-per-repo where
`IsAvailable=1`, sorted newest-first.

Define staleness:

```go
const FindNextStaleAfter = 24 * time.Hour
```

For each repo in the resolved set, compute:

| Cache state | Trigger probe? |
|-------------|----------------|
| No `VersionProbe` row at all | Yes |
| Latest row `ProbedAt` older than `FindNextStaleAfter` | Yes |
| Latest row `ProbedAt` within window | No (use cache) |
| `--refresh` flag | Yes (always) |
| `--no-probe` flag | Never |

### 3. Auto-Probe (Spec 103 Semantics)

When the bridge needs to probe, it invokes a **parallel** version of
`probe --depth 5`:

- Outer loop (per-repo): worker pool, default `runtime.NumCPU()`,
  configurable via `--probe-workers N`.
- Inner walk (per-repo): **sequential** as required by spec 103
  §Concurrency (clones share the same remote).
- Each worker executes the spec-103 algorithm:
  1. `git ls-remote --tags --sort=-v:refname <url>`
  2. Take strictly-newer slice vs current.
  3. Trim to `--depth 5` (capped at `ProbeMaxDepth=10`).
  4. Newest → oldest, `git clone --depth 1 --filter=blob:none
     --no-checkout --branch <tag> <url> <tmpdir>`.
  5. Stop on first failed shallow clone (per spec 103 — "stop on first
     gap" means "stop after first verification failure").
  6. Insert one `VersionProbe` row per verified tag.

This requires **spec 103 to ship first**. The bridge spec is blocked on
spec 103.

### 4. Interactive Summary

When stdout is a TTY and `--yes` is not set, render a Bubble Tea
summary table reusing `gitmap/tui` (spec 43):

```
Available updates (3)
─────────────────────────────────────────────────────────────
  [x] macro-ahk          v11 → v12  (ls-remote, 0.4s ago)
  [ ] coding-guidelines  v6  → v8   (shallow-clone-tag, 1.2s ago)
                                      ↳ also: v7 verified
  [x] scripts-fixer      v3  → v4   (cached 14m ago)
─────────────────────────────────────────────────────────────
[a] yes-all   [n] one-by-one   [enter] update selected
[r] refresh   [q] cancel
```

Key bindings:

| Key | Action |
|-----|--------|
| `↑/↓` `j/k` | Move cursor |
| `space` | Toggle selection |
| `a` | Select all + confirm |
| `n` | Sequential one-by-one mode (prompt per repo) |
| `enter` | Update currently-checked rows |
| `r` | Re-run probe walk |
| `q` `esc` `ctrl+c` | Cancel without changes |

When stdout is **not** a TTY (CI, pipes), default to "summary + exit"
unless `--yes` is set. `--yes` triggers the parallel update without any
prompt — matches your "default — update all (non-interactive)" mode.

### 5. Parallel Update Execution

For each selected repo:

- Worker pool with same sizing knob (`--update-workers N`, default
  `runtime.NumCPU()`).
- Each worker invokes the existing `runCloneNext` single-repo path
  with `v++` resolved against the highest verified tag from the cache.
- Stream per-worker progress lines through a channel into the TUI
  progress region.
- A worker failure is logged (per Code-Red zero-swallow policy) and
  recorded as a failed `RepoVersionHistory` row variant — see
  §"Persistence" below — but does **not** abort the other workers.

### 6. Persistence

- Every probe verification → one `VersionProbe` row (existing spec 90/103).
- Every successful clone → one `RepoVersionHistory` row (existing spec 87).
- Add an optional `TriggeredByProbeId INTEGER DEFAULT NULL` column to
  `RepoVersionHistory` so an operator can trace which probe row caused
  which upgrade. Migration via `addColumnIfNotExists` (idempotent).

### 7. React UI

The docs site must list `find-next` and the `cn` bridge. Concretely:

| Page | Path | Source |
|------|------|--------|
| `FindNext` | `/find-next` | New: `src/pages/FindNext.tsx` |
| Update to `CloneNext` | `/clone-next` | Add a "No-args bridge" section linking to `/find-next` |
| Update to `Commands` index | `/commands` | Add `find-next` (`fn`) entry + `cn` (no-args) callout |

The `FindNext` page must include:

- The full process diagram from this spec (mermaid).
- Flag table (`--yes`, `--select`, `--refresh`, `--no-probe`,
  `--probe-workers`, `--update-workers`).
- Workflow diagram (single-repo vs scan-folder branches).
- Examples for each of the four user prompt modes.

## Constants

`gitmap/constants/constants_cn_bridge.go` (new):

```go
const (
    FindNextStaleAfter        = 24 * time.Hour
    DefaultProbeWorkers       = 0  // 0 → runtime.NumCPU()
    DefaultUpdateWorkers      = 0
    DefaultBridgeProbeDepth   = 5

    FlagCnYes            = "--yes"
    FlagCnSelect         = "--select"
    FlagCnRefresh        = "--refresh"
    FlagCnNoProbe        = "--no-probe"
    FlagCnProbeWorkers   = "--probe-workers"
    FlagCnUpdateWorkers  = "--update-workers"
)
```

`gitmap/constants/constants_messages.go` additions:

```go
const (
    ErrCnBridgeNoScope     = "cn: not a git repo and not a registered scan folder (cwd=%s)"
    MsgCnBridgeProbing     = "→ Probing %d repo(s) (depth=%d, workers=%d)..."
    MsgCnBridgeCacheHit    = "✓ Using cached probe data for %d repo(s) (run `gitmap cn --refresh` to re-probe)"
    MsgCnBridgeNoUpdates   = "✓ All %d repo(s) up to date."
    MsgCnBridgeUpdating    = "→ Updating %d repo(s) (workers=%d)..."
    MsgCnBridgeDone        = "✓ Updated %d, failed %d, skipped %d."
)
```

## Error Handling (Code Red)

| Failure | Behavior |
|---------|----------|
| cwd is neither repo nor scan folder | Print `ErrCnBridgeNoScope` to stderr, exit 1. |
| Probe worker panics | Recover, log to stderr with worker index, mark repo as failed in summary, continue. |
| All probes fail | Print summary with failure counts, exit 1. |
| Update worker fails | Log to stderr, record failed transition, continue with other workers. |
| User Ctrl+C mid-probe or mid-update | Cancel context, drain workers, print partial summary, exit 130. |
| TTY detection ambiguous | Default to non-interactive; require `--select` or `--yes` to act. |

## Acceptance Criteria

1. `gitmap cn` (no args) inside a git repo enters single-repo bridge.
2. `gitmap cn` (no args) inside a registered scan folder enters multi-repo bridge.
3. `gitmap cn` (no args) outside both prints clear error and exits 1.
4. Empty/stale `VersionProbe` cache triggers an auto probe-walk.
5. Probe-walk uses spec-103 ls-remote+verify semantics with depth=5.
6. Per-repo probes execute in parallel (worker pool); per-repo walks stay sequential.
7. Interactive TUI offers yes-all / one-by-one / selected / cancel.
8. `--yes` performs default-all update without TTY.
9. Updates execute in parallel with streaming progress.
10. Every probe and every update is persisted (`VersionProbe`, `RepoVersionHistory`).
11. Subsequent runs hit the cache when within `FindNextStaleAfter`.
12. `gitmap find-next` remains strictly read-only — auto-probe lives **only** in the `cn` bridge.
13. React UI exposes `/find-next` page and lists `find-next` in `/commands`.
14. Existing `cn vN`, `cn v++`, `cn --all`, `cn --csv`, `cn <repo> vN` paths are unchanged.

## Out of Scope

- Probing intermediate versions for **changelog** content (covered by
  hypothetical `gitmap changelog-fetch`).
- Auto-resolving merge conflicts during update (always uses the same
  flatten-by-default workflow as spec 59).
- Cross-host probing (assumes one remote per repo, matching spec 59).

## See Also

- [59 — clone-next](59-clone-next.md)
- [90 — ScanFolder & VersionProbe](90-scan-folder-and-version-probe.md)
- [103 — probe --depth N](103-probe-depth.md) (**blocking dependency**)
- [43 — interactive TUI](43-interactive-tui.md)
- `mem://features/find-next` — read-only `find-next` reader contract
- `mem://features/parallel-pull` — worker-pool prior art
