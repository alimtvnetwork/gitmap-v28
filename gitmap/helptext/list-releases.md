# gitmap list-releases

List release metadata from the current git repo or stored database.

Builds a unified list from three sources: repo metadata files
(`.gitmap/release/v*.json`), git tags, and the SQLite database. All
discovered releases are cached to the DB on every invocation.

## Alias

lr, releases

## Usage

    gitmap list-releases [--json] [--limit N] [--source repo|release|import|tag]
    gitmap releases --all-repos [--json] [--limit N]

## Flags

| Flag | Default | Description |
|------|---------|-------------|
| --json | false | Output as structured JSON |
| --limit \<N\> | 0 | Show only the top N releases (0 = all) |
| --source \<type\> | — | Filter by release source (repo, release, import, or tag) |
| --all-repos | false | (v3.20.0) Multi-repo batch view: every Release row in the DB joined with its Repo, ordered by CreatedAt. Bypasses the cwd-bound repo/tag scan. |

## Prerequisites

- Inside a git repo with `.gitmap/release/v*.json` files or semver tags, **or**
- Run `gitmap scan` or `gitmap release` to populate the database

## Examples

### Example 1: List releases from current repo

    gitmap list-releases

**Output:**

    Releases (5 found)
    ────────────────────────────────────────────────────────────────────────
      VERSION    TAG          BRANCH              DRAFT  LATEST  SOURCE   DATE
      2.33.0     v2.33.0      release/v2.33.0     no     yes     repo     2026-03-26
      2.31.0     v2.31.0      release/v2.31.0     no     no      repo     2026-03-20
      2.30.0     v2.30.0      release/v2.30.0     no     no      repo     2026-03-15
      2.28.0     v2.28.0      release/v2.28.0     no     no      tag      2026-03-01
      2.25.0     v2.25.0      release/v2.25.0     no     no      tag      2026-02-10
      5 releases found

### Example 2: Show only tag-discovered releases

    gitmap lr --source tag

**Output:**

    Releases (2 found)
    ────────────────────────────────────────────────────────────────────────
      VERSION    TAG          BRANCH              DRAFT  LATEST  SOURCE   DATE
      2.28.0     v2.28.0      release/v2.28.0     no     no      tag      2026-03-01
      2.25.0     v2.25.0      release/v2.25.0     no     no      tag      2026-02-10

### Example 3: Top 3 releases as JSON

    gitmap lr --limit 3 --json

**Output:**

    [
      {"version":"2.33.0","tag":"v2.33.0","branch":"release/v2.33.0","source":"repo","draft":false,"isLatest":true},
      {"version":"2.31.0","tag":"v2.31.0","branch":"release/v2.31.0","source":"repo","draft":false,"isLatest":false},
      {"version":"2.30.0","tag":"v2.30.0","branch":"release/v2.30.0","source":"repo","draft":false,"isLatest":false}
    ]

## Notes

- Git tags without a matching `.gitmap/release/` metadata file are included
  with `source=tag`, containing version, tag, inferred branch, and tag date.
- All discovered releases are automatically upserted into the SQLite
  `Releases` table on every invocation, keeping the DB in sync.
- The database is used as a fallback only when no repo files or tags are found.

## See Also

- [list-versions](list-versions.md) — List Git release tags
- [changelog](changelog.md) — View release notes
- [release](release.md) — Create a release
- [scan](scan.md) — Scan to import release data

## Scripting (JSON)

Discover this command from a script using the machine-readable help payload:

```bash
gitmap help --json --filter list-releases
```

The JSON schema is published at `spec/08-json-schemas/help-json.schema.json` (v5.43.0+).
