---
slug: make-public-url-or-folder
status: open
scope: cmd/visibility
captured: 2026-06-06
---

# Command: `gitmap make-public` / `make-private` MUST accept URL or folder

**Verbatim user instruction:**
> Also make sure to have `gitmap make-public URL or folder` both should work fine, clear??

## Scope
Applies to all four visibility commands:
- `make-public` / `make-private` (existing, single + bulk-vN)
- `make-all-public` / `make-all-private` (new, see plan 01)
- Their shorthands `MAPUB` / `MAPRI`

## Requirement
The first positional argument MUST resolve correctly whether the user passes:
1. A full provider URL — `https://github.com/owner` or `https://github.com/owner/repo`
2. A bare host/owner path — `github.com/owner` or `github.com/owner/repo`
3. A local folder path (existing repo on disk) — resolve its `origin` remote → owner[/repo]
4. `.` (current working directory) — same as (3)

For the org-wide `make-all-*` family, forms (1) and (2) without a `/repo` segment
resolve to "owner/org root" so the repo list is fetched server-side. Form (3)/(4)
applied to `make-all-*` resolves the folder's origin owner, then operates on the
**whole owner**, not just that one repo (so `gitmap make-all-public .` inside any
repo of `alimtvnetwork` operates on `alimtvnetwork`).

## When it applies
Every invocation of the four commands. No flag needed — auto-detect by inspecting
the argument (`os.Stat` for folder, `url.Parse` for URL, else treat as bare host).

## Related
- `.lovable/plans/pending/01-bulk-visibility-mapub-mapri.md` step 6 + step 27
- `gitmap/cmd/visibilityresolve.go` (existing resolver to reuse/extend)
