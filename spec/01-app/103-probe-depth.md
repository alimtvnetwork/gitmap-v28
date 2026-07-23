# 103 — `gitmap probe --depth N` (Multi-Version Walk)

> Status: **planned** — targets v3.36.0
> Depends on: spec 90 (`VersionProbe` table, v3.7.0), `mem://features/version-probe` (v3.8.0)

## Goal

Today's `gitmap probe` reports **only the single newest tag** for each
repo. If a repo has jumped from `v1.4.0` → `v1.7.2` since the last
probe, the user sees `v1.7.2` and never learns that `v1.5.0` /
`v1.6.0` / `v1.7.0` / `v1.7.1` shipped in between (which may matter
for changelog review, security audits, or staged upgrades).

`--depth N` walks **up to N newer-than-current tags**, verifies each
via shallow clone, and persists one `VersionProbe` row per verified
version.

## CLI Surface

| Form                        | Behavior                                       |
|-----------------------------|------------------------------------------------|
| `gitmap probe`              | **Unchanged.** Single newest tag, as today.    |
| `gitmap probe --depth 5`    | Walk up to 5 newer tags, shallow-verify each.  |
| `gitmap probe --depth 5 --all` | Same, applied to every repo in DB.          |
| `gitmap probe <path> --depth 5` | Same, single repo by path.                 |

`--depth` is clamped to `[1, ProbeMaxDepth=10]`. Default when omitted is
`ProbeDefaultDepth=1` (preserves existing behavior — **fully backwards
compatible**).

User chose `--depth` (not a new `probe-deep` command and not
`--all-newer`). Rationale: zero new command surface, every existing
script keeps working, the depth knob is discoverable next to `--all`.

## Algorithm

```
1. ls-remote --tags --sort=-v:refname <url>
   → produce ordered tag list, newest first
2. Find current_tag (read from RepoVersionHistory or local `git describe`).
   If none: treat as "all returned tags are candidates".
3. Take the slice of tags STRICTLY NEWER than current_tag.
4. Trim to first N entries (where N = --depth).
5. For each candidate (newest → oldest):
       git clone --depth 1 --filter=blob:none --no-checkout --branch <tag>
                 <url> <tmpdir>
       if clone OK → record VersionProbe row { IsAvailable=1, Method="shallow-clone-tag" }
       if clone fails → record VersionProbe row { IsAvailable=0, Method="shallow-clone-tag", Err=... }
6. Cleanup tmpdirs (defer os.RemoveAll inside loop).
```

### Why shallow-clone every candidate (not ls-remote-only)

Chosen by user. Tags returned by `ls-remote` may be:

- Yanked / force-pushed-away (server returns the ref but cloning fails).
- Pointing to commits that no longer exist (rare but real on
  rewritten histories).
- Lightweight tag pointing into a detached / pruned object graph.

A successful `--branch <tag> --depth 1` clone is the only way to
**prove** the tag is fetchable. The `--filter=blob:none --no-checkout`
combo keeps each verification cheap (refs DB only, no working tree, no
blobs). Typical clone takes 100-400ms on a fast connection.

### "Strictly newer" comparison

Use the existing `probe.parseSemverInt` (`MAJOR*1e6 + MINOR*1e3 +
PATCH`) and compare against current. Pre-release suffixes (`-rc1`,
`-beta`) collapse to numeric prefix; we **include** pre-releases in
the candidate list but flag them with `IsPreRelease=1` (new column —
see schema below) so callers can filter.

## Storage — One VersionProbe Row Per Verified Version

Chosen by user. Schema unchanged at the table level, but each `probe
--depth N` invocation now inserts up to N rows into `VersionProbe`,
each with the same `RepoId` + a distinct `NextVersionTag` /
`NextVersionNum`.

### Optional column addition (non-breaking)

```sql
ALTER TABLE VersionProbe
  ADD COLUMN IsPreRelease INTEGER NOT NULL DEFAULT 0;
```

Added via `addColumnIfNotExists` per the idempotent migration pattern
(`mem://tech/database-architecture`).

