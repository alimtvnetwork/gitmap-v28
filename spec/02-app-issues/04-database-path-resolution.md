# Issue 04: Database Written to Scan Output Directory Instead of Binary Location

## Date
2026-03-15

## Severity
High — data integrity and portability issue.

## Symptom
Running `gitmap-v27 scan` creates the SQLite database at `<CWD>/.gitmap/output/data/gitmap.db`
inside the scan output directory, instead of at `<binary-location>/data/gitmap.db`. This
means:
- The database location changes depending on which directory the user runs the command from.
- Multiple unrelated databases can be created across different working directories.
- Commands that read from the database (list, group, status, exec) may not find repos
  scanned from a different CWD.

## Root Cause
`store.Open(outputDir)` accepted a directory parameter and resolved the DB path as
`outputDir/data/gitmap.db`. All 13 callers across the codebase passed CWD-relative
constants like `constants.DefaultOutputDir` ("./.gitmap/output") or
`constants.DefaultOutputFolder` (".gitmap/output"). This made the DB path depend
on the working directory rather than the binary's physical installation location.

### Affected Callers (13 total)
| File | Function | Previous Path |
|------|----------|---------------|
| `cmd/scan.go` | `upsertToDB` | `<outputDir>/data/gitmap.db` |
| `cmd/scan.go` | `alignRecordsWithDB` | `<outputDir>/data/gitmap.db` |
| `cmd/list.go` | `openDB` | `.gitmap/output/data/gitmap.db` |
| `cmd/release.go` | `persistRelease` | `.gitmap/output/data/gitmap.db` |
| `cmd/audit.go` | `openAuditDB` | `<CWD>/.gitmap/output/data/gitmap.db` |
| `cmd/amendaudit.go` | `saveAmendToDB` | `./.gitmap/output/data/gitmap.db` |
| `cmd/seowritetemplate.go` | `openSEODatabase` | `./.gitmap/output/data/gitmap.db` |
| `cmd/projectrepos.go` | `runProjectRepos` | `<resolved CWD>/.gitmap/output/data/` |
| `cmd/scanprojects.go` | `upsertProjectsToDB` | `<outputDir>/data/gitmap.db` |
| `cmd/scanimport.go` | `importReleases` | `<outputDir>/data/gitmap.db` |
| `cmd/diffprofiles.go` | `loadProfileRepos` | `.gitmap/output/data/gitmap.db` |
| `cmd/profileutil.go` | `initProfileDB` | `.gitmap/output/data/gitmap.db` |
| `cmd/interactive.go` | `runInteractive` | `constants.DefaultDBPath` (undefined) |

## Solution
1. Created `store/location.go` with `BinaryDataDir()` — resolves the executable's
   physical location via `os.Executable()` + `filepath.EvalSymlinks()`, then appends
   `/data/` to get the stable database directory.

2. Added `store.OpenDefault()` — opens the database from the binary's data directory
   without requiring any path argument.

3. Added `store.OpenDefaultProfile(name)` — opens a named profile's database from
   the same binary-relative directory.

4. Updated all 13 callers to use `store.OpenDefault()` or `store.OpenDefaultProfile()`
   instead of `store.Open(outputDir)`.

5. Removed now-unused helper functions:
   - `resolveAuditOutputDir()` in `cmd/audit.go`
   - `resolveDefaultOutputDir()` in `cmd/projectrepos.go`

## Correct Behavior After Fix
- Database always lives at `<binary-dir>/data/gitmap.db`
- Scan output files (CSV, JSON, scripts) still write to the user-specified
  output directory (unchanged behavior)
- All commands see the same database regardless of CWD

## Learnings
- Database paths must be anchored to the binary's installation directory, not CWD.
- When a path constant is used by many callers, a single wrong default propagates
  silently across the entire codebase.
- `os.Executable()` with `filepath.EvalSymlinks()` is the reliable way to find
  binary location (handles symlinks and PATH resolution).
