# Repository Split Database Commands

## 1. Storage Location & Schema Extensions

Repository split databases store file sequence maps, full-text indexes, and search caches under:
- **Location:** `<BinaryDataDir>/repo_search/<slug>-<id>.db`

### Table: `RepoScanLog`

Tracks repository synchronization events, indexing runs, and errors.

```sql
CREATE TABLE IF NOT EXISTS RepoScanLog (
    Id INTEGER PRIMARY KEY AUTOINCREMENT,
    RepoId INTEGER NOT NULL,
    RepoSlug TEXT NOT NULL,
    Action TEXT NOT NULL,
    Status TEXT NOT NULL,
    ErrorMessage TEXT NULL,
    Details TEXT NULL,
    Notes TEXT NULL,
    Comments TEXT NULL,
    CreatedAt TEXT DEFAULT CURRENT_TIMESTAMP
);
```

## 2. CLI Subcommands: `gitmap repo db`

| Subcommand | Description |
|---|---|
| `status` (default) | Shows repo split DB location, path, size, `RepoFile` rows, `SearchCache` rows, and `RepoScanLog` count |
| `log` | Shows recent scan and indexing event entries from `RepoScanLog` |
| `errorlogs` / `error-logs` | Queries entries where `Status == 'failure'` or `ErrorMessage IS NOT NULL` |
| `clear` | Clears search cache and file index rows with confirmation prompt (`-y` to skip) |
| `reset` | Drops repository index tables and re-initializes clean schema |
| `optimize` | Runs `VACUUM` and `PRAGMA optimize` on the repository split DB |
| `help` | Prints help documentation |