### Find-next compatibility

`gitmap find-next` (spec 9, v3.9.0) currently joins on
`MAX(ProbedAt)` per repo. With multi-row probes, all rows from a
single `probe --depth` invocation share the same `ProbedAt`
timestamp (set once at probe start, not per-clone), so the existing
`MAX(ProbedAt)` query naturally returns **all** rows from the most
recent walk. No `find-next` query change required.

A new flag `find-next --include-intermediate` toggles whether to show
just the highest version per repo (default, current behavior) or every
version from the latest walk. Default stays as today.

## Concurrency

Within a single repo's walk: **sequential** (clones share the same
remote, parallelizing them just hammers the server).

Across repos in `probe --all --depth N`: reuses the existing
sequential outer loop. Phase 2.7 (out of scope here) can later
parallelize the outer loop with the worker-pool pattern from spec 100.

## Output

```
✓ owner/repo (5 verified)
    v1.4.1  → v1.5.0  shallow-clone-tag  238ms
    v1.5.0  → v1.6.0  shallow-clone-tag  191ms
    v1.6.0  → v1.7.0  shallow-clone-tag  204ms
    v1.7.0  → v1.7.1  shallow-clone-tag  186ms
    v1.7.1  → v1.7.2  shallow-clone-tag  198ms
✗ other/repo  → v2.1.0 (clone failed: auth)
· third/repo  → no new version (depth=5)
```

JSON output (`--json` flag inherited from existing probe) emits one
top-level array element per repo, with a nested `versions: []` array
holding each verified row.

## Constants (per `mem://style/code-constraints`)

`gitmap/constants/constants_probe.go`:

```go
const (
    FlagProbeDepth     = "--depth"
    ProbeDefaultDepth  = 1
    ProbeMaxDepth      = 10
    ProbeMethodTag     = "shallow-clone-tag"
)
```

`gitmap/constants/constants_messages.go`:

```go
const (
    MsgProbeDepthClamped   = "[probe] --depth clamped to %d (max %d)"
    MsgProbeWalkSummary    = "✓ %s (%d verified)"
    MsgProbeWalkLine       = "    %s → %s  %s  %s"
    MsgProbeWalkNothing    = "· %s → no new version (depth=%d)"
)
```

No magic strings. No raw `"--depth"` / `"shallow-clone-tag"` literals
outside the constants file.

## Error Handling (Code Red)

- Per-candidate clone failure: log to `os.Stderr` with format
  `[probe-walk] %s@%s: %v`, insert `VersionProbe { IsAvailable=0,
  Err=<msg> }` row, continue with next candidate.
- ls-remote failure: aborts the walk for that repo (no candidates to
  iterate); same exit-code semantics as today's probe.
- Tmpdir cleanup failure: log warning, do not fail the walk
  (matches `mem://tech/security-hardening` deferred-cleanup pattern).
- Use `errors.Is(ctx.Err(), context.Canceled)` to honor Ctrl+C
  mid-walk.

## Testing

- Unit: `TestProbeDepth_ClampMaxDepth` (--depth 99 → clamped to 10
  with warning).
- Unit: `TestProbeDepth_StrictlyNewer` (current=v1.5.0, server returns
  v1.5.0/v1.4.0/v1.6.0 → only v1.6.0 is a candidate).
- Unit: `TestProbeDepth_PreReleaseFlagged` (server returns
  v2.0.0-rc1 → row with IsPreRelease=1).
- Integration (live remote): probe a real public repo with N=3,
  assert 3 VersionProbe rows inserted with distinct NextVersionTag.
- Integration (broken tag): use a test remote that returns a tag
  pointing to a missing commit, assert `IsAvailable=0` row recorded
  and walk continues.

## Out of Scope

- Parallel inter-repo walks (covered by future spec).
- Auto-pull / auto-checkout to the verified version (probe stays
  read-only).
- Reading remote `CHANGELOG.md` per verified version (separate
  `gitmap changelog-fetch` candidate).
- `--depth all` / unbounded walks (rejected by user; capped at 10).
