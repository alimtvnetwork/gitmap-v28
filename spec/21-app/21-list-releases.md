# 21 — list-releases Command

## Purpose

`gitmap list-releases` (`lr`) displays release records, sorted from newest to
oldest. It builds a unified list from three sources: repo metadata files, git
tags, and the SQLite database. All discovered releases are cached to the DB
on every invocation.

## Command Signature

```
gitmap list-releases [flags]
gitmap lr [flags]
```

## Flags

| Flag       | Short | Default | Description                                        |
|------------|-------|---------|----------------------------------------------------|
| `--json`   |       | false   | Output as JSON array                               |
| `--limit`  |       | 0       | Show only the top N releases (0 = all)             |
| `--source` |       | (all)   | Filter by source: `release`, `import`, `repo`, or `tag` |

## Data Source (Resolution Order)

1. **Repo-local `.gitmap/release/` files** (preferred): read all `.gitmap/release/v*.json`
   files via `release.ListReleaseMetaFiles()`, convert each `ReleaseMeta` to
   a `ReleaseRecord` with `Source = "repo"`, sort by `CreatedAt DESC`, and
   mark the latest using `.gitmap/release/latest.json`.
2. **Git tag discovery**: scan all local semver tags via `git for-each-ref`
   with `creatordate`. Tags that have no matching `.gitmap/release/` metadata
   get a minimal `ReleaseRecord` with `Source = "tag"`, containing version,
   tag, inferred branch name, and tag creation date.
3. **Database fallback**: if neither repo files nor git tags produce any
   records, open the SQLite database via `store.Open()` and call
   `db.ListReleases()`.

### DB Caching

After building the unified list, all records are upserted into the SQLite
`Releases` table on every `lr` invocation. This ensures the DB stays in sync
with repo-local metadata and newly created git tags without requiring a
separate `gitmap scan`.

## Behavior

1. Call `release.ListReleaseMetaFiles()` to read `.gitmap/release/v*.json`.
2. If results are found, convert to `[]model.ReleaseRecord`, sort by
   `CreatedAt DESC`, and mark `IsLatest` from `latest.json`.
3. Scan git tags via `release.ListVersionTags()`. For each tag not already
   covered by a repo metadata file, create a minimal `ReleaseRecord` with
   `Source = "tag"`.
4. Merge repo records and tag-only records.
5. If the merged list is empty, fall back to `db.ListReleases()`.
6. Sort the unified list by `CreatedAt DESC`.
7. Cache all records to the DB via `db.UpsertRelease()`.
8. Apply `--source` filter if provided.
9. Apply `--limit N` if provided and N > 0.
10. Render output in terminal or JSON format.
11. If no releases are found, print `"No releases found."` and exit 0.

## Terminal Output

Table format with columns: Version, Tag, Branch, Draft, Latest, Source, Date.

```
Releases (5 found)
──────────────────────────────────────────────────────────
  VERSION    TAG        BRANCH       DRAFT  LATEST  SOURCE   DATE
  2.33.0     v2.33.0    release/v2   no     yes     repo     2026-03-26
  2.31.0     v2.31.0    release/v2   no     no      repo     2026-03-20
  2.30.0     v2.30.0    release/v2   yes    no      repo     2026-03-15
  2.28.0     v2.28.0    release/v2   no     no      tag      2026-03-01
  2.25.0     v2.25.0    release/v2   no     no      tag      2026-02-10
```

## JSON Output Example

```json
[
  {
    "version": "2.33.0",
    "tag": "v2.33.0",
    "branch": "release/v2.33.0",
    "sourceBranch": "main",
    "commitSha": "abc123",
    "changelog": "Added --limit flag",
    "draft": false,
    "preRelease": false,
    "isLatest": true,
    "source": "repo",
    "createdAt": "2026-03-26T10:00:00Z"
  },
  {
    "version": "2.28.0",
    "tag": "v2.28.0",
    "branch": "release/v2.28.0",
    "sourceBranch": "",
    "commitSha": "",
    "changelog": "",
    "draft": false,
    "preRelease": false,
    "isLatest": false,
    "source": "tag",
    "createdAt": "2026-03-01T12:00:00+00:00"
  }
]
```

## Error Handling

| Condition              | Message                                          | Exit |
|------------------------|--------------------------------------------------|------|
| No repo + no tags + no DB | `"No database found. Run gitmap scan first.\n"` | 1    |
| DB open/migrate error  | `"failed to load releases: %v\n"`                | 1    |
| No releases            | `"No releases found.\n"`                         | 0    |

## Implementation Files

| File                              | Responsibility                                  |
|-----------------------------------|--------------------------------------------------|
| `cmd/listreleases.go`             | Command handler, filtering, output               |
| `cmd/listreleasesload.go`         | Repo loading, tag discovery, DB caching          |
| `release/metadata.go`             | `ListReleaseMetaFiles()`, `ReadLatest()`         |
| `release/gitopstags.go`           | `ListVersionTags()`, `TagEntry`                  |
| `model/release.go`                | `ReleaseRecord`, `SourceRepo`, `SourceTag`       |
| `constants/constants_cli.go`      | `CmdListReleases`, `CmdListReleasesAlias`        |
| `constants/constants_messages.go` | Terminal output format strings                   |
| `store/release.go`                | `UpsertRelease()`, `ListReleases()` (DB fallback)|

## Integration Points

- `cmd/root.go`: register `list-releases` / `lr` in `dispatchMisc`.
- Reuse `release.ListReleaseMetaFiles()` — no new filesystem code needed.
- `release.ListVersionTags()` uses `git for-each-ref` for tag dates.
- DB caching runs on every invocation (upsert by tag, no duplicates).
- DB path is only resolved for fallback loading when no local data exists.

## Code Style

All functions ≤ 15 lines. Positive logic. Blank line before every return.
No magic strings. No switch statements. PascalCase for SQL column names.
