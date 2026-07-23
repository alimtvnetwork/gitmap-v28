# Spec 56 — Unified `.gitmap/` Directory

## Overview

Consolidate all repo-local output directories under a single `.gitmap/` folder
at the repository root. This replaces the current split between `.release/` and
`gitmap-output/`.

## Current State

| Purpose            | Current Path          | Constant                  |
|--------------------|-----------------------|---------------------------|
| Release metadata   | `.release/`           | `DefaultReleaseDir`       |
| Scan output        | `gitmap-output/`      | `DefaultOutputDir`        |
| Output folder name | `gitmap-output`       | `DefaultOutputFolder`     |
| Verbose logs       | `gitmap-output/`      | `DefaultVerboseLogDir`    |

The SQLite database stays binary-relative (`data/`) — **not affected** by this change.

## Target State

```
.gitmap/
├── release/          ← was .release/
│   ├── latest.json
│   └── v2.34.0.json
└── output/           ← was gitmap-output/
    ├── gitmap.csv
    ├── gitmap.json
    ├── gitmap.txt
    ├── folder-structure.md
    └── clone.ps1
```

## New Constants

```go
// Root directory — all repo-local gitmap-v28 data lives here.
const GitMapDir = ".gitmap"

// Subdirectories.
const (
    ReleaseDirName = "release"
    OutputDirName  = "output"
)

// Resolved default paths.
var DefaultReleaseDir = filepath.Join(GitMapDir, ReleaseDirName)   // .gitmap/release
const DefaultOutputDir  = ".gitmap/output"
const DefaultOutputFolder = "output"
const DefaultVerboseLogDir = ".gitmap/output"
```

## Migration Checklist

### Phase 1 — Constants (1 file)

Update `constants/constants.go`:

| Old Value            | New Value              |
|----------------------|------------------------|
| `.release`           | `.gitmap/release`      |
| `./gitmap-output`    | `.gitmap/output`       |
| `gitmap-output`      | `.gitmap/output`       |

### Phase 2 — Code References (~27 files)

All Go source files already use `constants.DefaultReleaseDir`,
`constants.DefaultOutputDir`, etc. — no hardcoded paths to fix in
business logic. Files to verify:

- `release/metadata.go` — uses `DefaultReleaseDir` ✓
- `release/workflowfinalize.go` — uses `DefaultReleaseDir` ✓
- `release/autocommit.go` — prefix check on `DefaultReleaseDir` ✓
- `release/workflow.go` — uses `DefaultReleaseDir` ✓
- `cmd/scanimport.go` — uses `DefaultReleaseDir` ✓
- `cmd/clearreleasejson.go` — uses `DefaultReleaseDir` ✓
- `cmd/listreleases.go` — uses `DefaultReleaseDir` via release pkg ✓
- `config/config_test.go` — hardcoded `"./gitmap-output"` ⚠️ update
- `constants/constants_terminal.go` — hardcoded `gitmap-output` in display strings ⚠️ update
- `constants/constants_messages.go` — hardcoded `gitmap-output/` in error msg ⚠️ update
- `constants/constants_cli.go` — hardcoded in help text ⚠️ update

### Phase 3 — Display Strings & Help Text (~4 files)

Update user-facing messages that reference old paths:

- `constants_terminal.go` — `StatusRepoCountFmt`, `ExecRepoCountFmt`,
  `TermCloneCmd1`, `TermCloneCmd3b`
- `constants_messages.go` — `MsgNoOutputDir`
- `constants_cli.go` — `HelpOutputPath`
- `helptext/*.md` — any references to `gitmap-output/` or `.release/`

### Phase 4 — Tests

- `config/config_test.go` — update expected `OutputDir` value
- `release/metadata_test.go` — uses `DefaultReleaseDir` var (auto-inherits) ✓
- `release/workflow_test.go` — uses `DefaultReleaseDir` var ✓

### Phase 5 — Specs, Docs & Memory

- Update all spec files referencing `.release/` or `gitmap-output/`
- Update `helptext/` markdown files
- Update docs site (`src/data/`, `src/pages/`)
- Update `.lovable/memory/` files
- Add changelog entry

## .gitignore Consideration

Projects using gitmap-v28 should add `.gitmap/output/` to their `.gitignore`.
The `.gitmap/release/` directory should remain tracked (release metadata is
committed). The `.gitmap/` root itself should NOT be gitignored.

## Automatic Migration

When any gitmap-v28 command runs and detects a legacy directory in the current
working directory, it automatically moves the contents to `.gitmap/`:

| Legacy Directory   | Target                | Trigger                        |
|--------------------|-----------------------|--------------------------------|
| `gitmap-output/`   | `.gitmap/output/`     | Any command that writes output |
| `.release/`        | `.gitmap/release/`    | Any command that reads/writes release metadata |
| `.deployed/`       | `.gitmap/deployed/`   | Any command that reads/writes deploy state |

### Migration Rules

1. **Detection**: Check if the legacy directory exists at the working directory root.
2. **Create parent**: Ensure `.gitmap/` exists (create if missing).
3. **Move**: Rename the legacy directory to its new location under `.gitmap/`.
4. **Merge if target exists**: If the target directory already exists, merge files
   from the legacy directory into the target (skip files that already exist in
   the target), then **remove the legacy directory entirely**. This ensures the
   old folder never persists after migration runs.
5. **Log**: Print a single-line message per migration:
   - Clean move: `Migrated <old>/ -> .gitmap/<new>/`
   - Merge: `Merged <old>/ into .gitmap/<new>/ (N files copied, M skipped) and removed legacy folder`
6. **No database changes**: The SQLite database remains binary-relative in
   `data/` and is completely unaffected.

### Implementation

A shared `localdirs.MigrateLegacyDirs()` function handles all legacy directory
migration. It is called in two places:

1. Early in the root command startup flow (skipped for `version` / `v` to keep stdout clean)
2. Again in the release workflow after returning to the original branch

The second call is required because checking out the original branch can restore
tracked legacy `.release/` files from older branches. Re-running migration at
that point guarantees `.release/` is merged into `.gitmap/release/` and removed
before auto-commit runs.

```go
// After returning to the original branch, run migration again.
err = returnToBranch(originalBranch)
if err != nil {
    return err
}
localdirs.MigrateLegacyDirs()
```


### Doctor Check

Legacy directory warnings have been removed from the doctor command since
migration now handles cleanup automatically. No manual intervention is needed.

### Scan Output Path

The `scan` command always writes output to `.gitmap/output/` relative to the
scanned directory root. The `resolveOutputDir` function enforces this by
joining `scanDir + .gitmap + output`, ignoring the config `outputDir` value
unless it is an absolute path. This means users never see output in a
top-level `output/` or `gitmap-output/` folder.

## Not In Scope

- Database location (stays binary-relative)
- Config file location (stays in `data/`)
- Profile storage (stays in `data/`)
