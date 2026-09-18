---
name: cn-find-next-bridge
description: Planned `gitmap cn` no-args bridge that detects scope (single repo vs scan-folder root), hydrates VersionProbe cache, auto-probes (spec 103 ls-remote+verify, depth=5, parallel) on miss, and offers an interactive TUI to update all/one-by-one/selected/yes-all. `find-next` itself stays strictly read-only.
type: feature
---

# `cn`-no-args → `find-next` bridge (planned, v3.55.0)

## Status

**Planned.** Spec: `02-spec/01-app/107-cn-find-next-bridge.md`.
Plan: `.ai-memory/memory/plans/03-cn-find-next-bridge-plan.md`.
Blocked on spec 103 (`probe --depth N`) shipping first.

## Why

Today `gitmap cn` with no args errors out, and `gitmap find-next` is
read-only — users have to chain `probe --all` then `find-next` then
`cn vN` per repo. The bridge collapses that into one ergonomic command
with parallelism and an interactive picker.

## User-locked decisions

- **Probe walk semantics**: spec 103 model — `ls-remote
  --sort=-v:refname` then shallow-clone-verify newest→oldest with
  depth=5. NOT a strict `v+1, v+2, …` forward walk.
- **Auto-probe scope**: only the new `cn` bridge auto-probes on
  empty/stale cache. Plain `gitmap find-next` stays strictly
  read-only — its CI invariant ("never touches the network") is
  preserved.

## Workflow

```
gitmap cn (no args)
  ├── detect cwd: git repo | scan-folder root | neither
  ├── hydrate VersionProbe (FindNextStaleAfter = 24h)
  ├── empty or stale → parallel probe walk (depth=5, workers=NumCPU)
  ├── render TUI summary table
  ├── user keys: a (yes-all) / n (one-by-one) / space+enter (selected)
  │              / r (refresh) / q (cancel)
  ├── parallel `cn v++` per selected repo (workers=NumCPU)
  └── persist VersionProbe rows + RepoVersionHistory rows
      (with new TriggeredByProbeId trace column)
```

## Concurrency model

- **Outer (per repo)**: parallel worker pool. Configurable via
  `--probe-workers` and `--update-workers`. Default `runtime.NumCPU()`.
- **Inner (per repo's tag walk)**: sequential — clones share the same
  remote, parallelizing them just hammers the server (per spec 103
  §Concurrency).

## Non-interactive contract

- TTY + no `--yes` → interactive Bubble Tea picker.
- Non-TTY (CI, pipes) + no `--yes` → print summary and exit 0.
- `--yes` → default-all, no prompt, parallel update. Works in TTY and
  non-TTY identically.

## What stays unchanged

- `gitmap cn vN`, `gitmap cn v++`, `gitmap cn v+1` — spec 59 paths.
- `gitmap cn --all` / `--csv <path>` — existing batch path
  (`clonenextbatch.go`).
- `gitmap cn <repo> vN` — cross-dir form (`clonenextcrossdir.go`).
- `gitmap find-next` (`fn`) — read-only summary of latest available
  probes.
- `gitmap probe` (no `--depth`) — single-tag probe behavior.

## React UI requirement

`find-next` is currently **not listed** in the docs site at all. The
plan adds:

- `src/pages/FindNext.tsx` at `/find-next` (workflow diagram, flags,
  examples).
- New entry in `src/pages/Commands.tsx`.
- "No-args bridge" section in `src/pages/CloneNext.tsx`.

## Persistence trace

A new optional column `RepoVersionHistory.TriggeredByProbeId` lets
operators reconstruct which probe row caused each upgrade. Idempotent
migration via `addColumnIfNotExists`.

## See also

- `mem://features/find-next` — current read-only `find-next` reader
- `mem://features/version-probe` — single-tag probe
- `mem://features/clone-next-flatten` — flatten-by-default clone
- `mem://features/parallel-pull` — worker-pool prior art
- `mem://features/interactive-tui` — Bubble Tea infra to reuse
