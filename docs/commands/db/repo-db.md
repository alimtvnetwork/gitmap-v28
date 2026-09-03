# `gitmap db repo-db list`

Inspect per-repository split database metrics and tracking health.

## Usage

```bash
gitmap db repo-db list
```

## Columns

| Column | Description |
|---|---|
| `SLUG` | Repository slug identifier |
| `DB FILE` | Split database file basename in `repo_search/` |
| `SIZE` | File size on disk |
| `FILES` | Number of rows in `RepoFile` table |
| `CACHES` | Number of rows in `SearchCache` table |
| `STATUS` | Tracking state in master database (`TRACKED` / `UNTRACKED`) |
