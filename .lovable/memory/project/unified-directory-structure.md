# Memory: project/unified-directory-structure
Updated: now

All gitmap-related artifacts are consolidated under `.gitmap/` at the repository root:

| Subdirectory | Purpose |
|---|---|
| `.gitmap/release/` | Release metadata JSON files |
| `.gitmap/output/` | Scan output (CSV, JSON, scripts, folder-structure.md, scan cache) |
| `.gitmap/deployed/` | Deployment logs |

## Scan Output

The `scan` command **always** writes to `.gitmap/output/` relative to the scanned directory. The `resolveOutputDir` function enforces this by joining `scanDir + .gitmap + output`, ignoring config `outputDir` unless it is an absolute path.

## Migration

A shared `localdirs.MigrateLegacyDirs()` function automatically moves legacy directories (`.release/`, `gitmap-output/`, `.deployed/`) into `.gitmap/` using merge-and-remove. Called at CLI startup and after returning to the original branch during release. Doctor no longer warns about legacy directories since migration handles cleanup.

## Excluded

SQLite database and config remain binary-relative in `data/`.
