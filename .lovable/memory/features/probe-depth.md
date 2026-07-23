---
name: Probe Depth
description: gitmap probe --depth N walks up to N newer tags, shallow-clones each to verify, inserts one VersionProbe row per verified version. Default 1 (backwards compatible), max 10.
type: feature
---

# Feature: `gitmap probe --depth N` (multi-version walk)

**Spec:** `spec/01-app/103-probe-depth.md`
**Status:** planned for v3.36.0
**Depends on:** VersionProbe table (spec 90), existing probe (mem://features/version-probe v3.8.0)

## Behavior summary

- **Surface:** new `--depth N` flag on existing `gitmap probe`. No new command. Default `ProbeDefaultDepth=1` (preserves existing single-tag behavior). Clamped to `[1, ProbeMaxDepth=10]` with warning if exceeded.
- **Algorithm:** ls-remote → filter strictly newer than current tag → trim to first N → for each candidate, `git clone --depth 1 --filter=blob:none --no-checkout --branch <tag>` into tmpdir to **verify** the tag actually fetches. Yanked / force-pushed / broken tags get an `IsAvailable=0` row.
- **Storage:** insert one `VersionProbe` row per verified version (all sharing the same `ProbedAt` timestamp from the walk start). Schema gains optional `IsPreRelease INTEGER DEFAULT 0` column via idempotent `addColumnIfNotExists`.
- **Pre-releases:** included in walk, flagged via `IsPreRelease=1` so callers can filter. Numeric prefix used for ordering.
- **Find-next compat:** existing `MAX(ProbedAt)` join naturally returns all rows from the latest walk. New `find-next --include-intermediate` flag toggles between newest-only (default) and all-versions display. No SQL change to current default behavior.
- **Concurrency:** sequential per-repo (same remote → no parallel benefit). Outer cross-repo loop stays sequential like today; parallel-walk is a separate future spec.
- **JSON output:** existing `--json` flag emits nested `versions: []` array per repo.

## Constants (no magic strings)

`constants_probe.go`: `FlagProbeDepth="--depth"`, `ProbeDefaultDepth=1`, `ProbeMaxDepth=10`, `ProbeMethodTag="shallow-clone-tag"`.

`constants_messages.go`: `MsgProbeDepthClamped`, `MsgProbeWalkSummary`, `MsgProbeWalkLine`, `MsgProbeWalkNothing`.

## Why these decisions

- **`--depth` flag (not new command)** — chosen by user. Zero breakage; every existing script keeps working; the knob lives next to `--all`.
- **Shallow-clone every candidate** — chosen by user. Only way to prove a tag is fetchable (catches yanked / force-pushed / broken tags that ls-remote happily reports). `--filter=blob:none --no-checkout` keeps each verification ~200ms.
- **One VersionProbe row per verified version** — chosen by user. Richer history; future find-next variants can show full upgrade path. Schema impact is one optional column.
- **Default depth = 5** — chosen by user. Matches "up to five more versions" phrasing. Clamped at 10 for safety.

## Cross-references

- Original probe: `mem://features/version-probe`
- Find-next consumer: `mem://features/find-next`
- Worker-pool pattern (for future parallel walks): spec 100
