# Memory: features/scan-architecture
Updated: now

The 'scan' command is split into 'scan.go' (orchestration and folder opening) and 'scanoutput.go' (CSV, JSON, and terminal output generation) to comply with project limits.

## Output Directory

All generated output files (CSV, JSON, Markdown, PowerShell scripts, scan cache) are **always** written to `.gitmap/output/` relative to the scanned directory root. The `resolveOutputDir` function in `cmd/scan.go` enforces this: it joins `scanDir + constants.GitMapDir + constants.OutputDirName` unless an absolute `outputDir` is provided via config. Top-level `output/` or `gitmap-output/` folders are never used.

## Database Integration

The command integrates with the SQLite database via an `upsertToDB` hook that persists all discovered repository records after each scan completion. Record IDs are aligned with the database via `alignRecordsWithDB`.

## Project Detection

After scanning, the command detects project types (Go, React, Node, C#, C++) across repos, writes per-type JSON files, and upserts projects to the database.
